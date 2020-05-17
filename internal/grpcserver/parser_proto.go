package grpcserver

import (
	"sort"

	"github.com/golang/protobuf/ptypes"
	"github.com/mpuzanov/sysmonitor/internal/sysmonitor/domain/model"
	"github.com/mpuzanov/sysmonitor/pkg/sysmonitor/api"
)

func parserLoadDiskToProto(data *model.LoadDisk) (*api.DiskResponse, error) {
	var res api.DiskResponse

	queryTimeProto, err := ptypes.TimestampProto(data.QueryTime)
	if err != nil {
		return &res, nil
	}
	res.QueryTime = queryTimeProto

	for _, v := range data.IO {
		vProtoIO := api.DiskIO{}
		vProtoIO.Device = v.Device
		vProtoIO.Tps = v.Tps
		vProtoIO.KbReadS = v.KbReadS
		vProtoIO.KbWriteS = v.KbWriteS
		vProtoIO.KbRead = v.KbRead
		vProtoIO.KbWrite = v.KbWrite
		res.Io = append(res.Io, &vProtoIO)
	}

	for _, v := range data.FS {
		vProtoFS := api.DiskFS{}
		vProtoFS.FileSystem = v.FileSystem
		vProtoFS.MountedOn = v.MountedOn
		vProtoFS.Used = v.Used
		vProtoFS.Available = v.Available
		vProtoFS.UseProc = v.UseProc
		res.Fs = append(res.Fs, &vProtoFS)
	}
	for _, vProtoFS := range res.Fs {
		vProtoFSInode, ok := data.FSInode[vProtoFS.MountedOn]
		if ok {
			vProtoFS.UsedInode = vProtoFSInode.Used
			vProtoFS.AvailableInode = vProtoFSInode.Available
			vProtoFS.UseProcInode = vProtoFSInode.UseProc
		}
	}
	// сортируем таблицы по убыванию
	sort.Slice(res.Io, func(i, j int) bool {
		return res.Io[i].KbRead > res.Io[j].KbRead
	})
	sort.Slice(res.Fs, func(i, j int) bool {
		return res.Fs[i].Used > res.Fs[j].Used
	})

	return &res, nil
}

func parserTalkerNetToProto(data *model.TalkersNet) (*api.TalkersNetResponse, error) {
	var res api.TalkersNetResponse

	queryTimeProto, err := ptypes.TimestampProto(data.QueryTime)
	if err != nil {
		return &res, nil
	}
	res.QueryTime = queryTimeProto

	for _, v := range data.DevNet {

		vProto := api.DeviceNet{}
		vProto.NetInterface = v.NetInterface
		vProto.ReceiveBytes = int32(v.Receive.Bytes)
		vProto.ReceivePackets = int32(v.Receive.Packets)
		vProto.ReceiveErrs = int32(v.Receive.Errs)
		vProto.TransmitBytes = int32(v.Transmit.Bytes)
		vProto.TransmitPackets = int32(v.Transmit.Packets)
		vProto.TransmitErrs = int32(v.Transmit.Errs)

		res.Devnet = append(res.Devnet, &vProto)
	}

	return &res, nil
}

func parserNetworkStatisticsToProto(data *model.NetworkStatistics) (*api.NetworkStatisticsResponse, error) {
	var res api.NetworkStatisticsResponse

	queryTimeProto, err := ptypes.TimestampProto(data.QueryTime)
	if err != nil {
		return &res, nil
	}
	res.QueryTime = queryTimeProto

	for _, v := range data.StatNet {

		vProto := api.NetStatDetail{}
		vProto.State = v.State
		vProto.Recv = int32(v.Recv)
		vProto.Send = int32(v.Send)
		vProto.LocalAddress = v.LocalAddress
		vProto.PeerAddress = v.PeerAddress

		res.Netstat = append(res.Netstat, &vProto)
	}

	return &res, nil
}
