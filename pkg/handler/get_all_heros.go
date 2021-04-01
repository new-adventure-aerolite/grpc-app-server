package handler

import (
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/new-adventure-areolite/grpc-app-server/pd/fight"
)

// GetAllHeros ...
func GetAllHeros(fightSvcClient fight.FightSvcClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		stream, err := fightSvcClient.ListHeros(c.Request.Context(), &fight.ListHerosRequest{})
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, map[string]string{
				"error": err.Error(),
			})
			return
		}

		var result []*fight.Hero

		for {
			hero, err := stream.Recv()
			if err == io.EOF {
				break
			}

			if err != nil {
				c.AbortWithStatusJSON(http.StatusInternalServerError, map[string]string{
					"error": err.Error(),
				})
				return
			}

			result = append(result, hero)
		}

		c.JSON(http.StatusOK, result)
	} // return func
}
