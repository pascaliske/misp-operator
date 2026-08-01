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
	"github.com/pascaliske/misp-operator/internal/oidc"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

type AdminOrganisation struct {
	// Customize the initial organisations name
	// +kubebuilder:validation:Optional
	Name string `json:"name,omitempty"`

	// Pre-fill the initial organisations UUID
	// +kubebuilder:validation:Optional
	Uuid string `json:"uuid,omitempty"`
}

type Admin struct {
	// Optionally customize the initial admin organisation
	// +kubebuilder:validation:Optional
	Organisation *AdminOrganisation `json:"organisation,omitempty"`

	// Provide a secret in the same namespace with custom admin credentials
	// +kubebuilder:validation:Optional
	CredentialsSecretRef *corev1.LocalObjectReference `json:"credentialsSecretRef,omitempty"`
}

type Storage struct {
	// Customize the storage size of the PVC
	// +kubebuilder:validation:Type=string
	// +kubebuilder:validation:Optional
	Size string `json:"size,omitempty"`

	// Explicitly configure a storage class name for the PVC
	// +kubebuilder:validation:Type=string
	// +kubebuilder:validation:Optional
	StorageClass string `json:"storageClass,omitempty"`

	// Optionally provide your own PVC template - see here https://kubernetes.io/docs/reference/kubernetes-api/core/persistent-volume-claim-v1/#PersistentVolumeClaimSpec
	// +kubebuilder:validation:Optional
	PersistentVolumeClaimTemplate *corev1.PersistentVolumeClaimSpec `json:"pvcTemplate,omitempty"`
}

type Database struct {
	// Host address of the database
	// +required
	// +kubebuilder:validation:Type=string
	Host string `json:"host,omitempty"`

	// Port the database is listening on - defaults to "3306"
	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:validation:Maximum=65535
	// +kubebuilder:default=3306
	Port int32 `json:"port,omitempty"`

	// Name of the database for the instance - defaults to "misp"
	// +kubebuilder:validation:Optional
	// +kubebuilder:default=misp
	Name string `json:"name,omitempty"`

	// Reference a Kubernetes secret in the same namespace containing the database credentials
	// +required
	CredentialsSecretRef corev1.LocalObjectReference `json:"credentialsSecretRef"`
}

type Cache struct {
	// Host address of the cache
	// +required
	// +kubebuilder:validation:Type=string
	Host string `json:"host,omitempty"`

	// Port the cache is listening on - defaults to "6379"
	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:validation:Maximum=65535
	// +kubebuilder:default=6379
	Port int32 `json:"port,omitempty"`

	// Reference a Kubernetes secret in the same namespace containing the password
	// +kubebuilder:validation:Optional
	PasswordSecretRef *corev1.LocalObjectReference `json:"passwordSecretRef,omitempty"`

	// Explicitly disable the use of a password - this is not recommended for production setups
	// +kubebuilder:validation:Type=boolean
	// +kubebuilder:validation:Optional
	// +kubebuilder:default=false
	EnableEmptyPassword bool `json:"enableEmptyPassword,omitempty"`
}

type NginxForwardedHeaders struct {
	// Enable the usage of forwarded for headers
	// +kubebuilder:validation:Type=boolean
	// +kubebuilder:default=false
	Enabled bool `json:"enabled,omitempty"`

	// Configure the list of trusted proxies
	// +kubebuilder:default={}
	TrustedProxies []string `json:"trustedProxies,omitempty"`
}

type NginxSecurityHeaders struct {
	// Customize the value of the X-Frame-Options header
	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:Enum=DENY;SAMEORIGIN
	// +kubebuilder:default=SAMEORIGIN
	FrameOptions string `json:"frameOptions,omitempty"`

	// Customize the value of the Content-Security-Policy header
	// +kubebuilder:validation:Optional
	ContentSecurityPolicy string `json:"contentSecurityPolicy,omitempty"`

	// Customize the value of the Strict-Transport-Security header
	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:Minimum=0
	HstsMaxAge int32 `json:"hstsMaxAge,omitempty"`
}

type NginxFastCGI struct {
	// Customize the FastCGI read timeout - defaults to "300s"
	// +kubebuilder:validation:Optional
	// +kubebuilder:default="300s"
	ReadTimeout string `json:"readTimeout,omitempty"`

	// Customize the FastCGI send timeout - defaults to "300s"
	// +kubebuilder:validation:Optional
	// +kubebuilder:default="300s"
	SendTimeout string `json:"sendTimeout,omitempty"`

	// Customize the FastCGI connect timeout - defaults to "300s"
	// +kubebuilder:validation:Optional
	// +kubebuilder:default="300s"
	ConnectTimeout string `json:"connectTimeout,omitempty"`
}

type Nginx struct {
	// Override the nginx image - defaults to ghcr.io/misp/misp-docker/misp-nginx
	// +kubebuilder:validation:Type=string
	// +kubebuilder:validation:MinLength=1
	Image string `json:"image"`

	// Customize the max client body size of nginx
	// +kubebuilder:validation:Optional
	// +kubebuilder:default="50M"
	ClientMaxBodySize string `json:"clientMaxBodySize,omitempty"`

	// Configure forwarded headers
	// +kubebuilder:validation:Optional
	ForwardedHeaders *NginxForwardedHeaders `json:"forwardedHeaders,omitempty"`

	// Configure security headers
	// +kubebuilder:validation:Optional
	SecurityHeaders *NginxSecurityHeaders `json:"securityHeaders,omitempty"`

	// Customize FastCGI settings
	// +kubebuilder:validation:Optional
	FastCGI *NginxFastCGI `json:"fastcgi,omitempty"`
}

type OidcSettings struct {
	// Specify the client authentication method to use - defaults to "ClientSecretBasic"
	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:Enum=ClientSecretBasic;ClientSecretPost;ClientSecretJwt;PrivateKeyJwt
	// +kubebuilder:default=ClientSecretBasic
	Method oidc.AuthMethod `json:"method,omitempty"`

	// Request additional scopes from the Id
	// +kubebuilder:validation:Optional
	// +kubebuilder:default={profile,email}
	Scopes []string `json:"scopes,omitempty"`

	// Specify the default organisation for authenticated users
	// +kubebuilder:validation:Optional
	DefaultOrg string `json:"defaultOrg,omitempty"`

	// Specify the token claim which contains the user roles - defaults to "roles"
	// +kubebuilder:validation:Optional
	// +kubebuilder:default=roles
	RolesProperty string `json:"rolesProperty,omitempty"`

	// Specify a JSON mapping from OIDC roles to MISP role IDs
	// +kubebuilder:validation:Optional
	// +kubebuilder:default={}
	RolesMapping map[string]int `json:"rolesMapping,omitempty"`

	// Enable request of offline access tokens - defaults to "true"
	// +kubebuilder:validation:Optional
	// +kubebuilder:default=true
	OfflineAccess bool `json:"offlineAccess,omitempty"`

	// Allow both OIDC authentication and password authentication - defaults to "true"
	// +kubebuilder:validation:Optional
	// +kubebuilder:default=true
	MixedAuth bool `json:"mixedAuth,omitempty"`

	// Allow linking of existing users by email to OIDC - defaults to "false"
	// +kubebuilder:validation:Optional
	// +kubebuilder:default=false
	AllowEmailLinking bool `json:"allowEmailLinking,omitempty"`

	// Require `email_verified` token claim to be true - defaults to "true"
	// +kubebuilder:validation:Optional
	// +kubebuilder:default=true
	RequireEmailVerified bool `json:"requireEmailVerified,omitempty"`

	// Disable the use of the request object approach in authorization requests - defaults to "false"
	// +kubebuilder:validation:Optional
	// +kubebuilder:default=false
	DisableRequestObject bool `json:"disableRequestObject,omitempty"`

	// Disable the use of pushed authorization request - defaults to "false"
	// +kubebuilder:validation:Optional
	// +kubebuilder:default=false
	DisablePushedAuthorizationRequest bool `json:"disablePushedAuthorizationRequest,omitempty"`
}

type Oidc struct {
	// Enable OIDC auth mode
	// +kubebuilder:validation:Type=boolean
	// +kubebuilder:default=false
	Enabled bool `json:"enabled,omitempty"`

	// Provide OIDC client credentials
	CredentialsSecretRef corev1.LocalObjectReference `json:"credentialsSecretRef"`

	// Override default OIDC settings
	// +kubebuilder:validation:Optional
	Settings *OidcSettings `json:"settings,omitempty"`
}

type Smtp struct {
	// Provide SMTP host address
	// +kubebuilder:validation:Type=string
	Host string `json:"host,omitempty"`

	// Provide SMTP port number
	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:validation:Maximum=65535
	// +kubebuilder:default=25
	Port int32 `json:"port,omitempty"`
}

// +kubebuilder:validation:AtLeastOneOf=smtp
type Email struct {
	// Configure SMTP settings for email notifications
	// +kubebuilder:validation:Optional
	Smtp *Smtp `json:"smtp,omitempty"`
}

type Modules struct {
	// Enable optional modules component
	// +kubebuilder:validation:Type=boolean
	// +kubebuilder:default=true
	Enabled bool `json:"enabled,omitempty"`

	// Override the modules image - defaults to ghcr.io/misp/misp-docker/misp-modules
	// +kubebuilder:validation:Type=string
	// +kubebuilder:validation:MinLength=1
	Image string `json:"image"`
}

type MispInstanceSpec struct {
	// Suspend reconciliation of this resource through the operator
	// +kubebuilder:validation:Type=boolean
	// +kubebuilder:validation:Optional
	// +kubebuilder:default=false
	Suspend bool `json:"suspend,omitempty"`

	// Override the core image of the instance - defaults to ghcr.io/misp/misp-docker/misp-core
	// +kubebuilder:validation:Type=string
	// +kubebuilder:validation:Optional
	Image string `json:"image,omitempty"`

	// Override the image pull policy for the core image - defaults to "IfNotPresent"
	// +kubebuilder:validation:Optional
	// +kubebuilder:default=IfNotPresent
	ImagePullPolicy corev1.PullPolicy `json:"imagePullPolicy,omitempty"`

	// Provide image pull secrets for the core image
	// +kubebuilder:validation:Optional
	ImagePullSecrets []corev1.LocalObjectReference `json:"imagePullSecrets,omitempty"`

	// Provide the front-facing url of the instance
	// +required
	// +kubebuilder:validation:Type=string
	// +kubebuilder:validation:MinLength=1
	BaseUrl string `json:"baseUrl"`

	// Pre-fill the instance UUID - defaults to an auto generated one
	// +kubebuilder:validation:Type=string
	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:MinLength=1
	Uuid string `json:"uuid,omitempty"`

	// Override the time zone of the containers
	// +kubebuilder:validation:Type=string
	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:default=Etc/UTC
	TimeZone string `json:"timeZone"`

	// Optional admin account configuration
	// +kubebuilder:validation:Optional
	Admin *Admin `json:"admin,omitempty"`

	// Optional storage configuration - defaults to a 10Gi PVC with default storage class name
	// +kubebuilder:validation:Optional
	Storage *Storage `json:"storage,omitempty"`

	// Provide your database connection details
	// +required
	Database Database `json:"database"`

	// Provide your k/v based cache connection details - e.g. valkey or redis
	// +required
	Cache Cache `json:"cache"`

	// Optional nginx configuration
	// +kubebuilder:validation:Optional
	Nginx *Nginx `json:"nginx,omitempty"`

	// Configure optional OIDC auth settings
	// +kubebuilder:validation:Optional
	Oidc *Oidc `json:"oidc,omitempty"`

	// Configure optional email notification settings
	// +kubebuilder:validation:Optional
	Email *Email `json:"email"`

	// Optional modules component
	// +kubebuilder:validation:Optional
	Modules *Modules `json:"modules"`

	// Provide extra environment variables for the core container
	// +kubebuilder:validation:Optional
	// +kubebuilder:default={}
	ExtraEnvs []corev1.EnvVar `json:"extraEnvs,omitempty"`
}

type MispInstanceStatus struct {
	// +kubebuilder:validation:Type=string
	Image string `json:"image,omitempty"`

	// +kubebuilder:validation:Enum=Pending;Running;Failed;Suspended
	Phase string `json:"phase,omitempty"`

	// LastReconcileTime is the timestamp of the last reconciliation
	LastReconcileTime *metav1.Time `json:"lastReconcileTime,omitempty"`

	// The status of each condition is one of True, False, or Unknown.
	// +patchMergeKey=type
	// +patchStrategy=merge
	// +listType=map
	// +listMapKey=type
	// +optional
	Conditions []metav1.Condition `json:"conditions,omitempty" patchStrategy:"merge" patchMergeKey:"type"`

	// Message provides additional information about the current state
	Message string `json:"message,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Image",type=string,JSONPath=`.status.image`
// +kubebuilder:printcolumn:name="Base URL",type=string,JSONPath=`.spec.baseUrl`
// +kubebuilder:printcolumn:name="Phase",type=string,JSONPath=`.status.phase`
// +kubebuilder:resource:shortName=misp

// MispInstance is the Schema for the mispinstances API
type MispInstance struct {
	metav1.TypeMeta `json:",inline"`

	// metadata is a standard object metadata
	// +optional
	metav1.ObjectMeta `json:"metadata,omitzero"`

	// spec defines the desired state of MispInstance
	// +required
	Spec MispInstanceSpec `json:"spec"`

	// status defines the observed state of MispInstance
	// +optional
	Status MispInstanceStatus `json:"status,omitzero"`
}

// +kubebuilder:object:root=true

// MispInstanceList contains a list of MispInstance
type MispInstanceList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitzero"`
	Items           []MispInstance `json:"items"`
}

func init() {
	SchemeBuilder.Register(func(s *runtime.Scheme) error {
		s.AddKnownTypes(SchemeGroupVersion, &MispInstance{}, &MispInstanceList{})
		return nil
	})
}
