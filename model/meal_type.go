package model

import (
	"bytes"
	"net/http"
	"tp-api-common/data"
	"tp-api-common/util"

	"github.com/gorilla/context"
	gonanoid "github.com/matoous/go-nanoid"
)

//MealType - Module Name
var MealType = "MEAL_TYPE"

// AddMealType - Add Meal Type
func AddMealType(r *http.Request, reqMap data.MealType) bool {

	util.LogIt(r, "model - V_Meal_Type - AddMealType")

	var Qry bytes.Buffer
	nanoid, _ := gonanoid.Nanoid()

	Qry.WriteString("INSERT INTO cf_meal_type(id, meal_type, created_at, created_by) VALUES (?,?,?,?)")
	err := ExecuteNonQuery(Qry.String(), nanoid, reqMap.Type, util.GetIsoLocalDateTime(), context.Get(r, "UserId"))
	if util.CheckErrorLog(r, err) {
		return false
	}

	AddLog(r, "", MealType, "Create", nanoid, GetLogsValueMap(r, util.ToMap(reqMap), false, ""))

	return true
}

// UpdateMealType - Update Meal Type
func UpdateMealType(r *http.Request, reqMap data.MealType) bool {

	util.LogIt(r, "model - V_Meal_Type - UpdateMealType")

	var Qry bytes.Buffer

	BeforeUpdate, _ := GetModuleFieldByID(r, MealType, reqMap.ID, "meal_type")

	Qry.WriteString("UPDATE cf_meal_type SET meal_type = ? WHERE id = ?")
	err := ExecuteNonQuery(Qry.String(), reqMap.Type, reqMap.ID)
	if util.CheckErrorLog(r, err) {
		return false
	}

	AddLog(r, BeforeUpdate.(string), MealType, "Update", reqMap.ID, GetLogsValueMap(r, util.ToMap(reqMap), true, "ID"))

	return true
}

// MealTypeListing - Get Meal Type Listing
func MealTypeListing(r *http.Request, reqMap data.JQueryTableUI) (map[string]interface{}, error) {
	util.LogIt(r, "model - V_Meal_Type - MealTypeListing")

	var Qry, QryCnt, QryFilter, GroupBy, FromQry bytes.Buffer

	var status = make(map[string]interface{})

	var testColArrs [20]string
	testColArrs[0] = "CMT.id"
	testColArrs[1] = "CMT.meal_type"
	testColArrs[2] = "CMT.status"
	testColArrs[3] = "created_by"

	var testArrs []map[string]string
	testArrs = append(testArrs, map[string]string{
		"key":   "meal_type",
		"value": "CMT.meal_type",
	})
	testArrs = append(testArrs, map[string]string{
		"key":   "status",
		"value": "CMT.status",
	})
	testArrs = append(testArrs, map[string]string{
		"key":   "created_at",
		"value": "DATE(from_unixtime(CMT.created_at))",
	})

	QryCnt.WriteString(" COUNT(CMT.id) AS cnt ")
	QryFilter.WriteString(" COUNT(CMT.id) AS cnt ")

	Qry.WriteString(" CMT.id, CMT.meal_type, CONCAT(from_unixtime(CMT.created_at),' ',SUC.username) AS created_by, ST.status, ST.id AS status_id ")

	FromQry.WriteString(" FROM cf_meal_type AS CMT ")
	FromQry.WriteString(" INNER JOIN cf_user AS SUC ON SUC.id = CMT.created_by  ")
	FromQry.WriteString(" INNER JOIN status AS ST ON ST.id = CMT.status ")
	FromQry.WriteString(" WHERE CMT.status <> 3 ")
	Data, err := JQueryTable(r, reqMap, Qry, FromQry, QryCnt, QryFilter, GroupBy, testArrs, testColArrs, status, "")
	return Data, err
}

// GetMealTypeList - Get Meal Type List For Other Module
func GetMealTypeList(r *http.Request) ([]map[string]interface{}, error) {
	util.LogIt(r, "model - V_Meal_Type - GetMealTypeList")
	var Qry bytes.Buffer

	Qry.WriteString("SELECT id, meal_type FROM cf_meal_type WHERE status = 1")
	RetMap, err := ExecuteQuery(Qry.String())
	if util.CheckErrorLog(r, err) {
		return nil, err
	}

	return RetMap, nil
}
