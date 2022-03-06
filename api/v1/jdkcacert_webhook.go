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
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/validation/field"
	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
)

// log is for logging in this package.
var jdkcacertlog = logf.Log.WithName("jdkcacert-resource")

func (r *JdkCacert) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		Complete()
}

//+kubebuilder:webhook:path=/mutate-jvm-caltamirano-com-v1-jdkcacert,mutating=true,failurePolicy=fail,groups=jvm.caltamirano.com,resources=jdkcacerts,verbs=create;update,versions=v1,name=mjdkcacert.caltamirano.com,sideEffects=None,admissionReviewVersions=v1

var _ webhook.Defaulter = &JdkCacert{}

// Default implements webhook.Defaulter so a webhook will be registered for the type
func (r *JdkCacert) Default() {
	jdkcacertlog.Info("default", "name", r.Name)

	if r.Spec.OutputSecretName == "" {
		r.Spec.OutputSecretName = ""
	}

	if r.Spec.Secrets == nil {
		r.Spec.Secrets = make([]string, 0)
	}
}

//+kubebuilder:webhook:verbs=create;update;delete,path=/validate-jvm-caltamirano-com-v1-jdkcacert,mutating=false,failurePolicy=fail,groups=jvm.caltamirano.com,versions=v1,name=vjdkcacert.caltamirano.com,resources=jdkcacerts,sideEffects=None,admissionReviewVersions=v1

var _ webhook.Validator = &JdkCacert{}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type
func (r *JdkCacert) ValidateCreate() error {
	jdkcacertlog.Info("validate create", "name", r.Name)
	return r.ValidateJdkCacert()
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type
func (r *JdkCacert) ValidateUpdate(old runtime.Object) error {
	jdkcacertlog.Info("validate update", "name", r.Name)
	return r.ValidateJdkCacert()
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type
func (r *JdkCacert) ValidateDelete() error {
	jdkcacertlog.Info("validate delete", "name", r.Name)

	// TODO(user): fill in your validation logic upon object deletion.
	return nil
}

func (r *JdkCacert) ValidateJdkCacert() error {
	var errorsList field.ErrorList
	if len(r.Spec.Secrets) == 0 {
		errorsList = append(errorsList, &field.Error{Field: "secrets", Detail: "Secrets cannot be empty", Type: field.ErrorTypeInvalid, BadValue: r.Spec.Secrets})
	}

	if len(errorsList) == 0 {
		return nil
	}

	return apierrors.NewInvalid(schema.GroupKind{Kind: "JdkCacert", Group: "jvm.caltamirano.com"}, r.Name, errorsList)
}
