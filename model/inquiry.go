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

// AddListingInquiry - Adds Property Listing Inquiry
func AddListingInquiry(r *http.Request, reqMap data.HotelInquiryModel) bool {
	util.LogIt(r, "Model - Inquiry - AddListingInquiry")

	var Qry bytes.Buffer
	nanoid, _ := gonanoid.Nanoid()
	InquirerIP := context.Get(r, "Visitor_IP").(string)
	datetime := util.GetIsoLocalDateTime()

	Qry.WriteString(" INSERT INTO fd_hotel_inquiry(id, first_name, last_name, phone_code, phone, email, property_name, city_id, state_id, country_id, property_address, zip_code, inquiry_datetime, ip) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?) ")
	err := ExecuteNonQuery(Qry.String(), nanoid, reqMap.FirstName, reqMap.LastName, reqMap.PhoneCode, reqMap.Phone, reqMap.Email, reqMap.PropertyName, reqMap.City, reqMap.State, reqMap.Country, reqMap.Address, reqMap.ZipCode, datetime, InquirerIP)
	if util.CheckErrorLog(r, err) {
		return false
	}

	CountryName, _ := GetModuleFieldByID(r, "COUNTRY", fmt.Sprintf("%.0f", reqMap.Country), "name")
	State, _ := GetModuleFieldByID(r, "STATE", fmt.Sprintf("%.0f", reqMap.State), "name")
	City, _ := GetModuleFieldByID(r, "CITY", fmt.Sprintf("%.0f", reqMap.City), "name")

	reqStruct := util.ToMap(reqMap)
	reqStruct["Country"] = CountryName
	reqStruct["State"] = State
	reqStruct["City"] = City

	AddLog(r, "", "HOTEL_INQUIRY", "Create", nanoid, GetLogsValueMap(r, reqStruct, false, ""))

	return true
}

// OnBoardInquiryListing - Datatable On Board Inquiry Listing Data
func OnBoardInquiryListing(r *http.Request, reqMap data.JQueryTableUI) (map[string]interface{}, error) {
	util.LogIt(r, "Model - Inquiry - OnBoardInquiryListing")

	var Qry, QryCnt, QryFilter, GroupBy, FromQry bytes.Buffer
	var status = make(map[string]interface{})

	var testColArrs [20]string
	testColArrs[0] = "FHI.id"
	testColArrs[1] = "FHI.first_name"
	testColArrs[2] = "FHI.last_name"
	testColArrs[3] = "FHI.email"
	testColArrs[4] = "FHI.phone"
	testColArrs[5] = "FHI.property_name"
	testColArrs[6] = "FHI.country_id"
	testColArrs[7] = "FHI.state_id"
	testColArrs[8] = "FHI.city_name"
	testColArrs[9] = "FHI.ip"
	testColArrs[10] = "FHI.status"
	testColArrs[11] = "created_by"

	var testArrs []map[string]string
	testArrs = append(testArrs, map[string]string{
		"key":   "first_name",
		"value": "FHI.first_name",
	})
	testArrs = append(testArrs, map[string]string{
		"key":   "last_name",
		"value": "FHI.last_name",
	})
	testArrs = append(testArrs, map[string]string{
		"key":   "email",
		"value": "FHI.email",
	})
	testArrs = append(testArrs, map[string]string{
		"key":   "phone",
		"value": "FHI.phone",
	})
	testArrs = append(testArrs, map[string]string{
		"key":   "country_id",
		"value": "FHI.country_id",
	})
	testArrs = append(testArrs, map[string]string{
		"key":   "state_id",
		"value": "FHI.state_id",
	})
	testArrs = append(testArrs, map[string]string{
		"key":   "city_id",
		"value": "FHI.city_id",
	})
	testArrs = append(testArrs, map[string]string{
		"key":   "property_name",
		"value": "FHI.property_name",
	})
	testArrs = append(testArrs, map[string]string{
		"key":   "status",
		"value": "FHI.status",
	})
	testArrs = append(testArrs, map[string]string{
		"key":   "created_at",
		"value": "DATE(from_unixtime(FHI.created_at))",
	})

	QryCnt.WriteString(" COUNT(FHI.id) AS cnt ")
	QryFilter.WriteString(" COUNT(FHI.id) AS cnt ")

	Qry.WriteString(" FHI.id, FHI.first_name, FHI.last_name, CONCAT('+',FHI.phone_code,'-',FHI.phone) AS phone, FHI.email, FHI.property_name, FHI.ip, from_unixtime(FHI.inquiry_datetime) as inquiry_date, ST.status, ST.id AS status_id, ")
	Qry.WriteString(" CC.name as city_name, FHI.city_id, ")
	Qry.WriteString(" CST.name as state_name, FHI.state_id, ")
	Qry.WriteString(" CCN.name as country_name, FHI.country_id ")

	FromQry.WriteString(" FROM fd_hotel_inquiry AS FHI ")
	FromQry.WriteString(" LEFT JOIN ")
	FromQry.WriteString(" cf_city AS CC ON CC.id = FHI.city_id ")
	FromQry.WriteString(" LEFT JOIN ")
	FromQry.WriteString(" cf_states AS CST ON CST.id = FHI.state_id ")
	FromQry.WriteString(" LEFT JOIN ")
	FromQry.WriteString(" cf_country AS CCN ON CCN.id = FHI.country_id ")
	FromQry.WriteString(" INNER JOIN status AS ST ON ST.id = FHI.status ")
	FromQry.WriteString(" WHERE FHI.status <> 3 ")
	Data, err := JQueryTable(r, reqMap, Qry, FromQry, QryCnt, QryFilter, GroupBy, testArrs, testColArrs, status, "")

	var Qry1 bytes.Buffer
	Qry1.WriteString(" SELECT id, status FROM status WHERE id IN(4,5)")
	StatusList, err := ExecuteQuery(Qry1.String())
	if util.CheckErrorLog(r, err) {
		return nil, err
	}

	Data["status"] = StatusList

	return Data, err
}

// GetInquiryDetailInfo - Gets Email Config Info
func GetInquiryDetailInfo(r *http.Request, id string) (map[string]interface{}, error) {

	util.LogIt(r, "Model - Inquiry - GetInquiryDetailInfo")

	var Qry bytes.Buffer
	Qry.WriteString(" SELECT FHI.id, FHI.first_name, FHI.last_name, FHI.email, FHI.phone_code, FHI.phone, FHI.property_name, FHI.ip, from_unixtime(FHI.inquiry_datetime) as inquiry_date, ST.status, ST.id AS status_id, ")
	Qry.WriteString(" CC.name as city_name, FHI.city_id, ")
	Qry.WriteString(" CST.name as state_name, FHI.state_id, ")
	Qry.WriteString(" CCN.name as country_name, FHI.country_id ")
	Qry.WriteString(" FROM fd_hotel_inquiry AS FHI ")
	Qry.WriteString(" LEFT JOIN ")
	Qry.WriteString(" cf_city AS CC ON CC.id = FHI.city_id ")
	Qry.WriteString(" LEFT JOIN ")
	Qry.WriteString(" cf_states AS CST ON CST.id = FHI.state_id ")
	Qry.WriteString(" LEFT JOIN ")
	Qry.WriteString(" cf_country AS CCN ON CCN.id = FHI.country_id ")
	Qry.WriteString(" INNER JOIN ")
	Qry.WriteString(" status AS ST ON ST.id = FHI.status ")
	Qry.WriteString(" WHERE FHI.id = ?")

	RetMap, err := ExecuteRowQuery(Qry.String(), id)
	if util.CheckErrorLog(r, err) {
		return nil, err
	}

	return RetMap, nil
}
