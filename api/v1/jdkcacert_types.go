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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// JdkCacertSpec defines the desired state of JdkCacert
type JdkCacertSpec struct {

	// Secrets wich contain certificates name of the k8s secrets to be added to the cacert secret
	Secrets []string `json:"secrets,omitempty"`

	OutputSecretName string `json:"output_secret_name,omitempty"`
}

// JdkCacertStatus defines the observed state of JdkCacert
type JdkCacertStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// JdkCacert is the Schema for the jdkcacerts API
type JdkCacert struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   JdkCacertSpec   `json:"spec,omitempty"`
	Status JdkCacertStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// JdkCacertList contains a list of JdkCacert
type JdkCacertList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []JdkCacert `json:"items"`
}

func init() {
	SchemeBuilder.Register(&JdkCacert{}, &JdkCacertList{})
}
