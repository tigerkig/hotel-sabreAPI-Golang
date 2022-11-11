package model

import (
	"bytes"
	"net/http"
	"tp-api-common/data"
	"tp-api-common/util"
)

// GetSmsGatewayDetail -  Return SMS Config Details
func GetSmsGatewayDetail(r *http.Request) (map[string]interface{}, error) {
	util.LogIt(r, "model - sms - GetSmsGatewayDetail")
	var Qry bytes.Buffer

	Qry.WriteString("SELECT id,gateway_name,api_end_point,api_username,api_password FROM cf_sms_gateway WHERE status = 1")
	RetMap, err := ExecuteRowQuery(Qry.String())
	if util.CheckErrorLog(r, err) {
		return nil, err
	}

	return RetMap, nil
}

// GetSmsTemplate -  Get SMS Template
func GetSmsTemplate(r *http.Request, id string) (map[string]interface{}, error) {
	util.LogIt(r, "model - sms - GetSmsTemplate")
	var Qry bytes.Buffer

	Qry.WriteString("SELECT id,name,subject,short_code,template,template_variable FROM cf_sms_template WHERE status = 1 AND id = ?")
	RetMap, err := ExecuteRowQuery(Qry.String(), id)
	if util.CheckErrorLog(r, err) {
		return nil, err
	}

	return RetMap, nil
}

// SMSTemplateListing - Return Datatable Listing Of SMS Template
func SMSTemplateListing(r *http.Request, reqMap data.JQueryTableUI) (map[string]interface{}, error) {
	util.LogIt(r, "model - Email_Template - SMSTemplateListing")
	var Qry, QryCnt, QryFilter, GroupBy, FromQry bytes.Buffer
	var status = make(map[string]interface{})

	var testColArrs [20]string
	testColArrs[0] = "CET.id"
	testColArrs[1] = "CET.name"
	testColArrs[2] = "CET.subject"
	testColArrs[3] = "CET.short_code"

	var testArrs []map[string]string

	QryCnt.WriteString(" COUNT(CET.id) AS cnt ")
	QryFilter.WriteString(" COUNT(CET.id) AS cnt ")

	Qry.WriteString(" CET.id, CET.subject, CET.name, CET.short_code ")

	FromQry.WriteString(" FROM cf_sms_template AS CET ")
	FromQry.WriteString(" WHERE 1 = 1 ")
	Data, err := JQueryTable(r, reqMap, Qry, FromQry, QryCnt, QryFilter, GroupBy, testArrs, testColArrs, status, "")
	return Data, err
}

// UpdateSMSTemplate - Update SMS Template
func UpdateSMSTemplate(r *http.Request, reqMap data.SMSTemplate) bool {
	util.LogIt(r, "Model - SMS - UpdateSMSTemplate")
	var Qry bytes.Buffer
	OldData, _ := GetModuleFieldsByID(r, "SMS_TEMPLATE", reqMap.ID, "name")

	Qry.WriteString("UPDATE cf_sms_template SET name=?, subject=?, template=?,template_variable=? WHERE id = ?")
	err := ExecuteNonQuery(Qry.String(), reqMap.TemplateName, reqMap.Subject, reqMap.Template, reqMap.TemplateVariable, reqMap.ID)
	if util.CheckErrorLog(r, err) {
		return false
	}

	AddLog(r, OldData["name"].(string), "SMS_TEMPLATE", "Update", reqMap.ID, GetLogsValueMap(r, util.ToMap(reqMap), true, "ID"))

	return true
}

// UpdateSMSGateway - Update SMS GATEWAY
func UpdateSMSGateway(r *http.Request, reqMap data.SmsGateway) bool {
	util.LogIt(r, "Model - SMS - UpdateSMSGateway")
	var Qry bytes.Buffer
	OldData, _ := GetModuleFieldsByID(r, "SMS_GATEWAY", "1", "gateway_name")

	Qry.WriteString("UPDATE cf_sms_gateway SET api_end_point=?, api_username=?, api_password=? WHERE id = 1")
	err := ExecuteNonQuery(Qry.String(), reqMap.EndPoint, reqMap.Username, reqMap.Password)
	if util.CheckErrorLog(r, err) {
		return false
	}

	AddLog(r, OldData["gateway_name"].(string), "SMS_GATEWAY", "Update", "1", GetLogsValueMap(r, util.ToMap(reqMap), true, "GatewayID,GatewayName"))

	return true
}
