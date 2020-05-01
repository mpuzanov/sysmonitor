package sysmonitor

import (
	"context"
	"fmt"
	"log"
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
		go log.Fatal(s.workerLoadSystem(ctx, timoutCollection))
	}

	if s.cfg.Collector.Category.LoadCPU {
		go log.Fatal(s.workerLoadCPU(ctx, timoutCollection))
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

			res, err := GetInfoSystem()
			if err != nil {
				s.logger.Error("GetInfoSystem", zap.Error(err))
				return err
			}
			err = s.data.SaveLoadSystem(&res)
			if err != nil {
				return err
			}
			s.logger.Debug("GetInfoSystem", zap.Float64("systemLoadValue", res.SystemLoadValue))

		case <-ctx.Done():
			s.logger.Debug("завершаем сбор данных LoadSystem")
			return nil
		}
	}
}

// workerLoadCPU Считывание и сохранения информации по LoadSystem
func (s *SysMonitor) workerLoadCPU(ctx context.Context, timout int) error {

	for {
		d := time.Duration(int64(time.Second) * int64(timout))

		select {
		case <-time.After(d):

			res, err := GetInfoCPU()
			if err != nil {
				s.logger.Error("GetInfoCPU", zap.Error(err))
				return err
			}
			err = s.data.SaveLoadCPU(&res)
			if err != nil {
				return err
			}
			s.logger.Debug("GetInfoCPU", zap.Float64("Idle", res.Idle))

		case <-ctx.Done():
			s.logger.Debug("завершаем сбор данных LoadCPU")
			return nil
		}
	}
}

// GetInfoSystem Получение значения из системы по
func GetInfoSystem() (model.LoadSystem, error) {
	var res model.LoadSystem
	exitCode, txt, outerror := command.RunSystemLoad()
	if exitCode != 0 {
		return res, fmt.Errorf("%s. %s", outerror, errors.ErrRunReadInfoSystem)
	}
	res, err := parser.ParserSystemLoad(txt)
	if err != nil {
		return res, fmt.Errorf("%s. %w", errors.ErrParserReadInfoSystem, err)
	}
	return res, nil
}

// GetInfoCPU Получение значения из системы по
func GetInfoCPU() (model.LoadCPU, error) {
	var res model.LoadCPU
	exitCode, txt, outerror := command.RunLoadCPU()
	if exitCode != 0 {
		return res, fmt.Errorf("%s. %s", outerror, errors.ErrRunReadInfoCPU)
	}
	res, err := parser.ParserLoadCPU(txt)
	if err != nil {
		return res, fmt.Errorf("%s. %w", errors.ErrParserReadInfoCPU, err)
	}
	return res, nil
}
