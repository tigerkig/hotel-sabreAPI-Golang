package controller

import (
	"encoding/json"
	"net/http"
	"tp-api-common/data"
	"tp-api-common/util"
	"tp-system/model"
)

// SetDisplaySetting - Set display settings like timezone, timeformat, dateformat
func SetDisplaySetting(w http.ResponseWriter, r *http.Request) {
	util.LogIt(r, "Controller - Display_Setting - SetDisplaySetting")
	defer util.CommonDeferred(w, r, "Controller", "Display_Setting", "SetDisplaySetting")
	var reqMap data.DisplaySetting

	err := json.NewDecoder(r.Body).Decode(&reqMap)
	if err != nil {
		util.RespondBadRequest(r, w)
		return
	}

	if reqMap.DateFormat == "" || reqMap.TimeFormat == "" || reqMap.TimeZone == "" || reqMap.TimeZoneKey == "" {
		util.RespondBadRequest(r, w)
		return
	}

	err = model.SetParameter("date_format", reqMap.DateFormat)
	if err != nil {
		util.RespondWithError(r, w, "500")
		return
	}

	err = model.SetParameter("time_format", reqMap.TimeFormat)
	if err != nil {
		util.RespondWithError(r, w, "500")
		return
	}

	err = model.SetParameter("time_zone", reqMap.TimeZone)
	if err != nil {
		util.RespondWithError(r, w, "500")
		return
	}

	err = model.SetParameter("time_zone_state", reqMap.TimeZoneKey)
	if err != nil {
		util.RespondWithError(r, w, "500")
		return
	}

	util.Respond(r, w, nil, 204, "")
}

// GetDisplaySettings -  Return display setting variable like time format, date format, time zone format
func GetDisplaySettings(w http.ResponseWriter, r *http.Request) {
	util.LogIt(r, "Controller - Display_Setting - GetDisplaySettings")
	defer util.CommonDeferred(w, r, "Controller", "Display_Setting", "GetDisplaySettings")
	var retMap = make(map[string]interface{})
	retMap["time_format"], _ = model.GetParameter("time_format")
	retMap["date_format"], _ = model.GetParameter("date_format")
	retMap["time_zone"], _ = model.GetParameter("time_zone")
	retMap["time_zone_key"], _ = model.GetParameter("time_zone_state")
	util.RespondData(r, w, retMap, 200)
}

// DisplaySettingInit -  Return init list of display setting
func DisplaySettingInit(w http.ResponseWriter, r *http.Request) {
	util.LogIt(r, "Controller - Display_Setting - DisplaySettingInit")
	defer util.CommonDeferred(w, r, "Controller", "Display_Setting", "DisplaySettingInit")
	var retMap = make(map[string]interface{})
	retMap["time_format"] = util.StaticArrDateFormat
	retMap["date_format"] = util.StaticArrTimeFormat
	retMap["time_zone"] = util.TimeZone
	util.RespondData(r, w, retMap, 200)
}
