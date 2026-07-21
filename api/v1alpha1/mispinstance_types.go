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
	"github.com/pascaliske/misp-operator/internal/controller/oidc"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

type AdminOrganisation struct {
	// +kubebuilder:validation:Optional
	Name string `json:"name,omitempty"`

	// +kubebuilder:validation:Optional
	Uuid string `json:"uuid,omitempty"`
}

type Admin struct {
	// +kubebuilder:validation:Optional
	Organisation *AdminOrganisation `json:"organisation,omitempty"`

	// +kubebuilder:validation:Optional
	CredentialsSecretRef *corev1.LocalObjectReference `json:"credentialsSecretRef,omitempty"`
}

type Storage struct {
	// +kubebuilder:validation:Type=string
	// +kubebuilder:validation:Optional
	Size string `json:"size,omitempty"`

	// +kubebuilder:validation:Type=string
	// +kubebuilder:validation:Optional
	StorageClass string `json:"storageClass,omitempty"`

	// +kubebuilder:validation:Optional
	PersistentVolumeClaimTemplate *corev1.PersistentVolumeClaimSpec `json:"pvcTemplate,omitempty"`
}

type Database struct {
	// +kubebuilder:validation:Type=string
	Host string `json:"host,omitempty"`

	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:validation:Maximum=65535
	// +kubebuilder:default=3306
	Port int32 `json:"port,omitempty"`

	// +kubebuilder:validation:Optional
	// +kubebuilder:default=misp
	Name string `json:"name,omitempty"`

	CredentialsSecretRef corev1.LocalObjectReference `json:"credentialsSecretRef"`
}

type Cache struct {
	// +kubebuilder:validation:Type=string
	Host string `json:"host,omitempty"`

	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:validation:Maximum=65535
	// +kubebuilder:default=6379
	Port int32 `json:"port,omitempty"`

	// +kubebuilder:validation:Optional
	PasswordSecretRef *corev1.LocalObjectReference `json:"passwordSecretRef,omitempty"`

	// +kubebuilder:validation:Type=boolean
	// +kubebuilder:validation:Optional
	// +kubebuilder:default=false
	EnableEmptyPassword bool `json:"enableEmptyPassword,omitempty"`
}

type NginxForwardedHeaders struct {
	// +kubebuilder:validation:Type=boolean
	// +kubebuilder:default=false
	Enabled bool `json:"enabled,omitempty"`

	// +kubebuilder:default={}
	TrustedProxies []string `json:"trustedProxies,omitempty"`
}

type NginxSecurityHeaders struct {
	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:Enum=DENY;SAMEORIGIN
	// +kubebuilder:default=SAMEORIGIN
	FrameOptions string `json:"frameOptions,omitempty"`

	// +kubebuilder:validation:Optional
	ContentSecurityPolicy string `json:"contentSecurityPolicy,omitempty"`

	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:Minimum=0
	HstsMaxAge int32 `json:"hstsMaxAge,omitempty"`
}

type NginxFastCGI struct {
	// +kubebuilder:validation:Optional
	// +kubebuilder:default="300s"
	ReadTimeout string `json:"readTimeout,omitempty"`

	// +kubebuilder:validation:Optional
	// +kubebuilder:default="300s"
	SendTimeout string `json:"sendTimeout,omitempty"`

	// +kubebuilder:validation:Optional
	// +kubebuilder:default="300s"
	ConnectTimeout string `json:"connectTimeout,omitempty"`
}

type Nginx struct {
	// +kubebuilder:validation:Type=string
	// +kubebuilder:validation:MinLength=1
	Image string `json:"image"`

	// +kubebuilder:validation:Optional
	// +kubebuilder:default="50M"
	ClientMaxBodySize string `json:"clientMaxBodySize,omitempty"`

	// +kubebuilder:validation:Optional
	ForwardedHeaders *NginxForwardedHeaders `json:"forwardedHeaders,omitempty"`

	// +kubebuilder:validation:Optional
	SecurityHeaders *NginxSecurityHeaders `json:"securityHeaders,omitempty"`

	// +kubebuilder:validation:Optional
	FastCGI *NginxFastCGI `json:"fastcgi,omitempty"`
}

type OidcSettings struct {
	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:Enum=ClientSecretBasic;ClientSecretPost;ClientSecretJwt;PrivateKeyJwt
	// +kubebuilder:default=ClientSecretBasic
	Method oidc.AuthMethod `json:"method,omitempty"`

	// +kubebuilder:validation:Optional
	// +kubebuilder:default={profile,email}
	Scopes []string `json:"scopes,omitempty"`

	// +kubebuilder:validation:Optional
	DefaultOrg string `json:"defaultOrg,omitempty"`

	// +kubebuilder:validation:Optional
	// +kubebuilder:default=roles
	RolesProperty string `json:"rolesProperty,omitempty"`

	// +kubebuilder:validation:Optional
	// +kubebuilder:default={}
	RolesMapping map[string]int `json:"rolesMapping,omitempty"`

	// +kubebuilder:validation:Optional
	// +kubebuilder:default=true
	OfflineAccess bool `json:"offlineAccess,omitempty"`

	// +kubebuilder:validation:Optional
	// +kubebuilder:default=true
	MixedAuth bool `json:"mixedAuth,omitempty"`

	// +kubebuilder:validation:Optional
	// +kubebuilder:default=false
	AllowEmailLinking bool `json:"allowEmailLinking,omitempty"`

	// +kubebuilder:validation:Optional
	// +kubebuilder:default=true
	RequireEmailVerified bool `json:"requireEmailVerified,omitempty"`

	// +kubebuilder:validation:Optional
	// +kubebuilder:default=false
	DisableRequestObject bool `json:"disableRequestObject,omitempty"`

	// +kubebuilder:validation:Optional
	// +kubebuilder:default=false
	DisablePushedAuthorizationRequest bool `json:"disablePushedAuthorizationRequest,omitempty"`
}

type Oidc struct {
	// +kubebuilder:validation:Type=boolean
	// +kubebuilder:default=false
	Enabled bool `json:"enabled,omitempty"`

	CredentialsSecretRef corev1.LocalObjectReference `json:"credentialsSecretRef"`

	// +kubebuilder:validation:Optional
	Settings *OidcSettings `json:"settings,omitempty"`
}

type Smtp struct {
	// +kubebuilder:validation:Type=string
	Host string `json:"host,omitempty"`

	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:validation:Maximum=65535
	// +kubebuilder:default=25
	Port int32 `json:"port,omitempty"`
}

// +kubebuilder:validation:AtLeastOneOf=smtp
type Email struct {
	// +kubebuilder:validation:Optional
	Smtp *Smtp `json:"smtp,omitempty"`
}

type Modules struct {
	// +kubebuilder:validation:Type=boolean
	// +kubebuilder:default=true
	Enabled bool `json:"enabled,omitempty"`

	// +kubebuilder:validation:Type=string
	// +kubebuilder:validation:MinLength=1
	Image string `json:"image"`
}

type MispInstanceSpec struct {
	// +kubebuilder:validation:Type=boolean
	// +kubebuilder:validation:Optional
	// +kubebuilder:default=false
	Suspend bool `json:"suspend,omitempty"`

	// +kubebuilder:validation:Type=string
	// +kubebuilder:validation:Optional
	Image string `json:"image,omitempty"`

	// +kubebuilder:validation:Optional
	// +kubebuilder:default=IfNotPresent
	ImagePullPolicy corev1.PullPolicy `json:"imagePullPolicy,omitempty"`

	// +kubebuilder:validation:Optional
	ImagePullSecrets []corev1.LocalObjectReference `json:"imagePullSecrets,omitempty"`

	// +kubebuilder:validation:Type=string
	// +kubebuilder:validation:MinLength=1
	BaseUrl string `json:"baseUrl"`

	// +kubebuilder:validation:Type=string
	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:MinLength=1
	Uuid string `json:"uuid,omitempty"`

	// +kubebuilder:validation:Type=string
	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:default=Etc/UTC
	TimeZone string `json:"timeZone"`

	// +kubebuilder:validation:Optional
	Admin *Admin `json:"admin,omitempty"`

	// +kubebuilder:validation:Optional
	Storage *Storage `json:"storage,omitempty"`

	Database Database `json:"database"`

	Cache Cache `json:"cache"`

	// +kubebuilder:validation:Optional
	Nginx *Nginx `json:"nginx,omitempty"`

	// +kubebuilder:validation:Optional
	Oidc *Oidc `json:"oidc,omitempty"`

	// +kubebuilder:validation:Optional
	Email *Email `json:"email"`

	// +kubebuilder:validation:Optional
	Modules *Modules `json:"modules"`

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
