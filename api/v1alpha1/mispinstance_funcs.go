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

package v1alpha1

import (
	"fmt"
	"strings"
)

// default container images
const mispInstanceCoreImageDefault = "ghcr.io/misp/misp-docker/misp-core:v2.5.44"
const mispInstanceNginxImageDefault = "ghcr.io/misp/misp-docker/misp-nginx:v2.5.44"
const mispInstanceModulesImageDefault = "ghcr.io/misp/misp-docker/misp-modules:v3.0.9"

// Returns the custom misp-core image if set or the default misp-core image as fallback
func (mispInstance *MispInstance) GetCoreImage() string {
	if mispInstance.Spec.Image != "" {
		return mispInstance.Spec.Image
	}

	return mispInstanceCoreImageDefault
}

// Returns the custom misp-nginx image if set or the default misp-nginx image as fallback
func (mispInstance *MispInstance) GetNginxImage() string {
	if mispInstance.Spec.Nginx != nil && mispInstance.Spec.Nginx.Image != "" {
		return mispInstance.Spec.Nginx.Image
	}

	return mispInstanceNginxImageDefault
}

// Returns the custom misp-modules image if set or the default misp-modules image as fallback
func (mispInstance *MispInstance) GetModulesImage() string {
	if mispInstance.Spec.Modules != nil && mispInstance.Spec.Modules.Image != "" {
		return mispInstance.Spec.Modules.Image
	}

	return mispInstanceModulesImageDefault
}

func (mispInstance *MispInstance) GetName() string {
	return strings.ToLower(mispInstance.Name)
}

func (mispInstance *MispInstance) GetNameWithSuffix(suffix string) string {
	return strings.ToLower(fmt.Sprintf("%s-%s", mispInstance.Name, suffix))
}
