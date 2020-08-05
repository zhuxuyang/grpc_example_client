package handlers

import (
	"context"
	"github.com/labstack/echo/v4"
	"github.com/zhuxuyang/grpc_example_client/grpc"
	"github.com/zhuxuyang/grpc_example_client/model"
	"github.com/zhuxuyang/grpc_example_client/protos"
	"github.com/zhuxuyang/grpc_example_client/resource"
	"google.golang.org/grpc/metadata"
	"log"
	"net/http"
)

func Hello(ctx echo.Context) error {
	md := metadata.Pairs("request_id", "hehehehe")
	c := metadata.NewOutgoingContext(context.Background(), md)

	m := &model.UserA{}
	resource.GetDB().Model(&model.UserA{}).Last(m)
	log.Println(m)
	resp, err := grpc.ExampleClient.Hello(c, &protos.HelloRequest{})
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]interface{}{"msg": "err", "data": err.Error()})
	}
	return ctx.JSON(http.StatusOK, map[string]interface{}{"msg": "ok", "data": resp})
}
