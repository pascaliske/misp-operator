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
	"strings"
)

type AppLabelComponent string

const (
	AppLabelComponentMisp    AppLabelComponent = "misp"
	AppLabelComponentNginx   AppLabelComponent = "nginx"
	AppLabelComponentModules AppLabelComponent = "modules"
	AppLabelComponentStorage AppLabelComponent = "storage"
	AppLabelComponentRbac    AppLabelComponent = "rbac"
)

func BuildAppLabels(instance string, component AppLabelComponent) map[string]string {
	return map[string]string{
		"app.kubernetes.io/name":       strings.ToLower("misp"),
		"app.kubernetes.io/instance":   strings.ToLower(instance),
		"app.kubernetes.io/component":  strings.ToLower(string(component)),
		"app.kubernetes.io/managed-by": strings.ToLower("misp-operator"),
		"app.kubernetes.io/part-of":    strings.ToLower("misp"),
	}
}

func BuildSelectorLabels(instance string, component AppLabelComponent) map[string]string {
	return map[string]string{
		"app.kubernetes.io/name":      strings.ToLower("misp"),
		"app.kubernetes.io/instance":  strings.ToLower(instance),
		"app.kubernetes.io/component": strings.ToLower(string(component)),
	}
}
