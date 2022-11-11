package model

import (
	"bytes"
	"net/http"
	"tp-api-common/util"
)

// GetWebDefaultSettings - Returns settings front
func GetWebDefaultSettings(r *http.Request) (map[string]interface{}, error) {
	util.LogIt(r, "Model - web_setting - GetWebDefaultSettings")

	var Qry bytes.Buffer
	Qry.WriteString("SELECT id, `key`, par_value, description FROM cf_filteration_settings;")
	CntData, err := ExecuteQuery(Qry.String())
	if util.CheckErrorLog(r, err) {
		return nil, err
	}
	resBody := make(map[string]interface{})
	for _, v := range CntData {
		key := v["key"].(string)
		value := v["par_value"].(string)
		resBody[key] = value
	}
	return resBody, nil
}

// UpdateWebDefaultSettings - Returns settings front
func UpdateWebDefaultSettings(r *http.Request, reqMap map[string]interface{}) bool {
	util.LogIt(r, "Model - web_setting - UpdateWebDefaultSettings")

	var Qry bytes.Buffer
	Qry.WriteString("SELECT id, `key`, par_value, description FROM cf_filteration_settings;")
	oldData, err := ExecuteQuery(Qry.String())
	if util.CheckErrorLog(r, err) {
		return false
	}

	oldDataVal1 := make(map[string]interface{})
	oldDataVal2 := make(map[string]interface{})
	for _, v := range oldData {
		key := v["key"].(string)
		idKey := v["id"].(string)
		value := v["par_value"].(string)
		oldDataVal1[key] = value
		oldDataVal2[key] = idKey
	}

	for k1, v1 := range reqMap {
		if _, ok := oldDataVal1[k1]; ok {
			if v1.(string) != oldDataVal1[k1] {
				var Qry1 bytes.Buffer
				Qry1.WriteString("UPDATE cf_filteration_settings SET par_value = ? WHERE id = ?")
				err = ExecuteNonQuery(Qry1.String(), v1.(string), oldDataVal2[k1])
				if util.CheckErrorLog(r, err) {
					return false
				}
				AddLog(r, k1, "GEN_WEB_SETTING", k1+" Updated", oldDataVal2[k1].(string), map[string]interface{}{k1: v1.(string)})
			}
		}
	}
	return true
}
