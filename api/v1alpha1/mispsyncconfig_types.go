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
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

type Organisation struct {
	// Name of the organisation to sync between instances
	// +required
	// +kubebuilder:validation:Type=string
	Name string `json:"name"`
}

type User struct {
	// Provide a secret in the same namespace with sync user credentials
	// +required
	CredentialsSecretRef corev1.LocalObjectReference `json:"credentialsSecretRef,omitempty"`
}

type MispSyncConfigSpec struct {
	// Organisation to sync between instances
	// +required
	Organisation Organisation `json:"organisation"`

	// User to use for the sync
	// +required
	User User `json:"user"`

	// Reference an instance object
	// +required
	InstanceRef corev1.ObjectReference `json:"instanceRef"`

	// Reference an instance object
	// +required
	RemoteInstanceRef corev1.ObjectReference `json:"remoteInstanceRef"`
}

// MispSyncConfigStatus defines the observed state of MispSyncConfig.
type MispSyncConfigStatus struct {
	// The status of each condition is one of True, False, or Unknown.
	// +patchMergeKey=type
	// +patchStrategy=merge
	// +listType=map
	// +listMapKey=type
	// +optional
	Conditions []metav1.Condition `json:"conditions,omitempty" patchStrategy:"merge" patchMergeKey:"type"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Phase",type=string,JSONPath=`.status.phase`
// +kubebuilder:resource:shortName=mispsync

// MispSyncConfig is the Schema for the mispsyncconfigs API
type MispSyncConfig struct {
	metav1.TypeMeta `json:",inline"`

	// metadata is a standard object metadata
	// +optional
	metav1.ObjectMeta `json:"metadata,omitzero"`

	// spec defines the desired state of MispSyncConfig
	// +required
	Spec MispSyncConfigSpec `json:"spec"`

	// status defines the observed state of MispSyncConfig
	// +optional
	Status MispSyncConfigStatus `json:"status,omitzero"`
}

// +kubebuilder:object:root=true

// MispSyncConfigList contains a list of MispSyncConfig
type MispSyncConfigList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitzero"`
	Items           []MispSyncConfig `json:"items"`
}

func init() {
	SchemeBuilder.Register(func(s *runtime.Scheme) error {
		s.AddKnownTypes(SchemeGroupVersion, &MispSyncConfig{}, &MispSyncConfigList{})
		return nil
	})
}
