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

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	corev1apply "k8s.io/client-go/applyconfigurations/core/v1"
	metav1apply "k8s.io/client-go/applyconfigurations/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// create a generic secret tied to the instance
func (r *MispInstanceReconciler) createSecret(mispInstance *mispv1alpha1.MispInstance, name string, data map[string]string) *corev1apply.SecretApplyConfiguration {
	return corev1apply.
		Secret(name, mispInstance.Namespace).
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
		WithStringData(data)
}

func (r *MispInstanceReconciler) reconcileAdminSecret(ctx context.Context, mispInstance *mispv1alpha1.MispInstance) error {
	secretName := mispInstance.GetNameWithSuffix("admin")
	secretObjectKey := client.ObjectKey{Namespace: mispInstance.Namespace, Name: secretName}

	// skip re-creation if secret already exists
	if err := r.Get(ctx, secretObjectKey, &corev1.Secret{}); !errors.IsNotFound(err) {
		return nil
	}

	// generate default secret values
	defaultEmail := "admin@admin.test"
	defaultPassword := utils.RandomString(40)
	defaultApiKey := utils.RandomString(40)

	// build admin secret
	secret := r.createSecret(mispInstance, secretName, map[string]string{
		"email":    defaultEmail,
		"password": defaultPassword,
		"apiKey":   defaultApiKey,
	})

	return r.Apply(ctx, secret, &client.ApplyOptions{
		FieldManager: applyFieldManagerKey,
		Force:        new(true),
	})
}
