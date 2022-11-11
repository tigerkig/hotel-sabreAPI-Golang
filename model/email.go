package model

import (
	"bytes"
	b64 "encoding/base64"
	"net/http"
	"strconv"
	"tp-api-common/data"
	"tp-api-common/util"

	"github.com/gorilla/context"
	gonanoid "github.com/matoous/go-nanoid"
)

// AddEmailConfig - Adds Email Config
func AddEmailConfig(r *http.Request, reqMap data.EmailConfig) bool {
	util.LogIt(r, "Model - Email - AddEmailConfig")

	var Qry bytes.Buffer
	nanoid, _ := gonanoid.Nanoid()

	EncPasswd := b64.StdEncoding.EncodeToString([]byte(reqMap.SMTPPassword))

	Qry.WriteString(" INSERT INTO cf_email_config(id,email_id,email_name,smtp_host,smtp_port,smtp_user,smtp_password,signature,created_at,created_by) ")
	Qry.WriteString(" VALUES (?,?,?,?,?,?,?,?,?,?) ")
	err := ExecuteNonQuery(Qry.String(), nanoid, reqMap.EmailID, reqMap.EmailName, reqMap.SMTPHost, reqMap.SMTPPort, reqMap.SMTPUser, EncPasswd, reqMap.Signature, util.GetIsoLocalDateTime(), context.Get(r, "UserId"))
	if util.CheckErrorLog(r, err) {
		return false
	}

	AddLog(r, "", "EMAIL_CONFIG", "Create", nanoid, GetLogsValueMap(r, util.ToMap(reqMap), true, "SMTPUser,SMTPPassword,Signature"))

	return true
}

// UpdateEmailConfig - Updates Email Configuration
func UpdateEmailConfig(r *http.Request, reqMap data.EmailConfig) bool {

	util.LogIt(r, "Model - Email - UpdateEmailConfig")

	var Qry bytes.Buffer
	OldData, _ := GetModuleFieldsByID(r, "EMAIL_CONFIG", reqMap.ID, "email_id,email_name,smtp_host,smtp_port,smtp_user,smtp_password,signature")

	EncPasswd := b64.StdEncoding.EncodeToString([]byte(reqMap.SMTPPassword))
	Qry.WriteString("UPDATE cf_email_config SET email_id=?, email_name=?, smtp_host=?, smtp_port=?, smtp_user=?, smtp_password=?, signature=? WHERE id = ?")
	err := ExecuteNonQuery(Qry.String(), reqMap.EmailID, reqMap.EmailName, reqMap.SMTPHost, reqMap.SMTPPort, reqMap.SMTPUser, EncPasswd, reqMap.Signature, reqMap.ID)
	if util.CheckErrorLog(r, err) {
		return false
	}

	if OldData["email_id"].(string) != reqMap.EmailID {
		AddLog(r, OldData["email_id"].(string), "EMAIL_CONFIG", "Update", reqMap.ID, GetLogsValueMap(r, util.ToMap(reqMap), true, "ID,EmailName,SMTPHost,SMTPPort,SMTPUser,SMTPPassword,Signature"))
	}

	if OldData["email_name"].(string) != reqMap.EmailName {
		AddLog(r, OldData["email_name"].(string), "EMAIL_CONFIG", "Update", reqMap.ID, GetLogsValueMap(r, util.ToMap(reqMap), true, "ID,EmailID,SMTPHost,SMTPPort,SMTPUser,SMTPPassword,Signature"))
	}

	if OldData["smtp_host"].(string) != reqMap.SMTPHost {
		AddLog(r, OldData["smtp_host"].(string), "EMAIL_CONFIG", "Update", reqMap.ID, GetLogsValueMap(r, util.ToMap(reqMap), true, "ID,EmailID,EmailName,SMTPPort,SMTPUser,SMTPPassword,Signature"))
	}

	if OldData["smtp_port"].(int64) != int64(reqMap.SMTPPort) {
		AddLog(r, strconv.FormatInt(OldData["smtp_port"].(int64), 10), "EMAIL_CONFIG", "Update", reqMap.ID, GetLogsValueMap(r, util.ToMap(reqMap), true, "ID,EmailID,EmailName,SMTPHost,SMTPUser,SMTPPassword,Signature"))
	}

	if OldData["smtp_user"].(string) != reqMap.SMTPUser {
		AddLog(r, OldData["smtp_user"].(string), "EMAIL_CONFIG", "Update", reqMap.ID, GetLogsValueMap(r, util.ToMap(reqMap), true, "ID,EmailID,EmailName,SMTPPort,SMTPUser,SMTPPassword,Signature"))
	}

	DecOldPasswd, _ := b64.StdEncoding.DecodeString(OldData["smtp_password"].(string))
	DecOldPasswdCmp := string(DecOldPasswd)
	if DecOldPasswdCmp != reqMap.SMTPPassword {
		AddLog(r, OldData["smtp_password"].(string), "EMAIL_CONFIG", "Update", reqMap.ID, GetLogsValueMap(r, util.ToMap(reqMap), true, "ID,EmailID,EmailName,SMTPHost,SMTPPort,SMTPUser,Signature"))
	}

	if OldData["signature"].(string) != reqMap.Signature {
		AddLog(r, OldData["signature"].(string), "EMAIL_CONFIG", "Update", reqMap.ID, GetLogsValueMap(r, util.ToMap(reqMap), true, "ID,EmailID,EmailName,SMTPHost,SMTPPort,SMTPUser,SMTPPassword"))
	}

	return true
}

// GetEmailConfigInfo - Gets Email Config Info
func GetEmailConfigInfo(r *http.Request, id string) (map[string]interface{}, error) {

	util.LogIt(r, "Model - Email - GetEmailConfigInfo")

	var Qry bytes.Buffer
	Qry.WriteString("SELECT id,email_id,email_name,smtp_host,smtp_port,smtp_user,smtp_password,signature FROM cf_email_config WHERE id = ?")
	RetMap, err := ExecuteRowQuery(Qry.String(), id)
	if util.CheckErrorLog(r, err) {
		return nil, err
	}
	DecPasswd, _ := b64.StdEncoding.DecodeString(RetMap["smtp_password"].(string))
	DecPasswdVal := string(DecPasswd)
	RetMap["smtp_password"] = DecPasswdVal

	return RetMap, nil
}

// EmailListing - Get Meal Type Listing
func EmailListing(r *http.Request, reqMap data.JQueryTableUI) (map[string]interface{}, error) {
	util.LogIt(r, "Model - Email - EmailListing")

	var Qry, QryCnt, QryFilter, GroupBy, FromQry bytes.Buffer

	var status = make(map[string]interface{})

	var testColArrs [20]string
	testColArrs[0] = "CEC.id"
	testColArrs[1] = "CEC.email_id"
	testColArrs[2] = "CEC.status"
	testColArrs[3] = "created_by"

	var testArrs []map[string]string
	testArrs = append(testArrs, map[string]string{
		"key":   "email",
		"value": "CEC.email_id",
	})
	testArrs = append(testArrs, map[string]string{
		"key":   "status",
		"value": "CEC.status",
	})
	testArrs = append(testArrs, map[string]string{
		"key":   "created_at",
		"value": "DATE(from_unixtime(CEC.created_at))",
	})

	QryCnt.WriteString(" COUNT(CEC.id) AS cnt ")
	QryFilter.WriteString(" COUNT(CEC.id) AS cnt ")

	Qry.WriteString(" CEC.id, CEC.email_id, CONCAT(from_unixtime(CEC.created_at),' ',CU.username) AS created_by, ST.status, ST.id AS status_id ")

	FromQry.WriteString(" FROM cf_email_config AS CEC ")
	FromQry.WriteString(" INNER JOIN cf_user AS CU ON CU.id = CEC.created_by  ")
	FromQry.WriteString(" INNER JOIN status AS ST ON ST.id = CEC.status ")
	FromQry.WriteString(" WHERE CEC.status <> 3 ")
	Data, err := JQueryTable(r, reqMap, Qry, FromQry, QryCnt, QryFilter, GroupBy, testArrs, testColArrs, status, "")
	return Data, err
}
