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

package v1beta1

import (
	v1 "jdk-cacert-operator/api/v1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/conversion"
)

// ConvertTo converts this JdkCacert to the Hub version (v1).
func (src *JdkCacert) ConvertTo(dstRaw conversion.Hub) error {
	dst := dstRaw.(*v1.JdkCacert)
	dst.ObjectMeta = src.ObjectMeta
	dst.Spec.Secrets = src.Spec.Secrets
	dst.Spec.OutputSecretName = "secretfromconverter"
	return nil
}

// ConvertFrom converts from the Hub version (v1) to this version.
func (dst *JdkCacert) ConvertFrom(srcRaw conversion.Hub) error {
	src := srcRaw.(*v1.JdkCacert)
	dst.ObjectMeta = src.ObjectMeta
	dst.Spec.Secrets = src.Spec.Secrets
	return nil
}

func (r *JdkCacert) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		Complete()
}
