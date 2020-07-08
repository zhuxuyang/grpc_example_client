package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	_ "net/http/pprof"
	"sync"
	"time"

	"github.com/spf13/viper"
	"github.com/zhuxuyang/grpc_example_client/config"
	"github.com/zhuxuyang/grpc_example_client/protos"
	"google.golang.org/grpc"
)

var ExampleClient protos.ExampleClient

func initGrpcClient() {
	ctx, _ := context.WithTimeout(context.Background(), time.Second*10)
	conn, err := grpc.DialContext(ctx, viper.GetString("grpc_example_addr"), grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Panic(err)
	}
	ExampleClient = protos.NewExampleClient(conn)
}

func visit(w *sync.WaitGroup) {
	defer w.Done()
	r := rand.New(rand.NewSource(time.Now().Unix()))
	time.Sleep(time.Duration(r.Intn(10)/100) * time.Second)

	log.Println(ExampleClient.Hello(context.Background(), &protos.HelloRequest{}))
}

func main() {
	log.Println("client_example start...")
	config.InitConfig()
	initGrpcClient()

	go func() {
		http.HandleFunc("/test", handler) // each request calls handler
		err := http.ListenAndServe("localhost:"+viper.GetString("port"), nil)
		if err != nil {
			log.Panic(err)
		} else {
			log.Println("http 启动 ", viper.GetString("port"))
		}
	}()

	//w := sync.WaitGroup{}
	//	//for i := 0; i < 100; i++ {
	//	//	w.Add(1)
	//	//	go visit(&w)
	//	//}
	//	//w.Wait()

	select {}

}

func handler(w http.ResponseWriter, r *http.Request) {
	resp, err := ExampleClient.Hello(context.Background(), &protos.HelloRequest{})
	fmt.Fprintf(w, "resp %v err %v", resp, err)
}
func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "hello")
}
