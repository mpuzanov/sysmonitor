package sysmonitor_test

import (
	"testing"

	"github.com/mpuzanov/sysmonitor/internal/sysmonitor"
	"github.com/stretchr/testify/assert"
)

func TestQueryInfoDisk(t *testing.T) {
	got, err := sysmonitor.QueryInfoDisk()
	assert.Empty(t, err)
	assert.NotEqual(t, 0, len(got.IO))
	assert.NotEqual(t, 0, len(got.FS))
	assert.NotEqual(t, 0, len(got.FSInode))
}
