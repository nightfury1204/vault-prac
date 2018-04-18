package commands

import (
	"github.com/golang/glog"
	"github.com/nightfury1204/vault-prac/cert-issuer/pkg"
	"github.com/spf13/cobra"
)

func NewRunCmd() *cobra.Command {
	opts := pkg.NewOptions()

	cmd := &cobra.Command{
		Use:               "run",
		Short:             `cert-issuer`,
		DisableAutoGenTag: true,
		Run: func(cmd *cobra.Command, args []string) {
			if err := opts.Validate(); err != nil {
				glog.Fatal(err)
			}
			if err := opts.RunServer(); err != nil {
				glog.Fatal(err)
			}
		},
	}
	opts.AddFlags(cmd.Flags())

	return cmd
}
