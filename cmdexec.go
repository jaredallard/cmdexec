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
	"testing"
)

// Cmd is an interface to be used instead of [exec.Cmd] for mocking
// purposes.
type Cmd interface {
	Output() ([]byte, error)
	CombinedOutput() ([]byte, error)
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
func UseMockExecutor(t *testing.T, mock *MockExecutor) {
	// Prevent new mock executors from being used until this test has finished.
	if !executorWLock.TryLock() {
		t.Fatal("UseMockExecutor can only be called once per test")
	}

	// Lock the reader to prevent new commands from being created while we
	// swap out the executor.
	executorRLock.Lock()
	originalExecutor := executor
	executor = mock.executor
	executorRLock.Unlock()

	t.Cleanup(func() {
		// Lock the reader again to prevent new commands from being created
		// while we restore the original executor.
		executorRLock.Lock()

		// Unlock the reader and writer once we're done.
		defer executorRLock.Unlock()
		defer executorWLock.Unlock()

		// Restore the original executor.
		executor = originalExecutor
	})
}