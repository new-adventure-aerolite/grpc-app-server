package main

import (
	"flag"
	"fmt"
	"runtime"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/new-adventure-areolite/grpc-app-server/pd/auth"
	"github.com/new-adventure-areolite/grpc-app-server/pd/fight"
	auth_middle_ware "github.com/new-adventure-areolite/grpc-app-server/pkg/auth"
	"github.com/new-adventure-areolite/grpc-app-server/pkg/handler"
	"github.com/new-adventure-areolite/grpc-app-server/pkg/istio"
	"github.com/new-adventure-areolite/grpc-app-server/pkg/jaeger_service"
	"github.com/new-adventure-areolite/grpc-app-server/pkg/middleware/jaegerMiddleware"
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

// Compiling info.
var (
	branch    = "no version provided"
	buildTime = "no build time provided"
	commit    = "no git commit hash provided"
)

func init() {
	flag.StringVar(&port, "port", "8000", "listen port")
	flag.StringVar(&addr, "addr", "127.0.0.1:8001", "fight svc addr")
	flag.StringVar(&authServerAddr, "auth-server-addr", "127.0.0.1:6666", "auth svc addr")
	flag.StringVar(&tlsCert, "tls-cert", "", "tls cert")
	flag.StringVar(&tlsKey, "tls-key", "", "tls key")
}

func main() {
	printVersion()
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
	tracer, _, err := jaeger_service.NewJaegerTracer("app-server", "jaeger-collector.istio-system.svc.cluster.local:14268")
	if err != nil {
		klog.Fatal(err)
	}

	// add openTracing middleware
	r.Use(jaegerMiddleware.OpenTracingMiddleware())
	// r.Use(jaegerMiddleware.AfterOpenTracingMiddleware(tracer))

	// trace on grpc client
	dialOpts := []grpc.DialOption{grpc.WithInsecure()}
	if tracer != nil {
		dialOpts = append(dialOpts, grpc.WithUnaryInterceptor(jaeger_service.ClientInterceptor(tracer, "call fight gRPC client")))
	} else {
		klog.Fatal("tracer is nil, exist")
	}

	// create fight connection
	conn, err := grpc.Dial(addr, dialOpts...)
	if err != nil {
		klog.Fatal(err)
	}
	fightSvcClient := fight.NewFightSvcClient(conn)

	// create auth connection
	authConn, err := grpc.Dial(authServerAddr, grpc.WithInsecure(), grpc.WithUnaryInterceptor(jaeger_service.ClientInterceptor(tracer, "call auth gRPC client")))
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

func printVersion() {
	fmt.Printf("%-20s %s\n", "branch", branch)
	fmt.Printf("%-20s %s\n", "git commit hash", commit)
	fmt.Printf("%-20s %s\n", "go", runtime.Version())
	fmt.Printf("%-20s %s\n", "os", runtime.GOOS)
	fmt.Printf("%-20s %s\n", "arch", runtime.GOARCH)
	fmt.Printf("%-20s %s\n", "build", buildTime)
}
