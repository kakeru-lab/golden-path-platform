package cli

import (
	"errors"
	"flag"
	"fmt"

	"github.com/<your-org-or-user>/golden-path-platform/internal/kube"
)

func deployCmd(args []string) error {
	fs := flag.NewFlagSet("deploy", flag.ContinueOnError)
	cf := parseCommonFlags(fs)

	var file string
	fs.StringVar(&file, "f", "", "path to gpp.yaml")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if file == "" {
		return errors.New("usage: gpp deploy -f <gpp.yaml>")
	}

	kc, err := kube.NewClient(cf.Kubeconfig, cf.Context)
	if err != nil {
		return err
	}

	// 次ステップで:
	// 1) gpp.yaml読み込み
	// 2) Knative Serviceテンプレに差し込み
	// 3) apply
	_ = kc

	fmt.Printf("Deploy placeholder OK (file=%s)\n", file)
	return nil
}
