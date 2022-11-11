package model

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"
	"tp-api-common/data"
	"tp-api-common/util"
	"tp-system/config"

	"github.com/gorilla/context"
	gonanoid "github.com/matoous/go-nanoid"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

//VMongoSession := Mongo Db Connection Variable
var VMongoSession *mgo.Session

// CheckDuplicateRecords - It checks if any duplicate records exist or not.
func CheckDuplicateRecords(r *http.Request, module string, OrSet map[string]string, AndSet map[string]string, primKey string) (int64, error) {
	util.LogIt(r, "Models - Models - CheckDuplicateRecords")
	var cnt int64 = 1
	var Qry bytes.Buffer
	if _, ok := util.StaticModules[module]; ok {
		ModuleData := util.StaticModules[module].(map[string]string)
		if len(ModuleData) == 0 {
			return 2, nil
		}
		PrimaryKey := ModuleData["ColPrimary"]
		Tbl := ModuleData["tbl"]
		Status := ModuleData["ColStatus"]
		orcnt := 0
		andcnt := 0
		if module == "USER_PROFILE" {
			Qry.WriteString("SELECT CONVERT(count(SU." + PrimaryKey + "),UNSIGNED INTEGER) AS cnt FROM cf_user AS SU INNER JOIN cf_user_profile AS SUP ON SUP.user_id = SU.id WHERE SU.status != 3 AND ")
		} else {
			if Status == "" {
				Qry.WriteString("SELECT CONVERT(count(" + PrimaryKey + "),UNSIGNED INTEGER) AS cnt FROM " + Tbl + " WHERE 1 = 1 AND ")
			} else {
				Qry.WriteString("SELECT CONVERT(count(" + PrimaryKey + "),UNSIGNED INTEGER) AS cnt FROM " + Tbl + " WHERE " + Status + "!=3 AND ")
			}
		}

		if OrSet != nil {
			Qry.WriteString("(")
			for key, value := range OrSet {
				if orcnt == 0 {
					Qry.WriteString(key + " LIKE '" + value + "'")
				} else {
					Qry.WriteString(" AND " + key + " LIKE '" + value + "'")
				}
				orcnt++
			}
			Qry.WriteString(")")
		}
		if AndSet != nil {
			Qry.WriteString("(")
			for key, value := range AndSet {
				if andcnt == 0 {
					Qry.WriteString(key + " LIKE '" + value + "'")
				} else {
					Qry.WriteString(" OR " + key + " LIKE '" + value + "'")
				}
				andcnt++
			}
			Qry.WriteString(")")
		}
		if primKey != "0" {
			if module == "USER_PROFILE" {
				Qry.WriteString(" AND SU." + PrimaryKey + "!='" + primKey + "'")
			} else {
				Qry.WriteString(" AND " + PrimaryKey + "!='" + primKey + "'")
			}
		}

		util.LogIt(r, " Duplicate Checkings ")
		util.LogIt(r, Qry.String())
		CntData, err := ExecuteQuery(Qry.String())
		if err != nil {
			return 2, err
		}
		IntV, _ := strconv.Atoi(CntData[0]["cnt"].(string))
		cnt = int64(IntV)
	} else {
		return 2, nil
	}
	return cnt, nil
}

// UpdateStatusModuleWise - Update Status of master. Before use it please add db info of master module in StaticModules.
func UpdateStatusModuleWise(r *http.Request, module string, status int, id string) (int, int, error) {
	util.LogIt(r, fmt.Sprint("Models - Models - UpdateStatus - Module - ", module, " - Status - ", status, " id - ", id))
	var Qry, SQLQry bytes.Buffer
	var CommonOperation string
	if _, ok := util.StaticModules[module]; !ok {
		return 0, 500, nil
	}

	if status == 1 {
		if module == "ROOM_TYPE" {
			invFilled := CheckAndFillInvdataOnActiveStatus(r, context.Get(r, "HotelId").(string), id)
			if !invFilled {
				return 0, 500, nil
			}

			basicRoomInfo := UpdateAllRoomType(context.Get(r, "HotelId").(string), id) // 2020-06-24 - HK - Room Add Sync With Mongo Added - Partner Panel
			if !basicRoomInfo {
				util.LogIt(r, "common - common - Error While Syncing Room Basic Info")
				return 0, 500, nil
			}
			roomImgFlg := UpdateRoomImage(context.Get(r, "HotelId").(string), id)
			if !roomImgFlg {
				util.LogIt(r, "common - common - Error While Syncing Room Images")
				return 0, 500, nil
			}
			rroomAmntFlg := UpdateRoomAmenity(context.Get(r, "HotelId").(string), id)
			if !rroomAmntFlg {
				util.LogIt(r, "common - common - Error While Syncing Room Amenity")
				return 0, 500, nil
			}
		}

		if module == "RATE_PLAN" {

			rateRest, roomID := CheckAndFillRateRestDataOnActiveStatus(r, context.Get(r, "HotelId").(string), id)
			if !rateRest {
				return 0, 500, nil
			}
			rRateRestFlg := AddUpdateRateplanDetails(context.Get(r, "HotelId").(string), id) // 2020-06-24 - HK - Rateplan Add Sync With Mongo Added - Partner Panel
			if !rRateRestFlg {
				util.LogIt(r, "common - common - Error While Syncing Rate Plan Basic Data")
				return 0, 500, nil
			}

			if roomID != "" {
				rRateRestFlg := UpdateRatePlanDeals(context.Get(r, "HotelId").(string), roomID, id) // 2020-06-24 - HK - Rateplan Add Sync With Mongo Added - Partner Panel
				if !rRateRestFlg {
					util.LogIt(r, "common - common - Error While Syncing Rate Plan Rate Rest Data")
					return 0, 500, nil
				}
			}

		}

	} else if status == 2 || status == 3 {
		if module == "TAX" {
			CacheChn <- CacheObj{
				Type:       "tax",
				ID:         context.Get(r, "HotelId").(string),
				Additional: "",
			}
		}

		if module == "ROOM_TYPE" {
			if IsHotelApproved(id) {
				reqMap := data.DeleteCacheItem{
					HotelID:    context.Get(r, "HotelId").(string),
					RoomTypeID: id,
				}
				jsonList, _ := json.Marshal(reqMap)
				_, scode := SendMicroServiceRequest("POST", "deleteRoomType", string(jsonList))
				if scode != 200 {
					return 0, 500, nil
				}
			}

		}
	}

	Vmodule := util.StaticModules[module].(map[string]string)
	Vtable := Vmodule["tbl"]
	StatusCol := Vmodule["ColStatus"]
	ColPrimary := Vmodule["ColPrimary"]
	if Vtable != "" && StatusCol != "" && ColPrimary != "" {
		Qry.WriteString("UPDATE " + Vtable + " SET " + StatusCol + "=? WHERE " + ColPrimary + "=?;")
		err := ExecuteNonQuery(Qry.String(), status, id)
		if util.CheckErrorLog(r, err) {
			return 0, 500, err
		}
	} else {
		return 0, 0, nil
	}

	SQLQry.WriteString("SELECT status FROM status WHERE id=?")
	NewStatus, err := ExecuteRowQuery(SQLQry.String(), status)
	if util.CheckErrorLog(r, err) {
		return 0, 500, err
	}

	if status == 3 {
		CommonOperation = "Delete Status"
	} else {
		CommonOperation = "Update Status"
	}

	AddLog(r, "", module, CommonOperation, id, map[string]interface{}{"Status": NewStatus["status"]})

	return 1, 204, nil
}

// ChangeSortOrderOfModule - Change Sort Order Of Module
func ChangeSortOrderOfModule(r *http.Request, module string, SortOrder int, id string) (int, error) {
	util.LogIt(r, fmt.Sprint("Models - Models - ChangeSortOrderOfModule - Module - ", module, " - Sort Order - ", SortOrder, " id - ", id))
	var SQLQry, SiQry, SQLUpdate, SQLNewUpdate bytes.Buffer
	if _, ok := util.StaticModules[module]; !ok {
		return 0, nil
	}
	Vmodule := util.StaticModules[module].(map[string]string)
	Vtable := Vmodule["tbl"]
	ColPrimary := Vmodule["ColPrimary"]

	if Vtable != "" && ColPrimary != "" {
		SiQry.WriteString(" SELECT sortorder FROM " + Vtable + "  WHERE " + ColPrimary + "= ?")
		SortOrderOfGivenID, err := ExecuteRowQuery(SiQry.String(), id)
		if util.CheckErrorLog(r, err) {
			return 0, err
		}

		SQLQry.WriteString(" SELECT id FROM " + Vtable + " WHERE sortorder = ? AND hotel_id = ?")
		IDOfChangeSortOrder, err := ExecuteRowQuery(SQLQry.String(), SortOrder, context.Get(r, "HotelId"))
		if util.CheckErrorLog(r, err) {
			return 0, err
		}

		SQLUpdate.WriteString("UPDATE " + Vtable + " SET sortorder = ? WHERE " + ColPrimary + " = ?;")
		err = ExecuteNonQuery(SQLUpdate.String(), SortOrder, id)
		if util.CheckErrorLog(r, err) {
			return 0, err
		}

		SQLNewUpdate.WriteString("UPDATE " + Vtable + " SET sortorder = ? WHERE " + ColPrimary + " = ?;")
		err = ExecuteNonQuery(SQLNewUpdate.String(), SortOrderOfGivenID["sortorder"], IDOfChangeSortOrder["id"])
		if util.CheckErrorLog(r, err) {
			return 0, err
		}

	} else {
		return 0, nil
	}

	AddLog(r, "", "HOTEL", "Update Sort Order Of Image", context.Get(r, "HotelId").(string), map[string]interface{}{})
	return 1, nil
}

// GetLocalDateTime - Returns local datetime
func GetLocalDateTime(r *http.Request) string {
	hour, min := getTimeZoneParam(r)
	localDateTime := time.Now().UTC().Add(time.Hour*time.Duration(hour) + time.Minute*time.Duration(min)).Format("2006-01-02 15:04:05")
	return localDateTime
}

func getTimeZoneParam(r *http.Request) (int, int) {
	var timezone string
	timezone = getParameter(r, "timezone")
	if timezone == "" {
		timezone = "08:00"
	}
	timeSlice := strings.Split(timezone, ":")
	hour, _ := strconv.Atoi(timeSlice[0])
	min, _ := strconv.Atoi(timeSlice[1])
	return hour, min
}

func getParameter(r *http.Request, key string) string {
	util.LogIt(r, "model - functions - getParameter")
	var Qry bytes.Buffer
	// Qry.WriteString("SELECT value FROM cf_parameter WHERE keyname=$1")
	Qry.WriteString("SELECT value FROM cf_parameter WHERE `key` = ?")
	RetMap, err := ExecuteRowQuery(Qry.String(), key)
	if util.CheckErrorLog(r, err) {
		return ""
	}
	if _, ok := RetMap["value"].(string); ok {
		return RetMap["value"].(string)
	}
	return ""
}

func setParameter(r *http.Request, key string, value string) {
	util.LogIt(r, "model - functions - setParameter")
	var Qry bytes.Buffer
	// Qry.WriteString("UPDATE cf_parameter set value=$1 WHERE keyname=$2")
	Qry.WriteString("UPDATE cf_parameter set value = ? WHERE `key` = ?")
	err := ExecuteNonQuery(Qry.String(), value, key)
	util.CheckErrorLog(r, err)
}

//SendRequest - Send Curl Request
func SendRequest(Token string, jsonStr []byte) bool {
	url := "http://10.68.101.90:9015/api/v1/wallet"
	//fmt.Println("URL:>", url)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("X-Auth-Token", Token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 201 {
		return false
	}

	return true
}

// GetSelectData - Get Select All Data Created By Meet Soni 2Nd May 2020
func GetSelectData(r *http.Request, module string, id string, deleteFlag bool, status string) (interface{}, error) {
	var err error
	var Qry bytes.Buffer
	var store []map[string]interface{}
	var mainArr interface{}
	ModuleMap := util.StaticModules[module].(map[string]string)
	table := ModuleMap["tbl"]
	ColPrimary := ModuleMap["ColPrimary"]
	ColStatus := ModuleMap["ColStatus"]
	mainArrStatus := strings.Split(status, ",")
	values := []string{}
	for _, val := range mainArrStatus {
		values = append(values, val)
	}

	result := "'" + strings.Join(values, "','") + "'"

	Qry.WriteString("SELECT * FROM " + table + " WHERE 1 = 1 AND " + ColStatus + " in (" + result + ")")
	if id != "" {
		Qry.WriteString(" AND " + ColPrimary + " = ?")
	}

	if id != "" {
		store, err = ExecuteQuery(Qry.String(), id)
		if util.CheckErrorLog(r, err) {
			return nil, err
		}
	} else {
		store, err = ExecuteQuery(Qry.String())
		if util.CheckErrorLog(r, err) {
			return nil, err
		}
	}

	if deleteFlag {
		for i := 0; i < len(store); i++ {
			delete(store[i], "created_at")
			delete(store[i], "created_by")
			delete(store[i], "updated_at")
			delete(store[i], "updated_by")
			delete(store[i], "status")
		}
	}

	if len(store) == 1 && id != "" {
		mainArr = store[0]
	} else {
		mainArr = store
	}

	return mainArr, nil
}

// GetModuleFieldByID - Returns field from database.
func GetModuleFieldByID(r *http.Request, module string, id string, obj string) (interface{}, error) {
	util.LogIt(r, fmt.Sprint("Models - Models - GetModuleFieldByID"))
	var err error

	ModuleMap := util.StaticModules[module].(map[string]string)

	table := ModuleMap["tbl"]
	ColPrimary := ModuleMap["ColPrimary"]
	store := make(map[string]interface{})
	var Qry bytes.Buffer
	if module == "USER_PROFILE" {
		Qry.WriteString("SELECT CU.username FROM cf_user_profile AS CUP INNER JOIN cf_user AS CU ON CU.id = CUP.user_id WHERE CUP.user_id = ? ")
		store, err = ExecuteRowQuery(Qry.String(), id)
		if util.CheckErrorLog(r, err) {
			return nil, err
		}
	} else {
		Qry.WriteString("SELECT " + obj + " AS field FROM " + table + " WHERE " + ColPrimary + "=?")
		store, err = ExecuteRowQuery(Qry.String(), id)
		if util.CheckErrorLog(r, err) {
			return nil, err
		}
	}

	if store != nil || len(store) > 0 {
		return store["field"], nil
	}
	return nil, nil
}

// GetModuleData - Returns array or field from database.
func GetModuleData(r *http.Request, module string, id string) (map[string]interface{}, error) {
	util.LogIt(r, fmt.Sprint("Models - Models - GetModuleData"))
	var err error

	store := make(map[string]interface{})
	Data := []map[string]interface{}{}
	var Qry bytes.Buffer
	if module == "country" {
		Qry.WriteString("SELECT id, sortname,phonecode,name FROM cf_country ")
		Data, err = ExecuteQuery(Qry.String())
		if util.CheckErrorLog(r, err) {
			return nil, err
		}
	} else if module == "state" {
		Qry.WriteString("SELECT id,name FROM cf_states WHERE country_id = ? ")
		Data, err = ExecuteQuery(Qry.String(), id)
		if util.CheckErrorLog(r, err) {
			return nil, err
		}

	} else if module == "city" {
		Qry.WriteString("SELECT id,name FROM cf_city WHERE state_id = ? ")
		Data, err = ExecuteQuery(Qry.String(), id)
		if util.CheckErrorLog(r, err) {
			return nil, err
		}
	} else if module == "locality" {
		Qry.WriteString("SELECT id, locality as name FROM cf_locality WHERE city_id = ? ")
		Data, err = ExecuteQuery(Qry.String(), id) // note : change this with id
		if util.CheckErrorLog(r, err) {
			return nil, err
		}
	}

	if len(Data) > 0 {
		store["data"] = Data
	} else {
		store["data"] = []string{}
	}
	return store, nil
}

// AddLog - Add Activity Logs in MongoDB For Admin & Partner Panel
func AddLog(r *http.Request, Mainobject string, module string, Operation string, id string, dataMap map[string]interface{}) {
	util.LogIt(r, fmt.Sprint("Model - Model - AddLog"))
	var err error
	var AuthData = make(map[string]string)
	Panel := context.Get(r, "Side").(string)
	if VMongoSession == nil {
		util.PrintException(r, "Session Nil")
		return
	}

	err = VMongoSession.Ping()
	if util.CheckErrorLog(r, err) {
		util.PrintException(r, "Error While Connection Mongo")
		return
	}

	localDate := util.GetISODate()

	jsonByte, err := json.Marshal(dataMap)
	if util.CheckErrorLog(r, err) {
		return
	}

	var UserID string
	if _, ok := context.Get(r, "UserId").(string); ok {
		UserID = context.Get(r, "UserId").(string)
	}

	json := string(jsonByte)
	VisiterIP := context.Get(r, "Visitor_IP")

	Token := context.Get(r, "Request-Token")
	TokenStr := Token.(string)

	if Panel == "TP-BACKOFFICE" {
		AuthData, err = GetAuthDetails(r, TokenStr)
		if util.CheckErrorLog(r, err) {
			return
		}
	} else if Panel == "TP-PARTNER" {
		AuthData, err = GetConsoleAuthDetails(r, TokenStr)
		if util.CheckErrorLog(r, err) {
			return
		}
	} else if Panel == "TP-FRONT" {
		AuthData["username"] = "web-user"
		if util.CheckErrorLog(r, err) {
			return
		}
	}

	var UserName string
	UserName = AuthData["username"]
	if _, ok := util.StaticModules[module]; ok {
		ModuleMap := util.StaticModules[module].(map[string]string)
		obj := ModuleMap["object"]
		Collection := ModuleMap["collection"]

		Object, err := GetModuleFieldByID(r, module, id, obj)
		if err != nil || Object == nil {
			return
		}

		if Mainobject == "" {
			Mainobject = Object.(string)
		}

		var m *mgo.Collection
		if Panel == "TP-BACKOFFICE" || Panel == "TP-FRONT" {
			m = VMongoSession.DB(config.Env.Mongo.MongoDB).C(Collection)
		} else if Panel == "TP-PARTNER" {
			m = VMongoSession.DB(config.Env.Mongo.MongoDB).C("partner_" + Collection)
		}

		//c := VMongoSession.DB(config.Env.Mongo.MongoDB).C(Collection)
		MID := bson.NewObjectId()
		err = m.Insert(&data.Logs{MID, id, Mainobject, UserID, UserName, Operation, json, localDate, VisiterIP.(string)})
		if util.CheckErrorLog(r, err) {
			return
		}
		util.LogIt(r, fmt.Sprint("Log Added Successfully  : "+Object.(string)))
	}
	return
}

// GetActivityWiseLog - Returns Activity Logs  By Passing Id Of Object, Module
func GetActivityWiseLog(r *http.Request, module string, ObjectID string, limit int, offset int, dir string, sort int, search []data.JquerySearch) (map[string]interface{}, error) {
	util.LogIt(r, fmt.Sprint("Model - Functions - GetActivityWiseLog"))

	var err error
	var UserID, IP, Operation, FromDate, ToDate string
	var FromDateInt, ToDateInt int64
	var stuff = make(map[string]interface{})
	if VMongoSession == nil {
		util.PrintException(r, "Session Nil")
		return nil, err
	}
	err = VMongoSession.Ping()
	if util.CheckErrorLog(r, err) {
		return nil, err
	}

	var testColArrs [7]string
	testColArrs[0] = "objid"
	testColArrs[1] = "object"
	testColArrs[2] = "operation"
	testColArrs[3] = "data"
	testColArrs[4] = "uname"
	testColArrs[5] = "userip"
	testColArrs[6] = "datetime"

	if len(search) > 0 {
		for _, v := range search {
			if v.Field == "userId" {
				UserID = v.Value
			}
			if v.Field == "ip" {
				IP = v.Value
			}
			if v.Field == "operation" {
				Operation = v.Value
			}
			if v.Field == "datetime" && v.Operator == "sdate" {
				FromDate = v.Value
			}
			if v.Field == "datetime" && v.Operator == "edate" && v.Value != "" {
				ToDate = v.Value
				TodateCur, err := util.AddDaysToDate(r, ToDate, 1)
				if util.CheckErrorLog(r, err) {
					return nil, err
				}
				ToDate = TodateCur
			}
		}
	}

	var Logs []data.Logs

	if _, ok := util.StaticModules[module]; ok {
		ModuleMap := util.StaticModules[module].(map[string]string)
		Collection := ModuleMap["collection"]
		c := VMongoSession.DB(config.Env.Mongo.MongoDB).C(Collection)

		DateTime := make(map[string]interface{})
		Filtration := make(map[string]interface{})

		if ObjectID != "" {
			Filtration["objid"] = ObjectID
		}

		if UserID != "" {
			Filtration["userid"] = UserID
		}

		if IP != "" {
			Filtration["userip"] = IP
		}

		if Operation != "" {
			Filtration["operation"] = Operation
		}

		if FromDate != "" {
			FromDateInt = util.ConvertToISODateOnly(FromDate)
			DateTime["$gt"] = FromDateInt
		}

		if ToDate != "" {
			ToDateInt = util.ConvertToISODateOnly(ToDate)
			DateTime["$lt"] = ToDateInt
		}

		if FromDate != "" || ToDate != "" {
			Filtration["datetime"] = DateTime
		}

		direction := "-datetime"
		if dir == "asc" {
			direction = testColArrs[sort]
		} else if dir == "desc" {
			direction = "-" + testColArrs[sort]
		}

		c.Find(Filtration).Sort(direction).Limit(limit).Skip(offset).All(&Logs)
		FilterCnt, err := c.Find(Filtration).Count()
		if util.CheckErrorLog(r, err) {
			return nil, err
		}
		totalcount, err := c.Find(Filtration).Count()
		if util.CheckErrorLog(r, err) {
			return nil, err
		}

		res, err := json.Marshal(Logs)
		if util.CheckErrorLog(r, err) {
			return nil, err
		}

		var reqMap []map[string]interface{}
		json.Unmarshal(res, &reqMap)

		var retVal []map[string]interface{}
		if len(reqMap) > 0 {
			for _, val := range reqMap {
				var SubData = make(map[string]interface{})
				SubData["objectId"] = val["objid"]
				SubData["object"] = val["object"]
				SubData["userName"] = val["uname"]
				SubData["ip"] = val["userip"]
				SubData["operation"] = val["operation"]
				dataStr := val["data"].(string)
				if dataStr != "" {
					var DataMap = make(map[string]interface{})
					json.Unmarshal([]byte(dataStr), &DataMap)
					Strdata := ""
					if len(DataMap) > 0 {
						i := 0
						for key1, val1 := range DataMap {
							switch val1.(type) {
							case float64:
								val1 = fmt.Sprintf("%.0f", val1)
							case int:
								val1 = strconv.Itoa(val1.(int))
							case int64:
								val1 = strconv.Itoa(int(val1.(int64)))
							case interface{}:
								val1 = fmt.Sprintf("%v", val1)
							}

							if i == 0 {
								if val1 != nil {
									Strdata += key1 + ": " + val1.(string)
								} else {
									Strdata += key1 + ": " + " "
								}

							} else {
								if val1 != nil {
									Strdata += " , " + key1 + ": " + val1.(string)
								} else {
									Strdata += " , " + key1 + ": " + " "
								}

							}
							i = i + 1
						}
					}
					SubData["data"] = Strdata
				} else {
					SubData["data"] = ""
				}
				datetime, err := util.ISOToDateTime(int64(val["datetime"].(float64)))
				if util.CheckErrorLog(r, err) {
					return nil, err
				}
				SubData["datetime"] = datetime
				retVal = append(retVal, SubData)
			}
		}
		stuff["recordsTotal"] = totalcount
		stuff["recordsFiltered"] = FilterCnt

		if len(retVal) == 0 {
			stuff["data"] = []string{}
		} else {
			stuff["data"] = retVal
		}
		return stuff, nil
	}

	return nil, nil
}

// CommonChangeSortOrderOfModule - Change Sort Order Of Module
func CommonChangeSortOrderOfModule(r *http.Request, module string, SortOrder int, id string, apiModule string) (int, error) {
	util.LogIt(r, fmt.Sprint("Models - Models - CommonChangeSortOrderOfModule - Module - ", module, " - Sort Order - ", SortOrder, " id - ", id, " - apiModule - ", apiModule))
	var SQLQry, SiQry, SQLUpdate, SQLNewUpdate bytes.Buffer
	if _, ok := util.StaticModules[module]; !ok {
		return 0, nil
	}
	Vmodule := util.StaticModules[module].(map[string]string)
	Vtable := Vmodule["tbl"]
	ColPrimary := Vmodule["ColPrimary"]

	if Vtable != "" && ColPrimary != "" {
		SiQry.WriteString(" SELECT sortorder FROM " + Vtable + "  WHERE " + ColPrimary + "= ?")
		SortOrderOfGivenID, err := ExecuteRowQuery(SiQry.String(), id)
		if util.CheckErrorLog(r, err) {
			return 0, err
		}

		SQLQry.WriteString(" SELECT id FROM " + Vtable + " WHERE sortorder = ? AND hotel_id = ?")
		IDOfChangeSortOrder, err := ExecuteRowQuery(SQLQry.String(), SortOrder, context.Get(r, "HotelId"))
		if util.CheckErrorLog(r, err) {
			return 0, err
		}

		SQLUpdate.WriteString("UPDATE " + Vtable + " SET sortorder = ? WHERE " + ColPrimary + " = ?;")
		err = ExecuteNonQuery(SQLUpdate.String(), SortOrder, id)
		if util.CheckErrorLog(r, err) {
			return 0, err
		}

		SQLNewUpdate.WriteString("UPDATE " + Vtable + " SET sortorder = ? WHERE " + ColPrimary + " = ?;")
		err = ExecuteNonQuery(SQLNewUpdate.String(), SortOrderOfGivenID["sortorder"], IDOfChangeSortOrder["id"])
		if util.CheckErrorLog(r, err) {
			return 0, err
		}

	} else {
		return 0, nil
	}

	LogPurpose := util.DependentStaticAPIModules[apiModule].(map[string]string)
	IsSpecialLog := LogPurpose["special_log"]
	if IsSpecialLog == "true" {

		LogModule := Vmodule["dep_log_module"] // ROOM_TYPE

		// get dependent table name upon module
		MainTable := Vmodule["main_tbl"] // cf_room_image
		// get dependent column name upon id being sent in request
		MainTableRefColumn := Vmodule["main_tbl_ref_col"] // room_type_id

		SQLQry.WriteString(" SELECT " + MainTableRefColumn + " as depvalue FROM " + MainTable + " WHERE id = ?")
		DepColumnValue, err := ExecuteRowQuery(SQLQry.String(), id)
		if util.CheckErrorLog(r, err) {
			return 0, err
		}

		mainID := DepColumnValue["depvalue"]
		// get Original Table
		MainRefTable := Vmodule["ref_table"] // ref_table
		// get dependent column value upon first one's value
		MainTableOrigColumn := Vmodule["ref_table_column"] // ROOM_TYPE

		SQLQry.WriteString(" SELECT " + MainTableOrigColumn + " as log_column FROM " + MainRefTable + " WHERE id = ?")
		DepColumnValue, err1 := ExecuteRowQuery(SQLQry.String(), mainID)
		if util.CheckErrorLog(r, err1) {
			return 0, err1
		}
		AddLog(r, "", LogModule, "Update Sort Order Of Image", mainID.(string), map[string]interface{}{"Object Name": DepColumnValue["log_column"].(string)})
	} else {
		AddLog(r, "", "HOTEL", "Update Sort Order Of Image", context.Get(r, "HotelId").(string), map[string]interface{}{})
	}
	return 1, nil
}

// CheckDependentRecord - Checks If Record Exists In Depdendent Table - HK - 2020-05-14
func CheckDependentRecord(r *http.Request, id string, module string) (int64, error) {
	util.LogIt(r, fmt.Sprint("Model - Function - CheckDependentRecord"))
	var cnt int64

	if _, ok := util.DependentTables[module]; ok {

		ModuleData := util.DependentTables[module].(map[string]string)

		DependentTable := ModuleData["dependentTable"]
		DependentColumn := ModuleData["lnkid"]

		var dependentQry bytes.Buffer
		dependentQry.WriteString("SELECT count(id) AS cnt FROM " + DependentTable + " WHERE " + DependentColumn + " = '" + id + "' AND status = '1'")
		CntData, err := ExecuteRowQuery(dependentQry.String())
		if err != nil {
			return -1, err
		}
		cnt, _ = strconv.ParseInt(CntData["cnt"].(string), 10, 64)
	} else {
		return -1, nil
	}
	return cnt, nil
}

// GetYearMonthSliceBetweenTwoDates - makes unique string slice
func GetYearMonthSliceBetweenTwoDates(r *http.Request, startDate string, endDate string) []map[string]interface{} {

	noOfDays := DaysBetween(r, DateTime(r, endDate), DateTime(r, startDate))
	yearMonth := []string{}
	fullDate := []string{}
	for i := 0; i <= noOfDays; i++ {
		insertDate := DateTime(r, startDate).AddDate(0, 0, i).Format("2006-01-02")

		dateIno := strings.Split(insertDate, "-")

		yearMonthString := dateIno[0] + "-" + dateIno[1]
		yearMonth = append(yearMonth, yearMonthString)
		fullDate = append(fullDate, insertDate)
	}
	// log.Println(fullDate)

	keys := make(map[string]bool)
	list := []string{}
	for _, entry := range yearMonth {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}

	var yearMonthint []map[string]interface{}
	for k, v := range list {

		yearMonthint = append(yearMonthint, map[string]interface{}{
			"year_month": v,
		})

		/*var DateArr []map[string]interface{}
		for _, v1 := range fullDate {

			dateIno := strings.Split(v1, "-")
			yearMonthString := dateIno[0] + "-" + dateIno[1]

			if yearMonthString == v {
				DateArr = append(DateArr, map[string]interface{}{
					v1: "1",
				})
			}
		}*/

		var DateArr = make(map[string]int64)
		for _, v1 := range fullDate {

			dateIno := strings.Split(v1, "-")
			yearMonthString := dateIno[0] + "-" + dateIno[1]

			if yearMonthString == v {
				DateArr[v1] = 1
			}
		}
		yearMonthint[k]["data"] = DateArr
	}

	// return list
	// considering start date as 2020-12-31 and end date as 2020-02-01
	// [{	"data": {		"2020-12-31": "1"	},	"year_month": "2020-12"}, {	"data": {		"2021-01-01": "1",		"2021-01-02": "1",		"2021-01-03": "1",		"2021-01-04": "1",		"2021-01-05": "1",		"2021-01-06": "1",		"2021-01-07": "1",		"2021-01-08": "1",		"2021-01-09": "1",		"2021-01-10": "1",		"2021-01-11": "1",		"2021-01-12": "1",		"2021-01-13": "1",		"2021-01-14": "1",		"2021-01-15": "1",		"2021-01-16": "1",		"2021-01-17": "1",		"2021-01-18": "1",		"2021-01-19": "1",		"2021-01-20": "1",		"2021-01-21": "1",		"2021-01-22": "1",		"2021-01-23": "1",		"2021-01-24": "1",		"2021-01-25": "1",		"2021-01-26": "1",		"2021-01-27": "1",		"2021-01-28": "1",		"2021-01-29": "1",		"2021-01-30": "1",		"2021-01-31": "1"	},	"year_month": "2021-01"}, {	"data": {		"2021-02-01": "1"	},	"year_month": "2021-02"}]
	return yearMonthint

}

// UniqueString - makes unique string slice
func UniqueString(stringSlice []string) []string {
	keys := make(map[string]bool)
	list := []string{}
	for _, entry := range stringSlice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}

// DateTime - Returns Date Into Time
func DateTime(r *http.Request, s string) time.Time {
	d, _ := time.Parse("2006-01-02", s)
	return d
}

// DaysBetween - returns difference of days between 2 dates
func DaysBetween(r *http.Request, a, b time.Time) int {
	if a.After(b) {
		a, b = b, a
	}

	days := -a.YearDay()
	for year := a.Year(); year < b.Year(); year++ {
		days += time.Date(year, time.December, 31, 0, 0, 0, 0, time.UTC).YearDay()
	}
	days += b.YearDay()

	return days
}

// MonthInterval - Gets First And Last Date Of Month By Passing Year And Month
func MonthInterval(r *http.Request, y int, m time.Month) (firstDay, lastDay time.Time) {
	firstDay = time.Date(y, m, 1, 0, 0, 0, 0, time.UTC)
	lastDay = time.Date(y, m+1, 1, 0, 0, 0, -1, time.UTC)
	return firstDay, lastDay
}

// GetModuleFieldByIDFromDetailedTable - Returns latest record from table.
func GetModuleFieldByIDFromDetailedTable(r *http.Request, module string, id string, obj string) (interface{}, error) {

	util.LogIt(r, fmt.Sprint("Models - Models - GetModuleFieldByIDFromDetailedTable"))
	var err error

	ModuleMap := util.StaticModules[module].(map[string]string)

	table := ModuleMap["tbl"]
	ColPrimary := ModuleMap["ColPrimary"]
	store := make(map[string]interface{})

	var Qry bytes.Buffer
	Qry.WriteString("SELECT " + obj + " AS field FROM " + table + " WHERE " + ColPrimary + " = ? ORDER BY created_at DESC LIMIT 1")
	store, err = ExecuteRowQuery(Qry.String(), id)
	if util.CheckErrorLog(r, err) {
		return nil, err
	}

	if store != nil || len(store) > 0 {
		return store["field"], nil
	}
	return nil, nil
}

// GetModuleFieldsByID - Returns fields from table.
func GetModuleFieldsByID(r *http.Request, module string, id string, obj string) (map[string]interface{}, error) {
	util.LogIt(r, fmt.Sprint("Models - Models - GetModuleFieldByID"))
	var err error

	ModuleMap := util.StaticModules[module].(map[string]string)

	table := ModuleMap["tbl"]
	ColPrimary := ModuleMap["ColPrimary"]
	store := make(map[string]interface{})
	var Qry bytes.Buffer

	Qry.WriteString("SELECT " + obj + " FROM " + table + " WHERE " + ColPrimary + "=?")
	store, err = ExecuteRowQuery(Qry.String(), id)
	if util.CheckErrorLog(r, err) {
		return nil, err
	}

	if store != nil || len(store) > 0 {
		return store, nil
	}
	return nil, nil
}

// GetOneYearFirstLastDate - To get Current Month's First Date And Next Year's Last Month's Last Date - 2020-06-26 - HK
// Considering Date Today is 2020-05-19
// So result will be first : 2020-05-01 and last : 2021-05-31
func GetOneYearFirstLastDate(r *http.Request) (string, string) {
	util.LogIt(r, fmt.Sprint("Models - Functions - GetOneYearFirstLastDate"))

	currentYear, currentMonth, _ := time.Now().Date()
	first, _ := MonthInterval(r, currentYear, currentMonth)
	startDate := first.Format("2006-01-02")

	lastDate := time.Now().AddDate(1, 0, 0).Format("2006-01-02")
	newLastDate := strings.Split(lastDate, "-")
	endYear, _ := strconv.Atoi(newLastDate[0])
	endMonth, _ := strconv.Atoi(newLastDate[1])
	_, last := MonthInterval(r, endYear, time.Month(endMonth))
	endDate := last.Format("2006-01-02")

	return startDate, endDate

}

// CheckAndFillInvdataOnActiveStatus - Fills Inv Data Upon Active Status
func CheckAndFillInvdataOnActiveStatus(r *http.Request, hotelID string, roomID string) bool {
	util.LogIt(r, fmt.Sprint("Models - Functions - CheckAndFillInvdataOnActiveStatus"))

	startDate, endDate := GetOneYearFirstLastDate(r)

	var Qry bytes.Buffer
	Qry.WriteString("SELECT inventory FROM cf_room_type WHERE id = ? AND hotel_id = ?")
	RetMap, err := ExecuteRowQuery(Qry.String(), roomID, hotelID)
	if util.CheckErrorLog(r, err) {
		return false
	}

	baseInv := RetMap["inventory"].(int64)
	invSuccess := CheckAndFillInvData(r, hotelID, roomID, baseInv, startDate, endDate)
	if !invSuccess {
		return false
	}

	return true
}

// CheckAndFillInvData - Checks If Data Exists in DB and If not fills data
func CheckAndFillInvData(r *http.Request, HotelID string, roomID string, baseInv int64, startDate string, endDate string) bool {
	util.LogIt(r, "Model - Functions - CheckAndFillInvData")

	yearMonthDates := GetYearMonthSliceBetweenTwoDates(r, startDate, endDate)

	/* 2021-05-20 - Log Purpose */
	localDate := util.GetISODate()

	for _, val := range yearMonthDates {
		yearMonth := fmt.Sprintf("%v", val["year_month"]) // Convert interface to string
		dateData := val["data"]

		dateInfo := strings.Split(yearMonth, "-")
		tblYear := dateInfo[0]
		tblMonth := dateInfo[1]

		var Qry bytes.Buffer
		Qry.WriteString("SELECT inv_data FROM cf_inv_data WHERE hotel_id = ? AND room_id = ? AND year = ? AND month = ? ") // Add
		Data, err := ExecuteRowQuery(Qry.String(), HotelID, roomID, tblYear, tblMonth)
		if util.CheckErrorLog(r, err) {
			return false
		}

		if len(Data) == 0 {

			// https://medium.com/@prithvi_20863/interfaces-in-golang-a-short-anecdote-249d7c6f96f4
			dMap := reflect.ValueOf(dateData) // as dateData is interface{} type and to loop over through such type we need such conversion
			var DateArr = make(map[string]int64)
			for _, dateVal := range dMap.MapKeys() {
				// valTest := dMap.MapIndex(dateVal) // to access value of key use this line
				realDate, _ := reflect.Value(dateVal).Interface().(string)
				DateArr[realDate] = baseInv

				// Put Log Here - 2021-05-20 - HK - For Update Log
				var Qry1 bytes.Buffer
				nanoid, _ := gonanoid.Nanoid()
				Qry1.WriteString("INSERT INTO logs_inv (id, hotel_id, room_id, inventory, update_for_date, updated_at, updated_by, booking_id, ip) VALUES (?,?,?,?,?,?,?,?,?)")
				err = ExecuteNonQuery(Qry1.String(), nanoid, HotelID, roomID, baseInv, realDate, localDate, "", "", "")
				if util.CheckErrorLog(r, err) {
					return false
				}
				// Put Log Here - 2021-05-20 - HK - For Update Log
			}
			invJSON, err := json.Marshal(DateArr)

			var Qry1 bytes.Buffer
			nanoid, _ := gonanoid.Nanoid()
			Qry1.WriteString("INSERT INTO cf_inv_data (id, hotel_id, room_id, year, month, inv_data) VALUES (?,?,?,?,?,?)")
			err = ExecuteNonQuery(Qry1.String(), nanoid, HotelID, roomID, tblYear, tblMonth, string(invJSON))
			if util.CheckErrorLog(r, err) {
				return false
			}
		} else {

			existsData := Data["inv_data"].(string)
			// log.Println(existsData) fmt.Println(reflect.TypeOf(existsData)) // string

			// Declared an empty map interface
			var result map[string]interface{}

			// Unmarshal or Decode the JSON to the interface.
			json.Unmarshal([]byte(existsData), &result)

			// Print the data type of result variable
			// fmt.Println(result) fmt.Println(reflect.TypeOf(result)) // map[string]interface {}

			incomingData := reflect.ValueOf(dateData)
			for _, dateVal := range incomingData.MapKeys() {
				realDate, _ := reflect.Value(dateVal).Interface().(string)

				if _, ok := result[realDate]; !ok {
					//Update Json For Incoming Dates Which Not Exists In Table Inventory Data
					result[realDate] = baseInv

					// Put Log Here - 2021-05-20 - HK - For Update Log
					var Qry1 bytes.Buffer
					nanoid, _ := gonanoid.Nanoid()
					Qry1.WriteString("INSERT INTO logs_inv (id, hotel_id, room_id, inventory, update_for_date, updated_at, updated_by, booking_id, ip) VALUES (?,?,?,?,?,?,?,?,?)")
					err = ExecuteNonQuery(Qry1.String(), nanoid, HotelID, roomID, baseInv, realDate, localDate, "", "", "")
					if util.CheckErrorLog(r, err) {
						return false
					}
					// Put Log Here - 2021-05-20 - HK - For Update Log
				}
			}
			invJSON, err := json.Marshal(result)
			// fmt.Println(string(invJSON), err)

			var Qry1 bytes.Buffer
			Qry1.WriteString("UPDATE cf_inv_data SET inv_data = ? WHERE  hotel_id = ? AND  room_id = ? AND year = ? AND month = ?")
			err = ExecuteNonQuery(Qry1.String(), string(invJSON), HotelID, roomID, tblYear, tblMonth)
			if util.CheckErrorLog(r, err) {
				return false
			}
		}
	}
	return true
}

// CheckAndFillRateRestDataOnActiveStatus - Fills Rate Rest Data Upon Active Status
func CheckAndFillRateRestDataOnActiveStatus(r *http.Request, hotelID string, rateID string) (bool, string) {
	util.LogIt(r, fmt.Sprint("Models - Functions - CheckAndFillRateRestDataOnActiveStatus"))

	startDate, endDate := GetOneYearFirstLastDate(r)

	var Qry bytes.Buffer
	Qry.WriteString("SELECT CR.room_type_id as room_id, CRT.max_occupancy as max_occupancy, CR.rate as rate, CR.rate_plan_name as rate_name FROM cf_rateplan AS CR INNER JOIN cf_room_type AS CRT ON CRT.id = CR.room_type_id WHERE CR.id = ? AND CR.hotel_id = ?")
	RoomMap, err := ExecuteRowQuery(Qry.String(), rateID, hotelID)
	if util.CheckErrorLog(r, err) {
		return false, ""
	}

	RoomMaxOcc := RoomMap["max_occupancy"].(int64)
	roomID := RoomMap["room_id"].(string)
	var OccupancyArr []map[string]interface{}
	var j int64
	for j = 1; j <= RoomMaxOcc; j++ {
		joinStr := strconv.FormatInt(j, 16) //
		mainStr := "occ_" + joinStr
		OccupancyArr = append(OccupancyArr, map[string]interface{}{
			mainStr: RoomMap["rate"],
		})
	}

	occWiseRateInfoData := make(map[string]interface{})

	occWiseRateInfoData["rate_id"] = rateID
	occWiseRateInfoData["rate_name"] = RoomMap["rate_name"]
	occWiseRateInfoData["rate"] = OccupancyArr
	occWiseRateInfoData["start_date"] = startDate
	occWiseRateInfoData["end_date"] = endDate
	occWiseRateInfoData["min_night"] = 1
	occWiseRateInfoData["stop_sell"] = 0
	occWiseRateInfoData["cta"] = 0
	occWiseRateInfoData["ctd"] = 0

	rateSuccess := CheckAndFillRateRestrictionData(r, hotelID, roomID, occWiseRateInfoData)
	if !rateSuccess {
		return false, ""
	}
	return true, roomID
}

// CheckAndFillRateRestrictionData - Checks If Data Exists in DB and If not fills data
func CheckAndFillRateRestrictionData(r *http.Request, HotelID string, roomID string, rateInfo map[string]interface{}) bool {

	util.LogIt(r, "Model - Functions - CheckAndFillRateRestrictionData")

	startDate := rateInfo["start_date"].(string)
	endDate := rateInfo["end_date"].(string)
	rateID := rateInfo["rate_id"].(string)
	// occupancy := rateInfo["occupancy"].(int64)

	var dataDumpRatePlanWise = make(map[string]interface{})
	// dataDumpRatePlanWise["rate"] = rateInfo["rate"].(string)
	dataDumpRatePlanWise["rate"] = rateInfo["rate"].([]map[string]interface{})
	dataDumpRatePlanWise["min_night"] = rateInfo["min_night"].(int)
	dataDumpRatePlanWise["stop_sell"] = rateInfo["stop_sell"].(int)
	dataDumpRatePlanWise["cta"] = rateInfo["cta"].(int)
	dataDumpRatePlanWise["ctd"] = rateInfo["ctd"].(int)

	yearMonthDates := GetYearMonthSliceBetweenTwoDates(r, startDate, endDate)
	/* 2021-05-20 - Log Purpose */
	localDate := util.GetISODate()

	for _, val := range yearMonthDates {

		// insert rate, restriction data
		yearMonth := fmt.Sprintf("%v", val["year_month"]) // Convert interface to string

		dateData := val["data"]

		dateInfo := strings.Split(yearMonth, "-")
		tblYear := dateInfo[0]
		tblMonth := dateInfo[1]

		var Qry bytes.Buffer
		Qry.WriteString("SELECT rate_rest_data FROM cf_rate_restriction_data WHERE hotel_id = ? AND room_id = ? AND rateplan_id = ? AND year = ? AND month = ? ")
		RateRestDataFromTbl, err := ExecuteRowQuery(Qry.String(), HotelID, roomID, rateID, tblYear, tblMonth)
		if util.CheckErrorLog(r, err) {
			return false
		}

		if len(RateRestDataFromTbl) == 0 {

			// https://medium.com/@prithvi_20863/interfaces-in-golang-a-short-anecdote-249d7c6f96f4
			dMap := reflect.ValueOf(dateData) // as dateData is interface{} type and to loop over through such type we need such conversion
			var DateDataInsert = make(map[string]interface{})
			for _, dateVal := range dMap.MapKeys() {
				// valTest := dMap.MapIndex(dateVal) // to access value of key use this line
				realDate, _ := reflect.Value(dateVal).Interface().(string)
				DateDataInsert[realDate] = dataDumpRatePlanWise

				// Put Log Here - 2021-05-20 - HK - For Update Log
				rateDataForLog := dataDumpRatePlanWise["rate"]
				var rateStr string
				if x, ok := rateDataForLog.([]interface{}); ok {
					for _, e := range x {
						incomingData := reflect.ValueOf(e)
						for _, element := range incomingData.MapKeys() {
							valTest := incomingData.MapIndex(element) // to access value of key use this line
							occupancy, _ := reflect.Value(element).Interface().(string)
							occWiserate, _ := reflect.Value(valTest).Interface().(string)
							rateStr += occupancy + ":" + occWiserate + ", "
						}
					}
				}
				rateStr = strings.TrimRight(rateStr, ", ")

				var Qry2 bytes.Buffer
				nanoid, _ := gonanoid.Nanoid()
				Qry2.WriteString("INSERT INTO logs_rate_rest (id, hotel_id, room_id, rateplan_id,  update_for_date, rate, min_night, stop_sell, cta, ctd, updated_at, updated_by, ip) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?)")
				err = ExecuteNonQuery(Qry2.String(), nanoid, HotelID, roomID, rateID, realDate, rateStr, dataDumpRatePlanWise["min_night"], dataDumpRatePlanWise["stop_sell"], dataDumpRatePlanWise["cta"], dataDumpRatePlanWise["ctd"], localDate, "", "")
				if util.CheckErrorLog(r, err) {
					return false
				}
				// Put Log Here - 2021-05-20 - HK - For Update Log
			}

			rateRestJSON, err := json.Marshal(DateDataInsert)

			var Qry1 bytes.Buffer
			nanoid, _ := gonanoid.Nanoid()
			Qry1.WriteString("INSERT INTO cf_rate_restriction_data (id, hotel_id, room_id, rateplan_id, year, month, rate_rest_data)  VALUES (?,?,?,?,?,?,?)")
			err = ExecuteNonQuery(Qry1.String(), nanoid, HotelID, roomID, rateID, tblYear, tblMonth, string(rateRestJSON))
			if util.CheckErrorLog(r, err) {
				return false
			}

		} else {

			existsData := RateRestDataFromTbl["rate_rest_data"].(string)

			// Declared an empty map interface
			var result map[string]interface{}

			// Unmarshal or Decode the JSON to the interface.
			json.Unmarshal([]byte(existsData), &result)

			// Print the data type of result variable
			// fmt.Println(result) fmt.Println(reflect.TypeOf(result)) // map[string]interface {}

			incomingData := reflect.ValueOf(dateData)
			for _, dateVal := range incomingData.MapKeys() {

				realDate, _ := reflect.Value(dateVal).Interface().(string)
				if _, ok := result[realDate]; !ok {
					// Update Json For Incoming Dates Which Not Exists In Table Rate Restrcition Data
					result[realDate] = dataDumpRatePlanWise

					// Put Log Here - 2021-05-20 - HK - For Update Log
					rateDataForLog := dataDumpRatePlanWise["rate"]
					var rateStr string
					if x, ok := rateDataForLog.([]interface{}); ok {
						for _, e := range x {
							incomingData := reflect.ValueOf(e)
							for _, element := range incomingData.MapKeys() {
								valTest := incomingData.MapIndex(element) // to access value of key use this line
								occupancy, _ := reflect.Value(element).Interface().(string)
								occWiserate, _ := reflect.Value(valTest).Interface().(string)
								rateStr += occupancy + ":" + occWiserate + ", "
							}
						}
					}
					rateStr = strings.TrimRight(rateStr, ", ")

					var Qry2 bytes.Buffer
					nanoid, _ := gonanoid.Nanoid()
					Qry2.WriteString("INSERT INTO logs_rate_rest (id, hotel_id, room_id, rateplan_id,  update_for_date, rate, min_night, stop_sell, cta, ctd, updated_at, updated_by, ip) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?)")
					err = ExecuteNonQuery(Qry2.String(), nanoid, HotelID, roomID, rateID, realDate, rateStr, dataDumpRatePlanWise["min_night"], dataDumpRatePlanWise["stop_sell"], dataDumpRatePlanWise["cta"], dataDumpRatePlanWise["ctd"], localDate, "", "")
					if util.CheckErrorLog(r, err) {
						return false
					}
					// Put Log Here - 2021-05-20 - HK - For Update Log
				}
			}

			rateJSON, err := json.Marshal(result)

			var Qry1 bytes.Buffer
			Qry1.WriteString("UPDATE cf_rate_restriction_data SET rate_rest_data = ? WHERE  hotel_id = ? AND  room_id = ? AND rateplan_id = ? AND year = ? AND month = ?")
			err = ExecuteNonQuery(Qry1.String(), string(rateJSON), HotelID, roomID, rateID, tblYear, tblMonth)
			if util.CheckErrorLog(r, err) {
				return false
			}
		}
	}
	return true
}

// CheckDependentRecordWithoutStatus - Checks If Record Exists In Depdendent Table Or Not
// This function was created as we have common function for status update for all modules.
// But in the case of property tag, dependent table cf_hotel_tag doesn't have status column
// hence we have created same function as CheckDependentRecord function but without checking
// status column - HK - 2021-04-19
func CheckDependentRecordWithoutStatus(r *http.Request, id string, module string) (int64, error) {
	util.LogIt(r, fmt.Sprint("Model - Function - CheckDependentRecordWithoutStatus"))
	var cnt int64

	if _, ok := util.DependentTables[module]; ok {

		ModuleData := util.DependentTables[module].(map[string]string)

		DependentTable := ModuleData["dependentTable"]
		DependentColumn := ModuleData["lnkid"]

		var dependentQry bytes.Buffer
		dependentQry.WriteString("SELECT count(id) AS cnt FROM " + DependentTable + " WHERE " + DependentColumn + " = '" + id + "'")
		CntData, err := ExecuteRowQuery(dependentQry.String())
		if err != nil {
			return -1, err
		}
		cnt, _ = strconv.ParseInt(CntData["cnt"].(string), 10, 64)
	} else {
		return -1, nil
	}
	return cnt, nil
}

// CheckPrivilege - For checking priviledge assigned to user or not - 2021-05-04 - HK
func CheckPrivilege(r *http.Request, pid string) bool {
	util.LogIt(r, fmt.Sprint("Model - Function - CheckPrivilege"))

	var err error
	Privileges, err := GetRedisHashValueAdmin(r, "privileges")
	if util.CheckErrorLog(r, err) {
		return false
	}

	PArr := strings.Split(Privileges, ",")
	if len(PArr) > 0 {
		for _, n := range PArr {
			if pid == n {
				return true
			}
		}
		return false
	}
	return false
}

// GetRedisHashValueAdmin - Get Redis Single Hash by key  - 2021-05-04 - HK
func GetRedisHashValueAdmin(r *http.Request, key string) (string, error) {
	util.LogIt(r, fmt.Sprint("Model - Function - GetRedisHashValueAdmin"))

	Token := context.Get(r, "Request-Token")
	TokenStr := Token.(string)
	Hash, err := RedisClient.HGet("TP_Admin_Login_"+TokenStr, key).Result()
	if util.CheckErrorLog(r, err) {
		return "", err
	}

	return Hash, nil
}

// GetLocalDateSystem - Returns local date. - 2021-05-11 - HK
func GetLocalDateSystem() string {
	hour, min := GetTimeZoneParamDump()
	localDateTime := time.Now().UTC().Add(time.Hour*time.Duration(hour) + time.Minute*time.Duration(min)).Format("2006-01-02")
	return localDateTime
}

// GetTimeZoneParamDump - 2021-05-11 - HK
func GetTimeZoneParamDump() (int, int) {
	var timezone string
	timezone = GetParameterForDump("time_zone")
	if timezone == "" {
		timezone = "-10:00"
	}
	timeSlice := strings.Split(timezone, ":")
	hour, _ := strconv.Atoi(timeSlice[0])
	min, _ := strconv.Atoi(timeSlice[1])
	return hour, min
}

// GetParameterForDump - 2021-05-11 - HK
func GetParameterForDump(key string) string {
	var Qry bytes.Buffer
	Qry.WriteString("SELECT value FROM cf_parameter WHERE `key` = ?")
	RetMap, err := ExecuteRowQuery(Qry.String(), key)
	if chkDumpError(err) {
		util.SysLogIt("GetParameterForDump - Error Retrieving TimeZone")
		return ""
	}
	if _, ok := RetMap["value"].(string); ok {
		return RetMap["value"].(string)
	}
	return ""
}

// regex for slug
var re = regexp.MustCompile("[^a-z0-9]+")

// slug - to get slug for a string
func Slug(s string) string {
	return strings.Trim(re.ReplaceAllString(strings.ToLower(s), "-"), "-")
}
