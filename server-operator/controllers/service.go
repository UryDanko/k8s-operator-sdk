package controllers

import (
	"context"
	"fmt"
	"github.com/go-logr/logr"
	sandboxv1alpha1 "gitlab.com/sandbox/simple-server-operator/api/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/intstr"
	ctrl "sigs.k8s.io/controller-runtime"
)

// ReconcileService makes sure that the corresponding service exists
func (r *SimpleServerReconciler) ReconcileService(ctx context.Context, log logr.Logger, server *sandboxv1alpha1.SimpleServer) error {
	log.Info("Start reconciling service")
	service := &corev1.Service{}
	serviceName := getServiceName(server.Name)
	err := r.Get(ctx, types.NamespacedName{Name: serviceName, Namespace: server.Namespace}, service)
	if err != nil && !errors.IsNotFound(err) {
		return err
	} else if err != nil && errors.IsNotFound(err) {
		log.Info("Creating a new Service", "Service.Namespace", service.Namespace, "Service.Name", serviceName)
		service = createService(serviceName, server)
		log.Info("Service rendered")
		err = ctrl.SetControllerReference(server, service, r.Scheme)
		if err != nil {
			log.Info("Error failed to set owner reference", "error", err)
			return err
		}
		err = r.Create(ctx, service)
		if err != nil {
			log.Error(err, "Failed to create new Service", "Service.Namespace", service.Namespace, "Service.Name", service.Name)
			return err
		}
		// Service created successfully - return and requeue
		return nil
	}

	return err
}

// createService returns a SimpleServer Service object
func createService(serviceName string, server *sandboxv1alpha1.SimpleServer) *corev1.Service {
	service := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:        serviceName,
			Namespace:   server.Namespace,
			Labels:      make(map[string]string),
			Annotations: make(map[string]string),
		},
		Spec: corev1.ServiceSpec{
			Ports: []corev1.ServicePort{{
				Name: fmt.Sprintf("%s-service-tcp-port", server.Name),
				Port: 8080,
				TargetPort: intstr.IntOrString{
					Type:   intstr.Int,
					IntVal: 8080,
				},
			}},
			Selector: makeLabels(server.Name),
			Type:     corev1.ServiceTypeClusterIP,
		},
	}

	return service
}
