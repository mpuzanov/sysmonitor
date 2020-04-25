package main

import (
	"github.com/mpuzanov/sysmonitor/cmd/sysmonitor/grpc"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "sysmonitor",
	Short: "Sysmonitor microservice",
}

func init() {
	rootCmd.AddCommand(grpc.ServerCmd, grpc.GrpcClientCmd)
}
