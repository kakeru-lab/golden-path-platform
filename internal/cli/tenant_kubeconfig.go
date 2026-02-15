package cli

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/<your-module>/golden-path-platform/internal/kube"
)

func tenantKubeconfig(args []string) error {
	fs := flag.NewFlagSet("tenant kubeconfig", flag.ContinueOnError)
	cf := parseCommonFlags(fs)

	var out string
	var ttl string
	fs.StringVar(&out, "out", "", "output file path (e.g. tenant-a.kubeconfig)")
	fs.StringVar(&ttl, "ttl", "1h", "token TTL (e.g. 30m, 1h, 8h)")
	if err := fs.Parse(args); err != nil {
		return err
	}

	rest := fs.Args()
	if len(rest) < 1 {
		return errors.New("usage: gpp tenant kubeconfig <tenant> --out <file>")
	}
	tenant := rest[0]
	if out == "" {
		return errors.New("--out is required")
	}

	d, err := time.ParseDuration(ttl)
	if err != nil {
		return fmt.Errorf("invalid --ttl: %w", err)
	}

	kc, err := kube.NewClient(cf.Kubeconfig, cf.Context)
	if err != nil {
		return err
	}

	// SAは tenantテンプレで作ってる想定
	token, err := kube.RequestServiceAccountToken(kc, tenant, "gpp-deployer", d)
	if err != nil {
		return err
	}

	ci, err := kube.GetClusterInfo(kc)
	if err != nil {
		return err
	}

	kcfg, err := kube.RenderKubeconfig(kube.KubeconfigArgs{
		ClusterName: "gpp-cluster",
		UserName:    "gpp-" + tenant,
		ContextName: "gpp-" + tenant,
		Namespace:   tenant,
		Server:      ci.Server,
		CABase64:    ci.CABase64,
		Token:       token,
	})
	if err != nil {
		return err
	}

	if err := os.WriteFile(out, kcfg, 0o600); err != nil {
		return err
	}

	fmt.Printf("Wrote kubeconfig: %s\n", out)
	return nil
}
