package middleware

import (
	"fmt"
	"net"
	"net/http"
	"tp-api-common/util"
	"tp-system/config"
	"tp-system/model"

	"github.com/gorilla/context"
)

// WebMiddleware - Authentication API's Middleware. Body in JSON format is mandatory. Without body API is not allowed, it will throw error.
func WebMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip, _, _ := net.SplitHostPort(r.RemoteAddr)
		context.Set(r, "Visitor_IP", ip)
		context.Set(r, "Side", "TP-FRONT")
		context.Set(r, "UserId", "web-user")
		context.Set(r, "Request-Token", "web-token")

		next.ServeHTTP(w, r)
	})
}

// APIMiddleware - Authentication Admin API's Middleware.
func APIMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip, _, _ := net.SplitHostPort(r.RemoteAddr)
		context.Set(r, "Visitor_IP", ip)
		token := r.Header.Get("X-Auth-Token")
		SessionData, err := model.GetAuthDetails(r, token)
		if err != nil {
			util.RespondWithError(r, w, "500")
			return
		}

		if len(SessionData) > 0 && SessionData["token"] == token {
			context.Set(r, "Request-Token", token)
			context.Set(r, "UserId", SessionData["id"])
			context.Set(r, "Side", "TP-BACKOFFICE")
			context.Set(r, "Username", SessionData["username"])
			context.Set(r, "Privileges", SessionData["privileges"])
			RandVal := util.RandStringBytesMask(r, 8)
			context.Set(r, "ReqToken", RandVal)

			next.ServeHTTP(w, r)
		} else {
			util.Respond(r, w, nil, 403, "")
			return
		}
	})
}

// FormDataAPIMiddleware - Authentication Admin API's Middleware.
func FormDataAPIMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip, _, _ := net.SplitHostPort(r.RemoteAddr)
		context.Set(r, "Visitor_IP", ip)
		token := r.Header.Get("X-Auth-Token")
		SessionData, err := model.GetAuthDetails(r, token)
		if err != nil {
			util.RespondWithError(r, w, "500")
			return
		}

		if len(SessionData) > 0 && SessionData["token"] == token {
			context.Set(r, "Request-Token", token)
			context.Set(r, "UserId", SessionData["id"])
			context.Set(r, "Side", "TP-BACKOFFICE")
			context.Set(r, "Username", SessionData["username"])
			context.Set(r, "Privileges", SessionData["privileges"])
			RandVal := util.RandStringBytesMask(r, 8)
			context.Set(r, "ReqToken", RandVal)

			next.ServeHTTP(w, r)
		} else {
			util.LogIt(r, fmt.Sprint("FormDataAPIMiddleware Firewall - Token Expired - ", token))
			util.Respond(r, w, nil, 403, "")
			return
		}
	})
}

// CrossMiddleware - Authentication Cross API's Middleware With Static Token.
func CrossMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip, _, _ := net.SplitHostPort(r.RemoteAddr)
		context.Set(r, "Visitor_IP", ip)
		token := r.Header.Get("X-Auth-Token")
		if token == config.Env.InvAuthKey {
			next.ServeHTTP(w, r)
		} else {
			util.Respond(r, w, nil, 403, "")
			return
		}
	})
}
