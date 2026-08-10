package configtest_test

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	configtest "github.com/winebarrel/atlantis-configtest"
)

func TestServerRepoCmd(t *testing.T) {
	tests := []struct {
		name   string
		file   string
		errMsg string
	}{
		{
			name: "valid",
			file: "testdata/repos.yaml",
		},
		{
			name:   "unknown key",
			file:   "testdata/repos_unknown_key.yaml",
			errMsg: "field apply_requirement not found",
		},
		{
			name:   "unsupported override",
			file:   "testdata/repos_invalid_override.yaml",
			errMsg: `"nonexistent_key" is not a valid override`,
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
			cmd := &configtest.ServerRepoCmd{File: tt.file}

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
