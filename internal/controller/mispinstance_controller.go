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

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/tools/events"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/manager"

	mispv1alpha1 "github.com/pascaliske/misp-operator/api/v1alpha1"
)

const applyFieldManagerKeyInstance string = "mispinstance-controller"

// MispInstanceReconciler reconciles a MispInstance object
type MispInstanceReconciler struct {
	client.Client
	Scheme   *runtime.Scheme
	Recorder events.EventRecorder
}

func NewMispInstanceReconciler(mgr manager.Manager) *MispInstanceReconciler {
	return &MispInstanceReconciler{
		Client:   mgr.GetClient(),
		Scheme:   mgr.GetScheme(),
		Recorder: mgr.GetEventRecorder("mispinstance-controller"),
	}
}

// sorted alphabetically to prevent duplicates
// +kubebuilder:rbac:groups="",resources=events,verbs=create;patch
// +kubebuilder:rbac:groups="",resources=persistentvolumeclaims,verbs=get;list;watch;create;update;patch
// +kubebuilder:rbac:groups="",resources=secrets,verbs=get;list;watch;create;update;patch
// +kubebuilder:rbac:groups="",resources=serviceaccounts,verbs=get;list;watch;create;update;patch
// +kubebuilder:rbac:groups="",resources=services,verbs=get;list;watch;create;update;patch
// +kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch
// +kubebuilder:rbac:groups=events.k8s.io,resources=events,verbs=create;patch
// +kubebuilder:rbac:groups=misp.k8s.pascaliske.dev,resources=mispinstances,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=misp.k8s.pascaliske.dev,resources=mispinstances/finalizers,verbs=update
// +kubebuilder:rbac:groups=misp.k8s.pascaliske.dev,resources=mispinstances/status,verbs=get;update;patch

func (r *MispInstanceReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	logger.Info("Reconciliation started")
	defer func() {
		logger.Info("Reconciliation finished")
	}()

	// try to fetch the current resource, skip reconcile if not found
	var mispInstance mispv1alpha1.MispInstance
	if err := r.Get(ctx, req.NamespacedName, &mispInstance); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// object is being deleted -> handle external dependencies and finalizer
	if !mispInstance.DeletionTimestamp.IsZero() {
		return r.handleFinalizer(ctx, &mispInstance)
	}

	// initialize status conditions if not yet present
	if len(mispInstance.Status.Conditions) == 0 {
		meta.SetStatusCondition(&mispInstance.Status.Conditions, metav1.Condition{
			Type:    mispv1alpha1.ConditionTypeReady,
			Status:  metav1.ConditionUnknown,
			Reason:  mispv1alpha1.PhasePending,
			Message: "Starting reconciliation.",
		})

		if err := r.Status().Update(ctx, &mispInstance); err != nil {
			logger.Error(err, "Failed to update MispInstance status")
			return ctrl.Result{}, err
		}
	}

	// add finalizer if not present
	if !controllerutil.ContainsFinalizer(&mispInstance, mispInstanceFinalizer) {
		logger.Info("Added missing finalizer to this object")
		return r.setupFinalizer(ctx, &mispInstance)
	}

	// skip reconcile if resource is suspended
	if mispInstance.Spec.Suspend {
		logger.Info("Reconciliation is suspended for this object")

		// update resource status to reflect suspended state
		if err := r.updateStatus(ctx, &mispInstance); err != nil {
			return ctrl.Result{}, err
		}

		return ctrl.Result{}, nil
	}

	// reconcile instance
	if err := r.reconcileInstance(ctx, &mispInstance); err != nil {
		return ctrl.Result{}, err
	}

	// reconcile restart
	if err := r.reconcileRestart(ctx, &mispInstance); err != nil {
		return ctrl.Result{}, err
	}

	// update resource status
	if err := r.updateStatus(ctx, &mispInstance); err != nil {
		return ctrl.Result{}, err
	}

	// finish reconcile cycle
	return ctrl.Result{}, nil
}

func (r *MispInstanceReconciler) reconcileInstance(ctx context.Context, mispInstance *mispv1alpha1.MispInstance) error {
	// reconcile service account object
	if err := r.reconcileServiceAccount(ctx, mispInstance); err != nil {
		return err
	}

	// no admin secret referenced -> manage admin secret through controller
	if mispInstance.Spec.Admin == nil || mispInstance.Spec.Admin.CredentialsSecretRef == nil {
		if err := r.reconcileAdminSecret(ctx, mispInstance); err != nil {
			return err
		}
	}

	// reconcile persistent volume claim object
	if err := r.reconcilePersistentVolumeClaim(ctx, mispInstance); err != nil {
		return err
	}

	// reconcile deployment objects
	if err := r.reconcileDeployments(ctx, mispInstance); err != nil {
		return err
	}

	// reconcile service objects
	if err := r.reconcileServices(ctx, mispInstance); err != nil {
		return err
	}

	return nil
}

func (r *MispInstanceReconciler) reconcileRestart(ctx context.Context, mispInstance *mispv1alpha1.MispInstance) error {
	logger := log.FromContext(ctx)

	// skip if restart annotation is empty
	restartAt, ok := mispInstance.Annotations[RestartAnnotation]
	if !ok || restartAt == "" {
		return nil
	}

	// try to fetch the instance deployment
	var deployment appsv1.Deployment
	if err := r.Get(ctx, types.NamespacedName{Name: mispInstance.Name, Namespace: mispInstance.Namespace}, &deployment); err != nil {
		return client.IgnoreNotFound(err)
	}

	// ensure annotations are set
	if deployment.Spec.Template.Annotations == nil {
		deployment.Spec.Template.Annotations = map[string]string{}
	}

	// skip if restart value is already injected
	currentRestart := deployment.Spec.Template.Annotations[InternalRestartAnnotation]
	if currentRestart == restartAt {
		return nil
	}

	// inject restart value into main deployment
	logger.Info("Triggered restart of instance")
	deployment.Spec.Template.Annotations[InternalRestartAnnotation] = restartAt

	// update main deployment
	return r.Update(ctx, &deployment)
}

// SetupWithManager sets up the controller with the Manager.
func (r *MispInstanceReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&mispv1alpha1.MispInstance{}).
		Owns(&appsv1.Deployment{}).
		Owns(&corev1.PersistentVolumeClaim{}).
		Owns(&corev1.Service{}).
		Owns(&corev1.ServiceAccount{}).
		Owns(&corev1.Secret{}).
		Named("mispinstance").
		Complete(r)
}
