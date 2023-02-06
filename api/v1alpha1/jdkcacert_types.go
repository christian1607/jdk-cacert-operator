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

package v1alpha1

import (
	"errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"time"
)

// JdkCaCertSpec defines the desired state of JdkCaCert
type JdkCaCertSpec struct {

	// Secrets name of the k8s secrets to be added to the cacert secret
	Secrets []string `json:"secrets,omitempty"`

	// JdkVersion specify the JDK version which is used, valid option 8 or 11 .
	JdkVersion string `json:"jdkVersion,omitempty"`
}

// JdkCaCertStatus defines the observed state of JdkCaCert
type JdkCaCertStatus struct {

	// +kubebuilder:validation:Minimum:=0

	// TotalSecrets total certificates added in the cacert secret
	TotalSecrets int `json:"totalSecrets"`

	// +kubebuilder:validation:Optional

	// LastSync last time sync.
	LastSync string `json:"lastSync"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// JdkCaCert is the Schema for the jdkcacerts API
type JdkCaCert struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   JdkCaCertSpec   `json:"spec,omitempty"`
	Status JdkCaCertStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// JdkCaCertList contains a list of JdkCaCert
type JdkCaCertList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []JdkCaCert `json:"items"`
}

func init() {
	SchemeBuilder.Register(&JdkCaCert{}, &JdkCaCertList{})
}

func (jdk *JdkCaCert) UpdateStatus(status JdkCaCertStatus) error {

	if status.TotalSecrets < 0 {
		return errors.New("total secrets cannot be less than zero")
	}

	jdk.Status.TotalSecrets = status.TotalSecrets
	jdk.Status.LastSync = time.Now().String()
	return nil
}
