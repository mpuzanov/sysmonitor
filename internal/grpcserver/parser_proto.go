package grpcserver

import (
	"github.com/golang/protobuf/ptypes"
	"github.com/mpuzanov/sysmonitor/internal/sysmonitor/domain/model"
	"github.com/mpuzanov/sysmonitor/pkg/sysmonitor/api"
)

// ParserLoadDiskToProto .
func ParserLoadDiskToProto(data *model.LoadDisk) (*api.DiskResponse, error) {
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

	return &res, nil
}
