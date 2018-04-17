package commands

import (
	"flag"
	"github.com/spf13/cobra"
)

func NewRootCmd() *cobra.Command {

	cmd := &cobra.Command{
		Use:                "vault-bootstrapper",
		Short:              `vault initial setting`,
		DisableAutoGenTag:  true,
		DisableFlagParsing: true,
	}
	cmd.PersistentFlags().AddGoFlagSet(flag.CommandLine)
	// ref: https://github.com/kubernetes/kubernetes/issues/17162#issuecomment-225596212
	flag.CommandLine.Parse([]string{})

	cmd.AddCommand(NewRunCmd())

	return cmd
}
