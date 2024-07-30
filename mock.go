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
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"os/exec"
	"strings"
)

// MockExecutor provides an executor that returns mock data.
type MockExecutor struct {
	// cmd contains the commands that the executor should mock.
	cmds map[string]*MockCommand
}

// MockCommand is a command that can be executed by the MockExecutor.
type MockCommand struct {
	// Name is the name (or path) of the command that should be called to
	// trigger this mock.
	Name string

	// Args is the pair of arguments that the command should be called
	// with to trigger this mock.
	Args []string

	// Stdout is the expected output that the command should write to
	// stdout.
	Stdout []byte

	// Stderr is the expected output that the command should write to
	// stderr.
	Stderr []byte

	// Stdin is the expected input that the command should read from
	// stdin. If this is set, the command will check that the provided
	// stdin matches the expected input. SetStdin() must be called to set
	// the actual stdin data.
	Stdin []byte

	// Err is an error that will be returned when the command is executed.
	// If not set, the command will return nil.
	Err error

	// stdin is a reader that will be used to read from the command's
	// stdin if provided.
	stdin io.Reader
}

// checkStdin checks if the provided stdin matches the expected input.
// This is only checked if both SetStdin() was called on a given command
// and that we expected stdin to be provided.
func (c *MockCommand) checkStdin() error {
	if len(c.Stdin) == 0 {
		return nil
	}

	if c.stdin == nil {
		return fmt.Errorf("expected stdin to be provided but it was not (was SetStdin() called?)")
	}

	got := make([]byte, len(c.Stdin))
	if _, err := c.stdin.Read(got); err != nil {
		return err
	}

	if !bytes.Equal(got, c.Stdin) {
		// Read the rest of the input to provide a more friendly error
		// message.
		rest, err := io.ReadAll(c.stdin)
		if err != nil {
			return fmt.Errorf("failed to read remaining stdin: %w", err)
		}
		got = append(got, rest...)

		return fmt.Errorf("expected stdin set by SetStdin() to be %q but got %q", string(c.Stdin), got)
	}

	return nil
}

// Output implements the [Cmd] interface, see [Cmd.Output] for more
// information.
func (c *MockCommand) Output() ([]byte, error) {
	return c.Stdout, c.Run()
}

// CombinedOutput implements the [Cmd] interface, see
// [Cmd.CombinedOutput] for more information.
func (c *MockCommand) CombinedOutput() ([]byte, error) {
	return append(c.Stdout, c.Stderr...), c.Run()
}

// Run implements the [Cmd] interface, see [Cmd.Run] for more
// information.
func (c *MockCommand) Run() error {
	if err := c.checkStdin(); err != nil {
		return err
	}

	return c.Err
}

// String implements the [Cmd] interface, see [Cmd.String] for more
// information.
func (c *MockCommand) String() string {
	// If possible to look up the command in the PATH, we should return
	// the full path to the command. This is mostly to match the behavior
	// of [exec.Cmd.String].
	execPath := c.Name
	if realPath, err := exec.LookPath(c.Name); err == nil {
		execPath = realPath
	}

	return strings.Join(append([]string{execPath}, c.Args...), " ")
}

// SetEnviron implements the [Cmd] interface. For the MockCommand, this
// is a no-op because we do not actually execute any commands.
func (c *MockCommand) SetEnviron(_ []string) {}

// SetDir implements the [Cmd] interface. For the MockCommand, this is a
// no-op because we do not actually execute any commands.
func (c *MockCommand) SetDir(_ string) {}

// SetStdout implements the [Cmd] interface. For the MockCommand, this
// is a no-op because we do not actually execute any commands.
func (c *MockCommand) SetStdout(_ io.Writer) {}

// SetStderr implements the [Cmd] interface. For the MockCommand, this
// is a no-op because we do not actually execute any commands.
func (c *MockCommand) SetStderr(_ io.Writer) {}

// SetStdin sets the stdin of the command to the given reader. This is
// used for validation purposes to ensure that the provided stdin
// matches what was expected.
func (c *MockCommand) SetStdin(r io.Reader) {
	c.stdin = r
}

// UseOSStreams implements the [Cmd] interface. For the MockCommand,
// this is a no-op because we do not actually execute any commands.
func (c *MockCommand) UseOSStreams(_ bool) {}

// NewMockExecutor returns a new MockExecutor with the given commands. A
// [MockExecutor] contains various commands that should be mocked
// instead of actually executed.
//
// Once a [MockExecutor] is created and used with [UseMockExecutor], it
// will error if any commands are executed that were not registered with
// the executor.
//
// Commands that have had SetStdin() called and set Stdin data will also
// have their stdin checked to ensure that the provided input matches
// the expected input, this enables testing of commands that read from
// stdin.
func NewMockExecutor(cmds ...*MockCommand) *MockExecutor {
	me := &MockExecutor{}
	me.cmds = make(map[string]*MockCommand)
	for _, cmd := range cmds {
		me.AddCommand(cmd)
	}
	return me
}

// getCommandKey returns a unique key for a command based on its name
// and arguments.
func (e *MockExecutor) getCommandKey(name string, args ...string) string {
	return base64.StdEncoding.EncodeToString([]byte(name + " " + strings.Join(args, " ")))
}

// AddCommand adds a command to the executor. If the command has
// already been added, it will be replaced.
//
// Note: This is not thread-safe.
func (e *MockExecutor) AddCommand(cmd *MockCommand) {
	e.cmds[e.getCommandKey(cmd.Name, cmd.Args...)] = cmd
}

// executor implements the [executorFn] type, returning a Cmd based on
// the provided arguments. If no commands are available based on the
// provided input, this function will panic.
func (e *MockExecutor) executor(_ context.Context, name string, arg ...string) Cmd {
	key := e.getCommandKey(name, arg...)
	if cmd, ok := e.cmds[key]; ok {
		return cmd
	}

	panic(
		fmt.Errorf("cmdexec: no command registered for '%s %s' "+
			"missing call to MockExecutor.AddCommand?", name, strings.Join(arg, " "),
		),
	)
}
