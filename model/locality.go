package model

import (
	"bytes"
	"fmt"
	"net/http"
	"tp-api-common/data"
	"tp-api-common/util"

	"github.com/gorilla/context"
	gonanoid "github.com/matoous/go-nanoid"
)

// AddLocality - Add Locality
func AddLocality(r *http.Request, reqMap data.Locality) bool {

	util.LogIt(r, "model - V_Locality - AddLocality")
	nanoid, _ := gonanoid.Nanoid()

	var Qry bytes.Buffer
	Qry.WriteString("INSERT INTO cf_locality(id, locality, city_id, created_at, created_by) VALUES (?,?,?,?,?)")
	err := ExecuteNonQuery(Qry.String(), nanoid, reqMap.Locality, reqMap.City, util.GetIsoLocalDateTime(), context.Get(r, "UserId"))
	if util.CheckErrorLog(r, err) {
		return false
	}

	reqStruct := util.ToMap(reqMap)
	City, _ := GetModuleFieldByID(r, "CITY", fmt.Sprintf("%.0f", reqMap.City), "name")
	reqStruct["City"] = City

	AddLog(r, "", "LOCALITY", "Create", nanoid, GetLogsValueMap(r, reqStruct, true, ""))

	return true
}

// UpdateLocality - Update Locality
func UpdateLocality(r *http.Request, reqMap data.Locality) bool {

	util.LogIt(r, "model - V_Locality - UpdateLocality")

	var Qry bytes.Buffer
	Qry.WriteString("UPDATE cf_locality SET locality=?, city_id=? WHERE id=?")
	err := ExecuteNonQuery(Qry.String(), reqMap.Locality, reqMap.City, reqMap.ID)
	if util.CheckErrorLog(r, err) {
		return false
	}

	reqStruct := util.ToMap(reqMap)
	City, _ := GetModuleFieldByID(r, "CITY", fmt.Sprintf("%.0f", reqMap.City), "name")
	reqStruct["City"] = City

	AddLog(r, "", "LOCALITY", "Update", reqMap.ID, GetLogsValueMap(r, reqStruct, true, ""))

	return true
}

// GetLocality -  Get Locality Detail By ID
func GetLocality(r *http.Request, id string) (map[string]interface{}, error) {
	util.LogIt(r, "model - V_Locality - GetLocality")

	var Qry bytes.Buffer
	// Qry.WriteString("SELECT id, locality, city_id FROM cf_locality WHERE id = ?")
	Qry.WriteString("SELECT ")
	Qry.WriteString(" CL.id, CL.locality, CL.city_id, ")
	Qry.WriteString(" CC.name as city_name, ")
	Qry.WriteString(" CT.name as state_name, CT.id as state_id, ")
	Qry.WriteString(" CCN.name as country_name, CCN.id as country_id ")
	Qry.WriteString(" FROM ")
	Qry.WriteString(" cf_locality AS CL ")
	Qry.WriteString(" LEFT JOIN ")
	Qry.WriteString(" cf_city AS CC ON CC.id = CL.city_id ")
	Qry.WriteString(" LEFT JOIN ")
	Qry.WriteString(" cf_states AS CT ON CT.id = CC.state_id ")
	Qry.WriteString(" LEFT JOIN ")
	Qry.WriteString(" cf_country AS CCN ON CCN.id = CT.country_id ")
	Qry.WriteString(" WHERE CL.id = ?")
	RetMap, err := ExecuteRowQuery(Qry.String(), id)
	if util.CheckErrorLog(r, err) {
		return nil, err
	}

	return RetMap, nil
}

// LocalityListing - Return Datatable Listing Of Locality
func LocalityListing(r *http.Request, reqMap data.JQueryTableUI) (map[string]interface{}, error) {

	util.LogIt(r, "model - V_Locality - LocalityListing")

	var Qry, QryCnt, QryFilter, GroupBy, FromQry bytes.Buffer
	var status = make(map[string]interface{})

	var testColArrs [20]string
	testColArrs[0] = "CL.id"
	testColArrs[1] = "CL.locality"
	testColArrs[2] = "CC.name"
	testColArrs[3] = "CL.status"
	testColArrs[4] = "created_by"

	var testArrs []map[string]string
	testArrs = append(testArrs, map[string]string{
		"key":   "locality",
		"value": "CL.locality",
	})
	testArrs = append(testArrs, map[string]string{
		"key":   "city_id",
		"value": "CL.city_id",
	})
	testArrs = append(testArrs, map[string]string{
		"key":   "status",
		"value": "CL.status",
	})
	testArrs = append(testArrs, map[string]string{
		"key":   "created_at",
		"value": "DATE(from_unixtime(CL.created_at))",
	})

	QryCnt.WriteString(" COUNT(CL.id) AS cnt ")
	QryFilter.WriteString(" COUNT(CL.id) AS cnt ")

	Qry.WriteString(" CL.id, CL.locality, CC.name as city_name, ST.status, CONCAT(from_unixtime(CL.created_at),' ',CU.username) AS created_by,ST.id AS status_id ")

	FromQry.WriteString(" FROM cf_locality AS CL ")
	FromQry.WriteString(" INNER JOIN cf_city AS CC ON CC.id = CL.city_id ")
	FromQry.WriteString(" INNER JOIN cf_user AS CU ON CU.id = CL.created_by ")
	FromQry.WriteString(" INNER JOIN status AS ST ON ST.id = CL.status ")
	FromQry.WriteString(" WHERE ST.id <> 3 ")
	Data, err := JQueryTable(r, reqMap, Qry, FromQry, QryCnt, QryFilter, GroupBy, testArrs, testColArrs, status, "")
	return Data, err
}
