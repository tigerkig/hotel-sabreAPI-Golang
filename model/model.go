package model

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"runtime"
	"strconv"
	"strings"
	"time"
	"tp-api-common/data"
	"tp-api-common/util"
	"tp-system/config"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-redis/redis"
	"github.com/gorilla/context"
	gonanoid "github.com/matoous/go-nanoid"
	ua "github.com/mileusna/useragent"
	"github.com/tomasen/realip"
	"gopkg.in/mgo.v2/bson"
)

// RedisClient defined as redis object
var RedisClient *redis.Client

//JQueryTable - Common Structure For JQueryTable
func JQueryTable(r *http.Request, TabelMap data.JQueryTableUI, Qry bytes.Buffer, QryFrom bytes.Buffer, QryCnt bytes.Buffer, QryFilterCnt bytes.Buffer, GroupBy bytes.Buffer, col []map[string]string, colOrder [20]string, status map[string]interface{}, id string) (map[string]interface{}, error) {
	util.LogIt(r, "model - model - JQueryTable")
	var SelectQry, TotalCount, TotalFilterCount, Order, QryLimit, SqlConditionQry, StatusQry, PrimaryIdQry bytes.Buffer
	var MainQuery bytes.Buffer
	var Limit = TabelMap.Limit
	var Offset = TabelMap.Offset
	var search = TabelMap.Search
	SelectQry.WriteString(" SELECT ")

	if len(TabelMap.Order) > 0 {
		var Sort = TabelMap.Order[0]
		var Sort_Field = Sort.Field
		var direction = Sort.Direction
		if direction != "" {
			Order.WriteString("ORDER BY ")
			Order.WriteString("" + colOrder[Sort_Field] + "")
			Order.WriteString(" ")
			Order.WriteString(direction)
		}
	}

	if len(status) > 0 {
		StatusQry.WriteString(" AND '" + status["field"].(string) + status["field"].(string) + "' ")
		StatusQry.WriteString(" ' " + status["value"].(string) + " ' ")
	}

	if id != "" {
		PrimaryIdQry.WriteString(" AND id = ")
		PrimaryIdQry.WriteString(" ' " + id + " ' ")
	}

	if Offset != 0 {
		QryLimit.WriteString(" LIMIT ? OFFSET ?")
	}

	if len(search) > 0 {
		if len(col) > 0 {
			for I, v := range search {
				if I == 0 && v.Value != "" {
					if Flag, _ := InArray(v.Field, col); Flag {
						SqlConditionQry.WriteString(" AND ")
					}
				} else if v.Value != "" {
					if Flag, _ := InArray(v.Field, col); Flag {
						SqlConditionQry.WriteString(" AND ")
					}
				}

				if v.Operator == "begins" && v.Value != "" {
					if Flag, InD := InArray(v.Field, col); Flag {
						SqlConditionQry.WriteString(" " + col[InD]["value"] + " LIKE '" + v.Value + "%' ")
					}
				}
				if v.Operator == "is" && v.Value != "" {
					if Flag, InD := InArray(v.Field, col); Flag {
						SqlConditionQry.WriteString(" " + col[InD]["value"] + " = '" + v.Value + "' ")
					}
				}
				if v.Operator == "contains" && v.Value != "" {
					if Flag, InD := InArray(v.Field, col); Flag {
						SqlConditionQry.WriteString(" " + col[InD]["value"] + " LIKE '%" + v.Value + "%' ")
					}

				}
				if v.Operator == "ends" && v.Value != "" {
					if Flag, InD := InArray(v.Field, col); Flag {
						SqlConditionQry.WriteString(" " + col[InD]["value"] + " LIKE '%" + v.Value + "' ")
					}
				}
				if v.Operator == "date" && v.Value != "" {
					if Flag, InD := InArray(v.Field, col); Flag {
						SqlConditionQry.WriteString(" CAST(" + col[InD]["value"] + " AS DATE) = '" + v.Value + "' ")
					}
				}
				if v.Operator == "sdate" && v.Value != "" {
					if Flag, InD := InArray(v.Field, col); Flag {
						SqlConditionQry.WriteString(" CAST(" + col[InD]["value"] + " AS DATE) >= '" + v.Value + "' ")
					}
				}
				if v.Operator == "edate" && v.Value != "" {
					if Flag, InD := InArray(v.Field, col); Flag {
						SqlConditionQry.WriteString(" CAST(" + col[InD]["value"] + " AS DATE) <= '" + v.Value + "' ")
					}
				}
				if v.Operator == "between" && v.Value != "" {
					var SplitValue = strings.Split(v.Value, ",")
					if Flag, InD := InArray(v.Field, col); Flag {
						SqlConditionQry.WriteString(" CAST(" + col[InD]["value"] + " AS DATE) >= '" + SplitValue[0] + "' AND CAST(" + col[InD]["value"] + " AS DATE) <= '" + SplitValue[0] + "' ")
					}
				}
			}
		} else {
			for I, v := range search {
				if I == 0 && v.Value != "" {
					SqlConditionQry.WriteString(" AND ")
				} else if v.Value != "" {
					SqlConditionQry.WriteString(" AND ")
				}

				if v.Operator == "begins" && v.Value != "" {
					SqlConditionQry.WriteString(" " + v.Field + " LIKE '" + v.Value + "%' ")
				}
				if v.Operator == "is" && v.Value != "" {
					SqlConditionQry.WriteString(" " + v.Field + " = '" + v.Value + "' ")
				}
				if v.Operator == "contains" && v.Value != "" {
					SqlConditionQry.WriteString(" " + v.Field + " LIKE '%" + v.Value + "%' ")
				}
				if v.Operator == "ends" && v.Value != "" {
					SqlConditionQry.WriteString(" " + v.Field + " LIKE '%" + v.Value + "' ")
				}
				if v.Operator == "date" && v.Value != "" {
					SqlConditionQry.WriteString(" CAST(" + v.Field + " AS DATE) = '" + v.Value + "' ")
				}
				if v.Operator == "sdate" && v.Value != "" {
					SqlConditionQry.WriteString(" CAST(" + v.Field + " AS DATE) >= '" + v.Value + "' ")
				}
				if v.Operator == "edate" && v.Value != "" {
					SqlConditionQry.WriteString(" CAST(" + v.Field + " AS DATE) <= '" + v.Value + "' ")
				}
				if v.Operator == "between" && v.Value != "" {
					var SplitValue = strings.Split(v.Value, ",")
					SqlConditionQry.WriteString(" CAST(" + v.Field + " AS DATE) >= '" + SplitValue[0] + "' AND CAST(" + v.Field + " AS DATE) <= '" + SplitValue[0] + "' ")
				}
			}
		}
	}

	TotalCount.WriteString(SelectQry.String())
	TotalCount.WriteString(QryCnt.String())
	TotalCount.WriteString(QryFrom.String())
	TotalCount.WriteString(StatusQry.String())
	TotalCount.WriteString(PrimaryIdQry.String())
	TotalCount.WriteString(GroupBy.String())

	TotalFilterCount.WriteString(SelectQry.String())
	TotalFilterCount.WriteString(QryCnt.String())
	TotalFilterCount.WriteString(QryFrom.String())
	TotalFilterCount.WriteString(StatusQry.String())
	TotalFilterCount.WriteString(PrimaryIdQry.String())
	TotalFilterCount.WriteString(SqlConditionQry.String())
	TotalFilterCount.WriteString(GroupBy.String())

	MainQuery.WriteString(SelectQry.String())
	MainQuery.WriteString(Qry.String())
	MainQuery.WriteString(QryFrom.String())
	MainQuery.WriteString(StatusQry.String())
	MainQuery.WriteString(PrimaryIdQry.String())
	MainQuery.WriteString(SqlConditionQry.String())
	MainQuery.WriteString(GroupBy.String())
	MainQuery.WriteString(Order.String())
	MainQuery.WriteString(QryLimit.String())

	var LocaleData []map[string]interface{}
	stuff := make(map[string]interface{})
	var err error
	util.LogIt(r, MainQuery.String())
	if Offset != 0 {
		LocaleData, err = ExecuteQuery(MainQuery.String(), Offset, Limit)
		if util.CheckErrorLog(r, err) {
			return nil, err
		}
		TotalCntData, err := ExecuteRowQuery(TotalCount.String())
		if util.CheckErrorLog(r, err) {
			return nil, err
		}

		TotalFilterCntData, err := ExecuteRowQuery(TotalFilterCount.String())
		if util.CheckErrorLog(r, err) {
			return nil, err
		}
		stuff["recordsFiltered"] = TotalFilterCntData["cnt"]
		stuff["recordsTotal"] = TotalCntData["cnt"]
	} else {
		LocaleData, err = ExecuteQuery(MainQuery.String())
		if util.CheckErrorLog(r, err) {
			return nil, err
		}
		TotalCntData, err := ExecuteRowQuery(TotalCount.String())
		if util.CheckErrorLog(r, err) {
			return nil, err
		}

		TotalFilterCntData, err := ExecuteRowQuery(TotalFilterCount.String())
		if util.CheckErrorLog(r, err) {
			return nil, err
		}
		stuff["recordsFiltered"] = TotalFilterCntData["cnt"]
		stuff["recordsTotal"] = TotalCntData["cnt"]
	}

	if len(LocaleData) == 0 {
		stuff["data"] = []string{}
	} else {
		stuff["data"] = LocaleData
	}

	return stuff, nil
}

//InArray - Find In Array Data And Pass Flag And Its Index
func InArray(s interface{}, d []map[string]string) (bool, int) {
	for I, v := range d {
		if s == v["key"] {
			//log.Println(I, "III")
			return true, I
		}
	}
	return false, -1
}

// AdminLogin - Backoffice Authentication function
func AdminLogin(r *http.Request, username string, password string) (map[string]interface{}, int) {
	util.LogIt(r, "model - model - AdminLogin")
	userIP := context.Get(r, "Visitor_IP")
	var Qry bytes.Buffer
	Qry.WriteString(" SELECT CU.id, CU.username, CU.privileges, CU.status, CFU.name, CONCAT(phone_code,'',phone) AS phone, CFU.email, CFU.birthdate, CFU.email,CFU.address ")
	Qry.WriteString(" FROM cf_user AS CU ")
	Qry.WriteString(" INNER JOIN cf_user_profile AS CFU ON CFU.user_id = CU.id ")
	Qry.WriteString(" WHERE CU.username = ? AND CU.password= ? AND CU.status = ?")

	AuthMap, err := ExecuteRowQuery(Qry.String(), username, password, "1")
	if util.CheckErrorLog(r, err) {
		return nil, 0
	}

	if AuthMap == nil {
		return nil, 2
	}

	reqMap := make(map[string]interface{})
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.MapClaims{
		"username": username,
		"password": password,
		"exp":      time.Now().Add(time.Hour * 1).Unix(),
		"key":      config.Env.AppKey,
	})

	tokenString, error := token.SignedString([]byte(config.Env.AppKey))
	if error != nil {
		fmt.Println(error)
		return nil, 0
	}

	AuthMap["token"] = tokenString
	AuthMap["ip"] = userIP

	reqMap["token"] = tokenString
	reqMap["user"] = AuthMap["name"]
	chkFlag := SetAuth(r, tokenString, AuthMap)
	if !chkFlag {
		return nil, 0
	}
	if !AddAdminLoginHistory(r, AuthMap["id"].(string)) {
		fmt.Println(error)
		return nil, 0
	}
	return reqMap, 1
}

// AddAdminLoginHistory - Inserts Admin Panel's user's login information - Created By Meet Soni At 18Th April
func AddAdminLoginHistory(r *http.Request, UserID string) bool {
	util.LogIt(r, "Models - model - AddAdminLoginHistory")
	nanoid, _ := gonanoid.Nanoid()
	ip := realip.FromRequest(r)
	getsUserAgentInfo := ua.Parse(r.UserAgent())
	browser := getsUserAgentInfo.Name

	SQLQry := "INSERT INTO cf_backoffice_login_history(id, user_id, login_ip, login_time, user_agent) VALUES (?,?,?,?,?);"
	err := ExecuteNonQuery(SQLQry, nanoid, UserID, ip, util.GetIsoLocalDateTime(), browser)

	if util.CheckErrorLog(r, err) {
		return false
	}
	return true
}

// SetAuth - Set's Authentication details of user
func SetAuth(r *http.Request, token string, AuthMap map[string]interface{}) bool {
	util.LogIt(r, "Models - model - SetAuth")
	err := RedisClient.HMSet("TP_Admin_Login_"+token, AuthMap).Err()
	if util.CheckErrorLog(r, err) {
		return false
	}
	err = RedisClient.Expire("TP_Admin_Login_"+token, 360000*time.Second).Err()
	if util.CheckErrorLog(r, err) {
		return false
	}
	return true
}

// GetAuthDetails - Returns authentcation details by passing token. That is runs in every incoming request
func GetAuthDetails(r *http.Request, token string) (map[string]string, error) {
	AuthRec, err := RedisClient.HGetAll("TP_Admin_Login_" + token).Result()
	if util.CheckErrorLog(r, err) {
		return nil, err
	}
	return AuthRec, nil
}

// PartnerLogin - Partner Panel BackOffice Authentication function
func PartnerMultiHotelLogin(r *http.Request, username string, password string) (map[string]interface{}, int) {
	userIP := context.Get(r, "Visitor_IP")
	var Qry bytes.Buffer
	Qry.WriteString(" SELECT CHC.id, CHC.client_name, CHC.username,CONCAT(CHC.phone_code1,CHC.mobile1) AS mobile1,CONCAT(CHC.phone_code2,CHC.mobile2) AS mobile2, ")
	Qry.WriteString(" CHC.email, CHC.group_id, group_concat(CHI.id) AS hotel_ids,CHC.hotel_id AS default_hotel ")
	Qry.WriteString(" FROM cf_hotel_client AS CHC ")
	Qry.WriteString(" LEFT JOIN cf_hotel_info AS CHI ON CHI.group_id = CHC.group_id  ")
	Qry.WriteString(" WHERE CHC.username = ? AND CHC.password= ? AND CHC.status IN (1,7)")
	Qry.WriteString(" GROUP BY CHC.group_id ")
	AuthMap, err := ExecuteRowQuery(Qry.String(), username, password)
	if util.CheckErrorLog(r, err) {
		return nil, 0
	}

	if AuthMap == nil {
		return nil, 2
	}

	if AuthMap["group_id"] != "" {
		context.Set(r, "GroupId", AuthMap["group_id"])
	}
	HotelInfo := PartnerHotelList(r)

	reqMap := make(map[string]interface{})
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.MapClaims{
		"username": username,
		"password": password,
		"exp":      time.Now().Add(time.Hour * 1).Unix(),
		"key":      config.Env.AppKey,
		"hotel_id": AuthMap["hotel_ids"],
	})

	tokenString, error := token.SignedString([]byte(config.Env.AppKey))
	if util.CheckErrorLog(r, error) {
		fmt.Println(error)
		return nil, 0
	}

	AuthMap["token"] = tokenString
	AuthMap["ip"] = userIP
	AuthMap["hotels"] = HotelInfo

	reqMap["token"] = tokenString
	SetConsoleAuth(r, tokenString, AuthMap)
	// if !AddPartnerLoginHistory(r, AuthMap["id"].(string), "") {
	// 	fmt.Println(error)
	// 	return nil, 0
	// }
	return reqMap, 1
}

func PartnerHotelList(r *http.Request) []map[string]interface{} {
	util.LogIt(r, "model - model - PartnerHotelList")
	var Hotel bytes.Buffer
	Hotel.WriteString(" SELECT CHI.id, hotel_name, CPT.type, hotel_star, is_live, CHI.status AS status_id, ST.status, CHI.is_approved ")
	Hotel.WriteString(" FROM cf_hotel_info AS CHI  ")
	Hotel.WriteString(" INNER JOIN cf_property_type AS CPT ON CPT.id = CHI.property_type_id  ")
	Hotel.WriteString(" INNER JOIN status AS ST ON ST.id = CHI.status  ")
	Hotel.WriteString(" WHERE group_id = ? AND CHI.status IN (1,2,7) ")
	HotelList, err := ExecuteQuery(Hotel.String(), context.Get(r, "GroupId"))
	if util.CheckErrorLog(r, err) {
		return nil
	}

	return HotelList
}

// PartnerLogin - Partner Panel Backoffice Authentication function
func PartnerLogin(r *http.Request, username string, password string) (map[string]interface{}, int) {
	util.LogIt(r, "model - model - PartnerLogin")
	userIP := context.Get(r, "Visitor_IP")
	var Qry bytes.Buffer
	Qry.WriteString(" SELECT CHC.id, CHC.client_name, CHC.username,CONCAT(CHC.phone_code1,CHC.mobile1) AS mobile1,CONCAT(CHC.phone_code2,CHC.mobile2) AS mobile2, ")
	Qry.WriteString(" CHC.email, CFI.hotel_name, CFI.id AS hotel_id, CFI.hotel_star, CFI.latitude, CFI.longitude ")
	Qry.WriteString(" FROM cf_hotel_info AS CFI ")
	Qry.WriteString(" INNER JOIN cf_hotel_client AS CHC ON CHC.hotel_id = CFI.id AND CHC.status = 1 ")
	Qry.WriteString(" WHERE CHC.username = ? AND CHC.password= ? AND CFI.status IN (1,7)")
	AuthMap, err := ExecuteRowQuery(Qry.String(), username, password)
	if util.CheckErrorLog(r, err) {
		return nil, 0
	}

	if AuthMap == nil {
		return nil, 2
	}

	reqMap := make(map[string]interface{})
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.MapClaims{
		"username": username,
		"password": password,
		"exp":      time.Now().Add(time.Hour * 1).Unix(),
		"key":      config.Env.AppKey,
	})

	tokenString, error := token.SignedString([]byte(config.Env.AppKey))
	if error != nil {
		fmt.Println(error)
		return nil, 0
	}

	AuthMap["token"] = tokenString
	AuthMap["ip"] = userIP

	reqMap["token"] = tokenString
	SetConsoleAuth(r, tokenString, AuthMap)
	if !AddPartnerLoginHistory(r, AuthMap["id"].(string), AuthMap["hotel_id"].(string)) {
		fmt.Println(error)
		return nil, 0
	}
	return reqMap, 1
}

// AddPartnerLoginHistory - Inserts Partner Panel's user's login information - Created By Meet Soni At 20Th April
func AddPartnerLoginHistory(r *http.Request, UserID string, hotelID string) bool {
	util.LogIt(r, "Models - model - AddPartnerLoginHistory")
	nanoid, _ := gonanoid.Nanoid()
	ip := realip.FromRequest(r)
	getsUserAgentInfo := ua.Parse(r.UserAgent())
	browser := getsUserAgentInfo.Name

	SQLQry := "INSERT INTO cf_partner_login_history(id, hotel_client_id,hotel_id, login_ip, login_time, user_agent_info) VALUES (?,?,?,?,?,?);"
	err := ExecuteNonQuery(SQLQry, nanoid, UserID, hotelID, ip, util.GetIsoLocalDateTime(), browser)

	if util.CheckErrorLog(r, err) {
		return false
	}
	return true
}

// SetConsoleAuth - Set's Authentication details of user
func SetConsoleAuth(r *http.Request, token string, AuthMap map[string]interface{}) {
	//Hotels := []map[string]interface{}{}
	// if len(AuthMap["hotels"].([]map[string]interface{})) > 0 {
	// 	Hotels = AuthMap["hotels"].([]map[string]interface{})
	// }
	delete(AuthMap, "hotels")
	RedisClient.HMSet("TP_Partner_Login_"+token, AuthMap)
	RedisClient.Expire("TP_Partner_Login_"+token, 360000*time.Second)
	// for _, v := range Hotels {
	// 	RedisClient.SAdd("TP_Partner_Hotels_"+token, v["id"].(string)).Err()
	// }
}

func CheckConsoleHotel(r *http.Request, token, hotel string) bool {
	var tokenArr []string
	if hotel == "" {
		return true
	}
	tokenArr = RedisClient.SMembers("TP_Partner_Hotels_" + token).Val()
	if len(tokenArr) > 0 {
		return InSliceArray(hotel, tokenArr)
	}
	return true
}

func GetConsoleTotalHotel(r *http.Request, token string) (int, string) {
	var defaultHotel string
	var tokenArr []string
	tokenArr = RedisClient.SMembers("TP_Partner_Hotels_" + token).Val()
	if len(tokenArr) == 1 {
		defaultHotel = tokenArr[0]
	}
	return len(tokenArr), defaultHotel
}

//InSliceArray - Find In Array Data And Pass Flag And Its Index
func InSliceArray(s interface{}, d []string) bool {
	for _, a := range d {
		if a == s {
			return true
		}
	}
	return false
}

// GetConsoleAuthDetails - Returns authentcation details by passing token. That is runs in every incoming request
func GetConsoleAuthDetails(r *http.Request, token string) (map[string]string, error) {
	AuthRec, err := RedisClient.HGetAll("TP_Partner_Login_" + token).Result()
	if util.CheckErrorLog(r, err) {
		return nil, err
	}
	return AuthRec, nil
}

// GetConsoleHotelDetails - Returns authentcation details by passing token. That is runs in every incoming request
func GetConsoleHotelDetails(r *http.Request, token string) ([]map[string]interface{}, error) {
	AuthRec, err := RedisClient.HGetAll("TP_Partner_Login_" + token).Result()
	if util.CheckErrorLog(r, err) {
		return nil, err
	}
	HotelMap := []map[string]interface{}{}
	json.Unmarshal([]byte(AuthRec["hotels"]), &HotelMap)
	return HotelMap, nil
}

// GetRedisHashValue - Get Redis Single Hash by key
func GetRedisHashValue(r *http.Request, key string) (string, error) {
	Token := context.Get(r, "Request-Token")
	TokenStr := Token.(string)
	Hash, err := RedisClient.HGet("TP_Partner_Login_"+TokenStr, key).Result()
	if util.CheckErrorLog(r, err) {
		return "", err
	}
	return Hash, nil
}

// Logout - Logout method removes token
func Logout(r *http.Request) bool {
	token := context.Get(r, "Request-Token").(string)
	Panel := context.Get(r, "Side").(string)
	if Panel == "TP-BACKOFFICE" {
		RedisClient.Del("TP_Admin_Login_" + token)
	} else if Panel == "TP-PARTNER" {
		RedisClient.Del("TP_Partner_Login_" + token)
	}

	return true
}

// StatusList - - Returns List of status
func StatusList(r *http.Request) []map[string]interface{} {
	util.LogIt(r, "model - model - StatusList")
	var Qry bytes.Buffer
	Qry.WriteString(" SELECT id, status FROM status WHERE id IN (1,2)")
	AuthMap, err := ExecuteQuery(Qry.String())
	if util.CheckErrorLog(r, err) {
		return nil
	}

	return AuthMap
}

// GetLogsValueMap - Return Map Of Necessary Logs Value
func GetLogsValueMap(r *http.Request, reqMap interface{}, flag bool, removeKey string) map[string]interface{} {
	if flag {
		var SplitString = strings.Split(removeKey, ",")
		for _, val := range SplitString {
			delete(reqMap.(map[string]interface{}), val)
		}
	}

	for k, v := range reqMap.(map[string]interface{}) {
		delete(reqMap.(map[string]interface{}), "ID")
		if v == "" {
			delete(reqMap.(map[string]interface{}), k)
		}
	}

	return reqMap.(map[string]interface{})
}

//ShortDateFromString parse shot date from string - 2020-08-04 - HK
func ShortDateFromString(ds string) (time.Time, error) {
	const shortDate = "2006-01-02"

	t, err := time.Parse(shortDate, ds)
	if err != nil {
		return t, err
	}
	return t, nil
}

//CheckDataBoundariesStr checks is startdate <= enddate - 2020-08-04 - HK
func CheckDataBoundariesStr(startdate, enddate string) (bool, error) {

	tstart, err := ShortDateFromString(startdate)
	if err != nil {
		return false, fmt.Errorf("cannot parse startdate: %v", err)
	}

	tend, err := ShortDateFromString(enddate)
	if err != nil {
		return false, fmt.Errorf("cannot parse enddate: %v", err)
	}

	dt := time.Now()
	tcurrent, err := ShortDateFromString(dt.Format("2006-01-02"))
	if err != nil {
		return false, fmt.Errorf("cannot get current date: %v", err)
	}

	// checks if start date and end date are before current date or not, if yes returns false
	if tstart.Before(tcurrent) || tend.Before(tcurrent) {
		return false, fmt.Errorf("Invalid Start Date / End Date")
	}

	// checks if start date is after end date or not, if yes returns false
	if tstart.After(tend) {
		return false, fmt.Errorf("Start Date Is Greater Than End Date")
	}

	return true, err
}

// GetHotelListForDumpData - Gets Hotel List For Dumping inv, rates, restrictions data - 2021-05-07 - HK
func GetHotelListForDumpData() ([]map[string]interface{}, error) {
	util.SysLogIt("GetHotelListForDumpData Start")

	var Qry bytes.Buffer
	Qry.WriteString(" SELECT")
	Qry.WriteString(" CHI.id, CHI.hotel_name ")
	Qry.WriteString(" FROM ")
	Qry.WriteString(" cf_hotel_info AS CHI ")
	Qry.WriteString(" LEFT JOIN ")
	Qry.WriteString(" cf_hotel_client AS CHC ON CHC.group_id = CHI.group_id AND CHC.group_id <> '' ")
	Qry.WriteString(" INNER JOIN ")
	Qry.WriteString(" cf_hotel_settings AS CHS ON CHS.hotel_id = CHI.id ")
	Qry.WriteString(" INNER JOIN ")
	Qry.WriteString(" status AS ST ON ST.id = CHI.status ")
	Qry.WriteString(" WHERE ")
	Qry.WriteString(" CHI.status <> 3 ")
	// Qry.WriteString(" AND CHI.id = 'YkjA-EF6tUHT5n9ntIGCD' ") // remove / comment this before push

	RetMap, err := ExecuteQuery(Qry.String())
	if chkDumpError(err) {
		util.SysLogIt("GetHotelListForDumpData - Error Retrieving Hotel List")
		return nil, err
	}

	util.SysLogIt("GetHotelListForDumpData End")
	return RetMap, nil
}

// CheckRoomDataCountForDumpData - Checks Room Data Count For Hotel - 2021-05-07 - HK
func CheckRoomDataCountForDumpData(HotelID string) ([]map[string]interface{}, error) {
	util.SysLogIt("CheckRoomDataCountForDumpData Start")

	var Qry bytes.Buffer
	Qry.WriteString(" SELECT id, room_type_name, max_occupancy, inventory FROM cf_room_type WHERE hotel_id = ? AND status = ?")
	// Qry.WriteString(" AND id = '9yT_K5BN2bUWvzqfSmWOQ' ")
	Data, err := ExecuteQuery(Qry.String(), HotelID, 1)
	if chkDumpError(err) {
		util.SysLogIt("CheckRoomDataCountForDumpData - Error Retrieving Hotel Room Data")
		return nil, err
	}

	util.SysLogIt("CheckRoomDataCountForDumpData End")
	return Data, nil
}

// CheckRatePlanDataCountForDumpData - Checks Rate Plan Data Count For Provided Hotel And Room - 2021-05-07 - HK
func CheckRatePlanDataCountForDumpData(HotelID string, RoomID string) ([]map[string]interface{}, error) {
	util.SysLogIt("CheckRatePlanDataCountForDumpData Start")

	var Qry bytes.Buffer
	Qry.WriteString("SELECT id, rate_plan_name, rate FROM cf_rateplan WHERE hotel_id = ? AND room_type_id = ? AND status = ?")
	Data, err := ExecuteQuery(Qry.String(), HotelID, RoomID, 1)
	if chkDumpError(err) {
		util.SysLogIt("CheckRatePlanDataCountForDumpData - Error Retrieving Hotel Rate Plan Data")
		return nil, err
	}

	util.SysLogIt("CheckRatePlanDataCountForDumpData End")
	return Data, nil
}

// chkDumpError - Insertrs Error While Dumping Data - 2021-05-07 - HK
func chkDumpError(err error) bool {
	if err != nil {
		_, file, no, ok := runtime.Caller(1)
		if ok {
			content := fmt.Sprint("Exception -  on File - ", file, " Line - ", no)
			c := VMongoSession.DB(config.Env.Mongo.MongoDB).C("dump_error")
			MID := bson.NewObjectId()
			c.Insert(&data.ErrLog{ID: MID, Content: content + " /n " + err.Error()})
		}
		return true
	}
	return false
}

// MonthIntervalDump - Gets First And Last Date Of Month Based On Passed Arguments Year And Month
// For e.g. if we pass 2021 as year and 05 as month then it will return 2021-05-01, 2021-05-31
// 2021-05-08 - HK
func MonthIntervalDump(y int, m time.Month) (firstDay, lastDay time.Time) {
	firstDay = time.Date(y, m, 1, 0, 0, 0, 0, time.UTC)
	lastDay = time.Date(y, m+1, 1, 0, 0, 0, -1, time.UTC)
	return firstDay, lastDay
}

// Get365Plus1DayDate - To get Passed Year, Month's First Date And Date Of 366th day from today
// Considering passed arguments year and month is 2021 07 Then it will return  2021-07-01
// consideringto day date as 2021-05-08 will return 2022-05-09
// 2021-05-08 - HK
func Get365Plus1DayDate(fromYear int, fromMonth int) (string, string) {

	first, _ := MonthIntervalDump(fromYear, time.Month(fromMonth))
	startDate := first.Format("2006-01-02")

	t := time.Now()
	endDate := t.AddDate(0, 0, 366).Format("2006-01-02")
	return startDate, endDate
}

// DaysBetweenForDataDump - returns difference of days between 2 dates - 2021-05-08 - HK
func DaysBetweenForDataDump(a, b time.Time) int {
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

// DateTimeForDataDump - Returns Date Into Time - 2021-05-08 - HK
func DateTimeForDataDump(s string) time.Time {
	d, _ := time.Parse("2006-01-02", s)
	return d
}

// GetYearMonthSliceBetweenTwoDatesForDataDump - 2021-05-08 - HK
func GetYearMonthSliceBetweenTwoDatesForDataDump(startDate string, endDate string) []map[string]interface{} {
	noOfDays := DaysBetweenForDataDump(DateTimeForDataDump(endDate), DateTimeForDataDump(startDate))

	yearMonth := []string{}
	fullDate := []string{}
	for i := 0; i <= noOfDays; i++ {
		insertDate := DateTimeForDataDump(startDate).AddDate(0, 0, i).Format("2006-01-02")
		dateIno := strings.Split(insertDate, "-")
		yearMonthString := dateIno[0] + "-" + dateIno[1]
		yearMonth = append(yearMonth, yearMonthString)
		fullDate = append(fullDate, insertDate)
	}

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
	// [{"data":{"2020-12-31":"1"},"year_month":"2020-12"},{"data":{"2021-01-01":"1","2021-01-02":"1","2021-01-03":"1","2021-01-04":"1","2021-01-05":"1","2021-01-06":"1","2021-01-07":"1","2021-01-08":"1","2021-01-09":"1","2021-01-10":"1","2021-01-11":"1","2021-01-12":"1","2021-01-13":"1","2021-01-14":"1","2021-01-15":"1","2021-01-16":"1","2021-01-17":"1","2021-01-18":"1","2021-01-19":"1","2021-01-20":"1","2021-01-21":"1","2021-01-22":"1","2021-01-23":"1","2021-01-24":"1","2021-01-25":"1","2021-01-26":"1","2021-01-27":"1","2021-01-28":"1","2021-01-29":"1","2021-01-30":"1","2021-01-31":"1"},"year_month":"2021-01"},{"data":{"2021-02-01":"1"},"year_month":"2021-02"}]
	return yearMonthint
}

// MeasureTime - logs timing for data dumping and syncing - 2021-05-11 - HK
func MeasureTime(funcName string) func() {
	start := time.Now()
	return func() {
		// fmt.Printf("Time taken by %s function is %v \n", funcName, time.Since(start))
		util.SysLogIt(" Time taken by " + funcName + " function is : ")
		util.SysLogIt(time.Since(start))
	}
}

// FillInvRateRestData - Dump Inv Rate Restrictions Data For Each Hotel For Each 366th day
// with default respected data - 2021-05-08 - HK
func FillInvRateRestData() bool {
	util.SysLogIt("FillInvRateRestData Start")
	defer MeasureTime("FillInvRateRestData")()

	HotelList, err := GetHotelListForDumpData()
	if chkDumpError(err) {
		util.SysLogIt("FillInvRateRestData - Error Retrieving Hotel List")
		return false
	}

	if len(HotelList) > 0 {
		util.SysLogIt("FillInvRateRestData - Hotels Found ")

		// loop through all the hotels to fill inv, rates, restrictions data
		// for each present hotel and having proper room and rate plans
		for _, v := range HotelList {

			HotelID := v["id"].(string)
			HotelName := v["hotel_name"]
			util.SysLogIt(" Turn For hotel_name ")
			util.SysLogIt(HotelName)

			var finalArr = make(map[string]interface{})
			finalArr["hotel_id"] = HotelID
			finalArr["hotel_name"] = HotelName

			RoomListData, err := CheckRoomDataCountForDumpData(HotelID)
			if chkDumpError(err) {
				util.SysLogIt("FillInvRateRestData - Error Retrieving Hotel Room Data Count")
				continue
			}
			// j1, _ := json.Marshal(RoomListData)
			// util.SysLogIt(string(j1))

			// consider only those hotels who have rooms exists
			if len(RoomListData) > 0 {

				// Only Those Rooms Are Considered Whose Rate Plan Exists
				var FinalRoomListData []map[string]interface{}
				for _, j := range RoomListData {
					rateDATA, err := CheckRatePlanDataCountForDumpData(HotelID, j["id"].(string))
					if len(rateDATA) == 0 || err != nil {
						util.SysLogIt("FillInvRateRestData - No Rate Plans Found For This Room")
						util.SysLogIt(j["room_type_name"].(string))
					}
					if len(rateDATA) > 0 {
						FinalRoomListData = append(FinalRoomListData, j)
					}
				}

				if len(FinalRoomListData) > 0 {

					// var finalArr = make(map[string]interface{})
					var RoomPlanArr []map[string]interface{}

					for k1, v1 := range FinalRoomListData {

						RoomID := v1["id"].(string)
						RoomName := v1["room_type_name"].(string)
						ForMaxOcc := v1["max_occupancy"].(int64)
						RoomInv := v1["inventory"].(int64)

						var Qry bytes.Buffer
						Qry.WriteString(" SELECT ")
						Qry.WriteString(" max(year) as year, max(month) as month ")
						Qry.WriteString(" FROM ")
						Qry.WriteString(" cf_inv_data ")
						Qry.WriteString(" WHERE ")
						Qry.WriteString(" hotel_id = ? AND room_id = ? ")
						Qry.WriteString(" GROUP BY year, month ")
						Qry.WriteString(" ORDER BY year DESC, month DESC ")
						Qry.WriteString(" LIMIT 0, 1; ")
						Data, err := ExecuteRowQuery(Qry.String(), HotelID, RoomID)
						if chkDumpError(err) {
							util.SysLogIt("FillInvRateRestData - Error Retrieving Hotel Room Inv Data")
							continue
						}

						// startdate, enddate := Get365Plus1DayDate(2021, 07)
						startdate, enddate := Get365Plus1DayDate(int(Data["year"].(int64)), int(Data["month"].(int64)))
						/* START */

						rateDATA, err := CheckRatePlanDataCountForDumpData(HotelID, RoomID)
						if chkDumpError(err) {
							util.SysLogIt("FillInvRateRestData - Error Retrieving Hotel Room Data Count")
							continue
						}

						RoomPlanArr = append(RoomPlanArr, map[string]interface{}{
							"room_id":    RoomID,
							"room_name":  RoomName,
							"occupancy":  ForMaxOcc,
							"inventory":  RoomInv,
							"start_date": startdate,
							"end_date":   enddate,
						})

						var RatePlanArr []map[string]interface{}
						for _, v2 := range rateDATA {

							var OccupancyArr []map[string]interface{}
							var j int64
							for j = 1; j <= ForMaxOcc; j++ {
								joinStr := strconv.FormatInt(j, 16) //
								mainStr := "occ_" + joinStr
								OccupancyArr = append(OccupancyArr, map[string]interface{}{
									mainStr: v2["rate"],
								})
							}

							RatePlanArr = append(RatePlanArr, map[string]interface{}{
								"rate_id":    v2["id"],
								"rate_name":  v2["rate_plan_name"],
								"rate":       OccupancyArr,
								"start_date": startdate,
								"end_date":   enddate,
								"min_night":  1,
								"stop_sell":  0,
								"cta":        0,
								"ctd":        0,
							})
						}
						RoomPlanArr[k1]["rate_info"] = RatePlanArr

						// log.Println(RoomPlanArr)
						finalArr["room_info"] = RoomPlanArr
						/* END */
					} // end for k1, v1 := range RoomListData

					j, _ := json.Marshal(finalArr)
					util.SysLogIt(string(j))

					// prepared json needs to be sent for further operations start
					flg := DumpInvRateRestData(finalArr)
					if flg {
						util.SysLogIt("FillInvRateRestData - Inv, Rates, Restrictions Data Dumped Successfully For")
						util.SysLogIt(HotelName)
					} else {
						util.SysLogIt("FillInvRateRestData - Error Dumping Inv, Rates, Restrictions Data For")
						util.SysLogIt(HotelName)
					}
					// prepared json needs to be sent for further operations end

				} // end if len(FinalRoomListData) > 0 {
			} else {
				util.SysLogIt("FillInvRateRestData - No rooms found for hotel")
				util.SysLogIt(v["hotel_name"])
			}

			util.SysLogIt("===============================================")
		} // end for _, v := range HotelList

	} else {
		util.SysLogIt("FillInvRateRestData - No Hotels Found ")
	}

	util.SysLogIt("FillInvRateRestData End")
	return true
}

// CheckAndFillInvDataDump - Checks If Data Exists in DB and If not fills data for inventory - 2021-05-08 - HK
func CheckAndFillInvDataDump(HotelID string, roomID string, baseInv int64, startDate string, endDate string) bool {
	util.SysLogIt("CheckAndFillInvDataDump Start")

	yearMonthDates := GetYearMonthSliceBetweenTwoDatesForDataDump(startDate, endDate)
	j, _ := json.Marshal(yearMonthDates)
	util.SysLogIt("CheckAndFillInvDataDump - yearMonthDates")
	util.SysLogIt(string(j))

	/* Update Log Purpose */
	localDate := util.GetISODate()

	for _, val := range yearMonthDates {

		util.SysLogIt("CheckAndFillInvDataDump - Turn For Room Type")
		util.SysLogIt(roomID)

		yearMonth := fmt.Sprintf("%v", val["year_month"]) // Convert interface to string
		dateData := val["data"]

		dateInfo := strings.Split(yearMonth, "-")
		tblYear := dateInfo[0]
		tblMonth := dateInfo[1]

		var Qry bytes.Buffer
		Qry.WriteString("SELECT inv_data FROM cf_inv_data WHERE hotel_id = ? AND room_id = ? AND year = ? AND month = ? ") // Add
		Data, err := ExecuteRowQuery(Qry.String(), HotelID, roomID, tblYear, tblMonth)
		if chkDumpError(err) {
			util.SysLogIt("CheckAndFillInvDataDump - Error Retrieving Hotel Room Data For")
			util.SysLogIt(HotelID + "-" + roomID + "-" + tblYear + "-" + tblMonth)
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

				// Put Log Here - For Update Log
				var Qry1 bytes.Buffer
				nanoid, _ := gonanoid.Nanoid()
				Qry1.WriteString("INSERT INTO logs_inv (id, hotel_id, room_id, inventory, update_for_date, updated_at, updated_by, booking_id, ip) VALUES (?,?,?,?,?,?,?,?,?)")
				err = ExecuteNonQuery(Qry1.String(), nanoid, HotelID, roomID, baseInv, realDate, localDate, "", "", "")
				if chkDumpError(err) {
					util.SysLogIt("CheckAndFillInvDataDump - Error Inserting Hotel Room Data For")
					util.SysLogIt(HotelID + "-" + roomID + "-" + realDate)
					return false
				}
				// Put Log Here - For Update Log
			}
			invJSON, err := json.Marshal(DateArr)

			util.SysLogIt("if data == 0 invJSON")
			util.SysLogIt(string(invJSON))

			var Qry1 bytes.Buffer
			nanoid, _ := gonanoid.Nanoid()
			Qry1.WriteString("INSERT INTO cf_inv_data (id, hotel_id, room_id, year, month, inv_data) VALUES (?,?,?,?,?,?)")
			err = ExecuteNonQuery(Qry1.String(), nanoid, HotelID, roomID, tblYear, tblMonth, string(invJSON))
			if chkDumpError(err) {
				util.SysLogIt("CheckAndFillInvDataDump - Error Inserting Hotel Room Data For")
				util.SysLogIt(HotelID + "-" + roomID + "-" + tblYear + "-" + tblMonth)
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

					// Put Log Here For Update Log
					var Qry1 bytes.Buffer
					nanoid, _ := gonanoid.Nanoid()
					Qry1.WriteString("INSERT INTO logs_inv (id, hotel_id, room_id, inventory, update_for_date, updated_at, updated_by, booking_id, ip) VALUES (?,?,?,?,?,?,?,?,?)")
					err = ExecuteNonQuery(Qry1.String(), nanoid, HotelID, roomID, baseInv, realDate, localDate, "", "", "")
					if chkDumpError(err) {
						util.SysLogIt("CheckAndFillInvDataDump - Error Inserting Hotel Room Data For")
						util.SysLogIt(HotelID + "-" + roomID + "-" + realDate)
						return false
					}
					// Put Log Here For Update Log
				}
			}
			invJSON, err := json.Marshal(result)
			// fmt.Println(string(invJSON), err)

			util.SysLogIt("else invJSON For Update")
			util.SysLogIt(string(invJSON))

			var Qry1 bytes.Buffer
			Qry1.WriteString("UPDATE cf_inv_data SET inv_data = ? WHERE  hotel_id = ? AND  room_id = ? AND year = ? AND month = ?")
			err = ExecuteNonQuery(Qry1.String(), string(invJSON), HotelID, roomID, tblYear, tblMonth)
			if chkDumpError(err) {
				util.SysLogIt("CheckAndFillInvDataDump - Error Updating Hotel Room Data For")
				util.SysLogIt(HotelID + "-" + roomID + "-" + tblYear + "-" + tblMonth)
				return false
			}
		}
	}
	util.SysLogIt("CheckAndFillInvDataDump End")
	return true
}

// CheckAndFillRateRestrictionDataDump - Checks If Data Exists in DB and If not fills data for rates and restrictions - 2021-05-08 - HK
func CheckAndFillRateRestrictionDataDump(HotelID string, roomID string, rateInfo map[string]interface{}) bool {
	util.SysLogIt("CheckAndFillRateRestrictionDataDump Start")

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

	yearMonthDates := GetYearMonthSliceBetweenTwoDatesForDataDump(startDate, endDate)
	/* 2020-08-19 - Log Purpose */
	localDate := util.GetISODate()

	for _, val := range yearMonthDates {

		util.SysLogIt("CheckAndFillInvDataDump - Turn For Room Type And Rate Plan")
		util.SysLogIt(roomID + " - " + rateID)

		// insert rate, restriction data
		yearMonth := fmt.Sprintf("%v", val["year_month"]) // Convert interface to string

		dateData := val["data"]

		dateInfo := strings.Split(yearMonth, "-")
		tblYear := dateInfo[0]
		tblMonth := dateInfo[1]

		var Qry bytes.Buffer
		Qry.WriteString("SELECT rate_rest_data FROM cf_rate_restriction_data WHERE hotel_id = ? AND room_id = ? AND rateplan_id = ? AND year = ? AND month = ? ")
		RateRestDataFromTbl, err := ExecuteRowQuery(Qry.String(), HotelID, roomID, rateID, tblYear, tblMonth)
		if chkDumpError(err) {
			util.SysLogIt("CheckAndFillRateRestrictionDataDump - Error Retrieving Hotel RatePlan Data For")
			util.SysLogIt(HotelID + "-" + roomID + "-" + rateID + "-" + tblYear + "-" + tblMonth)
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

				// Put Log Here - For Update Log
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
				if chkDumpError(err) {
					util.SysLogIt("CheckAndFillRateRestrictionDataDump - If Error Inserting Hotel RatePlan Data For")
					util.SysLogIt(HotelID + "-" + roomID + "-" + rateID + "-" + realDate)
					return false
				}
				// Put Log Here - For Update Log
			}

			rateRestJSON, err := json.Marshal(DateDataInsert)
			util.SysLogIt("if data == 0 rateRestJSON")
			util.SysLogIt(string(rateRestJSON))

			var Qry1 bytes.Buffer
			nanoid, _ := gonanoid.Nanoid()
			Qry1.WriteString("INSERT INTO cf_rate_restriction_data (id, hotel_id, room_id, rateplan_id, year, month, rate_rest_data)  VALUES (?,?,?,?,?,?,?)")
			err = ExecuteNonQuery(Qry1.String(), nanoid, HotelID, roomID, rateID, tblYear, tblMonth, string(rateRestJSON))
			if chkDumpError(err) {
				util.SysLogIt("CheckAndFillRateRestrictionDataDump - If Error Inserting Hotel RatePlan Data For")
				util.SysLogIt(HotelID + "-" + roomID + "-" + rateID + "-" + tblYear + "-" + tblMonth)
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

					// Put Log Here - For Update Log
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
					if chkDumpError(err) {
						util.SysLogIt("CheckAndFillRateRestrictionDataDump - Else Error Inserting Hotel RatePlan Data For")
						util.SysLogIt(HotelID + "-" + roomID + "-" + rateID + "-" + realDate)
						return false
					}
					// Put Log Here - For Update Log
				}
			}

			rateJSON, err := json.Marshal(result)

			util.SysLogIt("else rateJSON For Update")
			util.SysLogIt(string(rateJSON))

			var Qry1 bytes.Buffer
			Qry1.WriteString("UPDATE cf_rate_restriction_data SET rate_rest_data = ? WHERE  hotel_id = ? AND  room_id = ? AND rateplan_id = ? AND year = ? AND month = ?")
			err = ExecuteNonQuery(Qry1.String(), string(rateJSON), HotelID, roomID, rateID, tblYear, tblMonth)
			if chkDumpError(err) {
				util.SysLogIt("CheckAndFillRateRestrictionDataDump - Else Error Inserting Hotel RatePlan Data For")
				util.SysLogIt(HotelID + "-" + roomID + "-" + rateID + "-" + tblYear + "-" + tblMonth)
				return false
			}
		}
	} // end for _, val := range yearMonthDates
	util.SysLogIt("CheckAndFillRateRestrictionDataDump End")
	return true
}

// DumpInvRateRestData - Fills Inv, Rate, Restrictions Data For Hotel - 2021-05-08 - HK
func DumpInvRateRestData(reqMap map[string]interface{}) bool {
	util.SysLogIt("DumpInvRateRestData Start")

	HotelID := reqMap["hotel_id"].(string)
	roomInfo := reqMap["room_info"].([]map[string]interface{})

	if len(roomInfo) > 0 {

		for i := 0; i < len(roomInfo); i++ {

			roomID := roomInfo[i]["room_id"].(string)
			baseInv := roomInfo[i]["inventory"].(int64)
			startDate := roomInfo[i]["start_date"].(string)
			endDate := roomInfo[i]["end_date"].(string)

			// log.Println(HotelID, roomID, baseInv, startDate, endDate)
			invSuccess := CheckAndFillInvDataDump(HotelID, roomID, baseInv, startDate, endDate)
			if !invSuccess {
				util.SysLogIt("DumpInvRateRestData - Error Filling Room Inventory")
				util.SysLogIt(HotelID + "-" + reqMap["hotel_name"].(string) + "-" + roomID + "-" + startDate + "-" + endDate)
				return false
			}

			rateInfo := roomInfo[i]["rate_info"].([]map[string]interface{})
			if len(rateInfo) > 0 {
				for i := 0; i < len(rateInfo); i++ {
					occWiseRateInfoData := rateInfo[i]
					rateSuccess := CheckAndFillRateRestrictionDataDump(HotelID, roomID, occWiseRateInfoData)
					if !rateSuccess {
						util.SysLogIt("DumpInvRateRestData - Error Filling Rate Restrictions Data")
						util.SysLogIt(HotelID + "-" + reqMap["hotel_name"].(string) + "-" + roomID + "-" + startDate + "-" + endDate)
						return false
					}

					rateID := occWiseRateInfoData["rate_id"].(string)
					CacheChn <- CacheObj{
						Type:        "updateDeals",
						ID:          HotelID,
						Additional:  roomID,
						Additional1: rateID,
					}

				}
			}
		} // end for i := 0; i < len(roomInfo); i++
	}
	util.SysLogIt("DumpInvRateRestData End")
	return true
}
