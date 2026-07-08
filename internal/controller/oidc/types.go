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

package oidc

type AuthMethod string

const (
	AuthMethodClientSecretBasic AuthMethod = "ClientSecretBasic"
	AuthMethodClientSecretPost  AuthMethod = "ClientSecretPost"
	AuthMethodClientSecretJwt   AuthMethod = "ClientSecretJwt"
	AuthMethodPrivateKeyJwt     AuthMethod = "PrivateKeyJwt"
)

func (method AuthMethod) String() string {
	switch method {
	case AuthMethodClientSecretBasic:
		return "client_secret_basic"
	case AuthMethodClientSecretPost:
		return "client_secret_post"
	case AuthMethodClientSecretJwt:
		return "client_secret_jwt"
	case AuthMethodPrivateKeyJwt:
		return "private_key_jwt"
	default:
		return ""
	}
}
