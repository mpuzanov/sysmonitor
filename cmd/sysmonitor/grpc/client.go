package grpc

import (
	"context"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/mpuzanov/sysmonitor/pkg/logger"
	"github.com/mpuzanov/sysmonitor/pkg/sysmonitor/api"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

const (
	layout = "2006.01.02 15.04.05 (MST)"
)

var (
	address string
	timeout int32
	period  int32

	sys     bool
	cpu     bool
	disk    bool
	toptalk bool
	netstat bool

	locZone *time.Location
)

var (
	// GrpcClientCmd .
	GrpcClientCmd = &cobra.Command{
		Use:     "grpc_client",
		Short:   "Run grpc client",
		Run:     grpcClientStart,
		Example: "sysmonitor grpc_client --address=':50051'",
	}
)

func init() {
	GrpcClientCmd.Flags().StringVarP(&address, "address", "", "localhost:50051", "host:port to connect to")
	GrpcClientCmd.Flags().Int32VarP(&timeout, "timeout", "", 5, "timeout(sec) for server")
	GrpcClientCmd.Flags().Int32VarP(&period, "period", "", 15, "period(sec) for info  for server")
	GrpcClientCmd.Flags().BoolVarP(&sys, "sys", "s", false, "collecting statistics on the load average")
	GrpcClientCmd.Flags().BoolVarP(&cpu, "cpu", "c", false, "collecting statistics on the CPU")
	GrpcClientCmd.Flags().BoolVarP(&disk, "disk", "d", false, "collecting statistics on the Disk")
	GrpcClientCmd.Flags().BoolVarP(&toptalk, "toptalk", "t", false, "collecting statistics top talkers")
	GrpcClientCmd.Flags().BoolVarP(&netstat, "netstat", "n", false, "collecting statistics on the network")
	err := viper.BindPFlags(GrpcClientCmd.Flags())
	if err != nil {
		logger.LogSugar.Fatal(err)
	}
	viper.AutomaticEnv()
	address = viper.GetString("address")
	timeout = viper.GetInt32("timeout")
	period = viper.GetInt32("period")
}

func grpcClientStart(cmd *cobra.Command, args []string) {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	conn, err := grpc.DialContext(ctx, address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		logger.LogSugar.Fatalf("fail to dial grpc-server: %s, %v", address, err)
	}
	defer conn.Close()
	logger.LogSugar.Infof("connected to %q, timeout: %d, period: %d", address, timeout, period)

	client := api.NewSysmonitorClient(conn)

	//если ни одна подсистема показа статистики не включена, то включаем все
	if !sys && !cpu && !disk && !toptalk && !netstat {
		sys = true
		cpu = true
		disk = true
		toptalk = true
		netstat = true
	}

	locZone, err = time.LoadLocation("Europe/Samara")
	if err != nil {
		logger.LogSugar.Fatalf("fail load location %v", err)
	}

	sysinfo(client, timeout, period)
}

func sysinfo(client api.SysmonitorClient, timeout int32, period int32) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	req := &api.Request{Timeout: timeout, Period: period}

	stream, err := client.SysInfo(ctx, req)
	if err != nil {
		logger.LogSugar.Fatalf("error stream %v", err)
	}
	for {
		msg, err := stream.Recv()
		if err == io.EOF {
			logger.LogSugar.Info("end stream")
			return
		}
		if err != nil {
			logger.LogSugar.Fatalf("error reading from stream: %v", err)
		}

		printResult(msg)
	}
}

// printResult вывод результатов в консоль
func printResult(msg *api.Result) {
	if sys && msg.SystemVal != nil {
		t, _ := ptypes.Timestamp(msg.SystemVal.GetQueryTime())
		fmt.Printf("\nInfoSystem: QueryTime: %s, SystemLoadValue:%v\n", t.In(locZone).Format(layout), msg.SystemVal.SystemLoadValue)
	}
	if cpu && msg.CpuVal != nil {
		t, _ := ptypes.Timestamp(msg.CpuVal.GetQueryTime())
		fmt.Printf("\nInfoCPU: QueryTime: %s, UserMode: %v, SystemMode: %v, Idle: %v\n", t.In(locZone).Format(layout),
			msg.GetCpuVal().GetUserMode(),
			msg.GetCpuVal().GetSystemMode(),
			msg.GetCpuVal().GetIdle(),
		)
	}

	if disk && msg.DiskVal != nil {
		t, _ := ptypes.Timestamp(msg.DiskVal.GetQueryTime())
		fmt.Printf("\nInfoDisk: QueryTime: %s\n", t.In(locZone).Format(layout))

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

	if toptalk && msg.TalkerNetVal != nil {
		t, _ := ptypes.Timestamp(msg.TalkerNetVal.GetQueryTime())
		fmt.Printf("\nInfoTalkerNet: QueryTime: %s", t.In(locZone).Format(layout))
		fmt.Printf("\n%-20v|%-30v  |%-30v", "", "Receive", "Transmit")
		fmt.Printf("\n%-20v|%10v|%10v|%10v|%10v|%10v|%10v\n", "Interface", "bytes", "packets", "errs", "bytes", "packets", "packets")
		fmt.Println(strings.Repeat("-", 86))
		dat := msg.GetTalkerNetVal().Devnet
		for _, val := range dat {
			fmt.Printf("%-20v|%10v|%10v|%10v|%10v|%10v|%10v\n", val.NetInterface, val.ReceiveBytes, val.ReceivePackets, val.ReceiveErrs,
				val.TransmitBytes, val.TransmitPackets, val.TransmitErrs)
		}
	}

	if netstat && msg.NetstatVal != nil {
		t, _ := ptypes.Timestamp(msg.NetstatVal.GetQueryTime())
		fmt.Printf("\nInfoNetworkStatistics: QueryTime: %s", t.In(locZone).Format(layout))
		fmt.Printf("\n%-15v|%10v|%10v|%25v|%20v\n", "State", "Recv", "Send", "LocalAddress", "PeerAddress")
		fmt.Println(strings.Repeat("-", 86))
		dat := msg.GetNetstatVal().Netstat
		for _, val := range dat {
			fmt.Printf("%-15v|%10v|%10v|%25v|%20v\n",
				val.State, val.Recv, val.Send, val.LocalAddress, val.PeerAddress)
		}
	}
}
