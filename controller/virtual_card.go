package controller

import (
	"encoding/json"
	"net/http"
	"tp-api-common/data"
	"tp-api-common/util"
	"tp-system/config"
	"tp-system/model"

	"github.com/gorilla/mux"
	"github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/issuing/card"
	"github.com/stripe/stripe-go/issuing/cardholder"
	"github.com/stripe/stripe-go/topup"
)

// VirtualCardList - Virtual Card List
func VirtualCardList(w http.ResponseWriter, r *http.Request) {
	util.LogIt(r, "Controller - V_VirtualCard - HotelierListing")
	defer util.CommonDeferred(w, r, "Controller", "V_VirtualCard", "HotelierListing")
	var reqMap data.JQueryTableUI
	err := json.NewDecoder(r.Body).Decode(&reqMap)
	if util.CheckErrorLog(r, err) {
		util.RespondBadRequest(r, w)
		return
	}

	Data, err := model.VirtualCardList(r, reqMap)
	if util.CheckErrorLog(r, err) {
		util.RespondWithError(r, w, "500")
		return
	}

	util.Respond(r, w, Data, 200, "")
}

// VirtualCardDetail - Virtual Card Detail
func VirtualCardDetail(w http.ResponseWriter, r *http.Request) {
	util.LogIt(r, "Controller - V_VirtualCard - VirtualCardDetail")
	defer util.CommonDeferred(w, r, "Controller", "V_VirtualCard", "VirtualCardDetail")

	vars := mux.Vars(r)
	id := vars["id"]

	ValidateString := ValidateNotNullStructString(id)
	if ValidateString == 0 {
		util.RespondBadRequest(r, w)
		return
	}

	Data, err := model.VirtualCardDetail(r, id)
	if util.CheckErrorLog(r, err) {
		util.RespondWithError(r, w, "500")
		return
	}

	util.Respond(r, w, Data, 200, "")
}

// VirtualCardActive - Set Virtual Card Active
func VirtualCardActive(w http.ResponseWriter, r *http.Request) {
	util.LogIt(r, "Controller - Virtual_Card - VirtualCardActive")
	defer util.CommonDeferred(w, r, "Controller", "Virtual_Card", "VirtualCardActive")

	vars := mux.Vars(r)
	bookingID := vars["id"]

	ValidateString := ValidateNotNullStructString(bookingID)
	if ValidateString == 0 {
		util.RespondBadRequest(r, w)
		return
	}

	data, errData := model.VirtualCardActive(r, bookingID)
	if errData != nil {
		util.RespondWithError(r, w, "500")
		return
	}

	if data == nil {
		util.RespondData(r, w, nil, 422)
		return
	}

	// Card Active
	activeVirtualCard, erractiveVirtualCard := StripeVirtualCardActive(w, r, data["card_id"].(string))
	if erractiveVirtualCard != nil {
		util.LogIt(r, "erractiveVirtualCard - ")
		util.LogIt(r, erractiveVirtualCard)
		util.RespondData(r, w, nil, 422)
		return
	}
	util.LogIt(r, "activeVirtualCard ID - ")
	util.LogIt(r, activeVirtualCard.ID)
	util.LogIt(r, "activeVirtualCard Status - ")
	util.LogIt(r, activeVirtualCard.Status)

	// Card Top Up
	topUpFlag, errTopUPFlag := CardTopUp(r, "Virtual Card Top Up", data["amount"].(int64))
	if errTopUPFlag != nil {
		util.LogIt(r, "errTopUPFlag - ")
		util.LogIt(r, errTopUPFlag)
		util.RespondData(r, w, nil, 422)
		return
	}

	if !topUpFlag {
		util.LogIt(r, "topUpFlag - ")
		util.LogIt(r, topUpFlag)
		util.RespondData(r, w, nil, 422)
		return
	}

	updFlag := model.VirtualCardUpdate(r, bookingID)
	if !updFlag {
		util.RespondWithError(r, w, "500")
		return
	}

	util.LogIt(r, "Controller - Virtual_Card - CardTopUp - Success")
	util.RespondData(r, w, data, 200)
}

// CardTopUp - Virtual Card Top Up
func CardTopUp(r *http.Request, subject string, amount int64) (bool, error) {
	util.LogIt(r, "Controller - Virtual_Card - CardTopUp")
	stripe.Key = config.Env.Stripe.Secret

	params := &stripe.TopupParams{
		Amount:              stripe.Int64(amount),
		Currency:            stripe.String(string(stripe.CurrencyUSD)),
		Description:         stripe.String("Top-up for Issuing - " + subject),
		StatementDescriptor: stripe.String("Top-up"),
	}
	params.AddExtra("destination_balance", "issuing")
	_, errorTu := topup.New(params)
	if errorTu != nil {
		return false, errorTu
	}
	return true, nil
}

// StripeVirtualCardActive - Stripe Virtual Card Active
func StripeVirtualCardActive(w http.ResponseWriter, r *http.Request, cardID string) (*stripe.IssuingCard, error) {
	util.LogIt(r, "Controller - Virtual_Card - StripeVirtualCardActive")
	defer util.CommonDeferred(w, r, "Controller", "Virtual_Card", "StripeVirtualCardActive")
	stripe.Key = config.Env.Stripe.Secret

	params := &stripe.IssuingCardParams{Status: stripe.String("active")}
	activeVirtualCard, erractiveVirtualCard := card.Update(cardID, params)
	return activeVirtualCard, erractiveVirtualCard
}

// VirtualCardInfo -  Get Virtual Card and return the details of it
func VirtualCardInfo(w http.ResponseWriter, r *http.Request) {
	util.LogIt(r, "Controller - Virtual_Card - VirtualCardInfo")
	defer util.CommonDeferred(w, r, "Controller", "Virtual_Card", "VirtualCardInfo")

	vars := mux.Vars(r)
	bookingID := vars["id"]

	ValidateString := ValidateNotNullStructString(bookingID)
	if ValidateString == 0 {
		util.RespondBadRequest(r, w)
		return
	}

	data, errData := model.VirtualCardInfo(r, bookingID)
	if errData != nil {
		util.RespondWithError(r, w, "500")
		return
	}

	util.RespondData(r, w, data, 200)
}

// CreateCard -  Create Virtual Card
func CreateCard(w http.ResponseWriter, r *http.Request) {
	util.LogIt(r, "Controller - Virtual_Card - CreateCard")
	defer util.CommonDeferred(w, r, "Controller", "Virtual_Card", "CreateCard")
	newCardHolderData := make(map[string]interface{})

	vars := mux.Vars(r)
	bookingID := vars["id"]

	ValidateString := ValidateNotNullStructString(bookingID)
	if ValidateString == 0 {
		util.RespondBadRequest(r, w)
		return
	}

	hotelData, errHotelData := model.GetHotelDataForVirtualCard(r, bookingID)
	if !errHotelData {
		util.RespondWithError(r, w, "422")
		return
	}

	// Check whether cardholder is present for Hotel
	cardholderFlag := model.CardholderDataFlag(r, hotelData)
	util.LogIt(r, "cardholderFlag - ")
	util.LogIt(r, cardholderFlag)
	if !cardholderFlag {

		// Create CardHolder For Hotelier
		cardHolder, errCardholder := CreateCardHolderForCustomer(r, hotelData)
		if errCardholder != nil {
			util.LogIt(r, "errCardholder")
			util.LogIt(r, errCardholder)
			util.RespondWithError(r, w, "422")
			return
		}

		newCardHolderData["id"] = string(cardHolder.ID)
		newCardHolderData["hotel_id"] = hotelData["hotel_id"].(string)
		newCardHolderData["email"] = hotelData["hotel_email"].(string)
		newCardHolderData["name"] = hotelData["hotel_name"].(string)
		newCardHolderData["phoneno"] = hotelData["hotel_phone"].(string)
		newCardHolderData["status"] = string(cardHolder.Status)
		newCardHolderData["type"] = string(cardHolder.Type)

		// Inserting the Cardholder data into the sytem
		errSync := model.CardholderDataInSystem(r, newCardHolderData)
		if !errSync {
			util.LogIt(r, "errSync")
			util.LogIt(r, errSync)
			util.RespondWithError(r, w, "422")
			return
		}
	} else {
		// Get CardHolder Data
		newCardHolderData = model.CardholderDataForHotel(r, hotelData)
	}

	// Create Virtual Credit Card For Hotelier Payment
	VirtualCreditCard, errVirtualCard := CreateVirtualCard(r, newCardHolderData["id"].(string), hotelData["amount"].(float64))
	if errVirtualCard != nil {
		util.LogIt(r, "errVirtualCard")
		util.LogIt(r, errVirtualCard)
		util.RespondWithError(r, w, "422")
		return
	}

	// Inserting the Virtual Card data into the
	virtualCardData := make(map[string]interface{})
	virtualCardData["id"] = VirtualCreditCard.ID
	virtualCardData["cardholder_id"] = newCardHolderData["id"].(string)
	virtualCardData["booking_id"] = bookingID
	virtualCardData["currency"] = string(VirtualCreditCard.Currency)
	virtualCardData["name"] = hotelData["hotel_name"].(string)
	virtualCardData["expmonth"] = VirtualCreditCard.ExpMonth
	virtualCardData["expyear"] = VirtualCreditCard.ExpYear
	virtualCardData["spending_limit"] = hotelData["amount"]
	virtualCardData["status"] = string(VirtualCreditCard.Status)
	virtualCardData["type"] = string(VirtualCreditCard.Type)

	var cardNumber *stripe.IssuingCard
	cardNumber, cardNoError := GetVirtualCard(r, VirtualCreditCard.ID)
	if cardNoError != nil {
		util.LogIt(r, "cardNoError")
		util.LogIt(r, cardNoError)
		util.RespondWithError(r, w, "422")
		return
	}

	virtualCardData["cardno"] = string(cardNumber.Number)

	errCardSync := model.VirtualCardDataInSystem(r, virtualCardData)
	if !errCardSync {
		util.LogIt(r, "errCardSync")
		util.LogIt(r, errCardSync)
		util.RespondWithError(r, w, "422")
		return
	}

	util.RespondData(r, w, nil, 200)
}

// CreateCardHolderForCustomer - Card Holder is created for Hotelier to create a Virtual Credit Card
func CreateCardHolderForCustomer(r *http.Request, hotelData map[string]interface{}) (*stripe.IssuingCardholder, error) {
	util.LogIt(r, "Controller - Virtual_Card - createCardHolderForCustomer")
	stripe.Key = config.Env.Stripe.Secret

	params := &stripe.IssuingCardholderParams{
		Billing: &stripe.IssuingCardholderBillingParams{
			Address: &stripe.AddressParams{
				Line1:      stripe.String(hotelData["long_address"].(string)),
				City:       stripe.String(hotelData["city_name"].(string)),
				State:      stripe.String(hotelData["state_name"].(string)),
				Country:    stripe.String(hotelData["country_name"].(string)),
				PostalCode: stripe.String("94111"),
			},
		},
		Email:       stripe.String(hotelData["hotel_email"].(string)),
		Name:        stripe.String(hotelData["hotel_name"].(string)),
		PhoneNumber: stripe.String(hotelData["hotel_phone"].(string)),
		Type:        stripe.String("individual"),
	}
	newCardholder, errCardholder := cardholder.New(params)
	return newCardholder, errCardholder
}

// CreateVirtualCard - Virtual Card is created for Hotelier payment
func CreateVirtualCard(r *http.Request, cardholderID string, amount float64) (*stripe.IssuingCard, error) {
	util.LogIt(r, "Controller - Virtual_Card - CreateVirtualCard")
	stripe.Key = config.Env.Stripe.Secret

	params := &stripe.IssuingCardParams{
		Cardholder: stripe.String(cardholderID),
		Currency:   stripe.String(string(stripe.CurrencyUSD)),
		Type:       stripe.String("virtual"),
		SpendingControls: &stripe.IssuingCardSpendingControlsParams{
			SpendingLimits: []*stripe.IssuingCardSpendingControlsSpendingLimitParams{
				{
					Amount:   stripe.Int64(int64(amount)),
					Interval: stripe.String(string(stripe.IssuingCardSpendingControlsSpendingLimitIntervalAllTime)),
				},
			},
		},
	}
	// params.AddExpand("number")

	newVirtualCard, errVirtualCard := card.New(params)
	return newVirtualCard, errVirtualCard
}

// GetVirtualCard - Virtual Card is created for Hotelier payment
func GetVirtualCard(r *http.Request, virtualCardID string) (*stripe.IssuingCard, error) {
	util.LogIt(r, "Controller - Virtual_Card - GetVirtualCard")
	stripe.Key = config.Env.Stripe.Secret

	params := &stripe.IssuingCardParams{}
	params.AddExpand("number")
	newVirtualCard, errVirtualCard := card.Get(
		virtualCardID,
		params,
	)
	return newVirtualCard, errVirtualCard
}
