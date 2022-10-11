package kubernetes

import (
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

// NewClient creates a new Kubernetes client using the default configuration loader
func NewClient() (kubernetes.Interface, error) {
	rules := clientcmd.NewDefaultClientConfigLoadingRules()
	overrides := &clientcmd.ConfigOverrides{}
	loader := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(rules, overrides)

	kubernetesConfig, err := loader.ClientConfig()
	if err != nil {
		return nil, err
	}

	client, err := kubernetes.NewForConfig(kubernetesConfig)
	if err != nil {
		return nil, err
	}

	return client, nil
}
