package controllers

import "fmt"

// makeLabels returns the labels for selecting the resources
// belonging to the given simpleServer CR name.
func makeLabels(name string) map[string]string {
	return map[string]string{"app": "simpleServer", "simpleServer_cr": name}
}

func getConfigMapName(serverName string) string {
	return fmt.Sprintf("%s-config", serverName)
}

func getServiceName(serverName string) string {
	return fmt.Sprintf("%s-service", serverName)
}

func getDeploymentName(serverName string) string {
	return serverName
}
