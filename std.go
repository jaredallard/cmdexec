// Copyright (C) 2026 cmdexec contributors
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.
//
// SPDX-License-Identifier: MPL-2.0

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
	//nolint:gosec // Why: acceptable
	return &stdExecutorCmd{exec.CommandContext(ctx, name, arg...)}
}

// String implements [Cmd.String].
func (c *stdExecutorCmd) String() string {
	return c.Cmd.String()
}

// SetEnviron implements [Cmd.SetEnviron].
func (c *stdExecutorCmd) SetEnviron(env []string) {
	c.Env = env
}

// SetDir implements [Cmd.SetDir].
func (c *stdExecutorCmd) SetDir(dir string) {
	c.Dir = dir
}

// SetStdout implements [Cmd.SetStdout].
func (c *stdExecutorCmd) SetStdout(w io.Writer) {
	c.Stdout = w
}

// SetStderr implements [Cmd.SetStderr].
func (c *stdExecutorCmd) SetStderr(w io.Writer) {
	c.Stderr = w
}

// SetStdin implements [Cmd.SetStdin].
func (c *stdExecutorCmd) SetStdin(r io.Reader) {
	c.Stdin = r
}

// UseOSStreams implements [Cmd.UseOSStreams].
func (c *stdExecutorCmd) UseOSStreams(stdin bool) {
	c.SetStdout(os.Stdout)
	c.SetStderr(os.Stderr)
	if stdin {
		c.SetStdin(os.Stdin)
	}
}
