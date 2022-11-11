package model

import (
	"bytes"
	"encoding/json"
	"net/http"
	"time"
	"tp-api-common/util"
	"tp-system/config"
)

// GetRoomTypeList - Get Room Type List For Other Module
func GetRoomTypeList(r *http.Request, id string) (map[string]interface{}, error) {
	util.LogIt(r, "model - Model - GetRoomTypeList")
	var Qry bytes.Buffer

	Qry.WriteString(" SELECT CFT.id,CFT.room_type_name,CBT.bed_type,CFT.room_size, CFT.description AS room_desc, ")
	Qry.WriteString(" CASE WHEN is_extra_bed = 1 THEN 'YES' ELSE 'NO' END AS extra_bed_type, ")
	Qry.WriteString(" IFNULL(CEBT.extra_bed_name,' ') AS extra_bed_name, max_occupancy,inventory,CRV.room_view_name ")
	Qry.WriteString(" FROM cf_room_type AS CFT ")
	Qry.WriteString(" INNER JOIN cf_bed_type AS CBT ON CBT.id = CFT.bed_type_id ")
	Qry.WriteString(" INNER JOIN cf_room_view AS CRV ON CRV.id = CFT.room_view_id ")
	Qry.WriteString(" LEFT JOIN cf_extra_bed_type AS CEBT ON CEBT.id = CFT.extra_bed_type_id ")
	Qry.WriteString(" WHERE CFT.status = 1 AND CFT.hotel_id = ?; ")
	RetMap, err := ExecuteQuery(Qry.String(), id)
	if util.CheckErrorLog(r, err) {
		return nil, err
	}

	var stuff = make(map[string]interface{})
	if len(RetMap) > 0 {
		for i, val := range RetMap {
			var Qry bytes.Buffer

			Qry.WriteString("SELECT id, CONCAT('" + config.Env.AwsBucketURL + "room/" + "',image) AS image, sort_order FROM cf_room_image WHERE room_type_id = ? AND hotel_id = ? ORDER BY sort_order")
			ImgMap, err := ExecuteQuery(Qry.String(), val["id"], id)
			if util.CheckErrorLog(r, err) {
				return nil, err
			}
			RetMap[i]["image"] = ImgMap
		}
		stuff["data"] = RetMap
	} else {
		stuff["data"] = []string{}
	}

	return stuff, nil
}

// GetRatePlanFromRoom - Return rateplan including cancellation policy and other stuff by passing room id
func GetRatePlanFromRoom(r *http.Request, id, roomID string) (map[string]interface{}, error) {
	util.LogIt(r, "model - Model - GetRatePlanFromRoom")
	var Qry bytes.Buffer

	Qry.WriteString(" SELECT CFR.id, rate_plan_name,CASE WHEN is_pay_at_hotel = 1 THEN 'YES' ELSE 'NO' END AS pay_at_hotel, ")
	Qry.WriteString(" CMT.meal_type, rate, sort_order, CFR.cancellation_policy_id ")
	Qry.WriteString(" FROM cf_rateplan AS CFR ")
	Qry.WriteString(" INNER JOIN cf_meal_type AS CMT ON CMT.id = CFR.meal_type_id ")
	Qry.WriteString(" WHERE CFR.status = 1 AND CFR.hotel_id = ? AND CFR.room_type_id = ?; ")
	RetMap, err := ExecuteQuery(Qry.String(), id, roomID)
	if util.CheckErrorLog(r, err) {
		return nil, err
	}

	var stuff = make(map[string]interface{})
	if len(RetMap) > 0 {
		for i, val := range RetMap {
			var Qry, Qry1 bytes.Buffer

			Qry.WriteString(" SELECT policy_name,CASE WHEN is_non_refundable = 1 THEN 'YES' ELSE 'NO' END AS is_non_refundable,before_day,before_day_charge,after_day_charge FROM cf_cancellation_policy WHERE id = ? AND hotel_id = ? ")
			ImgMap, err := ExecuteRowQuery(Qry.String(), val["cancellation_policy_id"], id)
			if util.CheckErrorLog(r, err) {
				return nil, err
			}
			RetMap[i]["cancellation_policy"] = ImgMap

			Qry1.WriteString("SELECT GROUP_CONCAT(inclusion) AS inclusion FROM cf_inclusion WHERE id IN (SELECT GROUP_CONCAT(distinct inclusion_id SEPARATOR ',') as inclusion FROM cf_rateplan_inclusion WHERE rateplan_id = ? AND hotel_id = ?);")
			RetMap1, err1 := ExecuteRowQuery(Qry1.String(), val["id"], id)
			if util.CheckErrorLog(r, err1) {
				return nil, err1
			}

			RetMap[i]["inclusion"] = RetMap1["inclusion"]
		}
		stuff["data"] = RetMap
	} else {
		stuff["data"] = []string{}
	}

	return stuff, nil
}

// GetHotelRoomRateData - Gets Inv, Rate, Restrictions Data Occupancy Wise For hotel
func GetHotelRoomRateData(r *http.Request, reqMap map[string]interface{}) ([]map[string]interface{}, error) {
	util.LogIt(r, "Model - Model - GetHotelRoomRateData")

	ReqMonth := reqMap["month"]
	ReqYear := reqMap["year"]
	roomID := reqMap["room_id"]
	rateID := reqMap["rate_id"]
	HotelID := reqMap["hotel_id"]

	if ReqMonth == "" || ReqYear == "" {
		currentTime := time.Now()
		ReqYear = currentTime.Year()
		ReqMonth = currentTime.Month()
	}

	var FinalArray []map[string]interface{}

	if roomID != "" && rateID != "" {
		var RQry bytes.Buffer
		// Get Room Info
		RQry.WriteString("SELECT id, room_type_name, bed_type_id, room_size, description, room_view_id, is_extra_bed, extra_bed_type_id, max_occupancy, inventory, sort_order, created_at, created_by, hotel_id FROM cf_room_type WHERE id = ? AND hotel_id = ?")
		roomInfo, err := ExecuteRowQuery(RQry.String(), roomID, HotelID)
		if util.CheckErrorLog(r, err) {
			return nil, err
		}

		RoomID := roomInfo["id"].(string)
		RoomName := roomInfo["room_type_name"].(string)
		RoomOccupancy := roomInfo["max_occupancy"].(int64)

		var Qry bytes.Buffer
		Qry.WriteString("SELECT inv_data as data FROM cf_inv_data WHERE hotel_id = ? AND room_id = ? AND year = ? AND month = ? ") // Add
		Data, err := ExecuteRowQuery(Qry.String(), HotelID.(string), RoomID, ReqYear, ReqMonth)
		if util.CheckErrorLog(r, err) {
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

		var Rate, Rate1 bytes.Buffer

		Rate.WriteString("SELECT id, rate_plan_name, room_type_id, is_pay_at_hotel, cancellation_policy_id, meal_type_id, rate, sort_order, status, created_at, created_by, hotel_id FROM cf_rateplan WHERE status = 1 AND id = ? AND hotel_id = ?")
		rateInfo, err := ExecuteRowQuery(Rate.String(), rateID, HotelID)
		if util.CheckErrorLog(r, err) {
			return nil, err
		}

		Rate1.WriteString("SELECT GROUP_CONCAT(distinct inclusion_id SEPARATOR ',') as inclusion FROM cf_rateplan_inclusion WHERE rateplan_id = ? AND hotel_id = ?;")
		RetMap1, err1 := ExecuteRowQuery(Rate1.String(), rateID, HotelID)
		if util.CheckErrorLog(r, err1) {
			return nil, err1
		}

		rateInfo["inclusion"] = RetMap1["inclusion"]
		RateID := rateInfo["id"].(string)
		RateName := rateInfo["rate_plan_name"].(string)

		var rateArr []map[string]interface{}

		var Qry1 bytes.Buffer
		Qry1.WriteString("SELECT rate_rest_data as data FROM cf_rate_restriction_data WHERE hotel_id = ? AND room_id = ? AND rateplan_id = ? AND year = ? AND month = ? ") // Add
		Data1, err := ExecuteRowQuery(Qry1.String(), HotelID.(string), RoomID, RateID, ReqYear, ReqMonth)
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
	}
	return FinalArray, nil
}
