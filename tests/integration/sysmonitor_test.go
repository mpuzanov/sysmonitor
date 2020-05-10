package integration

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/cucumber/godog"
	"github.com/mpuzanov/sysmonitor/pkg/sysmonitor/api"
	"google.golang.org/grpc"
)

var (
	grpcListen string = os.Getenv("GRPC_LISTEN")
	timeoutIn  string = os.Getenv("QUERY_TIMEOUT")
	periodIn   string = os.Getenv("QUERY_PERIOD")
	timeout    int    = 5
	period     int    = 10
	err        error
)

type testSysmonitor struct {
	conn        *grpc.ClientConn
	client      api.SysmonitorClient
	stream      api.Sysmonitor_SysInfoClient
	reqProto    *api.Request
	response    *api.Result
	responseErr error
}

func init() {
	if grpcListen == "" {
		grpcListen = "localhost:50051"
	}
	if timeoutIn != "" {
		timeout, err = strconv.Atoi(timeoutIn)
		if err != nil {
			log.Fatal(err)
		}
	}
	if periodIn != "" {
		period, err = strconv.Atoi(periodIn)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func (t *testSysmonitor) startSuite() {

	t.reqProto = &api.Request{Timeout: int32(timeout), Period: int32(period)}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	t.conn, err = grpc.DialContext(ctx, grpcListen, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("fail to dial : %s, %v", grpcListen, err)
	}
	t.client = api.NewSysmonitorClient(t.conn)
}

func (t *testSysmonitor) stopSuite() {
	t.conn.Close()
}

func (t *testSysmonitor) grpcClientCallMethodSysinfo() error {
	ctx := context.Background()
	t.stream, t.responseErr = t.client.SysInfo(ctx, t.reqProto)
	if t.responseErr != nil {
		return t.responseErr
	}
	return nil
}

func (t *testSysmonitor) theErrorShouldBeNil() error {
	return t.responseErr
}

func (t *testSysmonitor) theResponseDataStreamIsNotEmpty() error {
	count := 0
	// проверяем факт потока статистики за 30 секунд
	timer1 := time.NewTimer(30 * time.Second)
	defer timer1.Stop()
loop:
	for {
		select {
		case <-timer1.C:
			break loop
		default:
			{
				t.response, t.responseErr = t.stream.Recv()
				if t.responseErr == io.EOF {
					// end stream
					return nil
				}
				if t.responseErr != nil {
					return fmt.Errorf("error reading stream: %v", t.responseErr)
				}
				if t.response == nil {
					return fmt.Errorf("error query info")
				} else {
					count++
				}
			}
		}
	}
	if count == 0 {
		return fmt.Errorf("error query info. not stream info")
	}
	return nil
}

func FeatureContext(s *godog.Suite) {
	test := &testSysmonitor{}

	s.BeforeSuite(test.startSuite)

	s.Step(`^grpc client call method Sysinfo$`, test.grpcClientCallMethodSysinfo)
	s.Step(`^The error should be nil$`, test.theErrorShouldBeNil)
	s.Step(`^The response data stream is not empty$`, test.theResponseDataStreamIsNotEmpty)

	s.AfterSuite(test.stopSuite)
}
