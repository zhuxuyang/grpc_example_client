package handlers

import (
	"context"
	"github.com/labstack/echo/v4"
	"github.com/zhuxuyang/grpc_example_client/grpc"
	"github.com/zhuxuyang/grpc_example_client/model"
	"github.com/zhuxuyang/grpc_example_client/protos"
	"github.com/zhuxuyang/grpc_example_client/resource"
	"net/http"
	"time"
)

func Hello(c echo.Context) error {
	m := &model.UserA{}
	resource.GetDB().
		Model(&model.UserA{}).Last(m)
	return c.JSON(http.StatusOK, map[string]interface{}{"msg": "ok", "data": nil})
}

func HelloGrpc(c echo.Context) error {
	ctx, _ := context.WithTimeout(context.Background(), time.Second*20)

	resp, err := grpc.ExampleClient.Hello(ctx, &protos.HelloRequest{
		Name: "呵呵呵呵呵",
		Time: time.Now().Unix(),
	})
	return c.JSON(http.StatusOK, map[string]interface{}{"msg": resp, "data": err})
}

func HelloV2(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]interface{}{"msg": "ok", "data": nil})
}