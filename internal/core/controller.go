package core

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type Controller struct {
	*gin.Engine
}

func NewController() *Controller {
	c := &Controller{
		Engine: gin.Default(),
	}
	return c
}

func (c *Controller) RegisterRoutes() {
	c.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
	c.scoreGroup()
}

func (c *Controller) scoreGroup() {
	score := c.Group("/score")
	score.GET("/new", c.newScore)
}

func (c *Controller) newScore(ctx *gin.Context) {
	ctx.JSON(http.StatusNotImplemented, gin.H{})
}
