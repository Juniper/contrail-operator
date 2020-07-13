package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// FernetKeyManagerSpec defines the desired state of FernetKeyManager
type FernetKeyManagerSpec struct {
	TokenExpiration         int `json:"tokenExpiration"`
	TokenAllowExpiredWindow int `json:"tokenAllowExpiredWindow"`
	RotationInterval        int `json:"rotationFrequency"`
}

// FernetKeyManagerStatus defines the observed state of FernetKeyManager
type FernetKeyManagerStatus struct {
	MaxActiveKeys int `json:"maxActiveKeys"`
	SecretName string `json:"secretName"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// FernetKeyManager is the Schema for the fernetkeymanagers API
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=fernetkeymanagers,scope=Namespaced
type FernetKeyManager struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   FernetKeyManagerSpec   `json:"spec,omitempty"`
	Status FernetKeyManagerStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// FernetKeyManagerList contains a list of FernetKeyManager
type FernetKeyManagerList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []FernetKeyManager `json:"items"`
}

func init() {
	SchemeBuilder.Register(&FernetKeyManager{}, &FernetKeyManagerList{})
}
