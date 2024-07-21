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
	"fmt"
)

// MockExecutor provides an executor that returns mock data.
type MockExecutor struct {
	cmds []*MockCommand
}

// MockCommand is a command that can be executed by the MockExecutor.
type MockCommand struct {
	Name   string
	Args   []string
	Stdout []byte
	Stderr []byte
	Err    error
}

func (c *MockCommand) Output() ([]byte, error) {
	return c.Stdout, c.Err
}

func (c *MockCommand) CombinedOutput() ([]byte, error) {
	return append(c.Stdout, c.Stderr...), c.Err
}

// NewMockExecutor returns a new MockExecutor with the given commands.
func NewMockExecutor(cmds ...*MockCommand) *MockExecutor {
	return &MockExecutor{cmds}
}

// AddCommand adds a command to the executor.
//
// Note: This is not thread-safe.
func (e *MockExecutor) AddCommand(cmd *MockCommand) {
	e.cmds = append(e.cmds, cmd)
}

// executor implements the [executorFn] type, returning a Cmd based on
// the provided arguments. If no commands are available based on the
// provided input, this function will panic.
func (e *MockExecutor) executor(_ context.Context, name string, arg ...string) Cmd {
	if len(e.cmds) == 0 {
		panic("no commands to execute")
	}

	// argsEqual checks if two slices of strings are equal.
	var argsEqual = func(a, b []string) bool {
		for i := range a {
			if a[i] != b[i] {
				return false
			}
		}

		return true
	}

	// Check if we have a command that matches the input name and args.
	for i := range e.cmds {
		cmd := e.cmds[i]
		if cmd.Name == name && argsEqual(cmd.Args, arg) {
			return cmd
		}
	}

	// Did you forget to call [AddCommand]?
	panic(fmt.Errorf("no mocked output registered for %s %v", name, arg))
}
