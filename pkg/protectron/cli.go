package protectron

import (
	"time"

	gpflag "github.com/octago/sflags/gen/gpflag"
	cobra "github.com/spf13/cobra"
)

func Cli() *cobra.Command {
	var config = Config{
		Repost: 6 * time.Hour,
		Link:   3 * time.Hour,
	}
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
