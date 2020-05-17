package sysmonitor

import (
	"fmt"
	"time"

	"github.com/mpuzanov/sysmonitor/internal/sysmonitor/command"
	"github.com/mpuzanov/sysmonitor/internal/sysmonitor/domain/errors"
	"github.com/mpuzanov/sysmonitor/internal/sysmonitor/domain/model"
	"github.com/mpuzanov/sysmonitor/internal/sysmonitor/parser"
)

// QueryInfoSystem Получение информации из системы по LoadSystem
func QueryInfoSystem() (model.LoadSystem, error) {
	var res model.LoadSystem
	exitCode, txt, outerror := command.RunSystemLoad()
	if exitCode != 0 {
		return res, fmt.Errorf("%s. %s", outerror, errors.ErrRunReadInfoSystem)
	}
	res, err := parser.SystemLoad(txt)
	if err != nil {
		return res, fmt.Errorf("%s. %w", errors.ErrParserReadInfoSystem, err)
	}
	return res, nil
}

// QueryInfoCPU Получение информации из системы по LoadCPU
func QueryInfoCPU() (model.LoadCPU, error) {
	var res model.LoadCPU
	exitCode, txt, outerror := command.RunLoadCPU()
	if exitCode != 0 {
		return res, fmt.Errorf("%s. %s", outerror, errors.ErrRunReadInfoCPU)
	}
	res, err := parser.LoadCPU(txt)
	if err != nil {
		return res, fmt.Errorf("%s. %w", errors.ErrParserReadInfoCPU, err)
	}
	return res, nil
}

// QueryInfoDisk Получение информации по дисковой системе
func QueryInfoDisk() (model.LoadDisk, error) {
	var res model.LoadDisk

	exitCode, txt, outerror := command.RunCommand("iostat", "-d", "-k")
	if exitCode != 0 {
		return res, fmt.Errorf("%s. %s", outerror, errors.ErrRunLoadDiskDevice)
	}
	resIO, err := parser.LoadDiskDevice(txt)
	if err != nil {
		return res, fmt.Errorf("%s. %w", errors.ErrParserLoadDiskDevice, err)
	}

	exitCode, txt, outerror = command.RunCommand("df", "-k")
	if exitCode != 0 {
		return res, fmt.Errorf("%s. %s", outerror, errors.ErrRunLoadDiskFS)
	}
	resFS, err := parser.LoadDiskFS(txt)
	if err != nil {
		return res, fmt.Errorf("%s. %w", errors.ErrParserLoadDiskFS, err)
	}
	exitCode, txt, outerror = command.RunCommand("df", "-i")
	if exitCode != 0 {
		return res, fmt.Errorf("%s. %s", outerror, errors.ErrRunLoadDiskFSInode)
	}
	resFSInode, err := parser.LoadDiskFS(txt)
	if err != nil {
		return res, fmt.Errorf("%s. %w", errors.ErrParserLoadDiskFSInode, err)
	}
	res.IO = resIO
	res.FS = resFS
	res.FSInode = resFSInode
	res.QueryTime = time.Now()

	return res, nil
}

// QueryInfoTalkersNet Получение информации по трафику сети
func QueryInfoTalkersNet() (model.TalkersNet, error) {
	var res model.TalkersNet

	exitCode, txt, outerror := command.RunCommand("cat", "/proc/net/dev")
	if exitCode != 0 {
		return res, fmt.Errorf("%s. %s", outerror, errors.ErrRunDeviceNet)
	}
	resDevNet, err := parser.DeviceNet(txt)
	if err != nil {
		return res, fmt.Errorf("%s. %w", errors.ErrParserDeviceNet, err)
	}

	res.DevNet = resDevNet
	res.QueryTime = time.Now()

	return res, nil
}

// QueryInfoNetworkStatistics Получение информации по статистики сети
func QueryInfoNetworkStatistics() (model.NetworkStatistics, error) {
	var res model.NetworkStatistics

	exitCode, txt, outerror := command.RunCommand("ss", "-ta")
	if exitCode != 0 {
		return res, fmt.Errorf("%s. %s", outerror, errors.ErrRunNetworkStatistics)
	}
	resNet, err := parser.NetworkStatistics(txt)
	if err != nil {
		return res, fmt.Errorf("%s. %w", errors.ErrParserNetworkStatistics, err)
	}

	res.StatNet = resNet
	res.QueryTime = time.Now()

	return res, nil
}
