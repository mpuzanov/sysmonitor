package grpc

import (
	"context"
	"io"
	"log"
	"time"

	"github.com/mpuzanov/sysmonitor/pkg/sysmonitor/api"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

var (
	server  string
	timeout int32
	period  int32
)

var (
	// GrpcClientCmd .
	GrpcClientCmd = &cobra.Command{
		Use:     "grpc_client",
		Short:   "Run grpc client",
		Run:     grpcClientStart,
		Example: "sysmonitor grpc_client --server=':50051'",
	}
)

func init() {
	GrpcClientCmd.Flags().StringVarP(&server, "server", "s", "localhost:50051", "host:port to connect to")
	GrpcClientCmd.Flags().Int32VarP(&timeout, "timeout", "t", 5, "timeout(sec) for server")
	GrpcClientCmd.Flags().Int32VarP(&period, "period", "p", 15, "period(sec) for info  for server")
	err := viper.BindPFlags(GrpcClientCmd.Flags())
	if err != nil {
		log.Fatal(err)
	}
	viper.AutomaticEnv()
	server = viper.GetString("server")
	timeout = viper.GetInt32("timeout")
	period = viper.GetInt32("period")
}

func grpcClientStart(cmd *cobra.Command, args []string) {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	conn, err := grpc.DialContext(ctx, server, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("fail to dial grpc-server: %s, %v\n", server, err)
	}
	defer conn.Close()
	log.Printf("connected to %q, timeout: %d, period: %d", server, timeout, period)

	client := api.NewSysmonitorClient(conn)

	sysinfo(client, timeout, period)
}

func sysinfo(client api.SysmonitorClient, timeout int32, period int32) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	req := &api.Request{Timeout: timeout, Period: period}
	stream, err := client.SysInfo(ctx, req)
	if err != nil {
		log.Fatalf("error stream %v", err)
	}
	for {
		msg, err := stream.Recv()
		if err == io.EOF {
			log.Println("end stream")
			return
		}
		if err != nil {
			log.Fatalf("error reading stream: %v", err)
		}
		log.Println(msg)
		//log.Printf("System load value: %.4f", msg.GetSystemVal().GetSystemLoadValue())
	}
}
