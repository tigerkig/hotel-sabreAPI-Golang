package model

import (
	"bytes"
	"fmt"
	"net/http"
	"tp-api-common/data"
	"tp-api-common/util"
	"tp-system/config"

	"log"

	"github.com/gorilla/context"
	gonanoid "github.com/matoous/go-nanoid"
)

// AddHotelInfo - Add Hotel Info
func AddHotelInfo(r *http.Request, reqMap data.Hotel) bool {
	util.LogIt(r, "model - V_Hotel - AddHotelInfo")
	var Qry, Cqry, CsQry bytes.Buffer
	nanoid, _ := gonanoid.Nanoid()
	nanoidSub, _ := gonanoid.Nanoid()
	nanoidSubM, _ := gonanoid.Nanoid()
	Qry.WriteString(" INSERT INTO cf_hotel_info(id,hotel_name,short_address,long_address,description,latitude,longitude,policy,created_at,created_by,hotel_star,city_id,state_id,country_id,hotel_phone,account_manager,locality_id,property_type_id) ") // 2020-06-15 - HK - Locality Added
	Qry.WriteString(" VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?) ")
	err := ExecuteNonQuery(Qry.String(), nanoid, reqMap.Name, reqMap.ShortAddress, reqMap.LongAddress, reqMap.Description, reqMap.Latitude, reqMap.Longitude, reqMap.Policy, util.GetIsoLocalDateTime(), context.Get(r, "UserId"), reqMap.HotelStar, reqMap.City, reqMap.State, reqMap.Country, reqMap.HotelPhone, reqMap.AccountManager, reqMap.Locality, reqMap.PropertyType)
	if util.CheckErrorLog(r, err) {
		return false
	}

	var password = util.GeneratePasswordHash(reqMap.Password)
	Cqry.WriteString(" INSERT INTO cf_hotel_client(id,client_name,username,password,phone_code1,mobile1,phone_code2,mobile2,email,created_at,created_by,hotel_id) ")
	Cqry.WriteString(" VALUES (?,?,?,?,?,?,?,?,?,?,?,?) ")
	err = ExecuteNonQuery(Cqry.String(), nanoidSub, reqMap.Manager, reqMap.Username, password, reqMap.PhoneCode1, reqMap.Phone1, reqMap.PhoneCode2, reqMap.Phone2, reqMap.Email, util.GetIsoLocalDateTime(), context.Get(r, "UserId"), nanoid)
	if util.CheckErrorLog(r, err) {
		return false
	}

	CsQry.WriteString(" INSERT INTO cf_hotel_settings(id,hotel_id,checkin_time,checkout_time) ")
	CsQry.WriteString(" VALUES (?,?,?,?) ")
	err = ExecuteNonQuery(CsQry.String(), nanoidSubM, nanoid, reqMap.CheckInTime, reqMap.CheckOutTime)
	if util.CheckErrorLog(r, err) {
		return false
	}

	CountryName, _ := GetModuleFieldByID(r, "COUNTRY", fmt.Sprintf("%.0f", reqMap.Country), "name")
	State, _ := GetModuleFieldByID(r, "STATE", fmt.Sprintf("%.0f", reqMap.State), "name")
	City, _ := GetModuleFieldByID(r, "CITY", fmt.Sprintf("%.0f", reqMap.City), "name")
	Locality, _ := GetModuleFieldByID(r, "LOCALITY", reqMap.Locality, "locality")          // 2020-06-15 - HK - Locality Added
	PropertyType, _ := GetModuleFieldByID(r, "PROPERTY_TYPE", reqMap.PropertyType, "type") // 2020-06-18 - HK - Property Type Added
	Account, _ := GetModuleFieldByID(r, "USER", reqMap.AccountManager, "username")
	reqMap.AccountManager = Account.(string)

	reqStruct := util.ToMap(reqMap)
	reqStruct["Short Address"] = reqMap.ShortAddress
	reqStruct["Long Address"] = reqMap.LongAddress
	reqStruct["Hotel Star"] = reqMap.HotelStar
	reqStruct["Hotel Phone"] = reqMap.HotelPhone
	reqStruct["Check In Time"] = reqMap.CheckInTime
	reqStruct["Check Out Time"] = reqMap.CheckOutTime
	reqStruct["Country Name"] = CountryName
	reqStruct["State Name"] = State
	reqStruct["City Name"] = City
	reqStruct["Area"] = Locality              // 2020-06-15 - HK - Locality Added
	reqStruct["Property Type"] = PropertyType // 2020-06-18 - HK - Property Type Added
	reqStruct["Account Manager"] = Account.(string)
	reqStruct["Latitude"] = fmt.Sprintf("%.4f", reqMap.Latitude)
	reqStruct["Longitude"] = fmt.Sprintf("%.4f", reqMap.Longitude)

	UpdateHotelOnList(nanoid) // 2020-06-24 - HK - Sync With Mongo Added - Admin Panel

	AddLog(r, "", "HOTEL", "Create", nanoid, GetLogsValueMap(r, reqStruct, true, "ShortAddress,LongAddress,HotelStar,CheckInTime,CheckOutTime,Password,State,City,Country,AccountManager,HotelPhone"))

	return true
}

// HotelListing - Datatable Hotel listing with filter and order
func HotelListing(r *http.Request, reqMap data.JQueryTableUI) (map[string]interface{}, error) {
	util.LogIt(r, "model - V_Hotel - HotelListing")
	var Qry, QryCnt, QryFilter, GroupBy, FromQry bytes.Buffer
	var status = make(map[string]interface{})

	var testColArrs [20]string
	testColArrs[0] = "CHI.id"
	testColArrs[1] = "hotel_name"
	testColArrs[2] = "short_address"
	testColArrs[3] = "hotel_star"
	testColArrs[4] = "username"
	testColArrs[5] = "CMT.status"
	testColArrs[6] = "created_by"

	var testArrs []map[string]string
	testArrs = append(testArrs, map[string]string{
		"key":   "hotel_name",
		"value": "CHI.hotel_name",
	})
	testArrs = append(testArrs, map[string]string{
		"key":   "city",
		"value": "CHI.city_id",
	})
	testArrs = append(testArrs, map[string]string{
		"key":   "state",
		"value": "CHI.state_id",
	})
	testArrs = append(testArrs, map[string]string{
		"key":   "country",
		"value": "CHI.country_id",
	})
	testArrs = append(testArrs, map[string]string{
		"key":   "hotel_star",
		"value": "CHI.hotel_star",
	})
	testArrs = append(testArrs, map[string]string{
		"key":   "username",
		"value": "CHC.username",
	})
	testArrs = append(testArrs, map[string]string{
		"key":   "short_address",
		"value": "CHI.short_address",
	})
	testArrs = append(testArrs, map[string]string{
		"key":   "status",
		"value": "CHI.status",
	})
	testArrs = append(testArrs, map[string]string{
		"key":   "group_id",
		"value": "CHI.group_id",
	})
	testArrs = append(testArrs, map[string]string{
		"key":   "created_at",
		"value": "DATE(from_unixtime(CHI.created_at))",
	})

	QryCnt.WriteString(" COUNT(CHI.id) AS cnt ")
	QryFilter.WriteString(" COUNT(CHI.id) AS cnt ")

	Qry.WriteString(" CHI.is_approved,CONCAT(from_unixtime(CHI.approved_at),' ', CU.username) AS approved_by,CHI.is_live,CHS.commission_amount,CHS.is_auto_set_deal,CHS.auto_set_days, CHI.id,CHI.hotel_name, CHI.hotel_star, CHI.short_address, CHI.latitude, CHI.longitude, CHC.client_name, CHC.username,  ")
	Qry.WriteString(" CONCAT(CHC.phone_code1,'',CHC.mobile1) AS phone, CONCAT(from_unixtime(CHI.created_at)) AS created_by, CHI.status AS status_id,ST.status, CHC.email ")

	FromQry.WriteString(" FROM cf_hotel_info AS CHI ")
	//FromQry.WriteString(" INNER JOIN cf_hotel_client AS CHC ON CHC.hotel_id = CHI.id ")
	FromQry.WriteString(" INNER JOIN cf_hotel_client AS CHC ON CHC.group_id = CHI.group_id AND CHC.group_id <> '' ")
	FromQry.WriteString(" INNER JOIN cf_hotel_settings AS CHS ON CHS.hotel_id = CHI.id ")
	//FromQry.WriteString(" INNER JOIN cf_user AS CU ON CU.id = CHI.created_by ")
	FromQry.WriteString(" LEFT JOIN cf_user AS CU ON CU.id = CHI.approved_by ")
	FromQry.WriteString(" INNER JOIN status AS ST ON ST.id = CHI.status ")
	FromQry.WriteString(" WHERE CHI.status <> 3 ")
	Data, err := JQueryTable(r, reqMap, Qry, FromQry, QryCnt, QryFilter, GroupBy, testArrs, testColArrs, status, "")
	return Data, err
}

// ViewHotelInfo - View Hotel Info Pass By ID
func ViewHotelInfo(r *http.Request, id string) (map[string]interface{}, error) {
	util.LogIt(r, "model - V_Hotel - ViewHotelInfo")
	var SelectQry, Qry, HotelClientInfo, HotelInfoView, HotelSetting, HotelInfoQry, HotelClientInfoQry, HotelSettingQry, TagQry, MainTagQry bytes.Buffer
	SelectQry.WriteString(" SELECT ")

	HotelInfoView.WriteString(" CHI.is_approved,hotel_name,CHI.status AS status_id, ST.status, CHI.short_address, CHI.long_address,CHI.hotel_star, CHI.description, CHI.latitude, ")
	HotelInfoView.WriteString(" CHI.hotel_phone, CHI.longitude, CHI.policy, CHI.locality_id, CHI.city_id,CHI.state_id,CHI.country_id, CFC.name AS country, CFS.name AS state, CFCY.name AS city,  ")
	HotelInfoView.WriteString(" CHI.locality_id, CFLC.locality,  ") // 2020-06-18 - HK - Locality Returned
	HotelInfoView.WriteString(" CHI.property_type_id, CPT.type,  ") // 2020-06-18 - HK - Property Type ID Returned
	HotelInfoView.WriteString(" CU.username AS account_manager_name, CHI.account_manager ")
	HotelClientInfo.WriteString(" CHC.client_name,CHC.username,CONCAT(CHC.phone_code1,CHC.mobile1) AS phone1, CONCAT(CHC.phone_code2,CHC.mobile2) AS phone2, CHC.email ")
	HotelSetting.WriteString(" CHS.commission_amount, CHS.auto_set_days, CHS.is_auto_set_deal, CHS.checkin_time, CHS.checkout_time,CHS.account_holder_name,CHS.account_number,CHS.swift_code,CHS.bank_name ")
	TagQry.WriteString(" CASE WHEN GROUP_CONCAT(tag) IS NULL THEN '' ELSE GROUP_CONCAT(tag) END AS tag,CASE WHEN GROUP_CONCAT(tag_id) IS NULL THEN '' ELSE GROUP_CONCAT(tag_id) END AS tag_id FROM cf_hotel_tag AS CHT INNER JOIN cf_tags AS CFT ON CFT.id = CHT.tag_id WHERE hotel_id = ? ")

	Qry.WriteString(" FROM cf_hotel_info AS CHI ")
	//Qry.WriteString(" INNER JOIN cf_hotel_client AS CHC ON CHC.hotel_id = CHI.id ")
	Qry.WriteString(" INNER JOIN cf_hotel_client AS CHC ON CHC.group_id = CHI.group_id ")
	Qry.WriteString(" LEFT JOIN cf_hotel_settings AS CHS ON CHS.hotel_id = CHI.id ")
	Qry.WriteString(" LEFT JOIN cf_country AS CFC ON CFC.id = CHI.country_id ")
	Qry.WriteString(" LEFT JOIN cf_states AS CFS ON CFS.id = CHI.state_id ")
	Qry.WriteString(" LEFT JOIN cf_city AS CFCY ON CFCY.id = CHI.city_id ")
	Qry.WriteString(" LEFT JOIN cf_locality AS CFLC ON CFLC.id = CHI.locality_id ")
	Qry.WriteString(" LEFT JOIN cf_property_type AS CPT ON CPT.id = CHI.property_type_id ")
	Qry.WriteString(" LEFT JOIN cf_user AS CU ON CU.id = CHI.account_manager ")
	Qry.WriteString(" LEFT JOIN cf_user AS CU1 ON CU1.id = CHI.created_by ")
	Qry.WriteString(" LEFT JOIN status AS ST ON ST.id = CHI.status WHERE CHI.id = ? ")

	HotelInfoQry.WriteString(SelectQry.String())
	HotelInfoQry.WriteString(HotelInfoView.String())
	HotelInfoQry.WriteString(Qry.String())

	HotelClientInfoQry.WriteString(SelectQry.String())
	HotelClientInfoQry.WriteString(HotelClientInfo.String())
	HotelClientInfoQry.WriteString(Qry.String())

	HotelSettingQry.WriteString(SelectQry.String())
	HotelSettingQry.WriteString(HotelSetting.String())
	HotelSettingQry.WriteString(Qry.String())

	MainTagQry.WriteString(SelectQry.String())
	MainTagQry.WriteString(TagQry.String())

	var mainStuff = make(map[string]interface{})
	HotelInfo, err := ExecuteRowQuery(HotelInfoQry.String(), id)
	if util.CheckErrorLog(r, err) {
		util.LogIt(r, fmt.Sprintf("Hotel Info query error"))
		return nil, err
	}
	HotelClient, err := ExecuteRowQuery(HotelClientInfoQry.String(), id)
	if util.CheckErrorLog(r, err) {
		util.LogIt(r, fmt.Sprintf("Hotel Client query error"))
		return nil, err
	}
	HotelSettingData, err := ExecuteRowQuery(HotelSettingQry.String(), id)
	if util.CheckErrorLog(r, err) {
		util.LogIt(r, fmt.Sprintf("Hotel SettingData query error"))
		return nil, err
	}
	MainTagData, err := ExecuteRowQuery(MainTagQry.String(), id)
	if util.CheckErrorLog(r, err) {
		util.LogIt(r, fmt.Sprintf("Hotel MainTagData query error"))
		return nil, err
	}
	HotelInfo["tag"] = MainTagData["tag"]
	HotelInfo["tag_id"] = MainTagData["tag_id"]
	//Category Wise Image
	var CatQry bytes.Buffer
	CatQry.WriteString("SELECT category_id,CIC.name AS category FROM cf_hotel_image AS CFI INNER JOIN cf_image_category AS CIC ON CIC.id = CFI.category_id WHERE hotel_id = ? GROUP BY category_id")
	CatArray, err := ExecuteQuery(CatQry.String(), id)
	if util.CheckErrorLog(r, err) {
		util.LogIt(r, fmt.Sprintf("Hotel CatQry query error"))
		return nil, err
	}
	if len(CatArray) > 0 {
		for i := 0; i < len(CatArray); i++ {
			var CatQry bytes.Buffer
			CatQry.WriteString("SELECT id, CONCAT('" + config.Env.AwsBucketURL + "hotel/" + "',image) AS image, sortorder FROM cf_hotel_image WHERE hotel_id = ? AND category_id = ? ORDER BY sortorder")
			ImageArray, err := ExecuteQuery(CatQry.String(), id, CatArray[i]["category_id"])
			if util.CheckErrorLog(r, err) {
				util.LogIt(r, fmt.Sprintf("Hotel ImageArray query error"))
				return nil, err
			}

			CatArray[i]["image"] = ImageArray
			delete(CatArray[i], "category_id")
		}
		mainStuff["hotel_category_image"] = CatArray
	} else {
		mainStuff["hotel_category_image"] = make(map[string]interface{})
	}
	var AmenityArrData = make(map[string]interface{})

	if context.Get(r, "Side").(string) == "TP-PARTNER" {
		AmenityArrData, err = AmenityTypeWiseAmenity(r, id)
		if util.CheckErrorLog(r, err) {
			return nil, err
		}
	} else {
		AmenityArrData, err = AmenityTypeWiseAmenityAdmin(r, id)
		if util.CheckErrorLog(r, err) {
			return nil, err
		}
	}

	if len(AmenityArrData["data"].([]map[string]interface{})) > 0 {
		mainStuff["hotel_amenity"] = AmenityArrData
	} else {
		mainStuff["hotel_amenity"] = make(map[string]interface{})
	}

	//Policy Data
	var PolicyQry bytes.Buffer
	PolicyQry.WriteString("SELECT policy,CASE WHEN checkin_rules IS NULL THEN '' ELSE checkin_rules END AS checkin_rules  FROM cf_hotel_info WHERE id = ? ")
	PolicyData, err := ExecuteRowQuery(PolicyQry.String(), id)
	if util.CheckErrorLog(r, err) {
		return nil, err
	}

	PolicyData["checkin_time"] = HotelSettingData["checkin_time"]
	PolicyData["checkout_time"] = HotelSettingData["checkout_time"]
	delete(HotelSettingData, "checkin_time")
	delete(HotelSettingData, "checkout_time")

	if len(PolicyData) > 0 {
		mainStuff["hotel_policy"] = PolicyData
	} else {
		mainStuff["hotel_policy"] = make(map[string]interface{})
	}

	mainStuff["hotel_info"] = HotelInfo
	mainStuff["hotel_client"] = HotelClient
	mainStuff["hotel_setting"] = HotelSettingData

	return mainStuff, nil
}

// GetHotelInfo - Get Hotel Info Pass By ID
func GetHotelInfo(r *http.Request, id string) (map[string]interface{}, error) {
	util.LogIt(r, "model - V_Hotel - ViewHotelInfo")
	var SelectQry, Qry, HotelClientInfoGet, HotelInfoView, HotelSetting, HotelInfoQry, HotelClientInfoQry, HotelSettingQry, TagQry, MainTagQry bytes.Buffer
	SelectQry.WriteString(" SELECT ")

	HotelInfoView.WriteString(" hotel_name,CHI.status AS status_id, ST.status, CHI.short_address, CHI.long_address,CHI.hotel_star, CHI.description, CHI.latitude, ")
	HotelInfoView.WriteString(" CHI.hotel_phone, CHI.longitude, CHI.policy, CHI.city_id,CHI.state_id,CHI.country_id, CFC.name AS country, CFS.name AS state, CFCY.name AS city,  ")
	HotelInfoView.WriteString(" CHI.locality_id, CFLC.locality,  ") // 2020-06-15 - HK - Locality Returned
	HotelInfoView.WriteString(" CHI.property_type_id, CPT.type,  ") // 2020-06-18 - HK - Locality Returned
	HotelInfoView.WriteString(" CU.username AS account_manager_name, CHI.account_manager ")
	HotelClientInfoGet.WriteString(" CHC.client_name,CHC.username,CHC.phone_code1,CHC.mobile1,CHC.phone_code2,CHC.mobile2,CHC.email ")
	HotelSetting.WriteString(" CHS.commission_amount, CHS.auto_set_days, CHS.is_auto_set_deal, CHS.checkin_time, CHS.checkout_time,CHS.account_holder_name,CHS.account_number,CHS.swift_code,CHS.bank_name ")
	TagQry.WriteString(" CASE WHEN GROUP_CONCAT(tag) IS NULL THEN '' ELSE GROUP_CONCAT(tag) END AS tag,CASE WHEN GROUP_CONCAT(CFT.id) IS NULL THEN '' ELSE GROUP_CONCAT(CFT.id) END AS tag_id FROM cf_hotel_tag AS CHT INNER JOIN cf_tags AS CFT ON CFT.id = CHT.tag_id WHERE hotel_id = ? ")

	Qry.WriteString(" FROM cf_hotel_info AS CHI ")
	//Qry.WriteString(" INNER JOIN cf_hotel_client AS CHC ON CHC.hotel_id = CHI.id ")
	Qry.WriteString(" INNER JOIN cf_hotel_client AS CHC ON CHC.group_id = CHI.group_id ")
	Qry.WriteString(" LEFT JOIN cf_hotel_settings AS CHS ON CHS.hotel_id = CHI.id ")
	Qry.WriteString(" LEFT JOIN cf_country AS CFC ON CFC.id = CHI.country_id ")
	Qry.WriteString(" LEFT JOIN cf_states AS CFS ON CFS.id = CHI.state_id ")
	Qry.WriteString(" LEFT JOIN cf_city AS CFCY ON CFCY.id = CHI.city_id ")
	Qry.WriteString(" LEFT JOIN cf_locality AS CFLC ON CFLC.id = CHI.locality_id ")
	Qry.WriteString(" LEFT JOIN cf_property_type AS CPT ON CPT.id = CHI.property_type_id ")
	Qry.WriteString(" LEFT JOIN cf_user AS CU ON CU.id = CHI.account_manager ")
	//Qry.WriteString(" INNER JOIN cf_user AS CU1 ON CU1.id = CHI.created_by ")
	Qry.WriteString(" INNER JOIN status AS ST ON ST.id = CHI.status WHERE CHI.id = ? ")

	HotelInfoQry.WriteString(SelectQry.String())
	HotelInfoQry.WriteString(HotelInfoView.String())
	HotelInfoQry.WriteString(Qry.String())

	HotelClientInfoQry.WriteString(SelectQry.String())
	HotelClientInfoQry.WriteString(HotelClientInfoGet.String())
	HotelClientInfoQry.WriteString(Qry.String())

	HotelSettingQry.WriteString(SelectQry.String())
	HotelSettingQry.WriteString(HotelSetting.String())
	HotelSettingQry.WriteString(Qry.String())

	MainTagQry.WriteString(SelectQry.String())
	MainTagQry.WriteString(TagQry.String())

	var mainStuff = make(map[string]interface{})

	HotelInfo, err := ExecuteRowQuery(HotelInfoQry.String(), id)
	log.Println(HotelInfo, "1")
	if util.CheckErrorLog(r, err) {
		return nil, err
	}

	HotelClient, err := ExecuteRowQuery(HotelClientInfoQry.String(), id)
	if util.CheckErrorLog(r, err) {
		return nil, err
	}

	HotelSettingData, err := ExecuteRowQuery(HotelSettingQry.String(), id)
	if util.CheckErrorLog(r, err) {
		return nil, err
	}

	MainTagData, err := ExecuteRowQuery(MainTagQry.String(), id)
	if util.CheckErrorLog(r, err) {
		return nil, err
	}

	HotelInfo["tag"] = MainTagData["tag"]
	HotelInfo["tag_id"] = MainTagData["tag_id"]

	mainStuff["hotel_info"] = HotelInfo
	mainStuff["hotel_client"] = HotelClient
	mainStuff["hotel_setting"] = HotelSettingData

	return mainStuff, nil
}

// UpdateCommissionSetting - UpdateCommission Setting
func UpdateCommissionSetting(r *http.Request, reqMap data.HotelCommission) bool {
	util.LogIt(r, "model - V_Hotel - UpdateCommissionSetting")
	var Qry bytes.Buffer

	Qry.WriteString("UPDATE cf_hotel_settings SET commission_amount = ?, is_auto_set_deal = ?, auto_set_days = ? WHERE hotel_id = ?")
	err := ExecuteNonQuery(Qry.String(), reqMap.CommissionAmount, reqMap.AutoSetDeal, reqMap.AutoSetDays, reqMap.HotelID)
	if util.CheckErrorLog(r, err) {
		return false
	}

	reqStruct := util.ToMap(reqMap)
	if reqMap.AutoSetDeal == 1 {
		reqStruct["Is Auto Set Deal"] = "Yes"
	} else {
		reqStruct["Is Auto Set Deal"] = "No"
	}

	reqStruct["Commission Amount"] = reqMap.CommissionAmount
	reqStruct["Auto Set Days"] = reqMap.AutoSetDays

	AddLog(r, "", "HOTEL", "Update Commission", reqMap.HotelID, GetLogsValueMap(r, reqStruct, true, "CommissionAmount,AutoSetDeal,AutoSetDays,HotelID"))

	return true
}

// UpdateAccountManager - Update Assign Manager
func UpdateAccountManager(r *http.Request, reqMap data.Hotel) bool {
	util.LogIt(r, "model - V_Hotel - UpdateAccountManager")
	var Qry bytes.Buffer

	Qry.WriteString("UPDATE cf_hotel_info SET account_manager = ? WHERE id = ?")
	err := ExecuteNonQuery(Qry.String(), reqMap.AccountManager, reqMap.ID)
	if util.CheckErrorLog(r, err) {
		return false
	}

	Account, _ := GetModuleFieldByID(r, "USER", reqMap.AccountManager, "username")
	reqMap.AccountManager = Account.(string)
	retMap := make(map[string]interface{})
	retMap["Account Manager"] = Account.(string)

	AddLog(r, "", "HOTEL", "Update Account Manager", reqMap.ID, retMap)

	return true
}

// UpdateHotelInfo - Update Hotel Info
func UpdateHotelInfo(r *http.Request, reqMap data.Hotel) bool {
	util.LogIt(r, "model - V_Hotel - UpdateHotelInfo")
	var Qry, Cqry, CsQry bytes.Buffer
	BeforeUpdate, _ := GetModuleFieldByID(r, "HOTEL", reqMap.ID, "hotel_name")

	// 2020-06-18 - HK - PropertyType Added
	// 2020-06-15 - HK - Locality Added
	Qry.WriteString(" UPDATE cf_hotel_info SET hotel_name=?,short_address=?,long_address=?,description=?,latitude=?,longitude=?,policy=?,hotel_star=?,city_id=?,state_id=?,country_id=?,hotel_phone=?,locality_id=?,property_type_id=? ")
	Qry.WriteString(" WHERE id = ? ")
	err := ExecuteNonQuery(Qry.String(), reqMap.Name, reqMap.ShortAddress, reqMap.LongAddress, reqMap.Description, reqMap.Latitude, reqMap.Longitude, reqMap.Policy, reqMap.HotelStar, reqMap.City, reqMap.State, reqMap.Country, reqMap.HotelPhone, reqMap.Locality, reqMap.PropertyType, reqMap.ID)
	if util.CheckErrorLog(r, err) {
		return false
	}

	Cqry.WriteString(" UPDATE cf_hotel_client SET client_name=?,phone_code1=?,mobile1=?,phone_code2=?,mobile2=?,email=? WHERE hotel_id=? ")
	err = ExecuteNonQuery(Cqry.String(), reqMap.Manager, reqMap.PhoneCode1, reqMap.Phone1, reqMap.PhoneCode2, reqMap.Phone2, reqMap.Email, reqMap.ID)
	if util.CheckErrorLog(r, err) {
		return false
	}

	CsQry.WriteString(" UPDATE cf_hotel_settings SET checkin_time=?,checkout_time=? WHERE hotel_id=? ")
	err = ExecuteNonQuery(CsQry.String(), reqMap.CheckInTime, reqMap.CheckOutTime, reqMap.ID)
	if util.CheckErrorLog(r, err) {
		return false
	}

	CountryName, _ := GetModuleFieldByID(r, "COUNTRY", fmt.Sprintf("%.0f", reqMap.Country), "name")
	State, _ := GetModuleFieldByID(r, "STATE", fmt.Sprintf("%.0f", reqMap.State), "name")
	City, _ := GetModuleFieldByID(r, "CITY", fmt.Sprintf("%.0f", reqMap.City), "name")
	Locality, _ := GetModuleFieldByID(r, "LOCALITY", reqMap.Locality, "locality")          // 2020-06-15 - HK - Locality Added
	PropertyType, _ := GetModuleFieldByID(r, "PROPERTY_TYPE", reqMap.PropertyType, "type") // 2020-06-18 - HK - Property Type Added

	reqStruct := util.ToMap(reqMap)
	reqStruct["Short Address"] = reqMap.ShortAddress
	reqStruct["Long Address"] = reqMap.LongAddress
	reqStruct["Hotel Star"] = fmt.Sprintf("%.0f", reqMap.HotelStar)
	reqStruct["Hotel Phone"] = reqMap.HotelPhone
	reqStruct["Check In Time"] = reqMap.CheckInTime
	reqStruct["Check Out Time"] = reqMap.CheckOutTime
	reqStruct["Country Name"] = CountryName
	reqStruct["State Name"] = State
	reqStruct["City Name"] = City
	reqStruct["Area"] = Locality              // 2020-06-15 - HK - Locality Added
	reqStruct["Property Type"] = PropertyType // 2020-06-18 - HK - Property Type Added
	reqStruct["Latitude"] = fmt.Sprintf("%.4f", reqMap.Latitude)
	reqStruct["Longitude"] = fmt.Sprintf("%.4f", reqMap.Longitude)

	UpdateHotelOnList(reqMap.ID) // 2020-06-24 - HK - Sync With Mongo Added - Admin Panel

	AddLog(r, BeforeUpdate.(string), "HOTEL", "Update", reqMap.ID, GetLogsValueMap(r, reqStruct, true, "ShortAddress,LongAddress,HotelStar,CheckInTime,CheckOutTime,Password,State,City,Country,AccountManager,HotelPhone,PropertyType"))

	return true
}

// ResetHotelUserPwd - Reset Password Of Hotel User
func ResetHotelUserPwd(r *http.Request, reqMap data.Hotel) bool {
	util.LogIt(r, "model - V_Hotel - ResetHotelUserPwd")

	var SQry bytes.Buffer

	var password = util.GeneratePasswordHash(reqMap.Password)
	//SQry.WriteString("UPDATE cf_hotel_client SET password=? WHERE hotel_id = ?")
	SQry.WriteString("UPDATE cf_hotel_client SET password=? WHERE group_id = ?")
	err := ExecuteNonQuery(SQry.String(), password, reqMap.ID)
	if util.CheckErrorLog(r, err) {
		return false
	}

	AddLog(r, "", "HOTEL", "Update Password", reqMap.ID, map[string]interface{}{})

	return true
}

// UpdateBankDetails - Update Bank Details Setting
func UpdateBankDetails(r *http.Request, reqMap data.HotelBankDetails) bool {
	util.LogIt(r, "model - V_Hotel - UpdateBankDetails")
	var Qry bytes.Buffer

	Qry.WriteString("UPDATE cf_hotel_settings SET account_holder_name=?,account_number = ?, swift_code = ?, bank_name = ? WHERE hotel_id = ?")
	err := ExecuteNonQuery(Qry.String(), reqMap.AccountHolderName, reqMap.AccountNumber, reqMap.SwiftCode, reqMap.Bank, reqMap.HotelID)
	if util.CheckErrorLog(r, err) {
		return false
	}

	reqStruct := util.ToMap(reqMap)
	reqStruct["Account Number"] = reqMap.AccountNumber
	reqStruct["Swift Code"] = reqMap.SwiftCode
	reqStruct["Account Holder Name"] = reqMap.AccountHolderName
	reqStruct["Bank"] = reqMap.Bank

	AddLog(r, "", "HOTEL", "Update Bank Details", reqMap.HotelID, GetLogsValueMap(r, reqStruct, true, "AccountHolderName,AccountNumber,SwiftCode,HotelID"))

	return true
}

// GetHotelListForOtherModule - Get Property Type List For Other Module
func GetHotelListForOtherModule(r *http.Request) ([]map[string]interface{}, error) {
	util.LogIt(r, "Model - V_Hotel - GetHotelListForOtherModule")

	city := r.URL.Query().Get("city")
	state := r.URL.Query().Get("state")
	country := r.URL.Query().Get("country")

	var Qry bytes.Buffer
	Qry.WriteString(" SELECT CHI.id, CHI.hotel_name, ")
	Qry.WriteString(" CHI.city_id, CCT.name as city, ")
	Qry.WriteString(" CHI.state_id, CST.name as state, ")
	Qry.WriteString(" CHI.country_id, CNT.name as country ")
	Qry.WriteString(" FROM cf_hotel_info AS CHI ")
	Qry.WriteString(" INNER JOIN cf_hotel_group AS CHG ON CHG.id = CHI.group_id AND CHG.id <> '' ")
	Qry.WriteString(" INNER JOIN cf_city AS CCT ON CCT.id = CHI.city_id ")
	Qry.WriteString(" INNER JOIN cf_states AS CST ON CST.id = CHI.state_id ")
	Qry.WriteString(" INNER JOIN cf_country AS CNT ON CNT.id = CHI.country_id ")
	Qry.WriteString(" WHERE CHI.status=1 ")

	if city != "" {
		Qry.WriteString(" AND CHI.city_id = " + city)
	}

	if state != "" {
		Qry.WriteString(" AND CHI.state_id = " + state)
	}

	if country != "" {
		Qry.WriteString(" AND CHI.country_id = " + country)
	}

	RetMap, err := ExecuteQuery(Qry.String())
	if util.CheckErrorLog(r, err) {
		return nil, err
	}

	return RetMap, nil
}

// UpdateHotelStatusToLive - Updates Hotel Status As Live And It Gets Listed In Mongo Search
func UpdateHotelStatusToLive(r *http.Request, status int, hotelID string) (bool, error) {
	util.LogIt(r, "Model - V_Hotel - UpdateHotelStatusToLive")

	var Qry, SQLQry bytes.Buffer

	Qry.WriteString("UPDATE cf_hotel_info SET is_live = ? WHERE id = ?;")
	err := ExecuteNonQuery(Qry.String(), status, hotelID)
	if util.CheckErrorLog(r, err) {
		return false, err
	}

	SQLQry.WriteString("SELECT status FROM status WHERE id=?")
	NewStatus, err := ExecuteRowQuery(SQLQry.String(), status)
	if util.CheckErrorLog(r, err) {
		return false, err
	}

	var listFlg bool
	if status == 1 {
		listFlg = UpdateSearchList(hotelID, "List")
		if !listFlg {
			return false, err
		}
	} else {
		listFlg = UpdateSearchList(hotelID, "Unlist")
		if !listFlg {
			return false, err
		}
	}

	AddLog(r, "", "HOTEL", "Update Status", hotelID, map[string]interface{}{"Status": NewStatus["status"]})

	return true, nil
}

// UpdateHotelierLoginStatus - Update Hotelier Login Status
func UpdateHotelierLoginStatus(r *http.Request, id string, status int) bool {
	util.LogIt(r, "model - V_Hotel - UpdateHotelierLoginStatus")
	var Qry bytes.Buffer
	Qry.WriteString("UPDATE cf_hotel_client SET status = ? WHERE hotel_id = ?")
	if status == 1 {
		err := ExecuteNonQuery(Qry.String(), 1, id)
		if util.CheckErrorLog(r, err) {
			return false
		}
	} else if status == 2 {
		err := ExecuteNonQuery(Qry.String(), 2, id)
		if util.CheckErrorLog(r, err) {
			return false
		}
		//cache
		UpdateHotelStatusToLive(r, 2, id)
	}

	return true
}

// IsPartnerContainProperty - Check if pass property own by partner or not
var IsPartnerContainProperty = func(r *http.Request, hotelID, userID string) bool {
	var Qry bytes.Buffer
	var Panel = context.Get(r, "Side").(string)
	if hotelID == "" {
		return true
	}
	if Panel == "TP-BACKOFFICE" {
		return true
	} else {
		//check if user id is null or not
		if userID == "" {
			userID = context.Get(r, "UserId").(string)
		}
		Qry.WriteString(" SELECT CHI.id FROM cf_hotel_info AS CHI ")
		Qry.WriteString(" INNER JOIN cf_hotel_group AS CHG ON CHG.id = CHI.group_id ")
		Qry.WriteString(" WHERE CHG.client_id = ? AND CHI.id = ? ")
		Data, err := ExecuteQuery(Qry.String(), userID, hotelID)
		if util.CheckErrorLog(r, err) {
			return false
		}

		if len(Data) > 0 {
			return true
		}

		return false
	}
}

func ApprovedHotel(r *http.Request, reqMap data.ApprovedHotel) bool {
	util.LogIt(r, "model - V_Hotel - ApprovedHotel")
	var Qry bytes.Buffer
	Qry.WriteString(" UPDATE cf_hotel_info SET is_approved = ?, approved_at = ?, approved_by = ? WHERE id = ? ")
	err := ExecuteNonQuery(Qry.String(), reqMap.IsApproved, util.GetIsoLocalDateTime(), context.Get(r, "UserId"), reqMap.HotelID)
	if util.CheckErrorLog(r, err) {
		return false
	}

	if reqMap.IsApproved == 1 {
		if !HotelStripeAccountCreated(r, reqMap.HotelID) {
			if !AddAccountIDToHotelInfo(r, reqMap.HotelID) {
				err := fmt.Errorf("Error while creating stripe account id " + reqMap.HotelID)
				util.LogIt(r, fmt.Sprintf("Error while creating stripe account id "+reqMap.HotelID))
				util.CheckErrorLog(r, err)
			}
		}
	}

	return true
}

// HotelierListing - Datatable HotelierListing listing with filter and order
func HotelierListing(r *http.Request, reqMap data.JQueryTableUI) (map[string]interface{}, error) {
	util.LogIt(r, "model - V_Hotel - HotelierListing")
	var Qry, QryCnt, QryFilter, GroupBy, FromQry bytes.Buffer
	var status = make(map[string]interface{})

	var testColArrs [20]string
	testColArrs[0] = "CHC.id"
	testColArrs[1] = "CHC.client_name"
	testColArrs[2] = "CHC.username"
	testColArrs[3] = "CHC.email"
	testColArrs[4] = "CONCAT(CHC.phone_code1,'',CHC.mobile1)"
	testColArrs[5] = "hotel_cnt"
	testColArrs[6] = "CHC.status"
	testColArrs[7] = "created_at"

	var testArrs []map[string]string
	testArrs = append(testArrs, map[string]string{
		"key":   "client_name",
		"value": "CHC.client_name",
	})
	testArrs = append(testArrs, map[string]string{
		"key":   "username",
		"value": "CHC.username",
	})
	testArrs = append(testArrs, map[string]string{
		"key":   "email",
		"value": "CHC.username",
	})
	testArrs = append(testArrs, map[string]string{
		"key":   "status",
		"value": "CHC.status",
	})
	testArrs = append(testArrs, map[string]string{
		"key":   "created_at",
		"value": "DATE(from_unixtime(CHC.created_at))",
	})

	QryCnt.WriteString(" COUNT(CHC.id) AS cnt ")
	QryFilter.WriteString(" COUNT(CHC.id) AS cnt ")

	Qry.WriteString(" CHC.id, CHC.client_name, CHC.username, CHC.email,CHC.status,CHC.group_id,  ")
	Qry.WriteString(" CONCAT(CHC.phone_code1,'',CHC.mobile1) AS phone, from_unixtime(CHC.created_at) AS created_at,(SELECT COUNT(id) FROM cf_hotel_info WHERE group_id = CHC.group_id AND status <> 3) AS hotel_cnt ")

	FromQry.WriteString(" FROM cf_hotel_client AS CHC ")
	FromQry.WriteString(" INNER JOIN cf_hotel_group AS CHG ON CHG.id = CHC.group_id AND CHC.group_id <> '' ")
	FromQry.WriteString(" WHERE 1 = 1 ")
	Data, err := JQueryTable(r, reqMap, Qry, FromQry, QryCnt, QryFilter, GroupBy, testArrs, testColArrs, status, "")
	return Data, err
}

// UpdateBooking - Update booking status
func UpdateBooking(r *http.Request, id string, reqData data.Status) error {
	util.LogIt(r, "model - V_Hotel - UpdateBooking")

	var Qry bytes.Buffer
	Qry.WriteString(" UPDATE tp_front.fd_booking SET status = ? WHERE id = ? AND status = 3;")
	err := ExecuteNonQuery(Qry.String(), reqData.Status, reqData.ID)
	if chkError(err) {
		util.LogIt(r, "model - V_Hotel - UpdateBooking - Qry - 1 - Error")
		util.LogIt(r, err)
		return err
	}

	// log.Println(reqData.Status)

	if reqData.Status == 1 {

		// Sending Mail to Client
		MailChn <- MailObj{
			Type: "UserEmailTemplate",
			ID:   reqData.ID}

		// Sending Mail to Admin

	} else if reqData.Status == 2 {

		// Sending Mail to Client
		MailChn <- MailObj{
			Type: "UserEmailTemplate",
			ID:   reqData.ID}
	}

	return nil
}
