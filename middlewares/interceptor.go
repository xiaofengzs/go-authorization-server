package middlewares

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"strings"

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
		if strings.HasPrefix(c.Request.RequestURI, "/login") {
			if sessionId == "" && request.PostForm.Get("username") == "" {
				RedirectUserToLogin(rw, c, false)
				return
			}

			// authInterceptor.redisClient.Get(context.Background(), sessionId).
			value, err := authInterceptor.redisClient.Get(context.Background(), sessionId).Result()
			if err != nil && err.Error() != "redis: nil" {
				log.Printf("Get session if from redis failed, err: %s", err)
				log.Println("If connect redis error, may get session from RDS.")
			}

			if value == "" && request.PostForm.Get("username") == "" {
				RedirectUserToLogin(rw, c, false)
				return
			} else if value != "" {
				log.Println("Check expire time -> User is logging in now, go to index. otherwise, relogin.")
				RedirectUserToIndex(c)
				return
			}
		} else {
			// 在不是login的情况下，其他任何endpoint，如果没有login的情况下，是没有session id的
			// 如果session id为空，证明没有登录，需要去登录
			if sessionId == "" {
				RedirectUserToLogin(rw, c, true)
				return
			}

			// 如果session id不为空，去获取redis中的session，session如果不存在就去登录
			value, err := authInterceptor.redisClient.Get(context.Background(), sessionId).Result()
			if err != nil && err.Error() != "redis: nil" {
				// 假装有一个rds存session，并且可以查到
				log.Printf("Get session if from redis failed, err: %s", err)
				log.Println("If connect redis error, may get session from RDS.")
			}
			if value == "" {
				RedirectUserToLogin(rw, c, true)
				return
			}
		}

		// 走到这里说明 既不是login，访问其他endpoint时，session也在redis也存在
		c.Next()
	}
}

func RedirectUserToLogin(rw gin.ResponseWriter, c *gin.Context, redirectAfterLogin bool) {
	location := ""
	if redirectAfterLogin {
		url := base64.StdEncoding.EncodeToString([]byte(c.Request.URL.String()))
		location = fmt.Sprintf("?redirect_url=%s", url)
	} 

	rw.Header().Set("Content-Type", "text/html; charset=utf-8")
	rw.Write([]byte(`<h1>Login page</h1>`))
	rw.Write([]byte(fmt.Sprintf(`
						<p>Howdy! This is the log in page. For this example, it is enough to supply the username.</p>
						<form action="/login%s" method="post">
							<input type="text" name="username" /> <small>try peter</small><br>
							<input type="submit">
						</form>
					`, location)))
	
	toRedirectUrl := fmt.Sprintf("localhost:8080/login%s", location)
	c.Redirect(302, toRedirectUrl)
	c.Abort()
}

func RedirectUserToIndex(c *gin.Context) {
	c.Redirect(302, "localhost:8080/index")
	c.Abort()
}
