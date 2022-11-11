package model

import (
	"bytes"
	"net/http"
	"tp-api-common/data"
	"tp-api-common/util"

	"github.com/gorilla/context"
	gonanoid "github.com/matoous/go-nanoid"
)

//AmenityType - Module Name
var AmenityType = "AMENITY_TYPE"

// AddAmenityType - Add Amenity Type
func AddAmenityType(r *http.Request, reqMap data.AmenityType) bool {

	util.LogIt(r, "model - V_Amenity_Type - AddAmenityType")

	var Qry bytes.Buffer
	nanoid, _ := gonanoid.Nanoid()

	Qry.WriteString("INSERT INTO cf_amenity_type(id, type, amenity_of, created_at, created_by) VALUES (?,?,?,?,?)")
	err := ExecuteNonQuery(Qry.String(), nanoid, reqMap.Type, reqMap.AmenityOf, util.GetIsoLocalDateTime(), context.Get(r, "UserId"))
	if util.CheckErrorLog(r, err) {
		return false
	}

	reqStruct := util.ToMap(reqMap)
	if reqMap.AmenityOf == 1 {
		reqStruct["AmenityOf"] = "Property"
	} else {
		reqStruct["AmenityOf"] = "Room"
	}

	// AddLog(r, "", AmenityType, "Create", nanoid, GetLogsValueMap(r, structs.Map(reqMap), false, ""))
	AddLog(r, "", AmenityType, "Create", nanoid, GetLogsValueMap(r, reqStruct, false, ""))

	return true
}

// UpdateAmenityType - Update Amenity Type
func UpdateAmenityType(r *http.Request, reqMap data.AmenityType) bool {
	util.LogIt(r, "model - V_Amenity_Type - UpdateAmenityType")
	var Qry bytes.Buffer

	BeforeUpdate, _ := GetModuleFieldByID(r, AmenityType, reqMap.ID, "type")

	Qry.WriteString("UPDATE cf_amenity_type SET type = ?, amenity_of = ? WHERE id = ?")
	err := ExecuteNonQuery(Qry.String(), reqMap.Type, reqMap.AmenityOf, reqMap.ID)
	if util.CheckErrorLog(r, err) {
		return false
	}

	// AddLog(r, BeforeUpdate.(string), AmenityType, "Update", reqMap.ID, GetLogsValueMap(r, structs.Map(reqMap), true, "ID"))

	reqStruct := util.ToMap(reqMap)
	if reqMap.AmenityOf == 1 {
		reqStruct["AmenityOf"] = "Property"
	} else {
		reqStruct["AmenityOf"] = "Room"
	}

	AddLog(r, BeforeUpdate.(string), AmenityType, "Update", reqMap.ID, GetLogsValueMap(r, reqStruct, true, "ID"))

	return true
}

// GetAmenityTypeList - Get Amenity Type List For Other Module
func GetAmenityTypeList(r *http.Request) ([]map[string]interface{}, error) {
	util.LogIt(r, "model - V_Amenity_Type - GetAmenityTypeList")
	var Qry bytes.Buffer

	Qry.WriteString("SELECT id,type FROM cf_amenity_type WHERE status = 1")
	RetMap, err := ExecuteQuery(Qry.String())
	if util.CheckErrorLog(r, err) {
		return nil, err
	}

	return RetMap, nil
}

// AmenityTypeListing - Get Amenity Type Listing
func AmenityTypeListing(r *http.Request, reqMap data.JQueryTableUI) (map[string]interface{}, error) {
	util.LogIt(r, "model - V_Amenity_Type - AmenityTypeListing")
	var Qry, QryCnt, QryFilter, GroupBy, FromQry bytes.Buffer
	var status = make(map[string]interface{})

	var testColArrs [20]string
	testColArrs[0] = "CMT.id"
	testColArrs[1] = "type"
	testColArrs[2] = "CMT.amenity_of"
	testColArrs[3] = "CMT.status"
	testColArrs[4] = "created_by"

	var testArrs []map[string]string
	testArrs = append(testArrs, map[string]string{
		"key":   "type",
		"value": "CMT.type",
	})
	testArrs = append(testArrs, map[string]string{
		"key":   "amenity_of",
		"value": "CMT.amenity_of",
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

	Qry.WriteString(" CMT.id, type, CASE WHEN CMT.amenity_of = 1 THEN 'Property' ELSE 'Room' END AS amenity_of, CONCAT(from_unixtime(CMT.created_at),' ',SUC.username) AS created_by,ST.status,ST.id AS status_id ")

	FromQry.WriteString(" FROM cf_amenity_type AS CMT ")
	FromQry.WriteString(" INNER JOIN cf_user AS SUC ON SUC.id = CMT.created_by  ")
	FromQry.WriteString(" INNER JOIN status AS ST ON ST.id = CMT.status ")
	FromQry.WriteString(" WHERE CMT.status <> 3 ")
	Data, err := JQueryTable(r, reqMap, Qry, FromQry, QryCnt, QryFilter, GroupBy, testArrs, testColArrs, status, "")
	return Data, err
}

// GetAmenityTypeListCatgWise - Get Amenity Type List For Other Module
func GetAmenityTypeListCatgWise(r *http.Request, catgID string) ([]map[string]interface{}, error) {
	util.LogIt(r, "model - V_Amenity_Type - GetAmenityTypeListCatgWise")
	var Qry bytes.Buffer

	Qry.WriteString("SELECT id, type FROM cf_amenity_type WHERE status = 1 AND amenity_of = ?")
	RetMap, err := ExecuteQuery(Qry.String(), catgID)
	if util.CheckErrorLog(r, err) {
		return nil, err
	}

	return RetMap, nil
}

// GetAmenityType - Get Amenity Type Info
func GetAmenityType(r *http.Request, catgID string) (map[string]interface{}, error) {
	util.LogIt(r, "model - V_Amenity_Type - GetAmenityType")
	var Qry bytes.Buffer

	Qry.WriteString("SELECT id, type, amenity_of FROM cf_amenity_type WHERE id = ?")
	RetMap, err := ExecuteRowQuery(Qry.String(), catgID)
	if util.CheckErrorLog(r, err) {
		return nil, err
	}

	return RetMap, nil
}

// GetAmenityTypeListV1 - Get Amenity Type List For Other Module
func GetAmenityTypeListV1(r *http.Request, CatgID string) ([]map[string]interface{}, error) {
	util.LogIt(r, "model - V_Amenity_Type - GetAmenityTypeListV1")
	var Qry bytes.Buffer

	Qry.WriteString("SELECT id, type FROM cf_amenity_type WHERE status = 1 AND amenity_of = ?")
	RetMap, err := ExecuteQuery(Qry.String(), CatgID)
	if util.CheckErrorLog(r, err) {
		return nil, err
	}

	return RetMap, nil
}
