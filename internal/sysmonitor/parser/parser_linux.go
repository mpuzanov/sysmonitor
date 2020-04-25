//+build linux

package parser

import (
	"strconv"
	"strings"
)

// ParserSystemInfo Выдаёт заначение загрузки системы из строки с общей информацией
func ParserSystemInfo(in string) (float64, error) {
	//top - 19:30:09 up  4:30,  1 user,  load average: 1,02, 0,95, 0,80
	var res float64
	idx := strings.LastIndex(in, ":")
	in = strings.TrimLeft(in[idx+1:], " ") // 1,02, 0,95, 0,80
	arr := strings.Split(in, ", ")
	in = strings.Replace(arr[0], ",", ".", 1)
	res, err := strconv.ParseFloat(in, 64)
	if err != nil {
		return 0, err
	}
	return res, nil
}
