package resource

import (
	"errors"
	"fmt"
	"log"

	//_ "github.com/mbobakov/grpc-consul-resolver" // It's important
	"google.golang.org/grpc"
)

type GrpcClient struct {
	Conn    *grpc.ClientConn
	Host    string
	Port    int
	Token   string
	SrvName string
}

// 通过consul调用rpc
func (s *GrpcClient) RunConsulClient() error {
	//初始化 resolver 实例
	if s.Token == "" || s.Host == "" || s.Port == 0 {
		err := errors.New(fmt.Sprintf(
			"consul configuration empty [host: %v][port: %v][token: %v]",
			s.Host, s.Port, s.Token))
		log.Println(err)
		return err
	}
	Init(s.Host, s.Port, s.Token, s.SrvName)
	conn, err := grpc.Dial(
		fmt.Sprintf("%s://%s:%d/%s", "consul", s.Host, s.Port, s.SrvName),
		//不能block => blockkingPicker打开，在调用轮询时picker_wrapper => picker时若block则不进行robin操作直接返回失败
		//grpc.WithBlock(),
		grpc.WithInsecure(),
		//指定初始化round_robin => balancer (后续可以自行定制balancer和 register、resolver 同样的方式)
		grpc.WithDefaultServiceConfig(`{"LoadBalancingPolicy":"round_robin"}`),
	)
	if err != nil {
		return err
	}
	s.Conn = conn
	log.Printf(fmt.Sprintf("gRpc consul client [%s] start success", s.SrvName))
	return nil
}

func NewGrpClient(consulHost string, consulPort int, consulToken, srvName string) *GrpcClient {
	return &GrpcClient{
		Host:    consulHost,
		Port:    consulPort,
		Token:   consulToken,
		SrvName: srvName,
	}
}
