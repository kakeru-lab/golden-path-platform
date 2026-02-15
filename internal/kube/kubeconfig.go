package kube

import (
	"fmt"
)

type KubeconfigArgs struct {
	ClusterName string
	UserName    string
	ContextName string
	Namespace   string
	Server      string
	CABase64    string
	Token       string
}

func RenderKubeconfig(a KubeconfigArgs) ([]byte, error) {
	if a.ClusterName == "" {
		a.ClusterName = "gpp-cluster"
	}
	if a.UserName == "" {
		a.UserName = "gpp-user"
	}
	if a.ContextName == "" {
		a.ContextName = "gpp"
	}
	if a.Namespace == "" {
		return nil, fmt.Errorf("namespace is required")
	}
	if a.Server == "" {
		return nil, fmt.Errorf("server is required")
	}
	if a.Token == "" {
		return nil, fmt.Errorf("token is required")
	}

	// CA data が無いケースもある（その場合は insecure じゃなく、CAFile対応を後で足すのが安全）
	caBlock := ""
	if a.CABase64 != "" {
		caBlock = fmt.Sprintf("    certificate-authority-data: %s\n", a.CABase64)
	} else {
		// とりあえず空にして生成（後で CAFile 対応 or 明示エラーにしてもOK）
		// 本気運用なら CA 無しはNGにするのが良い
	}

	out := fmt.Sprintf(`apiVersion: v1
kind: Config
clusters:
- name: %s
  cluster:
    server: %s
%susers:
- name: %s
  user:
    token: %s
contexts:
- name: %s
  context:
    cluster: %s
    user: %s
    namespace: %s
current-context: %s
`, a.ClusterName, a.Server, caBlock, a.UserName, a.Token, a.ContextName, a.ClusterName, a.UserName, a.Namespace, a.ContextName)

	return []byte(out), nil
}
