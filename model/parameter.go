package model

import (
	"bytes"
)

// GetParameter - It returns parameter
func GetParameter(key string) (string, error) {
	var Qry bytes.Buffer
	Qry.WriteString("SELECT value FROM cf_parameter WHERE `key`= ?")
	data, err := ExecuteRowQuery(Qry.String(), key)
	if err != nil {
		return "", err
	}
	return data["value"].(string), nil
}

// SetParameter - It sets parameter
func SetParameter(key string, val string) error {
	var Qry bytes.Buffer
	Qry.WriteString("UPDATE cf_parameter SET value=? WHERE `key`=?")
	err := ExecuteNonQuery(Qry.String(), val, key)
	if err != nil {
		return err
	}
	return nil
}
