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
	"fmt"

	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	apiclient "github.com/pascaliske/misp-operator/internal/client"

	mispv1alpha1 "github.com/pascaliske/misp-operator/api/v1alpha1"
)

func (r *MispSyncConfigReconciler) getInstanceApiClient(ctx context.Context, namespace string, name string) (*apiclient.ApiClient, error) {
	var instance mispv1alpha1.MispInstance
	var adminSecret corev1.Secret

	// fetch instance
	if err := r.Get(ctx, client.ObjectKey{Namespace: namespace, Name: name}, &instance); err != nil {
		return nil, err
	}

	// fetch instance admin credentials
	if err := r.getInstanceAdminSecret(ctx, instance, &adminSecret); err != nil {
		return nil, err
	}

	// prepare connection details
	baseUrl := instance.Spec.BaseUrl
	apiKey := adminSecret.Data["apiKey"]

	fmt.Printf("==> New API client for %s with api key %s\n", baseUrl, string(apiKey))

	// return api client
	return apiclient.NewApiClient(baseUrl, string(apiKey))
}

func (r *MispSyncConfigReconciler) getInstanceAdminSecret(ctx context.Context, instance mispv1alpha1.MispInstance, adminSecret *corev1.Secret) error {
	// determine secret name
	secretName := instance.GetNameWithSuffix("admin")
	if instance.Spec.Admin != nil && instance.Spec.Admin.CredentialsSecretRef != nil {
		secretName = instance.Spec.Admin.CredentialsSecretRef.Name
	}

	// try to fetch secret
	secretObjectKey := client.ObjectKey{Namespace: instance.Namespace, Name: secretName}
	if err := r.Get(ctx, secretObjectKey, adminSecret); err != nil {
		return err
	}

	return nil
}
