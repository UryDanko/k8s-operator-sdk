/*
Copyright 2021.

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
	"fmt"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"time"

	"context"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	ctrllog "sigs.k8s.io/controller-runtime/pkg/log"

	sandboxv1alpha1 "gitlab.com/sandbox/simple-server-operator/api/v1alpha1"
)

// SimpleServerReconciler reconciles a SimpleServer object
type SimpleServerReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=sandbox.example.com,resources=simpleservers,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=sandbox.example.com,resources=simpleservers/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=sandbox.example.com,resources=simpleservers/finalizers,verbs=update
//+kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=apps,resources=deployments/status,verbs=get
//+kubebuilder:rbac:groups=core,resources=pods,verbs=get;list;watch
//+kubebuilder:rbac:groups="",resources=configmaps,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups="",resources=configmaps/status,verbs=get
//+kubebuilder:rbac:groups="",resources=services,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups="",resources=services/status,verbs=get

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the SimpleServer object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.8.3/pkg/reconcile
func (r *SimpleServerReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := ctrllog.FromContext(ctx)
	timeout := 60 * time.Second
	log.Info("Start reconciling loop")

	// Fetch the SimpleServer instance
	simpleServer := &sandboxv1alpha1.SimpleServer{}
	err := r.Get(ctx, req.NamespacedName, simpleServer)
	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			log.Info("SimpleServer resource not found. Ignoring since object must be deleted")
			return ctrl.Result{}, nil
		}
		// Error reading the object - requeue the request.
		log.Error(err, "Failed to get SimpleServer")
		return ctrl.Result{}, err
	}

	err = r.ReconcileConfigMap(ctx, log, simpleServer)
	if err != nil {
		log.Error(err, fmt.Sprintf("failed to reconcile Service, retrying in %ss", timeout))
		return ctrl.Result{RequeueAfter: timeout}, err
	}

	err = r.ReconcileDeployment(ctx, log, simpleServer)
	if err != nil {
		log.Error(err, fmt.Sprintf("failed to reconcile Deployment, retrying in %ss", timeout))
		return ctrl.Result{RequeueAfter: timeout}, err
	}

	err = r.ReconcileService(ctx, log, simpleServer)
	if err != nil {
		log.Error(err, fmt.Sprintf("failed to reconcile Service, retrying in %ss", timeout))
		return ctrl.Result{RequeueAfter: timeout}, err
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *SimpleServerReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&sandboxv1alpha1.SimpleServer{}).
		Owns(&appsv1.Deployment{}).
		Owns(&corev1.Service{}).
		Owns(&corev1.ConfigMap{}).
		WithOptions(controller.Options{MaxConcurrentReconciles: 1}).
		Complete(r)
}
