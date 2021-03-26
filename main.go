package main

import (
	"flag"

	"github.com/TianqiuHuang/grpc-client-app/pd/fight"
	"github.com/TianqiuHuang/grpc-client-app/pkg/handler"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"k8s.io/klog"
)

var (
	port string
	addr string
)

func init() {
	flag.StringVar(&port, "port", "8000", "listen port")
	flag.StringVar(&addr, "addr", "127.0.0.1:8001", "fight svc addr")
}

func main() {
	flag.Parse()

	gin.DisableConsoleColor()
	r := gin.Default()
	// gin.SetMode(gin.ReleaseMode)

	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		klog.Fatal(err)
	}
	fightSvcClient := fight.NewFightSvcClient(conn)

	r.GET("/heros", handler.GetAllHeros(fightSvcClient))
	r.GET("/session/:id", handler.LoadSession(fightSvcClient))
	r.PUT("/session/:id", handler.SelectHero(fightSvcClient))
	r.PUT("/session/:id/fight", handler.Fight(fightSvcClient))
	r.POST("/session/:id/archive", handler.Archive(fightSvcClient))
	r.POST("/session/:id/level", handler.Level(fightSvcClient))
	r.POST("/session/:id/quit", handler.Quit(fightSvcClient))

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
	r.POST("/admin/hero", handler.CreateHero())
	r.PUT("/admin/hero", handler.AdjustHero())

	r.Run(":" + port)
}
