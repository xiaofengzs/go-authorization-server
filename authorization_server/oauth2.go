package authorizationserver

import (
	"crypto/rand"
	"crypto/rsa"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ory/fosite"
	"github.com/ory/fosite/compose"
	"github.com/ory/fosite/storage"
)

var (
	secret = []byte("some-cool-secret-that-is-32bytes")

	config = &fosite.Config{
		AccessTokenLifespan: time.Minute * 30,
		GlobalSecret:        secret,
	}

	store = storage.NewExampleStore()

	privateKey, _ = rsa.GenerateKey(rand.Reader, 2048)

	strategy = compose.NewOAuth2HMACStrategy(config)

	oauth2 = compose.Compose(config, store, strategy, compose.OAuth2AuthorizeExplicitFactory)
)
func RegisterOAuth2Handlers(router *gin.Engine) {
	router.GET("/oauth2/auth", func(ctx *gin.Context) {
		authEndpoint(ctx.Writer, ctx.Request)
	})

	router.POST("oauth2/token", func(ctx *gin.Context) {
		tokenEndpoint(ctx.Writer, ctx.Request)
	})
}

