package model

import (
	"bytes"
	"net/http"
	"tp-api-common/data"
	"tp-api-common/util"

	"github.com/gorilla/context"
	gonanoid "github.com/matoous/go-nanoid"
)

// AddCms - Add CMS
func AddCms(r *http.Request, reqMap data.Cms) bool {
	util.LogIt(r, "Model - Cms - AddCms")

	var Qry bytes.Buffer
	nanoid, _ := gonanoid.Nanoid()

	Qry.WriteString("INSERT INTO cf_cms (id, page, sort_code, slug, content, created_at, created_by) VALUES (?, ?, ?, ?, ?, ?, ?);")
	err := ExecuteNonQuery(Qry.String(), nanoid, reqMap.Page, reqMap.SortCode, reqMap.Slug, reqMap.Content, util.GetLocalDateTime(), context.Get(r, "UserId"))
	if util.CheckErrorLog(r, err) {
		return false
	}

	AddLog(r, "", "CMS", "Create", nanoid, GetLogsValueMap(r, util.ToMap(reqMap), true, "Content"))

	return true
}

// UpdateCms - Update Cms
func UpdateCms(r *http.Request, reqMap data.Cms) bool {
	util.LogIt(r, "Model - Cms - UpdateCms")

	var Qry bytes.Buffer
	OldData, _ := GetModuleFieldsByID(r, "CMS", reqMap.ID, "slug")

	Qry.WriteString("UPDATE cf_cms SET page=?, content=?, updated_at=?, updated_by=? WHERE id = ?")
	err := ExecuteNonQuery(Qry.String(), reqMap.Page, reqMap.Content, util.GetLocalDateTime(), context.Get(r, "UserId"), reqMap.ID)
	if util.CheckErrorLog(r, err) {
		return false
	}

	AddLog(r, OldData["slug"].(string), "CMS", "Update", reqMap.ID, GetLogsValueMap(r, util.ToMap(reqMap), true, "Content"))

	return true
}

// GetCms -  Return Cms Details
func GetCms(r *http.Request, id string) (map[string]interface{}, error) {
	util.LogIt(r, "model - Cms - GetCms")
	var Qry bytes.Buffer

	Qry.WriteString("SELECT id, page, sort_code, slug, content FROM cf_cms WHERE id = ?")
	RetMap, err := ExecuteRowQuery(Qry.String(), id)
	if util.CheckErrorLog(r, err) {
		return nil, err
	}

	return RetMap, nil
}

// CmsListing - Return Datatable Listing Of Email Template
func CmsListing(r *http.Request, reqMap data.JQueryTableUI) (map[string]interface{}, error) {
	util.LogIt(r, "model - Cms - CmsListing")
	var Qry, QryCnt, QryFilter, GroupBy, FromQry bytes.Buffer
	var status = make(map[string]interface{})

	var testColArrs [20]string
	testColArrs[0] = "CC.id"
	testColArrs[1] = "CC.page"
	testColArrs[2] = "CC.sort_code"
	testColArrs[3] = "CC.slug"
	testColArrs[4] = "created_by"

	var testArrs []map[string]string
	testArrs = append(testArrs, map[string]string{
		"key":   "page",
		"value": "CC.page",
	})
	testArrs = append(testArrs, map[string]string{
		"key":   "sort_code",
		"value": "CC.sort_code",
	})
	testArrs = append(testArrs, map[string]string{
		"key":   "slug",
		"value": "CC.slug",
	})
	testArrs = append(testArrs, map[string]string{
		"key":   "created_at",
		"value": "DATE(CC.created_at)",
	})

	QryCnt.WriteString(" COUNT(CC.id) AS cnt ")
	QryFilter.WriteString(" COUNT(CC.id) AS cnt ")

	Qry.WriteString(" CC.id, CC.page, CC.sort_code, CC.slug, CU.username, CONCAT(CC.created_at,' ',CU.username) AS created_by ")

	FromQry.WriteString(" FROM cf_cms AS CC ")
	FromQry.WriteString(" INNER JOIN cf_user AS CU ON CU.id = CC.created_by ")
	FromQry.WriteString(" WHERE 1 = 1 ")
	Data, err := JQueryTable(r, reqMap, Qry, FromQry, QryCnt, QryFilter, GroupBy, testArrs, testColArrs, status, "")
	return Data, err
}
