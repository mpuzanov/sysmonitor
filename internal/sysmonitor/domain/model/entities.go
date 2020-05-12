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

// LoadDisk информация по дисковой системе
type LoadDisk struct {
	QueryTime time.Time
	IO        []DiskIO
	FS        map[string]DiskFS
	FSInode   map[string]DiskFS
}

// DiskIO информация по загрузке дисков
type DiskIO struct {
	Device   string
	Tps      float64
	KbReadS  float64
	KbWriteS float64
	KbRead   int32
	KbWrite  int32
}

// DiskFS информация о дисках по каждой файловой системе
type DiskFS struct {
	FileSystem string
	MountedOn  string
	Used       int32
	Available  int32
	UseProc    string
}

// TalkersNet информация по сети
type TalkersNet struct {
	QueryTime time.Time
	DevNet    []DeviceNet
}

// DeviceNet статистика по использованию сетевых интерфейсов
type DeviceNet struct {
	NetInterface string
	Receive      DevNetStat
	Transmit     DevNetStat
}

// DevNetStatistics количественные показатели
type DevNetStat struct {
	Bytes   int
	Packets int
	Errs    int
}
