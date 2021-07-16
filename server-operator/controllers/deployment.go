package controllers

import (
	"context"
	"github.com/go-logr/logr"
	sandboxv1alpha1 "gitlab.com/sandbox/simple-server-operator/api/v1alpha1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"reflect"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func (r *SimpleServerReconciler) ReconcileDeployment(ctx context.Context, log logr.Logger, server *sandboxv1alpha1.SimpleServer) error {
	log.Info("Start reconciling Deployment")

	// Check if the deployment already exists, if not create a new one
	depl := &appsv1.Deployment{}
	err := r.Get(ctx, types.NamespacedName{Name: server.Name, Namespace: server.Namespace}, depl)
	if err != nil && errors.IsNotFound(err) {
		// Define a new deployment
		dep := createDeployment(server)
		// Set instance as the owner and controller
		err := ctrl.SetControllerReference(server, dep, r.Scheme)
		if err != nil {
			log.Info("Error failed to set owner reference", "error", err)
			return err
		}

		log.Info("Creating a new Deployment", "Deployment.Namespace", dep.Namespace, "Deployment.Name", dep.Name)
		err = r.Create(ctx, dep)
		if err != nil {
			log.Error(err, "Failed to create new Deployment", "Deployment.Namespace", dep.Namespace, "Deployment.Name", dep.Name)
			return err
		}
		// Deployment created successfully - return and requeue
		return nil
	} else if err != nil {
		log.Error(err, "Failed to get Deployment")
		return err
	}

	// Ensure the deployment size is the same as the spec
	size := server.Spec.Size
	if *depl.Spec.Replicas != size {
		depl.Spec.Replicas = &size
		err = r.Update(ctx, depl)
		if err != nil {
			log.Error(err, "Failed to update Deployment", "Deployment.Namespace", depl.Namespace, "Deployment.Name", depl.Name)
			return err
		}
		// Ask to requeue after 1 minute in order to give enough time for the
		// pods be created on the cluster side and the operand be able
		// to do the next update step accurately.
		//return ctrl.Result{RequeueAfter: time.Minute}, nil  // todo: improve?
		return nil
	}

	// Update the SimpleServer status with the pod names
	// List the pods for this SimpleServer's deployment
	podList := &corev1.PodList{}
	listOpts := []client.ListOption{
		client.InNamespace(server.Namespace),
		client.MatchingLabels(makeLabels(server.Name)),
	}
	if err = r.List(ctx, podList, listOpts...); err != nil {
		log.Error(err, "Failed to list pods", "SimpleServer.Namespace", server.Namespace, "SimpleServer.Name", server.Name)
		return err
	}
	podNames := getPodNames(podList.Items)

	// Update status.Nodes if needed
	if !reflect.DeepEqual(podNames, server.Status.Nodes) {
		server.Status.Nodes = podNames
		err := r.Status().Update(ctx, server)
		if err != nil {
			log.Error(err, "Failed to update SimpleServer status")
			return err
		}
	}
	return nil
}

// createDeployment returns a SimpleServer Deployment object
func createDeployment(server *sandboxv1alpha1.SimpleServer) *appsv1.Deployment {
	ls := makeLabels(server.Name)
	replicas := server.Spec.Size

	dep := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      server.Name,
			Namespace: server.Namespace,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: ls,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: ls,
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{{
						Image: "ydanko/simple-server:latest",
						Name:  "simpleserver",
						Ports: []corev1.ContainerPort{{
							ContainerPort: 8080,
							Name:          "simpleserver",
						}},
					}},
				},
			},
		},
	}
	return dep
}

// getPodNames returns the pod names of the array of pods passed in
func getPodNames(pods []corev1.Pod) []string {
	var podNames []string
	for _, pod := range pods {
		podNames = append(podNames, pod.Name)
	}
	return podNames
}
