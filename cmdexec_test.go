package cmdexec_test

import (
	"testing"

	"github.com/jaredallard/cmdexec"
	"github.com/jaredallard/cmdexec/internal/mockt"
	"gotest.tools/v3/assert"
)

func TestTestIsFailedIfMultipleExecutors(t *testing.T) {
	subT := mockt.New()
	t.Cleanup(subT.RunCleanup) // ensure we don't cause other tests to fail.

	me := cmdexec.NewMockExecutor(&cmdexec.MockCommand{})
	cmdexec.UseMockExecutor(subT, me)
	cmdexec.UseMockExecutor(subT, me) // should fail

	assert.Equal(t, subT.Failed(), true, "expected sub-test to fail")
}
