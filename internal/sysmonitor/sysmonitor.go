package sysmonitor

import (
	"context"
	"time"

	"github.com/mpuzanov/sysmonitor/internal/config"
	"github.com/mpuzanov/sysmonitor/internal/repository"
	"github.com/mpuzanov/sysmonitor/internal/sysmonitor/domain/model"
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

	if s.cfg.Collector.Category.TopTalkers {
		go func() {
			err := s.workerTalkersNet(ctx, timoutCollection)
			if err != nil {
				s.logger.Error("Cannot start workerTalkersNet", zap.Error(err))
				return
			}
		}()
	}

	if s.cfg.Collector.Category.NetworkStat {
		go func() {
			err := s.workerNetworkStatistics(ctx, timoutCollection)
			if err != nil {
				s.logger.Error("Cannot start workerNetworkStatistics", zap.Error(err))
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

// SaveTalkersNet сохраняем информацию
func (s *Sysmonitor) SaveTalkersNet(data *model.TalkersNet) error {
	return s.data.SaveTalkersNet(data)
}

// SaveNetworkStatistics сохраняем информацию
func (s *Sysmonitor) SaveNetworkStatistics(data *model.NetworkStatistics) error {
	return s.data.SaveNetworkStatistics(data)
}

// GetAvgLoadSystem возвращаем среднее значение LoadSystem
func (s *Sysmonitor) GetAvgLoadSystem(period int32) (*model.LoadSystem, error) {
	return s.data.GetAvgLoadSystem(period)
}

// GetAvgLoadCPU возвращаем среднее значение LoadCPU
func (s *Sysmonitor) GetAvgLoadCPU(period int32) (*model.LoadCPU, error) {
	return s.data.GetAvgLoadCPU(period)
}

// GetAvgTalkersNet возвращаем среднее значение
func (s *Sysmonitor) GetAvgTalkersNet(period int32) (*model.TalkersNet, error) {
	return s.data.GetAvgTalkersNet(period)
}

// GetAvgNetworkStatistics возвращаем среднее значение
func (s *Sysmonitor) GetAvgNetworkStatistics(period int32) (*model.NetworkStatistics, error) {
	return s.data.GetAvgNetworkStatistics(period)
}

// GetInfoDisk возвращаем информацию по дискам
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

// workerTalkersNet Считывание и сохранение информации по TalkersNet
func (s *Sysmonitor) workerTalkersNet(ctx context.Context, timout int) error {
	s.logger.Info("starting collection TalkersNet")
	for {
		d := time.Duration(int64(time.Second) * int64(timout))

		select {
		case <-time.After(d):

			res, err := QueryInfoTalkersNet()
			if err != nil {
				s.logger.Error("QueryInfoTalkersNet", zap.Error(err))
				return err
			}
			err = s.data.SaveTalkersNet(&res)
			if err != nil {
				return err
			}

		case <-ctx.Done():
			s.logger.Info("completing data collection TalkersNet")
			return nil
		}
	}
}

// workerNetworkStatistics Считывание и сохранение информации по NetworkStatistics
func (s *Sysmonitor) workerNetworkStatistics(ctx context.Context, timout int) error {
	s.logger.Info("starting collection NetworkStatistics")
	for {
		d := time.Duration(int64(time.Second) * int64(timout))

		select {
		case <-time.After(d):

			res, err := QueryInfoNetworkStatistics()
			if err != nil {
				s.logger.Error("QueryInfoNetworkStatistics", zap.Error(err))
				return err
			}
			err = s.data.SaveNetworkStatistics(&res)
			if err != nil {
				return err
			}

		case <-ctx.Done():
			s.logger.Info("completing data collection NetworkStatistics")
			return nil
		}
	}
}
