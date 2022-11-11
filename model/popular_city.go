package model

import (
	"bytes"
	"fmt"
	"net/http"
	"strconv"
	"tp-api-common/data"
	"tp-api-common/util"
	"tp-system/config"

	"github.com/gorilla/context"
	gonanoid "github.com/matoous/go-nanoid"
)

// AddPopularCity - Add Popular City
func AddPopularCity(r *http.Request, reqMap map[string]interface{}) bool {

	util.LogIt(r, "Model - Popular_City - AddPopularCity")
	nanoid, _ := gonanoid.Nanoid()

	var Qry bytes.Buffer
	Qry.WriteString("INSERT INTO cms_popular_city(id, city_id, city_name, image, description, sort_order, created_at, created_by) VALUES (?,?,?,?,?,?,?,?)")
	err := ExecuteNonQuery(Qry.String(), nanoid, reqMap["city_id"], reqMap["city_name"], reqMap["image"], reqMap["description"], reqMap["sort_order"], util.GetIsoLocalDateTime(), context.Get(r, "UserId"))
	if util.CheckErrorLog(r, err) {
		return false
	}

	AddLog(r, "", "POPULARCITY", "Create", nanoid, GetLogsValueMap(r, reqMap, true, ""))

	return true
}

// UpdatePopularCity - Update Popular City
func UpdatePopularCity(r *http.Request, reqMap map[string]interface{}) bool {
	util.LogIt(r, "Model - Popular_City - UpdatePopularCity")

	BeforeUpdate, _ := GetModuleFieldByID(r, "POPULARCITY", reqMap["id"].(string), "city_name")

	var Qry bytes.Buffer
	Qry.WriteString("UPDATE cms_popular_city SET city_id=?, city_name=?, image=?, description=?, sort_order=? WHERE id=?")
	err := ExecuteNonQuery(Qry.String(), reqMap["city_id"], reqMap["city_name"], reqMap["image"], reqMap["description"], reqMap["sort_order"], reqMap["id"])
	if util.CheckErrorLog(r, err) {
		return false
	}

	AddLog(r, BeforeUpdate.(string), "POPULARCITY", "Update", reqMap["id"].(string), GetLogsValueMap(r, reqMap, true, "id"))

	return true
}

// PopularCityListing - Return Datatable Listing Of Locality
func PopularCityListing(r *http.Request, reqMap data.JQueryTableUI) (map[string]interface{}, error) {

	util.LogIt(r, "Model - Popular_City - PopularCityListing")

	var Qry, QryCnt, QryFilter, GroupBy, FromQry bytes.Buffer
	var status = make(map[string]interface{})

	var testColArrs [20]string
	testColArrs[0] = "CPC.id"
	testColArrs[1] = "CPC.city_name"
	testColArrs[3] = "CPC.status"
	testColArrs[4] = "created_by"

	var testArrs []map[string]string
	testArrs = append(testArrs, map[string]string{
		"key":   "city",
		"value": "CPC.city_name",
	})
	testArrs = append(testArrs, map[string]string{
		"key":   "status",
		"value": "CPC.status",
	})
	testArrs = append(testArrs, map[string]string{
		"key":   "created_at",
		"value": "DATE(from_unixtime(CPC.created_at))",
	})

	QryCnt.WriteString(" COUNT(CPC.id) AS cnt ")
	QryFilter.WriteString(" COUNT(CPC.id) AS cnt ")

	Qry.WriteString(" CPC.id, CPC.city_id, CPC.city_name, ST.status, CONCAT(from_unixtime(CPC.created_at),' ',CU.username) AS created_by,ST.id AS status_id ")

	FromQry.WriteString(" FROM cms_popular_city AS CPC ")
	FromQry.WriteString(" INNER JOIN cf_user AS CU ON CU.id = CPC.created_by ")
	FromQry.WriteString(" INNER JOIN status AS ST ON ST.id = CPC.status ")
	FromQry.WriteString(" WHERE ST.id <> 3 ")
	Data, err := JQueryTable(r, reqMap, Qry, FromQry, QryCnt, QryFilter, GroupBy, testArrs, testColArrs, status, "")
	return Data, err
}

// GetPopularCityinfo -  Get Popular City Detail By ID
func GetPopularCityinfo(r *http.Request, id string) (map[string]interface{}, error) {
	util.LogIt(r, "Model - Popular_City - GetPopularCityinfo")

	var Qry bytes.Buffer
	Qry.WriteString(" SELECT ")
	Qry.WriteString(" CPC.id, CPC.description, CPC.sort_order, CPC.city_name, CPC.city_id, CONCAT('" + config.Env.AwsBucketURL + "popular_city/" + "',image) AS image, ")
	Qry.WriteString(" CST.name as state_name, CST.id as state_id, ")
	Qry.WriteString(" CCN.name as country_name, CCN.id as country_id ")
	Qry.WriteString(" FROM ")
	Qry.WriteString(" cms_popular_city AS CPC ")
	Qry.WriteString(" LEFT JOIN ")
	Qry.WriteString(" cf_city AS CC ON CC.id = CPC.city_id ")
	Qry.WriteString(" LEFT JOIN ")
	Qry.WriteString(" cf_states AS CST ON CST.id = CC.state_id ")
	Qry.WriteString(" LEFT JOIN ")
	Qry.WriteString(" cf_country AS CCN ON CCN.id = CST.country_id ")
	Qry.WriteString(" WHERE CPC.id = ?")
	RetMap, err := ExecuteRowQuery(Qry.String(), id)
	if util.CheckErrorLog(r, err) {
		return nil, err
	}

	return RetMap, nil
}

// GetPopularCityCount - Get Popular City Total Count
func GetPopularCityCount(r *http.Request, considerStatus bool) (int64, error) {
	util.LogIt(r, "Model - Popular_City - GetPopularCityCount")

	var Qry bytes.Buffer
	Qry.WriteString("SELECT CONVERT(count(id),UNSIGNED INTEGER) AS cnt FROM cms_popular_city")
	if considerStatus {
		Qry.WriteString(" WHERE status = 1 ")
	}
	Data, err := ExecuteRowQuery(Qry.String())
	if util.CheckErrorLog(r, err) {
		return 0, err
	}

	IntV, _ := strconv.Atoi(Data["cnt"].(string))
	cnt := int64(IntV)

	return cnt, nil
}

// UpdatePopularCityStatus - Update Status - 2021-04-26 - HK
func UpdatePopularCityStatus(r *http.Request, module string, status int, id string) (int, int, error) {
	util.LogIt(r, fmt.Sprint("Models - Models - UpdatePopularCityStatus - Module - ", module, " - Status - ", status, " id - ", id))
	var Qry, SQLQry bytes.Buffer

	Qry.WriteString("UPDATE cms_popular_city SET status = ? WHERE id = ?;")
	err := ExecuteNonQuery(Qry.String(), status, id)
	if util.CheckErrorLog(r, err) {
		return 0, 500, err
	}

	SQLQry.WriteString("SELECT status FROM status WHERE id=?")
	NewStatus, err := ExecuteRowQuery(SQLQry.String(), status)
	if util.CheckErrorLog(r, err) {
		return 0, 500, err
	}

	CommonOperation := "Update Status"
	AddLog(r, "", module, CommonOperation, id, map[string]interface{}{"Status": NewStatus["status"]})

	return 1, 204, nil
}
