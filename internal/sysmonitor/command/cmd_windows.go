// +build windows

package command

import (
	"math/rand"
	"strconv"
	"time"

	"github.com/mpuzanov/sysmonitor/internal/rand"
)

// RunSystemLoad Возвращает строку с информацией о загрузке системы
// Заглушка для ОС Windows
func RunSystemLoad() (int, string, string) {
	var output, stderr string

	// формируем случайное значение
	val := randFloats(1, 2)

	output = strconv.FormatFloat(val, 'g', 4, 64)
	stderr = ""

	return 0, output, stderr
}

// RunLoadCPU Возвращает строку с информацией о загрузке CPU
// Заглушка для ОС Windows
func RunLoadCPU() (int, string, string) {
	var output, stderr string

	output = ""
	stderr = ""

	return 0, output, stderr
}

// randFloats случайное значение из диапазона
func randFloats(min, max float64) float64 {
	rand.Seed(time.Now().UnixNano())
	r := min + rand.Float64()*(max-min)
	return r
}
