package model

import (
	"bytes"
	"net/http"
	"tp-api-common/data"
	"tp-api-common/util"

	"github.com/gorilla/context"
	gonanoid "github.com/matoous/go-nanoid"
)

// AddBedType - Add Bed Type
func AddBedType(r *http.Request, reqMap data.BedTypes) bool {
	util.LogIt(r, "model - V_Bed_Type - AddBedType")
	var Qry bytes.Buffer
	nanoid, _ := gonanoid.Nanoid()
	Qry.WriteString("INSERT INTO cf_bed_type(id,bed_type,created_at,created_by) VALUES (?,?,?,?)")
	err := ExecuteNonQuery(Qry.String(), nanoid, reqMap.Type, util.GetIsoLocalDateTime(), context.Get(r, "UserId"))
	if util.CheckErrorLog(r, err) {
		return false
	}

	AddLog(r, "", "BED_TYPE", "Create", nanoid, GetLogsValueMap(r, util.ToMap(reqMap), false, ""))

	return true
}

// UpdateBedType - Update Bed Type
func UpdateBedType(r *http.Request, reqMap data.BedTypes) bool {
	util.LogIt(r, "model - V_Bed_Type - UpdateBedType")
	var Qry bytes.Buffer

	BeforeUpdate, _ := GetModuleFieldByID(r, "BED_TYPE", reqMap.ID, "bed_type")

	Qry.WriteString("UPDATE cf_bed_type SET bed_type = ? WHERE id = ?")
	err := ExecuteNonQuery(Qry.String(), reqMap.Type, reqMap.ID)
	if util.CheckErrorLog(r, err) {
		return false
	}

	AddLog(r, BeforeUpdate.(string), "BED_TYPE", "Update", reqMap.ID, GetLogsValueMap(r, util.ToMap(reqMap), false, ""))

	return true
}

// BedTypeListing - Datatable Bed Type listing with filter and order
func BedTypeListing(r *http.Request, reqMap data.JQueryTableUI) (map[string]interface{}, error) {
	util.LogIt(r, "model - V_Bed_Type - BedTypeListing")
	var Qry, QryCnt, QryFilter, GroupBy, FromQry bytes.Buffer
	var status = make(map[string]interface{})

	var testColArrs [20]string
	testColArrs[0] = "CMT.id"
	testColArrs[1] = "bed_type"
	testColArrs[2] = "CMT.status"
	testColArrs[3] = "created_by"

	var testArrs []map[string]string
	testArrs = append(testArrs, map[string]string{
		"key":   "bed_type",
		"value": "CMT.bed_type",
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

	Qry.WriteString(" CMT.id,bed_type,CONCAT(from_unixtime(CMT.created_at),' ',SUC.username) AS created_by,ST.status,ST.id AS status_id ")

	FromQry.WriteString(" FROM cf_bed_type AS CMT ")
	FromQry.WriteString(" INNER JOIN cf_user AS SUC ON SUC.id = CMT.created_by ")
	FromQry.WriteString(" INNER JOIN status AS ST ON ST.id = CMT.status ")
	FromQry.WriteString(" WHERE CMT.status <> 3 ")
	Data, err := JQueryTable(r, reqMap, Qry, FromQry, QryCnt, QryFilter, GroupBy, testArrs, testColArrs, status, "")
	return Data, err
}

// GetBedTypeList - Get Bed Type List For Other Module
func GetBedTypeList(r *http.Request) ([]map[string]interface{}, error) {
	util.LogIt(r, "model - V_Bed_Type - GetBedTypeList")
	var Qry bytes.Buffer

	Qry.WriteString("SELECT id,bed_type FROM cf_bed_type WHERE status = 1")
	RetMap, err := ExecuteQuery(Qry.String())
	if util.CheckErrorLog(r, err) {
		return nil, err
	}

	return RetMap, nil
}

// GetBedType -  Get Bed type Detail By ID - 2021-04-21 - HK
func GetBedType(r *http.Request, id string) (map[string]interface{}, error) {
	util.LogIt(r, "Model - V_Bed_Type - GetBedType")
	var Qry bytes.Buffer

	Qry.WriteString("SELECT id, bed_type FROM cf_bed_type WHERE id = ?")
	RetMap, err := ExecuteRowQuery(Qry.String(), id)
	if util.CheckErrorLog(r, err) {
		return nil, err
	}

	return RetMap, nil
}
