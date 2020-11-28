package cmd

import (
	"fmt"
	"path/filepath"

	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

func getKubeConfig() (*rest.Config, error) {
	kubeconfig := filepath.Join(homedir.HomeDir(), ".kube", "config")

	// use the current context in kubeconfig
	return clientcmd.BuildConfigFromFlags("", kubeconfig)
}

func buildResourceName(testName string, resourceType ...string) string {
	name := fmt.Sprintf("lokust-%s", testName)

	if len(resourceType) > 1 {
		name += "-" + resourceType[0]
	}

	return name
}
