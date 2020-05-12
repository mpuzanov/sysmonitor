package grpcserver

import (
	"context"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/mpuzanov/sysmonitor/internal/config"
	"github.com/mpuzanov/sysmonitor/internal/sysmonitor"
	"github.com/mpuzanov/sysmonitor/pkg/sysmonitor/api"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

var grpcServer *grpc.Server

//GRPCServer структура сервера
type GRPCServer struct {
	cfg    *config.Config
	logger *zap.Logger
	sysmon *sysmonitor.Sysmonitor
}

// Start GRPC service
func Start(conf *config.Config, logger *zap.Logger, smon *sysmonitor.Sysmonitor) error {

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)
	defer signal.Stop(interrupt)

	wg := sync.WaitGroup{}

	wg.Add(1)
	// запуск Sysmonitor
	go func() {
		defer wg.Done()
		err := smon.Run(ctx)
		if err != nil {
			logger.Error("Cannot start Sysmonitor", zap.Error(err))
			return
		}
	}()

	// запуск gRPC server
	wg.Add(1)
	go func() {
		defer wg.Done()

		s := &GRPCServer{
			cfg:    conf,
			logger: logger,
			sysmon: smon,
		}
		GRPCAddr := s.cfg.Host + ":" + s.cfg.Port
		lisn, err := net.Listen("tcp", GRPCAddr)
		if err != nil {
			s.logger.Error("Cannot listen", zap.Error(err))
			return
		}
		grpcServer = grpc.NewServer()
		api.RegisterSysmonitorServer(grpcServer, s)

		s.logger.Info("Starting gRPC server", zap.String("address", GRPCAddr))

		if err := grpcServer.Serve(lisn); err != nil {
			s.logger.Error("Cannot start gRPC server", zap.Error(err))
			return
		}
	}()

	select {
	case <-interrupt:
		break
	case <-ctx.Done():
		break
	}
	logger.Debug("received shutdown signal")

	cancel()
	if grpcServer != nil {
		grpcServer.GracefulStop()
	}

	wg.Wait()
	//logger.Debug("The end grpc")
	return nil
}

// SysInfo Раздача клиентам сервиса GRPC информации по системе
func (s *GRPCServer) SysInfo(req *api.Request, stream api.Sysmonitor_SysInfoServer) error {

	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, os.Interrupt)

	ctx := stream.Context()

	s.logger.Info("Connected client", zap.Int32("timeout", req.Timeout), zap.Int32("period", req.Period))
	for {
		d := time.Duration(int64(time.Second) * int64(req.Timeout))
		select {
		case <-time.After(d):

			// если подсистема load_system включена
			if s.cfg.Collector.Category.LoadSystem {
				data, err := s.sysmon.GetAvgLoadSystem(req.Period)
				if err != nil {
					s.logger.Error("GetAvgLoadSystem", zap.Error(err))
					return err
				}
				err = stream.Send(&api.Result{SystemVal: &api.SystemResponse{
					QueryTime:       ptypes.TimestampNow(),
					SystemLoadValue: data.SystemLoadValue}})
				if err != nil {
					return err
				}
			}
			if s.cfg.Collector.Category.LoadCPU {
				data, err := s.sysmon.GetAvgLoadCPU(req.Period)
				if err != nil {
					s.logger.Error("GetAvgLoadCPU", zap.Error(err))
					return err
				}
				queryTimeProto, err := ptypes.TimestampProto(data.QueryTime)
				if err != nil {
					s.logger.Error("convert QueryTime", zap.Error(err))
					return err
				}
				err = stream.Send(&api.Result{CpuVal: &api.CPUResponse{
					QueryTime:  queryTimeProto,
					UserMode:   data.UserMode,
					SystemMode: data.SystemMode,
					Idle:       data.Idle,
				},
				})
				if err != nil {
					return err
				}
			}

			if s.cfg.Collector.Category.LoadDisk {
				data := s.sysmon.GetInfoDisk()
				valueProto, err := ParserLoadDiskToProto(data)
				if err != nil {
					s.logger.Error("ParserLoadDiskToProto", zap.Error(err))
					return err
				}
				err = stream.Send(&api.Result{DiskVal: valueProto})
				if err != nil {
					return err
				}
			}

			if s.cfg.Collector.Category.TopTalkers {
				data, err := s.sysmon.GetAvgTalkersNet(req.Period)
				if err != nil {
					s.logger.Error("GetAvgTalkersNet", zap.Error(err))
					return err
				}
				valueProto, err := ParserTalkerNetToProto(data)
				if err != nil {
					s.logger.Error("ParserTalkerNetToProto", zap.Error(err))
					return err
				}
				err = stream.Send(&api.Result{TalkerNetVal: valueProto})
				if err != nil {
					return err
				}
			}

		case <-ctx.Done():
			s.logger.Error("end stream to client", zap.Error(ctx.Err()))
			return ctx.Err()

		case <-stopChan:
			s.logger.Debug("exit stream")
			return nil
		}
	}

}
