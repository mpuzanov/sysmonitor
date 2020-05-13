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
		{NetInterface: "enp0s8", Receive: model.DevNetDetail{Bytes: 12483477, Packets: 15144, Errs: 0},
			Transmit: model.DevNetDetail{Bytes: 886761, Packets: 7317, Errs: 0}},
		{NetInterface: "lo", Receive: model.DevNetDetail{Bytes: 1490304, Packets: 1553, Errs: 0},
			Transmit: model.DevNetDetail{Bytes: 1490304, Packets: 1553, Errs: 0}},
		{NetInterface: "enp0s3", Receive: model.DevNetDetail{Bytes: 36811, Packets: 322, Errs: 0},
			Transmit: model.DevNetDetail{Bytes: 12964, Packets: 96, Errs: 0}},
	}

	testNetwork = `
	State         Recv-Q     Send-Q          Local Address:Port              Peer Address:Port 
	LISTEN        0          128                   0.0.0.0:http                   0.0.0.0:*                
	LISTEN        0          5                   127.0.0.1:ipp                    0.0.0.0:*         
	CLOSE-WAIT    734        0                   10.0.3.15:47704            68.232.34.200:https     
	ESTAB         0          0                   10.0.3.15:45116           51.144.164.215:https     
	LISTEN        0          128                      [::]:http                      [::]:*         
	LISTEN        0          5                       [::1]:ipp                       [::]:*
`
	wantNetwork = []model.NetStatDetail{
		{State: "LISTEN", Recv: 0, Send: 128, LocalAddress: "0.0.0.0:http", PeerAddress: "0.0.0.0:*"},
		{State: "LISTEN", Recv: 0, Send: 5, LocalAddress: "127.0.0.1:ipp", PeerAddress: "0.0.0.0:*"},
		{State: "CLOSE-WAIT", Recv: 734, Send: 0, LocalAddress: "10.0.3.15:47704", PeerAddress: "68.232.34.200:https"},
		{State: "ESTAB", Recv: 0, Send: 0, LocalAddress: "10.0.3.15:45116", PeerAddress: "51.144.164.215:https"},
		{State: "LISTEN", Recv: 0, Send: 128, LocalAddress: "[::]:http", PeerAddress: "[::]:*"},
		{State: "LISTEN", Recv: 0, Send: 5, LocalAddress: "[::1]:ipp", PeerAddress: "[::]:*"},
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
				assert.Equal(t, tC.want[0].Receive.Packets, 15144)
			}
		})
	}
}

func TestParserNetworkStatistics(t *testing.T) {
	testCases := []struct {
		desc string
		in   string
		want []model.NetStatDetail
		err  error
	}{
		{
			desc: "Test 1",
			in:   testNetwork,
			want: wantNetwork,
			err:  nil,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			got, err := parser.ParserNetworkStatistics(tC.in)
			assert.Equal(t, tC.err, err)
			if err == nil {
				assert.Equal(t, len(tC.want), len(got))
				assert.Equal(t, tC.want[0].Send, 128)
			}
		})
	}
}
