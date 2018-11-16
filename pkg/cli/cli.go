package cli

import (
	"log"

	"github.com/ninedraft/protectron/pkg/bot"
	gpflag "github.com/octago/sflags/gen/gpflag"
	cobra "github.com/spf13/cobra"
)

type Flags struct {
	// https://github.com/octago/sflags#flags-based-on-structures------
	UseLocalTorProxy bool
	Token            string
}

func Command() *cobra.Command {
	var flags = Flags{}
	var cmd = &cobra.Command{
		Use: "command",
		Run: func(cmd *cobra.Command, args []string) {
			var proxy *bot.SOCKS5ProxyConfig
			if flags.UseLocalTorProxy {
				proxy = &bot.SOCKS5ProxyConfig{
					Host: "localhost",
					Port: 9050,
				}
			}
			var protector = bot.New(flags.Token, proxy)
			if err := protector.Run(); err != nil {
				log.Fatal(err)
			}
		},
	}
	if err := gpflag.ParseTo(&flags, cmd.PersistentFlags()); err != nil {
		panic(err)
	}

	return cmd
}
