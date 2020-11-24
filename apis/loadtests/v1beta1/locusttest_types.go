/*

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

package v1beta1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.
type LocustTestResources struct {
	Master  corev1.ResourceRequirements `json:"master"`
	Workers corev1.ResourceRequirements `json:"workers"`
}

// LocustTestSpec defines the desired state of LocustTest
type LocustTestSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// +kubebuilder:default:1
	// Number of workers
	Replicas *int32 `json:"replicas,omitempty"`

	// +optional
	// Test resources
	Resources LocustTestResources `json:"resources,omitempty"`
}

// LocustTestStatus defines the observed state of LocustTest
type LocustTestStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Number of running workers
	Replicas *int32 `json:"replicas,omitempty"`
}

// +genclient

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:subresource:scale:specpath=.spec.replicas,statuspath=.status.replicas,selectorpath=.status.selector
// LocustTest is the Schema for the locusttests API
type LocustTest struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   LocustTestSpec   `json:"spec,omitempty"`
	Status LocustTestStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// LocustTestList contains a list of LocustTest
type LocustTestList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []LocustTest `json:"items"`
}

func init() {
	SchemeBuilder.Register(&LocustTest{}, &LocustTestList{})
}
