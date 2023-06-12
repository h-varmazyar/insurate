package core

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type Controller struct {
	*gin.Engine
	service *Service
}

func NewController(service *Service) *Controller {
	c := &Controller{
		Engine:  gin.Default(),
		service: service,
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
	score.GET("/:score_id/download", c.downloadScoreReport)
}

func (c *Controller) newScore(ctx *gin.Context) {
	req := new(NewScoreReq)
	err := ctx.ShouldBind(req)
	if err != nil {
		ctx.String(http.StatusBadRequest, err.Error())
		return
	}
	resp, err := c.service.NewScore(ctx, req)
	if err != nil {
		ctx.String(http.StatusNotAcceptable, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, resp)
}

func (c *Controller) downloadScoreReport(ctx *gin.Context) {
	scoreID := ctx.Param("score_id")

	score, err := c.service.DownloadScore(ctx, &DownloadScoreReq{ScoreID: scoreID})
	if err != nil {
		ctx.String(http.StatusNotAcceptable, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, score.Score)
}
