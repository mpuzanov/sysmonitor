package grpc

import (
	"context"
	"fmt"
	"io"
	"log"
	"strings"
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/mpuzanov/sysmonitor/pkg/sysmonitor/api"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

const (
	layout = "2006.01.02 15.04.05 (MST)"
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

	locZone, _ := time.LoadLocation("Europe/Samara")
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
		if msg.SystemVal != nil {
			t, _ := ptypes.Timestamp(msg.SystemVal.GetQueryTime())
			fmt.Printf("InfoSystem: QueryTime: %s, SystemLoadValue:%v\n", t.In(locZone).Format(layout), msg.SystemVal.SystemLoadValue)
		}
		if msg.CpuVal != nil {
			t, _ := ptypes.Timestamp(msg.CpuVal.GetQueryTime())
			fmt.Printf("InfoCPU: QueryTime: %s, UserMode: %v, SystemMode: %v, Idle: %v\n", t.In(locZone).Format(layout),
				msg.GetCpuVal().GetUserMode(),
				msg.GetCpuVal().GetSystemMode(),
				msg.GetCpuVal().GetIdle(),
			)
		}

		if msg.DiskVal != nil {
			t, _ := ptypes.Timestamp(msg.DiskVal.GetQueryTime())
			fmt.Printf("InfoDisk: QueryTime: %s\n", t.In(locZone).Format(layout))

			fmt.Printf("\n%-10v  %10v %10v %10v %10v %10v\n", "Device", "Tps", "KbReadS", "KbWriteS", "KbRead", "KbWrite")
			fmt.Println(strings.Repeat("-", 80))
			io := msg.GetDiskVal().Io
			for _, val := range io {
				fmt.Printf("%-10v  %10v %10v %10v %10v %10v\n", val.Device, val.Tps, val.KbReadS, val.KbWriteS, val.KbRead, val.KbWrite)
			}

			fmt.Printf("\n%-15v  %10v %10v %10v %10v %10v %10v %30v\n", "FileSystem",
				"Used", "Available", "Use%", "Used_Inode", "Available_Inode", "Use%_Inode",
				"MountedOn")
			fmt.Println(strings.Repeat("-", 100))
			fs := msg.GetDiskVal().Fs
			for _, val := range fs {
				fmt.Printf("%-15v  %10v %10v %10v %10v %10v %10v %30v\n", val.FileSystem,
					val.Used, val.Available, val.UseProc, val.UsedInode, val.AvailableInode, val.UseProcInode,
					val.MountedOn)
			}
		}

	}
}
