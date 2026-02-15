package kube

import (
	"context"
	"fmt"
	"time"

	authenticationv1 "k8s.io/api/authentication/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// RequestServiceAccountToken: TokenRequest で短命トークンを発行
func RequestServiceAccountToken(c *Client, namespace, serviceAccount string, ttl time.Duration) (string, error) {
	if ttl <= 0 {
		ttl = 1 * time.Hour
	}
	sec := int64(ttl.Seconds())

	tr := &authenticationv1.TokenRequest{
		Spec: authenticationv1.TokenRequestSpec{
			ExpirationSeconds: &sec,
		},
	}

	resp, err := c.Kube.CoreV1().
		ServiceAccounts(namespace).
		CreateToken(context.Background(), serviceAccount, tr, metav1.CreateOptions{})
	if err != nil {
		return "", fmt.Errorf("create token for sa %s/%s: %w", namespace, serviceAccount, err)
	}
	return resp.Status.Token, nil
}
