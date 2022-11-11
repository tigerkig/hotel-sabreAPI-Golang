package partner

import (
	"bytes"
	"log"
	"net/http"
	"strings"
	"tp-api-common/data"
	"tp-api-common/util"
	"tp-system/model"

	"github.com/gorilla/context"
	gonanoid "github.com/matoous/go-nanoid"
)

func AddInv(r *http.Request, hotelID string) bool {

	util.LogIt(r, "Model - V_Partner_HotelInv - AddInv")

	var Qry bytes.Buffer
	Qry.WriteString("SELECT * FROM cf_rateplan order by hotel_id desc")
	log.Println(Qry.String())
	RetMap, err := model.ExecuteQuery(Qry.String())
	if util.CheckErrorLog(r, err) {
		return false
	}

	for _, val := range RetMap {
		RoutineInsertData(r, val)
	}

	return true

}

func RoutineInsertData(r *http.Request, val map[string]interface{}) bool {
	log.Println(val["id"].(string))
	rateRestFilled, _ := model.CheckAndFillRateRestDataOnActiveStatus(r, val["hotel_id"].(string), val["id"].(string))
	if !rateRestFilled {
		util.LogIt(r, "Issue on sync data on add rate plan - "+val["id"].(string))
	}
	// 2020-08-20 - HK - RatePlan Details Sync On Mongo
	invFilled := model.CheckAndFillInvdataOnActiveStatus(r, val["hotel_id"].(string), val["room_type_id"].(string))
	return !invFilled
}

// AddRatePlan - Add Rate Plan For Room
func AddRatePlan(r *http.Request, reqMap data.RatePlan) bool {

	util.LogIt(r, "Model - V_Partner_RatePlan - AddRatePlan")

	var Qry bytes.Buffer
	ratePlanID, _ := gonanoid.Nanoid()

	Qry.WriteString("INSERT INTO cf_rateplan(id, rate_plan_name, room_type_id, is_pay_at_hotel, cancellation_policy_id, meal_type_id, rate, sort_order, created_at, created_by, hotel_id) VALUES (?,?,?,?,?,?,?,?,?,?,?)")
	// err := model.ExecuteNonQuery(Qry.String(), ratePlanID, reqMap.Name, reqMap.RoomType, reqMap.IsPayAtHotel, reqMap.CancelPolicy, reqMap.MealPlan, reqMap.Rate, reqMap.SortOrder, util.GetIsoLocalDateTime(), context.Get(r, "UserId"), context.Get(r, "HotelId"))
	err := model.ExecuteNonQuery(Qry.String(), ratePlanID, reqMap.Name, reqMap.RoomType, reqMap.IsPayAtHotel, reqMap.CancelPolicy, reqMap.MealPlan, reqMap.Rate, reqMap.SortOrder, util.GetIsoLocalDateTime(), context.Get(r, "UserId"), reqMap.HotelID)
	if util.CheckErrorLog(r, err) {
		return false
	}

	model.AddLog(r, "", "RATE_PLAN", "Create", ratePlanID, model.GetLogsValueMap(r, util.ToMap(reqMap), false, ""))

	if reqMap.Inclusion != "" {
		var InclusionArr = strings.Split(reqMap.Inclusion, ",")
		if len(InclusionArr) > 0 {
			for _, val := range InclusionArr {
				var TagQry bytes.Buffer
				nanoid, _ := gonanoid.Nanoid()
				TagQry.WriteString("INSERT INTO cf_rateplan_inclusion SET id = ?, rateplan_id = ?, inclusion_id = ?, hotel_id = ?")
				// err = model.ExecuteNonQuery(TagQry.String(), nanoid, ratePlanID, val, context.Get(r, "HotelId"))
				err = model.ExecuteNonQuery(TagQry.String(), nanoid, ratePlanID, val, reqMap.HotelID)
				if util.CheckErrorLog(r, err) {
					return false
				}
			}
		}
	}

	// 2020-08-20 - HK - RatePlan Details Sync On Mongo
	model.CacheChn <- model.CacheObj{
		Type: "ratePlanDetails",
		// ID:         context.Get(r, "HotelId").(string),
		ID:         reqMap.HotelID,
		Additional: ratePlanID,
	}

	// rateRestFilled, _ := model.CheckAndFillRateRestDataOnActiveStatus(r, context.Get(r, "HotelId").(string), ratePlanID)
	rateRestFilled, _ := model.CheckAndFillRateRestDataOnActiveStatus(r, reqMap.HotelID, ratePlanID)
	if !rateRestFilled {
		util.LogIt(r, "Issue on sync data on add rate plan - "+ratePlanID)
	}
	// 2020-08-20 - HK - RatePlan Details Sync On Mongo

	// invFilled := model.CheckAndFillInvdataOnActiveStatus(r, reqMap.HotelID, reqMap.RoomType)
	// if !invFilled {
	// 	return false
	// }

	return true
}

// UpdateRatePlan - Updates Rate Plan
func UpdateRatePlan(r *http.Request, reqMap data.RatePlan) bool {

	util.LogIt(r, "Model - V_Partner_RatePlan - UpdateRatePlan")

	var Qry bytes.Buffer

	Qry.WriteString("UPDATE cf_rateplan SET rate_plan_name = ?,  room_type_id = ?,  is_pay_at_hotel = ?,  cancellation_policy_id = ?, meal_type_id = ?, rate = ?, sort_order = ? WHERE id = ? AND hotel_id = ?")
	// err := model.ExecuteNonQuery(Qry.String(), reqMap.Name, reqMap.RoomType, reqMap.IsPayAtHotel, reqMap.CancelPolicy, reqMap.MealPlan, reqMap.Rate, reqMap.SortOrder, reqMap.ID, context.Get(r, "HotelId"))
	err := model.ExecuteNonQuery(Qry.String(), reqMap.Name, reqMap.RoomType, reqMap.IsPayAtHotel, reqMap.CancelPolicy, reqMap.MealPlan, reqMap.Rate, reqMap.SortOrder, reqMap.ID, reqMap.HotelID)
	if util.CheckErrorLog(r, err) {
		return false
	}

	model.AddLog(r, "", "RATE_PLAN", "UPDATE", reqMap.ID, model.GetLogsValueMap(r, util.ToMap(reqMap), false, ""))

	if reqMap.Inclusion != "" {

		// var InclusionStringArr []string
		var InclusionArr = strings.Split(reqMap.Inclusion, ",")

		if len(InclusionArr) > 0 {

			var DelQry bytes.Buffer
			DelQry.WriteString("DELETE FROM cf_rateplan_inclusion WHERE rateplan_id = ? AND hotel_id = ?")
			// err = model.ExecuteNonQuery(DelQry.String(), reqMap.ID, context.Get(r, "HotelId"))
			err = model.ExecuteNonQuery(DelQry.String(), reqMap.ID, reqMap.HotelID)
			if util.CheckErrorLog(r, err) {
				return false
			}

			for _, val := range InclusionArr {
				var TagQry bytes.Buffer
				nanoid, _ := gonanoid.Nanoid()
				TagQry.WriteString("INSERT INTO cf_rateplan_inclusion SET id = ?, rateplan_id = ?, inclusion_id = ?, hotel_id = ?")
				// err = model.ExecuteNonQuery(TagQry.String(), nanoid, reqMap.ID, val, context.Get(r, "HotelId"))
				err = model.ExecuteNonQuery(TagQry.String(), nanoid, reqMap.ID, val, reqMap.HotelID)
				if util.CheckErrorLog(r, err) {
					return false
				}

				// BeforeUpdate, _ := model.GetModuleFieldByID(r, "PROPERTY_TAG", val, "tag")
				// InclusionStringArr = append(InclusionStringArr, BeforeUpdate.(string))
			}
		}
	}

	// 2020-08-20 - HK - RatePlan Details Sync On Mongo
	model.CacheChn <- model.CacheObj{
		Type: "ratePlanDetails",
		// ID:         context.Get(r, "HotelId").(string),
		ID:         reqMap.HotelID,
		Additional: reqMap.ID,
	}
	// 2020-08-20 - HK - RatePlan Details Sync On Mongo

	return true
}

// GetRatePlan - Get Rate Plan Info
func GetRatePlan(r *http.Request, rateID string, HotelID string) (map[string]interface{}, error) {

	util.LogIt(r, "Model - V_Partner_RatePlan - GetRatePlan")

	var Qry, Qry1 bytes.Buffer
	//Changed made by Meet Soni(Remove status field from where caluse)
	//Qry.WriteString("SELECT id, rate_plan_name, room_type_id, is_pay_at_hotel, cancellation_policy_id, meal_type_id, rate, sort_order, status, created_at, created_by, hotel_id FROM cf_rateplan WHERE status = 1 AND id = ? AND hotel_id = ?")
	Qry.WriteString("SELECT id, rate_plan_name, room_type_id, is_pay_at_hotel, cancellation_policy_id, meal_type_id, rate, sort_order, status, created_at, created_by, hotel_id FROM cf_rateplan WHERE id = ? AND hotel_id = ?")
	RetMap, err := model.ExecuteRowQuery(Qry.String(), rateID, HotelID)
	if util.CheckErrorLog(r, err) {
		return nil, err
	}

	Qry1.WriteString("SELECT GROUP_CONCAT(distinct inclusion_id SEPARATOR ',') as inclusion FROM cf_rateplan_inclusion WHERE rateplan_id = ? AND hotel_id = ?;")
	RetMap1, err1 := model.ExecuteRowQuery(Qry1.String(), rateID, HotelID)
	if util.CheckErrorLog(r, err1) {
		return nil, err1
	}

	RetMap["inclusion"] = RetMap1["inclusion"]

	return RetMap, nil
}

// RatePlanListing - Return Datatable Listing Of Rate Plan
func RatePlanListing(r *http.Request, reqMap data.JQueryTableUI) (map[string]interface{}, error) {

	util.LogIt(r, "Model - V_Partner_RatePlan - RatePlanListing")

	var Qry, QryCnt, QryFilter, GroupBy, FromQry bytes.Buffer
	var status = make(map[string]interface{})

	var testColArrs [20]string
	testColArrs[0] = "CRP.id"
	testColArrs[1] = "CRP.rate_plan_name"
	testColArrs[2] = "ST.status"
	testColArrs[3] = "CRT.room_type_name"
	testColArrs[4] = "CRP.rate"
	testColArrs[5] = "CMT.meal_type"
	testColArrs[6] = "CRP.is_pay_at_hotel"
	testColArrs[7] = "CRP.sort_order"
	testColArrs[8] = "created_by"

	var testArrs []map[string]string
	testArrs = append(testArrs, map[string]string{
		"key":   "rate_plan_name",
		"value": "CRP.rate_plan_name",
	})
	testArrs = append(testArrs, map[string]string{
		"key":   "room_type_name",
		"value": "CRP.room_type_id",
	})
	testArrs = append(testArrs, map[string]string{
		"key":   "is_pay_at_hotel",
		"value": "CRP.is_pay_at_hotel",
	})
	testArrs = append(testArrs, map[string]string{
		"key":   "meal_type",
		"value": "CRP.meal_type_id",
	})
	testArrs = append(testArrs, map[string]string{
		"key":   "status",
		"value": "CRP.status",
	})
	testArrs = append(testArrs, map[string]string{
		"key":   "created_at",
		"value": "DATE(from_unixtime(CRP.created_at))",
	})

	QryCnt.WriteString(" COUNT(CRP.id) AS cnt ")

	QryFilter.WriteString(" select count(tbl.cnt) FROM (")
	QryFilter.WriteString(" COUNT(CRP.id) AS cnt ")

	Qry.WriteString(" CRP.id, CRP.rate_plan_name, CRP.is_pay_at_hotel, CRP.rate, ST.status, CONCAT(from_unixtime(CRP.created_at),' ',CHC.username) AS created_by, ST.id AS status_id, CRP.sort_order, ")
	Qry.WriteString(" CRP.room_type_id, CRT.room_type_name, ")
	Qry.WriteString(" CRP.cancellation_policy_id, CCP.policy_name, ")
	Qry.WriteString(" CRP.meal_type_id, CMT.meal_type ")

	FromQry.WriteString(" FROM cf_rateplan AS CRP ")
	FromQry.WriteString(" INNER JOIN cf_room_type AS CRT ON CRT.id = CRP.room_type_id ")
	FromQry.WriteString(" INNER JOIN cf_cancellation_policy AS CCP ON CCP.id = CRP.cancellation_policy_id ")
	FromQry.WriteString(" INNER JOIN cf_meal_type AS CMT ON CMT.id = CRP.meal_type_id ")

	FromQry.WriteString(" INNER JOIN cf_hotel_client AS CHC ON CHC.id = CRP.created_by ")
	FromQry.WriteString(" INNER JOIN status AS ST ON ST.id = CRP.status ")
	FromQry.WriteString(" WHERE ST.id <> 3 AND CRP.hotel_id = '")
	// FromQry.WriteString(context.Get(r, "HotelId").(string))
	FromQry.WriteString(reqMap.HotelID)
	FromQry.WriteString("'")

	Data, err := model.JQueryTable(r, reqMap, Qry, FromQry, QryCnt, QryFilter, GroupBy, testArrs, testColArrs, status, "")
	return Data, err
}

// GetRatePlanList - Get Rate Type List For Other Module
func GetRatePlanList(r *http.Request, RoomID string, HotelID string) ([]map[string]interface{}, error) {

	util.LogIt(r, "Model - V_Partner_RatePlan - GetRatePlanList")

	var Qry bytes.Buffer

	Qry.WriteString("SELECT id, rate_plan_name FROM cf_rateplan WHERE status = 1 AND hotel_id = ? AND room_type_id = ?")
	RetMap, err := model.ExecuteQuery(Qry.String(), HotelID, RoomID)
	if util.CheckErrorLog(r, err) {
		return nil, err
	}

	return RetMap, nil
}
