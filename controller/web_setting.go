package controller

import (
	"net/http"
	"tp-api-common/util"
	"tp-system/model"
)

// GetWebDefaultSettings - Returns all web front filtration and sorting settings
func GetWebDefaultSettings(w http.ResponseWriter, r *http.Request) {
	util.LogIt(r, "Controller - web_setting - GetWebDefaultSettings")
	defer util.CommonDeferred(w, r, "Controller", "web_setting", "GetWebDefaultSettings")

	retBody, err := model.GetWebDefaultSettings(r)
	if err != nil {
		util.RespondWithError(r, w, "500")
		return
	}
	util.RespondData(r, w, retBody, 200)
}

// UpdateWebDefaultSettings - Updates all web front filtration and sorting settings
func UpdateWebDefaultSettings(w http.ResponseWriter, r *http.Request) {

	util.LogIt(r, "Controller - web_setting - UpdateWebDefaultSettings")
	defer util.CommonDeferred(w, r, "Controller", "web_setting", "UpdateWebDefaultSettings")

	reqBody, err := util.ExtractRequestBody(r)
	if err != nil {
		util.RespondBadRequest(r, w)
		return
	}

	validateFlag := ValidateNotNullAndString(reqBody, []string{"default_adult", "default_checkin_checkout_range", "default_checkin_day", "default_child", "default_room"})
	validateFlag1 := ValidateNotNullAndString(reqBody, []string{"filter_max_price", "filter_min_price"})
	validateFlag2 := ValidateNotNullAndString(reqBody, []string{"show_amenities", "show_customer_ratings", "show_hotel_ratings", "show_popular_filtration", "show_price_range", "show_property_type"})

	if validateFlag == 0 || validateFlag1 == 0 || validateFlag2 == 0 {
		util.RespondBadRequest(r, w)
		return
	}

	flg := model.UpdateWebDefaultSettings(r, reqBody)
	if !flg {
		util.RespondWithError(r, w, "500")
		return
	}

	util.Respond(r, w, nil, 204, "")
}
