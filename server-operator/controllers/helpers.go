package controllers

// makeLabels returns the labels for selecting the resources
// belonging to the given simpleServer CR name.
func makeLabels(name string) map[string]string {
	return map[string]string{"app": "simpleServer", "simpleServer_cr": name}
}
