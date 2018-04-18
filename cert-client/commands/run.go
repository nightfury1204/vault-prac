package commands

import (
	"github.com/golang/glog"
	"github.com/nightfury1204/vault-prac/cert-client/pkg"
	"github.com/spf13/cobra"
)

func NewRunCmd() *cobra.Command {
	opts := pkg.NewOptions()

	cmd := &cobra.Command{
		Use:               "run",
		Short:             `get client cert`,
		DisableAutoGenTag: true,
		Run: func(cmd *cobra.Command, args []string) {
			if err := opts.Validate(); err != nil {
				glog.Fatal(err)
			}
			if err := opts.Run(); err != nil {
				glog.Fatal(err)
			}
		},
	}
	opts.AddFlags(cmd.Flags())

	return cmd
}
