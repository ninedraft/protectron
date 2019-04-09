package protectron

import (
	gpflag "github.com/octago/sflags/gen/gpflag"
	cobra "github.com/spf13/cobra"
)

func Cli() *cobra.Command {
	var config Config // default values can be defined here
	var cmd = &cobra.Command{
		Use: "protectron",
		Run: func(cmd *cobra.Command, args []string) {
			Run(config)
		},
	}
	if err := gpflag.ParseTo(&config, cmd.PersistentFlags()); err != nil {
		panic(err)
	}
	return cmd
}
