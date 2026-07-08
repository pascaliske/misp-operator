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
	"reflect"

	mispv1alpha1 "github.com/pascaliske/misp-operator/api/v1alpha1"

	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

func isRolloutPending(deploy appsv1.Deployment) bool {
	// fetch desired replicas
	replicas := int32(1)
	if deploy.Spec.Replicas != nil {
		replicas = *deploy.Spec.Replicas
	}

	return deploy.Status.ObservedGeneration < deploy.Generation ||
		deploy.Status.Replicas != replicas ||
		deploy.Status.UpdatedReplicas != replicas ||
		deploy.Status.AvailableReplicas != replicas
}

func (r *MispInstanceReconciler) updateStatus(ctx context.Context, mispInstance *mispv1alpha1.MispInstance) error {
	logger := log.FromContext(ctx)

	// save previous status for comparison
	previousStatus := mispInstance.Status

	// reflect actual image in status
	if image := mispInstance.GetCoreImage(); mispInstance.Status.Image != image {
		mispInstance.Status.Image = image
	}

	// fetch current deployment
	var deploy appsv1.Deployment
	if err := r.Get(ctx, client.ObjectKey{Namespace: mispInstance.Namespace, Name: mispInstance.Name}, &deploy); err != nil {
		return client.IgnoreNotFound(err)
	}

	// reflect suspended/pending/running state in phase
	if mispInstance.Spec.Suspend {
		mispInstance.Status.Phase = mispv1alpha1.PhaseSuspended
		mispInstance.Status.Message = "Reconciliation is suspended."
	} else if isRolloutPending(deploy) {
		mispInstance.Status.Phase = mispv1alpha1.PhasePending
		mispInstance.Status.Message = "MispInstance is pending."
	} else {
		mispInstance.Status.Phase = mispv1alpha1.PhaseRunning
		mispInstance.Status.Message = "MispInstance is running."
	}

	// build ready condition
	now := metav1.Now()
	condition := metav1.Condition{
		Type:               mispv1alpha1.ConditionTypeReady,
		Status:             metav1.ConditionFalse,
		Reason:             mispInstance.Status.Phase,
		Message:            mispInstance.Status.Message,
		LastTransitionTime: now,
	}

	// reflect phase in ready condition
	if mispInstance.Status.Phase == mispv1alpha1.PhaseRunning {
		condition.Status = metav1.ConditionTrue
	}

	// update status subresource
	if !reflect.DeepEqual(previousStatus, mispInstance.Status) {
		logger.Info("Updated resource status")

		// update status condition
		meta.SetStatusCondition(&mispInstance.Status.Conditions, condition)

		// update last reconciliation time
		mispInstance.Status.LastReconcileTime = &now

		return r.Status().Update(ctx, mispInstance)
	}

	return nil
}
