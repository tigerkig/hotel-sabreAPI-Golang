package model

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"html/template"
	"strconv"
	"strings"
	"tp-api-common/util"
	"tp-system/config"

	gonanoid "github.com/matoous/go-nanoid"
	"gopkg.in/gomail.v2"
)

//MailObj - Mail object
type MailObj struct {
	Type         string
	ID           string
	Additional   string
	InterfaceObj interface{}
}

//MailChn - Mail channel
var MailChn = make(chan MailObj)

// HandleSendMailObject - Handle send mail
func HandleSendMailObject() {
	for {
		select {
		case message := <-MailChn:
			mtype := message.Type
			switch mtype {
			case "EmailTemplate":
				SendMailUsingTemplate(message.ID, message.Additional, message.InterfaceObj)
				break
			case "UserEmailTemplate":
				SendUserMailUsingTemplate(message.ID)
				break
			default:
				break
			}
		}
	}
}

type SMTP struct {
	Host      string `json:"smtp_host"`
	Port      int    `json:"smtp_port"`
	Username  string `json:"smtp_user"`
	Password  string `json:"smtp_password"`
	Email     string `json:"email_id"`
	EmailName string `json:"email_name"`
}

// GetEmailTemplateInfo -  Return Email Template Details
func GetEmailTemplateInfo(id string) (map[string]interface{}, error) {
	util.SysLogIt("model - Email_Template - GetEmailTemplateInfo")
	var Qry bytes.Buffer

	Qry.WriteString(" SELECT CET.subject,CET.id,CET.email_template_name,CEC.email_id,CET.short_code,CEC.email_name,CEC.smtp_host,CEC.smtp_port,CEC.smtp_user,from_base64(CEC.smtp_password) AS smtp_password,CEC.signature, IFNULL(GROUP_CONCAT(CEC1.email_id),'') AS bcc, CET.template ")
	Qry.WriteString(" FROM cf_email_template AS CET ")
	Qry.WriteString(" INNER JOIN cf_email_config AS CEC ON CEC.id = CET.email_from ")
	Qry.WriteString(" LEFT JOIN cf_email_config AS CEC1 ON find_in_set(CEC1.id, CET.email_bcc) ")
	Qry.WriteString(" WHERE CET.id = ? GROUP BY CET.id")

	RetMap, err := ExecuteRowQuery(Qry.String(), id)
	if err != nil {
		return nil, err
	}

	return RetMap, nil
}

//SendPartnerActivationMail - send partner activation mail
func SendMailUsingTemplate(templateType, to string, funcMap interface{}) bool {
	Data, _ := GetEmailTemplateInfo(templateType)
	m := SMTP{
		Host:      Data["smtp_host"].(string),
		Port:      int(Data["smtp_port"].(int64)),
		Username:  Data["smtp_user"].(string),
		Password:  Data["smtp_password"].(string),
		Email:     Data["email_id"].(string),
		EmailName: Data["email_name"].(string),
	}

	subject := Data["subject"].(string)
	defaultTemplate := Data["template"].(string)
	var ReqMap = funcMap.(map[string]interface{})
	if mailErr := m.Mail(to, subject, defaultTemplate, ReqMap); mailErr != nil {
		util.SysLogIt(mailErr)
	}

	return true
}

// Mail sends a templated mail. It will try to load the template from a URL, and
// otherwise fall back to the default
func (m SMTP) Mail(to, subjectTemplate, defaultTemplate string, templateData map[string]interface{}) error {
	tmp, err := template.New("Subject").Funcs(template.FuncMap(make(map[string]interface{}))).Parse(subjectTemplate)
	if err != nil {
		return err
	}

	subject := &bytes.Buffer{}
	err = tmp.Execute(subject, templateData)
	if err != nil {
		return err
	}

	body, err := MailBody(defaultTemplate, templateData)
	if err != nil {
		return err
	}

	mail := gomail.NewMessage()
	mail.SetHeader("From", m.Email, m.EmailName)
	mail.SetHeader("To", to)
	mail.SetHeader("Subject", subject.String())
	mail.SetBody("text/html", body)

	dial := gomail.NewPlainDialer(m.Host, m.Port, m.Email, m.Password)
	//dial := gomail.NewPlainDialer("smtp.gmail.com", 587, "mallbackendall@gmail.com", "Abcd1298.")
	return dial.DialAndSend(mail)

}

//MailBody - Mail Body
func MailBody(defaultTemplate string, data map[string]interface{}) (string, error) {

	var temp *template.Template
	var err error
	nanoID, _ := gonanoid.Nanoid()
	parsed, err := template.New(nanoID).Funcs(map[string]interface{}{}).Parse(defaultTemplate)
	if err != nil {
		return "", err
	}

	temp = parsed

	buf := &bytes.Buffer{}
	err = temp.Execute(buf, data)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

// SendMail -- email send
func SendMail(smtp map[string]interface{}, bccMail string, toMail string, toName string, subjectLine string, mailBody string) bool {
	var err error
	//bccEmail := []string{}
	smtpPort := int(smtp["port"].(int64))
	d := gomail.NewDialer(smtp["host"].(string), smtpPort, smtp["user"].(string), smtp["password"].(string))
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	m := gomail.NewMessage()
	m.SetAddressHeader("From", smtp["email_id"].(string), smtp["email_name"].(string))
	m.SetAddressHeader("To", toMail, toName) // To Email
	if bccMail != "" {
		str := strings.Split(bccMail, ",")
		addresses := make([]string, len(str))
		if len(str) > 0 {
			for i, val := range str {
				addresses[i] = m.FormatAddress(val, "")
			}
		}
		m.SetHeader("Bcc", addresses...)
		//bccEmail = addresses
	}
	m.SetHeader("Subject", subjectLine) //Subject line
	m.SetBody("text/html", mailBody)    // Email Body

	// Send the email
	if err = d.DialAndSend(m); err != nil {
		util.SysLogIt("Error while sending email")
		util.SysLogIt(err)
		return false
	}

	return true
}

//SendUserMailUsingTemplate - send mail from template of booking
func SendUserMailUsingTemplate(bookingID string) bool {
	util.SysLogIt("model - SendUserMailUsingTemplate")

	BookingData, _ := GetBookingInfoForSendMail(bookingID)
	if len(BookingData) == 0 {
		return false
	}

	var bookingData = BookingData["booking_info"].(map[string]interface{})

	templateID := "1"
	if bookingData["status"].(int64) == 2 {
		templateID = "16"
	}

	GetBookingTemData, err := GetUserEmailTemplate(templateID)
	if err != nil {
		util.SysLogIt(err)
		return false
	}

	//Company Info
	GetCompanyInfo, errCompanyInfo := CompanyInfo("1")
	if errCompanyInfo != nil {
		util.SysLogIt(errCompanyInfo)
		return false
	}

	smtp := make(map[string]interface{})
	smtp["host"] = GetBookingTemData["smtp_host"]
	smtp["user"] = GetBookingTemData["smtp_user"]
	smtp["password"] = GetBookingTemData["smtp_password"]
	smtp["port"] = GetBookingTemData["smtp_port"]
	smtp["email_name"] = GetBookingTemData["email_name"]
	smtp["email_id"] = GetBookingTemData["email_id"]

	var guestData = BookingData["guest_info"].(map[string]interface{})

	BookingData["company_logo"] = GetCompanyInfo["image"].(string)
	BookingData["company_name"] = GetCompanyInfo["company_name"].(string)

	HTML := ReplaceLabelOfBookingHTML(GetBookingTemData["template"].(string), BookingData)
	SendMail(smtp, GetBookingTemData["bcc"].(string), guestData["email"].(string), guestData["guest_name"].(string), GetBookingTemData["subject"].(string), HTML)
	// if mailErr := m.Mail(guestData["email"].(string), GetBookingTemData["subject"].(string), HTML, nil); mailErr != nil {
	// 	util.SysLogIt(mailErr)
	// }
	return true
}

//ReplaceLabelOfBookingHTML - Replace static lable of html and make it dynamic
func ReplaceLabelOfBookingHTML(html string, bookingInfo map[string]interface{}) string {
	html = strings.ReplaceAll(html, "{currency}", "USD")
	var bookingData = bookingInfo["booking_info"].(map[string]interface{})
	if bookingData != nil || len(bookingData) != 0 {
		html = strings.ReplaceAll(html, "{hotel_name}", bookingData["hotel_name"].(string))
		html = strings.ReplaceAll(html, "{booking_id}", bookingData["booking_no"].(string))
		html = strings.ReplaceAll(html, "{booking_date_time}", bookingData["booking_date_time"].(string))
		html = strings.ReplaceAll(html, "{checkin_date_time}", bookingData["checkin_date_time"].(string))
		html = strings.ReplaceAll(html, "{check_in_day_date}", bookingData["check_in_day_date"].(string))
		html = strings.ReplaceAll(html, "{check_in_time}", bookingData["check_in_time"].(string))
		html = strings.ReplaceAll(html, "{check_out_day_date}", bookingData["check_out_day_date"].(string))
		html = strings.ReplaceAll(html, "{check_out_time}", bookingData["check_out_time"].(string))
		html = strings.ReplaceAll(html, "{no_of_nights}", strconv.Itoa(int(bookingData["no_of_night"].(int64))))
		html = strings.ReplaceAll(html, "{no_of_room}", strconv.Itoa(int(bookingData["no_of_room"].(int64))))
		html = strings.ReplaceAll(html, "{no_of_guest}", strconv.Itoa(int(bookingData["adult"].(int64)+bookingData["child"].(int64))))
		//html = strings.ReplaceAll(html, "{hotel_long_address}", bookingData["hotel_long_address"].(string))
		html = strings.ReplaceAll(html, "{hotel_city}", bookingData["hotel_city"].(string))
		html = strings.ReplaceAll(html, "{hotel_state}", bookingData["hotel_state"].(string))
		//html = strings.ReplaceAll(html, "{hotel_phone}", bookingData["hotel_phone"].(string))
		html = strings.ReplaceAll(html, "{latitude}", fmt.Sprintf("%f", bookingData["latitude"].(float64)))
		html = strings.ReplaceAll(html, "{longitude}", fmt.Sprintf("%f", bookingData["longitude"].(float64)))
		//html = strings.ReplaceAll(html, "{hotel_image}", config.Env.AwsBucketURL+config.Env.HotelFolder+"/"+bookingData["hotel_image"].(string))
		html = strings.ReplaceAll(html, "{hotel_image}", bookingData["hotel_image"].(string))

		// Status = 2 For booking rejection
		if bookingData["status"].(int64) == 2 {
			html = strings.ReplaceAll(html, "is Confirmed!", "is Cancelled by Hotel!")
		}
	}

	var guestData = bookingInfo["guest_info"].(map[string]interface{})
	if guestData != nil || len(guestData) != 0 {
		html = strings.ReplaceAll(html, "{guest_name}", guestData["guest_name"].(string))
		html = strings.ReplaceAll(html, "{guest_phone}", guestData["mobile"].(string))
		html = strings.ReplaceAll(html, "{guest_email}", guestData["email"].(string))
	}

	var cancelData = bookingInfo["cancellation_policy"].(map[string]interface{})
	if cancelData != nil || len(cancelData) != 0 {
		html = strings.ReplaceAll(html, "{cancellation_date_time}", cancelData["cancellation_term"].(string))
		if cancelData["is_norefund"].(int64) == 1 {
			html = strings.ReplaceAll(html, "{is_no_refund}", "block")
			html = strings.ReplaceAll(html, "{is_refund}", "none")
			html = strings.ReplaceAll(html, "{cancellation_term}", cancelData["cancellation_term"].(string))
			html = strings.ReplaceAll(html, "{before_term_charge}", cancelData["before_term_charge"].(string))
			html = strings.ReplaceAll(html, "{after_term_charge}", cancelData["after_term_charge"].(string))
		} else {
			html = strings.ReplaceAll(html, "{is_no_refund}", "none")
			html = strings.ReplaceAll(html, "{is_refund}", "block")
		}
	}

	var roomData = bookingInfo["booking_room_info"].([]map[string]interface{})
	if roomData != nil || len(roomData) != 0 {
		html = strings.ReplaceAll(html, "{room_name}", roomData[0]["room_type_name"].(string))
		html = strings.ReplaceAll(html, "{no_of_adults}", strconv.Itoa(int(bookingData["adult"].(int64))))
		html = strings.ReplaceAll(html, "{rateplan_name}", roomData[0]["rateplan_name"].(string))
		//html = strings.ReplaceAll(html, "{room_image}", config.Env.AwsBucketURL+config.Env.RoomFolder+"/"+roomData[0]["room_image"].(string))
		html = strings.ReplaceAll(html, "{room_image}", roomData[0]["room_image"].(string))
	}

	var accDetails = bookingInfo["room_info"].(map[string]interface{})
	if accDetails != nil || len(accDetails) != 0 {
		var totalAmount = accDetails["room_payment"].(map[string]interface{})
		if totalAmount != nil || len(totalAmount) > 0 {
			html = strings.ReplaceAll(html, "{total_amount}", totalAmount["amount"].(string))
			var totalTax = accDetails["room_tax"].(map[string]interface{})
			html = strings.ReplaceAll(html, "{total_tax}", totalTax["amount"].(string))
			var totalRoomCharge = accDetails["room_charge"].(map[string]interface{})
			html = strings.ReplaceAll(html, "{total_room_charges}", totalRoomCharge["amount"].(string))
			html = strings.ReplaceAll(html, "{company_logo}", bookingInfo["company_logo"].(string))
			html = strings.ReplaceAll(html, "{company_name}", bookingInfo["company_name"].(string))
			html = strings.ReplaceAll(html, "{website_link}", config.Env.FrontURL)
		}
	}

	return html
}

// GetUserEmailTemplate -  Return Email Template Details
func GetUserEmailTemplate(id string) (map[string]interface{}, error) {
	util.SysLogIt("model - mail - GetUserEmailTemplate")
	var Qry bytes.Buffer

	Qry.WriteString(" SELECT CET.subject,CET.id,CET.email_template_name,CEC.email_id,CET.short_code,CEC.email_name,CEC.smtp_host,CEC.smtp_port,CEC.smtp_user,from_base64(CEC.smtp_password) AS smtp_password,CEC.signature, IFNULL(GROUP_CONCAT(CEC1.email_id),'') AS bcc, CET.template ")
	Qry.WriteString(" FROM cf_email_template AS CET ")
	Qry.WriteString(" INNER JOIN cf_email_config AS CEC ON CEC.id = CET.email_from ")
	Qry.WriteString(" LEFT JOIN cf_email_config AS CEC1 ON find_in_set(CEC1.id, CET.email_bcc) ")
	Qry.WriteString(" WHERE CET.id = ? GROUP BY CET.id")

	RetMap, err := ExecuteRowQuery(Qry.String(), id)
	if err != nil {
		util.SysLogIt(err)
		return nil, err
	}

	return RetMap, nil
}

// CompanyInfo -  Get Company Info
func CompanyInfo(id string) (map[string]interface{}, error) {
	util.SysLogIt("model - mail - CompanyInfo")

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
	Qry.WriteString(" WHERE CCI.id = ?")
	RetMap, err = ExecuteRowQuery(Qry.String(), 1)
	if err != nil {
		util.SysLogIt(err)
		return nil, err
	}

	return RetMap, nil
}
