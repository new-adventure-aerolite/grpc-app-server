package jaegerMiddleware

import (
	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
)

func OpenTracingMiddleware() gin.HandlerFunc {

	return func(c *gin.Context) {
		span, ctx := opentracing.StartSpanFromContext(c, c.Request.URL.Path)

		defer span.Finish()

		opentracing.GlobalTracer().Inject(span.Context(), opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(c.Request.Header))
		c.Set("SpanContext", ctx)

		c.Next()
	}
}
