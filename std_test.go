//go:build !windows

package cmdexec_test

import (
	"testing"

	"github.com/jaredallard/cmdexec"
	"gotest.tools/v3/assert"
)

// TestCanExecuteACommand ensures that the default executor is the
// standard executor and that it can execute a command.
func TestCanExecuteACommand(t *testing.T) {
	cmd := cmdexec.Command("echo", "hello")
	out, err := cmd.Output()
	assert.NilError(t, err)
	assert.Equal(t, string(out), "hello\n")
}

func ExampleCmd_UseOSStreams() {
	// This example demonstrates how to use the UseOSStreams function to
	// set the stdin, stdout, and stderr of a command to the OS streams.
	cmd := cmdexec.Command("echo", "hello")
	cmd.UseOSStreams(false)
	if err := cmd.Run(); err != nil {
		panic(err)
	}

	// Output:
	// hello
}