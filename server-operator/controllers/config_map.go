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
	ctrl "sigs.k8s.io/controller-runtime"
)

func (r *SimpleServerReconciler) ReconcileConfigMap(ctx context.Context, log logr.Logger, server *sandboxv1alpha1.SimpleServer) error {
	log.Info("Start reconciling config map")
	configMap := &corev1.ConfigMap{}
	configMapName := fmt.Sprintf("%s-config", server.Name)
	err := r.Get(ctx, types.NamespacedName{Name: configMapName, Namespace: server.Namespace}, configMap)
	if err != nil && !errors.IsNotFound(err) {
		return err
	} else if err != nil && errors.IsNotFound(err) {
		log.Info("Creating a new ConfigMap", "ConfigMap.Namespace", configMap.Namespace, "ConfigMap.Name", configMap)
		configMap = createConfigMap(configMapName, server)
		log.Info("Service rendered")
		err = ctrl.SetControllerReference(server, configMap, r.Scheme)
		if err != nil {
			log.Info("Error failed to set owner reference", "error", err)
			return err
		}
		err = r.Create(ctx, configMap)
		if err != nil {
			log.Error(err, "Failed to create new Service", "Service.Namespace", configMap.Namespace, "Service.Name", configMap.Name)
			return err
		}
		// ConfigMap created successfully - return and requeue
		return nil
	}

	return err
}

func createConfigMap(configMapName string, server *sandboxv1alpha1.SimpleServer) *corev1.ConfigMap {
	configMap := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:        configMapName,
			Namespace:   server.Namespace,
			Labels:      make(map[string]string),
			Annotations: make(map[string]string),
		},
		Data: map[string]string{
			"property0": "value0",
			"property1": "value1",
			"property2": "value2",
		},
	}

	return configMap
}
