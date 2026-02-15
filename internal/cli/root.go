package cli

import (
	"errors"
	"flag"
	"fmt"
	"os"
)

func Execute() error {
	if len(os.Args) < 2 {
		usage()
		return errors.New("missing command")
	}

	switch os.Args[1] {
	case "version":
		fmt.Println("gpp v0.0.1")
		return nil
	case "tenant":
		return tenantCmd(os.Args[2:])
	case "deploy":
		return deployCmd(os.Args[2:])
	case "status":
		return statusCmd(os.Args[2:])
	default:
		usage()
		return fmt.Errorf("unknown command: %s", os.Args[1])
	}
}

func usage() {
	fmt.Println(`gpp - golden path platform CLI

Usage:
  gpp version
  gpp tenant create <tenant> [--kubeconfig <path>] [--context <ctx>]
  gpp deploy -f <gpp.yaml>   [--kubeconfig <path>] [--context <ctx>]
`)
}

type commonFlags struct {
	Kubeconfig string
	Context    string
}

func parseCommonFlags(fs *flag.FlagSet) commonFlags {
	var c commonFlags
	fs.StringVar(&c.Kubeconfig, "kubeconfig", "", "path to kubeconfig (optional)")
	fs.StringVar(&c.Context, "context", "", "kubeconfig context (optional)")
	return c
}
