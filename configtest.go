// Package configtest validates the structure of Atlantis configuration files.
//
// Each file is checked on its own. Anything that can only be decided by
// comparing a repo-level file with a server-side one is out of scope.
//
// The rules come from the Atlantis version pinned in go.mod.
package configtest

import (
	"fmt"
	"io"
)

// Context carries the stream commands report on. It is stderr, so that the
// report never mixes into output a caller may be piping somewhere.
type Context struct {
	ErrOutput io.Writer
}

// validate checks every file and reports the result of each one. A failure does
// not stop the run, so one broken file cannot hide the others.
func (cmdCtx *Context) validate(files []string, check func(file string) error) error {
	failed := 0

	for _, file := range files {
		err := check(file)

		if err != nil {
			failed++
			fmt.Fprintf(cmdCtx.ErrOutput, "%s: %s\n", file, err) //nolint:errcheck
		} else {
			fmt.Fprintf(cmdCtx.ErrOutput, "%s: ok\n", file) //nolint:errcheck
		}
	}

	if failed > 0 {
		return fmt.Errorf("%d of %d file(s) invalid", failed, len(files))
	}

	return nil
}
