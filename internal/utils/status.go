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
	appsv1 "k8s.io/api/apps/v1"
)

func IsRolloutPending(deploy appsv1.Deployment) bool {
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
