package main

import (
	"context"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	authorizationserver "github.com/xiaofengzs/go-authorization-server/authorization_server"
	"github.com/xiaofengzs/go-authorization-server/cache"
	"github.com/xiaofengzs/go-authorization-server/web"
)

func main() {

	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	val, err := client.Get(context.Background(), "").Result()

	// err := client.Set(context.Background(), "person", nil, 0).Err()
	if err != nil {
		log.Println(err)
	}

	// // val, err := client.Get(context.Background(), "person").Result()
	// if err != nil {
	// 	log.Println(err)
	// }
	fmt.Println(val)

	redisClient := cache.BuildRedisClient()
	router := gin.Default()
	web.RegisterHttpRequestHandlers(router, redisClient)
	web.RegisterClientHandlers(router)

	// ### oauth2 server ###
	authorizationserver.RegisterOAuth2Handlers(router) // the authorization server (fosite)

	router.Run()
}
