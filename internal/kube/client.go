package kube

import (
	"path/filepath"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

type Client struct {
	Kube   kubernetes.Interface
	Rest   *rest.Config
}

func NewClient(kubeconfigPath, context string) (*Client, error) {
	var cfg *rest.Config
	var err error

	if kubeconfigPath == "" {
		// in-cluster or default loading rules
		cfg, err = rest.InClusterConfig()
		if err != nil {
			// fallback to default kubeconfig (~/.kube/config)
			loadingRules := clientcmd.NewDefaultClientConfigLoadingRules()
			overrides := &clientcmd.ConfigOverrides{}
			if context != "" {
				overrides.CurrentContext = context
			}
			cfg, err = clientcmd.NewNonInteractiveDeferredLoadingClientConfig(loadingRules, overrides).ClientConfig()
		}
	} else {
		abs, _ := filepath.Abs(kubeconfigPath)
		loadingRules := &clientcmd.ClientConfigLoadingRules{ExplicitPath: abs}
		overrides := &clientcmd.ConfigOverrides{}
		if context != "" {
			overrides.CurrentContext = context
		}
		cfg, err = clientcmd.NewNonInteractiveDeferredLoadingClientConfig(loadingRules, overrides).ClientConfig()
	}
	if err != nil {
		return nil, err
	}

	k, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		return nil, err
	}

	return &Client{Kube: k, Rest: cfg}, nil
}
