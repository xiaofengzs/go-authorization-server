package authorizationserver

import (
	"fmt"
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
	if (err != nil) {
		log.Printf("Error in building AuthorizeRequest, err: %s", err)
		oauth2.WriteAuthorizeError(ctx, w, authorizeRequest, err)
		return 
	}
	
	if !checkIfSessionValidated(w, r) {
		return
	}

	for _, scope := range r.PostForm["scopes"] {
		authorizeRequest.GrantScope(scope)
	}

	mySessionData := newSession("peter")

	response, err := oauth2.NewAuthorizeResponse(ctx, authorizeRequest, mySessionData)

	if err != nil {
		log.Printf("Error occurred in NewAuthorizeResponse: %+v", err)
		oauth2.WriteAuthorizeError(ctx, w, authorizeRequest, err)
		return
	}

	// Last but not least, send the response!
	oauth2.WriteAuthorizeResponse(ctx, w, authorizeRequest, response)

}

func checkIfSessionValidated(w http.ResponseWriter, request *http.Request) bool {
	cookie, err := request.Cookie("sessionId")
	if err != nil {
		log.Println("Not found the session for ", request.Header)
		buildLoginPage(w, request)
		return false
	}

	// Get the session ID from the cookie value
	sessionID := cookie.Value
	if (sessionID == "") {
		buildLoginPage(w, request)
	}
	_, ok := sessions[sessionID]
	if (!ok) {
		buildLoginPage(w, request)
		return false
	}
	return true
}

func buildLoginPage(w http.ResponseWriter, request *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write([]byte(`<h1>Login page</h1>`))
		w.Write([]byte(fmt.Sprintf(`
			<p>Howdy! This is the log in page. For this example, it is enough to supply the username.</p>
			<form method="post">
				<p>
					By logging in, you consent to grant these scopes:
					<ul>%s</ul>
				</p>
				<input type="text" name="username" /> <small>try peter</small><br>
				<input type="submit">
			</form>
		`, "profile")))
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