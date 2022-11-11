package model

import (
	"bytes"
	"net/http"
	"tp-api-common/data"
	"tp-api-common/util"

	"github.com/gorilla/context"
	gonanoid "github.com/matoous/go-nanoid"
)

// AddInclusion - Add Inclusion That's Free BreakFast etc
func AddInclusion(r *http.Request, reqMap data.Inclusion) bool {
	util.LogIt(r, "model - V_Inclusion - AddInclusion")
	var Qry bytes.Buffer
	nanoid, _ := gonanoid.Nanoid()
	Qry.WriteString("INSERT INTO cf_inclusion(id,inclusion,description,created_at,created_by) VALUES (?,?,?,?,?)")
	err := ExecuteNonQuery(Qry.String(), nanoid, reqMap.Inclusion, reqMap.Description, util.GetIsoLocalDateTime(), context.Get(r, "UserId"))
	if util.CheckErrorLog(r, err) {
		return false
	}

	AddLog(r, "", "INCLUSION", "Create", nanoid, GetLogsValueMap(r, util.ToMap(reqMap), false, ""))

	return true
}

// UpdateInclusion - Update Inclusion That's Free BreakFast etc
func UpdateInclusion(r *http.Request, reqMap data.Inclusion) bool {
	util.LogIt(r, "model - V_Inclusion - UpdateInclusion")
	var Qry bytes.Buffer

	BeforeUpdate, _ := GetModuleFieldByID(r, "INCLUSION", reqMap.ID, "inclusion")

	Qry.WriteString("UPDATE cf_inclusion SET inclusion = ?, description = ? WHERE id = ?")
	err := ExecuteNonQuery(Qry.String(), reqMap.Inclusion, reqMap.Description, reqMap.ID)
	if util.CheckErrorLog(r, err) {
		return false
	}

	AddLog(r, BeforeUpdate.(string), "INCLUSION", "Update", reqMap.ID, GetLogsValueMap(r, util.ToMap(reqMap), false, ""))

	return true
}

// InclusionListing - Datatable Inclusion listing with filter and order
func InclusionListing(r *http.Request, reqMap data.JQueryTableUI) (map[string]interface{}, error) {
	util.LogIt(r, "model - V_Inclusion - InclusionListing")
	var Qry, QryCnt, QryFilter, GroupBy, FromQry bytes.Buffer
	var status = make(map[string]interface{})

	var testColArrs [20]string
	testColArrs[0] = "CMT.id"
	testColArrs[1] = "inclusion"
	testColArrs[2] = "description"
	testColArrs[3] = "CMT.status"
	testColArrs[4] = "created_by"

	var testArrs []map[string]string
	testArrs = append(testArrs, map[string]string{
		"key":   "inclusion",
		"value": "CMT.inclusion",
	})
	testArrs = append(testArrs, map[string]string{
		"key":   "description",
		"value": "CMT.description",
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

	Qry.WriteString(" CMT.id,inclusion,description,CONCAT(from_unixtime(CMT.created_at),' ',SUC.username) AS created_by,ST.status,ST.id AS status_id ")

	FromQry.WriteString(" FROM cf_inclusion AS CMT ")
	FromQry.WriteString(" INNER JOIN cf_user AS SUC ON SUC.id = CMT.created_by ")
	FromQry.WriteString(" INNER JOIN status AS ST ON ST.id = CMT.status ")
	FromQry.WriteString(" WHERE CMT.status <> 3 ")
	Data, err := JQueryTable(r, reqMap, Qry, FromQry, QryCnt, QryFilter, GroupBy, testArrs, testColArrs, status, "")
	return Data, err
}

// GetInclusionList - Get Inclusion List For Other Module
func GetInclusionList(r *http.Request) ([]map[string]interface{}, error) {
	util.LogIt(r, "model - V_Inclusion - GetInclusionList")
	var Qry bytes.Buffer

	Qry.WriteString("SELECT id,inclusion,description FROM cf_inclusion WHERE status = 1")
	RetMap, err := ExecuteQuery(Qry.String())
	if util.CheckErrorLog(r, err) {
		return nil, err
	}

	return RetMap, nil
}

// GetInclusion -  Get Inclusion Detail By ID - 2021-04-20 - HK
func GetInclusion(r *http.Request, id string) (map[string]interface{}, error) {
	util.LogIt(r, "model - V_Inclusion - GetInclusion")
	var Qry bytes.Buffer

	Qry.WriteString("SELECT id, inclusion, description FROM cf_inclusion WHERE id = ?")
	RetMap, err := ExecuteRowQuery(Qry.String(), id)
	if util.CheckErrorLog(r, err) {
		return nil, err
	}

	return RetMap, nil
}
