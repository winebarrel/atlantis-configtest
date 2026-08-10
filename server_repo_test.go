package configtest_test

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	configtest "github.com/winebarrel/atlantis-configtest"
)

func TestServerRepoCmd(t *testing.T) {
	tests := []struct {
		name    string
		files   []string
		errMsg  string
		reports []string
	}{
		{
			name:    "valid",
			files:   []string{"testdata/repos.yaml"},
			reports: []string{"testdata/repos.yaml: ok\n"},
		},
		{
			name:    "unknown key",
			files:   []string{"testdata/repos_unknown_key.yaml"},
			errMsg:  "1 of 1 file(s) invalid",
			reports: []string{"field apply_requirement not found"},
		},
		{
			name:    "unsupported override",
			files:   []string{"testdata/repos_invalid_override.yaml"},
			errMsg:  "1 of 1 file(s) invalid",
			reports: []string{`"nonexistent_key" is not a valid override`},
		},
		{
			name:    "missing file",
			files:   []string{"testdata/nonexistent.yaml"},
			errMsg:  "1 of 1 file(s) invalid",
			reports: []string{"no such file or directory"},
		},
		{
			name:   "every file is checked",
			files:  []string{"testdata/repos_unknown_key.yaml", "testdata/repos.yaml", "testdata/repos_invalid_override.yaml"},
			errMsg: "2 of 3 file(s) invalid",
			reports: []string{
				"testdata/repos_unknown_key.yaml: ",
				"testdata/repos.yaml: ok\n",
				"testdata/repos_invalid_override.yaml: ",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert := assert.New(t)
			out := &bytes.Buffer{}
			cmd := &configtest.ServerRepoCmd{Files: tt.files}

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
