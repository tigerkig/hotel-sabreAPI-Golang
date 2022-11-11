package partner

import (
	"encoding/json"
	"net/http"
	"tp-api-common/data"
	"tp-api-common/util"
	"tp-system/controller"
	"tp-system/model"
	"tp-system/model/partner"

	"github.com/gorilla/mux"
)

// AddRoomType - Add Room Type
func AddRoomType(w http.ResponseWriter, r *http.Request) {

	util.LogIt(r, "Controller - V_Room_Type - AddRoomType")
	defer util.CommonDeferred(w, r, "Controller", "V_Room_Type", "AddRoomType")

	var reqMap data.RoomType

	err := json.NewDecoder(r.Body).Decode(&reqMap)
	if err != nil {
		util.RespondBadRequest(r, w)
		return
	}

	// ValidateString := controller.ValidateNotNullStructString(reqMap.Name, reqMap.BedType, reqMap.Description, reqMap.RoomView)
	ValidateString := controller.ValidateNotNullStructString(reqMap.Name, reqMap.BedType, reqMap.Description, reqMap.RoomView, reqMap.HotelID)
	ValidateFloat := controller.ValidateNotNullStructFloat(reqMap.MaxOccupancy, reqMap.Inventory, reqMap.IsExtraBed, reqMap.SortOrder)
	if ValidateString == 0 || ValidateFloat == 0 {
		util.RespondBadRequest(r, w)
		return
	}

	if reqMap.IsExtraBed == 1 {
		ValidateString := controller.ValidateNotNullStructString(reqMap.ExtraBed)
		if ValidateString == 0 {
			util.RespondBadRequest(r, w)
			return
		}
	}

	HotelID := reqMap.HotelID
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
	// cnt, err := model.CheckDuplicateRecords(r, "ROOM_TYPE", map[string]string{"room_type_name": reqMap.Name, "hotel_id": HotelID.(string)}, nil, "0")
	cnt, err := model.CheckDuplicateRecords(r, "ROOM_TYPE", map[string]string{"room_type_name": reqMap.Name, "hotel_id": HotelID}, nil, "0")
	if err != nil || cnt != 0 {
		if err != nil {
			util.RespondBadRequest(r, w)
		} else {
			util.Respond(r, w, nil, 409, "10010")
		}
		return
	}

	retMap, err := partner.AddRoomType(r, reqMap)
	if err != nil {
		util.RespondWithError(r, w, "500")
		return
	}

	util.Respond(r, w, retMap, 201, "")
}

// UpdateRoomType - Update Room Type
func UpdateRoomType(w http.ResponseWriter, r *http.Request) {

	util.LogIt(r, "Controller - V_Room_Type - UpdateRoomType")
	defer util.CommonDeferred(w, r, "Controller", "V_Room_Type", "UpdateRoomType")

	var reqMap data.RoomType

	err := json.NewDecoder(r.Body).Decode(&reqMap)
	if err != nil {
		util.RespondBadRequest(r, w)
		return
	}

	vars := mux.Vars(r)
	reqMap.ID = vars["id"]

	// ValidateString := controller.ValidateNotNullStructString(reqMap.ID, reqMap.Name, reqMap.BedType, reqMap.Description, reqMap.RoomView)
	ValidateString := controller.ValidateNotNullStructString(reqMap.ID, reqMap.Name, reqMap.BedType, reqMap.Description, reqMap.RoomView, reqMap.HotelID)
	ValidateFloat := controller.ValidateNotNullStructFloat(reqMap.MaxOccupancy, reqMap.Inventory, reqMap.IsExtraBed, reqMap.SortOrder)
	if ValidateString == 0 || ValidateFloat == 0 {
		util.RespondBadRequest(r, w)
		return
	}

	if reqMap.IsExtraBed == 1 {
		ValidateString := controller.ValidateNotNullStructString(reqMap.ExtraBed)
		if ValidateString == 0 {
			util.RespondBadRequest(r, w)
			return
		}
	}

	HotelID := reqMap.HotelID
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
	// cnt, err := model.CheckDuplicateRecords(r, "ROOM_TYPE", map[string]string{"room_type_name": reqMap.Name, "hotel_id": HotelID.(string)}, nil, reqMap.ID)
	cnt, err := model.CheckDuplicateRecords(r, "ROOM_TYPE", map[string]string{"room_type_name": reqMap.Name, "hotel_id": HotelID}, nil, reqMap.ID)
	if err != nil || cnt != 0 {
		if err != nil {
			util.RespondBadRequest(r, w)
		} else {
			util.Respond(r, w, nil, 409, "10010")
		}
		return
	}

	flag := partner.UpdateRoomType(r, reqMap)
	if !flag {
		util.RespondWithError(r, w, "500")
		return
	}

	util.Respond(r, w, nil, 204, "")
}

// RoomTypeListing - Return Datatable Listing Of Room Type
func RoomTypeListing(w http.ResponseWriter, r *http.Request) {

	util.LogIt(r, "Controller - V_Room_Type - RoomTypeListing")
	defer util.CommonDeferred(w, r, "Controller", "V_Room_Type", "RoomTypeListing")

	var reqMap data.JQueryTableUI

	err := json.NewDecoder(r.Body).Decode(&reqMap)
	if err != nil {
		util.RespondBadRequest(r, w)
		return
	}

	if reqMap.HotelID == "" {
		util.LogIt(r, "Controller - V_Room_Type - RoomTypeListing - Hotel Id Missing")
		util.RespondBadRequest(r, w)
		return
	}

	Data, err := partner.RoomTypeListing(r, reqMap)
	if err != nil {
		util.RespondWithError(r, w, "500")
		return
	}

	util.Respond(r, w, Data, 200, "")
}

// GetRoomType - Get Room Type Info
func GetRoomType(w http.ResponseWriter, r *http.Request) {

	util.LogIt(r, "Controller - V_Room_Type - GetRoomType")
	defer util.CommonDeferred(w, r, "Controller", "V_Room_Type", "GetRoomType")

	vars := mux.Vars(r)
	ID := vars["id"]

	HotelID := r.URL.Query().Get("hotelid")

	ValidateString := controller.ValidateNotNullStructString(ID, HotelID)
	if ValidateString == 0 {
		util.RespondBadRequest(r, w)
		return
	}

	// retMap, err := partner.GetRoomType(r, ID, HotelID)
	retMap, err := partner.GetRoomTypeWithHotel(r, ID, HotelID)
	if err != nil {
		util.RespondWithError(r, w, "500")
		return
	}

	util.RespondData(r, w, retMap, 200)
}

// UpdateRoomAmenity - Update Amenity Of Room
func UpdateRoomAmenity(w http.ResponseWriter, r *http.Request) {

	util.LogIt(r, "Controller - V_Room_Type - UpdateRoomAmenity")
	defer util.CommonDeferred(w, r, "Controller", "V_Room_Type", "UpdateRoomAmenity")

	var reqMap data.RoomAmenity

	err := json.NewDecoder(r.Body).Decode(&reqMap)
	if err != nil {
		util.RespondBadRequest(r, w)
		return
	}

	vars := mux.Vars(r)
	reqMap.RoomID = vars["id"]
	// reqMap.HotelID = context.Get(r, "HotelId").(string)
	ValidateString := controller.ValidateNotNullStructString(reqMap.HotelID, reqMap.RoomID)
	var Arr = reqMap.Amenity
	if len(Arr) == 0 || ValidateString == 0 {
		util.RespondBadRequest(r, w)
		return
	}

	flag := partner.UpdateRoomAmenity(r, reqMap)
	if !flag {
		util.RespondWithError(r, w, "500")
		return
	}

	util.Respond(r, w, nil, 204, "")
}

// UploadHotelRoomImage - Upload Hotel Room Image Max 5 At One Time
func UploadHotelRoomImage(w http.ResponseWriter, r *http.Request) {

	util.LogIt(r, "Controller - V_Partner_RoomType - UploadHotelRoomImage")
	defer util.CommonDeferred(w, r, "Controller", "V_Partner_RoomType", "UploadHotelRoomImage")

	var reqMap = make(map[string]interface{})

	err := r.ParseMultipartForm(100000)
	if err != nil {
		util.RespondBadRequest(r, w)
		return
	}

	m := r.MultipartForm
	// reqMap["HotelID"] = context.Get(r, "HotelId").(string)

	HotelID := r.FormValue("hotel_id")
	reqMap["HotelID"] = HotelID

	var RoomID = m.Value["room_type_id"]
	reqMap["RoomID"] = RoomID[0]

	if len(m.File["hotel_room_image"]) == 0 || len(RoomID) == 0 || HotelID == "" {
		util.RespondBadRequest(r, w)
		return
	}

	ImageCount, err := partner.GetHotelRoomImageCount(r, reqMap["RoomID"].(string))
	if err != nil {
		util.RespondBadRequest(r, w)
		return
	}

	var ImageTotalSum = int(ImageCount) + len(m.File["hotel_room_image"])
	if ImageTotalSum > 10 {
		util.Respond(r, w, nil, 406, "100014")
		return
	}

	if len(m.File["hotel_room_image"]) > 5 {
		util.Respond(r, w, nil, 406, "100013")
		return
	}

	files := m.File["hotel_room_image"]
	var ImageArr []string
	for i := range files {
		file, err := files[i].Open()
		if err != nil {
			util.RespondBadRequest(r, w)
			return
		}
		ImageName, err := controller.UploadImageFormData(r, "room", files[i])
		if err != nil {
			util.RespondBadRequest(r, w)
			return
		}
		ImageArr = append(ImageArr, ImageName)
		defer file.Close()
	}

	reqMap["Image"] = ImageArr

	flag := partner.UploadHotelRoomImage(r, reqMap)
	if !flag {
		util.RespondWithError(r, w, "500")
		return
	}

	util.Respond(r, w, nil, 201, "")
}

// DeleteHotelRoomImage - Delete Hotel Room Image
func DeleteHotelRoomImage(w http.ResponseWriter, r *http.Request) {

	util.LogIt(r, "Controller - V_Partner_RoomType - DeleteHotelRoomImage")
	defer util.CommonDeferred(w, r, "Controller", "V_Partner_RoomType", "DeleteHotelRoomImage")

	var reqMap data.RoomImageName

	vars := mux.Vars(r)
	reqMap.Image = vars["id"]

	reqMap.HotelID = r.URL.Query().Get("hotelid")
	reqMap.RoomID = r.URL.Query().Get("roomid")

	ValidateString := controller.ValidateNotNullStructString(reqMap.Image, reqMap.HotelID, reqMap.RoomID)
	if ValidateString == 0 {
		util.RespondBadRequest(r, w)
		return
	}

	ImageName, err := partner.GetHotelRoomImageName(r, reqMap)
	if err != nil {
		util.RespondWithError(r, w, "500")
		return
	}

	if !controller.DeleteImageFromS3(r, "room", ImageName) {
		util.RespondWithError(r, w, "500")
		return
	}

	flag := partner.DeleteHotelRoomImage(r, reqMap)
	if !flag {
		util.RespondWithError(r, w, "500")
		return
	}

	util.Respond(r, w, nil, 204, "")
}

// GetRoomImageList - Get Room Type Info
func GetRoomImageList(w http.ResponseWriter, r *http.Request) {

	util.LogIt(r, "Controller - V_Room_Type - GetRoomImageList")
	defer util.CommonDeferred(w, r, "Controller", "V_Room_Type", "GetRoomImageList")

	vars := mux.Vars(r)
	ID := vars["id"]
	HotelID := r.URL.Query().Get("hotelid")

	ValidateString := controller.ValidateNotNullStructString(ID, HotelID)
	if ValidateString == 0 {
		util.RespondBadRequest(r, w)
		return
	}

	retMap, err := partner.GetRoomImageList(r, ID, HotelID)
	if err != nil {
		util.RespondWithError(r, w, "500")
		return
	}

	util.RespondData(r, w, retMap, 200)
}

// GetRoomTypeList - Get Room Type List For Other Module
func GetRoomTypeList(w http.ResponseWriter, r *http.Request) {

	util.LogIt(r, "Controller - V_Room_Type - GetRoomTypeList")
	defer util.CommonDeferred(w, r, "Controller", "V_Room_Type", "GetRoomTypeList")

	HotelID := r.URL.Query().Get("hotelid")
	if HotelID == "" {
		util.RespondBadRequest(r, w)
		return
	}

	retMap, err := partner.GetRoomTypeList(r, HotelID)
	if err != nil {
		util.RespondWithError(r, w, "500")
		return
	}

	util.RespondData(r, w, retMap, 200)
}

// SortRoomImage - Sorts Room Image - 2021-04-27 - HK
func SortRoomImage(w http.ResponseWriter, r *http.Request) {
	util.LogIt(r, "Controller - V_Room_Type - SortRoomImage")
	defer util.CommonDeferred(w, r, "Controller", "V_Room_Type", "SortRoomImage")

	reqBody, err := util.ExtractRequestBody(r)
	if err != nil {
		util.RespondBadRequest(r, w)
		return
	}

	vars := mux.Vars(r)
	id := vars["id"]

	validateFlag := controller.ValidateNotNullAndString(reqBody, []string{"hotel_id", "room_id"})
	validateFlag1 := controller.ValidateNotNullAndFloat(reqBody, []string{"sortorder"})
	if validateFlag == 0 || validateFlag1 == 0 || id == "" {
		util.RespondBadRequest(r, w)
		return
	}

	reqBody["id"] = id
	flag, err := partner.SortRoomImage(r, reqBody)
	if flag == 0 || err != nil {
		util.RespondWithError(r, w, "500")
		return
	}

	util.Respond(r, w, nil, 204, "")
}
