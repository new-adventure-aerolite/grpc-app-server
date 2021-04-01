package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/new-adventure-areolite/grpc-app-server/pd/fight"
)

// SelectHero ...
func SelectHero(fightSvcClient fight.FightSvcClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		sid := c.GetString(("id"))
		heroName := c.Query("hero")

		resp, err := fightSvcClient.SelectHero(c.Request.Context(), &fight.SelectHeroRequest{
			Id:       sid,
			HeroName: heroName,
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
