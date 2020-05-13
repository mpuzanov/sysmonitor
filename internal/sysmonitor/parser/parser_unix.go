// +build linux darwin

package parser

import (
	"bufio"
	"regexp"
	"strconv"
	"strings"

	"github.com/mpuzanov/sysmonitor/internal/sysmonitor/domain/errors"
	"github.com/mpuzanov/sysmonitor/internal/sysmonitor/domain/model"
)

const (
	patternLoadCPU = `(\d*[.,]?\d+) us\,\s+(\d*[.,]?\d+) sy\,\s+(\d*[.,]?\d+) ni\,\s+(\d*[.,]?\d+) id`
)

var (
	regexpLoadCPU *regexp.Regexp
)

func init() {
	regexpLoadCPU = regexp.MustCompile(patternLoadCPU)
}

// ParserSystemLoad Выдаёт значение загрузки системы из строки с общей информацией
func ParserSystemLoad(in string) (model.LoadSystem, error) {
	//top - 19:30:09 up  4:30,  1 user,  load average: 1,02, 0,95, 0,80
	var res model.LoadSystem
	idx := strings.LastIndex(in, ":")
	in = strings.TrimLeft(in[idx+1:], " ") // 1,02, 0,95, 0,80
	arr := strings.Split(in, ", ")
	in = strings.Replace(arr[0], ",", ".", 1)
	val, err := strconv.ParseFloat(in, 64)
	if err != nil {
		return res, err
	}
	res.SystemLoadValue = val

	return res, nil
}

// ParserLoadCPU Выдаёт значение загрузки системы из строки с общей информацией
func ParserLoadCPU(in string) (model.LoadCPU, error) {
	//%Cpu(s):  1,3 us,  0,6 sy,  0,0 ni, 96,7 id,  1,3 wa,  0,0 hi,  0,1 si,  0,0 st
	var res model.LoadCPU
	var err error

	matchall := regexpLoadCPU.FindAllStringSubmatch(in, -1)
	if len(matchall) == 0 {
		return res, errors.ErrParserReadInfoCPU
	}

	for _, elements := range matchall {

		for key, elem := range elements {
			elem = strings.Replace(elem, ",", ".", 1)
			switch key {
			case 1:
				{
					res.UserMode, err = strconv.ParseFloat(elem, 64)
					if err != nil {
						return res, err
					}
				}
			case 2:
				{
					res.SystemMode, err = strconv.ParseFloat(elem, 64)
					if err != nil {
						return res, err
					}
				}
			case 4:
				{
					res.Idle, err = strconv.ParseFloat(elem, 64)
					if err != nil {
						return res, err
					}
				}
			}
		}
	}
	return res, nil
}

// ParserLoadDiskDevice анализируем результаты команды iostat -d -k
func ParserLoadDiskDevice(in string) ([]model.DiskIO, error) {
	var (
		res = []model.DiskIO{}
		v   model.DiskIO
		err error
		i   int
	)

	scanner := bufio.NewScanner(strings.NewReader(in))
	for scanner.Scan() {
		// ищем шапку
		if !strings.HasPrefix(scanner.Text(), "Device") {
			continue
		}
		// нашли, далее читаем данные
		for scanner.Scan() {
			s := strings.Replace(scanner.Text(), ",", ".", -1)
			data := strings.Fields(s)
			if len(data) < 6 {
				continue
			}
			v.Device = data[0]
			v.Tps, err = strconv.ParseFloat(data[1], 64)
			if err != nil {
				return res, err
			}
			v.KbReadS, err = strconv.ParseFloat(data[2], 64)
			if err != nil {
				return res, err
			}
			v.KbWriteS, err = strconv.ParseFloat(data[3], 64)
			if err != nil {
				return res, err
			}
			i, err = strconv.Atoi(data[4])
			if err != nil {
				return res, err
			}
			v.KbRead = int32(i)
			i, err = strconv.Atoi(data[5])
			if err != nil {
				return res, err
			}
			v.KbWrite = int32(i)

			res = append(res, v)
		}
	}
	return res, nil
}

// ParserLoadDiskFS анализируем результаты команды df -k   df -i
func ParserLoadDiskFS(in string) (map[string]model.DiskFS, error) {
	var (
		res = map[string]model.DiskFS{}
		v   model.DiskFS
		err error
		i   int
	)

	scanner := bufio.NewScanner(strings.NewReader(in))
	for scanner.Scan() {
		// ищем шапку
		if !strings.HasPrefix(scanner.Text(), "Filesystem") {
			continue
		}
		// нашли, далее читаем данные
		for scanner.Scan() {
			s := strings.Replace(scanner.Text(), ",", ".", -1)
			data := strings.Fields(s)
			if len(data) < 6 {
				continue
			}
			v.FileSystem = data[0]
			i, err = strconv.Atoi(data[2])
			if err != nil {
				return res, err
			}
			v.Used = int32(i)
			i, err = strconv.Atoi(data[3])
			if err != nil {
				return res, err
			}
			v.Available = int32(i)
			v.UseProc = data[4]
			v.MountedOn = data[5]

			res[v.MountedOn] = v
		}
	}

	return res, nil
}

// ParserDeviceNet анализируем результаты команды cat /proc/net/dev
func ParserDeviceNet(in string) ([]model.DeviceNet, error) {
	var (
		res = []model.DeviceNet{}
		v   model.DeviceNet
		err error
	)

	scanner := bufio.NewScanner(strings.NewReader(in))
	for scanner.Scan() {
		// ищем 2 строку шапки
		if !strings.HasPrefix(strings.TrimSpace(scanner.Text()), "face") {
			continue
		}
		// нашли, далее читаем данные
		for scanner.Scan() {
			s := scanner.Text()
			data := strings.Fields(s)
			if len(data) < 11 {
				continue
			}
			v.NetInterface = data[0]
			// далее 8 полей Receive берём первые 3 (1,2,3)
			v.Receive.Bytes, err = strconv.Atoi(data[1])
			if err != nil {
				return res, err
			}
			v.Receive.Packets, err = strconv.Atoi(data[2])
			if err != nil {
				return res, err
			}
			v.Receive.Errs, err = strconv.Atoi(data[3])
			if err != nil {
				return res, err
			}
			// далее 8 полей Transmit берём первые 3  (9,10,11)
			v.Transmit.Bytes, err = strconv.Atoi(data[9])
			if err != nil {
				return res, err
			}
			v.Transmit.Packets, err = strconv.Atoi(data[10])
			if err != nil {
				return res, err
			}
			v.Transmit.Errs, err = strconv.Atoi(data[11])
			if err != nil {
				return res, err
			}
			res = append(res, v)
		}
	}
	return res, nil
}

// ParserNetworkStatistics анализируем результаты команды: ss -ta
func ParserNetworkStatistics(in string) ([]model.NetStatDetail, error) {
	var (
		res = []model.NetStatDetail{}
		v   model.NetStatDetail
		err error
	)

	scanner := bufio.NewScanner(strings.NewReader(in))
	for scanner.Scan() {
		// ищем строку шапки
		if !strings.HasPrefix(strings.TrimSpace(scanner.Text()), "State") {
			continue
		}
		// нашли, далее читаем данные
		for scanner.Scan() {
			s := scanner.Text()
			data := strings.Fields(s)
			if len(data) < 5 {
				continue
			}
			v.State = data[0]
			v.Recv, err = strconv.Atoi(data[1])
			if err != nil {
				return res, err
			}
			v.Send, err = strconv.Atoi(data[2])
			if err != nil {
				return res, err
			}
			v.LocalAddress = data[3]
			v.PeerAddress = data[4]
			res = append(res, v)
		}
	}
	return res, nil
}
