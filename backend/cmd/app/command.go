package app

import (
	"github.com/urfave/cli/v2"
	"wallet/cmd/app/example"
	"wallet/cmd/app/migrate"
	"wallet/cmd/app/rewrite"
	"wallet/cmd/app/server"
)

var commands = []*cli.Command{
	{
		Name:                   "server",
		Usage:                  "HTTP Server",
		UseShortOptionHandling: true,
		Action:                 server.Action,
		Flags:                  server.Flags,
	},
	{
		Name:                   "migrate",
		Usage:                  "Migrate Database",
		UseShortOptionHandling: true,
		Action:                 migrate.Action,
		Flags:                  migrate.Flags,
	},
	{
		Name:                   "rewrite",
		Usage:                  "Rewrite Database",
		UseShortOptionHandling: true,
		Action:                 rewrite.Action,
		Flags:                  rewrite.Flags,
	},
	{
		Name:                   "example",
		Usage:                  "Generate Example",
		UseShortOptionHandling: true,
		Action:                 example.Action,
		Flags:                  example.Flags,
	},
}
