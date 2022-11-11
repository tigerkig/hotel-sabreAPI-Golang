package middleware

import (
	"fmt"
	"net"
	"net/http"
	"tp-api-common/util"
	"tp-system/model"

	"github.com/gorilla/context"
)

// PartnerMiddleware - Authentication Partner Panel API's Middleware.
func PartnerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip, _, _ := net.SplitHostPort(r.RemoteAddr)
		context.Set(r, "Visitor_IP", ip)
		token := r.Header.Get("X-Auth-Token")
		SessionData, err := model.GetConsoleAuthDetails(r, token)
		if err != nil {
			util.RespondWithError(r, w, "500")
			return
		}

		if len(SessionData) > 0 && SessionData["token"] == token {
			context.Set(r, "Request-Token", token)
			context.Set(r, "UserId", SessionData["id"])
			context.Set(r, "HotelId", SessionData["hotel_id"])
			context.Set(r, "GroupId", SessionData["group_id"])
			context.Set(r, "BusinessName", SessionData["hotel_name"])
			context.Set(r, "Side", "TP-PARTNER")
			context.Set(r, "Username", SessionData["username"])
			RandVal := util.RandStringBytesMask(r, 8)
			context.Set(r, "ReqToken", RandVal)

			next.ServeHTTP(w, r)
		} else {
			util.Respond(r, w, nil, 403, "")
			return
		}
	})
}

// FormDataPartnerAPIMiddleware - Partner Panel Form Data API
func FormDataPartnerAPIMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var err error
		ip, _, _ := net.SplitHostPort(r.RemoteAddr)
		context.Set(r, "Visitor_IP", ip)
		token := r.Header.Get("X-Auth-Token")
		SessionData, err := model.GetConsoleAuthDetails(r, token)
		if err != nil {
			util.RespondWithError(r, w, "500")
			return
		}
		if len(SessionData) > 0 && SessionData["token"] == token {
			context.Set(r, "Request-Token", token)
			context.Set(r, "UserId", SessionData["id"])
			context.Set(r, "HotelId", SessionData["hotel_id"])
			context.Set(r, "BusinessName", SessionData["hotel_name"])
			context.Set(r, "Side", "TP-PARTNER")
			context.Set(r, "Username", SessionData["username"])
			RandVal := util.RandStringBytesMask(r, 8)
			context.Set(r, "ReqToken", RandVal)
			next.ServeHTTP(w, r)
		} else {
			util.LogIt(r, fmt.Sprint("Firewall - Token Expired - ", token))
			util.Respond(r, w, nil, 403, "")
			return
		}
	})
}
