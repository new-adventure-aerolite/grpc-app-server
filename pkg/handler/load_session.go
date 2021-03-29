package handler

import (
	"context"
	"net/http"

	"github.com/TianqiuHuang/grpc-client-app/pd/fight"
	"github.com/gin-gonic/gin"
)

// LoadSession ...
func LoadSession(fightSvcClient fight.FightSvcClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetString(("id"))

		ctx, _ := c.Get("SpanContext")
		resp, err := fightSvcClient.LoadSession(ctx.(context.Context), &fight.LoadSessionRequest{
			Id: userID,
		})

		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, map[string]string{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, resp)
	} // return func
}
