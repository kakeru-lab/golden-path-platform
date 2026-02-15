package cli

import (
	"context"
	"errors"
	"flag"
	"fmt"

	"github.com/<your-module>/golden-path-platform/internal/config"
	"github.com/<your-module>/golden-path-platform/internal/kube"
	"github.com/<your-module>/golden-path-platform/internal/render"
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

	app, err := config.LoadApp(file)
	if err != nil {
		return err
	}

	kc, err := kube.NewClient(cf.Kubeconfig, cf.Context)
	if err != nil {
		return err
	}

	// tenant namespaceが無ければ作る（既にあるなら何もしない）
	if err := kube.EnsureNamespace(kc, app.Metadata.Namespace); err != nil {
		return err
	}

	// tenant defaults（Quota/NetPol/RBAC）も入れておきたいなら有効化
	if err := kube.ApplyTenantDefaults(kc, app.Metadata.Namespace); err != nil {
		return err
	}

	ksvcYAML, err := render.RenderKnativeService(app)
	if err != nil {
		return err
	}

	if err := kube.ApplyYAML(context.Background(), kc, app.Metadata.Namespace, ksvcYAML); err != nil {
		return err
	}

	fmt.Printf("Deployed Knative Service: %s/%s\n", app.Metadata.Namespace, app.Metadata.Name)
	return nil
}
