package main

import (
	"flag"
	"log"
	"time"

	"github.com/TianqiuHuang/grpc-client-app/pd/auth"
	"github.com/TianqiuHuang/grpc-client-app/pd/fight"
	auth_middle_ware "github.com/TianqiuHuang/grpc-client-app/pkg/auth"
	"github.com/TianqiuHuang/grpc-client-app/pkg/handler"
	"github.com/TianqiuHuang/grpc-client-app/pkg/istio"
	"github.com/TianqiuHuang/grpc-client-app/pkg/jaeger_service"
	"github.com/TianqiuHuang/grpc-client-app/pkg/middleware/jaegerMiddleware"
	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
	"google.golang.org/grpc"
	"k8s.io/klog"
)

var (
	port           string
	addr           string
	authServerAddr string
	tlsCert        string
	tlsKey         string
)

func init() {
	flag.StringVar(&port, "port", "8000", "listen port")
	flag.StringVar(&addr, "addr", "127.0.0.1:8001", "fight svc addr")
	flag.StringVar(&authServerAddr, "auth-server-addr", "127.0.0.1:6666", "auth svc addr")
	flag.StringVar(&tlsCert, "tls-cert", "", "tls cert")
	flag.StringVar(&tlsKey, "tls-key", "", "tls key")
}

func main() {
	flag.Parse()

	proxy := istio.New(10, 10*time.Second, 30*time.Second)
	if err := proxy.Wait(); err != nil {
		klog.Fatal(err)
	}

	defer func() {
		if err := proxy.Close(); err != nil {
			klog.Error(err)
		}
	}()

	gin.DisableConsoleColor()
	r := gin.Default()
	gin.SetMode(gin.DebugMode)

	// new jaeger tracer
	tracer, _, err := jaeger_service.NewJaegerTracer("grpc-app-server", "127.0.0.1:6831")
	if err != nil {
		klog.Fatal(err)
	}

	opentracing.SetGlobalTracer(tracer)

	// add openTracing middleware
	r.Use(jaegerMiddleware.OpenTracingMiddleware())
	// r.Use(jaegerMiddleware.AfterOpenTracingMiddleware(tracer))

	// trace on grpc client
	dialOpts := []grpc.DialOption{grpc.WithInsecure()}
	if tracer != nil {
		dialOpts = append(dialOpts, grpc.WithUnaryInterceptor(jaeger_service.ClientInterceptor(tracer, "call fight gRPC")))
	} else {
		log.Fatal("tracer is nil, exist")
	}

	// create fight connection
	conn, err := grpc.Dial(addr, dialOpts...)
	if err != nil {
		klog.Fatal(err)
	}
	fightSvcClient := fight.NewFightSvcClient(conn)

	// create auth connection
	authConn, err := grpc.Dial(authServerAddr, grpc.WithInsecure(), grpc.WithUnaryInterceptor(jaeger_service.ClientInterceptor(tracer, "call auth gRPC")))
	if err != nil {
		klog.Fatal(err)
	}

	authSvcClient := auth.NewAuthServiceClient(authConn)
	authClient := auth_middle_ware.New(authSvcClient)
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET,HEAD,OPTIONS,POST,PUT")
		c.Writer.Header().Set("Access-Control-Allow-Headers",
			"Access-Control-Allow-Headers, Origin,Accept, X-Requested-With, Content-Type, Authorization, Access-Control-Request-Method, Access-Control-Request-Headers")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(200)
		}
		c.Next()
	})

	group := r.Group("/", auth_middle_ware.AuthMiddleWare(authClient))
	group.GET("/heros", handler.GetAllHeros(fightSvcClient))
	group.GET("/session", handler.LoadSession(fightSvcClient))
	group.PUT("/session", handler.SelectHero(fightSvcClient))
	group.PUT("/session/fight", handler.Fight(fightSvcClient))
	group.POST("/session/archive", handler.Archive(fightSvcClient))
	group.POST("/session/level", handler.Level(fightSvcClient))
	group.POST("/session/quit", handler.Quit(fightSvcClient))
	group.POST("/session/clear", handler.ClearSession(fightSvcClient))

	go func() {
		if err := handler.InitTop10Client(fightSvcClient); err != nil {
			klog.Warning(err)
		}
	}()

	r.GET("/top10", handler.Top10())

	// Admin rest api
	go func() {
		if err := handler.InitAdminClient(fightSvcClient); err != nil {
			klog.Warning(err)
		}
	}()
	adminGroup := r.Group("/admin", auth_middle_ware.AdminAuthMiddleWare(authClient))
	adminGroup.POST("/hero", handler.CreateHero())
	adminGroup.PUT("/hero", handler.AdjustHero())

	if tlsCert != "" && tlsKey != "" {
		r.RunTLS(":"+port, tlsCert, tlsKey)
	} else {
		r.Run(":" + port)
	}
}
