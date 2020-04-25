// +build windows

package command

import (
	"strconv"

	"github.com/mpuzanov/sysmonitor/internal/utils"
)

// RunSystem Возвращает строку с информацией о загрузке системы
// Заглушка для ОС Windows
func RunSystem() (int, string, string) {
	var output, stderr string

	// формируем случайное значение
	val := utils.RandFloats(1, 2)

	output = strconv.FormatFloat(val, 'g', 4, 64)
	stderr = ""

	return 0, output, stderr
}
