package model

import (
	"bytes"
	"net/http"
	"tp-api-common/data"
	"tp-api-common/util"

	"github.com/gorilla/context"
	gonanoid "github.com/matoous/go-nanoid"
)

// AddExtraBedType - Add Extra Bed Type
func AddExtraBedType(r *http.Request, reqMap data.ExtraBedType) bool {
	util.LogIt(r, "model - V_Extra_Bed_Type - AddExtraBedType")
	var Qry bytes.Buffer
	nanoid, _ := gonanoid.Nanoid()
	Qry.WriteString("INSERT INTO cf_extra_bed_type(id,extra_bed_name,created_at,created_by) VALUES (?,?,?,?)")
	err := ExecuteNonQuery(Qry.String(), nanoid, reqMap.Type, util.GetIsoLocalDateTime(), context.Get(r, "UserId"))
	if util.CheckErrorLog(r, err) {
		return false
	}

	AddLog(r, "", "EXTRA_BED_TYPE", "Create", nanoid, GetLogsValueMap(r, util.ToMap(reqMap), false, ""))

	return true
}

// UpdateExtraBedType - Update Extra Bed Type
func UpdateExtraBedType(r *http.Request, reqMap data.ExtraBedType) bool {
	util.LogIt(r, "model - V_Extra_Bed_Type - UpdateExtraBedType")
	var Qry bytes.Buffer

	BeforeUpdate, _ := GetModuleFieldByID(r, "EXTRA_BED_TYPE", reqMap.ID, "extra_bed_name")

	Qry.WriteString("UPDATE cf_extra_bed_type SET extra_bed_name = ? WHERE id = ?")
	err := ExecuteNonQuery(Qry.String(), reqMap.Type, reqMap.ID)
	if util.CheckErrorLog(r, err) {
		return false
	}

	AddLog(r, BeforeUpdate.(string), "EXTRA_BED_TYPE", "Update", reqMap.ID, GetLogsValueMap(r, util.ToMap(reqMap), false, ""))

	return true
}

// ExtraBedTypeListing - Datatable Extra Bed Type listing with filter and order
func ExtraBedTypeListing(r *http.Request, reqMap data.JQueryTableUI) (map[string]interface{}, error) {
	util.LogIt(r, "model - V_Extra_Bed_Type - ExtraBedTypeListing")
	var Qry, QryCnt, QryFilter, GroupBy, FromQry bytes.Buffer
	var status = make(map[string]interface{})

	var testColArrs [20]string
	testColArrs[0] = "CMT.id"
	testColArrs[1] = "extra_bed_name"
	testColArrs[2] = "CMT.status"
	testColArrs[3] = "created_by"

	var testArrs []map[string]string
	testArrs = append(testArrs, map[string]string{
		"key":   "extra_bed_name",
		"value": "CMT.extra_bed_name",
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

	Qry.WriteString(" CMT.id,extra_bed_name,CONCAT(from_unixtime(CMT.created_at),' ',SUC.username) AS created_by,ST.status,ST.id AS status_id ")

	FromQry.WriteString(" FROM cf_extra_bed_type AS CMT ")
	FromQry.WriteString(" INNER JOIN cf_user AS SUC ON SUC.id = CMT.created_by ")
	FromQry.WriteString(" INNER JOIN status AS ST ON ST.id = CMT.status ")
	FromQry.WriteString(" WHERE CMT.status <> 3 ")
	Data, err := JQueryTable(r, reqMap, Qry, FromQry, QryCnt, QryFilter, GroupBy, testArrs, testColArrs, status, "")
	return Data, err
}

// GetExtraBedTypeList - Get Extra Bed Type List For Other Module
func GetExtraBedTypeList(r *http.Request) ([]map[string]interface{}, error) {
	util.LogIt(r, "model - V_Extra_Bed_Type - GetExtraBedTypeList")
	var Qry bytes.Buffer

	Qry.WriteString("SELECT id,extra_bed_name FROM cf_extra_bed_type WHERE status = 1")
	RetMap, err := ExecuteQuery(Qry.String())
	if util.CheckErrorLog(r, err) {
		return nil, err
	}

	return RetMap, nil
}
