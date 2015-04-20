package handlers

import (
	"net/http"
)

func UserAuthorizationMiddleware(req *http.Request) {
	// This isn't a real authorization middleware.
	// It will act like one, in the sense that if will look for some key,
	// and use it to restrict access to certain resources
	userKey := req.Header.Get("X-USER-API-KEY")
	if userKey != "" {
		userId = userKey
	} else {
		userId = "0"
	}
}
