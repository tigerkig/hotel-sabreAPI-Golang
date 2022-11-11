package model

import (
	"database/sql"
	"errors"
	"log"
)

// VDB - Database Variable, It is defined global
var VDB *sql.DB

// ExecuteQuery - Select Query.
func ExecuteQuery(SQLQry string, params ...interface{}) (retVal []map[string]interface{}, err error) {
	defer func() {
		if r := recover(); r != nil {
			retVal = nil
			err = errors.New("DB Side Error")
			log.Println("DB Exception Log : ", r)
		}
	}()

	rows, Qerr := VDB.Query(SQLQry, params...)
	if Qerr != nil {
		return nil, Qerr
	}

	columns, _ := rows.Columns()

	count := len(columns)
	values := make([]interface{}, count)

	valuePtrs := make([]interface{}, count)

	rowCnt := 0
	for rows.Next() {
		retObj := make(map[string]interface{}, count)

		for i := range columns {
			valuePtrs[i] = &values[i]
		}

		rows.Scan(valuePtrs...)

		for i, col := range columns {

			var v interface{}

			val := values[i]

			b, ok := val.([]byte)

			if ok {
				v = string(b)
			} else {
				v = val
			}
			retObj[col] = v
		}
		retVal = append(retVal, retObj)
		rowCnt = rowCnt + 1
	}

	if len(retVal) == 0 {
		retVal = []map[string]interface{}{}
	}

	return retVal, err
}

// ExecuteNonQuery - Update, Insert Query
func ExecuteNonQuery(SQLQry string, params ...interface{}) error {
	_, err := VDB.Exec(SQLQry, params...)
	return err
}

// ExecuteInsertWithID - Insert and return id
func ExecuteInsertWithID(SQLQry string, params ...interface{}) (string, error) {
	var ID string

	err := VDB.QueryRow(SQLQry, params...).Scan(&ID)
	if err != nil {
		return "0", err
	}
	return ID, err
}

// ExecuteRowQuery - Select Query return only one row.
func ExecuteRowQuery(SQLQry string, params ...interface{}) (retVal map[string]interface{}, err error) {
	rows, Qerr := VDB.Query(SQLQry, params...)
	if Qerr != nil {
		return nil, Qerr
	}
	defer func() {
		rows.Close()
		if r := recover(); r != nil {
			retVal = nil
			log.Println("DB Exception Log : ", r)
		}
	}()
	columns, _ := rows.Columns()
	count := len(columns)
	values := make([]interface{}, count)
	valuePtrs := make([]interface{}, count)

	for rows.Next() {
		retObj := make(map[string]interface{}, count)
		for i := range columns {
			valuePtrs[i] = &values[i]
		}
		rows.Scan(valuePtrs...)
		for i, col := range columns {
			var v interface{}
			val := values[i]
			b, ok := val.([]byte)
			if ok {
				v = string(b)
			} else {
				v = val
			}
			retObj[col] = v
		}
		retVal = retObj
		break
	}
	return retVal, err
}
