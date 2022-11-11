package model

import (
	"bytes"
	"net/http"
	"tp-api-common/data"
	"tp-api-common/util"
)

// UpdatePaymentConfiguration - Update Payment Configuration
func UpdatePaymentConfiguration(r *http.Request, reqMap data.PaymentGatewayConfiguration) bool {
	util.LogIt(r, "Model - Payment_Configuration - UpdatePaymentConfiguration")
	var Qry bytes.Buffer
	OldData, _ := GetModuleFieldsByID(r, "PAYMENT_GATEWAY", reqMap.ID, "payment_type")

	Qry.WriteString("UPDATE cf_payment_type SET payment_type=?, URL=?, auth_key=?, secret_key=?,is_surcharge=?,is_surcharge_inclusive=?,surcharge=?,country=? WHERE id = ?")
	err := ExecuteNonQuery(Qry.String(), reqMap.PaymentType, reqMap.URL, reqMap.AuthKey, reqMap.SecretKey, reqMap.IsSurcharge, reqMap.IsSurchargeInclusive, reqMap.Surcharge, reqMap.Country, reqMap.ID)
	if util.CheckErrorLog(r, err) {
		return false
	}

	//Logger map for logging values to mongoDB
	var logMap = make(map[string]interface{})
	var IsSurcharge string
	var IsSurchargeInclusive string
	logMap["Payment Type"] = reqMap.PaymentType
	logMap["URL"] = reqMap.URL
	logMap["Auth Key"] = reqMap.AuthKey
	logMap["Secret Key"] = reqMap.SecretKey
	if reqMap.IsSurcharge == 0 {
		IsSurcharge = "No"
	} else {
		IsSurcharge = "Yes"
	}
	logMap["Is Surcharge"] = IsSurcharge
	if reqMap.IsSurchargeInclusive == 0 {
		IsSurchargeInclusive = "No"
	} else {
		IsSurchargeInclusive = "Yes"
	}
	logMap["Is Surcharge Inclusive"] = IsSurchargeInclusive
	logMap["Surcharge"] = reqMap.Surcharge
	logMap["Country"] = reqMap.Country

	AddLog(r, OldData["payment_type"].(string), "PAYMENT_GATEWAY", "Update", reqMap.ID, GetLogsValueMap(r, logMap, true, "ID"))

	return true
}

// GetPaymentConfigDetail -  Return Payment Config Details
func GetPaymentConfigDetail(r *http.Request, id string) (map[string]interface{}, error) {
	util.LogIt(r, "model - Payment_Configuration - GetPaymentConfigDetail")
	var Qry bytes.Buffer

	Qry.WriteString("SELECT id,payment_type,URL,auth_key,secret_key,is_surcharge,is_surcharge_inclusive, surcharge,country FROM cf_payment_type WHERE id = ?")
	RetMap, err := ExecuteRowQuery(Qry.String(), id)
	if util.CheckErrorLog(r, err) {
		return nil, err
	}

	return RetMap, nil
}

// GetPaymentGatewayList -  Return Payment Config Details List
func GetPaymentGatewayList(r *http.Request) (map[string]interface{}, error) {
	util.LogIt(r, "model - Payment_Configuration - GetPaymentGatewayList")
	var Qry bytes.Buffer

	Qry.WriteString("SELECT id,payment_type,URL,auth_key,secret_key,is_surcharge,is_surcharge_inclusive, surcharge,country, status AS is_selected FROM cf_payment_type")
	RetMap, err := ExecuteQuery(Qry.String())
	if util.CheckErrorLog(r, err) {
		return nil, err
	}

	var retStuff = make(map[string]interface{})
	if RetMap == nil || len(RetMap) == 0 {
		retStuff["data"] = []string{}
	} else {
		retStuff["data"] = RetMap
	}

	return retStuff, nil
}

// ActivatePaymentGateway -  Activate Payment Gateway
func ActivatePaymentGateway(r *http.Request, id string) bool {
	util.LogIt(r, "model - Payment_Configuration - ActivatePaymentGateway")
	var Qry, Qry1 bytes.Buffer

	Qry.WriteString("UPDATE cf_payment_type SET status = 2")
	err := ExecuteNonQuery(Qry.String())
	if util.CheckErrorLog(r, err) {
		return false
	}

	Qry1.WriteString("UPDATE cf_payment_type SET status = 1 WHERE id = ?")
	err = ExecuteNonQuery(Qry1.String(), id)
	if util.CheckErrorLog(r, err) {
		return false
	}

	err = SetParameter("payment_type", id)
	if util.CheckErrorLog(r, err) {
		return false
	}

	return true
}

// GetActivatePaymentGateway -  Return Activate Payment Gateway
// Get active payment gateway data
func GetActivatePaymentGateway(r *http.Request) (map[string]interface{}, error) {
	util.LogIt(r, "model - Payment_Configuration - GetActivatePaymentGateway")
	var Qry bytes.Buffer
	ActivatePayment, _ := GetParameter("payment_type")
	Qry.WriteString("SELECT id,payment_type,URL,auth_key,secret_key,is_surcharge,is_surcharge_inclusive, surcharge,country, status AS is_selected FROM cf_payment_type WHERE id = ?")
	RetMap, err := ExecuteRowQuery(Qry.String(), ActivatePayment)
	if util.CheckErrorLog(r, err) {
		return nil, err
	}

	return RetMap, nil
}

// PaymentGatewayListing - Return Datatable Listing Of Payment Gateway Configuration
func PaymentGatewayListing(r *http.Request, reqMap data.JQueryTableUI) (map[string]interface{}, error) {
	util.LogIt(r, "model - Email_Template - EmailTemplateListing")
	var Qry, QryCnt, QryFilter, GroupBy, FromQry bytes.Buffer
	var status = make(map[string]interface{})

	var testColArrs [20]string
	testColArrs[0] = "id"
	testColArrs[1] = "payment_type"
	testColArrs[2] = "CASE WHEN is_surcharge = 0 THEN 'No' ELSE 'Yes' END"
	testColArrs[3] = "CASE WHEN is_surcharge_inclusive = 0 THEN 'No' ELSE 'Yes' END"
	testColArrs[4] = "status"

	var testArrs []map[string]string

	QryCnt.WriteString(" COUNT(id) AS cnt ")
	QryFilter.WriteString(" COUNT(id) AS cnt ")

	Qry.WriteString(" id,payment_type,CASE WHEN is_surcharge = 0 THEN 'No' ELSE 'Yes' END AS is_surcharge,CASE WHEN is_surcharge_inclusive = 0 THEN 'No' ELSE 'Yes' END AS is_surcharge_inclusive,surcharge,status ")

	FromQry.WriteString(" FROM cf_payment_type ")
	FromQry.WriteString(" WHERE 1 = 1 ")
	Data, err := JQueryTable(r, reqMap, Qry, FromQry, QryCnt, QryFilter, GroupBy, testArrs, testColArrs, status, "")
	return Data, err
}
