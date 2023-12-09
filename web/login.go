package web

import (
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
		redisClient.Set(ctx.Request.Context(), sessionId, nil, time.Minute * 30)
		ctx.SetCookie("sessionId", sessionId, 0, "", "localhost", false, false )
	})

	router.GET("/index", func(ctx *gin.Context) {
		ctx.JSON(200, "hello world")
	})

}