//+build linux darwin

package parser_test

import (
	"testing"

	"github.com/mpuzanov/sysmonitor/internal/sysmonitor/domain/model"
	"github.com/mpuzanov/sysmonitor/internal/sysmonitor/parser"
	"github.com/stretchr/testify/assert"
)

var (
	testIn = `
	Inter-|   Receive                                                |  Transmit
	face |bytes    packets errs drop fifo frame compressed multicast|bytes    packets errs drop fifo colls carrier compressed
   enp0s8: 12483477   15144    0    0    0     0          0         0   886761    7317    0    0    0     0       0          0
	   lo: 1490304    1553    0    0    0     0          0         0  1490304    1553    0    0    0     0       0          0
   enp0s3:   36811     322    0    0    0     0          0       161    12964      96    0    0    0     0       0          0   

`
	wantInfo = []model.DeviceNet{
		{NetInterface: "enp0s8", Receive: model.DevNetStat{Bytes: 12483477, Packets: 15144, Errs: 0},
			Transmit: model.DevNetStat{Bytes: 886761, Packets: 7317, Errs: 0}},
		{NetInterface: "lo", Receive: model.DevNetStat{Bytes: 1490304, Packets: 1553, Errs: 0},
			Transmit: model.DevNetStat{Bytes: 1490304, Packets: 1553, Errs: 0}},
		{NetInterface: "enp0s3", Receive: model.DevNetStat{Bytes: 36811, Packets: 322, Errs: 0},
			Transmit: model.DevNetStat{Bytes: 12964, Packets: 96, Errs: 0}},
	}
)

func TestParserDeviceNet(t *testing.T) {
	testCases := []struct {
		desc string
		in   string
		want []model.DeviceNet
		err  error
	}{
		{
			desc: "Test 1",
			in:   testIn,
			want: wantInfo,
			err:  nil,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			got, err := parser.ParserDeviceNet(tC.in)
			assert.Equal(t, tC.err, err)
			if err == nil {
				assert.Equal(t, len(tC.want), len(got))
			}
		})
	}
}
