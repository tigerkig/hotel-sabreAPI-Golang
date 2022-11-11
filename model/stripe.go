package model

import (
	"bytes"
	"net/http"
	"tp-api-common/util"
	"tp-system/config"

	"github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/account"
	"github.com/stripe/stripe-go/accountlink"
	"github.com/stripe/stripe-go/loginlink"
)

//CreateStripeAccount - Create stripe connected account
func CreateStripeAccount(name, email string) (string, error) {
	params := &stripe.AccountParams{
		Capabilities: &stripe.AccountCapabilitiesParams{
			CardPayments: &stripe.AccountCapabilitiesCardPaymentsParams{
				Requested: stripe.Bool(true),
			},
			Transfers: &stripe.AccountCapabilitiesTransfersParams{
				Requested: stripe.Bool(true),
			},
		},
		Country: stripe.String("US"),
		Email:   stripe.String(email),
		Type:    stripe.String("express"),
		BusinessProfile: &stripe.AccountBusinessProfileParams{
			Name: stripe.String(name),
		},
	}
	a, err := account.New(params)
	if err != nil {
		util.SysLogIt("Error found in create stripe account")
		util.SysLogIt("Error while create stripe account for email " + email)
	}
	return a.ID, err
}

func AccountLinks(accID string) (string, error) {
	params := &stripe.AccountLinkParams{
		Account:    stripe.String(accID),
		RefreshURL: stripe.String(config.Env.Stripe.APIURL + "reauth/" + accID),
		ReturnURL:  stripe.String(config.Env.Stripe.APIURL + "return/" + accID),
		Type:       stripe.String("account_onboarding"),
	}
	al, err := accountlink.New(params)
	if err != nil {
		util.SysLogIt("Error found in create stripe account links")
		util.SysLogIt("Error while create stripe account links for accID " + accID)
	}
	return al.URL, err
}

// AddAccountIDToHotelInfo - Mapped stripe connected account with hotel
func AddAccountIDToHotelInfo(r *http.Request, hotelID string) bool {
	var Qry, GQry bytes.Buffer
	GQry.WriteString(" SELECT email,hotel_name  ")
	GQry.WriteString(" FROM cf_hotel_client AS HC  ")
	GQry.WriteString(" INNER JOIN cf_hotel_info AS HI ON HI.group_id = HC.group_id AND HC.group_id <> '' ")
	GQry.WriteString(" WHERE HI.id = ?  ")
	HotelInfo, err := ExecuteRowQuery(GQry.String(), hotelID)
	if util.CheckErrorLog(r, err) {
		return false
	}

	if HotelInfo == nil {
		return false
	}

	email := HotelInfo["email"].(string)
	hname := HotelInfo["hotel_name"].(string)

	accID, err := CreateStripeAccount(hname, email)
	if util.CheckErrorLog(r, err) {
		return false
	}

	Qry.WriteString("UPDATE cf_hotel_info SET stripe_account = ? WHERE id = ?")
	err = ExecuteNonQuery(Qry.String(), accID, hotelID)
	if util.CheckErrorLog(r, err) {
		return false
	}

	onBoardingURL, err := AccountLinks(accID)
	if util.CheckErrorLog(r, err) {
		return false
	}

	resMap := make(map[string]interface{})
	resMap["link"] = onBoardingURL
	resMap["hotel_name"] = hname
	//Send stripe account connects link to hotelier
	MailChn <- MailObj{
		Type:         "EmailTemplate",
		ID:           "10",
		Additional:   email,
		InterfaceObj: resMap,
	}

	return true
}

func HotelStripeAccountCreated(r *http.Request, hotelID string) bool {
	var Qry bytes.Buffer
	Qry.WriteString(" SELECT *  ")
	Qry.WriteString(" FROM cf_hotel_info ")
	Qry.WriteString(" WHERE id = ? AND stripe_account != '' ")
	HotelInfo, err := ExecuteQuery(Qry.String(), hotelID)
	if util.CheckErrorLog(r, err) {
		return false
	}

	if len(HotelInfo) > 0 {
		return true
	}

	return false
}

//CheckConnectStripeAccExists - This function is used to check if passing stripe connected
// account exists in DB or not
func CheckConnectStripeAccExists(r *http.Request, accID string, chkFlag bool) bool {
	var Qry bytes.Buffer
	Qry.WriteString("SELECT * FROM cf_hotel_info WHERE stripe_account = ?")
	if chkFlag {
		Qry.WriteString(" AND is_stripe_bank_created = 1 ")
	}
	Data, err := ExecuteQuery(Qry.String(), accID)
	if util.CheckErrorLog(r, err) {
		return false
	}

	if len(Data) == 0 {
		return false
	}

	return true
}

//CreateLoginLink - Create Login Link of existing account
func CreateLoginLink(accID string) (string, error) {
	params := &stripe.LoginLinkParams{
		Account: stripe.String(accID),
	}
	ll, err := loginlink.New(params)
	return ll.URL, err
}

// UpdateHotelBankStatus - This function use to update bank status flag if
// Bank status is created or not in stripe end
func UpdateHotelBankStatus(r *http.Request, accID string) bool {
	var Qry, EmailQry bytes.Buffer
	Qry.WriteString("UPDATE cf_hotel_info SET is_stripe_bank_created = 1 WHERE stripe_account = ?")
	err := ExecuteNonQuery(Qry.String(), accID)
	if util.CheckErrorLog(r, err) {
		return false
	}

	EmailQry.WriteString(" SELECT email,hotel_name,HI.id AS hotel_id  ")
	EmailQry.WriteString(" FROM cf_hotel_client AS HC  ")
	EmailQry.WriteString(" INNER JOIN cf_hotel_info AS HI ON HI.group_id = HC.group_id AND HC.group_id <> '' ")
	EmailQry.WriteString(" WHERE stripe_account = ? ")
	EmailData, err := ExecuteRowQuery(EmailQry.String(), accID)
	if util.CheckErrorLog(r, err) {
		return false
	}

	if _, ok := EmailData["email"].(string); ok {
		connectLoginURL, err := CreateLoginLink(accID)
		if util.CheckErrorLog(r, err) {
			return false
		}

		resMap := make(map[string]interface{})
		resMap["link"] = connectLoginURL
		resMap["hotel_name"] = EmailData["hotel_name"].(string)

		//Send stripe account connects link to hotelier
		MailChn <- MailObj{
			Type:         "EmailTemplate",
			ID:           "13",
			Additional:   EmailData["email"].(string),
			InterfaceObj: resMap,
		}
		go UpdateHotelOnList(EmailData["hotel_id"].(string))
	}

	return true
}
