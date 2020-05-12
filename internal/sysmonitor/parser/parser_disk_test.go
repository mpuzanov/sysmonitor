//+build linux darwin

package parser_test

import (
	"testing"

	"github.com/mpuzanov/sysmonitor/internal/sysmonitor/domain/model"
	"github.com/mpuzanov/sysmonitor/internal/sysmonitor/parser"
	"github.com/stretchr/testify/assert"
)

var (
	testIostat = `
Linux 5.3.0-51-generic (adm2) 	07.05.2020 	_x86_64_	(2 CPU)

Device             tps    kB_read/s    kB_wrtn/s    kB_read    kB_wrtn
loop0             3,39         3,76         0,00      10473          0
sda              31,72       597,45       108,62    1665365     302760
sdb               0,35         1,68        21,00       4688      58536
`
	wantIostat = []model.DiskIO{
		{Device: "loop0", Tps: 3.39, KbReadS: 3.76, KbWriteS: 0.00, KbRead: 10473, KbWrite: 0},
		{Device: "sda", Tps: 31.72, KbReadS: 597.45, KbWriteS: 108.62, KbRead: 1665365, KbWrite: 302760},
		{Device: "sdb", Tps: 0.35, KbReadS: 1.68, KbWriteS: 21.00, KbRead: 10473, KbWrite: 58536},
	}

	testFS = `
Filesystem     1K-blocks      Used Available Use% Mounted on
/dev/sda1       32121664  15114024  15470980  50% /
tmpfs            1017628     20396    997232   3% /dev/shm
tmpfs               5120         4      5116   1% /run/lock
tmpfs             203524        12    203512   1% /run/user/1000
`
	wantFS = map[string]model.DiskFS{
		"/":              {FileSystem: "/dev/sda1", MountedOn: "/", Used: 15114024, Available: 15470980, UseProc: "50%"},
		"/dev/shm":       {FileSystem: "tmpfs", MountedOn: "/dev/shm", Used: 20396, Available: 997232, UseProc: "3%"},
		"/run/lock":      {FileSystem: "tmpfs", MountedOn: "/run/lock", Used: 4, Available: 5116, UseProc: "1%"},
		"/run/user/1000": {FileSystem: "tmpfs", MountedOn: "/run/user/1000", Used: 12, Available: 203512, UseProc: "1%"},
	}

	testFSInode = `
Filesystem       Inodes   IUsed    IFree IUse% Mounted on
/dev/sda1       2048000  512210  1535790   26% /
tmpfs            254407      63   254344    1% /dev/shm
tmpfs            254407       5   254402    1% /run/lock
tmpfs            254407      29   254378    1% /run/user/1000
`
	wantFSInode = map[string]model.DiskFS{
		"/":              {FileSystem: "/dev/sda1", MountedOn: "/", Used: 512210, Available: 1535790, UseProc: "26%"},
		"/dev/shm":       {FileSystem: "tmpfs", MountedOn: "/dev/shm", Used: 63, Available: 254344, UseProc: "1%"},
		"/run/lock":      {FileSystem: "tmpfs", MountedOn: "/run/lock", Used: 5, Available: 254402, UseProc: "1%"},
		"/run/user/1000": {FileSystem: "tmpfs", MountedOn: "/run/user/1000", Used: 29, Available: 254378, UseProc: "1%"},
	}
)

func TestParserLoadDiskDevice(t *testing.T) {
	testCases := []struct {
		desc string
		in   string
		want []model.DiskIO
		err  error
	}{
		{
			desc: "Test 1",
			in:   testIostat,
			want: wantIostat,
			err:  nil,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			got, err := parser.ParserLoadDiskDevice(tC.in)
			assert.Equal(t, tC.err, err)
			if err == nil {
				assert.Equal(t, len(tC.want), len(got))
			}
		})
	}
}

func TestParserLoadDiskFS(t *testing.T) {
	testCases := []struct {
		desc string
		in   string
		want map[string]model.DiskFS
		err  error
	}{
		{
			desc: "Test 1",
			in:   testFS,
			want: wantFS,
			err:  nil,
		},
		{
			desc: "Test 2 FSInode",
			in:   testFSInode,
			want: wantFSInode,
			err:  nil,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			got, err := parser.ParserLoadDiskFS(tC.in)
			assert.Equal(t, tC.err, err)
			if err == nil {
				assert.Equal(t, len(tC.want), len(got))
			}
		})
	}
}
