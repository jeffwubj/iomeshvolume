/*
Copyright 2022.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	cdiv1 "kubevirt.io/containerized-data-importer-api/pkg/apis/core/v1beta1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// IOMeshVolumeSpec defines the desired state of IOMeshVolume
type IOMeshVolumeSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	Node string                            `json:"node,omitempty"`
	PVC  *corev1.PersistentVolumeClaimSpec `json:"pvc,omitempty"`
}

// IOMeshVolumeStatus defines the observed state of IOMeshVolume
type IOMeshVolumeStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	Node       string                 `json:"node,omitempty"`
	DevicePath string                 `json:"devicePath,omitempty"`
	Phase      *cdiv1.DataVolumePhase `json:"phase,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// IOMeshVolume is the Schema for the iomeshvolumes API
type IOMeshVolume struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   IOMeshVolumeSpec   `json:"spec,omitempty"`
	Status IOMeshVolumeStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// IOMeshVolumeList contains a list of IOMeshVolume
type IOMeshVolumeList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []IOMeshVolume `json:"items"`
}

func init() {
	SchemeBuilder.Register(&IOMeshVolume{}, &IOMeshVolumeList{})
}
