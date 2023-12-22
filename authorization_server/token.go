package authorizationserver

import (
	"log"
	"net/http"
)

func tokenEndpoint(rw http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	mySessionData := newSession("")
	accessRequest, err := oauth2.NewAccessRequest(ctx, req, mySessionData)

	if err != nil {
		log.Printf("Error occurred in NewAccessRequest: %+v", err)
		oauth2.WriteAccessError(ctx, rw, accessRequest, err)
		return
	}

	response, err := oauth2.NewAccessResponse(ctx, accessRequest)
	if err != nil {
		log.Printf("Error occurred in NewAccessResponse: %+v", err)
		oauth2.WriteAccessError(ctx, rw, accessRequest, err)
		return
	}

	oauth2.WriteAccessResponse(ctx, rw, accessRequest, response)
}

