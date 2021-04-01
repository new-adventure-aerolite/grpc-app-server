package handler

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/new-adventure-areolite/grpc-app-server/pd/fight"
	"google.golang.org/grpc/metadata"
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

func ClearSession(fightSvcClient fight.FightSvcClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		ID := c.GetString("id")
		resp, err := fightSvcClient.ClearSession(context.Background(), &fight.ClearSessionRequest{
			Id: ID,
		})
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, map[string]string{
				"error": err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, resp)
	}
}

// game ...
func game(fightSvcClient fight.FightSvcClient, eventType fight.Type) gin.HandlerFunc {
	return func(c *gin.Context) {
		sid := c.GetString(("id"))
		ctx := metadata.NewOutgoingContext(c.Request.Context(), metadata.Pairs(
			"user-type", c.GetString("user-type"),
		))

		resp, err := fightSvcClient.Game(ctx, &fight.GameRequest{
			Type: eventType,
			Id:   sid,
		})

		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, map[string]string{
				"error": err.Error(),
			})
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
			c.AbortWithStatusJSON(http.StatusNotFound, map[string]string{
				"error": fmt.Sprintf("event type: '%T' doesn't exist", fight.Type_name[int32(resp.Type)]),
			})
			return
		}
	} // return func
}
