//+build linux

package command_test

import (
	"testing"

	"github.com/mpuzanov/sysmonitor/internal/sysmonitor/command"

	"github.com/stretchr/testify/assert"
)

func TestRunRunSystemLoad(t *testing.T) {
	exitCode, outresult, outerror := command.RunSystemLoad()
	assert.Equal(t, 0, exitCode)
	assert.NotEmpty(t, outresult)
	assert.Empty(t, outerror)
}
