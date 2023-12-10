package web

import (
	"encoding/base64"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/xiaofengzs/go-authorization-server/middlewares"
)

func RegisterHttpRequestHandlers(router *gin.Engine, redisClient *redis.Client) {

	sessionInterceptor := middlewares.NewAuthInterceptor(redisClient).GetSessionInterceptor()

	router.Use(sessionInterceptor)
	router.POST("/login", func(ctx *gin.Context) {
		sessionId := uuid.New().String()
		redisClient.Set(ctx.Request.Context(), sessionId, sessionId, time.Minute * 30)
		ctx.SetCookie("sessionId", sessionId, 1800, "", "localhost", false, false )
		
		
		redirectUrl := ctx.Query("redirect_url")
		location, err := base64.StdEncoding.DecodeString(redirectUrl)
	
		if err != nil {
			ctx.String(http.StatusBadRequest, "Invalid base64 data")
			return
		}

		ctx.Redirect(302, string(location))
	})

	router.GET("/index", func(ctx *gin.Context) {
		ctx.JSON(200, "hello world")
	})

}