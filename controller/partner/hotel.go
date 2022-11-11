package partner

import (
	"encoding/json"
	"net/http"
	"net/http/httputil"
	"strconv"
	"strings"
	"time"
	"tp-api-common/data"
	"tp-api-common/util"
	"tp-system/config"
	"tp-system/controller"
	"tp-system/model"
	"tp-system/model/partner"

	"github.com/gorilla/context"
	"github.com/gorilla/mux"
)

// AddHotelByHotelier - Update Hotel Basic Info
func AddHotelByHotelier(w http.ResponseWriter, r *http.Request) {
	util.LogIt(r, "Controller - V_Partner_Hotel - AddHotelByHotelier")
	defer util.CommonDeferred(w, r, "Controller", "V_Partner_Hotel", "AddHotelByHotelier")
	var reqMap data.Hotel
	err := json.NewDecoder(r.Body).Decode(&reqMap)
	if util.CheckErrorLog(r, err) {
		util.RespondBadRequest(r, w)
		return
	}

	ValidateString := controller.ValidateNotNullStructString(reqMap.HotelPhone, reqMap.Name, reqMap.Description, reqMap.Tag, reqMap.PropertyType)
	ValidateFloat := controller.ValidateNotNullStructFloat(reqMap.HotelStar)
	if ValidateString == 0 || ValidateFloat == 0 {
		util.RespondBadRequest(r, w)
		return
	}

	cnt, err := model.CheckDuplicateRecords(r, "HOTEL", map[string]string{"hotel_name": reqMap.Name, "group_id": context.Get(r, "GroupId").(string)}, nil, "0")
	if err != nil || cnt != 0 {
		if err != nil {
			util.RespondBadRequest(r, w)
		} else {
			util.Respond(r, w, nil, 409, "10010")
		}
		return
	}

	flag, PID := partner.AddHotelByHotelier(r, reqMap)
	if !flag {
		util.RespondWithError(r, w, "500")
		return
	}

	resMap := make(map[string]string)
	resMap["hotelid"] = PID

	util.Respond(r, w, resMap, 201, "")
}

//HotelListOfHotelier - Get hotelier all hotel list i.e status wise
func HotelListOfHotelier(w http.ResponseWriter, r *http.Request) {
	util.LogIt(r, "Controller - V_Partner_Hotel - UpdateLocation")
	defer util.CommonDeferred(w, r, "Controller", "V_Partner_Hotel", "UpdateLocation")
	HotelList := model.PartnerHotelList(r)
	if HotelList == nil {
		HotelList = []map[string]interface{}{}
	}
	util.Respond(r, w, HotelList, 200, "")
}

// UpdateLocation - Update Location Of Hotel
func UpdateLocation(w http.ResponseWriter, r *http.Request) {
	util.LogIt(r, "Controller - V_Partner_Hotel - UpdateLocation")
	defer util.CommonDeferred(w, r, "Controller", "V_Partner_Hotel", "UpdateLocation")
	var reqMap data.Hotel
	err := json.NewDecoder(r.Body).Decode(&reqMap)
	if util.CheckErrorLog(r, err) {
		util.RespondBadRequest(r, w)
		return
	}

	//reqMap.ID = context.Get(r, "HotelId").(string)
	vars := mux.Vars(r)
	reqMap.ID = vars["id"]

	//check of property associated with partner or not
	if !model.IsPartnerContainProperty(r, reqMap.ID, "") {
		util.Respond(r, w, nil, 406, "")
		return
	}

	ValidateString := controller.ValidateNotNullStructString(reqMap.ID, reqMap.Locality, reqMap.ShortAddress, reqMap.LongAddress)
	ValidateFloat := controller.ValidateNotNullStructFloat(reqMap.Latitude, reqMap.Longitude, reqMap.City, reqMap.State, reqMap.Country)
	if ValidateString == 0 || ValidateFloat == 0 {
		util.RespondBadRequest(r, w)
		return
	}

	flag := partner.UpdateLocation(r, reqMap)
	if !flag {
		util.RespondWithError(r, w, "500")
		return
	}

	util.Respond(r, w, nil, 204, "")
}

// UpdateAmenity - Update Amenity Of Hotel
func UpdateAmenity(w http.ResponseWriter, r *http.Request) {
	util.LogIt(r, "Controller - V_Partner_Hotel - UpdateAmenity")
	defer util.CommonDeferred(w, r, "Controller", "V_Partner_Hotel", "UpdateAmenity")
	var reqMap data.HotelAmenity
	err := json.NewDecoder(r.Body).Decode(&reqMap)
	if util.CheckErrorLog(r, err) {
		util.RespondBadRequest(r, w)
		return
	}

	//reqMap.HotelID = context.Get(r, "HotelId").(string)
	vars := mux.Vars(r)
	reqMap.HotelID = vars["id"]

	//check of property associated with partner or not
	if !model.IsPartnerContainProperty(r, reqMap.HotelID, "") {
		util.Respond(r, w, nil, 406, "")
		return
	}

	ValidateString := controller.ValidateNotNullStructString(reqMap.HotelID)
	var Arr = reqMap.Amenity
	if len(Arr) == 0 || ValidateString == 0 {
		util.RespondBadRequest(r, w)
		return
	}

	flag := partner.UpdateAmenity(r, reqMap)
	if !flag {
		util.RespondWithError(r, w, "500")
		return
	}

	util.Respond(r, w, nil, 204, "")
}

// UpdateHotelBasicInfo - Update Hotel Basic Info
func UpdateHotelBasicInfo(w http.ResponseWriter, r *http.Request) {
	util.LogIt(r, "Controller - V_Partner_Hotel - UpdateHotelBasicInfo")
	defer util.CommonDeferred(w, r, "Controller", "V_Partner_Hotel", "UpdateHotelBasicInfo")
	var reqMap data.Hotel
	err := json.NewDecoder(r.Body).Decode(&reqMap)
	if util.CheckErrorLog(r, err) {
		util.RespondBadRequest(r, w)
		return
	}

	//reqMap.ID = context.Get(r, "HotelId").(string)
	vars := mux.Vars(r)
	reqMap.ID = vars["id"]

	//check of property associated with partner or not
	if !model.IsPartnerContainProperty(r, reqMap.ID, "") {
		util.Respond(r, w, nil, 406, "")
		return
	}
	ValidateString := controller.ValidateNotNullStructString(reqMap.ID, reqMap.HotelPhone, reqMap.Name, reqMap.Description, reqMap.Tag, reqMap.PropertyType)
	ValidateFloat := controller.ValidateNotNullStructFloat(reqMap.HotelStar) // 2020-05-17 - HK - Country, State, City, Short Long Address Info Will Be Added With Location
	if ValidateString == 0 || ValidateFloat == 0 {
		util.RespondBadRequest(r, w)
		return
	}

	flag := partner.UpdateHotelBasicInfo(r, reqMap)
	if !flag {
		util.RespondWithError(r, w, "500")
		return
	}

	util.Respond(r, w, nil, 204, "")
}

// UpdatePolicyRules - Update Policy Rules
func UpdatePolicyRules(w http.ResponseWriter, r *http.Request) {
	util.LogIt(r, "Controller - V_Partner_Hotel - UpdatePolicyRules")
	defer util.CommonDeferred(w, r, "Controller", "V_Partner_Hotel", "UpdatePolicyRules")
	var reqMap data.Hotel
	err := json.NewDecoder(r.Body).Decode(&reqMap)
	if util.CheckErrorLog(r, err) {
		util.RespondBadRequest(r, w)
		return
	}

	//reqMap.ID = context.Get(r, "HotelId").(string)
	vars := mux.Vars(r)
	reqMap.ID = vars["id"]

	//check of property associated with partner or not
	if !model.IsPartnerContainProperty(r, reqMap.ID, "") {
		util.Respond(r, w, nil, 406, "")
		return
	}
	ValidateString := controller.ValidateNotNullStructString(reqMap.ID, reqMap.ChekinRules, reqMap.Policy, reqMap.CheckInTime, reqMap.CheckOutTime)
	if ValidateString == 0 {
		util.RespondBadRequest(r, w)
		return
	}

	flag := partner.UpdatePolicyRules(r, reqMap)
	if !flag {
		util.RespondWithError(r, w, "500")
		return
	}

	util.Respond(r, w, nil, 204, "")
}

// UploadHotelImage - Upload Hotel Image Max 5 At One Time
func UploadHotelImage(w http.ResponseWriter, r *http.Request) {
	util.LogIt(r, "Controller - V_Partner_Hotel - UploadHotelImage")
	defer util.CommonDeferred(w, r, "Controller", "V_Partner_Hotel", "UploadHotelImage")
	var reqMap = make(map[string]interface{})
	err := r.ParseMultipartForm(100000)
	if util.CheckErrorLog(r, err) {
		util.RespondBadRequest(r, w)
		return
	}
	var hotelID = r.FormValue("hotelid")

	m := r.MultipartForm
	var ImageCategory = m.Value["image_category"]

	reqMap["ImageCategory"] = ImageCategory[0]
	// if context.Get(r, "HotelId").(string) != "" {
	// 	hotelID = context.Get(r, "HotelId").(string)
	// }
	if len(m.File["hotel_image"]) == 0 || len(ImageCategory) == 0 || hotelID == "" {
		util.RespondBadRequest(r, w)
		return
	}
	reqMap["HotelID"] = hotelID

	// checks if hotel id exists in system with optional status check
	chkIfHotelExists, err := partner.CheckIfHotelExistsOptional(r, hotelID, false, false)
	if err != nil {
		util.SysLogIt("Error While Checking Hotel Existance")
		util.RespondWithError(r, w, "500")
		return
	}
	if chkIfHotelExists == 0 {
		util.SysLogIt("No Such Hotel Exists")
		util.Respond(r, w, nil, 406, "Selected property inactivated by Admin")
		return
	}

	ImageCount, err := partner.GetHotelImageCount(r)
	if err != nil {
		util.RespondBadRequest(r, w)
		return
	}

	var ImageTotalSum = int(ImageCount) + len(m.File["hotel_image"])
	if ImageTotalSum > 30 {
		util.Respond(r, w, nil, 406, "100014")
		return
	}

	if len(m.File["hotel_image"]) > 5 {
		util.Respond(r, w, nil, 406, "100013")
		return
	}

	files := m.File["hotel_image"]
	var ImageArr []string
	for i, _ := range files {
		file, err := files[i].Open()
		if err != nil {
			util.RespondBadRequest(r, w)
			return
		}
		ImageName, err := controller.UploadImageFormData(r, "hotel", files[i])
		if err != nil {
			util.RespondBadRequest(r, w)
			return
		}
		ImageArr = append(ImageArr, ImageName)
		defer file.Close()
	}

	reqMap["Image"] = ImageArr

	flag := partner.UploadHotelImage(r, reqMap)
	if !flag {
		util.RespondWithError(r, w, "500")
		return
	}

	util.Respond(r, w, nil, 201, "")
}

// DeleteHotelImage - Delete Hotel Image
func DeleteHotelImage(w http.ResponseWriter, r *http.Request) {
	util.LogIt(r, "Controller - V_Partner_Hotel - DeleteHotelImage")
	defer util.CommonDeferred(w, r, "Controller", "V_Partner_Hotel", "DeleteHotelImage")
	var reqMap data.ImageName

	vars := mux.Vars(r)
	reqMap.Image = vars["id"]
	reqMap.ID = vars["hotelid"]

	ValidateString := controller.ValidateNotNullStructString(reqMap.Image, reqMap.ID)
	if ValidateString == 0 {
		util.RespondBadRequest(r, w)
		return
	}

	ImageName, err := partner.GetHotelImageName(r, reqMap)
	if err != nil {
		util.RespondWithError(r, w, "500")
		return
	}

	if !controller.DeleteImageFromS3(r, "hotel", ImageName) {
		util.RespondWithError(r, w, "500")
		return
	}

	flag := partner.DeleteHotelImage(r, reqMap)
	if !flag {
		util.RespondWithError(r, w, "500")
		return
	}

	util.Respond(r, w, nil, 204, "")
}

// GetHotelImageList - Get Hotel Image List Info
func GetHotelImageList(w http.ResponseWriter, r *http.Request) {
	util.LogIt(r, "Controller - V_Partner_Hotel - GetHotelImageList")
	defer util.CommonDeferred(w, r, "Controller", "V_Partner_Hotel", "GetHotelImageList")
	//HotelID := context.Get(r, "HotelId")
	HotelID := r.URL.Query().Get("hotelid")
	if HotelID == "" {
		util.RespondBadRequest(r, w)
		return
	}
	retMap, err := partner.GetHotelImageList(r, HotelID)
	if util.CheckErrorLog(r, err) {
		util.RespondWithError(r, w, "500")
		return
	}

	util.RespondData(r, w, retMap, 200)
}

// GetHotelRoomRateData - Get All Data like Inv, Rate, Min, SS, CTA, CTD Of Room Type And Rate Plan Data Of Hotel
func GetHotelRoomRateData(w http.ResponseWriter, r *http.Request) {

	util.LogIt(r, "Controller - V_Partner_Hotel - GetHotelRoomRateData")
	defer util.CommonDeferred(w, r, "Controller", "V_Partner_Hotel", "GetHotelRoomRateData")

	reqMap, err := util.ExtractRequestBody(r)
	if err != nil {
		util.RespondBadRequest(r, w)
		return
	}

	validateFlag := controller.ValidateNotNullAndString(reqMap, []string{"hotel_id"})
	if validateFlag == 0 {
		util.RespondBadRequest(r, w)
		return
	}

	context.Set(r, "HotelId", reqMap["hotel_id"].(string))

	// checks if hotel id exists in system with optional status check
	chkIfHotelExists, err := partner.CheckIfHotelExistsOptional(r, reqMap["hotel_id"].(string), true, false)
	if util.CheckErrorLog(r, err) {
		util.SysLogIt("Error While Checking Hotel Existance")
		util.RespondWithError(r, w, "500")
		return
	}
	if chkIfHotelExists == 0 {
		util.SysLogIt("No Such Hotel Exists")
		util.Respond(r, w, nil, 406, "Selected property inactivated by Admin")
		return
	}

	Data, err := partner.GetHotelRoomRateData(r, reqMap)
	if util.CheckErrorLog(r, err) {
		util.RespondWithError(r, w, "500")
		return
	}

	util.Respond(r, w, Data, 200, "")

}

// FillInvRateData - Dump Inv Rate Data For Hotel
func FillInvRateData(w http.ResponseWriter, r *http.Request) {

	util.LogIt(r, "Controller - V_Partner_Hotel - FillInvRateData")
	defer util.CommonDeferred(w, r, "Controller", "V_Partner_Hotel", "FillInvRateData")

	reqBody, err := util.ExtractRequestBody(r)
	if err != nil {
		util.RespondBadRequest(r, w)
		return
	}
	validateFlag := controller.ValidateNotNullAndString(reqBody, []string{"hotel_id"})
	if validateFlag == 0 {
		util.RespondBadRequest(r, w)
		return
	}

	HotelID := reqBody["hotel_id"].(string)

	// checks if hotel id exists in system with optional status check
	chkIfHotelExists, err := partner.CheckIfHotelExistsOptional(r, HotelID, true, false)
	if err != nil {
		util.SysLogIt("Error While Checking Hotel Existance")
		util.RespondWithError(r, w, "500")
		return
	}
	if chkIfHotelExists == 0 {
		util.SysLogIt("No Such Hotel Exists")
		util.Respond(r, w, nil, 406, "Selected property inactivated by Admin")
		return
	}

	// HotelID := context.Get(r, "HotelId")
	roomDATA, err := partner.CheckRoomCount(r, HotelID)
	if len(roomDATA) == 0 || err != nil {
		util.LogIt(r, "Controller - V_Partner_Hotel - FillInvRateData No Rooms Found For This Hotel")
		util.Respond(r, w, nil, 403, "")
		return
	}

	// To get Current Month's First Date And Next Year's Last Month's Last Date
	// Considering Date Today is 2020-05-19
	// So result will be first : 2020-05-01 and last : 2021-05-31
	currentYear, currentMonth, _ := time.Now().Date()
	first, _ := model.MonthInterval(r, currentYear, currentMonth)
	startDate := first.Format("2006-01-02")

	lastDate := time.Now().AddDate(1, 0, 0).Format("2006-01-02")
	newLastDate := strings.Split(lastDate, "-")
	endYear, _ := strconv.Atoi(newLastDate[0])
	endMonth, _ := strconv.Atoi(newLastDate[1])
	_, last := model.MonthInterval(r, endYear, time.Month(endMonth))
	endDate := last.Format("2006-01-02")

	// util.Respond(r, w, nil, 201, "")
	// return

	var finalArr = make(map[string]interface{})
	var RoomPlanArr []map[string]interface{}

	for k, v := range roomDATA {

		RoomID := v["id"].(string)
		RoomName := v["room_type_name"].(string)
		ForMaxOcc := v["max_occupancy"].(int64)
		RoomInv := v["inventory"].(int64)

		rateDATA, err := partner.CheckRatePlanCount(r, HotelID, RoomID)
		if len(rateDATA) == 0 || err != nil {
			util.LogIt(r, "Controller - V_Partner_Hotel - FillInvRateData No Rate Plans Found For This Hotel")
			util.Respond(r, w, nil, 403, "")
			return
		}

		RoomPlanArr = append(RoomPlanArr, map[string]interface{}{
			"room_id":    RoomID,
			"room_name":  RoomName,
			"occupancy":  ForMaxOcc,
			"inventory":  RoomInv,
			"start_date": startDate,
			"end_date":   endDate,
		})

		var RatePlanArr []map[string]interface{}
		for _, v1 := range rateDATA {

			var OccupancyArr []map[string]interface{}
			var j int64
			for j = 1; j <= ForMaxOcc; j++ {
				joinStr := strconv.FormatInt(j, 16) //
				mainStr := "occ_" + joinStr
				OccupancyArr = append(OccupancyArr, map[string]interface{}{
					mainStr: v1["rate"],
				})
			}

			RatePlanArr = append(RatePlanArr, map[string]interface{}{
				"rate_id":    v1["id"],
				"rate_name":  v1["rate_plan_name"],
				"rate":       OccupancyArr,
				"start_date": startDate,
				"end_date":   endDate,
				"min_night":  1,
				"stop_sell":  0,
				"cta":        0,
				"ctd":        0,
			})
		}
		RoomPlanArr[k]["rate_info"] = RatePlanArr
	}
	// log.Println(RoomPlanArr)
	finalArr["room_info"] = RoomPlanArr

	// allData, err := json.Marshal(finalArr)

	flg := partner.FillInvRateData(r, HotelID, finalArr)
	if !flg {
		util.RespondWithError(r, w, "500")
		return
	}

	util.Respond(r, w, finalArr, 201, "")
}

// UpdateHotelRoomRateData - Updates Inv Rate Data For Hotel
func UpdateHotelRoomRateData(w http.ResponseWriter, r *http.Request) {

	util.LogIt(r, "Controller - V_Partner_Hotel - UpdateHotelRoomRateData")
	defer util.CommonDeferred(w, r, "Controller", "V_Partner_Hotel", "UpdateHotelRoomRateData")

	reqBody, err := util.ExtractRequestBody(r)
	if err != nil {
		util.RespondBadRequest(r, w)
		return
	}

	validateFlag := controller.ValidateNotNullAndFloat(reqBody, []string{"year", "month"})
	validateFlag1 := controller.ValidateExists(reqBody, []string{"data"})
	validateFlagHtl := controller.ValidateNotNullAndString(reqBody, []string{"hotel_id"})

	if validateFlag == 0 || validateFlag1 == 0 || validateFlagHtl == 0 {
		util.RespondBadRequest(r, w)
		return
	}

	// checks if hotel id exists in system with optional status check
	chkIfHotelExists, err := partner.CheckIfHotelExistsOptional(r, reqBody["hotel_id"].(string), true, false)
	if err != nil {
		util.SysLogIt("Error While Checking Hotel Existance")
		util.RespondWithError(r, w, "500")
		return
	}
	if chkIfHotelExists == 0 {
		util.SysLogIt("No Such Hotel Exists")
		util.Respond(r, w, nil, 406, "Selected property inactivated by Admin")
		return
	}

	updateData := reqBody["data"].([]interface{})
	if len(updateData) == 0 {
		util.RespondBadRequest(r, w)
		return
	}

	// HotelID := context.Get(r, "HotelId")

	flag := partner.UpdateHotelRoomRateData(r, reqBody)
	if !flag {
		util.RespondWithError(r, w, "500")
		return
	}

	util.Respond(r, w, nil, 204, "")

}

// UpdateHotelRoomRateDataForBooking - Updates Inv Data For Booking
func UpdateHotelRoomRateDataForBooking(w http.ResponseWriter, r *http.Request) {

	util.LogIt(r, "Controller - V_Partner_Hotel - UpdateHotelRoomRateDataForBooking")
	defer util.CommonDeferred(w, r, "Controller", "V_Partner_Hotel", "UpdateHotelRoomRateDataForBooking")

	token := r.Header.Get("X-Auth-Token")
	if token == "" {
		util.SysLogIt("X-Auth-Token Not Provided")
		util.RespondBadRequest(r, w)
		return
	}

	// checks if authorization token provided is valid or not
	if token != config.Env.InvAuthKey {
		util.SysLogIt("Invalid X-Auth-Token")
		util.Respond(r, w, nil, 401, "")
		return
	}

	// Save a copy of this request for debugging.
	requestDump, err := httputil.DumpRequest(r, true)
	if err != nil {
		util.SysLogIt("Error While Dumping Request Data")
		util.RespondWithError(r, w, "500")
		return
	}
	util.SysLogIt("Request ::")
	util.SysLogIt(string(requestDump))

	reqBody, err := util.ExtractRequestBody(r)
	if err != nil {
		util.RespondBadRequest(r, w)
		return
	}

	validateFlag := controller.ValidateExists(reqBody, []string{"data", "hotel_id"})
	validateFlag1 := controller.ValidateNotNullAndString(reqBody, []string{"hotel_id"})
	if validateFlag == 0 || validateFlag1 == 0 {
		util.SysLogIt("Invalid Request HotelId Or Data Not Provided")
		util.RespondBadRequest(r, w)
		return
	}

	HotelID := reqBody["hotel_id"].(string)
	// checks if hotel id exists in system and is active or not
	chkIfHotelExists, err := partner.CheckIfHotelExists(r, HotelID)
	if err != nil {
		util.SysLogIt("Error While Checking Hotel Existance")
		util.RespondWithError(r, w, "500")
		return
	}
	if chkIfHotelExists == 0 {
		util.SysLogIt("No Such Hotel Exists")
		util.Respond(r, w, nil, 406, "")
		return
	}

	// checks if update data are provided or not in request body
	updateData := reqBody["data"].([]interface{})
	if len(updateData) == 0 {
		util.SysLogIt("No Data Provided")
		util.RespondBadRequest(r, w)
		return
	}

	for i := 0; i < len(updateData); i++ {

		updatePartData := updateData[i].(map[string]interface{})
		validateFlag2 := controller.ValidateExists(updatePartData, []string{"room_id", "start_date", "end_date", "room_cnt"})
		validateFlag3 := controller.ValidateNotNullAndString(updatePartData, []string{"room_id", "start_date", "end_date"})
		validateFlag4 := controller.ValidateNotNullAndFloat(updatePartData, []string{"room_cnt"})
		if validateFlag2 == 0 || validateFlag3 == 0 || validateFlag4 == 0 {
			util.SysLogIt("Invalid Request Room Id OR Date Or Room Count Not Provided")
			util.RespondBadRequest(r, w)
			return
		}

		startDate := updatePartData["start_date"].(string)
		endDate := updatePartData["end_date"].(string)
		roomID := updatePartData["room_id"].(string)

		// checkes if start date is valid or not
		_, err := model.ShortDateFromString(startDate)
		if err != nil {
			util.SysLogIt("Invalid Start Date")
			util.SysLogIt(err)
			util.RespondBadRequest(r, w)
			return
		}

		// checkes if end date is valid or not
		_, err = model.ShortDateFromString(endDate)
		if err != nil {
			util.SysLogIt("Invalid End Date")
			util.SysLogIt(err)
			util.RespondBadRequest(r, w)
			return
		}

		// checks if
		// 1. start date is greater than end date then returns false
		// 2. if start date / end date are of past than current date then returns false
		dateCheck, err := model.CheckDataBoundariesStr(startDate, endDate)
		if err != nil || dateCheck == false {
			util.SysLogIt("date validation error")
			util.SysLogIt(err)
			util.RespondBadRequest(r, w)
			return
		}

		// Check If Room Id Is Valid Or Not and Also Room Exists With The Provided Hotel ID or not
		chkIfHotelRoomExists, err := partner.CheckIfRoomExists(r, roomID, HotelID)
		if err != nil {
			util.SysLogIt("Error While Checking Hotel Room Existance")
			util.RespondWithError(r, w, "500")
			return
		}
		if chkIfHotelRoomExists == 0 {
			util.SysLogIt("No Such Hotel Room Exists")
			util.Respond(r, w, nil, 406, "")
			return
		}

	}

	flag := partner.UpdateHotelRoomRateDataForBooking(r, reqBody)
	if !flag {
		util.SysLogIt("Error While Updating Inventory")
		util.RespondWithError(r, w, "500")
		return
	}

	util.Respond(r, w, nil, 201, "")

}

// GetUpdateLogsOfProperty - Gets All Update Logs Data like Inv, Rate, Min, SS, CTA, CTD Of Room Type And Rate Plan Data Of Hotel
func GetUpdateLogsOfProperty(w http.ResponseWriter, r *http.Request) {
	util.LogIt(r, "Controller - V_Partner_Hotel - GetUpdateLogsOfProperty")
	defer util.CommonDeferred(w, r, "Controller", "V_Partner_Hotel", "GetUpdateLogsOfProperty")
	var hotelID string
	reqBody, err := util.ExtractRequestBody(r)
	if err != nil {
		util.RespondBadRequest(r, w)
		return
	}

	/*
		if reqBody["hotel_id"] != "" {
			context.Set(r, "HotelId", reqBody["hotel_id"])
			hotelID = context.Get(r, "HotelId").(string)
		}
	*/

	validateFlag := controller.ValidateExists(reqBody, []string{"room_id", "rate_id", "date", "logs_for", "hotel_id"})
	validateFlag1 := controller.ValidateNotNullAndString(reqBody, []string{"logs_for"})
	if validateFlag == 0 || validateFlag1 == 0 {
		util.SysLogIt("Invalid Request Data Provided")
		util.RespondBadRequest(r, w)
		return
	}

	logsFor := reqBody["logs_for"].(string)
	if logsFor != "inv" && logsFor != "rates" {
		util.SysLogIt("Invalid Log Identifier Provided")
		util.RespondBadRequest(r, w)
		return
	}

	hotelID = reqBody["hotel_id"].(string)

	var roomID string
	// var rateID string
	if logsFor == "inv" {

		roomID = reqBody["room_id"].(string)
		if roomID == "" {
			util.SysLogIt("For Page Landing Purpose")
		} else {
			validateFlag3 := controller.ValidateNotNullAndString(reqBody, []string{"room_id"})
			if validateFlag3 == 0 {
				util.SysLogIt("Invalid Request Data Provided 2")
				util.RespondBadRequest(r, w)
				return
			}
			// roomID = reqBody["room_id"].(string)

			// Check If Room Id Is Valid Or Not and Also Room Exists With The Hotel ID or not
			chkIfHotelRoomExists, err := partner.CheckIfRoomExists(r, roomID, hotelID)
			if err != nil {
				util.SysLogIt("Error While Checking Hotel Room Existance")
				util.RespondWithError(r, w, "500")
				return
			}
			if chkIfHotelRoomExists == 0 {
				util.SysLogIt("No Such Hotel Room Exists")
				util.Respond(r, w, nil, 406, "")
				return
			}
		}
	} else if logsFor == "rates" {
		// validateFlag3 := controller.ValidateNotNullAndString(reqBody, []string{"rate_id", "room_id"})
		roomID = reqBody["room_id"].(string)
		if roomID == "" {
			util.SysLogIt("For Rates Page Landing Purpose")
		} else {
			validateFlag3 := controller.ValidateNotNullAndString(reqBody, []string{"room_id"})
			if validateFlag3 == 0 {
				util.SysLogIt("Invalid Request Data Provided 3")
				util.RespondBadRequest(r, w)
				return
			}
			// roomID = reqBody["room_id"].(string)
			// rateID = reqBody["rate_id"].(string)

			// Check If Room Id, Rate Id Is Valid Or Not and Also RatePlan Exists With The Hotel ID or not
			/*chkIfHotelRoomRateplanExists, err := partner.CheckIfRatePlanExists(r, roomID, rateID, hotelID)
			if err != nil {
				util.SysLogIt("Error While Checking Hotel Room RatePlan Existance")
				util.RespondWithError(r, w, "500")
				return
			}
			if chkIfHotelRoomRateplanExists == 0 {
				util.SysLogIt("No Such Hotel Room RatePlan Exists")
				util.Respond(r, w, nil, 406, "")
				return
			}*/
		}
	}

	Data, err := partner.GetUpdateLogsOfProperty(r, reqBody)
	if err != nil {
		util.RespondWithError(r, w, "500")
		return
	}

	if len(Data["data"].([]map[string]interface{})) == 0 {
		Data["data"] = []string{}
	}

	/*var code int
	retData := Data["data"].([]map[string]interface{})
	if len(retData) > 0 {
		code = 200
	} else {
		code = 204
	}*/
	util.Respond(r, w, Data, 200, "")
}

// GetReviewOfHotel - Get Review Flag for Hotel
func GetReviewOfHotel(w http.ResponseWriter, r *http.Request) {

	util.LogIt(r, "Controller - V_Partner_Hotel - GetReviewOfHotel")
	defer util.CommonDeferred(w, r, "Controller", "V_Partner_Hotel", "GetReviewOfHotel")
	var hotelID string

	vars := mux.Vars(r)
	hotelID = vars["id"]

	if hotelID == "" {
		util.RespondBadRequest(r, w)
		return
	}

	Data, err := partner.GetReviewOfHotel(r, hotelID)
	if util.CheckErrorLog(r, err) {
		util.RespondWithError(r, w, "500")
		return
	}

	util.Respond(r, w, Data, 200, "")

}

// VerifyHotel - Verify Hotel data submit to Admin
func VerifyHotel(w http.ResponseWriter, r *http.Request) {

	util.LogIt(r, "Controller - V_Partner_Hotel - VerifyHotel")
	defer util.CommonDeferred(w, r, "Controller", "V_Partner_Hotel", "VerifyHotel")
	var hotelID string

	vars := mux.Vars(r)
	hotelID = vars["id"]

	if hotelID == "" {
		util.RespondBadRequest(r, w)
		return
	}

	// Check whether the hotel is already reviewed or not
	ReviewData, err := partner.GetReviewOfHotel(r, hotelID)
	if util.CheckErrorLog(r, err) {
		util.RespondWithError(r, w, "500")
		return
	}

	if len(ReviewData) > 0 {
		if ReviewData["is_review"].(int64) == 1 {
			util.RespondData(r, w, nil, 406)
			return
		}
	} else {
		util.RespondWithError(r, w, "500")
		return
	}

	Flg, statusCode, err, errCodeString := partner.VerifyHotel(r, hotelID)
	if util.CheckErrorLog(r, err) {
		util.RespondWithError(r, w, "500")
		return
	}

	if !Flg {
		if statusCode == "verification" {
			util.Respond(r, w, nil, 409, errCodeString)
			return
		}
		util.RespondWithError(r, w, "500")
		return
	}

	util.Respond(r, w, nil, 200, "")

}
