package grpc

import (
	"context"
	"github.com/spf13/viper"
	"github.com/zhuxuyang/grpc_example_client/protos"
	"google.golang.org/grpc"
	"log"
	"time"
)

var ExampleClient protos.ExampleClient

func InitGrpcClient() {
	ctx, _ := context.WithTimeout(context.Background(), time.Second*2)
	log.Println("address ", viper.GetString("grpc_example_addr"))
	conn, err := grpc.DialContext(ctx, viper.GetString("grpc_example_addr"), grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Panic(err)
	}
	ExampleClient = protos.NewExampleClient(conn)
}
