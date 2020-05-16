package grpc

import (
	"github.com/mpuzanov/sysmonitor/internal/config"
	"github.com/mpuzanov/sysmonitor/internal/grpcserver"
	"github.com/mpuzanov/sysmonitor/internal/storage"
	"github.com/mpuzanov/sysmonitor/internal/sysmonitor"
	"github.com/mpuzanov/sysmonitor/pkg/logger"
	"go.uber.org/zap"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgPath string
	port    string
)

var (
	// ServerCmd .
	ServerCmd = &cobra.Command{
		Use:     "grpc_server",
		Short:   "Run grpc server",
		Run:     grpcServerStart,
		Example: "sysmonitor grpc_server --config=configs/prod.yaml",
	}
)

func init() {

	ServerCmd.Flags().StringVarP(&cfgPath, "config", "c", "", "path to the configuration file")
	ServerCmd.Flags().StringVarP(&port, "port", "p", "50051", "port grpc server")
	err := viper.BindPFlag("port", ServerCmd.Flags().Lookup("port"))
	if err != nil {
		logger.LogSugar.Fatal("viper.BindPFlag", zap.Error(err))
	}

}

func grpcServerStart(cmd *cobra.Command, args []string) {

	cfg, err := config.LoadConfig(cfgPath)
	if err != nil {
		logger.LogSugar.Fatal("Fail load", zap.String("config path", cfgPath), zap.Error(err))
	}

	l := logger.NewLogger(cfg.Log)

	db := storage.NewMemoryStorage()

	sysmonitor := sysmonitor.NewSysmonitor(db, cfg, l)

	if err := grpcserver.Start(cfg, l, sysmonitor); err != nil {
		logger.LogSugar.Fatal("fail start grpc-server", zap.Error(err))
	}
}
