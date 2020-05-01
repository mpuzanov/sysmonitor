package grpcserver

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

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
	sysmon *sysmonitor.SysMonitor
}

// Start GRPC service
func Start(conf *config.Config, logger *zap.Logger, smon *sysmonitor.SysMonitor) error {

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
			log.Fatalf("Cannot start Sysmonitor: %v\n", err)
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
			log.Fatalf("Cannot listen: %s\n", err)
		}
		grpcServer = grpc.NewServer()
		api.RegisterSysmonitorServer(grpcServer, s)
		log.Printf("Starting gRPC server %s, file log: %s\n", GRPCAddr, s.cfg.Log.File)
		s.logger.Info("Starting gRPC server", zap.String("address", GRPCAddr))

		if err := grpcServer.Serve(lisn); err != nil {
			log.Fatalf("Cannot start gRPC server: %s\n", err)
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

// SysInfo .
func (s *GRPCServer) SysInfo(req *api.Request, stream api.Sysmonitor_SysInfoServer) error {
	var systemLoadValue float64

	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, os.Interrupt)

	ctx := stream.Context()

	s.logger.Info("Connected client", zap.Int32("timeout", req.Timeout), zap.Int32("period", req.Period))
	for i := 1; ; i++ {
		d := time.Duration(int64(time.Second) * int64(req.Timeout))
		select {
		case <-time.After(d):

			// если подсистема включена load_system
			if s.cfg.Collector.Category.LoadSystem {
				dataLoadSystem, err := s.sysmon.GetAvgLoadSystem(req.Period)
				if err != nil {
					s.logger.Error("GetAvgLoadSystem", zap.Error(err))
					return err
				}
				systemLoadValue = dataLoadSystem.SystemLoadValue
				err = stream.Send(&api.Result{SystemVal: &api.SystemResponse{SystemLoadValue: systemLoadValue}})
				if err != nil {
					return err
				}
			}
			if s.cfg.Collector.Category.LoadCPU {
				dataLoadCPU, err := s.sysmon.GetAvgLoadCPU(req.Period)
				if err != nil {
					s.logger.Error("GetAvgLoadCPU", zap.Error(err))
					return err
				}
				err = stream.Send(&api.Result{CpuVal: &api.CPUResponse{
					UserMode:   dataLoadCPU.UserMode,
					SystemMode: dataLoadCPU.SystemMode,
					Idle:       dataLoadCPU.Idle,
				},
				})
				if err != nil {
					return err
				}
			}

		case <-ctx.Done():
			s.logger.Error("ctx.Done() stream", zap.Error(ctx.Err()))
			return ctx.Err()

		case <-stopChan:
			s.logger.Debug("exit stream")
			return nil
		}
	}

}
