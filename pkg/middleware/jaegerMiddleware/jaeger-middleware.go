package jaegerMiddleware

import (
	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
)

func OpenTracingMiddleware(tracer opentracing.Tracer) gin.HandlerFunc {

	return func(c *gin.Context) {

		// var childSpan opentracing.Span

		// spCtx, err := opentracing.GlobalTracer().Extract(opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(c.Request.Header))
		span, ctx := opentracing.StartSpanFromContext(c, c.Request.URL.Path)
		// if err != nil {
		// 	// childSpan = tracer.StartSpan(c.Request.URL.Path)
		// 	opentracing.StartSpanFromContext()
		// 	defer childSpan.Finish()
		// } else {
		// 	childSpan = opentracing.StartSpan(
		// 		c.Request.URL.Path,
		// 		opentracing.ChildOf(spCtx),
		// 		opentracing.Tag{Key: string(ext.Component), Value: "HTTP"},
		// 		ext.SpanKindRPCServer,
		// 	)
		// 	defer childSpan.Finish()
		// }
		defer span.Finish()

		opentracing.GlobalTracer().Inject(span.Context(), opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(c.Request.Header))
		// c.Set("Tracer", tracer)
		c.Set("SpanContext", ctx)
		//injectErr := childSpan.Tracer().Inject(parentSpan.Context(), opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(c.Request.Header))
		// if injectErr != nil {
		// 	klog.Fatalf("%s: Couldn't inject headers", err)
		// }
		// c.Request = c.Request.WithContext(opentracing.ContextWithSpan(c.Request.Context(), parentSpan))
		// c.Next()
		// injectErr := tracer.Inject(parentSpan.Context(), opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(c.Request.Header))
		// if injectErr != nil {
		// 	log.Fatalf("%s: Couldn't inject headers", err)
		// }

		// c.Request = c.Request.WithContext(opentracing.ContextWithSpan(c.Request.Context(), childSpan))
		c.Next()
	}
}

// func AfterOpenTracingMiddleware(tracer opentracing.Tracer) gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		if span := opentracing.SpanFromContext(c); span != nil {
// 			span := tracer.StartSpan("xxx", opentracing.ChildOf(span.Context()))
// 			defer span.Finish()
// 			//c = opentracing.ContextWithSpan(c, span)
// 			c.Request = c.Request.WithContext(opentracing.ContextWithSpan(c.Request.Context(), span))
// 		} else {
// 			log.Println("xxxxxxxx, span is nil")
// 		}
// 	}
// }
