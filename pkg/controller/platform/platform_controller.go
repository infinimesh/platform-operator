/*
Copyright 2019 infinimesh, inc.
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

package platform

import (
	"context"
	"time"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	infinimeshv1beta1 "github.com/infinimesh/enterprise-operator/pkg/apis/infinimesh/v1beta1"
)

// Add creates a new Platform Controller and adds it to the Manager with default RBAC. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcilePlatform{Client: mgr.GetClient(), scheme: mgr.GetScheme()}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	return ctrl.NewControllerManagedBy(mgr).
		Named("platform-controller").
		For(&infinimeshv1beta1.Platform{}).
		Owns(&appsv1.Deployment{}).
		Owns(&corev1.Service{}).
		Owns(&networkingv1.Ingress{}).
		Owns(&corev1.Secret{}).
		Complete(r)
}

var _ reconcile.Reconciler = &ReconcilePlatform{}

// ReconcilePlatform reconciles a Platform object
type ReconcilePlatform struct {
	client.Client
	scheme *runtime.Scheme
}

// Reconcile reads that state of the cluster for a Platform object and makes changes based on the state read
// and what is in the Platform.Spec
// TODO(user): Modify this Reconcile function to implement your Controller logic.  The scaffolding writes
// a Deployment as an example
// Automatically generate RBAC rules to allow the Controller to read and write Deployments
// +kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=core,resources=configmaps,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=apps,resources=statefulsets,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=core,resources=secrets;services,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=networking.k8s.io,resources=ingresses,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=apps,resources=deployments/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=infinimesh.infinimesh.io,resources=platforms,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=infinimesh.infinimesh.io,resources=platforms/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=kubedb.com,resources=postgreses,verbs=get;list;watch;create;update;patch;delete
func (r *ReconcilePlatform) Reconcile(ctx context.Context, request reconcile.Request) (reconcile.Result, error) {
	instance := &infinimeshv1beta1.Platform{}

	err := r.Get(ctx, request.NamespacedName, instance)
	if err != nil {
		if errors.IsNotFound(err) {
			// Object not found, return.  Created objects are automatically garbage collected.
			// For additional cleanup logic use finalizers.
			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return reconcile.Result{}, err
	}

	now := time.Now().UTC()
	tag, write, after, err := resolveCoreVersion(instance.Spec.Version, instance.Status.Version, instance.Status.CheckedAt, now, func() (string, error) {
		return fetchNewestCore(ctx)
	})
	if err != nil {
		return reconcile.Result{}, err
	}
	if write {
		instance.Status.Version = tag
		if instance.Spec.Version == autoVersion {
			instance.Status.CheckedAt = now.Format(time.RFC3339)
		}
		if err := r.Status().Update(ctx, instance); err != nil {
			return reconcile.Result{}, err
		}
	} else {
		instance.Status.Version = tag
	}

	if err := r.reconcileCore(instance); err != nil {
		return reconcile.Result{}, err
	}

	return reconcile.Result{RequeueAfter: after}, nil
}
