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

package utils

import (
	corev1 "k8s.io/api/core/v1"
	corev1apply "k8s.io/client-go/applyconfigurations/core/v1"
)

func EnvVarToApplyConfiguration(env corev1.EnvVar) *corev1apply.EnvVarApplyConfiguration {
	ac := corev1apply.EnvVar().WithName(env.Name)

	if env.Value != "" {
		ac = ac.WithValue(env.Value)
	}

	if env.ValueFrom != nil {
		vf := corev1apply.EnvVarSource()

		if env.ValueFrom.FieldRef != nil {
			vf = vf.WithFieldRef(
				corev1apply.ObjectFieldSelector().
					WithAPIVersion(env.ValueFrom.FieldRef.APIVersion).
					WithFieldPath(env.ValueFrom.FieldRef.FieldPath),
			)
		}

		if env.ValueFrom.ResourceFieldRef != nil {
			vf = vf.WithResourceFieldRef(
				corev1apply.ResourceFieldSelector().
					WithContainerName(env.ValueFrom.ResourceFieldRef.ContainerName).
					WithResource(env.ValueFrom.ResourceFieldRef.Resource),
			)
		}

		if env.ValueFrom.ConfigMapKeyRef != nil {
			cmRef := corev1apply.
				ConfigMapKeySelector().
				WithName(env.ValueFrom.ConfigMapKeyRef.Name).
				WithKey(env.ValueFrom.ConfigMapKeyRef.Key)

			if env.ValueFrom.ConfigMapKeyRef.Optional != nil {
				cmRef = cmRef.WithOptional(*env.ValueFrom.ConfigMapKeyRef.Optional)
			}

			vf = vf.WithConfigMapKeyRef(cmRef)
		}

		if env.ValueFrom.SecretKeyRef != nil {
			skRef := corev1apply.
				SecretKeySelector().
				WithName(env.ValueFrom.SecretKeyRef.Name).
				WithKey(env.ValueFrom.SecretKeyRef.Key)

			if env.ValueFrom.SecretKeyRef.Optional != nil {
				skRef = skRef.WithOptional(*env.ValueFrom.SecretKeyRef.Optional)
			}

			vf = vf.WithSecretKeyRef(skRef)
		}

		ac = ac.WithValueFrom(vf)
	}

	return ac
}
