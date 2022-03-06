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

package controllers

import (
	"context"
	"fmt"
	"k8s.io/apimachinery/pkg/api/errors"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	jvmv1alpha1 "jdk-cacert-operator/api/v1alpha1"
	jvmv1beta1 "jdk-cacert-operator/api/v1beta1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// JdkCacertReconciler reconciles a JdkCacert object
type JdkCacertReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=jvm.caltamirano.com,resources=jdkcacerts,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=jvm.caltamirano.com,resources=jdkcacerts/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=jvm.caltamirano.com,resources=jdkcacerts/finalizers,verbs=update
//+kubebuilder:rbac:groups=core,resources=secrets,verbs=get;list

func (r *JdkCacertReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)
	//logger:=zap.New()

	jdkCacert := &jvmv1beta1.JdkCacert{}
	err := r.Get(ctx, req.NamespacedName, jdkCacert)

	if err != nil {
		if errors.IsNotFound(err) {
			logger.Info(fmt.Sprintf("jdk-cacert %v not found", req.Name))
			return reconcile.Result{}, nil
		}
		return reconcile.Result{}, err
	}

	logger.Info(fmt.Sprintf("secret output: %v", jdkCacert.Spec.OutputSecretName))
	for _, s := range jdkCacert.Spec.Secrets {
		logger.Info(fmt.Sprintf("Secret: %v", s))
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *JdkCacertReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&jvmv1alpha1.JdkCacert{}).
		//Owns(&jvmv1beta1.JdkCacert{}).
		Complete(r)
}
