package kube

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func EnsureNamespace(c *Client, name string) error {
	_, err := c.Kube.CoreV1().Namespaces().Get(context.Background(), name, metav1.GetOptions{})
	if err == nil {
		return nil
	}
	if !apierrors.IsNotFound(err) {
		return err
	}

	ns := &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
			Labels: map[string]string{
				"gpp.dev/tenant": name,
			},
		},
	}
	_, err = c.Kube.CoreV1().Namespaces().Create(context.Background(), ns, metav1.CreateOptions{})
	return err
}
