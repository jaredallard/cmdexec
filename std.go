// Copyright (C) 2024 Jared Allard <jaredallard@users.noreply.github.com>
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program. If not, see <https://www.gnu.org/licenses/>.

package cmdexec

import (
	"context"
	"io"
	"os"
	"os/exec"
)

// stdExecutorCmd is a simple wrapper around [exec.Cmd] to implement the
// [Cmd] interface.
//
// All functions on this struct are not thread-safe.
type stdExecutorCmd struct {
	*exec.Cmd
}

// stdExecutor creates a new [Cmd] using [exec.CommandContext] as the
// underlying executor.
func stdExecutor(ctx context.Context, name string, arg ...string) Cmd {
	return &stdExecutorCmd{exec.CommandContext(ctx, name, arg...)}
}

// String implements [Cmd.String].
func (c *stdExecutorCmd) String() string {
	return c.Cmd.String()
}

// SetEnviron implements [Cmd.SetEnviron].
func (c *stdExecutorCmd) SetEnviron(env []string) {
	c.Cmd.Env = env
}

// SetDir implements [Cmd.SetDir].
func (c *stdExecutorCmd) SetDir(dir string) {
	c.Cmd.Dir = dir
}

// SetStdout implements [Cmd.SetStdout].
func (c *stdExecutorCmd) SetStdout(w io.Writer) {
	c.Cmd.Stdout = w
}

// SetStderr implements [Cmd.SetStderr].
func (c *stdExecutorCmd) SetStderr(w io.Writer) {
	c.Cmd.Stderr = w
}

// SetStdin implements [Cmd.SetStdin].
func (c *stdExecutorCmd) SetStdin(r io.Reader) {
	c.Cmd.Stdin = r
}

// UseOSStreams implements [Cmd.UseOSStreams].
func (c *stdExecutorCmd) UseOSStreams(stdin bool) {
	c.SetStdout(os.Stdout)
	c.SetStderr(os.Stderr)
	if stdin {
		c.SetStdin(os.Stdin)
	}
}
