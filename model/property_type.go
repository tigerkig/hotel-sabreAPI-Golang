package model

import (
	"bytes"
	"net/http"
	"tp-api-common/data"
	"tp-api-common/util"
	"tp-system/config"

	"github.com/gorilla/context"
	gonanoid "github.com/matoous/go-nanoid"
)

//PropertyType - Module Name
var PropertyType = "PROPERTY_TYPE"

// AddPropertyType - Add Property Type
func AddPropertyType(r *http.Request, reqMap data.PropertyType) bool {

	util.LogIt(r, "model - V_Property_Type - AddPropertyType")

	var Qry bytes.Buffer
	nanoid, _ := gonanoid.Nanoid()

	Qry.WriteString("INSERT INTO cf_property_type(id, type, created_at, created_by) VALUES (?,?,?,?)")
	err := ExecuteNonQuery(Qry.String(), nanoid, reqMap.Type, util.GetIsoLocalDateTime(), context.Get(r, "UserId"))
	if util.CheckErrorLog(r, err) {
		return false
	}

	AddLog(r, "", PropertyType, "Create", nanoid, GetLogsValueMap(r, util.ToMap(reqMap), false, ""))

	return true
}

// AddPropertyTypeNew - Add Property Type
func AddPropertyTypeNew(r *http.Request, reqMap map[string]interface{}) bool {

	util.LogIt(r, "Model - V_Property_Type - AddPropertyTypeNew")

	var Qry bytes.Buffer
	nanoid, _ := gonanoid.Nanoid()

	Qry.WriteString("INSERT INTO cf_property_type(id, type, image, created_at, created_by) VALUES (?,?,?,?,?)")
	err := ExecuteNonQuery(Qry.String(), nanoid, reqMap["type"], reqMap["image"], util.GetIsoLocalDateTime(), context.Get(r, "UserId"))
	if util.CheckErrorLog(r, err) {
		return false
	}

	AddLog(r, "", PropertyType, "Create", nanoid, GetLogsValueMap(r, reqMap, false, ""))

	return true
}

// UpdatePropertyType - Update Property Type
func UpdatePropertyType(r *http.Request, reqMap data.PropertyType) bool {
	util.LogIt(r, "model - V_Property_Type - UpdatePropertyType")
	var Qry bytes.Buffer

	BeforeUpdate, _ := GetModuleFieldByID(r, PropertyType, reqMap.ID, "type")

	Qry.WriteString("UPDATE cf_property_type SET type = ? WHERE id = ?")
	err := ExecuteNonQuery(Qry.String(), reqMap.Type, reqMap.ID)
	if util.CheckErrorLog(r, err) {
		return false
	}

	AddLog(r, BeforeUpdate.(string), PropertyType, "Update", reqMap.ID, GetLogsValueMap(r, util.ToMap(reqMap), true, "ID"))

	return true
}

// UpdatePropertyTypeNew - Update Property Type
func UpdatePropertyTypeNew(r *http.Request, reqMap map[string]interface{}) bool {
	util.LogIt(r, "model - V_Property_Type - UpdatePropertyTypeNew")
	var Qry bytes.Buffer

	BeforeUpdate, _ := GetModuleFieldByID(r, PropertyType, reqMap["id"].(string), "type")

	Qry.WriteString("UPDATE cf_property_type SET type = ?, image = ? WHERE id = ?")
	err := ExecuteNonQuery(Qry.String(), reqMap["type"], reqMap["image"], reqMap["id"])
	if util.CheckErrorLog(r, err) {
		return false
	}

	AddLog(r, BeforeUpdate.(string), PropertyType, "Update", reqMap["id"].(string), GetLogsValueMap(r, reqMap, true, "id"))

	return true
}

// PropertyTypeListing - Get Property Type Listing
func PropertyTypeListing(r *http.Request, reqMap data.JQueryTableUI) (map[string]interface{}, error) {
	util.LogIt(r, "model - V_Property_Type - PropertyTypeListing")
	var Qry, QryCnt, QryFilter, GroupBy, FromQry bytes.Buffer
	var status = make(map[string]interface{})

	var testColArrs [20]string
	testColArrs[0] = "CFT.id"
	testColArrs[1] = "type"
	testColArrs[2] = "CFT.status"
	testColArrs[3] = "created_by"

	var testArrs []map[string]string
	testArrs = append(testArrs, map[string]string{
		"key":   "type",
		"value": "CFT.type",
	})
	testArrs = append(testArrs, map[string]string{
		"key":   "status",
		"value": "CFT.status",
	})
	testArrs = append(testArrs, map[string]string{
		"key":   "created_at",
		"value": "DATE(from_unixtime(CFT.created_at))",
	})

	QryCnt.WriteString(" COUNT(CFT.id) AS cnt ")
	QryFilter.WriteString(" COUNT(CFT.id) AS cnt ")

	Qry.WriteString(" CFT.id, CFT.type, CONCAT(from_unixtime(CFT.created_at),' ',SUC.username) AS created_by, ST.status, ST.id AS status_id ")

	FromQry.WriteString(" FROM cf_property_type AS CFT ")
	FromQry.WriteString(" INNER JOIN cf_user AS SUC ON SUC.id = CFT.created_by  ")
	FromQry.WriteString(" INNER JOIN status AS ST ON ST.id = CFT.status ")
	FromQry.WriteString(" WHERE CFT.status <> 3 ")
	Data, err := JQueryTable(r, reqMap, Qry, FromQry, QryCnt, QryFilter, GroupBy, testArrs, testColArrs, status, "")
	return Data, err
}

// GetPropertyType - Update Property Type
func GetPropertyType(r *http.Request) ([]map[string]interface{}, bool) {
	util.LogIt(r, "model - V_Property_Type - GetPropertyType")
	var Qry bytes.Buffer
	Qry.WriteString("SELECT id, type FROM cf_property_type WHERE status=1")
	data, err := ExecuteQuery(Qry.String())
	if util.CheckErrorLog(r, err) {
		return nil, false
	}
	return data, true
}

// GetPropertyTypeInfo - Get Room Image List
func GetPropertyTypeInfo(r *http.Request, PropertyTypeID string) (map[string]interface{}, error) {
	util.LogIt(r, "Model - V_Property_Type - GetPropertyTypeInfo")

	var Qry bytes.Buffer

	Qry.WriteString("SELECT id, type, CONCAT('" + config.Env.AwsBucketURL + "property_type/" + "',image) AS image FROM cf_property_type WHERE id = ?")
	RetMap, err := ExecuteRowQuery(Qry.String(), PropertyTypeID)
	if util.CheckErrorLog(r, err) {
		return nil, err
	}

	return RetMap, nil
}

// GetPropertyTypeList - Get Property Type List For Other Module
func GetPropertyTypeList(r *http.Request) ([]map[string]interface{}, error) {

	util.LogIt(r, "model - V_Property_Type - GetPropertyTypeList")

	var Qry bytes.Buffer
	Qry.WriteString("SELECT id, type FROM cf_property_type WHERE status = 1")
	RetMap, err := ExecuteQuery(Qry.String())
	if util.CheckErrorLog(r, err) {
		return nil, err
	}

	return RetMap, nil
}
