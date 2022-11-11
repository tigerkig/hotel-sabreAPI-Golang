package model

import (
	"bytes"
	"net/http"
	"tp-api-common/data"
	"tp-api-common/util"

	"github.com/gorilla/context"
	gonanoid "github.com/matoous/go-nanoid"
)

// AddRoomView - Add Room View
func AddRoomView(r *http.Request, reqMap data.RoomView) bool {
	util.LogIt(r, "model - V_Room_View - AddRoomView")
	var Qry bytes.Buffer
	nanoid, _ := gonanoid.Nanoid()
	Qry.WriteString("INSERT INTO cf_room_view(id,room_view_name,created_at,created_by) VALUES (?,?,?,?)")
	err := ExecuteNonQuery(Qry.String(), nanoid, reqMap.RoomView, util.GetIsoLocalDateTime(), context.Get(r, "UserId"))
	if util.CheckErrorLog(r, err) {
		return false
	}

	AddLog(r, "", "ROOM_VIEW", "Create", nanoid, map[string]interface{}{"View Name": reqMap.RoomView})

	return true
}

// UpdateRoomView - Update Room View
func UpdateRoomView(r *http.Request, reqMap data.RoomView) bool {
	util.LogIt(r, "model - V_Room_View - UpdateRoomView")
	var Qry bytes.Buffer

	BeforeUpdate, _ := GetModuleFieldByID(r, "ROOM_VIEW", reqMap.ID, "room_view_name")

	Qry.WriteString("UPDATE cf_room_view SET room_view_name = ? WHERE id = ?")
	err := ExecuteNonQuery(Qry.String(), reqMap.RoomView, reqMap.ID)
	if util.CheckErrorLog(r, err) {
		return false
	}

	AddLog(r, BeforeUpdate.(string), "ROOM_VIEW", "Update", reqMap.ID, map[string]interface{}{"View Name": reqMap.RoomView})

	return true
}

// RoomViewListing - Datatable Room View listing with filter and order
func RoomViewListing(r *http.Request, reqMap data.JQueryTableUI) (map[string]interface{}, error) {
	util.LogIt(r, "model - V_Room_View - RoomViewListing")
	var Qry, QryCnt, QryFilter, GroupBy, FromQry bytes.Buffer
	var status = make(map[string]interface{})

	var testColArrs [20]string
	testColArrs[0] = "CRV.id"
	testColArrs[1] = "room_view_name"
	testColArrs[2] = "CRV.status"
	testColArrs[3] = "created_by"

	var testArrs []map[string]string
	testArrs = append(testArrs, map[string]string{
		"key":   "room_view_name",
		"value": "CRV.room_view_name",
	})
	testArrs = append(testArrs, map[string]string{
		"key":   "status",
		"value": "CRV.status",
	})
	testArrs = append(testArrs, map[string]string{
		"key":   "created_at",
		"value": "DATE(from_unixtime(CRV.created_at))",
	})

	QryCnt.WriteString(" COUNT(CRV.id) AS cnt ")
	QryFilter.WriteString(" COUNT(CRV.id) AS cnt ")

	Qry.WriteString(" CRV.id, room_view_name, CONCAT(from_unixtime(CRV.created_at),' ',SUC.username) AS created_by,ST.status,ST.id AS status_id ")

	FromQry.WriteString(" FROM cf_room_view AS CRV ")
	FromQry.WriteString(" INNER JOIN cf_user AS SUC ON SUC.id = CRV.created_by ")
	FromQry.WriteString(" INNER JOIN status AS ST ON ST.id = CRV.status ")
	FromQry.WriteString(" WHERE CRV.status <> 3 ")
	Data, err := JQueryTable(r, reqMap, Qry, FromQry, QryCnt, QryFilter, GroupBy, testArrs, testColArrs, status, "")
	return Data, err
}

// GetRoomViewList - Get Room View List For Other Module
func GetRoomViewList(r *http.Request) ([]map[string]interface{}, error) {
	util.LogIt(r, "model - V_Room_View - GetRoomViewList")
	var Qry bytes.Buffer

	Qry.WriteString("SELECT id,room_view_name FROM cf_room_view WHERE status = 1")
	RetMap, err := ExecuteQuery(Qry.String())
	if util.CheckErrorLog(r, err) {
		return nil, err
	}

	return RetMap, nil
}

// GetRoomView -  Get Room View Detail By ID - 2021-04-21 - HK
func GetRoomView(r *http.Request, id string) (map[string]interface{}, error) {
	util.LogIt(r, "Model - V_Room_View - GetRoomView")
	var Qry bytes.Buffer

	Qry.WriteString("SELECT id, room_view_name FROM cf_room_view WHERE id = ?")
	RetMap, err := ExecuteRowQuery(Qry.String(), id)
	if util.CheckErrorLog(r, err) {
		return nil, err
	}

	return RetMap, nil
}
