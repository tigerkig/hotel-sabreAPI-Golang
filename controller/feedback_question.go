package controller

import (
	"encoding/json"
	"net/http"
	"tp-api-common/data"
	"tp-api-common/util"
	"tp-system/model"

	"github.com/gorilla/mux"
)

// AddQuestion - Add FeedBack Question
func AddQuestion(w http.ResponseWriter, r *http.Request) {

	util.LogIt(r, "Controller - V_Feedback_Question - AddQuestion")
	defer util.CommonDeferred(w, r, "Controller", "V_Feedback_Question", "AddQuestion")

	var reqMap data.FeedBackQuestion

	err := json.NewDecoder(r.Body).Decode(&reqMap)
	if err != nil {
		util.RespondBadRequest(r, w)
		return
	}

	ValidateString := ValidateNotNullStructString(reqMap.Question)
	if ValidateString == 0 {
		util.RespondBadRequest(r, w)
		return
	}

	cnt, err := model.CheckDuplicateRecords(r, "FEEDBACK_QUESTION", map[string]string{"question": reqMap.Question}, nil, "0")
	if err != nil || cnt != 0 {
		if err != nil {
			util.RespondBadRequest(r, w)
		} else {
			util.Respond(r, w, nil, 409, "10010")
		}
		return
	}

	flag := model.AddQuestion(r, reqMap)
	if !flag {
		util.RespondWithError(r, w, "500")
		return
	}

	util.Respond(r, w, nil, 201, "")
}

// UpdateQuestion - Update FeedBack Question
func UpdateQuestion(w http.ResponseWriter, r *http.Request) {

	util.LogIt(r, "Controller - V_Feedback_Question - UpdateQuestion")
	defer util.CommonDeferred(w, r, "Controller", "V_Feedback_Question", "UpdateQuestion")

	var reqMap data.FeedBackQuestion

	err := json.NewDecoder(r.Body).Decode(&reqMap)
	if err != nil {
		util.RespondBadRequest(r, w)
		return
	}

	vars := mux.Vars(r)
	reqMap.ID = vars["id"]

	ValidateString := ValidateNotNullStructString(reqMap.ID, reqMap.Question)
	if ValidateString == 0 {
		util.RespondBadRequest(r, w)
		return
	}

	cnt, err := model.CheckDuplicateRecords(r, "FEEDBACK_QUESTION", map[string]string{"question": reqMap.Question}, nil, reqMap.ID)
	if err != nil || cnt != 0 {
		if err != nil {
			util.RespondBadRequest(r, w)
		} else {
			util.Respond(r, w, nil, 409, "10010")
		}
		return
	}

	flag := model.UpdateQuestion(r, reqMap)
	if !flag {
		util.RespondWithError(r, w, "500")
		return
	}

	util.Respond(r, w, nil, 204, "")
}

// QuestionListing - Datatable FeedBack Question listing with filter and order
func QuestionListing(w http.ResponseWriter, r *http.Request) {

	util.LogIt(r, "Controller - V_Feedback_Question - QuestionListing")
	defer util.CommonDeferred(w, r, "Controller", "V_Feedback_Question", "QuestionListing")

	var reqMap data.JQueryTableUI

	err := json.NewDecoder(r.Body).Decode(&reqMap)
	if err != nil {
		util.RespondBadRequest(r, w)
		return
	}

	Data, err := model.QuestionListing(r, reqMap)
	if err != nil {
		util.RespondWithError(r, w, "500")
		return
	}

	util.Respond(r, w, Data, 200, "")
}
