// Package configtest validates the structure of Atlantis configuration files.
//
// A file is judged on its own. Anything that can only be decided by comparing a
// repo-level file with a server-side one is out of scope.
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

// validate runs the check and reports the result.
func (cmdCtx *Context) validate(check func() error) error {
	err := check()

	if err != nil {
		return err
	}

	fmt.Fprintln(cmdCtx.ErrOutput, "ok") //nolint:errcheck

	return nil
}
