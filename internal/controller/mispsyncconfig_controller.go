/*
MISP-Operator - A Kubernetes operator for simplified deployments of MISP at scale.
Copyright (C) 2026 Pascal Iske

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU Affero General Public License as published
by the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU Affero General Public License for more details.

You should have received a copy of the GNU Affero General Public License
along with this program.  If not, see <https://www.gnu.org/licenses/>.
*/

package controller

import (
	"context"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/events"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/manager"

	apiclient "github.com/pascaliske/misp-operator/internal/client"

	mispv1alpha1 "github.com/pascaliske/misp-operator/api/v1alpha1"
)

// MispSyncConfigReconciler reconciles a MispSyncConfig object
type MispSyncConfigReconciler struct {
	client.Client
	Scheme   *runtime.Scheme
	Recorder events.EventRecorder
}

func NewMispSyncConfigReconciler(mgr manager.Manager) *MispSyncConfigReconciler {
	return &MispSyncConfigReconciler{
		Client:   mgr.GetClient(),
		Scheme:   mgr.GetScheme(),
		Recorder: mgr.GetEventRecorder("mispsyncconfig-controller"),
	}
}

// sorted alphabetically to prevent duplicates
// +kubebuilder:rbac:groups=misp.k8s.pascaliske.dev,resources=mispsyncconfigs,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=misp.k8s.pascaliske.dev,resources=mispsyncconfigs/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=misp.k8s.pascaliske.dev,resources=mispsyncconfigs/finalizers,verbs=update

func (r *MispSyncConfigReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	logger.Info("Reconciliation started")
	defer func() {
		logger.Info("Reconciliation finished")
	}()

	// try to fetch the current resource, skip reconcile if not found
	var mispSyncConfig mispv1alpha1.MispSyncConfig
	if err := r.Get(ctx, req.NamespacedName, &mispSyncConfig); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// object is being deleted -> handle external dependencies and finalizer
	if !mispSyncConfig.DeletionTimestamp.IsZero() {
		return r.handleFinalizer(ctx, &mispSyncConfig)
	}

	// TODO: update status conditions

	// add finalizer if not present
	if !controllerutil.ContainsFinalizer(&mispSyncConfig, mispSyncConfigFinalizer) {
		logger.Info("Added missing finalizer to this object")
		return r.setupFinalizer(ctx, &mispSyncConfig)
	}

	// fetch local instance api client
	localInstanceApi, err := r.getInstanceApiClient(ctx, mispSyncConfig.Spec.InstanceRef.Namespace, mispSyncConfig.Spec.InstanceRef.Name)
	if err != nil {
		return ctrl.Result{}, err
	}

	// fetch remote instance api client
	remoteInstanceApi, err := r.getInstanceApiClient(ctx, mispSyncConfig.Spec.RemoteInstanceRef.Namespace, mispSyncConfig.Spec.RemoteInstanceRef.Name)
	if err != nil {
		return ctrl.Result{}, err
	}

	// TODO: check is ready for local & remote instances

	// step 1: local instance - get local org for sync config
	localOrganisation, err := localInstanceApi.GetOrganisationByName(mispSyncConfig.Spec.Organisation.Name)
	if err != nil {
		logger.Error(err, "Local organisation not found")
		return ctrl.Result{}, err
	}

	// step 2: remote instance - ensure org exists on remote instance with same uuid
	remoteOrganisation, err := remoteInstanceApi.EnsureOrganisation(localOrganisation)
	if err != nil {
		logger.Error(err, "Remote organisation not found")
		return ctrl.Result{}, err
	}

	// step 3: remote instance - ensure sync user exists in remote instances org
	syncUser, err := remoteInstanceApi.EnsureUser(&apiclient.User{OrgId: remoteOrganisation.Id, RoleId: "5", Email: "misp-sync@localhost"})
	if err != nil {
		logger.Error(err, "Remote sync user not found")
		return ctrl.Result{}, err
	}

	// step 4: remote instance - ensure auth key for sync user exists in remote instance
	syncAuthKey, err := remoteInstanceApi.EnsureAuthKey(&apiclient.AuthKey{UserId: syncUser.Id, Comment: mispSyncConfig.Name})
	if err != nil {
		logger.Error(err, "Remote auth key for sync user not found")
		return ctrl.Result{}, err
	}

	// step 5: local instance - ensure sync server for remote instance exists in local instance
	// TODO: use actual base url - remoteInstanceApi.GetBaseUrl().String()
	server, err := localInstanceApi.EnsureServer(&apiclient.Server{Name: mispSyncConfig.Name, Url: "http://misp-instance-b:8080", AuthKey: syncAuthKey.Raw, RemoteOrgId: localOrganisation.Id})
	if err != nil {
		logger.Error(err, "Local sync server not found")
		return ctrl.Result{}, err
	}

	// step 6: local instance - test server connection
	err = localInstanceApi.TestConnection(server)
	if err != nil {
		logger.Error(err, "Testing sync server failed")
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *MispSyncConfigReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&mispv1alpha1.MispSyncConfig{}).
		Named("mispsyncconfig").
		Complete(r)
}
