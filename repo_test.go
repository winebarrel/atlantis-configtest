package configtest_test

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	configtest "github.com/winebarrel/atlantis-configtest"
)

func TestRepoCmd(t *testing.T) {
	tests := []struct {
		name   string
		file   string
		errMsg string
	}{
		{
			// The file names a workflow defined in repos.yaml and sets keys that
			// need allowed_overrides. None of that can be judged from this file
			// alone, so it must pass.
			name: "settings that depend on the server side",
			file: "testdata/atlantis.yaml",
		},
		{
			name:   "unknown key",
			file:   "testdata/atlantis_unknown_key.yaml",
			errMsg: "field directory not found",
		},
		{
			name:   "duplicate project",
			file:   "testdata/atlantis_duplicate_project.yaml",
			errMsg: "there are two or more projects with dir",
		},
		{
			name:   "broken yaml",
			file:   "testdata/atlantis_broken.yaml",
			errMsg: "did not find expected key",
		},
		{
			name:   "missing file",
			file:   "testdata/nonexistent.yaml",
			errMsg: "no such file or directory",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert := assert.New(t)
			out := &bytes.Buffer{}
			cmd := &configtest.RepoCmd{File: tt.file}

			err := cmd.Run(&configtest.Context{ErrOutput: out})

			if tt.errMsg == "" {
				assert.NoError(err)
				assert.Equal("ok\n", out.String())
			} else {
				assert.ErrorContains(err, tt.errMsg)
				assert.Empty(out.String())
			}
		})
	}
}
