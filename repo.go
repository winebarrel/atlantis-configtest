package configtest

import (
	"os"

	"github.com/runatlantis/atlantis/server/core/config"
	"github.com/runatlantis/atlantis/server/core/config/valid"
	"gopkg.in/yaml.v3"
)

// RepoCmd validates a repo-level config file (atlantis.yaml).
type RepoCmd struct {
	File string `arg:"" name:"file" help:"Repo-level config file (atlantis.yaml)."`
}

func (cmd *RepoCmd) Run(cmdCtx *Context) error {
	parserValidator := &config.ParserValidator{}

	return cmdCtx.validate(func() error {
		data, err := os.ReadFile(cmd.File)

		if err != nil {
			return err
		}

		_, err = parserValidator.ParseRepoCfgData(data, standaloneGlobalCfg(data), "", "")

		return err
	})
}

// standaloneGlobalCfg builds a server-side config that permits everything a
// repo-level file may contain, so that the file is judged on its own structure.
// Whether a real server permits those settings is a property of its repos.yaml,
// which we deliberately do not look at.
func standaloneGlobalCfg(repoCfgData []byte) valid.GlobalCfg {
	globalCfg := valid.NewGlobalCfgFromArgs(valid.GlobalCfgArgs{AllowAllRepoSettings: true})

	for i := range globalCfg.Repos {
		// AllowAllRepoSettings leaves custom_policy_check out of the list.
		globalCfg.Repos[i].AllowedOverrides = append(globalCfg.Repos[i].AllowedOverrides, valid.CustomPolicyCheckKey)
	}

	// A workflow referenced by name is usually defined in repos.yaml. Register
	// the referenced names so the reference does not look dangling.
	for _, name := range referencedWorkflows(repoCfgData) {
		if _, ok := globalCfg.Workflows[name]; !ok {
			globalCfg.Workflows[name] = valid.Workflow{Name: name}
		}
	}

	return globalCfg
}

// referencedWorkflows returns the workflow names the projects refer to. It is a
// loose scan: anything malformed is left for the real parser to report.
func referencedWorkflows(repoCfgData []byte) []string {
	var repoCfg struct {
		Projects []struct {
			Workflow *string `yaml:"workflow"`
		} `yaml:"projects"`
	}

	if err := yaml.Unmarshal(repoCfgData, &repoCfg); err != nil {
		return nil
	}

	names := []string{}

	for _, project := range repoCfg.Projects {
		if project.Workflow != nil {
			names = append(names, *project.Workflow)
		}
	}

	return names
}
