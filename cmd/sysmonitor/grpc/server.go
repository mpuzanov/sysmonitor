package grpc

import (
	"log"

	"github.com/mpuzanov/sysmonitor/internal/config"
	"github.com/mpuzanov/sysmonitor/internal/grpcserver"
	"github.com/mpuzanov/sysmonitor/internal/storage"
	"github.com/mpuzanov/sysmonitor/internal/sysmonitor"
	"github.com/mpuzanov/sysmonitor/pkg/logger"

	"github.com/spf13/cobra"
)

var cfgPath string

var (
	// ServerCmd .
	ServerCmd = &cobra.Command{
		Use:   "grpc_server",
		Short: "Run grpc server",
		Run:   grpcServerStart,
	}
)

func init() {
	ServerCmd.Flags().StringVarP(&cfgPath, "config", "c", "", "path to the configuration file")
}

func grpcServerStart(cmd *cobra.Command, args []string) {
	cfg, err := config.LoadConfig(cfgPath)
	if err != nil {
		log.Fatalf("unable to load %s: %s", cfgPath, err)
	}
	logger := logger.NewLogger(cfg.Log)

	db, err := storage.NewStorage()

	sysmonitor := sysmonitor.NewSysmonitor(db, cfg, logger)

	if err := grpcserver.Start(cfg, logger, sysmonitor); err != nil {
		log.Fatal(err)
	}
}
