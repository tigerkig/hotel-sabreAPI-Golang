package controller

import (
	"net/http"
	"tp-api-common/util"
	"tp-system/config"
	"tp-system/model"

	"github.com/gorilla/mux"
)

// GetFrontSettings - Returns all front filtration and sorting settings
func GetFrontSettings(w http.ResponseWriter, r *http.Request) {
	util.LogIt(r, "Controller - front - GetFrontSettings")
	defer util.CommonDeferred(w, r, "Controller", "front", "GetFrontSettings")

	retBody, err := model.GetFrontSettings(r)
	if err != nil {
		util.RespondWithError(r, w, "500")
		return
	}
	util.RespondData(r, w, retBody, 200)
}

// GetHomePageData - Returns all front filtration and sorting settings
func GetHomePageData(w http.ResponseWriter, r *http.Request) {
	util.LogIt(r, "Controller - front - GetHomePageData")
	defer util.CommonDeferred(w, r, "Controller", "front", "GetHomePageData")

	retBody, err := model.GetHomePageData(r)
	if err != nil {
		util.RespondWithError(r, w, "500")
		return
	}
	util.RespondData(r, w, retBody, 200)
}

// GetRatingQuestions - Returns rating questions
func GetRatingQuestions(w http.ResponseWriter, r *http.Request) {
	util.LogIt(r, "Controller - front - GetRatingQuestions")
	defer util.CommonDeferred(w, r, "Controller", "front", "GetRatingQuestions")
	retBody, err := model.GetRatingQuestions(r)
	if err != nil {
		util.RespondWithError(r, w, "500")
		return
	}
	util.RespondData(r, w, retBody, 200)
}

// GetFrontStaticData - Returns all front static data
func GetFrontStaticData(w http.ResponseWriter, r *http.Request) {
	util.LogIt(r, "Controller - front - GetFrontStaticData")
	defer util.CommonDeferred(w, r, "Controller", "front", "GetFrontStaticData")
	var resBody = make(map[string]interface{})
	resBody["bucket_url"] = config.Env.AwsBucketURL
	resBody["website_url"] = config.Env.FrontWebURL
	resBody["room_folder"] = config.Env.RoomFolder
	resBody["hotel_folder"] = config.Env.HotelFolder
	resBody["profile_folder"] = config.Env.ProfileFolder
	resBody["property_type_folder"] = config.Env.PropertyTypeFolder
	resBody["popular_city_folder"] = config.Env.PopularCityFolder
	util.RespondData(r, w, resBody, 200)
}

// GetCmsData - Returns Cms Data
func GetCmsData(w http.ResponseWriter, r *http.Request) {
	util.LogIt(r, "Controller - front - GetCmsData")
	defer util.CommonDeferred(w, r, "Controller", "front", "GetCmsData")

	vars := mux.Vars(r)
	slug := vars["slug"]

	ValidateString := ValidateNotNullStructString(slug)
	if ValidateString == 0 {
		util.RespondBadRequest(r, w)
		return
	}

	retBody, err := model.GetCmsData(r, slug)
	if err != nil {
		util.RespondWithError(r, w, "500")
		return
	}
	util.RespondData(r, w, retBody, 200)
}

// GetCmsListData - Returns Cms List Data
func GetCmsListData(w http.ResponseWriter, r *http.Request) {
	util.LogIt(r, "Controller - front - GetCmsListData")
	defer util.CommonDeferred(w, r, "Controller", "front", "GetCmsListData")

	retBody, err := model.GetCmsListData(r)
	if err != nil {
		util.RespondWithError(r, w, "500")
		return
	}
	util.RespondData(r, w, retBody, 200)
}
