package grpc

import (
	"context"
	"io"
	"log"

	"github.com/mpuzanov/sysmonitor/pkg/sysmonitor/api"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
)

var (
	server  string
	timeOut int32
	period  int32
)

var (
	// GrpcClientCmd .
	GrpcClientCmd = &cobra.Command{
		Use:   "grpc_client",
		Short: "Run grpc client",
		Run:   grpcClientStart,
	}
)

func init() {
	GrpcClientCmd.Flags().StringVar(&server, "server", ":50051", "host:port to connect to")
	GrpcClientCmd.Flags().Int32VarP(&timeOut, "timeout", "t", 5, "timeout(sec) for server")
	GrpcClientCmd.Flags().Int32VarP(&period, "period", "p", 15, "period(sec) for info  for server")
}

func grpcClientStart(cmd *cobra.Command, args []string) {

	conn, err := grpc.Dial(server, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("fail to dial : %s, %v\n", server, err)
	}
	defer conn.Close()
	client := api.NewSysmonitorClient(conn)

	sysinfo(client, timeOut, period)
}

func sysinfo(client api.SysmonitorClient, timeOut int32, period int32) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	req := &api.Request{Timeout: timeOut, Period: period}
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
