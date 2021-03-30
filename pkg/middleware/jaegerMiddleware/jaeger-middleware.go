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
		// log.Println(parentSpanCtx, err)

		serverSpan := opentracing.GlobalTracer().StartSpan("rpc-app-server", ext.RPCServerOption(wireSpanCtx))
		defer serverSpan.Finish()

		//span, ctx := opentracing.StartSpanFromContext(c.Request.Context(), c.Request.URL.Path)
		//defer span.Finish()

		opentracing.GlobalTracer().Inject(serverSpan.Context(), opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(c.Request.Header))
		c.Request = c.Request.WithContext(
			opentracing.ContextWithSpan(c.Request.Context(), serverSpan))

		//c.Set("SpanContext", c.Request.Context())

		c.Next()
	}
}
