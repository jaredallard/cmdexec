// Copyright (C) 2026 cmdexec contributors
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as
// published by the Free Software Foundation, either version 3 of the
// License, or (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public
// License along with this program. If not, see
// <https://www.gnu.org/licenses/>.
//
// SPDX-License-Identifier: LGPL-3.0

// Package cmdexec provides a way to execute commands using the exec
// package while supporting mocking for testing purposes. The default
// behaviour of the package is to simply wrap [exec.Command] and it's
// context accepting counterpart, [exec.CommandContext]. However, when
// running in tests, the package can be configured to use a mock
// executor that allows for controlling the output and behaviour of the
// commands executed for testing purposes.
package cmdexec

import (
	"context"
	"io"

	"github.com/jaredallard/cmdexec/internal/mockt"
)

// Cmd is an interface to be used instead of [exec.Cmd] for mocking
// purposes.
type Cmd interface {
	// Output matches [exec.Cmd.Output].
	Output() ([]byte, error)
	// CombinedOutput matches [exec.Cmd.CombinedOutput].
	CombinedOutput() ([]byte, error)
	// Run matches [exec.Cmd.Run].
	Run() error
	// String returns the command line string that will be executed.
	String() string

	// Below are non-standard functions (no present in the [exec.Cmd])
	// that are provided for convenience.

	// SetEnviron sets the environment variables of the command. Matches
	// the behavior of setting [exec.Cmd.Environ] directly.
	SetEnviron([]string)

	// SetDir sets the working directory of the command.
	SetDir(string)

	// SetStdout, SetStderr, and SetStdin set the stdout, stderr, and
	// stdin of the command respectively.
	SetStdout(io.Writer)
	SetStderr(io.Writer)
	SetStdin(io.Reader)

	// UseOSStreams sets Stdout, Stderr, and Stdin to the OS streams
	// (os.Stdout, os.Stderr, and os.Stdin respectively). If stdin is
	// false then Stdin is not set.
	UseOSStreams(stdin bool)
}

// Command returns a new Cmd that will call the given command with the
// given arguments. See [exec.Command] for more information.
func Command(name string, arg ...string) Cmd {
	return CommandContext(context.Background(), name, arg...)
}

// CommandContext returns a new Cmd that will call the given command with
// the given arguments and the given context. See [exec.CommandContext]
// for more information.
func CommandContext(ctx context.Context, name string, arg ...string) Cmd {
	executorRLock.Lock()
	defer executorRLock.Unlock()

	return executor(ctx, name, arg...)
}

// LookPath searches for an executable named file in the directories
// named by the PATH environment variable. See [exec.LookPath] for more
// information.
func LookPath(file string) (string, error) {
	executorRLock.Lock()
	defer executorRLock.Unlock()

	return lookPath(file)
}

// UseMockExecutor replaces the executor used by cmdexec with a mock
// executor that can be used to control the output of all commands
// created after this function is called. A cleanup function is added
// to the test to ensure that the original executor is restored after
// the test has finished.
//
// Note: This function can only ever be called once per test. If called
// again in the same test, it will cause the test to fail.
//
// Usage:
//
//	func TestSomething(t *testing.T) {
//	    mock := cmdexec.NewMockExecutor()
//	    mock.AddCommand(&cmdexec.MockCommand{
//	        Name:   "echo",
//	        Args:   []string{"hello", "world"},
//	        Stdout: []byte("hello world\n"),
//	    })
//
//	    cmdexec.UseMockExecutor(t, mock)
//
//	    // Your test code here.
//	}
func UseMockExecutor(t mockt.T, mock *MockExecutor) {
	// Prevent new mock executors from being used until this test has finished.
	if !executorWLock.TryLock() {
		t.Fatal("UseMockExecutor can only be called once per test")
		return
	}

	// Lock the reader to prevent new commands from being created while we
	// swap out the executor.
	executorRLock.Lock()
	originalExecutor := executor
	originalLookPath := lookPath
	executor = mock.executor
	lookPath = mock.lookPath
	executorRLock.Unlock()

	t.Cleanup(func() {
		// Lock the reader again to prevent new commands from being created
		// while we restore the original executor.
		executorRLock.Lock()

		// Unlock the reader and writer once we're done.
		defer executorRLock.Unlock()
		defer executorWLock.Unlock()

		// Restore the original executor and lookPath.
		executor = originalExecutor
		lookPath = originalLookPath
	})
}
