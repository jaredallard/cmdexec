package cmdexec_test

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"testing"

	"github.com/jaredallard/cmdexec"
	"gotest.tools/v3/assert"
)

// TestCanMockACommand ensures that if we mock a command, it actually
// gets ran and returns the expected output.
func TestCanMockACommand(t *testing.T) {
	cmdexec.UseMockExecutor(t, cmdexec.NewMockExecutor(&cmdexec.MockCommand{
		Name: "echo",
		Args: []string{"hello", "world"},
		// We use different output here because real 'echo' command would've
		// actually printed the text verbatim.
		Stdout: []byte("hello world ---"),
	}))

	cmd := cmdexec.Command("echo", "hello", "world")
	out, err := cmd.Output()
	assert.NilError(t, err)
	assert.Equal(t, string(out), "hello world ---")
}

// TestCanMockStdin ensures that if we mock a command that reads from
// stdin, it actually reads the expected input and validates it.
func TestCanMockStdin(t *testing.T) {
	mock := cmdexec.NewMockExecutor(&cmdexec.MockCommand{
		Name:  "cat",
		Stdin: []byte("hello world"),
	})

	cmdexec.UseMockExecutor(t, mock)

	// shouldn't fail when setting stdin and expecting it
	cmd := cmdexec.Command("cat")
	cmd.SetStdin(bytes.NewBuffer([]byte("hello world")))
	_, err := cmd.Output()
	assert.NilError(t, err)

	// ensure that it actually validated
	cmd.SetStdin(bytes.NewBuffer([]byte("goodbye world")))
	_, err = cmd.Output()
	assert.Error(t, err, fmt.Sprintf("expected stdin set by SetStdin() to be %q but got %q", "hello world", "goodbye world"))
}

// TestCanReadCombinedOutput ensures that we can read the combined
// output of a command.
func TestCanReadCombinedOutput(t *testing.T) {
	mock := cmdexec.NewMockExecutor(&cmdexec.MockCommand{
		Name:   "echo",
		Args:   []string{"hello", "world"},
		Stdout: []byte("hello"),
		Stderr: []byte("world"),
	})

	cmdexec.UseMockExecutor(t, mock)

	cmd := cmdexec.Command("echo", "hello", "world")
	out, err := cmd.CombinedOutput()
	assert.NilError(t, err)
	assert.Equal(t, string(out), "helloworld")
}

// TestPanicsIfCommandNotRegistered ensures that if we try to run a
// command that hasn't been registered with the mock executor, a panic
// is raised.
func TestPanicsIfCommandNotRegistered(t *testing.T) {
	defer func() {
		r := recover()
		assert.Assert(t, r != nil, "expected a panic to be raised")
		assert.Error(t,
			r.(error),
			"cmdexec: no command registered for 'echo hello' missing call to MockExecutor.AddCommand?",
		)
	}()

	cmdexec.UseMockExecutor(t, cmdexec.NewMockExecutor(&cmdexec.MockCommand{}))

	// run the command. This should panic.
	cmdexec.Command("echo", "hello").Run()
}

// TestStringMatchesExecString ensures that the string representation
// of a command matches that of the standard library's exec.Cmd.
func TestStringMatchesExecString(t *testing.T) {
	echoPath, err := exec.LookPath("echo")
	assert.NilError(t, err)

	mock := cmdexec.NewMockExecutor(&cmdexec.MockCommand{
		Name: "echo",
		Args: []string{"hello", "world"},
	})
	mock.AddLookPath("echo", echoPath)
	cmdexec.UseMockExecutor(t, mock)

	mockStr := cmdexec.Command("echo", "hello", "world").String()
	stdStr := exec.Command("echo", "hello", "world").String()
	assert.Equal(t, mockStr, stdStr)
}

// TestStringFallsBackWithoutLookPath ensures that String() gracefully
// falls back to the raw command name when no LookPath is registered.
func TestStringFallsBackWithoutLookPath(t *testing.T) {
	cmdexec.UseMockExecutor(t, cmdexec.NewMockExecutor(&cmdexec.MockCommand{
		Name: "mycmd",
		Args: []string{"arg1"},
	}))

	str := cmdexec.Command("mycmd", "arg1").String()
	assert.Equal(t, str, "mycmd arg1")
}

func TestMockUnusedDontPanic(_ *testing.T) {
	cmd := &cmdexec.MockCommand{}

	// These are all noops, so just ensure they don't panic.
	cmd.SetEnviron(os.Environ())
	cmd.SetDir("")
	cmd.SetStderr(nil)
	cmd.SetStdout(nil)
	cmd.UseOSStreams(false)
}

// TestCanMockLookPath ensures that a registered LookPath returns the
// expected path.
func TestCanMockLookPath(t *testing.T) {
	mock := cmdexec.NewMockExecutor()
	mock.AddLookPath("git", "/usr/bin/git")
	cmdexec.UseMockExecutor(t, mock)

	path, err := cmdexec.LookPath("git")
	assert.NilError(t, err)
	assert.Equal(t, path, "/usr/bin/git")
}

// TestCanMockLookPathError ensures that a registered LookPath error is
// returned correctly.
func TestCanMockLookPathError(t *testing.T) {
	mock := cmdexec.NewMockExecutor()
	mock.AddLookPathError("nonexistent", exec.ErrNotFound)
	cmdexec.UseMockExecutor(t, mock)

	_, err := cmdexec.LookPath("nonexistent")
	assert.ErrorIs(t, err, exec.ErrNotFound)
}

// TestPanicsIfLookPathNotRegistered ensures that if we try to look up
// a command that hasn't been registered, a panic is raised.
func TestPanicsIfLookPathNotRegistered(t *testing.T) {
	defer func() {
		r := recover()
		assert.Assert(t, r != nil, "expected a panic to be raised")
		assert.Error(t,
			r.(error),
			`cmdexec: no LookPath registered for "unknown" missing call to MockExecutor.AddLookPath?`,
		)
	}()

	cmdexec.UseMockExecutor(t, cmdexec.NewMockExecutor())

	// This should panic.
	cmdexec.LookPath("unknown")
}
