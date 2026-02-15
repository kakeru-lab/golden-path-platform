package kube

import (
	"context"
	"strings"

	"github.com/<your-org-or-user>/golden-path-platform/internal/templates"
)

func ApplyTenantDefaults(c *Client, tenant string) error {
	ctx := context.Background()

	repl := func(b []byte) []byte {
		return []byte(strings.ReplaceAll(string(b), "{{NAMESPACE}}", tenant))
	}

	// 依存の薄い順（Quota/Limit/NetPol/RBAC）でOK
	if err := ApplyYAML(ctx, c, tenant, repl(templates.TenantQuota)); err != nil {
		return err
	}
	if err := ApplyYAML(ctx, c, tenant, repl(templates.TenantLimitRange)); err != nil {
		return err
	}
	if err := ApplyYAML(ctx, c, tenant, repl(templates.TenantNetworkPolicy)); err != nil {
		return err
	}
	if err := ApplyYAML(ctx, c, tenant, repl(templates.TenantRBAC)); err != nil {
		return err
	}

	return nil
}
