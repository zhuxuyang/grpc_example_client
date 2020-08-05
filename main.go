package main

import (
	"github.com/zhuxuyang/grpc_example_client/config"
	"github.com/zhuxuyang/grpc_example_client/grpc"
	"github.com/zhuxuyang/grpc_example_client/resource"
	"github.com/zhuxuyang/grpc_example_client/router"
	_ "net/http/pprof"
)

func main() {
	config.InitConfig()
	grpc.InitGrpcClient()
	resource.InitDB()
	resource.InitLogger()
	router.InitEcho()
}
