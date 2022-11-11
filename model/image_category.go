package model

import (
	"bytes"
	"net/http"
	"tp-api-common/data"
	"tp-api-common/util"

	"github.com/gorilla/context"
	gonanoid "github.com/matoous/go-nanoid"
)

// AddImageCategory - Add Image Category
func AddImageCategory(r *http.Request, reqMap data.ImageCategory) bool {
	util.LogIt(r, "model - V_Image_Category - AddImageCategory")
	var Qry bytes.Buffer
	nanoid, _ := gonanoid.Nanoid()
	Qry.WriteString("INSERT INTO cf_image_category(id,name,created_at,created_by) VALUES (?,?,?,?)")
	err := ExecuteNonQuery(Qry.String(), nanoid, reqMap.Category, util.GetIsoLocalDateTime(), context.Get(r, "UserId"))
	if util.CheckErrorLog(r, err) {
		return false
	}

	AddLog(r, "", "IMAGE_CATEGORY", "Create", nanoid, GetLogsValueMap(r, util.ToMap(reqMap), false, ""))

	return true
}

// UpdateImageCategory - Update Image Category
func UpdateImageCategory(r *http.Request, reqMap data.ImageCategory) bool {
	util.LogIt(r, "model - V_Image_Category - UpdateImageCategory")
	var Qry bytes.Buffer
	BeforeUpdate, _ := GetModuleFieldByID(r, "IMAGE_CATEGORY", reqMap.ID, "name")

	Qry.WriteString("UPDATE cf_image_category SET name = ? WHERE id = ?")
	err := ExecuteNonQuery(Qry.String(), reqMap.Category, reqMap.ID)
	if util.CheckErrorLog(r, err) {
		return false
	}

	AddLog(r, BeforeUpdate.(string), "IMAGE_CATEGORY", "Update", reqMap.ID, GetLogsValueMap(r, util.ToMap(reqMap), false, ""))

	return true
}

// ImageCategoryListing - Datatable Image Category listing with filter and order
func ImageCategoryListing(r *http.Request, reqMap data.JQueryTableUI) (map[string]interface{}, error) {
	util.LogIt(r, "model - V_Image_Category - ImageCategoryListing")
	var Qry, QryCnt, QryFilter, GroupBy, FromQry bytes.Buffer
	var status = make(map[string]interface{})

	var testColArrs [20]string
	testColArrs[0] = "CMT.id"
	testColArrs[1] = "CMT.name"
	testColArrs[2] = "CMT.status"
	testColArrs[3] = "created_by"

	var testArrs []map[string]string
	testArrs = append(testArrs, map[string]string{
		"key":   "category",
		"value": "CMT.name",
	})
	testArrs = append(testArrs, map[string]string{
		"key":   "status",
		"value": "CMT.status",
	})
	testArrs = append(testArrs, map[string]string{
		"key":   "created_at",
		"value": "DATE(from_unixtime(CMT.created_at))",
	})

	QryCnt.WriteString(" COUNT(CMT.id) AS cnt ")
	QryFilter.WriteString(" COUNT(CMT.id) AS cnt ")

	Qry.WriteString(" CMT.id,name AS category,CONCAT(from_unixtime(CMT.created_at),' ',SUC.username) AS created_by,ST.status,ST.id AS status_id ")

	FromQry.WriteString(" FROM cf_image_category AS CMT ")
	FromQry.WriteString(" INNER JOIN cf_user AS SUC ON SUC.id = CMT.created_by ")
	FromQry.WriteString(" INNER JOIN status AS ST ON ST.id = CMT.status ")
	FromQry.WriteString(" WHERE CMT.status <> 3 ")
	Data, err := JQueryTable(r, reqMap, Qry, FromQry, QryCnt, QryFilter, GroupBy, testArrs, testColArrs, status, "")
	return Data, err
}

// GetImageCategoryList - Get Image Category List For Other Module
func GetImageCategoryList(r *http.Request) ([]map[string]interface{}, error) {
	util.LogIt(r, "model - V_Image_Category - GetImageCategoryList")
	var Qry bytes.Buffer

	Qry.WriteString("SELECT id,name FROM cf_image_category WHERE status = 1")
	RetMap, err := ExecuteQuery(Qry.String())
	if util.CheckErrorLog(r, err) {
		return nil, err
	}

	return RetMap, nil
}
