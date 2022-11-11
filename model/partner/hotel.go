package partner

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"strings"
	"time"
	"tp-api-common/data"
	"tp-api-common/util"
	"tp-system/config"
	"tp-system/model"

	"github.com/gorilla/context"
	gonanoid "github.com/matoous/go-nanoid"
)

func AddHotelByHotelier(r *http.Request, reqMap data.Hotel) (bool, string) {
	util.LogIt(r, "model - V_Partner_Hotel - AddHotelByHotelier")
	var Qry, SQry bytes.Buffer
	nanoHotelID, _ := gonanoid.Nanoid()
	nanoSettingID, _ := gonanoid.Nanoid()
	Qry.WriteString("INSERT INTO cf_hotel_info SET id=?, hotel_name=?, hotel_star=?, description=?, hotel_phone=?, property_type_id=?, group_id=?, status=1, created_at=?, created_by=?")
	err := model.ExecuteNonQuery(Qry.String(), nanoHotelID, reqMap.Name, reqMap.HotelStar, reqMap.Description, reqMap.HotelPhone, reqMap.PropertyType, context.Get(r, "GroupId"), util.GetIsoLocalDateTime(), 0)
	if util.CheckErrorLog(r, err) {
		return false, ""
	}

	SQry.WriteString(" INSERT INTO cf_hotel_settings SET id = ?, hotel_id=? ")
	err = model.ExecuteNonQuery(SQry.String(), nanoSettingID, nanoHotelID)

	var TagStringArr []string
	var TagArr = strings.Split(reqMap.Tag, ",")
	if len(TagArr) > 0 {
		for _, val := range TagArr {
			var TagQry bytes.Buffer
			nanoid, _ := gonanoid.Nanoid()
			TagQry.WriteString("INSERT INTO cf_hotel_tag SET id=?, tag_id=?, hotel_id=?")
			err = model.ExecuteNonQuery(TagQry.String(), nanoid, val, nanoHotelID)
			if util.CheckErrorLog(r, err) {
				return false, ""
			}

			BeforeUpdate, _ := model.GetModuleFieldByID(r, "PROPERTY_TAG", val, "tag")
			TagStringArr = append(TagStringArr, BeforeUpdate.(string))
		}
	}

	reqStruct := util.ToMap(reqMap)
	reqStruct["HotelStar"] = reqMap.HotelStar
	reqStruct["HotelPhone"] = reqMap.HotelPhone
	reqStruct["Tag"] = strings.Join(TagStringArr, ",")
	UpdatePartnerDefaultHotelID(r, nanoHotelID)
	model.AddLog(r, "", "HOTEL", "Create", nanoHotelID, model.GetLogsValueMap(r, reqStruct, true, "State,City,Country,Latitude,Longitude"))
	model.CacheChn <- model.CacheObj{
		Type: "hotel",
		ID:   nanoHotelID,
	}
	return true, nanoHotelID
}

// UpdatePartnerDefaultHotelID - Update Partner Default Hotel
func UpdatePartnerDefaultHotelID(r *http.Request, nanoID string) bool {
	util.LogIt(r, "model - V_Partner_Hotel - UpdatePartnerDefaultHotelID")
	var UQry bytes.Buffer

	UQry.WriteString(" UPDATE cf_hotel_client SET hotel_id = ? WHERE group_id = ?")
	err := model.ExecuteNonQuery(UQry.String(), nanoID, context.Get(r, "GroupId"))
	if util.CheckErrorLog(r, err) {
		return false
	}

	return true
}

// UpdateLocation - Update Location Of Hotel
func UpdateLocation(r *http.Request, reqMap data.Hotel) bool {
	util.LogIt(r, "model - V_Partner_Hotel - UpdateLocation")
	var Qry bytes.Buffer

	Qry.WriteString("UPDATE cf_hotel_info SET short_address = ?,long_address=?, latitude=?, longitude=?, locality_id=?, city_id=?, state_id=?, country_id=? WHERE id = ?")
	err := model.ExecuteNonQuery(Qry.String(), reqMap.ShortAddress, reqMap.LongAddress, reqMap.Latitude, reqMap.Longitude, reqMap.Locality, reqMap.City, reqMap.State, reqMap.Country, reqMap.ID)
	if util.CheckErrorLog(r, err) {
		return false
	}

	CountryName, _ := model.GetModuleFieldByID(r, "COUNTRY", fmt.Sprintf("%.0f", reqMap.Country), "name")
	State, _ := model.GetModuleFieldByID(r, "STATE", fmt.Sprintf("%.0f", reqMap.State), "name")
	City, _ := model.GetModuleFieldByID(r, "CITY", fmt.Sprintf("%.0f", reqMap.City), "name")
	Locality, _ := model.GetModuleFieldByID(r, "LOCALITY", reqMap.Locality, "locality")

	reqStruct := util.ToMap(reqMap)
	reqStruct["Country"] = CountryName
	reqStruct["State"] = State
	reqStruct["City"] = City
	reqStruct["Locality"] = Locality

	model.UpdateHotelOnList(reqMap.ID) // 2020-06-24 - HK - Sync With Mongo Added - Partner Panel

	model.AddLog(r, "", "HOTEL", "Update Location", reqMap.ID, model.GetLogsValueMap(r, reqStruct, true, "HotelStar"))

	return true
}

// UpdateAmenity - Update Amenity Of Hotel
func UpdateAmenity(r *http.Request, reqMap data.HotelAmenity) bool {
	util.LogIt(r, "model - V_Partner_Hotel - UpdateAmenity")
	var DelQry bytes.Buffer
	var AmenityArr = reqMap.Amenity

	DelQry.WriteString("DELETE FROM cf_hotel_amenities WHERE hotel_id = ?")
	err := model.ExecuteNonQuery(DelQry.String(), reqMap.HotelID)
	if util.CheckErrorLog(r, err) {
		return false
	}

	for _, val := range AmenityArr {
		var Qry bytes.Buffer
		nanoid, _ := gonanoid.Nanoid()
		Qry.WriteString("INSERT INTO cf_hotel_amenities SET id=?, amenity_id=?,instruction = ?,hotel_id = ?")
		err = model.ExecuteNonQuery(Qry.String(), nanoid, val.ID, val.Description, reqMap.HotelID)
		if util.CheckErrorLog(r, err) {
			return false
		}
	}

	model.UpdateHotelAmenity(reqMap.HotelID) // 2020-06-24 - HK - Amenity Sync With Mongo Added - Partner Panel

	model.AddLog(r, "", "HOTEL", "Update Hotel Amenities", reqMap.HotelID, map[string]interface{}{})

	return true
}

// UpdateHotelBasicInfo - Update Hotel Basic Info
func UpdateHotelBasicInfo(r *http.Request, reqMap data.Hotel) bool {
	util.LogIt(r, "model - V_Partner_Hotel - UpdateHotelBasicInfo")
	var Qry, DelQry bytes.Buffer

	Qry.WriteString("UPDATE cf_hotel_info SET hotel_name=?, hotel_star=?, description=?, hotel_phone=?, property_type_id=? WHERE id = ?")
	err := model.ExecuteNonQuery(Qry.String(), reqMap.Name, reqMap.HotelStar, reqMap.Description, reqMap.HotelPhone, reqMap.PropertyType, reqMap.ID)

	if util.CheckErrorLog(r, err) {
		return false
	}
	var TagStringArr []string
	var TagArr = strings.Split(reqMap.Tag, ",")
	if len(TagArr) > 0 {
		DelQry.WriteString("DELETE FROM cf_hotel_tag WHERE hotel_id = ?")
		err = model.ExecuteNonQuery(DelQry.String(), reqMap.ID)
		if util.CheckErrorLog(r, err) {
			return false
		}

		for _, val := range TagArr {
			var TagQry bytes.Buffer
			nanoid, _ := gonanoid.Nanoid()
			TagQry.WriteString("INSERT INTO cf_hotel_tag SET id=?, tag_id=?, hotel_id=?")
			err = model.ExecuteNonQuery(TagQry.String(), nanoid, val, reqMap.ID)
			if util.CheckErrorLog(r, err) {
				return false
			}

			BeforeUpdate, _ := model.GetModuleFieldByID(r, "PROPERTY_TAG", val, "tag")
			TagStringArr = append(TagStringArr, BeforeUpdate.(string))
		}
	}

	reqStruct := util.ToMap(reqMap)
	reqStruct["HotelStar"] = reqMap.HotelStar
	reqStruct["HotelPhone"] = reqMap.HotelPhone
	reqStruct["Tag"] = strings.Join(TagStringArr, ",")

	model.UpdateHotelOnList(reqMap.ID) // 2020-06-24 - HK - Sync With Mongo Added - Partner Panel
	model.UpdateHotelTag(reqMap.ID)    // 2020-06-25 - HK - Sync With Mongo Added - Partner Panel

	model.AddLog(r, "", "HOTEL", "Update", reqMap.ID, model.GetLogsValueMap(r, reqStruct, true, "State,City,Country,Latitude,Longitude"))

	return true
}

// UpdatePolicyRules - Update Policy Rules
func UpdatePolicyRules(r *http.Request, reqMap data.Hotel) bool {
	util.LogIt(r, "model - V_Partner_Hotel - UpdatePolicyRules")
	var Qry, SelQry bytes.Buffer

	Qry.WriteString("UPDATE cf_hotel_info SET policy=?,checkin_rules=? WHERE id = ?")
	err := model.ExecuteNonQuery(Qry.String(), reqMap.Policy, reqMap.ChekinRules, reqMap.ID)
	if util.CheckErrorLog(r, err) {
		return false
	}

	SelQry.WriteString("UPDATE cf_hotel_settings SET checkin_time=?,checkout_time = ? WHERE hotel_id = ?")
	err = model.ExecuteNonQuery(SelQry.String(), reqMap.CheckInTime, reqMap.CheckOutTime, reqMap.ID)
	if util.CheckErrorLog(r, err) {
		return false
	}

	reqStruct := make(map[string]interface{})
	reqStruct["Policy"] = reqMap.Policy
	reqStruct["Checkin Rules"] = reqMap.ChekinRules
	reqStruct["Check In Time"] = util.ParseTime(reqMap.CheckInTime)
	reqStruct["Check Out Time"] = util.ParseTime(reqMap.CheckOutTime)

	model.UpdateHotelOnList(reqMap.ID) // 2020-06-24 - HK - Policy Sync With Mongo Added - Partner Panel

	model.AddLog(r, "", "HOTEL", "Update Policy Rules", reqMap.ID, reqStruct)

	return true
}

// UploadHotelImage - Upload Hotel Image Max 5 At One Time
func UploadHotelImage(r *http.Request, reqMap map[string]interface{}) bool {
	util.LogIt(r, "model - V_Partner_Hotel - UploadHotelImage")
	var SQLGet bytes.Buffer
	SQLGet.WriteString(" SELECT CASE WHEN MAX(sortorder) IS NULL THEN 0 ELSE MAX(sortorder) END AS sortorder FROM cf_hotel_image WHERE hotel_id = ? ")
	SortOrder, err := model.ExecuteRowQuery(SQLGet.String(), reqMap["HotelID"])
	if util.CheckErrorLog(r, err) {
		return false
	}

	var latestSortOrder = int(SortOrder["sortorder"].(int64))

	for i, val := range reqMap["Image"].([]string) {
		var Qry bytes.Buffer
		nanoid, _ := gonanoid.Nanoid()
		Qry.WriteString("INSERT INTO cf_hotel_image SET id=?, category_id=?,image = ?,hotel_id = ?,sortorder=?")
		err := model.ExecuteNonQuery(Qry.String(), nanoid, reqMap["ImageCategory"], val, reqMap["HotelID"], latestSortOrder+i+1)
		if util.CheckErrorLog(r, err) {
			return false
		}
	}

	util.LogIt(r, " General Image Start ")
	genImg := model.UpdateHotelImage(reqMap["HotelID"].(string)) // 2020-06-24 - HK - Image Sync With Mongo Added - Partner Panel
	if !genImg {
		return false
	}

	util.LogIt(r, " Category Image Start ")
	catImg := model.UpdateDetailedImage(reqMap["HotelID"].(string)) // 2020-06-24 - HK - Image Sync With Mongo Added - Partner Panel
	if !catImg {
		return false
	}

	BeforeUpdate, _ := model.GetModuleFieldByID(r, "IMAGE_CATEGORY", reqMap["ImageCategory"].(string), "name")

	model.AddLog(r, "", "HOTEL", "Upload Hotel Images", reqMap["HotelID"].(string), map[string]interface{}{"Image Category": BeforeUpdate, "Total Image": len(reqMap["Image"].([]string))})

	return true
}

// DeleteHotelImage - Delete Hotel Image
func DeleteHotelImage(r *http.Request, reqMap data.ImageName) bool {
	util.LogIt(r, "model - V_Partner_Hotel - DeleteHotelImage")
	var Qry bytes.Buffer

	Qry.WriteString("DELETE FROM cf_hotel_image WHERE id=? AND hotel_id = ?")
	err := model.ExecuteNonQuery(Qry.String(), reqMap.Image, reqMap.ID)
	if util.CheckErrorLog(r, err) {
		return false
	}

	model.UpdateHotelImage(reqMap.ID)    // 2020-06-24 - HK - Image Sync With Mongo Added - Partner Panel
	model.UpdateDetailedImage(reqMap.ID) // 2020-06-24 - HK - Image Sync With Mongo Added - Partner Panel
	model.AddLog(r, "", "HOTEL", "Delete Hotel Images", reqMap.ID, map[string]interface{}{})

	return true
}

// GetHotelImageName - Get Hotel Image Name For Delete Function
func GetHotelImageName(r *http.Request, reqMap data.ImageName) (string, error) {
	util.LogIt(r, "model - V_Partner_Hotel - GetHotelImageName")
	var Qry bytes.Buffer

	Qry.WriteString("SELECT image FROM cf_hotel_image WHERE id=? AND hotel_id = ?")
	Data, err := model.ExecuteRowQuery(Qry.String(), reqMap.Image, reqMap.ID)
	if util.CheckErrorLog(r, err) {
		return "", err
	}

	return Data["image"].(string), nil
}

// GetHotelImageCount - Get Hotel Image Count
func GetHotelImageCount(r *http.Request) (int64, error) {
	util.LogIt(r, "model - V_Partner_Hotel - GetHotelImageCount")
	var Qry bytes.Buffer

	Qry.WriteString("SELECT COUNT(id) AS id FROM cf_hotel_image WHERE hotel_id = ?")
	Data, err := model.ExecuteRowQuery(Qry.String(), context.Get(r, "HotelId"))
	if util.CheckErrorLog(r, err) {
		return 0, err
	}

	return Data["id"].(int64), nil
}

// GetHotelImageList - Get Hotel Image List
func GetHotelImageList(r *http.Request, hotelID string) (map[string]interface{}, error) {

	util.LogIt(r, "model - V_Partner_Hotel - GetRoomImageList")

	var mainStuff = make(map[string]interface{})

	//Category Wise Image
	var CatQry bytes.Buffer
	CatQry.WriteString("SELECT category_id,CIC.name AS category FROM cf_hotel_image AS CFI INNER JOIN cf_image_category AS CIC ON CIC.id = CFI.category_id WHERE hotel_id = ? GROUP BY category_id")
	CatArray, err := model.ExecuteQuery(CatQry.String(), hotelID)
	if util.CheckErrorLog(r, err) {
		return nil, err
	}

	if len(CatArray) > 0 {
		for i := 0; i < len(CatArray); i++ {
			var CatQry bytes.Buffer
			CatQry.WriteString("SELECT id, CONCAT('" + config.Env.AwsBucketURL + "hotel/" + "',image) AS image, sortorder FROM cf_hotel_image WHERE hotel_id = ? AND category_id = ? ORDER BY sortorder")
			ImageArray, err := model.ExecuteQuery(CatQry.String(), hotelID, CatArray[i]["category_id"])
			if util.CheckErrorLog(r, err) {
				return nil, err
			}

			CatArray[i]["image"] = ImageArray
			delete(CatArray[i], "category_id")
		}
		mainStuff["hotel_category_image"] = CatArray
	} else {
		mainStuff["hotel_category_image"] = make(map[string]interface{})
	}
	return mainStuff, nil
}

// CheckRoomCount - Checks Room Type Count For Hotel
func CheckRoomCount(r *http.Request, HotelID string) ([]map[string]interface{}, error) {

	util.LogIt(r, "Model - V_Partner_Hotel - CheckRoomCount")
	var Qry bytes.Buffer

	Qry.WriteString("SELECT id, room_type_name, max_occupancy, inventory FROM cf_room_type WHERE hotel_id = ? AND status = ?")
	Data, err := model.ExecuteQuery(Qry.String(), HotelID, 1)
	if util.CheckErrorLog(r, err) {
		return nil, err
	}

	return Data, nil
}

// CheckRatePlanCount - Checks Rate Plan Count For Hotel
func CheckRatePlanCount(r *http.Request, HotelID string, RoomID string) ([]map[string]interface{}, error) {

	util.LogIt(r, "Model - V_Partner_Hotel - CheckRatePlanCount")
	var Qry bytes.Buffer

	Qry.WriteString("SELECT id, rate_plan_name, rate FROM cf_rateplan WHERE hotel_id = ? AND room_type_id = ? AND status = ?")
	Data, err := model.ExecuteQuery(Qry.String(), HotelID, RoomID, 1)
	if util.CheckErrorLog(r, err) {
		return nil, err
	}

	return Data, nil
}

// CheckAndFillInvData - Checks If Data Exists in DB and If not fills data
func CheckAndFillInvData(r *http.Request, HotelID string, roomID string, baseInv int64, startDate string, endDate string) bool {

	util.LogIt(r, "Model - V_Partner_Hotel - CheckAndFillInvData")

	yearMonthDates := model.GetYearMonthSliceBetweenTwoDates(r, startDate, endDate)
	/* 2020-08-11 - Log Purpose */
	localDate := util.GetISODate()

	for _, val := range yearMonthDates {

		yearMonth := fmt.Sprintf("%v", val["year_month"]) // Convert interface to string

		dateData := val["data"]

		dateInfo := strings.Split(yearMonth, "-")
		tblYear := dateInfo[0]
		tblMonth := dateInfo[1]

		var Qry bytes.Buffer
		Qry.WriteString("SELECT inv_data FROM cf_inv_data WHERE hotel_id = ? AND room_id = ? AND year = ? AND month = ? ") // Add
		Data, err := model.ExecuteRowQuery(Qry.String(), HotelID, roomID, tblYear, tblMonth)
		if util.CheckErrorLog(r, err) {
			return false
		}

		if len(Data) == 0 {

			// https://medium.com/@prithvi_20863/interfaces-in-golang-a-short-anecdote-249d7c6f96f4
			dMap := reflect.ValueOf(dateData) // as dateData is interface{} type and to loop over through such type we need such conversion
			var DateArr = make(map[string]int64)
			for _, dateVal := range dMap.MapKeys() {
				// valTest := dMap.MapIndex(dateVal) // to access value of key use this line
				realDate, _ := reflect.Value(dateVal).Interface().(string)
				DateArr[realDate] = baseInv

				// Put Log Here - 2020-08-11 - HK - For Update Log
				var Qry1 bytes.Buffer
				nanoid, _ := gonanoid.Nanoid()
				Qry1.WriteString("INSERT INTO logs_inv (id, hotel_id, room_id, inventory, update_for_date, updated_at, updated_by, booking_id, ip) VALUES (?,?,?,?,?,?,?,?,?)")
				err = model.ExecuteNonQuery(Qry1.String(), nanoid, HotelID, roomID, baseInv, realDate, localDate, "", "", "")
				if util.CheckErrorLog(r, err) {
					return false
				}
				// Put Log Here - 2020-08-11 - HK - For Update Log
			}
			invJSON, err := json.Marshal(DateArr)

			var Qry1 bytes.Buffer
			nanoid, _ := gonanoid.Nanoid()
			Qry1.WriteString("INSERT INTO cf_inv_data (id, hotel_id, room_id, year, month, inv_data) VALUES (?,?,?,?,?,?)")
			err = model.ExecuteNonQuery(Qry1.String(), nanoid, HotelID, roomID, tblYear, tblMonth, string(invJSON))
			if util.CheckErrorLog(r, err) {
				return false
			}
		} else {

			existsData := Data["inv_data"].(string)
			// log.Println(existsData) fmt.Println(reflect.TypeOf(existsData)) // string

			// Declared an empty map interface
			var result map[string]interface{}

			// Unmarshal or Decode the JSON to the interface.
			json.Unmarshal([]byte(existsData), &result)

			// Print the data type of result variable
			// fmt.Println(result) fmt.Println(reflect.TypeOf(result)) // map[string]interface {}

			incomingData := reflect.ValueOf(dateData)
			for _, dateVal := range incomingData.MapKeys() {
				realDate, _ := reflect.Value(dateVal).Interface().(string)

				if _, ok := result[realDate]; !ok {
					//Update Json For Incoming Dates Which Not Exists In Table Inventory Data
					result[realDate] = baseInv

					// Put Log Here - 2020-08-11 - HK - For Update Log
					var Qry1 bytes.Buffer
					nanoid, _ := gonanoid.Nanoid()
					Qry1.WriteString("INSERT INTO logs_inv (id, hotel_id, room_id, inventory, update_for_date, updated_at, updated_by, booking_id, ip) VALUES (?,?,?,?,?,?,?,?,?)")
					err = model.ExecuteNonQuery(Qry1.String(), nanoid, HotelID, roomID, baseInv, realDate, localDate, "", "", "")
					if util.CheckErrorLog(r, err) {
						return false
					}
					// Put Log Here - 2020-08-11 - HK - For Update Log
				}
			}
			invJSON, err := json.Marshal(result)
			// fmt.Println(string(invJSON), err)

			var Qry1 bytes.Buffer
			Qry1.WriteString("UPDATE cf_inv_data SET inv_data = ? WHERE  hotel_id = ? AND  room_id = ? AND year = ? AND month = ?")
			err = model.ExecuteNonQuery(Qry1.String(), string(invJSON), HotelID, roomID, tblYear, tblMonth)
			if util.CheckErrorLog(r, err) {
				return false
			}
		}
		// insert room inventory
	}

	return true
}

// CheckAndFillRateRestrictionData - Checks If Data Exists in DB and If not fills data
func CheckAndFillRateRestrictionData(r *http.Request, HotelID string, roomID string, rateInfo map[string]interface{}) bool {

	util.LogIt(r, "Model - V_Partner_Hotel - CheckAndFillRateRestrictionData")

	startDate := rateInfo["start_date"].(string)
	endDate := rateInfo["end_date"].(string)
	rateID := rateInfo["rate_id"].(string)
	// occupancy := rateInfo["occupancy"].(int64)

	var dataDumpRatePlanWise = make(map[string]interface{})
	// dataDumpRatePlanWise["rate"] = rateInfo["rate"].(string)
	dataDumpRatePlanWise["rate"] = rateInfo["rate"].([]map[string]interface{})
	dataDumpRatePlanWise["min_night"] = rateInfo["min_night"].(int)
	dataDumpRatePlanWise["stop_sell"] = rateInfo["stop_sell"].(int)
	dataDumpRatePlanWise["cta"] = rateInfo["cta"].(int)
	dataDumpRatePlanWise["ctd"] = rateInfo["ctd"].(int)

	yearMonthDates := model.GetYearMonthSliceBetweenTwoDates(r, startDate, endDate)
	/* 2020-08-19 - Log Purpose */
	localDate := util.GetISODate()

	for _, val := range yearMonthDates {

		// insert rate, restriction data
		yearMonth := fmt.Sprintf("%v", val["year_month"]) // Convert interface to string

		dateData := val["data"]

		dateInfo := strings.Split(yearMonth, "-")
		tblYear := dateInfo[0]
		tblMonth := dateInfo[1]

		var Qry bytes.Buffer
		Qry.WriteString("SELECT rate_rest_data FROM cf_rate_restriction_data WHERE hotel_id = ? AND room_id = ? AND rateplan_id = ? AND year = ? AND month = ? ")
		RateRestDataFromTbl, err := model.ExecuteRowQuery(Qry.String(), HotelID, roomID, rateID, tblYear, tblMonth)
		if util.CheckErrorLog(r, err) {
			return false
		}

		if len(RateRestDataFromTbl) == 0 {

			// https://medium.com/@prithvi_20863/interfaces-in-golang-a-short-anecdote-249d7c6f96f4
			dMap := reflect.ValueOf(dateData) // as dateData is interface{} type and to loop over through such type we need such conversion
			var DateDataInsert = make(map[string]interface{})
			for _, dateVal := range dMap.MapKeys() {
				// valTest := dMap.MapIndex(dateVal) // to access value of key use this line
				realDate, _ := reflect.Value(dateVal).Interface().(string)
				DateDataInsert[realDate] = dataDumpRatePlanWise

				// Put Log Here - 2020-08-19 - HK - For Update Log
				rateDataForLog := dataDumpRatePlanWise["rate"]
				var rateStr string
				if x, ok := rateDataForLog.([]interface{}); ok {
					for _, e := range x {
						incomingData := reflect.ValueOf(e)
						for _, element := range incomingData.MapKeys() {
							valTest := incomingData.MapIndex(element) // to access value of key use this line
							occupancy, _ := reflect.Value(element).Interface().(string)
							occWiserate, _ := reflect.Value(valTest).Interface().(string)
							rateStr += occupancy + ":" + occWiserate + ", "
						}
					}
				}
				rateStr = strings.TrimRight(rateStr, ", ")

				var Qry2 bytes.Buffer
				nanoid, _ := gonanoid.Nanoid()
				Qry2.WriteString("INSERT INTO logs_rate_rest (id, hotel_id, room_id, rateplan_id,  update_for_date, rate, min_night, stop_sell, cta, ctd, updated_at, updated_by, ip) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?)")
				err = model.ExecuteNonQuery(Qry2.String(), nanoid, HotelID, roomID, rateID, realDate, rateStr, dataDumpRatePlanWise["min_night"], dataDumpRatePlanWise["stop_sell"], dataDumpRatePlanWise["cta"], dataDumpRatePlanWise["ctd"], localDate, "", "")
				if util.CheckErrorLog(r, err) {
					return false
				}
				// Put Log Here - 2020-08-19 - HK - For Update Log
			}

			rateRestJSON, err := json.Marshal(DateDataInsert)

			var Qry1 bytes.Buffer
			nanoid, _ := gonanoid.Nanoid()
			Qry1.WriteString("INSERT INTO cf_rate_restriction_data (id, hotel_id, room_id, rateplan_id, year, month, rate_rest_data)  VALUES (?,?,?,?,?,?,?)")
			err = model.ExecuteNonQuery(Qry1.String(), nanoid, HotelID, roomID, rateID, tblYear, tblMonth, string(rateRestJSON))
			if util.CheckErrorLog(r, err) {
				return false
			}

		} else {

			existsData := RateRestDataFromTbl["rate_rest_data"].(string)

			// Declared an empty map interface
			var result map[string]interface{}

			// Unmarshal or Decode the JSON to the interface.
			json.Unmarshal([]byte(existsData), &result)

			// Print the data type of result variable
			// fmt.Println(result) fmt.Println(reflect.TypeOf(result)) // map[string]interface {}

			incomingData := reflect.ValueOf(dateData)
			for _, dateVal := range incomingData.MapKeys() {

				realDate, _ := reflect.Value(dateVal).Interface().(string)
				if _, ok := result[realDate]; !ok {
					// Update Json For Incoming Dates Which Not Exists In Table Rate Restrcition Data
					result[realDate] = dataDumpRatePlanWise

					// Put Log Here - 2020-08-19 - HK - For Update Log
					rateDataForLog := dataDumpRatePlanWise["rate"]
					var rateStr string
					if x, ok := rateDataForLog.([]interface{}); ok {
						for _, e := range x {
							incomingData := reflect.ValueOf(e)
							for _, element := range incomingData.MapKeys() {
								valTest := incomingData.MapIndex(element) // to access value of key use this line
								occupancy, _ := reflect.Value(element).Interface().(string)
								occWiserate, _ := reflect.Value(valTest).Interface().(string)
								rateStr += occupancy + ":" + occWiserate + ", "
							}
						}
					}
					rateStr = strings.TrimRight(rateStr, ", ")

					var Qry2 bytes.Buffer
					nanoid, _ := gonanoid.Nanoid()
					Qry2.WriteString("INSERT INTO logs_rate_rest (id, hotel_id, room_id, rateplan_id,  update_for_date, rate, min_night, stop_sell, cta, ctd, updated_at, updated_by, ip) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?)")
					err = model.ExecuteNonQuery(Qry2.String(), nanoid, HotelID, roomID, rateID, realDate, rateStr, dataDumpRatePlanWise["min_night"], dataDumpRatePlanWise["stop_sell"], dataDumpRatePlanWise["cta"], dataDumpRatePlanWise["ctd"], localDate, "", "")
					if util.CheckErrorLog(r, err) {
						return false
					}
					// Put Log Here - 2020-08-19 - HK - For Update Log
				}
			}

			rateJSON, err := json.Marshal(result)

			var Qry1 bytes.Buffer
			Qry1.WriteString("UPDATE cf_rate_restriction_data SET rate_rest_data = ? WHERE  hotel_id = ? AND  room_id = ? AND rateplan_id = ? AND year = ? AND month = ?")
			err = model.ExecuteNonQuery(Qry1.String(), string(rateJSON), HotelID, roomID, rateID, tblYear, tblMonth)
			if util.CheckErrorLog(r, err) {
				return false
			}
		}
	}
	return true
}

// FillInvRateData - Fills Inv, Rate, Restrictions Data Occupancy Wise For hotel
func FillInvRateData(r *http.Request, HotelID string, reqMap map[string]interface{}) bool {

	util.LogIt(r, "Model - V_Partner_Hotel - FillInvRateData")

	roomInfo := reqMap["room_info"].([]map[string]interface{})
	// log.Println(roomInfo)

	if len(roomInfo) > 0 {

		for i := 0; i < len(roomInfo); i++ {

			roomID := roomInfo[i]["room_id"].(string)
			baseInv := roomInfo[i]["inventory"].(int64)
			startDate := roomInfo[i]["start_date"].(string)
			endDate := roomInfo[i]["end_date"].(string)

			// log.Println(HotelID, roomID, baseInv, startDate, endDate)
			invSuccess := CheckAndFillInvData(r, HotelID, roomID, baseInv, startDate, endDate)
			if !invSuccess {
				return false
			}

			rateInfo := roomInfo[i]["rate_info"].([]map[string]interface{})
			if len(rateInfo) > 0 {
				for i := 0; i < len(rateInfo); i++ {
					occWiseRateInfoData := rateInfo[i]
					rateSuccess := CheckAndFillRateRestrictionData(r, HotelID, roomID, occWiseRateInfoData)
					if !rateSuccess {
						return false
					}
				}
			}
		} // end for i := 0; i < len(roomInfo); i++
	}
	return true
}

// GetHotelRoomRateData - Gets Inv, Rate, Restrictions Data Occupancy Wise For hotel
func GetHotelRoomRateData(r *http.Request, reqMap map[string]interface{}) ([]map[string]interface{}, error) {

	util.LogIt(r, "Model - V_Partner_Hotel - GetHotelRoomRateData")

	ReqMonth := reqMap["month"]
	ReqYear := reqMap["year"]
	roomID := reqMap["room_id"]
	rateID := reqMap["rate_id"]

	if ReqMonth == "" || ReqYear == "" {
		currentTime := time.Now()
		ReqYear = currentTime.Year()
		ReqMonth = currentTime.Month()
	}

	// HotelID := context.Get(r, "HotelId")
	HotelID := reqMap["hotel_id"].(string)
	context.Set(r, "HotelId", HotelID)
	var FinalArray []map[string]interface{}

	if roomID != "" && rateID != "" {

		// Get Room Info
		roomInfo, err := GetRoomType(r, roomID.(string))
		if err != nil {
			return nil, err
		}

		RoomID := roomInfo["id"].(string)
		RoomName := roomInfo["room_type_name"].(string)
		RoomOccupancy := roomInfo["max_occupancy"].(int64)

		var Qry bytes.Buffer
		Qry.WriteString("SELECT inv_data as data FROM cf_inv_data WHERE hotel_id = ? AND room_id = ? AND year = ? AND month = ? ") // Add
		Data, err := model.ExecuteRowQuery(Qry.String(), HotelID, RoomID, ReqYear, ReqMonth)
		if util.CheckErrorLog(r, err) {
			util.LogIt(r, "Model - V_Partner_Hotel - GetHotelRoomRateData if Error Getting Data For This Hotel And Room And Year And Month")
			util.LogIt(r, err)
			util.LogIt(r, HotelID+" - "+RoomID+" - "+fmt.Sprintf("%v", ReqYear)+" - "+fmt.Sprintf("%v", ReqMonth))
			return nil, err
		}

		if len(Data) == 0 {
			util.LogIt(r, "Model - V_Partner_Hotel - GetHotelRoomRateData if No Data Found For This Hotel And Room And Year And Month")
			util.LogIt(r, err)
			util.LogIt(r, HotelID+" - "+RoomID+" - "+fmt.Sprintf("%v", ReqYear)+" - "+fmt.Sprintf("%v", ReqMonth))
			return nil, err
		}

		jsonMap := make(map[string]int64)
		err = json.Unmarshal([]byte(Data["data"].(string)), &jsonMap)
		if util.CheckErrorLog(r, err) {
			return nil, err

		}
		FinalArray = append(FinalArray, map[string]interface{}{
			"room_id":   RoomID,
			"room_name": RoomName,
			"occupancy": RoomOccupancy,
			"data":      jsonMap,
		})

		rateInfo, err := GetRatePlan(r, rateID.(string), HotelID)
		if err != nil {
			return nil, err
		}

		RateID := rateInfo["id"].(string)
		RateName := rateInfo["rate_plan_name"].(string)

		var rateArr []map[string]interface{}

		var Qry1 bytes.Buffer
		Qry1.WriteString("SELECT rate_rest_data as data FROM cf_rate_restriction_data WHERE hotel_id = ? AND room_id = ? AND rateplan_id = ? AND year = ? AND month = ? ") // Add
		Data1, err := model.ExecuteRowQuery(Qry1.String(), HotelID, RoomID, RateID, ReqYear, ReqMonth)
		if util.CheckErrorLog(r, err) {
			return nil, err
		}

		jsonMap1 := make(map[string]interface{})
		err = json.Unmarshal([]byte(Data1["data"].(string)), &jsonMap1)
		if util.CheckErrorLog(r, err) {
			return nil, err

		}

		rateArr = append(rateArr, map[string]interface{}{
			"rate_id":   RateID,
			"rate_name": RateName,
			"data":      jsonMap1,
		})

		FinalArray[0]["rate_info"] = rateArr

	} else if roomID != "" && rateID == "" {
		// Get Room Info
		roomInfo, err := GetRoomType(r, roomID.(string))
		if err != nil {
			return nil, err
		}

		RoomID := roomInfo["id"].(string)
		RoomName := roomInfo["room_type_name"].(string)
		RoomOccupancy := roomInfo["max_occupancy"].(int64)

		var Qry bytes.Buffer
		Qry.WriteString("SELECT inv_data as data FROM cf_inv_data WHERE hotel_id = ? AND room_id = ? AND year = ? AND month = ? ") // Add
		Data, err := model.ExecuteRowQuery(Qry.String(), HotelID, RoomID, ReqYear, ReqMonth)
		if util.CheckErrorLog(r, err) {
			util.LogIt(r, "Model - V_Partner_Hotel - GetHotelRoomRateData else if Error Getting Data For This Hotel And Room And Year And Month")
			util.LogIt(r, err)
			util.LogIt(r, HotelID+" - "+RoomID+" - "+fmt.Sprintf("%v", ReqYear)+" - "+fmt.Sprintf("%v", ReqMonth))
			return nil, err
		}

		if len(Data) == 0 {
			util.LogIt(r, "Model - V_Partner_Hotel - GetHotelRoomRateData else if No Data Found For This Hotel And Room And Year And Month")
			util.LogIt(r, err)
			util.LogIt(r, HotelID+" - "+RoomID+" - "+fmt.Sprintf("%v", ReqYear)+" - "+fmt.Sprintf("%v", ReqMonth))
			return nil, err
		}

		jsonMap := make(map[string]int64)
		err = json.Unmarshal([]byte(Data["data"].(string)), &jsonMap)
		if util.CheckErrorLog(r, err) {
			return nil, err

		}
		FinalArray = append(FinalArray, map[string]interface{}{
			"room_id":   RoomID,
			"room_name": RoomName,
			"occupancy": RoomOccupancy,
			"data":      jsonMap,
		})

		rateDATA, err := CheckRatePlanCount(r, HotelID, RoomID)
		if len(rateDATA) == 0 || err != nil {
			util.LogIt(r, "Controller - V_Partner_Hotel - GetHotelRoomRateData No Rate Plans Found For This Hotel")
			return nil, err
		}

		var rateArr []map[string]interface{}
		for _, v2 := range rateDATA {

			RateID := v2["id"].(string)
			RateName := v2["rate_plan_name"].(string)

			var Qry1 bytes.Buffer
			Qry1.WriteString("SELECT rate_rest_data as data FROM cf_rate_restriction_data WHERE hotel_id = ? AND room_id = ? AND rateplan_id = ? AND year = ? AND month = ? ") // Add
			Data1, err := model.ExecuteRowQuery(Qry1.String(), HotelID, RoomID, RateID, ReqYear, ReqMonth)
			if util.CheckErrorLog(r, err) {
				return nil, err
			}

			jsonMap1 := make(map[string]interface{})
			err = json.Unmarshal([]byte(Data1["data"].(string)), &jsonMap1)
			if util.CheckErrorLog(r, err) {
				return nil, err

			}

			rateArr = append(rateArr, map[string]interface{}{
				"rate_id":   RateID,
				"rate_name": RateName,
				"data":      jsonMap1,
			})

		}
		FinalArray[0]["rate_info"] = rateArr
	} else {

		// Room Array Preparation - START
		roomDATA, err := CheckRoomCount(r, HotelID)
		if len(roomDATA) == 0 || err != nil {
			util.LogIt(r, "Model - V_Partner_Hotel - GetHotelRoomRateData No Rooms Found For This Hotel")
			return nil, err
		}

		// 2021-04-27 - HK - START
		// Purpose : Only Those Rooms Are Considered Whose Rate Plan Exists. So Deals Page
		// Can Work With All Drop Down Value
		var froomDATA []map[string]interface{}
		for _, v := range roomDATA {
			rateDATA, err := CheckRatePlanCount(r, HotelID, v["id"].(string))
			if len(rateDATA) == 0 || err != nil {
				util.LogIt(r, "Controller - V_Partner_Hotel - GetHotelRoomRateData No Rate Plans Found For This Room")
				util.LogIt(r, v["room_type_name"].(string))
			}
			if len(rateDATA) > 0 {
				froomDATA = append(froomDATA, v)
			}
		}
		// 2021-04-27 - HK - END

		for _, v := range froomDATA {

			RoomID := v["id"].(string)
			RoomName := v["room_type_name"].(string)
			RoomOccupancy := v["max_occupancy"].(int64)

			var Qry bytes.Buffer
			Qry.WriteString("SELECT inv_data as data FROM cf_inv_data WHERE hotel_id = ? AND room_id = ? AND year = ? AND month = ? ") // Add
			Data, err := model.ExecuteRowQuery(Qry.String(), HotelID, RoomID, ReqYear, ReqMonth)
			if util.CheckErrorLog(r, err) {
				util.LogIt(r, "Model - V_Partner_Hotel - GetHotelRoomRateData else Error Getting Data For This Hotel And Room And Year And Month")
				util.LogIt(r, err)
				util.LogIt(r, HotelID+" - "+RoomID+" - "+fmt.Sprintf("%v", ReqYear)+" - "+fmt.Sprintf("%v", ReqMonth))
				return nil, err
			}
			if len(Data) == 0 {
				util.LogIt(r, "Model - V_Partner_Hotel - GetHotelRoomRateData else No Data Found For This Hotel And Room And Year And Month")
				util.LogIt(r, err)
				util.LogIt(r, HotelID+" - "+RoomID+" - "+fmt.Sprintf("%v", ReqYear)+" - "+fmt.Sprintf("%v", ReqMonth))
				return nil, err
			}

			jsonMap := make(map[string]int64)
			err = json.Unmarshal([]byte(Data["data"].(string)), &jsonMap)
			if util.CheckErrorLog(r, err) {
				return nil, err

			}
			FinalArray = append(FinalArray, map[string]interface{}{
				"room_id":   RoomID,
				"room_name": RoomName,
				"occupancy": RoomOccupancy,
				"data":      jsonMap,
			})
		}

		for k1, v1 := range FinalArray {

			fRoomID := fmt.Sprintf("%v", fmt.Sprintf("%v", v1["room_id"]))
			rateDATA, err := CheckRatePlanCount(r, HotelID, fRoomID)
			if len(rateDATA) == 0 || err != nil {
				util.LogIt(r, "Controller - V_Partner_Hotel - GetHotelRoomRateData No Rate Plans Found For This Hotel")
				return nil, err
			}

			var rateArr []map[string]interface{}
			for _, v2 := range rateDATA {
				// log.Println(v2)

				RateID := v2["id"].(string)
				RateName := v2["rate_plan_name"].(string)

				var Qry bytes.Buffer
				Qry.WriteString("SELECT rate_rest_data as data FROM cf_rate_restriction_data WHERE hotel_id = ? AND room_id = ? AND rateplan_id = ? AND year = ? AND month = ? ") // Add
				Data, err := model.ExecuteRowQuery(Qry.String(), HotelID, fRoomID, RateID, ReqYear, ReqMonth)
				if util.CheckErrorLog(r, err) {
					return nil, err
				}

				jsonMap := make(map[string]interface{})
				err = json.Unmarshal([]byte(Data["data"].(string)), &jsonMap)
				if util.CheckErrorLog(r, err) {
					return nil, err

				}

				rateArr = append(rateArr, map[string]interface{}{
					"rate_id":   RateID,
					"rate_name": RateName,
					"data":      jsonMap,
				})

			}
			FinalArray[k1]["rate_info"] = rateArr
		}
	}
	return FinalArray, nil
}

// UpdateHotelRoomRateData - Update Hotel's Room Rate Data
func UpdateHotelRoomRateData(r *http.Request, reqMap map[string]interface{}) bool {

	util.LogIt(r, "Model - V_Partner_Hotel - UpdateHotelRoomRateData")

	tblYear := reqMap["year"]
	tblMonth := reqMap["month"]
	updateData := reqMap["data"].([]interface{})
	// HotelID := context.Get(r, "HotelId")
	HotelID := reqMap["hotel_id"].(string)

	/* 2020-08-08 - Log Purpose */
	localDate := util.GetISODate()
	var UserID string
	if _, ok := context.Get(r, "UserId").(string); ok {
		UserID = context.Get(r, "UserId").(string)
	}
	VisiterIP := context.Get(r, "Visitor_IP")
	/* 2020-08-08 - Log Purpose */

	if len(updateData) > 0 {

		for i := 0; i < len(updateData); i++ {

			allData := updateData[i].(map[string]interface{})

			// --------------------------------------------		Inv Update Start	-------------------------------------------- //
			RoomID := allData["room_id"].(string)
			// RoomName := allData["room_name"].(string)
			ReqRoomData := allData["data"].(map[string]interface{})

			var Qry bytes.Buffer
			Qry.WriteString("SELECT inv_data as tblData FROM cf_inv_data WHERE hotel_id = ? AND room_id = ? AND year = ? AND month = ? ") // Add
			QryRoomData, err := model.ExecuteRowQuery(Qry.String(), HotelID, RoomID, tblYear, tblMonth)
			if util.CheckErrorLog(r, err) {
				return false
			}

			TblRoomData := QryRoomData["tblData"].(string)

			// Declared an empty map interface
			var invFinalResult map[string]interface{}

			// Unmarshal or Decode the JSON to the interface.
			json.Unmarshal([]byte(TblRoomData), &invFinalResult)

			// Loop Through Request Inventory Data And Check If Inv Data Is Not Same As Table Inv Data Then Replace Inv To Table Data Of That Date
			for reqDate, reqDateInv := range ReqRoomData {
				if ReqRoomData[reqDate] != invFinalResult[reqDate] {
					invFinalResult[reqDate] = reqDateInv
					// Put Log Here - 2020-08-07 - HK - For Update Log
					var Qry1 bytes.Buffer
					nanoid, _ := gonanoid.Nanoid()
					Qry1.WriteString("INSERT INTO logs_inv (id, hotel_id, room_id, inventory, update_for_date, updated_at, updated_by, booking_id, ip) VALUES (?,?,?,?,?,?,?,?,?)")
					err = model.ExecuteNonQuery(Qry1.String(), nanoid, HotelID, RoomID, reqDateInv, reqDate, localDate, UserID, "", VisiterIP)
					if util.CheckErrorLog(r, err) {
						return false
					}
					// Put Log Here - 2020-08-07 - HK - For Update Log
				}
			}

			invJSON, err := json.Marshal(invFinalResult)
			// fmt.Println(string(invJSON), err)

			var Qry1 bytes.Buffer
			Qry1.WriteString("UPDATE cf_inv_data SET inv_data = ? WHERE  hotel_id = ? AND  room_id = ? AND year = ? AND month = ?")
			err = model.ExecuteNonQuery(Qry1.String(), string(invJSON), HotelID, RoomID, tblYear, tblMonth)
			if util.CheckErrorLog(r, err) {
				return false
			}
			// --------------------------------------------		Inv Update END	-------------------------------------------- //

			// --------------------------------------------		Rate Res Start	-------------------------------------------- //
			ReqRateRestData := allData["rate_info"].([]interface{})
			// log.Println(ReqRateRestData)

			if len(ReqRateRestData) > 0 {

				for i := 0; i < len(ReqRateRestData); i++ {

					ReqParticularRateRestData := ReqRateRestData[i].(map[string]interface{})

					RateID := ReqParticularRateRestData["rate_id"].(string)
					// RateName := ReqParticularRateRestData["rate_name"].(string)
					ReqRateRestData := ReqParticularRateRestData["data"].(map[string]interface{})

					var Qry2 bytes.Buffer
					Qry2.WriteString("SELECT rate_rest_data as tblData FROM cf_rate_restriction_data WHERE hotel_id = ? AND room_id = ? AND rateplan_id = ? AND year = ? AND month = ? ")
					QryRateRestData, err := model.ExecuteRowQuery(Qry2.String(), HotelID, RoomID, RateID, tblYear, tblMonth)
					if util.CheckErrorLog(r, err) {
						return false
					}

					TblRateRestData := QryRateRestData["tblData"].(string)

					// Declared an empty map interface
					var rateRestFinalresult map[string]interface{}

					// Unmarshal or Decode the JSON to the interface.
					json.Unmarshal([]byte(TblRateRestData), &rateRestFinalresult)

					// Loop Through Request Rate Rest Data And Check If Rate Rest Data Is Not Same As Table Rate Rest Data Then Replace Those Data To Table Data Of That Date
					for reqDate := range ReqRateRestData {

						DateWiseReqRateRestData := ReqRateRestData[reqDate].(map[string]interface{})
						DateWiseTblRateRestData := rateRestFinalresult[reqDate].(map[string]interface{})

						// Checkes If Date Have Same Data Of Rate, StopSell, CTA, CTD, Min Nights
						res := reflect.DeepEqual(DateWiseReqRateRestData, DateWiseTblRateRestData)
						if res == false {
							rateRestFinalresult[reqDate] = DateWiseReqRateRestData

							// Put Log Here - 2020-08-08 - HK - For Update Log
							rateDataForLog := DateWiseReqRateRestData["rate"]
							var rateStr string
							if x, ok := rateDataForLog.([]interface{}); ok {
								for _, e := range x {
									incomingData := reflect.ValueOf(e)
									for _, element := range incomingData.MapKeys() {
										valTest := incomingData.MapIndex(element) // to access value of key use this line
										occupancy, _ := reflect.Value(element).Interface().(string)
										occWiserate, _ := reflect.Value(valTest).Interface().(string)
										rateStr += occupancy + ":" + occWiserate + ", "
									}
								}
							}
							rateStr = strings.TrimRight(rateStr, ", ")

							// rateJSON, _ := json.Marshal(DateWiseReqRateRestData["rate"])
							var Qry2 bytes.Buffer
							nanoid, _ := gonanoid.Nanoid()
							Qry2.WriteString("INSERT INTO logs_rate_rest (id, hotel_id, room_id, rateplan_id,  update_for_date, rate, min_night, stop_sell, cta, ctd, updated_at, updated_by, ip) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?)")
							err = model.ExecuteNonQuery(Qry2.String(), nanoid, HotelID, RoomID, RateID, reqDate, rateStr, DateWiseReqRateRestData["min_night"], DateWiseReqRateRestData["stop_sell"], DateWiseReqRateRestData["cta"], DateWiseReqRateRestData["ctd"], localDate, UserID, VisiterIP)
							if util.CheckErrorLog(r, err) {
								return false
							}
							// Put Log Here - 2020-08-08 - HK - For Update Log
						}
					}
					rateRestJSON, err := json.Marshal(rateRestFinalresult)
					// fmt.Println(string(rateRestJSON), err)

					var Qry4 bytes.Buffer
					Qry4.WriteString("UPDATE cf_rate_restriction_data SET rate_rest_data = ? WHERE hotel_id = ? AND room_id = ? AND rateplan_id = ? AND year = ? AND month = ?")
					err = model.ExecuteNonQuery(Qry4.String(), string(rateRestJSON), HotelID, RoomID, RateID, tblYear, tblMonth)
					if util.CheckErrorLog(r, err) {
						return false
					}

					model.CacheChn <- model.CacheObj{
						Type:        "updateDeals",
						ID:          HotelID,
						Additional:  RoomID,
						Additional1: RateID,
					}
				}

			}
			// --------------------------------------------		Rate Res Start	-------------------------------------------- //
		}
	}

	return true
}

// CheckIfHotelExists - Checks If Hotel Exists With Active Status - 2020-08-04 - HK
func CheckIfHotelExists(r *http.Request, HotelID string) (int64, error) {

	util.LogIt(r, "model - V_Partner_Hotel - CheckIfHotelExists")
	var Qry bytes.Buffer

	Qry.WriteString("SELECT COUNT(id) AS id FROM cf_hotel_info WHERE id = ? AND status = ? AND is_live = ?")
	Data, err := model.ExecuteRowQuery(Qry.String(), HotelID, 1, 1)
	if util.CheckErrorLog(r, err) {
		return -1, err
	}

	return Data["id"].(int64), nil
}

// CheckIfRoomExists - Checks If Hotel Exists With Active Status - 2020-08-04 - HK
func CheckIfRoomExists(r *http.Request, roomID string, HotelID string) (int64, error) {
	util.LogIt(r, "model - V_Partner_Hotel - CheckIfRoomExists")
	var Qry bytes.Buffer

	Qry.WriteString("SELECT COUNT(id) AS id FROM cf_room_type WHERE id = ? AND hotel_id = ? AND status = ?")
	Data, err := model.ExecuteRowQuery(Qry.String(), roomID, HotelID, 1)
	if util.CheckErrorLog(r, err) {
		return -1, err
	}

	return Data["id"].(int64), nil
}

// UpdateHotelRoomRateDataForBooking - Update Hotel's Room Inv Data - 2020-08-04 - HK
func UpdateHotelRoomRateDataForBooking(r *http.Request, reqMap map[string]interface{}) bool {
	util.LogIt(r, "Model - V_Partner_Hotel - UpdateHotelRoomRateDataForBooking")

	HotelID := reqMap["hotel_id"].(string)
	updateData := reqMap["data"].([]interface{})

	// 2020-08-11 - To Store Update Log Details
	localDate := util.GetISODate()
	var BookingID string
	if _, ok := reqMap["booking_id"]; ok {
		BookingID = reqMap["booking_id"].(string)
	}
	/*var UserID string
	if _, ok := context.Get(r, "UserId").(string); ok {
		UserID = context.Get(r, "UserId").(string)
	}
	VisiterIP := context.Get(r, "Visitor_IP")*/
	// 2020-08-11 - Log Purpose

	if len(updateData) > 0 {
		for i := 0; i < len(updateData); i++ {
			allData := updateData[i].(map[string]interface{})
			RoomID := allData["room_id"].(string)
			startDate := allData["start_date"].(string)
			endDate := allData["end_date"].(string)
			roomCnt := allData["room_cnt"].(float64)
			actionPerf := allData["action"].(string)
			yearMonthDates := model.GetYearMonthSliceBetweenTwoDates(r, startDate, endDate)
			updateJSON, _ := json.Marshal(yearMonthDates)
			util.SysLogIt("YearMonthDatesPrepared For Update ::")
			util.SysLogIt(string(updateJSON))
			// e.g. O/P [{"data":{"2020-08-30":1,"2020-08-31":1},"year_month":"2020-08"},{"data":{"2020-09-01":1,"2020-09-02":1},"year_month":"2020-09"}]
			for _, val := range yearMonthDates {
				yearMonth := fmt.Sprintf("%v", val["year_month"]) // Convert interface to string
				dateData := val["data"]
				dateInfo := strings.Split(yearMonth, "-")
				tblYear := dateInfo[0]
				tblMonth := dateInfo[1]

				var Qry bytes.Buffer
				Qry.WriteString("SELECT inv_data FROM cf_inv_data WHERE hotel_id = ? AND room_id = ? AND year = ? AND month = ? ")
				Data, err := model.ExecuteRowQuery(Qry.String(), HotelID, RoomID, tblYear, tblMonth)
				if util.CheckErrorLog(r, err) {
					return false
				}
				existsData := Data["inv_data"].(string)
				util.SysLogIt("Data From Database For Year :: " + tblYear + " And Month :: " + tblMonth + " For Hotel ID ::" + HotelID + " Room ID :: " + RoomID)
				util.SysLogIt(existsData)

				// Declared an empty map interface
				var result map[string]interface{}

				// Unmarshal or Decode the JSON to the interface.
				json.Unmarshal([]byte(existsData), &result)

				incomingData := reflect.ValueOf(dateData)
				for _, dateVal := range incomingData.MapKeys() {
					realDate, _ := reflect.Value(dateVal).Interface().(string)
					if _, ok := result[realDate]; ok {
						if actionPerf == "block" {
							result[realDate] = result[realDate].(float64) - roomCnt
						} else if actionPerf == "unblock" {
							result[realDate] = result[realDate].(float64) + roomCnt
						}

						// 2020-08-11 - HK - START
						// To Store Inventory Update Logs Affected By Booking Dates
						var Qry1 bytes.Buffer
						nanoid, _ := gonanoid.Nanoid()
						Qry1.WriteString("INSERT INTO logs_inv (id, hotel_id, room_id, inventory, update_for_date, updated_at, updated_by, booking_id, ip) VALUES (?,?,?,?,?,?,?,?,?)")
						err = model.ExecuteNonQuery(Qry1.String(), nanoid, HotelID, RoomID, result[realDate], realDate, localDate, "", BookingID, "")
						if util.CheckErrorLog(r, err) {
							return false
						}
						// 2020-08-11 - HK - END

					}
				}
				invJSON, err := json.Marshal(result)
				util.SysLogIt("Updated Data To Database For Year :: " + tblYear + " And Month :: " + tblMonth + " For Hotel ID ::" + HotelID + " Room ID :: " + RoomID)
				util.SysLogIt(string(invJSON))

				var Qry1 bytes.Buffer
				Qry1.WriteString("UPDATE cf_inv_data SET inv_data = ? WHERE  hotel_id = ? AND  room_id = ? AND year = ? AND month = ?")
				err = model.ExecuteNonQuery(Qry1.String(), string(invJSON), HotelID, RoomID, tblYear, tblMonth)
				if util.CheckErrorLog(r, err) {
					return false
				}
			} // for _, val := range yearMonthDates
		} // for i := 0; i < len(updateData); i++ {
	} // if len(updateData) > 0
	return true
}

// CheckIfRatePlanExists - Checks If Hotel RatePlan Exists With Active Status - 2020-08-08 - HK
func CheckIfRatePlanExists(r *http.Request, roomID string, rateID string, HotelID string) (int64, error) {
	util.LogIt(r, "model - V_Partner_Hotel - CheckIfRatePlanExists")
	var Qry bytes.Buffer

	Qry.WriteString("SELECT COUNT(id) AS id FROM cf_room_type WHERE id = ? AND hotel_id = ? AND room_type_id = ? AND status = ?")
	Data, err := model.ExecuteRowQuery(Qry.String(), roomID, HotelID, 1)
	if util.CheckErrorLog(r, err) {
		return -1, err
	}

	return Data["id"].(int64), nil
}

// GetUpdateLogsOfProperty - Gets All Update Logs Data like Inv, Rate, Min, SS, CTA, CTD Of Room Type And Rate Plan Data Of Hotel
func GetUpdateLogsOfProperty(r *http.Request, reqMap map[string]interface{}) (map[string]interface{}, error) {
	util.LogIt(r, "Model - V_Partner_Hotel - GetUpdateLogsOfProperty")

	stuff := make(map[string]interface{})
	reqFor := reqMap["logs_for"].(string)

	// HotelID := context.Get(r, "HotelId").(string)
	HotelID := reqMap["hotel_id"].(string)

	var results []map[string]interface{}

	if reqFor == "inv" {

		RoomID := reqMap["room_id"].(string)
		if RoomID == "" {
			stuff["data"] = results
			stuff["recordsFiltered"] = 0
			stuff["recordsTotal"] = 0
		} else {
			var SQLQrySel bytes.Buffer
			SQLQrySel.WriteString("SELECT")

			var SQLQryCnt bytes.Buffer
			SQLQryCnt.WriteString(" count(LI.id) AS count ")

			var SQLQryCol bytes.Buffer
			SQLQryCol.WriteString(" LI.id, LI.update_for_date, LI.inventory, from_unixtime(LI.updated_at) as updated_at, LI.ip, ")
			SQLQryCol.WriteString(" CRT.room_type_name, ")
			SQLQryCol.WriteString(" CASE WHEN LI.booking_id IS NULL OR LI.booking_id = '' THEN CHC.client_name ELSE CONCAT('Booking ',' ',LI.booking_id) END AS updated_by ")

			var SQLQryFrom bytes.Buffer
			SQLQryFrom.WriteString(" FROM ")
			SQLQryFrom.WriteString(" logs_inv AS LI ")
			SQLQryFrom.WriteString(" LEFT JOIN ")
			SQLQryFrom.WriteString(" cf_hotel_client AS CHC on CHC.id = LI.updated_by ")
			SQLQryFrom.WriteString(" LEFT JOIN ")
			SQLQryFrom.WriteString(" cf_room_type AS CRT on CRT.id = LI.room_id ")
			SQLQryFrom.WriteString(" WHERE LI.hotel_id = ")
			SQLQryFrom.WriteString("'" + HotelID + "'")
			SQLQryFrom.WriteString(" AND LI.room_id = ")
			SQLQryFrom.WriteString("'" + reqMap["room_id"].(string) + "'")

			var SQLCondition bytes.Buffer
			if reqMap["date"].(string) != "" {
				SQLCondition.WriteString(" AND LI.update_for_date = ")
				SQLCondition.WriteString("'" + reqMap["date"].(string) + "'")
			}

			var SQLGroupOrder bytes.Buffer
			SQLGroupOrder.WriteString(" ORDER BY updated_at DESC ")

			var SQLLimit bytes.Buffer
			SQLLimit.WriteString(" LIMIT ? OFFSET ?")

			var SelectQry bytes.Buffer
			SelectQry.WriteString(SQLQrySel.String())
			SelectQry.WriteString(SQLQryCol.String())
			SelectQry.WriteString(SQLQryFrom.String())
			SelectQry.WriteString(SQLCondition.String())
			SelectQry.WriteString(SQLGroupOrder.String())
			SelectQry.WriteString(SQLLimit.String())

			var SelectFilterQry bytes.Buffer
			SelectFilterQry.WriteString(SQLQrySel.String())
			SelectFilterQry.WriteString(SQLQryCnt.String())
			SelectFilterQry.WriteString(SQLQryFrom.String())
			SelectFilterQry.WriteString(SQLCondition.String())

			var TotalCnt bytes.Buffer
			TotalCnt.WriteString(SQLQrySel.String())
			TotalCnt.WriteString(SQLQryCnt.String())
			TotalCnt.WriteString(SQLQryFrom.String())

			// invUpdateData, err := model.ExecuteQuery(SelectQry.String(), reqMap["length"], reqMap["start"])
			invUpdateData, err := model.ExecuteQuery(SelectQry.String(), fmt.Sprintf("%.0f", reqMap["length"]), fmt.Sprintf("%.0f", reqMap["start"]))
			if util.CheckErrorLog(r, err) {
				return nil, err
			}
			if len(invUpdateData) == 0 {
				stuff["data"] = results
			} else {
				stuff["data"] = invUpdateData
			}

			// log.Println("SelectFilterQry :: ", SelectFilterQry.String())
			FilterTotalCnt, err := model.ExecuteRowQuery(SelectFilterQry.String())
			if util.CheckErrorLog(r, err) {
				return nil, err
			}
			stuff["recordsFiltered"] = FilterTotalCnt["count"]

			// log.Println("TotalCnt :: ", TotalCnt.String())
			TotalCount, err := model.ExecuteRowQuery(TotalCnt.String())
			if util.CheckErrorLog(r, err) {
				return nil, err
			}
			stuff["recordsTotal"] = TotalCount["count"]
		}
	} else if reqFor == "rates" {

		RoomID := reqMap["room_id"].(string)
		if RoomID == "" {
			var results []map[string]interface{}
			stuff["data"] = results
			stuff["recordsFiltered"] = 0
			stuff["recordsTotal"] = 0
		} else {
			var SQLQrySel bytes.Buffer
			SQLQrySel.WriteString("SELECT")

			var SQLQryCnt bytes.Buffer
			SQLQryCnt.WriteString(" count(LRR.id) AS count ")

			var SQLQryCol bytes.Buffer
			SQLQryCol.WriteString(" LRR.id, LRR.update_for_date, LRR.rate, LRR.min_night, LRR.stop_sell, LRR.cta, LRR.ctd, from_unixtime(LRR.updated_at) as updated_at, LRR.ip, ")
			SQLQryCol.WriteString(" CRT.room_type_name, ")
			SQLQryCol.WriteString(" CR.rate_plan_name, ")
			SQLQryCol.WriteString(" CHC.client_name AS updated_by ")

			var SQLQryFrom bytes.Buffer
			SQLQryFrom.WriteString(" FROM ")
			SQLQryFrom.WriteString(" logs_rate_rest AS LRR ")
			SQLQryFrom.WriteString(" LEFT JOIN ")
			SQLQryFrom.WriteString(" cf_hotel_client AS CHC on CHC.id = LRR.updated_by ")
			SQLQryFrom.WriteString(" LEFT JOIN ")
			SQLQryFrom.WriteString(" cf_room_type AS CRT on CRT.id = LRR.room_id ")
			SQLQryFrom.WriteString(" LEFT JOIN ")
			SQLQryFrom.WriteString(" cf_rateplan AS CR on CR.id = LRR.rateplan_id AND CR.room_type_id = LRR.room_id ")
			SQLQryFrom.WriteString(" WHERE LRR.hotel_id = ")
			SQLQryFrom.WriteString("'" + HotelID + "'")
			SQLQryFrom.WriteString(" AND LRR.room_id = ")
			SQLQryFrom.WriteString("'" + reqMap["room_id"].(string) + "'")

			var SQLCondition bytes.Buffer
			if reqMap["rate_id"].(string) != "" {
				SQLCondition.WriteString(" AND LRR.rateplan_id = ")
				SQLCondition.WriteString("'" + reqMap["rate_id"].(string) + "'")
			}
			if reqMap["date"].(string) != "" {
				SQLCondition.WriteString(" AND LRR.update_for_date = ")
				SQLCondition.WriteString("'" + reqMap["date"].(string) + "'")
			}

			var SQLGroupOrder bytes.Buffer
			SQLGroupOrder.WriteString(" ORDER BY updated_at DESC ")

			var SQLLimit bytes.Buffer
			SQLLimit.WriteString(" LIMIT ? OFFSET ?")

			var SelectQry bytes.Buffer
			SelectQry.WriteString(SQLQrySel.String())
			SelectQry.WriteString(SQLQryCol.String())
			SelectQry.WriteString(SQLQryFrom.String())
			SelectQry.WriteString(SQLCondition.String())
			SelectQry.WriteString(SQLGroupOrder.String())
			SelectQry.WriteString(SQLLimit.String())

			var SelectFilterQry bytes.Buffer
			SelectFilterQry.WriteString(SQLQrySel.String())
			SelectFilterQry.WriteString(SQLQryCnt.String())
			SelectFilterQry.WriteString(SQLQryFrom.String())
			SelectFilterQry.WriteString(SQLCondition.String())

			var TotalCnt bytes.Buffer
			TotalCnt.WriteString(SQLQrySel.String())
			TotalCnt.WriteString(SQLQryCnt.String())
			TotalCnt.WriteString(SQLQryFrom.String())

			// rateRestUpdateData, err := model.ExecuteQuery(SelectQry.String(), reqMap["length"], reqMap["start"])
			rateRestUpdateData, err := model.ExecuteQuery(SelectQry.String(), fmt.Sprintf("%.0f", reqMap["length"]), fmt.Sprintf("%.0f", reqMap["start"]))
			if util.CheckErrorLog(r, err) {
				return nil, err
			}
			if len(rateRestUpdateData) == 0 {
				stuff["data"] = results
			} else {
				stuff["data"] = rateRestUpdateData
			}

			// log.Println("SelectFilterQry :: ", SelectFilterQry.String())
			FilterTotalCnt, err := model.ExecuteRowQuery(SelectFilterQry.String())
			if util.CheckErrorLog(r, err) {
				return nil, err
			}
			stuff["recordsFiltered"] = FilterTotalCnt["count"]

			// log.Println("TotalCnt :: ", TotalCnt.String())
			TotalCount, err := model.ExecuteRowQuery(TotalCnt.String())
			if util.CheckErrorLog(r, err) {
				return nil, err
			}
			stuff["recordsTotal"] = TotalCount["count"]
		}
	}
	return stuff, nil
}

// CheckIfHotelExistsOptional - Checks If Hotel Exists With Optional Active, Live Status Check - 2021-04-27 - HK
func CheckIfHotelExistsOptional(r *http.Request, HotelID string, statusChk bool, liveChk bool) (int64, error) {
	util.LogIt(r, "model - V_Partner_Hotel - CheckIfHotelExistsOptional")
	var Qry bytes.Buffer

	Qry.WriteString("SELECT COUNT(id) AS id FROM cf_hotel_info WHERE id = ?")
	if statusChk {
		Qry.WriteString(" AND status = 1 ")
	}
	if liveChk {
		Qry.WriteString(" AND is_live = 1 ")
	}
	Data, err := model.ExecuteRowQuery(Qry.String(), HotelID)
	if util.CheckErrorLog(r, err) {
		return -1, err
	}

	return Data["id"].(int64), nil
}

// GetReviewOfHotel - Get Review Flag for Hotel
func GetReviewOfHotel(r *http.Request, hotelID string) (map[string]interface{}, error) {
	util.LogIt(r, "Model - V_Partner_Hotel - GetReviewOfHotel")
	var Qry bytes.Buffer

	Qry.WriteString("SELECT is_review FROM cf_hotel_info WHERE id = ?")
	Data, err := model.ExecuteRowQuery(Qry.String(), hotelID)
	if util.CheckErrorLog(r, err) {
		return nil, err
	}

	return Data, nil
}

// VerifyHotel - Verify Hotel data submit to Admin
func VerifyHotel(r *http.Request, hotelID string) (bool, string, error, string) {
	util.LogIt(r, "Model - V_Partner_Hotel - VerifyHotel")

	// Check whether the Hotel Data is completely filled or not
	hotelVerificationFlg, errCodeString := model.CheckHotelVerificationEligibility(hotelID)
	if !hotelVerificationFlg {
		return false, "verification", nil, errCodeString
	}

	// Send mail for Admin and Partner
	if hotelVerificationFlg {

		// Hotel Data
		var Qry bytes.Buffer
		Qry.WriteString("SELECT hotel_name FROM cf_hotel_info WHERE id = ?")
		hotelData, err := model.ExecuteRowQuery(Qry.String(), hotelID)
		if util.CheckErrorLog(r, err) {
			return false, "hotel_data", err, ""
		}

		resMap := make(map[string]interface{})
		resMap["hotel_name"] = hotelData["hotel_name"].(string)

		HotelClientID := context.Get(r, "UserId")
		var Qry2 bytes.Buffer
		Qry2.WriteString("SELECT email, client_name FROM cf_hotel_client WHERE id = ?")
		hotelClientData, err := model.ExecuteRowQuery(Qry2.String(), HotelClientID)
		if util.CheckErrorLog(r, err) {
			return false, "hotel_client_data", err, ""
		}

		resMap["email_id"] = hotelClientData["email"].(string)
		resMap["client_name"] = hotelClientData["client_name"].(string)

		// Admin Data for getting admin email
		var AdmQry bytes.Buffer
		AdmQry.WriteString(" SELECT email FROM cf_user_profile WHERE user_id = 1")
		AdminData, err := model.ExecuteRowQuery(AdmQry.String())
		if util.CheckErrorLog(r, err) {
			return false, "admin_data", err, ""
		}

		// Admin mail template
		model.MailChn <- model.MailObj{
			Type:         "EmailTemplate",
			ID:           "11",
			Additional:   AdminData["email"].(string),
			InterfaceObj: resMap,
		}

		// Partner mail template
		model.MailChn <- model.MailObj{
			Type:         "EmailTemplate",
			ID:           "12",
			Additional:   hotelClientData["email"].(string),
			InterfaceObj: resMap,
		}

		var UpdQry bytes.Buffer
		UpdQry.WriteString("UPDATE cf_hotel_info SET is_review = 1 WHERE id = ?")
		err = model.ExecuteNonQuery(UpdQry.String(), hotelID)
		if util.CheckErrorLog(r, err) {
			return true, "mail_delivered", nil, ""
		}
	}

	return true, "mail_delivered", nil, ""
}
