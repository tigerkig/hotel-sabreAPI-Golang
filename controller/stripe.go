package controller

import (
	"net/http"
	"tp-api-common/util"
	"tp-system/config"
	"tp-system/model"

	"github.com/gorilla/mux"
)

//ReCreateAccountLinks - Recreate account links for hotelier account
func ReCreateAccountLinks(w http.ResponseWriter, r *http.Request) {
	util.LogIt(r, "controller - controller - ReCreateAccountLinks")
	defer util.CommonDeferred(w, r, "Controller", "controller", "ReCreateAccountLinks")
	vars := mux.Vars(r)
	accID := vars["id"]
	if accID == "" {
		util.RespondBadRequest(r, w)
		return
	}

	if model.CheckConnectStripeAccExists(r, accID, true) {
		LoginLink, _ := model.CreateLoginLink(accID)
		http.Redirect(w, r, LoginLink, http.StatusSeeOther)
		return
	} else {
		onBoardLinks, err := model.AccountLinks(accID)
		if util.CheckErrorLog(r, err) {
			util.RespondBadRequest(r, w)
			return
		}
		http.Redirect(w, r, onBoardLinks, http.StatusSeeOther)
		return
	}
}

//ReturnAccountURL - Stripe connects accounts return URL
func ReturnAccountURL(w http.ResponseWriter, r *http.Request) {
	util.LogIt(r, "controller - controller - ReturnAccountURL")
	defer util.CommonDeferred(w, r, "Controller", "controller", "ReturnAccountURL")
	vars := mux.Vars(r)
	accID := vars["id"]
	if accID == "" {
		util.RespondBadRequest(r, w)
		return
	}

	if !model.CheckConnectStripeAccExists(r, accID, false) {
		util.Respond(r, w, nil, 404, "")
		return
	}

	if !model.UpdateHotelBankStatus(r, accID) {
		util.RespondBadRequest(r, w)
		return
	}

	http.Redirect(w, r, config.Env.PartnerWebURL, http.StatusSeeOther)
	return
}
