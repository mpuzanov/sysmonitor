// +build windows

package parser

import (
	"strconv"

	"github.com/mpuzanov/sysmonitor/internal/sysmonitor/domain/model"
)

// ParserSystemLoad Выдаёт заначение загрузки системы из строки с общей информацией
func ParserSystemLoad(in string) (model.LoadSystem, error) {
	//
	var res model.LoadSystem
	val, err := strconv.ParseFloat(in, 64)
	if err != nil {
		return res, err
	}
	res.SystemLoadValue = val
	return res, nil
}

// ParserLoadCPU Выдаёт заначение загрузки системы из строки с общей информацией
func ParserLoadCPU(in string) (model.LoadCPU, error) {
	var res model.LoadCPU

	return res, nil
}
