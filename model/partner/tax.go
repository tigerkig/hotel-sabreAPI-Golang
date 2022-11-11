package partner

import (
	"bytes"
	"fmt"
	"net/http"
	"tp-api-common/data"
	"tp-api-common/util"
	"tp-system/model"

	"github.com/gorilla/context"
	gonanoid "github.com/matoous/go-nanoid"
)

// AddTax - Adds Tax For Hotel Into Table
func AddTax(r *http.Request, reqMap data.Tax) bool {
	util.LogIt(r, "Model - V_Tax - AddTax")
	var Qry bytes.Buffer

	// Creates Nano ID
	mainTaxID, _ := gonanoid.Nanoid()

	// Insert Tax Info Into Table
	Qry.WriteString("INSERT INTO cf_hotel_tax(id, tax, description, type, amount, hotel_id, created_at, created_by) VALUES (?,?,?,?,?,?,?,?)")
	err := model.ExecuteNonQuery(Qry.String(), mainTaxID, reqMap.Name, reqMap.Description, reqMap.Type, reqMap.Amount, reqMap.HotelID, util.GetIsoLocalDateTime(), context.Get(r, "UserId"))
	if util.CheckErrorLog(r, err) {
		return false
	}

	// Insert Log Of Added Tax Into MongoDB
	model.AddLog(r, "", "TAX", "Create", mainTaxID, model.GetLogsValueMap(r, util.ToMap(reqMap), false, ""))

	model.CacheChn <- model.CacheObj{
		Type: "tax",
		// ID:         context.Get(r, "HotelId").(string),
		ID:         reqMap.HotelID,
		Additional: "",
	}

	return true
}

// UpdateTax - Updates Tax For Hotel
func UpdateTax(r *http.Request, reqMap data.Tax) bool {
	util.LogIt(r, "Model - V_Tax - UpdateTax")

	// Retrieves Old Values Of Passed Parameter For Log Purpose
	NameBeforeUpdate, _ := model.GetModuleFieldByID(r, "TAX", reqMap.ID, "tax")
	DescBeforeUpdate, _ := model.GetModuleFieldByID(r, "TAX", reqMap.ID, "description")
	TypeBeforeUpdate, _ := model.GetModuleFieldByID(r, "TAX", reqMap.ID, "type")
	AmountBeforeUpdate, _ := model.GetModuleFieldByID(r, "TAX", reqMap.ID, "amount")

	var Qry bytes.Buffer
	Qry.WriteString("UPDATE cf_hotel_tax SET tax = ?, description = ?, type = ?, amount = ?  WHERE id = ? AND hotel_id = ?")
	err := model.ExecuteNonQuery(Qry.String(), reqMap.Name, reqMap.Description, reqMap.Type, reqMap.Amount, reqMap.ID, reqMap.HotelID)
	if util.CheckErrorLog(r, err) {
		return false
	}

	// Converts request Map Stuct To Map
	reqStruct := util.ToMap(reqMap)

	if NameBeforeUpdate.(string) != reqMap.Name {
		model.AddLog(r, NameBeforeUpdate.(string), "TAX", "Update", reqMap.ID, model.GetLogsValueMap(r, reqStruct, true, "ID,Type,Amount,Description"))
	}

	if DescBeforeUpdate.(string) != reqMap.Description {
		model.AddLog(r, NameBeforeUpdate.(string), "TAX", "Update", reqMap.ID, model.GetLogsValueMap(r, reqStruct, true, "ID,Name,Type,Amount"))
	}

	if TypeBeforeUpdate.(string) != reqMap.Type {
		model.AddLog(r, NameBeforeUpdate.(string), "TAX", "Update", reqMap.ID, model.GetLogsValueMap(r, reqStruct, true, "ID,Name,Amount,Description"))
	}

	AmountAfterUpdate := fmt.Sprintf("%.2f", reqMap.Amount)
	if AmountBeforeUpdate != AmountAfterUpdate {
		model.AddLog(r, NameBeforeUpdate.(string), "TAX", "Update", reqMap.ID, model.GetLogsValueMap(r, reqStruct, true, "ID,Type,Description,Name"))
	}

	model.CacheChn <- model.CacheObj{
		Type: "tax",
		// ID:         context.Get(r, "HotelId").(string),
		ID:         reqMap.HotelID,
		Additional: "",
	}

	return true
}

// GetTaxInfo - Get Tax Info
func GetTaxInfo(r *http.Request, TaxID string, HotelID string) (map[string]interface{}, error) {
	util.LogIt(r, "Model - V_Tax - GetTaxInfo")

	var Qry bytes.Buffer
	Qry.WriteString(" SELECT ")
	Qry.WriteString(" CHT.id, CHT.tax, CHT.description, CHT.type, CHT.amount ")
	Qry.WriteString(" FROM cf_hotel_tax AS CHT ")
	Qry.WriteString(" WHERE CHT.id = ? AND CHT.hotel_id = ? ")

	RetMap, err := model.ExecuteRowQuery(Qry.String(), TaxID, HotelID)
	if util.CheckErrorLog(r, err) {
		return nil, err
	}
	return RetMap, nil
}

// TaxListing - Return Datatable Listing Of Tax
func TaxListing(r *http.Request, reqMap data.JQueryTableUI) (map[string]interface{}, error) {

	util.LogIt(r, "Model - V_Tax - TaxListing")

	var Qry, QryCnt, QryFilter, GroupBy, FromQry bytes.Buffer
	var status = make(map[string]interface{})

	var testColArrs [20]string
	testColArrs[0] = "CHT.id"
	testColArrs[1] = "CHT.tax"
	testColArrs[2] = "ST.status"
	testColArrs[3] = "CHT.description"
	testColArrs[4] = "type_id"
	testColArrs[5] = "amount"
	testColArrs[6] = "created_by"

	var testArrs []map[string]string
	testArrs = append(testArrs, map[string]string{
		"key":   "tax",
		"value": "CHT.tax",
	})
	testArrs = append(testArrs, map[string]string{
		"key":   "description",
		"value": "CHT.description",
	})
	testArrs = append(testArrs, map[string]string{
		"key":   "status",
		"value": "CHT.status",
	})
	testArrs = append(testArrs, map[string]string{
		"key":   "created_at",
		"value": "DATE(from_unixtime(CHT.created_at))",
	})

	QryCnt.WriteString(" COUNT(CHT.id) AS cnt ")

	QryFilter.WriteString(" SELECT COUNT(tbl.cnt) FROM (")
	QryFilter.WriteString(" COUNT(CHT.id) AS cnt ")

	Qry.WriteString(" CHT.id, CHT.tax, CHT.description, CHT.type, CHT.amount, ST.status, CONCAT(from_unixtime(CHT.created_at),' ',CHC.username) AS created_by, ST.id AS status_id ")

	FromQry.WriteString(" FROM cf_hotel_tax AS CHT ")

	FromQry.WriteString(" INNER JOIN cf_hotel_client AS CHC ON CHC.id = CHT.created_by ")
	FromQry.WriteString(" INNER JOIN status AS ST ON ST.id = CHT.status ")
	FromQry.WriteString(" WHERE ST.id <> 3 AND CHT.hotel_id = '")
	// FromQry.WriteString(context.Get(r, "HotelId").(string))
	FromQry.WriteString(reqMap.HotelID)
	FromQry.WriteString("'")

	Data, err := model.JQueryTable(r, reqMap, Qry, FromQry, QryCnt, QryFilter, GroupBy, testArrs, testColArrs, status, "")
	return Data, err
}
