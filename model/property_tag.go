package model

import (
	"bytes"
	"net/http"
	"tp-api-common/data"
	"tp-api-common/util"

	"github.com/gorilla/context"
	gonanoid "github.com/matoous/go-nanoid"
)

// AddPropertyTag - Add Property Tag
func AddPropertyTag(r *http.Request, reqMap data.PropertyTags) bool {
	util.LogIt(r, "model - V_Property_Tag - AddPropertyTag")
	var Qry bytes.Buffer
	nanoid, _ := gonanoid.Nanoid()
	Qry.WriteString("INSERT INTO cf_tags(id,tag,created_at,created_by) VALUES (?,?,?,?)")
	err := ExecuteNonQuery(Qry.String(), nanoid, reqMap.Tag, util.GetIsoLocalDateTime(), context.Get(r, "UserId"))
	if util.CheckErrorLog(r, err) {
		return false
	}

	AddLog(r, "", "PROPERTY_TAG", "Create", nanoid, GetLogsValueMap(r, util.ToMap(reqMap), false, ""))

	return true
}

// UpdatePropPertyTag - Update Property Tag
func UpdatePropPertyTag(r *http.Request, reqMap data.PropertyTags) bool {
	util.LogIt(r, "model - V_Property_Tag - UpdatePropPertyTag")
	var Qry bytes.Buffer

	BeforeUpdate, _ := GetModuleFieldByID(r, "PROPERTY_TAG", reqMap.ID, "tag")

	Qry.WriteString("UPDATE cf_tags SET tag = ? WHERE id = ?")
	err := ExecuteNonQuery(Qry.String(), reqMap.Tag, reqMap.ID)
	if util.CheckErrorLog(r, err) {
		return false
	}

	AddLog(r, BeforeUpdate.(string), "PROPERTY_TAG", "Update", reqMap.ID, GetLogsValueMap(r, util.ToMap(reqMap), false, ""))

	return true
}

// PropPertyTagListing - Datatable Property Tag listing with filter and order
func PropPertyTagListing(r *http.Request, reqMap data.JQueryTableUI) (map[string]interface{}, error) {
	util.LogIt(r, "model - V_Property_Tag - PropPertyTagListing")
	var Qry, QryCnt, QryFilter, GroupBy, FromQry bytes.Buffer
	var status = make(map[string]interface{})

	var testColArrs [20]string
	testColArrs[0] = "CMT.id"
	testColArrs[1] = "tag"
	testColArrs[2] = "CMT.status"
	testColArrs[3] = "created_by"

	var testArrs []map[string]string
	testArrs = append(testArrs, map[string]string{
		"key":   "tag",
		"value": "CMT.tag",
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

	Qry.WriteString(" CMT.id,tag,CONCAT(from_unixtime(CMT.created_at),' ',SUC.username) AS created_by,ST.status,ST.id AS status_id ")

	FromQry.WriteString(" FROM cf_tags AS CMT ")
	FromQry.WriteString(" INNER JOIN cf_user AS SUC ON SUC.id = CMT.created_by ")
	FromQry.WriteString(" INNER JOIN status AS ST ON ST.id = CMT.status ")
	FromQry.WriteString(" WHERE CMT.status <> 3 ")
	Data, err := JQueryTable(r, reqMap, Qry, FromQry, QryCnt, QryFilter, GroupBy, testArrs, testColArrs, status, "")
	return Data, err
}

// GetPropertyTagList - Get Property Tag List For Other Module
func GetPropertyTagList(r *http.Request) ([]map[string]interface{}, error) {
	util.LogIt(r, "model - V_Property_Tag - GetPropertyTagList")
	var Qry bytes.Buffer

	Qry.WriteString("SELECT id,tag as text FROM cf_tags WHERE status = 1")
	RetMap, err := ExecuteQuery(Qry.String())
	if util.CheckErrorLog(r, err) {
		return nil, err
	}

	return RetMap, nil
}
