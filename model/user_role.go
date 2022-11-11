package model

import (
	"bytes"
	"net/http"
	"tp-api-common/data"
	"tp-api-common/util"

	"github.com/gorilla/context"
	gonanoid "github.com/matoous/go-nanoid"
)

// GetPrivilegeList - Get List Of privilege
func GetPrivilegeList(r *http.Request) ([]map[string]interface{}, error) {
	util.LogIt(r, "model - V_User_Role - GetPrivilegeList")
	var Qry bytes.Buffer
	Qry.WriteString("SELECT id, module FROM cf_privilege_module")
	AuthMap, err := ExecuteQuery(Qry.String())
	if util.CheckErrorLog(r, err) {
		return nil, err
	}
	if len(AuthMap) > 0 {
		for i := 0; i < len(AuthMap); i++ {
			var Qry bytes.Buffer
			Qry.WriteString("SELECT id, privilege FROM cf_privilege WHERE module_id = ? ")
			ipMap, err := ExecuteQuery(Qry.String(), AuthMap[i]["id"])
			if util.CheckErrorLog(r, err) {
				return nil, err
			}
			AuthMap[i]["privilege"] = ipMap
		}
	}

	return AuthMap, nil
}

// AddUserRole - Add User Role
func AddUserRole(r *http.Request, role string, privileges string) bool {
	util.LogIt(r, "model - V_User_Role - AddUserRole")
	var Qry bytes.Buffer
	nanoid, _ := gonanoid.Nanoid()
	Qry.WriteString("INSERT INTO cf_user_role(id, role, privileges,created_at,created_by) VALUES (?,?,?,?,?)")
	err := ExecuteNonQuery(Qry.String(), nanoid, role, privileges, util.GetIsoLocalDateTime(), context.Get(r, "UserId"))
	if util.CheckErrorLog(r, err) {
		return false
	}

	AddLog(r, "", "ROLE", "Create", nanoid, map[string]interface{}{"Role": role})

	return true
}

// UpdateUserRole - Update User Role
func UpdateUserRole(r *http.Request, id string, role string, privileges string) bool {
	util.LogIt(r, "model - V_User_Role - UpdateUserRole")
	var Qry bytes.Buffer

	RoleName, _ := GetModuleFieldByID(r, "ROLE", id, "role")

	Qry.WriteString("UPDATE cf_user_role SET role=?, privileges=? WHERE id=?")
	err := ExecuteNonQuery(Qry.String(), role, privileges, id)
	if util.CheckErrorLog(r, err) {
		return false
	}

	AddLog(r, RoleName.(string), "ROLE", "Permission Update", id, map[string]interface{}{"Role": role})

	return true
}

// GetUserRole - Get User Role Detail By ID
func GetUserRole(r *http.Request, id string) (map[string]interface{}, error) {
	util.LogIt(r, "model - V_User_Role - GetUserRole")
	var Qry bytes.Buffer

	Qry.WriteString("SELECT id,role,privileges FROM cf_user_role WHERE id = ?")
	RetMap, err := ExecuteRowQuery(Qry.String(), id)
	if util.CheckErrorLog(r, err) {
		return nil, err
	}

	return RetMap, nil
}

// GetUserRoleList - Get User Role List For User Module
func GetUserRoleList(r *http.Request) ([]map[string]interface{}, error) {
	util.LogIt(r, "model - V_User_Role - GetUserRoleList")
	var Qry bytes.Buffer

	Qry.WriteString("SELECT id,role,privileges FROM cf_user_role WHERE status = 1")
	RetMap, err := ExecuteQuery(Qry.String())
	if util.CheckErrorLog(r, err) {
		return nil, err
	}

	return RetMap, nil
}

// UserRoleListing - Get User Role Listing
func UserRoleListing(r *http.Request, reqMap data.JQueryTableUI) (map[string]interface{}, error) {
	util.LogIt(r, "model - V_User_Role - UserRoleListing")
	var Qry, QryCnt, QryFilter, GroupBy, FromQry bytes.Buffer
	var status = make(map[string]interface{})

	var testColArrs [20]string
	testColArrs[0] = "SUR.id"
	testColArrs[1] = "role"
	testColArrs[2] = "SUR.status"
	testColArrs[3] = "created_by"

	var testArrs []map[string]string
	testArrs = append(testArrs, map[string]string{
		"key":   "role",
		"value": "SUR.role",
	})
	testArrs = append(testArrs, map[string]string{
		"key":   "status",
		"value": "SUR.status",
	})
	testArrs = append(testArrs, map[string]string{
		"key":   "created_at",
		"value": "DATE(from_unixtime(SUR.created_at))",
	})

	QryCnt.WriteString(" COUNT(SUR.id) AS cnt ")
	QryFilter.WriteString(" COUNT(SUR.id) AS cnt ")

	Qry.WriteString(" SUR.id,SUR.role, CONCAT(from_unixtime(SUR.created_at),' ',SUC.username) AS created_by, SUR.status AS status_id, ST.status ")

	FromQry.WriteString(" FROM cf_user_role AS SUR ")
	FromQry.WriteString(" INNER JOIN cf_user AS SUC ON SUC.id = SUR.created_by ")
	FromQry.WriteString(" INNER JOIN status AS ST ON ST.id = SUR.status ")
	FromQry.WriteString(" WHERE SUR.status != 3 ")
	Data, err := JQueryTable(r, reqMap, Qry, FromQry, QryCnt, QryFilter, GroupBy, testArrs, testColArrs, status, "")
	return Data, err
}

// OnflyAddUserRole - Add User Role
func OnflyAddUserRole(r *http.Request, role string, privileges string) (map[string]interface{}, error) {
	util.LogIt(r, "model - V_User_Role - OnflyAddUserRole")
	nanoid, _ := gonanoid.Nanoid()

	var Qry bytes.Buffer
	Qry.WriteString("INSERT INTO cf_user_role(id, role, privileges,created_at,created_by) VALUES (?,?,?,?,?)")
	err := ExecuteNonQuery(Qry.String(), nanoid, role, privileges, util.GetIsoLocalDateTime(), context.Get(r, "UserId"))
	if util.CheckErrorLog(r, err) {
		return nil, err
	}

	AddLog(r, "", "ROLE", "Create", nanoid, map[string]interface{}{"Role": role})

	var retMap = make(map[string]interface{})
	retMap["id"] = nanoid
	return retMap, nil
}
