package model

import (
	"bytes"
	"net/http"
	"strconv"
	"tp-api-common/data"
	"tp-api-common/util"
	"tp-system/config"
)

// VirtualCardList - Datatable Virtual Card listing with filter and order
func VirtualCardList(r *http.Request, reqMap data.JQueryTableUI) (map[string]interface{}, error) {
	util.LogIt(r, "model - Virtual_Card - VirtualCardList")
	var Qry, QryCnt, QryFilter, GroupBy, FromQry bytes.Buffer
	var status = make(map[string]interface{})

	var testColArrs [20]string
	testColArrs[0] = "id"
	testColArrs[1] = "booking_no"
	testColArrs[2] = "guest_name"
	testColArrs[3] = "amount"
	testColArrs[4] = "booking_status"
	testColArrs[5] = "cardno"

	var testArrs []map[string]string
	testArrs = append(testArrs, map[string]string{
		"key":   "booking_status",
		"value": "FDB.booking_status",
	})
	testArrs = append(testArrs, map[string]string{
		"key":   "guest_name",
		"value": "CONCAT(FDBG.first_name, ' ', FDBG.last_name)",
	})
	testArrs = append(testArrs, map[string]string{
		"key":   "booking_no",
		"value": "FDB.booking_no",
	})

	QryCnt.WriteString(" COUNT(CVC.id) AS cnt ")
	QryFilter.WriteString(" COUNT(CVC.id) AS cnt ")

	Qry.WriteString(" CVC.id, FDB.booking_no, CONCAT(FDBG.first_name, ' ', FDBG.last_name) AS guest_name, ROUND((spending_limit/100),2) AS amount, ")
	Qry.WriteString(" FDB.booking_status, CVC.cardno ")

	FromQry.WriteString(" FROM cf_virtual_card AS CVC ")
	FromQry.WriteString(" INNER JOIN tp_front.fd_booking AS FDB ON FDB.id = CVC.booking_id ")
	FromQry.WriteString(" INNER JOIN tp_front.fd_booking_guest AS FDBG ON FDBG.booking_id = FDB.id AND FDBG.is_primary=1 ")
	//FromQry.WriteString(" INNER JOIN tp_front.fd_account_detail AS FAD ON FAD.booking_id = FDB.id")
	FromQry.WriteString(" WHERE 1=1 ")
	if reqMap.ID != "" {
		FromQry.WriteString(" NAD FDB.hotel_id='" + reqMap.ID + "' ")
	}

	Data, err := JQueryTable(r, reqMap, Qry, FromQry, QryCnt, QryFilter, GroupBy, testArrs, testColArrs, status, "")
	return Data, err
}

// VirtualCardDetail - Get Virtual Card data
func VirtualCardDetail(r *http.Request, bookingID string) (map[string]interface{}, error) {
	util.LogIt(r, "Model - Virtual_Card - VirtualCardDetail")
	var cardSql, bookingSql bytes.Buffer

	cardSql.WriteString(" SELECT id, cardholder_id, cardno, name, currency, CAST(expmonth AS CHAR) as expmonth, CAST(expyear AS CHAR) as expyear, ")
	cardSql.WriteString(" CAST(ROUND((spending_limit / 100)) AS CHAR) as spending_limit, booking_id  FROM new_tp_system2.cf_virtual_card WHERE id = ?;")
	cardData, err := ExecuteRowQuery(cardSql.String(), bookingID)
	if chkError(err) {
		util.LogIt(r, "Model - Virtual_Card - VirtualCardDetail - cardSql - err")
		util.LogIt(r, err)
		return nil, err
	}

	if cardData == nil {
		return nil, err
	}

	bookingSql.WriteString(" SELECT booking_no, is_topup, booking_status, status FROM tp_front.fd_booking WHERE id = ?;")
	bookingData, err := ExecuteRowQuery(bookingSql.String(), cardData["booking_id"])
	if chkError(err) {
		util.LogIt(r, "Model - Virtual_Card - VirtualCardDetail - bookingSql - err")
		util.LogIt(r, err)
		return nil, err
	}

	if bookingData == nil {
		return nil, err
	}

	cardData["booking_no"] = bookingData["booking_no"].(string)
	cardData["is_topup"] = bookingData["is_topup"].(int64)
	cardData["booking_status"] = bookingData["booking_status"].(string)
	cardData["status"] = bookingData["status"].(int64)

	return cardData, nil
}

// VirtualCardUpdate - Virtual Card Status Update
func VirtualCardUpdate(r *http.Request, bookingID string) bool {
	util.LogIt(r, "Model - Virtual_Card - VirtualCardUpdate")
	var cardUpdateSQL, bookingUpdateSQL bytes.Buffer

	cardUpdateSQL.WriteString(" UPDATE new_tp_system2.cf_virtual_card SET type = 'active' WHERE booking_id = ?;")
	err := ExecuteNonQuery(cardUpdateSQL.String(), bookingID)
	if chkError(err) {
		util.LogIt(r, "Model - Virtual_Card - VirtualCardUpdate - cardUpdateSQL - err")
		util.LogIt(r, err)
		return false
	}

	bookingUpdateSQL.WriteString(" UPDATE tp_front.fd_booking SET is_topup = 1 WHERE id = ?; ")
	errBookingUpd := ExecuteNonQuery(bookingUpdateSQL.String(), bookingID)
	if chkError(errBookingUpd) {
		util.LogIt(r, "Model - Virtual_Card - VirtualCardUpdate - bookingUpdateSQL - errBookingUpd")
		util.LogIt(r, errBookingUpd)
		return false
	}
	return true
}

// VirtualCardActive - Set Virtual Card Active
func VirtualCardActive(r *http.Request, bookingID string) (map[string]interface{}, error) {
	util.LogIt(r, "Model - Virtual_Card - VirtualCardActive")
	var virtualCardActiveSql bytes.Buffer

	virtualCardActiveSql.WriteString(" SELECT FB.id as booking_id, FB.booking_no as booking_no, VC.id as card_id, VC.cardno, VC.spending_limit as amount FROM tp_front.fd_booking as FB ")
	virtualCardActiveSql.WriteString(" LEFT JOIN new_tp_system2.cf_virtual_card as VC ON VC.booking_id = FB.id ")
	virtualCardActiveSql.WriteString(" WHERE FB.id = ? AND FB.is_topup = 0 AND FB.booking_status = 'CHECKEDOUT' AND VC.type = 'inactive';")
	cardData, err := ExecuteRowQuery(virtualCardActiveSql.String(), bookingID)
	if chkError(err) {
		util.LogIt(r, "Model - Virtual_Card - VirtualCardActive - virtualCardActiveSql - err")
		util.LogIt(r, err)
		return nil, err
	}

	if cardData == nil {
		return nil, nil
	}
	return cardData, nil
}

// VirtualCardInfo - Get Card Holder data
func VirtualCardInfo(r *http.Request, bookingID string) (map[string]interface{}, error) {
	util.LogIt(r, "Model - Virtual_Card - VirtualCardInfo")
	var cardSql, cardHolderSql, companyInfoSql bytes.Buffer

	cardSql.WriteString(" SELECT id, cardholder_id, cardno, name, currency, CAST(expmonth AS CHAR) as expmonth, CAST(expyear AS CHAR) as expyear, CAST(ROUND((spending_limit / 100)) AS CHAR) as spending_limit  FROM new_tp_system2.cf_virtual_card WHERE booking_id = ?;")
	cardData, err := ExecuteRowQuery(cardSql.String(), bookingID)
	if chkError(err) {
		util.LogIt(r, "Model - Virtual_Card - VirtualCardInfo - cardSql - err")
		util.LogIt(r, err)
		return nil, err
	}

	if cardData == nil {
		return nil, err
	}

	cardHolderSql.WriteString(" SELECT * FROM new_tp_system2.cf_cardholder WHERE id = ?;")
	cardHolderData, errCardHolderData := ExecuteRowQuery(cardHolderSql.String(), cardData["cardholder_id"].(string))
	if chkError(errCardHolderData) {
		util.LogIt(r, "Model - Virtual_Card - VirtualCardInfo - cardHolderSql - errCardHolderData")
		util.LogIt(r, errCardHolderData)
		return nil, err
	}

	if cardHolderData == nil {
		return nil, err
	}

	companyInfoSql.WriteString(" SELECT CONCAT('" + config.Env.AwsBucketURL + "company_logo/" + "',image) AS image FROM new_tp_system2.cf_company_info WHERE id = 1;")
	companyInfoData, errCompanyInfo := ExecuteRowQuery(companyInfoSql.String())
	if chkError(errCompanyInfo) {
		util.LogIt(r, "Model - Virtual_Card - VirtualCardInfo - companyInfoSql - errCompanyInfo")
		util.LogIt(r, errCompanyInfo)
		return nil, err
	}

	cardData["email"] = cardHolderData["email"].(string)
	cardData["image"] = companyInfoData["image"].(string)
	return cardData, nil
}

// GetHotelDataForVirtualCard - Get Hotel Data for Card Holder Creation
func GetHotelDataForVirtualCard(r *http.Request, bookingID string) (map[string]interface{}, bool) {
	util.LogIt(r, "Model - Virtual_Card - GetHotelDataForVirtualCard")
	var hotelData = make(map[string]interface{})
	var amount float64
	var errPrice error
	var hotelSQL, hotelIDSQL, citySQL, stateSQL, countrySQL, transactionSQL bytes.Buffer

	hotelIDSQL.WriteString(" SELECT hotel_id FROM tp_front.fd_booking WHERE id = ?;")
	hotelID, err := ExecuteRowQuery(hotelIDSQL.String(), bookingID)
	if chkError(err) {
		return nil, false
	}

	hotelSQL.WriteString(" SELECT * FROM new_tp_system2.cf_hotel_info WHERE id = ?;")
	hotelDetails, err := ExecuteRowQuery(hotelSQL.String(), hotelID["hotel_id"].(string))
	if chkError(err) {
		return nil, false
	}

	citySQL.WriteString(" SELECT * FROM new_tp_system2.cf_city WHERE id = ?;")
	cityDetails, err := ExecuteRowQuery(citySQL.String(), hotelDetails["city_id"].(int64))
	if chkError(err) {
		return nil, false
	}

	stateSQL.WriteString(" SELECT * FROM new_tp_system2.cf_states WHERE id = ?;")
	stateDetails, err := ExecuteRowQuery(stateSQL.String(), hotelDetails["state_id"].(int64))
	if chkError(err) {
		return nil, false
	}

	countrySQL.WriteString(" SELECT * FROM new_tp_system2.cf_country WHERE id = ?;")
	countryDetails, err := ExecuteRowQuery(countrySQL.String(), hotelDetails["country_id"].(int64))
	if chkError(err) {
		return nil, false
	}

	transactionSQL.WriteString(" SELECT amount FROM tp_front.fd_transaction_info WHERE booking_id = ?;")
	transactionDetails, err := ExecuteRowQuery(transactionSQL.String(), bookingID)
	if chkError(err) {
		return nil, false
	}

	amount, errPrice = strconv.ParseFloat(transactionDetails["amount"].(string), 64)
	if chkError(errPrice) {
		return nil, false
	}

	hotelDetails["city_name"] = cityDetails["name"]
	hotelDetails["state_name"] = stateDetails["name"]
	hotelDetails["country_name"] = countryDetails["sortname"]
	hotelDetails["amount"] = amount * 100
	hotelDetails["hotel_id"] = hotelID["hotel_id"].(string)
	hotelDetails["hotel_name"] = util.TruncateString(hotelDetails["hotel_name"].(string), 24)

	hotelData = hotelDetails
	return hotelData, true
}

// CardholderDataFlag - Check whether Card Holder data for Hotelier is created or not
func CardholderDataFlag(r *http.Request, hotelData map[string]interface{}) bool {
	util.LogIt(r, "Model - Virtual_Card - CardholderDataFlag")
	var cardHolderSQL bytes.Buffer

	cardHolderSQL.WriteString(" SELECT * FROM new_tp_system2.cf_cardholder WHERE hotel_id = ?;")
	cardHolderData, err := ExecuteRowQuery(cardHolderSQL.String(), hotelData["hotel_id"].(string))
	if chkError(err) {
		util.LogIt(r, "Model - Virtual_Card - CardholderDataFlag - cardHolderSQL - err")
		util.LogIt(r, err)
		return false
	}

	if cardHolderData == nil {
		return false
	}
	return true
}

// CardholderDataInSystem - Sync cardholder data
func CardholderDataInSystem(r *http.Request, hotelData map[string]interface{}) bool {
	util.LogIt(r, "Model - Virtual_Card - CardholderDataInSystem")
	var cardHolderSQL bytes.Buffer

	cardHolderSQL.WriteString(" INSERT INTO new_tp_system2.cf_cardholder (id, hotel_id, name, email, phone, type, status) ")
	cardHolderSQL.WriteString(" VALUES (?, ?, ?, ?, ?, ?, ?); ")
	err := ExecuteNonQuery(cardHolderSQL.String(), hotelData["id"].(string), hotelData["hotel_id"].(string), hotelData["name"].(string), hotelData["email"].(string), hotelData["phoneno"].(string), hotelData["status"].(string), hotelData["type"].(string))
	if chkError(err) {
		util.LogIt(r, "Model - Virtual_Card - CardholderDataInSystem - cardHolderSQL - err")
		util.LogIt(r, err)
		return false
	}
	return true
}

// CardholderDataForHotel - Get Card Holder data
func CardholderDataForHotel(r *http.Request, hotelData map[string]interface{}) map[string]interface{} {
	util.LogIt(r, "Model - Virtual_Card - CardholderDataForHotel")
	var cardHolderSQL bytes.Buffer

	cardHolderSQL.WriteString(" SELECT * FROM new_tp_system2.cf_cardholder WHERE hotel_id = ?;")
	cardHolderData, err := ExecuteRowQuery(cardHolderSQL.String(), hotelData["hotel_id"].(string))
	if chkError(err) {
		util.LogIt(r, "Model - Virtual_Card - CardholderDataForHotel - cardHolderSQL - err")
		util.LogIt(r, err)
		return nil
	}

	if cardHolderData == nil {
		return nil
	}
	return cardHolderData
}

// VirtualCardDataInSystem - Sync virtual card data
func VirtualCardDataInSystem(r *http.Request, hotelData map[string]interface{}) bool {
	util.LogIt(r, "Model - Virtual_Card - VirtualCardDataInSystem")
	var virtualCardSQL bytes.Buffer
	/*log.Println(hotelData)
	log.Println("id")
	log.Println(reflect.TypeOf(hotelData["id"]))
	log.Println("cardholder_id")
	log.Println(reflect.TypeOf(hotelData["cardholder_id"]))
	log.Println("booking_id")
	log.Println(reflect.TypeOf(hotelData["booking_id"]))
	log.Println("currency")
	log.Println(reflect.TypeOf(hotelData["currency"]))
	log.Println("name")
	log.Println(reflect.TypeOf(hotelData["name"]))
	log.Println("expmonth")
	log.Println(reflect.TypeOf(hotelData["expmonth"]))
	log.Println("expyear")
	log.Println(reflect.TypeOf(hotelData["expyear"]))
	log.Println("spending_limit")
	log.Println(reflect.TypeOf(hotelData["spending_limit"]))
	log.Println("cardno")
	log.Println(reflect.TypeOf(hotelData["cardno"]))
	log.Println("status")
	log.Println(reflect.TypeOf(hotelData["status"]))
	log.Println("type")
	log.Println(reflect.TypeOf(hotelData["type"]))*/
	virtualCardSQL.WriteString(" INSERT INTO new_tp_system2.cf_virtual_card (id, booking_id, cardholder_id, name, currency, expmonth, expyear, spending_limit, cardno, type, status) ")
	virtualCardSQL.WriteString(" VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?); ")
	err := ExecuteNonQuery(virtualCardSQL.String(), hotelData["id"].(string), hotelData["booking_id"].(string), hotelData["cardholder_id"].(string), hotelData["name"].(string), hotelData["currency"].(string), hotelData["expmonth"].(int64), hotelData["expyear"].(int64), hotelData["spending_limit"].(float64), hotelData["cardno"].(string), hotelData["status"].(string), hotelData["type"].(string))
	if chkError(err) {
		util.LogIt(r, "Model - Virtual_Card - VirtualCardDataInSystem - virtualCardSQL - err")
		util.LogIt(r, err)
		return false
	}
	return true
}
