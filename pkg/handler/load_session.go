package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/new-adventure-areolite/grpc-app-server/pd/fight"
)

// LoadSession ...
func LoadSession(fightSvcClient fight.FightSvcClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetString("id")

		resp, err := fightSvcClient.LoadSession(c.Request.Context(), &fight.LoadSessionRequest{
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
