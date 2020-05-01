//+build linux

package command_test

import (
	"testing"

	"github.com/mpuzanov/sysmonitor/internal/sysmonitor/command"

	"github.com/stretchr/testify/assert"
)

func TestRunSystemLoad(t *testing.T) {
	exitCode, outresult, outerror := command.RunSystemLoad()
	assert.Equal(t, 0, exitCode)
	assert.NotEmpty(t, outresult)
	assert.Empty(t, outerror)
}

func TestRunSystemCPU(t *testing.T) {
	exitCode, outresult, outerror := command.RunLoadCPU()
	assert.Equal(t, 0, exitCode)
	assert.NotEmpty(t, outresult)
	assert.Empty(t, outerror)
}
