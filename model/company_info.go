package model

import (
	"bytes"
	"net/http"
	"tp-api-common/util"
	"tp-system/config"
)

// GetCompanyInfo -  Get Popular City Detail By ID
func GetCompanyInfo(r *http.Request, id string) (map[string]interface{}, error) {
	util.LogIt(r, "Model - Company - GetCompanyInfo")

	var RetMap = make(map[string]interface{})
	var err error

	var Qry bytes.Buffer
	Qry.WriteString(" SELECT ")
	Qry.WriteString(" CCI.id, CCI.company_name, CCI.zip_code, CCI.address, CCI.registered_office_address, CONCAT('" + config.Env.AwsBucketURL + "company_logo/" + "',image) AS image, ")
	Qry.WriteString(" CC.name as city_name, CCI.city_id, ")
	Qry.WriteString(" CST.name as state_name, CCI.state_id, ")
	Qry.WriteString(" CCN.name as country_name, CCI.country_id ")
	Qry.WriteString(" FROM ")
	Qry.WriteString(" cf_company_info AS CCI ")
	Qry.WriteString(" LEFT JOIN ")
	Qry.WriteString(" cf_city AS CC ON CC.id = CCI.city_id ")
	Qry.WriteString(" LEFT JOIN ")
	Qry.WriteString(" cf_states AS CST ON CST.id = CCI.state_id ")
	Qry.WriteString(" LEFT JOIN ")
	Qry.WriteString(" cf_country AS CCN ON CCN.id = CCI.country_id ")
	if id != "" {
		Qry.WriteString(" WHERE CCI.id = ?")
	}

	if id != "" {
		RetMap, err = ExecuteRowQuery(Qry.String(), id)
	} else {
		RetMap, err = ExecuteRowQuery(Qry.String())
	}

	if util.CheckErrorLog(r, err) {
		return nil, err
	}

	return RetMap, nil
}

// UpdateCompanyInfo - Update Company Info
func UpdateCompanyInfo(r *http.Request, reqMap map[string]interface{}) bool {
	util.LogIt(r, "Model - Company - UpdateCompanyInfo")

	var Qry bytes.Buffer
	Qry.WriteString("UPDATE cf_company_info SET company_name=?, country_id=?, state_id=?, city_id=?, zip_code=?, address=?, registered_office_address=?, image=? WHERE id=?")
	err := ExecuteNonQuery(Qry.String(), reqMap["company_name"], reqMap["country_id"], reqMap["state_id"], reqMap["city_id"], reqMap["zip_code"], reqMap["address"], reqMap["registered_office_address"], reqMap["image"], reqMap["id"])
	if util.CheckErrorLog(r, err) {
		return false
	}
	return true
}
