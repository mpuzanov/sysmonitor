package sysmonitor

import (
	"context"
	"fmt"
	"time"

	"github.com/mpuzanov/sysmonitor/internal/config"
	"github.com/mpuzanov/sysmonitor/internal/repository"
	"github.com/mpuzanov/sysmonitor/internal/sysmonitor/command"
	"github.com/mpuzanov/sysmonitor/internal/sysmonitor/domain/errors"
	"github.com/mpuzanov/sysmonitor/internal/sysmonitor/domain/model"
	"github.com/mpuzanov/sysmonitor/internal/sysmonitor/parser"
	"go.uber.org/zap"
)

//Sysmonitor сервис для сбора и выдачи информации по системе
type Sysmonitor struct {
	data   repository.Storage
	cfg    *config.Config
	logger *zap.Logger
}

// NewSysmonitor - конструктор сервиса
func NewSysmonitor(repo repository.Storage, conf *config.Config, log *zap.Logger) *Sysmonitor {
	return &Sysmonitor{data: repo, cfg: conf, logger: log}
}

// Run запускаем сбор данных
func (s *Sysmonitor) Run(ctx context.Context) error {

	timoutCollection := s.cfg.Collector.Timeout
	s.logger.Debug("запускаем сбор данных")

	if s.cfg.Collector.Category.LoadSystem {
		go func() {
			err := s.workerLoadSystem(ctx, timoutCollection)
			if err != nil {
				s.logger.Error("Cannot start workerLoadSystem", zap.Error(err))
				return
			}
		}()
	}

	if s.cfg.Collector.Category.LoadCPU {
		go func() {
			err := s.workerLoadCPU(ctx, timoutCollection)
			if err != nil {
				s.logger.Error("Cannot start workerLoadCPU", zap.Error(err))
				return
			}
		}()
	}

	if s.cfg.Collector.Category.LoadDisk {
		go func() {
			err := s.workerLoadDisk(ctx, timoutCollection)
			if err != nil {
				s.logger.Error("Cannot start workerLoadDisk", zap.Error(err))
				return
			}
		}()
	}

	return nil
}

// SaveLoadSystem сохраняем информацию по загрузке системы
func (s *Sysmonitor) SaveLoadSystem(data *model.LoadSystem) error {
	return s.data.SaveLoadSystem(data)
}

// SaveLoadCPU сохраняем информацию по загрузке CPU
func (s *Sysmonitor) SaveLoadCPU(data *model.LoadCPU) error {
	return s.data.SaveLoadCPU(data)
}

// GetAvgLoadSystem возвращаем среднее заначение LoadSystem
func (s *Sysmonitor) GetAvgLoadSystem(period int32) (*model.LoadSystem, error) {
	return s.data.GetAvgLoadSystem(period)
}

// GetAvgLoadCPU возвращаем среднее заначение LoadCPU
func (s *Sysmonitor) GetAvgLoadCPU(period int32) (*model.LoadCPU, error) {
	return s.data.GetAvgLoadCPU(period)
}

// GetInfoDisk .
func (s *Sysmonitor) GetInfoDisk() *model.LoadDisk {
	return s.data.GetInfoDisk()
}

// workerLoadSystem Считывание и сохранение информации по LoadSystem
func (s *Sysmonitor) workerLoadSystem(ctx context.Context, timout int) error {
	s.logger.Info("starting collection LoadSystem")
	for {
		d := time.Duration(int64(time.Second) * int64(timout))

		select {
		case <-time.After(d):

			res, err := QueryInfoSystem()
			if err != nil {
				s.logger.Error("QueryInfoSystem", zap.Error(err))
				return err
			}
			err = s.data.SaveLoadSystem(&res)
			if err != nil {
				return err
			}
			s.logger.Debug("QueryInfoSystem", zap.Float64("systemLoadValue", res.SystemLoadValue))

		case <-ctx.Done():
			s.logger.Info("completing data collection LoadSystem")
			return nil
		}
	}
}

// workerLoadCPU Считывание и сохранение информации по LoadSystem
func (s *Sysmonitor) workerLoadCPU(ctx context.Context, timout int) error {
	s.logger.Info("starting collection LoadCPU")
	for {
		d := time.Duration(int64(time.Second) * int64(timout))

		select {
		case <-time.After(d):

			res, err := QueryInfoCPU()
			if err != nil {
				s.logger.Error("QueryInfoCPU", zap.Error(err))
				return err
			}
			err = s.data.SaveLoadCPU(&res)
			if err != nil {
				return err
			}
			s.logger.Debug("QueryInfoCPU", zap.Float64("Idle", res.Idle))

		case <-ctx.Done():
			s.logger.Info("completing data collection LoadCPU")
			return nil
		}
	}
}

// workerLoadDisk Считывание и сохранение информации по дисковой системе
func (s *Sysmonitor) workerLoadDisk(ctx context.Context, timout int) error {
	s.logger.Info("starting collection LoadDisk")
	for {
		d := time.Duration(int64(time.Second) * int64(timout))

		select {
		case <-time.After(d):

			res, err := QueryInfoDisk()
			if err != nil {
				s.logger.Error("QueryInfoDisk", zap.Error(err))
				return err
			}
			s.data.SaveLoadDisk(&res)
			s.logger.Debug("Query QueryInfoDisk")

		case <-ctx.Done():
			s.logger.Info("completing data collection LoadDisk")
			return nil
		}
	}
}

// QueryInfoSystem Получение информации из системы по LoadSystem
func QueryInfoSystem() (model.LoadSystem, error) {
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

// QueryInfoCPU Получение информации из системы по LoadCPU
func QueryInfoCPU() (model.LoadCPU, error) {
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

// QueryInfoDisk Получение информации по дисковой системе
func QueryInfoDisk() (model.LoadDisk, error) {
	var res model.LoadDisk

	exitCode, txt, outerror := command.RunCommand("iostat", "-d", "-k")
	if exitCode != 0 {
		return res, fmt.Errorf("%s. %s", outerror, errors.ErrRunLoadDiskDevice)
	}
	resIO, err := parser.ParserLoadDiskDevice(txt)
	if err != nil {
		return res, fmt.Errorf("%s. %w", errors.ErrParserLoadDiskDevice, err)
	}

	exitCode, txt, outerror = command.RunCommand("df", "-k")
	if exitCode != 0 {
		return res, fmt.Errorf("%s. %s", outerror, errors.ErrRunLoadDiskFS)
	}
	resFS, err := parser.ParserLoadDiskFS(txt)
	if err != nil {
		return res, fmt.Errorf("%s. %w", errors.ErrParserLoadDiskFS, err)
	}
	exitCode, txt, outerror = command.RunCommand("df", "-i")
	if exitCode != 0 {
		return res, fmt.Errorf("%s. %s", outerror, errors.ErrRunLoadDiskFSInode)
	}
	resFSInode, err := parser.ParserLoadDiskFS(txt)
	if err != nil {
		return res, fmt.Errorf("%s. %w", errors.ErrParserLoadDiskFSInode, err)
	}
	res.IO = resIO
	res.FS = resFS
	res.FSInode = resFSInode
	res.QueryTime = time.Now()
	return res, nil
}
