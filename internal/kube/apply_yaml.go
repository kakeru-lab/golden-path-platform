package kube

import (
	"bytes"
	"context"
	"fmt"
	"strings"

	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/yaml"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/discovery/cached/memory"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/restmapper"
)

func (c *Client) dynamic() (dynamic.Interface, error) {
	return dynamic.NewForConfig(c.Rest)
}

func (c *Client) restMapper() (meta.RESTMapper, error) {
	dc, err := discovery.NewDiscoveryClientForConfig(c.Rest)
	if err != nil {
		return nil, err
	}
	cache := memory.NewMemCacheClient(dc)
	return restmapper.NewDeferredDiscoveryRESTMapper(cache), nil
}

func ApplyYAML(ctx context.Context, c *Client, defaultNamespace string, yamlBytes []byte) error {
	dyn, err := c.dynamic()
	if err != nil {
		return err
	}
	rm, err := c.restMapper()
	if err != nil {
		return err
	}

	dec := yaml.NewYAMLOrJSONDecoder(bytes.NewReader(yamlBytes), 4096)

	for {
		var obj map[string]any
		if err := dec.Decode(&obj); err != nil {
			// EOFで終了
			if strings.Contains(err.Error(), "EOF") {
				break
			}
			return fmt.Errorf("decode yaml: %w", err)
		}
		if len(obj) == 0 {
			continue
		}

		u := &unstructured.Unstructured{Object: obj}

		gvk := u.GroupVersionKind()
		mapping, err := rm.RESTMapping(schema.GroupKind{Group: gvk.Group, Kind: gvk.Kind}, gvk.Version)
		if err != nil {
			return fmt.Errorf("rest mapping %s: %w", gvk.String(), err)
		}

		// namespace補完
		if mapping.Scope.Name() == meta.RESTScopeNameNamespace {
			if u.GetNamespace() == "" {
				u.SetNamespace(defaultNamespace)
			}
		}

		// Server-Side Apply
		data, err := u.MarshalJSON()
		if err != nil {
			return err
		}

		var ri dynamic.ResourceInterface
		if mapping.Scope.Name() == meta.RESTScopeNameNamespace {
			ri = dyn.Resource(mapping.Resource).Namespace(u.GetNamespace())
		} else {
			ri = dyn.Resource(mapping.Resource)
		}

		_, err = ri.Patch(
			ctx,
			u.GetName(),
			types.ApplyPatchType,
			data,
			metav1.PatchOptions{
				FieldManager: "gpp",
				Force:        ptr(true),
			},
		)
		if err != nil {
			return fmt.Errorf("apply %s/%s: %w", gvk.Kind, u.GetName(), err)
		}
	}

	return nil
}

func ptr[T any](v T) *T { return &v }
