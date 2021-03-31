package jaegerMiddleware

import (
	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
)

func OpenTracingMiddleware() gin.HandlerFunc {

	return func(c *gin.Context) {
		carrier := opentracing.HTTPHeadersCarrier(c.Request.Header)
		wireSpanCtx, _ := opentracing.GlobalTracer().Extract(opentracing.HTTPHeaders, carrier)
		// FIXME: handle err?

		serverSpan := opentracing.GlobalTracer().StartSpan("app-server-backend", ext.RPCServerOption(wireSpanCtx))
		defer serverSpan.Finish()

		c.Request = c.Request.WithContext(
			opentracing.ContextWithSpan(c.Request.Context(), serverSpan))

		// if we bring the c.Request.Context() which contains the serverSpan already
		// to the Gin's internal map by below code:
		//     c.Set("SpanContext", c.Request.Context())
		// then we don't have to Inject the server span's context to the
		// carrier as below line of code does

		// opentracing.GlobalTracer().Inject(serverSpan.Context(), opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(c.Request.Header))

		c.Set("SpanContext", c.Request.Context())

		c.Next()
	}
}
