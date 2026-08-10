package configtest

import (
	"github.com/runatlantis/atlantis/server/core/config"
	"github.com/runatlantis/atlantis/server/core/config/valid"
)

// ServerRepoCmd validates a server-side repo config file (repos.yaml).
type ServerRepoCmd struct {
	File string `arg:"" name:"file" help:"Server-side repo config file (repos.yaml)."`
}

func (cmd *ServerRepoCmd) Run(cmdCtx *Context) error {
	parserValidator := &config.ParserValidator{}

	// This is the same call Atlantis makes on startup, in server.NewServer().
	// The defaults it merges in do not affect whether a file is valid.
	defaultCfg := valid.NewGlobalCfgFromArgs(valid.GlobalCfgArgs{})

	return cmdCtx.validate(func() error {
		_, err := parserValidator.ParseGlobalCfg(cmd.File, defaultCfg)

		return err
	})
}
