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

package controllers

import (
	"context"
	"fmt"
	"github.com/pavel-v-chernykh/keystore-go/v4"
	"io/ioutil"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"os"
	"path"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"strings"
	"time"

	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	deliveryv1alpha1 "github.com/christian1607/jdk-cacerts-operator/api/v1alpha1"
)

const (
	CacertSecret     = "cacerts"
	CacertFileName   = "cacerts"
	CacertFolderName = "cacerts-tmp"
	CacertPassword   = "changeit"
	x509             = "X.509"
)

// JdkCaCertReconciler reconciles a JdkCaCert object
type JdkCaCertReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=delivery.caltamirano.com,namespace=jdk-operator,resources=jdkcacerts,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=delivery.caltamirano.com,namespace=jdk-operator,resources=jdkcacerts/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=core,namespace=jdk-operator,resources=secrets,verbs=get;update;patch;create;delete;list
func (r *JdkCaCertReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	ctx := context.Background()
	_ = r.Log.WithValues("jdkcacert", req.NamespacedName)

	jdkCacert := &deliveryv1alpha1.JdkCaCert{}
	err := r.Get(ctx, req.NamespacedName, jdkCacert)

	if err != nil {
		if errors.IsNotFound(err) {
			r.Log.Info(fmt.Sprintf("jdk-cacert %v not found", req.Name))
			return reconcile.Result{}, nil
		}
		return reconcile.Result{}, err
	}

	ks := keystore.New()
	crtsAdded := 0
	for _, s := range jdkCacert.Spec.Secrets {
		secret := &v1.Secret{}
		err = r.Get(ctx, types.NamespacedName{Namespace: jdkCacert.Namespace, Name: s}, secret)
		if err != nil {
			if errors.IsNotFound(err) {
				r.Log.Info(fmt.Sprintf("secret  %v not found, it won't be added to cacerts", s))
				continue
			}
			return reconcile.Result{}, err
		}

		for name, cert := range secret.Data {
			if strings.Contains(name, ".crt") {
				alias := strings.Split(name, ".")
				err = r.addCertificateKeystore(alias[0], cert, &ks)
				if err != nil {
					r.Log.Error(err,fmt.Sprintf("certificate  %v cannot be added,", s))
					continue
				}
				crtsAdded += 1
			}
		}
	}

	cacertFile, err := r.saveCacert(&ks)
	if err != nil {
		return reconcile.Result{}, err
	}

	cacertSecret := &v1.Secret{}
	err = r.Get(ctx, types.NamespacedName{Namespace: jdkCacert.Namespace, Name: CacertSecret}, cacertSecret)
	if err != nil {
		if errors.IsNotFound(err) {
			r.Log.Info(fmt.Sprintf("Secret '%v' not found, creating...", CacertSecret))
			err = r.Create(ctx, newJdkSecret(jdkCacert.Namespace, cacertFile))
			if err != nil {
				return reconcile.Result{}, err
			}
			err = jdkCacert.UpdateStatus(deliveryv1alpha1.JdkCaCertStatus{TotalSecrets: crtsAdded})
			if err != nil {
				return reconcile.Result{}, err
			}
			return reconcile.Result{}, r.Status().Update(ctx, jdkCacert)
		}
		return reconcile.Result{}, err
	}

	// the cacert exists, so just update it
	r.Log.Info(fmt.Sprintf("Updating secret  '%v'.", CacertSecret))
	m := map[string][]byte{CacertFileName: cacertFile}
	cacertSecret.Data = m
	err = r.Update(ctx, cacertSecret)
	if err != nil {
		return reconcile.Result{}, err
	}

	if jdkCacert.Status.TotalSecrets != crtsAdded {
		err = jdkCacert.UpdateStatus(deliveryv1alpha1.JdkCaCertStatus{TotalSecrets: crtsAdded})
		if err != nil {
			return reconcile.Result{}, err
		}

		return reconcile.Result{}, r.Status().Update(ctx, jdkCacert)
	}

	return reconcile.Result{}, err

}

func (r *JdkCaCertReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&deliveryv1alpha1.JdkCaCert{}).
		WithOptions(controller.Options{MaxConcurrentReconciles: 1}).
		Complete(r)
}

func newJdkSecret(ns string, data []byte) *v1.Secret {
	m := make(map[string][]byte)
	m[CacertFileName] = data
	return &v1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      CacertSecret,
			Namespace: ns,
			Labels:  map[string]string{"type":"jks"},
		},
		Type: v1.SecretTypeOpaque,
		Data: m,
	}
}

func (r *JdkCaCertReconciler) addCertificateKeystore(alias string, cert []byte, ks *keystore.KeyStore) error {

	return ks.SetTrustedCertificateEntry(alias, keystore.TrustedCertificateEntry{
		Certificate: keystore.Certificate{
			Content: cert,
			Type:    x509,
		},
		CreationTime: time.Now(),
	})
}

// Save all certificates added to ks and return the cacert file in []byte
func (r *JdkCaCertReconciler) saveCacert(ks *keystore.KeyStore) ([]byte, error) {

	r.Log.Info("Creating JDK tmp folder")
	err := os.Mkdir(CacertFolderName, 0775)

	if !os.IsExist(err) {
		return nil, err
	}

	r.Log.Info("Creating JDK truststore")
	o, err := os.Create(path.Join(CacertFolderName, CacertFileName))
	if err != nil {
		return nil, err
	}
	defer o.Close()

	e := ks.Store(o, []byte(CacertPassword))
	if e != nil {
		return nil, err
	}

	r.Log.Info("Reading JDK truststore")
	return ioutil.ReadFile(path.Join(CacertFolderName, CacertFileName))
}
