package jaegerMiddleware

import (
	"github.com/TianqiuHuang/grpc-client-app/pkg/trace"
	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
)

func OpenTracingMiddleware() gin.HandlerFunc {

	return func(c *gin.Context) {
		tracer, closer := trace.NewJaegerTracer("grpc-app-server")
		defer closer.Close()

		var parentSpan opentracing.Span

		spCtx, err := opentracing.GlobalTracer().Extract(opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(c.Request.Header))
		if err != nil {
			parentSpan = tracer.StartSpan(c.Request.URL.Path)
			defer parentSpan.Finish()
		} else {
			parentSpan = opentracing.StartSpan(
				c.Request.URL.Path,
				opentracing.ChildOf(spCtx),
				opentracing.Tag{Key: string(ext.Component), Value: "HTTP"},
				ext.SpanKindRPCServer,
			)
			defer parentSpan.Finish()
		}

		c.Set("Tracer", tracer)
		c.Set("ParentSpanContext", parentSpan.Context())

		c.Next()
	}
}
