package kube

import (
	"encoding/base64"
	"fmt"
)

type ClusterInfo struct {
	Server   string // https://x.x.x.x:6443
	CABase64 string // kubeconfigに入れる CA data (base64)
}

func GetClusterInfo(c *Client) (*ClusterInfo, error) {
	if c.Rest == nil {
		return nil, fmt.Errorf("rest config is nil")
	}
	if c.Rest.Host == "" {
		return nil, fmt.Errorf("rest config host is empty")
	}

	var caB64 string
	if len(c.Rest.CAData) > 0 {
		caB64 = base64.StdEncoding.EncodeToString(c.Rest.CAData)
	}
	// CAFile の場合は、ここでは省略（必要なら後で対応）
	return &ClusterInfo{
		Server:   c.Rest.Host,
		CABase64: caB64,
	}, nil
}
