package router

import (
	"context"
	"github.com/labstack/echo/v4"
	"github.com/zhuxuyang/grpc_example_client/trace"
)

//opentracing中间件
func Opentracing(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx, span, err := trace.TraceFromHeader(context.Background(), "api:"+c.Request().URL.Path, c.Request().Header)
		if err == nil {
			defer span.Finish()
			c.Set("TRACE_CONTEXT", ctx)
		} else {
			c.Set("TRACE_CONTEXT", context.Background())
		}
		return next(c)
	}
}

//func OpenTracing(comp string) echo.MiddlewareFunc {
//	return func(next echo.HandlerFunc) echo.HandlerFunc {
//		return func(c echo.Context) error {
//			var span opentracing.Span
//			opName := comp + ":" + c.Request().URL.Path
//			// 监测Header中是否有Trace信息
//			wireContext, err := opentracing.GlobalTracer().Extract(
//				opentracing.TextMap,
//				opentracing.HTTPHeadersCarrier(c.Request().Header))
//			if err != nil {
//				// 启动新Span
//				span = opentracing.StartSpan(opName)
//			} else {
//				span = opentracing.StartSpan(opName, opentracing.ChildOf(wireContext))
//			}
//
//			defer span.Finish()
//			span.SetTag("component", comp)
//			span.SetTag("span.kind", "server")
//			span.SetTag("http.url", c.Request().Host+c.Request().RequestURI)
//			span.SetTag("http.method", c.Request().Method)
//
//			if err := next(c); err != nil {
//				span.SetTag("error", true)
//				c.Error(err)
//			}
//
//			span.SetTag("error", false)
//			span.SetTag("http.status_code", c.Response().Status)
//
//			return nil
//		}
//	}
//}
