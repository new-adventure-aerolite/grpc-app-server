package handler

import (
	"context"
	"net/http"

	"github.com/TianqiuHuang/grpc-client-app/pd/fight"
	"github.com/gin-gonic/gin"
)

// SelectHero ...
func SelectHero(fightSvcClient fight.FightSvcClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		sid := c.GetString(("id"))
		heroName := c.Query("hero")

		resp, err := fightSvcClient.SelectHero(context.Background(), &fight.SelectHeroRequest{
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
