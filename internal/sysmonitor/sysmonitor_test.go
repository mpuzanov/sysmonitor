package sysmonitor_test

import (
	"testing"

	"github.com/mpuzanov/sysmonitor/internal/sysmonitor"
	"github.com/stretchr/testify/assert"
)

func TestQueryInfoSystem(t *testing.T) {
	got, err := sysmonitor.QueryInfoSystem()
	assert.Empty(t, err)
	assert.NotEmpty(t, got.SystemLoadValue)
}

func TestQueryInfoCPU(t *testing.T) {
	got, err := sysmonitor.QueryInfoCPU()
	assert.Empty(t, err)
	assert.NotEmpty(t, got)
}
func TestQueryInfoDisk(t *testing.T) {
	got, err := sysmonitor.QueryInfoDisk()
	assert.Empty(t, err)
	assert.NotEqual(t, 0, len(got.IO))
	assert.NotEqual(t, 0, len(got.FS))
	assert.NotEqual(t, 0, len(got.FSInode))
}

func TestQueryInfoDeviceNet(t *testing.T) {
	got, err := sysmonitor.QueryInfoTalkersNet()
	assert.Empty(t, err)
	assert.NotEqual(t, 0, len(got.DevNet))
}
