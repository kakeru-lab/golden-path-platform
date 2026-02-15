package kube

import (
	"context"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

func GetUnstructured(ctx context.Context, c *Client, group, version, resource, namespace, name string) (*unstructured.Unstructured, error) {
	dyn, err := c.dynamic()
	if err != nil {
		return nil, err
	}

	gvr := schema.GroupVersionResource{Group: group, Version: version, Resource: resource}

	if namespace == "" {
		return dyn.Resource(gvr).Get(ctx, name, metav1.GetOptions{})
	}
	return dyn.Resource(gvr).Namespace(namespace).Get(ctx, name, metav1.GetOptions{})
}
