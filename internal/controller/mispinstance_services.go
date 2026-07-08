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

	mispv1alpha1 "github.com/pascaliske/misp-operator/api/v1alpha1"
	"github.com/pascaliske/misp-operator/internal/utils"

	"k8s.io/apimachinery/pkg/util/intstr"
	corev1apply "k8s.io/client-go/applyconfigurations/core/v1"
	metav1apply "k8s.io/client-go/applyconfigurations/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func (r *MispInstanceReconciler) createInstanceService(mispInstance *mispv1alpha1.MispInstance) *corev1apply.ServiceApplyConfiguration {
	return corev1apply.
		Service(mispInstance.Name, mispInstance.Namespace).
		WithLabels(utils.BuildAppLabels(mispInstance.Name, utils.AppLabelComponentMisp)).
		WithOwnerReferences(
			metav1apply.
				OwnerReference().
				WithAPIVersion(mispInstance.APIVersion).
				WithKind(mispInstance.Kind).
				WithName(mispInstance.Name).
				WithUID(mispInstance.UID).
				WithController(true).
				WithBlockOwnerDeletion(true),
		).
		WithSpec(
			corev1apply.
				ServiceSpec().
				WithSelector(utils.BuildSelectorLabels(mispInstance.Name, utils.AppLabelComponentMisp)).
				WithPorts(
					corev1apply.
						ServicePort().
						WithName("http").
						WithProtocol("TCP").
						WithPort(8080).
						WithTargetPort(intstr.FromInt32(8080)),
				),
		)
}

func (r *MispInstanceReconciler) createModulesService(mispInstance *mispv1alpha1.MispInstance) *corev1apply.ServiceApplyConfiguration {
	return corev1apply.
		Service(mispInstance.GetNameWithSuffix("modules"), mispInstance.Namespace).
		WithLabels(utils.BuildAppLabels(mispInstance.Name, utils.AppLabelComponentModules)).
		WithOwnerReferences(
			metav1apply.
				OwnerReference().
				WithAPIVersion(mispInstance.APIVersion).
				WithKind(mispInstance.Kind).
				WithName(mispInstance.Name).
				WithUID(mispInstance.UID).
				WithController(true).
				WithBlockOwnerDeletion(true),
		).
		WithSpec(
			corev1apply.
				ServiceSpec().
				WithSelector(utils.BuildSelectorLabels(mispInstance.Name, utils.AppLabelComponentModules)).
				WithPorts(
					corev1apply.
						ServicePort().
						WithName("http").
						WithProtocol("TCP").
						WithPort(6666).
						WithTargetPort(intstr.FromInt32(6666)),
				),
		)
}

func (r *MispInstanceReconciler) reconcileServices(ctx context.Context, mispInstance *mispv1alpha1.MispInstance) error {
	options := &client.ApplyOptions{
		FieldManager: applyFieldManagerKey,
		Force:        new(true),
	}

	// reconcile misp modules service if enabled
	if mispInstance.Spec.Modules != nil && mispInstance.Spec.Modules.Enabled {
		if err := r.Apply(ctx, r.createModulesService(mispInstance), options); err != nil {
			return err
		}
	}

	// reconcile misp instance service
	return r.Apply(ctx, r.createInstanceService(mispInstance), options)
}
