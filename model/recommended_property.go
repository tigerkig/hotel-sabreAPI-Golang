package model

import (
	"bytes"
	"net/http"
	"strconv"
	"tp-api-common/data"
	"tp-api-common/util"

	"github.com/gorilla/context"
	gonanoid "github.com/matoous/go-nanoid"
)

// AddHotelToRecommendedList - Adds Hotel Into Recommended List
func AddHotelToRecommendedList(r *http.Request, reqMap data.RecommendedHotel) bool {
	util.LogIt(r, "Model - Recommended_Hotel - AddHotelToRecommendedList")

	var Qry bytes.Buffer
	nanoid, _ := gonanoid.Nanoid()

	Qry.WriteString(" INSERT INTO cms_recommended_property(id, hotel_id, sort_order, created_at, created_by) VALUES (?,?,?,?,?) ")
	err := ExecuteNonQuery(Qry.String(), nanoid, reqMap.HotelID, reqMap.SortOrder, util.GetIsoLocalDateTime(), context.Get(r, "UserId"))
	if util.CheckErrorLog(r, err) {
		return false
	}

	HotelName, _ := GetModuleFieldByID(r, "HOTEL", reqMap.HotelID, "hotel_name")
	reqStruct := util.ToMap(reqMap)
	reqStruct["Hotel"] = HotelName

	AddLog(r, "", "RECOMMENDED_HOTEL", "Create", nanoid, GetLogsValueMap(r, reqStruct, true, "HotelID"))

	return true
}

// UpdateHotelToRecommendedList - Updates Hotel Into Recommended List
func UpdateHotelToRecommendedList(r *http.Request, reqMap data.RecommendedHotel) bool {
	util.LogIt(r, "Model - Recommended_Hotel - UpdateHotelToRecommendedList")

	var Qry bytes.Buffer
	Qry.WriteString("UPDATE cms_recommended_property SET hotel_id=?, sort_order=? WHERE id=?")
	err := ExecuteNonQuery(Qry.String(), reqMap.HotelID, reqMap.SortOrder, reqMap.ID)
	if util.CheckErrorLog(r, err) {
		return false
	}

	HotelName, _ := GetModuleFieldByID(r, "HOTEL", reqMap.HotelID, "hotel_name")
	reqStruct := util.ToMap(reqMap)
	reqStruct["Hotel"] = HotelName

	AddLog(r, "", "RECOMMENDED_HOTEL", "Update", reqMap.ID, GetLogsValueMap(r, reqStruct, true, "HotelID"))

	return true
}

// RecommendedHotelList - Datatable Listing Of Recommended Hotels
func RecommendedHotelList(r *http.Request, reqMap data.JQueryTableUI) (map[string]interface{}, error) {
	util.LogIt(r, "Model - Recommended_Hotel - RecommendedHotelList")

	var Qry, QryCnt, QryFilter, GroupBy, FromQry bytes.Buffer
	var status = make(map[string]interface{})

	var testColArrs [20]string
	testColArrs[0] = "CHI.id"
	testColArrs[1] = "CHI.hotel_name"
	testColArrs[2] = "country_name"
	testColArrs[3] = "state_name"
	testColArrs[4] = "city_name"
	testColArrs[5] = "CRP.sort_order"
	testColArrs[6] = "CRP.status"
	testColArrs[7] = "created_by"

	var testArrs []map[string]string
	testArrs = append(testArrs, map[string]string{
		"key":   "hotel_name",
		"value": "CHI.hotel_name",
	})
	testArrs = append(testArrs, map[string]string{
		"key":   "country_id",
		"value": "CHI.country_id",
	})
	testArrs = append(testArrs, map[string]string{
		"key":   "state_id",
		"value": "CHI.state_id",
	})
	testArrs = append(testArrs, map[string]string{
		"key":   "city_id",
		"value": "CHI.city_id",
	})
	testArrs = append(testArrs, map[string]string{
		"key":   "status",
		"value": "CRP.status",
	})
	testArrs = append(testArrs, map[string]string{
		"key":   "created_at",
		"value": "DATE(from_unixtime(CRP.created_at))",
	})

	QryCnt.WriteString(" COUNT(CRP.id) AS cnt ")
	QryFilter.WriteString(" COUNT(CRP.id) AS cnt ")

	Qry.WriteString(" CRP.id, CHI.hotel_name, CONCAT(from_unixtime(CRP.created_at),' ',SUC.username) AS created_by, ST.status, ST.id AS status_id, ")
	Qry.WriteString(" CCI.name as city_name, CHI.city_id, ")
	Qry.WriteString(" CFS.name as state_name, CHI.state_id, ")
	Qry.WriteString(" CCO.name as country_name, CHI.country_id ")

	FromQry.WriteString(" FROM cms_recommended_property AS CRP ")
	FromQry.WriteString(" INNER JOIN cf_hotel_info AS CHI ON CHI.id = CRP.hotel_id ")
	FromQry.WriteString(" INNER JOIN cf_country AS CCO ON CCO.id = CHI.country_id")
	FromQry.WriteString(" INNER JOIN cf_city AS CCI ON CCI.id = CHI.city_id")
	FromQry.WriteString(" INNER JOIN cf_states AS CFS ON CFS.id = CHI.state_id")
	FromQry.WriteString(" INNER JOIN cf_user AS SUC ON SUC.id = CRP.created_by ")
	FromQry.WriteString(" INNER JOIN status AS ST ON ST.id = CRP.status ")
	FromQry.WriteString(" WHERE CRP.status <> 3 ")

	Data, err := JQueryTable(r, reqMap, Qry, FromQry, QryCnt, QryFilter, GroupBy, testArrs, testColArrs, status, "")
	return Data, err
}

// GetRecommendedHotelInfo - Gets Recommended Hotel Info
func GetRecommendedHotelInfo(r *http.Request, id string) (map[string]interface{}, error) {

	util.LogIt(r, "Model - Recommended_Hotel - GetRecommendedHotelInfo")

	var Qry bytes.Buffer
	Qry.WriteString(" SELECT CRP.id, CRP.hotel_id, CHI.hotel_name, CRP.sort_order, ")
	Qry.WriteString(" CCI.name as city_name, CHI.city_id, ")
	Qry.WriteString(" CFS.name as state_name, CHI.state_id, ")
	Qry.WriteString(" CCO.name as country_name, CHI.country_id ")
	Qry.WriteString(" FROM cms_recommended_property AS CRP ")
	Qry.WriteString(" INNER JOIN cf_hotel_info AS CHI ON CHI.id = CRP.hotel_id ")
	Qry.WriteString(" INNER JOIN cf_country AS CCO ON CCO.id = CHI.country_id")
	Qry.WriteString(" INNER JOIN cf_city AS CCI ON CCI.id = CHI.city_id")
	Qry.WriteString(" INNER JOIN cf_states AS CFS ON CFS.id = CHI.state_id")
	Qry.WriteString(" INNER JOIN status AS ST ON ST.id = CRP.status ")
	Qry.WriteString(" WHERE CRP.id = ?")

	RetMap, err := ExecuteRowQuery(Qry.String(), id)
	if util.CheckErrorLog(r, err) {
		return nil, err
	}

	return RetMap, nil
}

// GetRecommendedPropertyCount - Get Recommended Property Count
func GetRecommendedPropertyCount(r *http.Request) (int64, error) {
	util.LogIt(r, "Model - Recommended_Hotel - GetRecommendedPropertyCount")

	var Qry bytes.Buffer
	Qry.WriteString("SELECT CONVERT(count(id),UNSIGNED INTEGER) AS cnt FROM cms_recommended_property")
	Data, err := ExecuteRowQuery(Qry.String())
	if util.CheckErrorLog(r, err) {
		return 0, err
	}

	IntV, _ := strconv.Atoi(Data["cnt"].(string))
	cnt := int64(IntV)

	return cnt, nil
}
