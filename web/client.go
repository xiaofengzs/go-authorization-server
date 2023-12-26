package web

import (
	"github.com/gin-gonic/gin"
	"github.com/xiaofengzs/go-authorization-server/client"
	"net/http"
)

type ClientRequest struct {
	Id string `json:"id"`
}

func RegisterClientHandlers(router *gin.Engine) {
	router.POST("/client", func(ctx *gin.Context) {
		var body ClientRequest
		if err := ctx.ShouldBindJSON(&body); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		clientManagement := client.NewClientManagement()
		clientManagement.SaveClient(ctx.Request.Context(), body.Id)

		ctx.Status(http.StatusCreated)
	})
}
