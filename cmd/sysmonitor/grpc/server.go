package grpc

import (
	"log"

	"github.com/mpuzanov/sysmonitor/internal/config"
	"github.com/mpuzanov/sysmonitor/internal/grpcserver"
	"github.com/mpuzanov/sysmonitor/internal/storage"
	"github.com/mpuzanov/sysmonitor/internal/sysmonitor"
	"github.com/mpuzanov/sysmonitor/pkg/logger"

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
		log.Fatal(err)
	}

}

func grpcServerStart(cmd *cobra.Command, args []string) {

	cfg, err := config.LoadConfig(cfgPath)
	if err != nil {
		log.Fatalf("Fail load %s: %s", cfgPath, err)
	}

	logger := logger.NewLogger(cfg.Log)

	db := storage.NewStorage()

	sysmonitor := sysmonitor.NewSysmonitor(db, cfg, logger)

	if err := grpcserver.Start(cfg, logger, sysmonitor); err != nil {
		log.Fatal(err)
	}
}
