package middlewares

import (
	"context"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

type AuthInterceptor struct {
	redisClient *redis.Client
}

func NewAuthInterceptor(redisClient *redis.Client) AuthInterceptor {
	return AuthInterceptor{redisClient}
}

func (authInterceptor AuthInterceptor) GetSessionInterceptor() gin.HandlerFunc {
	return func(c *gin.Context) {
		request := c.Request
		request.ParseForm()
		rw := c.Writer
		sessionId, _ := c.Cookie("sessionId")

		log.Printf("full path: %s", c.Request.RequestURI)
		if c.Request.RequestURI == "/login" {
			if sessionId == "" && request.PostForm.Get("username") == "" {
				RedirectUserToLogin(rw, c)
				return
			}
	
			// authInterceptor.redisClient.Get(context.Background(), sessionId).
			value, err := authInterceptor.redisClient.Get(context.Background(), sessionId).Result()
			if err != nil && err.Error() != "redis: nil" {
				log.Printf("Get session if from redis failed, err: %s", err)
				log.Println("If connect redis error, may get session from RDS.")
			}

			if value == "" && request.PostForm.Get("username") == "" {
				RedirectUserToLogin(rw, c)
				return
			} else if value != "" {
				log.Println("Check expire time -> User is logging in now, go to index. otherwise, relogin.")
				RedirectUserToIndex(c)
				return
			}
		} else {
			value, err := authInterceptor.redisClient.Get(context.Background(), sessionId).Result()
			if err != nil && err.Error() != "redis: nil" {
				log.Printf("Get session if from redis failed, err: %s", err)
				log.Println("If connect redis error, may get session from RDS.")
			}
			if (value == "") {
				RedirectUserToLogin(rw, c)
			}
		}
		// if no session id, go to login page
		// has session id
		// 1. session id is not in redis, if not login path, go to login page
		// 2. session id is not in redis, if login path, check username, if username == "", go to login page
		// if username != "", go to login handler,
		c.Next()
	}
}

func RedirectUserToLogin(rw gin.ResponseWriter, c *gin.Context) {
	rw.Header().Set("Content-Type", "text/html; charset=utf-8")
	rw.Write([]byte(`<h1>Login page</h1>`))
	rw.Write([]byte(fmt.Sprintf(`
						<p>Howdy! This is the log in page. For this example, it is enough to supply the username.</p>
						<form action="/login" method="post">
							<input type="text" name="username" /> <small>try peter</small><br>
							<input type="submit">
						</form>
					`)))
	c.Redirect(302, "localhost:8080/login")
	c.Abort()
}

func RedirectUserToIndex(c *gin.Context) {
	c.Redirect(302, "localhost:8080/index")
	c.Abort()
}
