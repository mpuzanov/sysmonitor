package model

import "time"

// LoadSystem Средняя загрузка системы
type LoadSystem struct {
	QueryTime       time.Time
	SystemLoadValue float64
}

// LoadCPU Средняя загрузка CPU
type LoadCPU struct {
	QueryTime  time.Time
	UserMode   float64
	SystemMode float64
	Idle       float64
}
