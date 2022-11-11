package model

import (
	"bytes"
	"net/http"
	"tp-api-common/data"
	"tp-api-common/util"

	"github.com/gorilla/context"
	gonanoid "github.com/matoous/go-nanoid"
)

//SpecialRequest - Module Name
var SpecialRequest = "SPECIAL_REQUEST"

// AddSpecialRequest - Add Special Request
func AddSpecialRequest(r *http.Request, reqMap data.SpecialRequest) bool {
	util.LogIt(r, "model - V_Special_Request - AddSpecialRequest")
	var Qry bytes.Buffer
	nanoid, _ := gonanoid.Nanoid()
	Qry.WriteString("INSERT INTO cf_special_request(id,special_request,created_at,created_by) VALUES (?,?,?,?)")
	err := ExecuteNonQuery(Qry.String(), nanoid, reqMap.Request, util.GetIsoLocalDateTime(), context.Get(r, "UserId"))
	if util.CheckErrorLog(r, err) {
		return false
	}

	AddLog(r, "", "SPECIAL_REQUEST", "Create", nanoid, GetLogsValueMap(r, util.ToMap(reqMap), false, ""))

	return true
}

// UpdateSpecialRequest - Update Special Request
func UpdateSpecialRequest(r *http.Request, reqMap data.SpecialRequest) bool {
	util.LogIt(r, "model - V_Special_Request - UpdateSpecialRequest")
	var Qry bytes.Buffer

	BeforeUpdate, _ := GetModuleFieldByID(r, SpecialRequest, reqMap.ID, "special_request")

	Qry.WriteString("UPDATE cf_special_request SET special_request = ? WHERE id = ?")
	err := ExecuteNonQuery(Qry.String(), reqMap.Request, reqMap.ID)
	if util.CheckErrorLog(r, err) {
		return false
	}

	AddLog(r, BeforeUpdate.(string), SpecialRequest, "Update", reqMap.ID, GetLogsValueMap(r, util.ToMap(reqMap), true, "ID"))

	return true
}

// SpecialRequestListing - Get Special Request Listing
func SpecialRequestListing(r *http.Request, reqMap data.JQueryTableUI) (map[string]interface{}, error) {
	util.LogIt(r, "model - V_Special_Request - SpecialRequestListing")
	var Qry, QryCnt, QryFilter, GroupBy, FromQry bytes.Buffer
	var status = make(map[string]interface{})

	var testColArrs [20]string
	testColArrs[0] = "CPR.id"
	testColArrs[1] = "CPR.special_request"
	testColArrs[2] = "CPR.status"
	testColArrs[3] = "created_by"

	var testArrs []map[string]string
	testArrs = append(testArrs, map[string]string{
		"key":   "special_request",
		"value": "CPR.special_request",
	})
	testArrs = append(testArrs, map[string]string{
		"key":   "status",
		"value": "CPR.status",
	})
	testArrs = append(testArrs, map[string]string{
		"key":   "created_at",
		"value": "DATE(from_unixtime(CPR.created_at))",
	})

	QryCnt.WriteString(" COUNT(CPR.id) AS cnt ")
	QryFilter.WriteString(" COUNT(CPR.id) AS cnt ")

	Qry.WriteString(" CPR.id, CPR.special_request, CONCAT(from_unixtime(CPR.created_at),' ',SUC.username) AS created_by, ST.status, ST.id AS status_id ")

	FromQry.WriteString(" FROM cf_special_request AS CPR ")
	FromQry.WriteString(" INNER JOIN cf_user AS SUC ON SUC.id = CPR.created_by  ")
	FromQry.WriteString(" INNER JOIN status AS ST ON ST.id = CPR.status ")
	FromQry.WriteString(" WHERE CPR.status <> 3 ")
	Data, err := JQueryTable(r, reqMap, Qry, FromQry, QryCnt, QryFilter, GroupBy, testArrs, testColArrs, status, "")
	return Data, err
}

// GetSpecialRequest - This function returns special request which are active. Function made for web.
func GetSpecialRequest(r *http.Request) ([]map[string]interface{}, bool) {
	util.LogIt(r, "model - V_Special_Request - GetSpecialRequest")
	var Qry bytes.Buffer

	Qry.WriteString(" SELECT id, special_request FROM cf_special_request WHERE status=1;")
	data, err := ExecuteQuery(Qry.String())
	if util.CheckErrorLog(r, err) {
		return nil, false
	}
	return data, true
}
