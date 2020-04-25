//+build windows

package parser

import (
	"strconv"
)

// ParserSystemInfo Выдаёт заначение загрузки системы из строки с общей информацией
func ParserSystemInfo(in string) (float64, error) {
	//
	var res float64
	res, err := strconv.ParseFloat(in, 64)
	if err != nil {
		return 0, err
	}
	return res, nil
}
