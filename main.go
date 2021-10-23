package main

import (
	"context"
	"fmt"

	"os"
	"os/signal"

	"github.com/haardikdharma10/kubearmor-adapter/pkg/report"
	pb "github.com/kubearmor/KubeArmor/protobuf"
	"google.golang.org/grpc"

	"syscall"
)

func GetOSSigChannel() chan os.Signal {
	c := make(chan os.Signal, 1)

	signal.Notify(c,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
		os.Interrupt)

	return c
}

func main() {

	StopChan := make(chan struct{})

	conn, err := grpc.Dial("localhost:32767", grpc.WithInsecure()) //make it configurable

	if err != nil {
		fmt.Print(err.Error())
		return
	}

	client := pb.NewLogServiceClient(conn)

	req := pb.RequestMessage{}

	//Stream Alerts
	go func() {
		defer conn.Close()
		if stream, err := client.WatchAlerts(context.Background(), &req); err == nil {
			for {
				select {
				case <-StopChan:
					return

				default:
					res, err := stream.Recv()
					//error checking for stream.Recv()
					if err != nil {
						fmt.Print("system alerts stream stopped: " + err.Error())
					}

					//fmt.Printf("Alert:  %v\n", res) //TODO : Not print here, comment this line later;
					//Put something like a debug flag and print it (pick a level logger) zap/glog/klog
					report.Create(res) //Push res to a channel and then have the workers
				}
			}
		} else {
			fmt.Print("unable to stream systems alerts: " + err.Error())
		}
	}()
	sigChan := GetOSSigChannel()
	<-sigChan
	close(StopChan)
}

// func convert(jsonString string) (*LogsFromKubeArmor, error) {
// 	jsonDataReader := strings.NewReader(jsonString)
// 	decoder := json.NewDecoder(jsonDataReader)
// 	var controls LogsFromKubeArmor
// 	if err := decoder.Decode(&controls); err != nil {
// 		return nil, err
// 	}
// 	return &controls, nil
// }
