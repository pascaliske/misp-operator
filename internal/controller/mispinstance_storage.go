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
	"strconv"

	mispv1alpha1 "github.com/pascaliske/misp-operator/api/v1alpha1"
	"github.com/pascaliske/misp-operator/internal/utils"

	"sigs.k8s.io/controller-runtime/pkg/client"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	corev1apply "k8s.io/client-go/applyconfigurations/core/v1"
	metav1apply "k8s.io/client-go/applyconfigurations/meta/v1"
)

func (r *MispInstanceReconciler) createPersistentVolumeClaim(mispInstance *mispv1alpha1.MispInstance) *corev1apply.PersistentVolumeClaimApplyConfiguration {
	pvcSpec := corev1apply.PersistentVolumeClaimSpec()

	// inject pvc template if set
	if mispInstance.Spec.Storage != nil && mispInstance.Spec.Storage.PersistentVolumeClaimTemplate != nil {
		pvcSpec = pvcSpec.
			WithAccessModes(mispInstance.Spec.Storage.PersistentVolumeClaimTemplate.AccessModes...).
			WithVolumeMode(*mispInstance.Spec.Storage.PersistentVolumeClaimTemplate.VolumeMode).
			WithStorageClassName(*mispInstance.Spec.Storage.PersistentVolumeClaimTemplate.StorageClassName).
			WithResources(
				corev1apply.
					VolumeResourceRequirements().
					WithRequests(corev1.ResourceList{
						corev1.ResourceStorage: resource.MustParse(strconv.Itoa(mispInstance.Spec.Storage.PersistentVolumeClaimTemplate.Resources.Size())),
					}),
			)
	} else {
		pvcSpec = pvcSpec.
			WithAccessModes(corev1.ReadWriteOnce).
			WithVolumeMode(corev1.PersistentVolumeFilesystem)
	}

	// inject storage class name if set
	if mispInstance.Spec.Storage != nil && mispInstance.Spec.Storage.StorageClass != "" {
		pvcSpec = pvcSpec.WithStorageClassName(mispInstance.Spec.Storage.StorageClass)
	}

	// inject storage size if set
	if mispInstance.Spec.Storage != nil && mispInstance.Spec.Storage.Size != "" {
		pvcSpec = pvcSpec.WithResources(
			corev1apply.
				VolumeResourceRequirements().
				WithRequests(corev1.ResourceList{
					corev1.ResourceStorage: resource.MustParse(mispInstance.Spec.Storage.Size),
				}),
		)
	} else {
		pvcSpec = pvcSpec.WithResources(
			corev1apply.
				VolumeResourceRequirements().
				WithRequests(corev1.ResourceList{
					corev1.ResourceStorage: resource.MustParse("10Gi"),
				}),
		)
	}

	return corev1apply.
		PersistentVolumeClaim(mispInstance.GetNameWithSuffix("storage"), mispInstance.Namespace).
		WithLabels(utils.BuildAppLabels(mispInstance.Name, utils.AppLabelComponentStorage)).
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
			pvcSpec,
		)

}

func (r *MispInstanceReconciler) reconcilePersistentVolumeClaim(ctx context.Context, mispInstance *mispv1alpha1.MispInstance) error {
	pvc := r.createPersistentVolumeClaim(mispInstance)

	return r.Apply(ctx, pvc, &client.ApplyOptions{
		FieldManager: applyFieldManagerKeyInstance,
		Force:        new(true),
	})
}
