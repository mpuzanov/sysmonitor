package sysmonitor

import (
	"context"
	"fmt"
	"time"

	"github.com/mpuzanov/sysmonitor/internal/config"
	"github.com/mpuzanov/sysmonitor/internal/interfaces"
	"github.com/mpuzanov/sysmonitor/internal/sysmonitor/command"
	"github.com/mpuzanov/sysmonitor/internal/sysmonitor/domain/errors"
	"github.com/mpuzanov/sysmonitor/internal/sysmonitor/domain/model"
	"github.com/mpuzanov/sysmonitor/internal/sysmonitor/parser"
	"go.uber.org/zap"
)

//SysMonitor сервис календаря
type SysMonitor struct {
	data   interfaces.Storage
	cfg    *config.Config
	logger *zap.Logger
}

// NewSysmonitor - конструктор календаря
func NewSysmonitor(repository interfaces.Storage, conf *config.Config, log *zap.Logger) *SysMonitor {
	return &SysMonitor{data: repository, cfg: conf, logger: log}
}

// Run запускаем сбор данных
func (s *SysMonitor) Run(ctx context.Context) error {

	timoutCollection := s.cfg.Collector.Timeout
	s.logger.Debug("запускаем сбор данных")

	if s.cfg.Collector.Category.LoadSystem {
		go s.workerLoadSystem(ctx, timoutCollection)
	}

	return nil
}

// SaveLoadSystem .
func (s *SysMonitor) SaveLoadSystem(data *model.LoadSystem) error {
	return s.data.SaveLoadSystem(data)
}

// SaveLoadCPU .
func (s *SysMonitor) SaveLoadCPU(data *model.LoadCPU) error {
	return s.data.SaveLoadCPU(data)
}

// GetAvgLoadSystem возвращаем среднее заначение LoadSystem
func (s *SysMonitor) GetAvgLoadSystem(period int32) (*model.LoadSystem, error) {
	return s.data.GetAvgLoadSystem(period)
}

// GetAvgLoadCPU возвращаем среднее заначение LoadCPU
func (s *SysMonitor) GetAvgLoadCPU(period int32) (*model.LoadCPU, error) {
	return s.data.GetAvgLoadCPU(period)
}

// workerLoadSystem Считывание и сохранения информации по LoadSystem
func (s *SysMonitor) workerLoadSystem(ctx context.Context, timout int) error {
	for {
		d := time.Duration(int64(time.Second) * int64(timout))
		select {
		case <-time.After(d):

			systemLoadValue, err := GetSystemInfo()
			if err != nil {
				s.logger.Error("GetSystemInfo", zap.Error(err))
				return err
			}
			s.data.SaveLoadSystem(&model.LoadSystem{QueryTime: time.Now(), SystemLoadValue: systemLoadValue})
			s.logger.Debug("GetSystemInfo", zap.Float64("systemLoadValue", systemLoadValue))

		case <-ctx.Done():
			s.logger.Debug("завершаем сбор данных LoadSystem")
			return nil
		}
	}
}

// GetSystemInfo Получение значения из системы по
func GetSystemInfo() (float64, error) {
	res, txt, outerror := command.RunSystemLoad()
	if res != 0 {
		return 0, fmt.Errorf("%s. %s", outerror, errors.ErrRunReadSystemInfo)
	}
	v, err := parser.ParserSystemLoad(txt)
	if err != nil {
		return 0, fmt.Errorf("%s. %w", errors.ErrParserReadSystemInfo, err)
	}
	return v, nil
}
