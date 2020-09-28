package grpc

import (
	"context"
	"github.com/grpc-ecosystem/go-grpc-middleware/retry"
	"github.com/spf13/viper"
	"github.com/zhuxuyang/grpc_example_client/protos"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"log"
	"time"
)

var ExampleClient protos.ExampleClient

func InitGrpcClient() {
	ctx, _ := context.WithTimeout(context.Background(), time.Second*10)
	log.Println("address ", viper.GetString("grpc_example_addr"))
	//conn, err := grpc.DialContext(ctx, viper.GetString("grpc_example_addr"), grpc.WithInsecure(), grpc.WithBlock())
	//if err != nil {
	//	log.Panic(err)
	//}

	conn, err := grpc.DialContext(ctx, viper.GetString("grpc_example_addr"),
		grpc.WithInsecure(),
		grpc.WithBlock(),
		grpc.WithUnaryInterceptor(grpc_retry.UnaryClientInterceptor(
			grpc_retry.WithCodes(codes.Unavailable,codes.DeadlineExceeded),
			grpc_retry.WithBackoff(func(attempt uint) time.Duration {
				log.Println("重试了")
				return 5 * time.Second
			}),

			grpc_retry.WithMax(3),
			),
		))
	if err != nil {
		log.Panic(err)
	}

	ExampleClient = protos.NewExampleClient(conn)
}
