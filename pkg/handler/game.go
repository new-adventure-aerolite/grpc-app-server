package handler

import (
	"context"
	"fmt"
	"net/http"

	"github.com/TianqiuHuang/grpc-client-app/pd/fight"
	"github.com/gin-gonic/gin"
)

// Fight ...
func Fight(fightSvcClient fight.FightSvcClient) gin.HandlerFunc {
	return game(fightSvcClient, fight.Type_FIGHT)
}

// Quit ...
func Quit(fightSvcClient fight.FightSvcClient) gin.HandlerFunc {
	return game(fightSvcClient, fight.Type_QUIT)
}

// Archive ...
func Archive(fightSvcClient fight.FightSvcClient) gin.HandlerFunc {
	return game(fightSvcClient, fight.Type_ARCHIVE)
}

// Level ...
func Level(fightSvcClient fight.FightSvcClient) gin.HandlerFunc {
	return game(fightSvcClient, fight.Type_LEVEL)
}

// game ...
func game(fightSvcClient fight.FightSvcClient, eventType fight.Type) gin.HandlerFunc {
	return func(c *gin.Context) {
		sid := c.Params.ByName("id")

		resp, err := fightSvcClient.Game(context.Background(), &fight.GameRequest{
			Type: eventType,
			Id:   sid,
		})

		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		switch resp.Type {
		case fight.Type_ARCHIVE:
			c.JSON(http.StatusOK, resp.GetArchive())
		case fight.Type_FIGHT:
			c.JSON(http.StatusOK, resp.GetFight())
		case fight.Type_LEVEL:
			c.JSON(http.StatusOK, resp.GetLevel())
		case fight.Type_QUIT:
			c.JSON(http.StatusOK, resp.GetQuit())
		default:
			c.AbortWithError(http.StatusNotFound, fmt.Errorf("event type: '%T' doesn't exist", fight.Type_name[int32(resp.Type)]))
			return
		}
	} // return func
}
