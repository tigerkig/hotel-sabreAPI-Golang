package model

import (
	"bytes"
	"net/http"
	"strings"
	"tp-api-common/data"
	"tp-api-common/util"
)

// UpdateEmailTemplate - Update Email Template
func UpdateEmailTemplate(r *http.Request, reqMap data.EmailTemplate) bool {
	util.LogIt(r, "Model - Email_Template - UpdateEmailTemplate")
	var BccStr string
	var Qry bytes.Buffer
	OldData, _ := GetModuleFieldsByID(r, "EMAIL_TEMPLATE", reqMap.ID, "email_from,email_bcc,template,email_template_name,subject,short_code")

	Qry.WriteString("UPDATE cf_email_template SET email_from=?, email_bcc=?, template=?, subject=?, template_variable=? WHERE id = ?")
	err := ExecuteNonQuery(Qry.String(), reqMap.Email, reqMap.Bcc, reqMap.Template, reqMap.Subject, reqMap.TemplateVariable, reqMap.ID)
	if util.CheckErrorLog(r, err) {
		return false
	}

	if reqMap.Bcc != "" {
		str := strings.Split(reqMap.Bcc, ",")
		if len(str) > 0 {
			for i, val := range str {
				BccStrVal, _ := GetModuleFieldsByID(r, "EMAIL_CONFIG", val, "email_id")
				if i == 1 {
					BccStr = BccStrVal["email_id"].(string)
				} else {
					BccStr = BccStr + " , " + BccStrVal["email_id"].(string)
				}
			}
		}
	}

	EmailStrVal, _ := GetModuleFieldsByID(r, "EMAIL_CONFIG", reqMap.Email, "email_id")

	reqMap.Bcc = BccStr
	reqMap.Email = EmailStrVal["email_id"].(string)

	AddLog(r, OldData["email_template_name"].(string), "EMAIL_TEMPLATE", "Update", reqMap.ID, GetLogsValueMap(r, util.ToMap(reqMap), true, "ID,Template"))

	return true
}

// GetEmailTemplate -  Return Email Template Details
func GetEmailTemplate(r *http.Request, id string) (map[string]interface{}, error) {
	util.LogIt(r, "model - Email_Template - GetEmailTemplate")
	var Qry bytes.Buffer

	Qry.WriteString("SELECT id,email_template_name,email_from,email_bcc,subject,short_code,template, template_variable FROM cf_email_template WHERE id = ?")
	RetMap, err := ExecuteRowQuery(Qry.String(), id)
	if util.CheckErrorLog(r, err) {
		return nil, err
	}

	return RetMap, nil
}

// EmailTemplateListing - Return Datatable Listing Of Email Template
func EmailTemplateListing(r *http.Request, reqMap data.JQueryTableUI) (map[string]interface{}, error) {
	util.LogIt(r, "model - Email_Template - EmailTemplateListing")
	var Qry, QryCnt, QryFilter, GroupBy, FromQry bytes.Buffer
	var status = make(map[string]interface{})

	var testColArrs [20]string
	testColArrs[0] = "CET.id"
	testColArrs[1] = "CET.email_template_name"
	testColArrs[2] = "CET.subject"
	testColArrs[3] = "CET.short_code"
	testColArrs[4] = "CEC.email_id"

	var testArrs []map[string]string

	QryCnt.WriteString(" COUNT(CET.id) AS cnt ")
	QryFilter.WriteString(" COUNT(CET.id) AS cnt ")

	Qry.WriteString(" CET.id, CET.subject, CET.email_template_name, CEC.email_id, CET.short_code, CET.template_variable ")

	FromQry.WriteString(" FROM cf_email_template AS CET ")
	FromQry.WriteString(" INNER JOIN cf_email_config AS CEC ON CEC.id = CET.email_from ")
	FromQry.WriteString(" WHERE 1 = 1 ")
	Data, err := JQueryTable(r, reqMap, Qry, FromQry, QryCnt, QryFilter, GroupBy, testArrs, testColArrs, status, "")
	return Data, err
}

// GetEmailTemplateDetailInfo -  Return Email Template Details
func GetEmailTemplateDetailInfo(r *http.Request, id string) (map[string]interface{}, error) {
	util.LogIt(r, "model - Email_Template - GetEmailTemplateDetailInfo")
	var Qry bytes.Buffer

	Qry.WriteString(" SELECT CET.subject,CET.id,CET.email_template_name,CEC.email_id,CET.short_code,CEC.email_name,CEC.smtp_host,CEC.smtp_port,CEC.smtp_user,from_base64(CEC.smtp_password) AS smtp_password,CEC.signature, IFNULL(GROUP_CONCAT(CEC1.email_id),'') AS bcc, CET.template ")
	Qry.WriteString(" FROM cf_email_template AS CET ")
	Qry.WriteString(" INNER JOIN cf_email_config AS CEC ON CEC.id = CET.email_from ")
	Qry.WriteString(" LEFT JOIN cf_email_config AS CEC1 ON find_in_set(CEC1.id, CET.email_bcc) ")
	Qry.WriteString(" WHERE CET.id = ? GROUP BY CET.id")

	RetMap, err := ExecuteRowQuery(Qry.String(), id)
	if util.CheckErrorLog(r, err) {
		return nil, err
	}

	return RetMap, nil
}

// GetEmailList -  Return Email List
func GetEmailList(r *http.Request) ([]map[string]interface{}, error) {
	util.LogIt(r, "model - Email_Template - GetEmailList")
	var Qry bytes.Buffer

	Qry.WriteString(" SELECT * ")
	Qry.WriteString(" FROM cf_email_config")
	Qry.WriteString(" WHERE status = 1")

	RetMap, err := ExecuteQuery(Qry.String())
	if util.CheckErrorLog(r, err) {
		return nil, err
	}

	return RetMap, nil
}
