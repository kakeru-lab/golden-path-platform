package kube

import (
	"context"
	"fmt"

	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

func GetUnstructured(ctx context.Context, c *Client, group, version, resource, namespace, name string) (*unstructured.Unstructured, error) {
	dyn, err := c.dynamic()
	if err != nil {
		return nil, err
	}
	rm, err := c.restMapper()
	if err != nil {
		return nil, err
	}

	gvr := schema.GroupVersionResource{Group: group, Version: version, Resource: resource}
	// RESTMapperで scope を見て namespace有無を合わせたい場合は mapping を使う
	_ , _ = rm, meta.RESTScopeNameNamespace

	if namespace == "" {
		return dyn.Resource(gvr).Get(ctx, name, metav1.GetOptions{})
	}
	return dyn.Resource(gvr).Namespace(namespace).Get(ctx, name, metav1.GetOptions{})
}
