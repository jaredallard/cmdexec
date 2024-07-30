//go:build !windows

package cmdexec_test

import (
	"os/exec"
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

func Test_stdExecutorString(t *testing.T) {
	execPath, err := exec.LookPath("echo")
	assert.NilError(t, err)

	// /bin/echo hello (usually)
	assert.Equal(t, cmdexec.Command("echo", "hello").String(), execPath+" hello")
}

func Test_stdExecutorSetEnviron(t *testing.T) {
	cmd := cmdexec.Command("printenv", "STENCIL_TEST_ENV")
	cmd.SetEnviron([]string{"STENCIL_TEST_ENV=1"})

	out, err := cmd.Output()
	assert.NilError(t, err)
	assert.Equal(t, string(out), "1\n")
}
