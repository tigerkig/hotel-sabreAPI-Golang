package partner

import (
	"bytes"
	"net/http"
	"tp-api-common/data"
	"tp-api-common/util"
	"tp-system/model"

	"github.com/gorilla/context"
	gonanoid "github.com/matoous/go-nanoid"
)

// AddCancelPolicy - Add Cancel Policy For Hotel
func AddCancelPolicy(r *http.Request, reqMap data.CancelPolicy) bool {
	util.LogIt(r, "model - V_Cancel_Policy - AddCancelPolicy")

	var Qry bytes.Buffer
	nanoid, _ := gonanoid.Nanoid()

	Qry.WriteString("INSERT INTO cf_cancellation_policy(id, policy_name, is_non_refundable, before_day, before_day_charge, after_day_charge, created_at, created_by, hotel_id) VALUES (?,?,?,?,?,?,?,?,?)")
	err := model.ExecuteNonQuery(Qry.String(), nanoid, reqMap.Policy, reqMap.IsNonRefundable, reqMap.BeforeDay, reqMap.BeforeDayCharge, reqMap.AfterDayCharge, util.GetIsoLocalDateTime(), context.Get(r, "UserId"), reqMap.HotelID)
	if util.CheckErrorLog(r, err) {
		return false
	}

	model.AddLog(r, "", "CANCEL_POLICY", "Create", nanoid, model.GetLogsValueMap(r, util.ToMap(reqMap), false, ""))

	return true
}

// UpdateCancelPolicy - Update Amenity
func UpdateCancelPolicy(r *http.Request, reqMap data.CancelPolicy) bool {
	util.LogIt(r, "model - V_Cancel_Policy - UpdateCancelPolicy")

	var Qry bytes.Buffer

	Qry.WriteString("UPDATE cf_cancellation_policy SET policy_name=?,is_non_refundable=?,before_day=?, before_day_charge=?,after_day_charge=? WHERE id=? AND hotel_id=?")
	err := model.ExecuteNonQuery(Qry.String(), reqMap.Policy, reqMap.IsNonRefundable, reqMap.BeforeDay, reqMap.BeforeDayCharge, reqMap.AfterDayCharge, reqMap.ID, reqMap.HotelID)
	if util.CheckErrorLog(r, err) {
		return false
	}

	reqStruct := util.ToMap(reqMap)
	model.AddLog(r, "", "CANCEL_POLICY", "UPDATE", reqMap.ID, model.GetLogsValueMap(r, reqStruct, false, ""))

	//Update cancellation policy in all deals & rateplan details in sync DB
	GetRatePlanInfoFromCancelPolicy(r, reqMap.ID)

	return true
}

// CancelPolicyListing - Return Datatable Listing Of Cancel Policy
func CancelPolicyListing(r *http.Request, reqMap data.JQueryTableUI) (map[string]interface{}, error) {
	util.LogIt(r, "model - V_Amenity - CancelPolicyListing")

	var Qry, QryCnt, QryFilter, GroupBy, FromQry bytes.Buffer
	var status = make(map[string]interface{})

	var testColArrs [20]string
	testColArrs[0] = "CCP.id"
	testColArrs[1] = "CCP.policy_name"
	testColArrs[2] = "CCP.is_non_refundable"
	testColArrs[3] = "CCP.before_day"
	testColArrs[4] = "CCP.before_day_charge"
	testColArrs[5] = "CCP.after_day_charge"
	testColArrs[6] = "CFA.status"
	testColArrs[7] = "created_by"

	var testArrs []map[string]string
	testArrs = append(testArrs, map[string]string{
		"key":   "policy_name",
		"value": "CCP.policy_name",
	})
	testArrs = append(testArrs, map[string]string{
		"key":   "is_non_refundable",
		"value": "CCP.is_non_refundable",
	})
	testArrs = append(testArrs, map[string]string{
		"key":   "before_day",
		"value": "CCP.before_day",
	})
	testArrs = append(testArrs, map[string]string{
		"key":   "before_day_charge",
		"value": "CCP.before_day_charge",
	})
	testArrs = append(testArrs, map[string]string{
		"key":   "after_day_charge",
		"value": "CCP.after_day_charge",
	})
	testArrs = append(testArrs, map[string]string{
		"key":   "status",
		"value": "CCP.status",
	})
	testArrs = append(testArrs, map[string]string{
		"key":   "created_at",
		"value": "DATE(from_unixtime(CCP.created_at))",
	})

	QryCnt.WriteString(" COUNT(CCP.id) AS cnt ")
	QryFilter.WriteString(" COUNT(CCP.id) AS cnt ")

	Qry.WriteString(" CCP.id, CCP.policy_name, CCP.is_non_refundable, CCP.before_day, ST.status, CONCAT(from_unixtime(CCP.created_at),' ',CHC.username) AS created_by,ST.id AS status_id ")

	FromQry.WriteString(" FROM cf_cancellation_policy AS CCP ")
	FromQry.WriteString(" INNER JOIN cf_hotel_client AS CHC ON CHC.id = CCP.created_by ")
	FromQry.WriteString(" INNER JOIN status AS ST ON ST.id = CCP.status ")
	FromQry.WriteString(" WHERE ST.id <> 3 AND CCP.hotel_id = '")
	//FromQry.WriteString(context.Get(r, "HotelId").(string))
	FromQry.WriteString(reqMap.HotelID)
	FromQry.WriteString("'")
	Data, err := model.JQueryTable(r, reqMap, Qry, FromQry, QryCnt, QryFilter, GroupBy, testArrs, testColArrs, status, "")
	return Data, err
}

// GetCancelPolicyList - Get Cancel Policy List For Other Module
func GetCancelPolicyList(r *http.Request, hotelID string) ([]map[string]interface{}, error) {
	util.LogIt(r, "Model - V_Cancel_Policy - GetCancelPolicyList")
	var Qry bytes.Buffer
	Qry.WriteString("SELECT id, policy_name FROM cf_cancellation_policy WHERE status = 1 AND hotel_id = ?")
	RetMap, err := model.ExecuteQuery(Qry.String(), hotelID)
	if util.CheckErrorLog(r, err) {
		return nil, err
	}

	return RetMap, nil
}

// GetCancelPolicyInfo  - Get Room Type Info
func GetCancelPolicyInfo(r *http.Request, canID, HotelID string) (map[string]interface{}, error) {
	util.LogIt(r, "Model - V_Cancel_Policy - GetRoomType")

	var Qry bytes.Buffer

	Qry.WriteString("SELECT id, policy_name, is_non_refundable, before_day, before_day_charge, after_day_charge FROM cf_cancellation_policy WHERE id = ? AND hotel_id = ?")
	RetMap, err := model.ExecuteRowQuery(Qry.String(), canID, HotelID)
	if util.CheckErrorLog(r, err) {
		return nil, err
	}

	return RetMap, nil
}

// GetRatePlanInfoFromCancelPolicy  - Get rate plan info from cancellation policy id and update cancellation policy info in sync DB MS 2020-08-04
func GetRatePlanInfoFromCancelPolicy(r *http.Request, canPolicyID string) bool {
	util.LogIt(r, "Model - V_Cancel_Policy - GetRoomType")

	var Qry bytes.Buffer

	Qry.WriteString("SELECT id,cancellation_policy_id,hotel_id FROM cf_rateplan WHERE cancellation_policy_id = ? AND status = 1;")
	RetMap, err := model.ExecuteQuery(Qry.String(), canPolicyID)
	if util.CheckErrorLog(r, err) {
		return false
	}

	if len(RetMap) > 0 {
		for _, val := range RetMap {
			model.CacheChn <- model.CacheObj{
				Type:       "ratePlanDetails",
				ID:         val["hotel_id"].(string),
				Additional: val["id"].(string),
			}
		}
	}

	return true
}
