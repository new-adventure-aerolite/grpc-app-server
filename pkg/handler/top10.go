package handler

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/TianqiuHuang/grpc-client-app/pd/fight"
	"github.com/gin-gonic/gin"
)

var (
	top10Players       = []*fight.Top10Response_Player{}
	foreverTop10Stream fight.FightSvc_Top10Client
)

// InitTop10Client ...
func InitTop10Client(fightSvcClient fight.FightSvcClient) error {
	var err error
	foreverTop10Stream, err = fightSvcClient.Top10(context.Background(), &fight.Top10Request{})
	if err != nil {
		return err
	}
	for {
		top10, err := foreverTop10Stream.Recv()
		if err == io.EOF {
			return fmt.Errorf("connection is closed by the server")
		}
		if err != nil {
			return err
		}
		top10Players = top10.GetPlayers()
	}
}

// Top10 ...
func Top10() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, top10Players)
	}
}
