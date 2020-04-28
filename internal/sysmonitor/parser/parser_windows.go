// +build windows

package parser

import (
	"strconv"
)

// ParserSystemLoad Выдаёт заначение загрузки системы из строки с общей информацией
func ParserSystemLoad(in string) (float64, error) {
	//
	var res float64
	res, err := strconv.ParseFloat(in, 64)
	if err != nil {
		return 0, err
	}
	return res, nil
}