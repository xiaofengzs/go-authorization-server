package authorizationserver

import (
	"log"
	"net/http"
	"time"

	"github.com/ory/fosite/handler/openid"
	"github.com/ory/fosite/token/jwt"
)

var sessions map[string]bool = make(map[string]bool)

func authEndpoint(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	authorizeRequest, err := oauth2.NewAuthorizeRequest(ctx, r)
	if err != nil {
		log.Printf("Error in building AuthorizeRequest, err: %s", err)
		oauth2.WriteAuthorizeError(ctx, w, authorizeRequest, err)
		return
	}

	// for _, scope := range r.PostForm["scopes"] {
	// 	authorizeRequest.GrantScope(scope)
	// }

	cookie, err := r.Cookie("sessionId")
	mySessionData := newSession(cookie.Value)

	response, err := oauth2.NewAuthorizeResponse(ctx, authorizeRequest, mySessionData)

	if err != nil {
		log.Printf("Error occurred in NewAuthorizeResponse: %+v", err)
		oauth2.WriteAuthorizeError(ctx, w, authorizeRequest, err)
		return
	}

	// Last but not least, send the response!
	oauth2.WriteAuthorizeResponse(ctx, w, authorizeRequest, response)

}

func newSession(user string) *openid.DefaultSession {
	return &openid.DefaultSession{
		Claims: &jwt.IDTokenClaims{
			Issuer:      "https://fosite.my-application.com",
			Subject:     user,
			Audience:    []string{"https://my-client.my-application.com"},
			ExpiresAt:   time.Now().Add(time.Hour * 6),
			IssuedAt:    time.Now(),
			RequestedAt: time.Now(),
			AuthTime:    time.Now(),
		},
		Headers: &jwt.Headers{
			Extra: make(map[string]interface{}),
		},
	}
}
