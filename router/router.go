package router

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/zhuxuyang/grpc_example_client/router/handlers"
	"log"
)

func InitEcho() {
	e := echo.New()
	e.Use(middleware.RequestID())
	e.Use(middleware.BodyDump(func(context echo.Context, bytes []byte, bytes2 []byte) {
		log.Println("request: ", context.Request().URL.Path, context.Request().URL.Query().Encode(), string(bytes))
		log.Println("response: ", string(bytes2))
	}))
	e.GET("/hello", handlers.Hello)
	e.Logger.Fatal(e.Start(":8080"))
}
