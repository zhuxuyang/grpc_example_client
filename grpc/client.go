package grpc

import (
	"log"
	"strconv"

	"github.com/spf13/viper"
	"github.com/zhuxuyang/grpc_example_client/protos"
	"github.com/zhuxuyang/grpc_example_client/resource"
)

var ExampleClient protos.ExampleClient

func InitGrpcClient() {
	consulConf := viper.GetStringMapString("golang-consul")

	port, _ := strconv.Atoi(consulConf["port"])
	sClient := resource.NewGrpClient(consulConf["host"], port, consulConf["token"], "xuyang")
	err := sClient.RunConsulClient()
	if err != nil {
		log.Fatalf("dial grpc tcp is failed, err is %v,config is %v ", err, consulConf)
		return
	}
	conn := sClient.Conn

	//conn, err := grpc.DialContext(context.Background(), viper.GetString("grpc_example_addr"),
	//	grpc.WithInsecure(),
	//	grpc.WithBlock(),
	//grpc.WithUnaryInterceptor(grpc_retry.UnaryClientInterceptor(
	//	grpc_retry.WithCodes(codes.Unavailable,codes.DeadlineExceeded),
	//	grpc_retry.WithBackoff(func(attempt uint) time.Duration {
	//		log.Println("重试了")
	//		return 5 * time.Second
	//	}),
	//
	//	grpc_retry.WithMax(3),
	//	),
	//)
	if err != nil {
		log.Panic(err)
	}

	ExampleClient = protos.NewExampleClient(conn)
}
