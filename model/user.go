package model

import (
	"bytes"
	"net/http"
	"tp-api-common/data"
	"tp-api-common/util"

	"github.com/gorilla/context"
	gonanoid "github.com/matoous/go-nanoid"
)

//User - Module Name
var User = "USER"

func ResetPassword(r *http.Request, reqMap data.ResetPassword) bool {
	util.LogIt(r, "model - V_User - ResetPassword")
	var SQry bytes.Buffer
	var password = util.GeneratePasswordHash(reqMap.NewPassword)
	switch context.Get(r, "Side") {
	case "TP-BACKOFFICE":
		SQry.WriteString("UPDATE cf_user SET password=? WHERE id = ?")
	case "TP-PARTNER":
		SQry.WriteString("UPDATE cf_hotel_client SET password=? WHERE id = ?")
	}
	err := ExecuteNonQuery(SQry.String(), password, context.Get(r, "UserId"))
	if util.CheckErrorLog(r, err) {
		return false
	}

	return true
}

func CheckOldPassword(r *http.Request, reqMap data.ResetPassword) bool {
	util.LogIt(r, "model - V_User - ResetPassword")
	var SQry bytes.Buffer
	var password = util.GeneratePasswordHash(reqMap.OldPassword)
	switch context.Get(r, "Side") {
	case "TP-BACKOFFICE":
		SQry.WriteString("SELECT id FROM cf_user WHERE password=? AND id = ?")
	case "TP-PARTNER":
		SQry.WriteString("SELECT id FROM cf_hotel_client WHERE password=? AND id = ?")
	}
	Cnt, err := ExecuteQuery(SQry.String(), password, context.Get(r, "UserId"))
	if util.CheckErrorLog(r, err) {
		return false
	}

	if len(Cnt) == 0 {
		return false
	}
	return true
}

// AddUser - Add User
func AddUser(r *http.Request, reqMap data.User) bool {
	util.LogIt(r, "model - V_User - AddUser")
	var Qry, SQry bytes.Buffer
	nanoid, _ := gonanoid.Nanoid()
	nanoidSys, _ := gonanoid.Nanoid()

	var password = util.GeneratePasswordHash(reqMap.Password)
	Qry.WriteString("INSERT INTO cf_user(id, username, password,user_role_id,privileges,created_at,created_by) VALUES (?,?,?,?,?,?,?)")
	err := ExecuteNonQuery(Qry.String(), nanoid, reqMap.Username, password, reqMap.Role, reqMap.Privileges, util.GetIsoLocalDateTime(), context.Get(r, "UserId"))
	if util.CheckErrorLog(r, err) {
		return false
	}

	SQry.WriteString("INSERT INTO cf_user_profile(id, user_id, name, phone, phone_code, email, cc, description, birthdate, address,created_at,created_by) VALUES (?,?,?,?,?,?,?,?,?,?,?,?)")
	err = ExecuteNonQuery(SQry.String(), nanoidSys, nanoid, reqMap.Name, reqMap.Phone, reqMap.PhoneCode, reqMap.Email, reqMap.CC, reqMap.Description, reqMap.BirthDate, reqMap.Address, util.GetIsoLocalDateTime(), context.Get(r, "UserId"))
	if util.CheckErrorLog(r, err) {
		return false
	}

	RoleName, _ := GetModuleFieldByID(r, "ROLE", reqMap.Role, "role")
	reqMap.Role = RoleName.(string)

	AddLog(r, "", User, "Create", nanoid, GetLogsValueMap(r, util.ToMap(reqMap), true, "Privileges,Password"))

	return true
}

// UpdateUser - Update User
func UpdateUser(r *http.Request, reqMap data.User) bool {
	util.LogIt(r, "model - V_User - UpdateUser")
	var SQry bytes.Buffer

	SQry.WriteString("UPDATE cf_user_profile SET name=?, phone=?, phone_code=?, email=?, description=?, birthdate=?,address=?, cc=? WHERE user_id = ?")
	err := ExecuteNonQuery(SQry.String(), reqMap.Name, reqMap.Phone, reqMap.PhoneCode, reqMap.Email, reqMap.Description, reqMap.BirthDate, reqMap.Address, reqMap.CC, reqMap.ID)
	if util.CheckErrorLog(r, err) {
		return false
	}

	// RoleName, _ := GetModuleFieldByID(r, "ROLE", reqMap.Role, "role")
	// reqMap.Role = RoleName.(string)

	AddLog(r, reqMap.Username, User, "Update", reqMap.ID, GetLogsValueMap(r, util.ToMap(reqMap), true, "Privileges"))

	return true
}

// GetUserInfo - Get User Detail By ID
func GetUserInfo(r *http.Request, id string) (map[string]interface{}, error) {
	util.LogIt(r, "model - V_User - GetUserInfo")
	var Qry bytes.Buffer
	Qry.WriteString(" SELECT ")
	Qry.WriteString(" SU.id, SUP.name, SU.user_role_id AS role, SUP.email, SUP.cc, SUP.phone AS phone, SUP.phone_code, SU.username, SUP.description, SU.status, SUP.birthdate, SU.privileges, SUP.address ")
	Qry.WriteString(" FROM cf_user_profile AS SUP ")
	Qry.WriteString(" INNER JOIN cf_user AS SU ON SU.id = SUP.user_id ")
	Qry.WriteString(" WHERE SU.status != 3 AND SU.id = ?")
	RetMap, err := ExecuteRowQuery(Qry.String(), id)
	if util.CheckErrorLog(r, err) {
		return nil, err
	}

	return RetMap, nil
}

// UserListing - Get User Listing
func UserListing(r *http.Request, reqMap data.JQueryTableUI) (map[string]interface{}, error) {
	util.LogIt(r, "model - V_User - UserListing")
	var Qry, QryCnt, QryFilter, GroupBy, FromQry bytes.Buffer
	var status = make(map[string]interface{})

	var testColArrs [20]string
	testColArrs[0] = "CU.id"
	testColArrs[1] = "name"
	testColArrs[2] = "username"
	testColArrs[3] = "phone"
	testColArrs[4] = "email"
	testColArrs[5] = "role"
	testColArrs[6] = "CU.status"
	testColArrs[7] = "created_by"

	var testArrs []map[string]string
	testArrs = append(testArrs, map[string]string{
		"key":   "name",
		"value": "CUP.name",
	})
	testArrs = append(testArrs, map[string]string{
		"key":   "username",
		"value": "CU.username",
	})
	testArrs = append(testArrs, map[string]string{
		"key":   "email",
		"value": "email",
	})
	testArrs = append(testArrs, map[string]string{
		"key":   "phone",
		"value": "CONCAT(phone_code,'',phone)",
	})
	testArrs = append(testArrs, map[string]string{
		"key":   "role",
		"value": "CUR.id",
	})
	testArrs = append(testArrs, map[string]string{
		"key":   "status",
		"value": "CU.status",
	})
	testArrs = append(testArrs, map[string]string{
		"key":   "created_at",
		"value": "DATE(from_unixtime(CU.created_at))",
	})

	QryCnt.WriteString(" COUNT(CU.id) AS cnt ")
	QryFilter.WriteString(" COUNT(CU.id) AS cnt ")

	Qry.WriteString(" CU.privileges, CU.id,CUP.name, CU.username,CUP.email,CONCAT(phone_code,'',phone) AS phone, CUR.role, ST.status, CONCAT(from_unixtime(CU.created_at),' ',CUC.username) AS created_by,ST.id AS status_id ")

	FromQry.WriteString(" FROM cf_user AS CU ")
	FromQry.WriteString(" INNER JOIN cf_user_profile AS CUP ON CUP.user_id = CU.id ")
	FromQry.WriteString(" INNER JOIN cf_user_role AS CUR ON CUR.id = CU.user_role_id ")
	FromQry.WriteString(" INNER JOIN cf_user AS CUC ON CUC.id = CU.created_by  ")
	FromQry.WriteString(" INNER JOIN status AS ST ON ST.id = CU.status WHERE CU.status != 3 ")
	Data, err := JQueryTable(r, reqMap, Qry, FromQry, QryCnt, QryFilter, GroupBy, testArrs, testColArrs, status, "")
	return Data, err
}

// ResetPwdAndPrivileges - Reset Password And Privileges
func ResetPwdAndPrivileges(r *http.Request, reqMap data.ResetPwdAndPrivileges) bool {
	util.LogIt(r, "model - V_User - ResetPwdAndPrivileges")

	var SQry bytes.Buffer

	if reqMap.Password != "" {
		var password = util.GeneratePasswordHash(reqMap.Password)
		SQry.WriteString("UPDATE cf_user SET password=?, privileges=?, user_role_id=? WHERE id = ?")
		err := ExecuteNonQuery(SQry.String(), password, reqMap.Privileges, reqMap.Role, reqMap.ID)
		if util.CheckErrorLog(r, err) {
			return false
		}
		// AddLog(r, reqMap.Username, User, "Update", reqMap.ID, GetLogsValueMap(r, util.ToMap(reqMap), true, "Role,Privileges"))
	} else {
		SQry.WriteString("UPDATE cf_user SET privileges=?, user_role_id=? WHERE id = ?")
		err := ExecuteNonQuery(SQry.String(), reqMap.Privileges, reqMap.Role, reqMap.ID)
		if util.CheckErrorLog(r, err) {
			return false
		}
		RoleName, _ := GetModuleFieldByID(r, "ROLE", reqMap.Role, "role")
		reqMap.Role = RoleName.(string)
		AddLog(r, reqMap.Username, User, "Update", reqMap.ID, GetLogsValueMap(r, util.ToMap(reqMap), true, "Privileges"))
	}

	return true
}

// GetUserList - Get User List For Another Module
func GetUserList(r *http.Request) ([]map[string]interface{}, error) {
	util.LogIt(r, "model - V_User - GetUserList")
	var Qry bytes.Buffer

	Qry.WriteString("SELECT id,username FROM cf_user WHERE status = 1")
	RetMap, err := ExecuteQuery(Qry.String())
	if util.CheckErrorLog(r, err) {
		return nil, err
	}

	return RetMap, nil
}
