package configtest_test

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	configtest "github.com/winebarrel/atlantis-configtest"
)

func TestRepoCmd(t *testing.T) {
	tests := []struct {
		name    string
		files   []string
		errMsg  string
		reports []string
	}{
		{
			// The file names a workflow defined in repos.yaml and sets keys that
			// need allowed_overrides. None of that can be judged from this file
			// alone, so it must pass.
			name:    "settings that depend on the server side",
			files:   []string{"testdata/atlantis.yaml"},
			reports: []string{"testdata/atlantis.yaml: ok\n"},
		},
		{
			name:    "unknown key",
			files:   []string{"testdata/atlantis_unknown_key.yaml"},
			errMsg:  "1 of 1 file(s) invalid",
			reports: []string{"field directory not found"},
		},
		{
			name:    "duplicate project",
			files:   []string{"testdata/atlantis_duplicate_project.yaml"},
			errMsg:  "1 of 1 file(s) invalid",
			reports: []string{"there are two or more projects with dir"},
		},
		{
			name:    "broken yaml",
			files:   []string{"testdata/atlantis_broken.yaml"},
			errMsg:  "1 of 1 file(s) invalid",
			reports: []string{"did not find expected key"},
		},
		{
			name:    "missing file",
			files:   []string{"testdata/nonexistent.yaml"},
			errMsg:  "1 of 1 file(s) invalid",
			reports: []string{"no such file or directory"},
		},
		{
			name:   "every file is checked",
			files:  []string{"testdata/atlantis_unknown_key.yaml", "testdata/atlantis.yaml", "testdata/atlantis_broken.yaml"},
			errMsg: "2 of 3 file(s) invalid",
			reports: []string{
				"testdata/atlantis_unknown_key.yaml: ",
				"testdata/atlantis.yaml: ok\n",
				"testdata/atlantis_broken.yaml: ",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert := assert.New(t)
			out := &bytes.Buffer{}
			cmd := &configtest.RepoCmd{Files: tt.files}

			err := cmd.Run(&configtest.Context{ErrOutput: out})

			if tt.errMsg == "" {
				assert.NoError(err)
			} else {
				assert.EqualError(err, tt.errMsg)
			}

			for _, report := range tt.reports {
				assert.Contains(out.String(), report)
			}
		})
	}
}
