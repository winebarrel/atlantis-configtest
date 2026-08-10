package main

import (
	"os"

	"github.com/alecthomas/kong"
	configtest "github.com/winebarrel/atlantis-configtest"
)

var version string

var cli struct {
	Version    kong.VersionFlag
	ServerRepo configtest.ServerRepoCmd `cmd:"" help:"Validate server-side repo config files (repos.yaml)."`
	Repo       configtest.RepoCmd       `cmd:"" help:"Validate repo-level config files (atlantis.yaml)."`
}

func main() {
	kctx := kong.Parse(&cli,
		kong.Name("atlantis-configtest"),
		kong.Description("Validate the structure of Atlantis configuration files."),
		kong.Vars{"version": version},
	)

	err := kctx.Run(&configtest.Context{ErrOutput: os.Stderr})
	kctx.FatalIfErrorf(err)
}
