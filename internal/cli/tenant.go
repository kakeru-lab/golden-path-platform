package cli

import (
	"errors"
	"flag"
	"fmt"

	"github.com/<your-org-or-user>/golden-path-platform/internal/kube"
)

func tenantCmd(args []string) error {
	if len(args) < 1 {
		return errors.New("missing tenant subcommand")
	}

	switch args[0] {
	case "create":
		return tenantCreate(args[1:])
	case "kubeconfig":
        return tenantKubeconfig(args[1:])
	default:
		return fmt.Errorf("unknown tenant subcommand: %s", args[0])
	}
}

func tenantCreate(args []string) error {
	fs := flag.NewFlagSet("tenant create", flag.ContinueOnError)
	cf := parseCommonFlags(fs)
	if err := fs.Parse(args); err != nil {
		return err
	}

	rest := fs.Args()
	if len(rest) < 1 {
		return errors.New("usage: gpp tenant create <tenant>")
	}
	tenant := rest[0]

	kc, err := kube.NewClient(cf.Kubeconfig, cf.Context)
	if err != nil {
		return err
	}

	// ここでは「namespace作成」だけ先に実装（次でRBAC/Quota/NetPolも適用する）
	if err := kube.EnsureNamespace(kc, tenant); err != nil {
		return err
	}

	fmt.Printf("Tenant created: %s\n", tenant)
	return nil
}
