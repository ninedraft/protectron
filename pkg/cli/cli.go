package cli

import (
	gpflag "github.com/octago/sflags/gen/gpflag"
	cobra "github.com/spf13/cobra"
)

type Flags struct {
	// flag definitions here
	// https://github.com/octago/sflags#flags-based-on-structures------
}

func Command() *cobra.Command {
	var flags Flags // default values can be defined here
	var cmd = &cobra.Command{
		Use: "command",
		Run: func(cmd *cobra.Command, args []string) {
			// write your code here
		},
	}
	if err := gpflag.ParseTo(&flags, cmd.PersistentFlags()); err != nil {
		panic(err)
	}
	return cmd
}
