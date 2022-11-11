package front

import (
	"bytes"
	"fmt"
	"net/http"
	"strconv"
	"time"
	"tp-api-common/data"
	"tp-api-common/util"
	"tp-system/config"
	"tp-system/model"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/context"
	gonanoid "github.com/matoous/go-nanoid"
)

// AddListingInquiry - Adds Property Listing Inquiry
func AddListingInquiry(r *http.Request, reqMap data.ListYourProperty) bool {
	util.LogIt(r, "Model - Inquiry - AddListingInquiry")
	var Qry, AdmQry bytes.Buffer
	nanoid, _ := gonanoid.Nanoid()
	InquirerIP := context.Get(r, "Visitor_IP").(string)

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.MapClaims{
		"email": reqMap.Email,
		"exp":   time.Now().Add(time.Hour * 48).Unix(),
	})

	tokenString, error := token.SignedString([]byte(config.Env.AppKey))
	if util.CheckErrorLog(r, error) {
		fmt.Println(error)
		return false
	}

	Qry.WriteString(" INSERT INTO cf_list_your_property SET id=?, first_name=?,last_name=?, phone_code=?, phone =?, email =?, password=?, token = ?, token_expired_at = ?, inquiry_at =?, ip = ?")
	err := model.ExecuteNonQuery(Qry.String(), nanoid, reqMap.FirstName, reqMap.LastName, reqMap.PhoneCode, reqMap.Phone, reqMap.Email, util.GeneratePasswordHash(reqMap.Password), tokenString, util.GetExpiredISODateTime(2880), util.GetIsoLocalDateTime(), InquirerIP)
	if util.CheckErrorLog(r, err) {
		return false
	}

	CompanyInfo, _ := model.GetCompanyInfo(r, "1")
	resMap := make(map[string]interface{})
	resMap["logo"] = CompanyInfo["image"].(string)
	resMap["activation_link"] = config.Env.FrontWebURL + "activation/" + tokenString

	model.MailChn <- model.MailObj{
		Type:         "EmailTemplate",
		ID:           "7",
		Additional:   reqMap.Email,
		InterfaceObj: resMap,
	}

	// Mail Reminder to Admin for Partner Registration
	var phoneCode = strconv.FormatFloat(reqMap.PhoneCode, 'f', 0, 64)
	var phoneNo = strconv.FormatFloat(reqMap.Phone, 'f', 0, 64)

	resMap["first_name"] = reqMap.FirstName
	resMap["last_name"] = reqMap.LastName
	resMap["phone_no"] = phoneCode + phoneNo
	resMap["email_id"] = reqMap.Email

	// Admin Data for getting admin email
	AdmQry.WriteString(" SELECT email FROM cf_user_profile WHERE user_id = 1")
	AdminData, err := model.ExecuteRowQuery(AdmQry.String())
	if util.CheckErrorLog(r, err) {
		return false
	}

	model.MailChn <- model.MailObj{
		Type:         "EmailTemplate",
		ID:           "9",
		Additional:   AdminData["email"].(string),
		InterfaceObj: resMap,
	}

	return true
}

func CheckPartnerAlreadyRegister(r *http.Request, email, flag string) int64 {
	util.LogIt(r, "Model - Inquiry - CheckPartnerAlreadyRegister")
	var Qry bytes.Buffer
	Qry.WriteString(" SELECT ")
	Qry.WriteString(" CASE  ")
	Qry.WriteString(" WHEN COUNT(id) = 0 THEN 1  ")
	Qry.WriteString(" WHEN is_active = 1 THEN 2  ")
	Qry.WriteString(" WHEN is_active = 0 AND NOW() > FROM_UNIXTIME(token_expired_at) THEN 3 ")
	Qry.WriteString(" WHEN is_active = 0 AND FROM_UNIXTIME(token_expired_at) > NOW() THEN 4 ")
	Qry.WriteString(" ELSE 5 ")
	Qry.WriteString(" END AS is_token_exist ")
	Qry.WriteString(" FROM cf_list_your_property ")
	Qry.WriteString(" WHERE 1 = 1 ")
	if flag == "EMAIL" {
		Qry.WriteString(" AND email = ? ")
	} else if flag == "TOKEN" {
		Qry.WriteString(" AND token = ? ")
	}
	Data, err := model.ExecuteRowQuery(Qry.String(), email)
	if util.CheckErrorLog(r, err) {
		return 0
	}

	return Data["is_token_exist"].(int64)

}

func CheckPartnerTokenIsValid(r *http.Request, reqMap data.VerifyPartnerToken) bool {
	util.LogIt(r, "Model - Inquiry - AddListingInquiry")
	var Qry, GrpQry, UpdateGrpQry, UpdateFlag bytes.Buffer
	nanoID, _ := gonanoid.Nanoid()
	grpID, _ := gonanoid.Nanoid()

	Qry.WriteString(" INSERT INTO cf_hotel_client(id,client_name,username,password,phone_code1,mobile1,email,created_at,created_by) SELECT ?,CONCAT(first_name,' ',last_name),email,password,phone_code,phone,email,?,1 FROM cf_list_your_property where token = ? ")
	err := model.ExecuteNonQuery(Qry.String(), nanoID, util.GetIsoLocalDateTime(), reqMap.Token)
	if util.CheckErrorLog(r, err) {
		return false
	}

	GrpQry.WriteString(" INSERT INTO cf_hotel_group SET id = ?, client_id = ?, created_at = ? ")
	err = model.ExecuteNonQuery(GrpQry.String(), grpID, nanoID, util.GetIsoLocalDateTime())
	if util.CheckErrorLog(r, err) {
		return false
	}

	UpdateGrpQry.WriteString(" UPDATE cf_hotel_client SET group_id = ? WHERE id = ? ")
	err = model.ExecuteNonQuery(UpdateGrpQry.String(), grpID, nanoID)
	if util.CheckErrorLog(r, err) {
		return false
	}

	UpdateFlag.WriteString(" UPDATE cf_list_your_property SET is_active = 1 WHERE token = ? ")
	err = model.ExecuteNonQuery(UpdateFlag.String(), reqMap.Token)
	if util.CheckErrorLog(r, err) {
		return false
	}

	//go routine for sending mail
	go SendWelcomeMailToPartner(r, reqMap.Token)

	return true
}

var SendWelcomeMailToPartner = func(r *http.Request, token string) bool {
	var Qry bytes.Buffer
	Qry.WriteString(" SELECT email FROM cf_list_your_property WHERE token = ? ")
	Email, err := model.ExecuteRowQuery(Qry.String(), token)
	if util.CheckErrorLog(r, err) {
		return false
	}
	if _, ok := Email["email"].(string); ok {
		CompanyInfo, _ := model.GetCompanyInfo(r, "1")
		resMap := make(map[string]interface{})
		resMap["logo"] = CompanyInfo["image"].(string)
		resMap["email"] = Email["email"].(string)
		resMap["company"] = CompanyInfo["company_name"].(string)
		resMap["partner_link"] = config.Env.PartnerWebURL
		resMap["website_url"] = config.Env.FrontWebURL

		model.MailChn <- model.MailObj{
			Type:         "EmailTemplate",
			ID:           "8",
			Additional:   Email["email"].(string),
			InterfaceObj: resMap,
		}

	}

	return true
}

func VerifyJWTToken(r *http.Request, reqMap data.VerifyPartnerToken) bool {
	token, err := jwt.Parse(reqMap.Token, func(token *jwt.Token) (interface{}, error) {
		//Make sure that the token method conform to "SigningMethodHMAC"
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(config.Env.AppKey), nil
	})

	if util.CheckErrorLog(r, err) {
		return false
	}

	if _, ok := token.Claims.(jwt.Claims); !ok && !token.Valid {
		return false
	} else {
		if !CheckPartnerTokenIsValid(r, reqMap) {
			return false
		}
	}
	return true
}

// ListYourPropertyListing - Return Datatable Listing Of List Your Property
func ListYourPropertyListing(r *http.Request, reqMap data.JQueryTableUI) (map[string]interface{}, error) {
	util.LogIt(r, "model - Inquiry - ListYourPropertyListing")
	var Qry, QryCnt, QryFilter, GroupBy, FromQry bytes.Buffer
	var status = make(map[string]interface{})

	var testColArrs [20]string
	testColArrs[0] = "id"
	testColArrs[1] = "name"
	testColArrs[2] = "phone"
	testColArrs[3] = "email"
	testColArrs[4] = "is_active"
	testColArrs[5] = "ip"
	testColArrs[6] = "inquiry_at"

	var testArrs []map[string]string
	testArrs = append(testArrs, map[string]string{
		"key":   "name",
		"value": "CONCAT(first_name,' ',last_name)",
	})
	testArrs = append(testArrs, map[string]string{
		"key":   "phone",
		"value": "CONCAT(phone_code,phone)",
	})
	testArrs = append(testArrs, map[string]string{
		"key":   "email",
		"value": "email",
	})
	testArrs = append(testArrs, map[string]string{
		"key":   "is_active",
		"value": "is_active",
	})
	testArrs = append(testArrs, map[string]string{
		"key":   "inquiry_at",
		"value": "DATE(from_unixtime(inquiry_at))",
	})

	QryCnt.WriteString(" COUNT(id) AS cnt ")
	QryFilter.WriteString(" COUNT(id) AS cnt ")

	Qry.WriteString(" id,CONCAT(first_name,' ',last_name) AS name,CONCAT(phone_code,phone) AS phone,email,is_active,ip,from_unixtime(inquiry_at) AS inquiry_at ")

	FromQry.WriteString(" FROM cf_list_your_property ")
	FromQry.WriteString(" WHERE 1 = 1 ")
	Data, err := model.JQueryTable(r, reqMap, Qry, FromQry, QryCnt, QryFilter, GroupBy, testArrs, testColArrs, status, "")
	return Data, err
}
