// +build linux

package parser

import (
	"strconv"
	"strings"
	"time"

	"github.com/mpuzanov/sysmonitor/internal/sysmonitor/domain/model"
)

// ParserSystemLoad Выдаёт заначение загрузки системы из строки с общей информацией
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
	res.QueryTime = time.Now()
	res.SystemLoadValue = val

	return res, nil
}

// ParserLoadCPU Выдаёт заначение загрузки системы из строки с общей информацией
func ParserLoadCPU(in string) (model.LoadCPU, error) {
	//%Cpu(s):  1,3 us,  0,6 sy,  0,0 ni, 96,7 id,  1,3 wa,  0,0 hi,  0,1 si,  0,0 st
	var res model.LoadCPU
	// idx := strings.LastIndex(in, ":")
	// in = strings.TrimLeft(in[idx+1:], " ") // 1,02, 0,95, 0,80
	// arr := strings.Split(in, ", ")
	// in = strings.Replace(arr[0], ",", ".", 1)
	// res, err := strconv.ParseFloat(in, 64)
	// if err != nil {
	// 	return 0, err
	// }
	return res, nil
}
