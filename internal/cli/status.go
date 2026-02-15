package cli

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"strings"

	"github.com/kakeru-lab/golden-path-platform/internal/kube"
)

func statusCmd(args []string) error {
	fs := flag.NewFlagSet("status", flag.ContinueOnError)
	cf := parseCommonFlags(fs)

	var ns string
	fs.StringVar(&ns, "n", "", "namespace (tenant)")
	if err := fs.Parse(args); err != nil {
		return err
	}

	rest := fs.Args()
	if len(rest) < 1 {
		return errors.New("usage: gpp status -n <tenant> <service>")
	}
	name := rest[0]
	if ns == "" {
		return errors.New("namespace is required: -n <tenant>")
	}

	kc, err := kube.NewClient(cf.Kubeconfig, cf.Context)
	if err != nil {
		return err
	}

	u, err := kube.GetUnstructured(context.Background(), kc, "serving.knative.dev", "v1", "services", ns, name)
	if err != nil {
		return err
	}

	url := digString(u.Object, "status", "url")
	ready := digCondition(u.Object, "status", "conditions", "Ready")
	latest := digString(u.Object, "status", "latestReadyRevisionName")

	fmt.Printf("KService: %s/%s\n", ns, name)
	fmt.Printf("Ready  : %s\n", ready)
	if latest != "" {
		fmt.Printf("Revision: %s\n", latest)
	}
	if url != "" {
		fmt.Printf("URL    : %s\n", url)
	}

	return nil
}

func digString(obj map[string]any, path ...string) string {
	var cur any = obj
	for _, p := range path {
		m, ok := cur.(map[string]any)
		if !ok {
			return ""
		}
		cur, ok = m[p]
		if !ok {
			return ""
		}
	}
	s, _ := cur.(string)
	return s
}

func digCondition(obj map[string]any, statusKey, condKey, condType string) string {
	status, ok := obj[statusKey].(map[string]any)
	if !ok {
		return "unknown"
	}
	conds, ok := status[condKey].([]any)
	if !ok {
		return "unknown"
	}
	for _, c := range conds {
		m, ok := c.(map[string]any)
		if !ok {
			continue
		}
		if t, _ := m["type"].(string); t == condType {
			s, _ := m["status"].(string)
			r, _ := m["reason"].(string)
			msg, _ := m["message"].(string)
			if r != "" {
				if msg != "" {
					return fmt.Sprintf("%s (%s) %s", s, r, shorten(msg, 80))
				}
				return fmt.Sprintf("%s (%s)", s, r)
			}
			return s
		}
	}
	return "unknown"
}

func shorten(s string, n int) string {
	s = strings.TrimSpace(s)
	if len(s) <= n {
		return s
	}
	return s[:n] + "..."
}
