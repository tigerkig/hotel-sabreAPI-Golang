package model

import (
	"bytes"
	"mallbackend/util"
)

//GetBookingInfoForSendMail - Get booking info all details from mail
func GetBookingInfoForSendMail(bookingID string) (map[string]interface{}, bool) {
	var bookingStuff = make(map[string]interface{})
	var CancelledPolicyQry, BookingRoomQry, BookingInfo, GuestInfo, RoomCharge, RoomTax, Payment, Discount bytes.Buffer

	BookingInfo.WriteString(" SELECT hotel_phone, latitude, longitude, booking_no, DATE_FORMAT(booking_time, '%d %b %Y %h:%i %p') AS booking_date_time,DATE_FORMAT(checkin_time, '%d %b %y %h:%i %p') AS checkin_date_time,DATE_FORMAT(checkin_time, '%a, %d %b %y') AS check_in_day_date,DATE_FORMAT(checkin_time, ' %h:%i %p') AS check_in_time,DATE_FORMAT(checkout_time, '%a, %d %b %y') AS check_out_day_date,DATE_FORMAT(checkout_time, ' %h:%i %p') AS check_out_time, checkin_time,checkout_time,no_of_night,no_of_room,adult,child,booking_status,hotel_name,hotel_id,property_type,hotel_short_address,hotel_long_address,hotel_locality,hotel_state,hotel_city,hotel_image,status FROM tp_front.fd_booking WHERE id = ? ")
	BookingInfoData, err := ExecuteRowQuery(BookingInfo.String(), bookingID)
	if err != nil {
		util.SysLogIt(err)
		return nil, false
	}

	bookingStuff["booking_info"] = BookingInfoData

	GuestInfo.WriteString(" SELECT CONCAT(first_name,' ',last_name) AS guest_name, email, CONCAT(phone_code,mobile) AS mobile FROM tp_front.fd_booking_guest WHERE is_primary = 1 AND booking_id = ? ")
	GuestInfoData, err := ExecuteRowQuery(GuestInfo.String(), bookingID)
	if err != nil {
		util.SysLogIt(err)
		return nil, false
	}

	bookingStuff["guest_info"] = GuestInfoData

	CancelledPolicyQry.WriteString(" SELECT is_norefund,DATE_FORMAT(cancellation_term,'%d %b %Y') AS cancellation_term,before_term_charge,after_term_charge FROM tp_front.fd_booking_cancellation_policy WHERE booking_id = ?; ")
	CancelledPolicy, err := ExecuteRowQuery(CancelledPolicyQry.String(), bookingID)
	if err != nil {
		util.SysLogIt(err)
		return nil, false
	}

	bookingStuff["cancellation_policy"] = CancelledPolicy

	BookingRoomQry.WriteString(" SELECT room_type_id,room_type_name,rateplan_id,rateplan_name,adult,child,room_image FROM tp_front.fd_booking_room WHERE booking_id = ?; ")
	BookingRoom, err := ExecuteQuery(BookingRoomQry.String(), bookingID)
	if err != nil {
		util.SysLogIt(err)
		return nil, false
	}

	bookingStuff["booking_room_info"] = BookingRoom

	accountDetails := make(map[string]interface{})
	RoomCharge.WriteString(" SELECT master_name,ROUND(amount) AS amount,description FROM tp_front.fd_account_detail WHERE booking_id = ? AND master_id = 1;")
	RoomChargeData, err := ExecuteRowQuery(RoomCharge.String(), bookingID)
	if err != nil {
		util.SysLogIt(err)
		return nil, false
	}

	accountDetails["room_charge"] = RoomChargeData

	RoomTax.WriteString(" SELECT master_name,ROUND(amount) AS amount,description FROM tp_front.fd_account_detail WHERE booking_id = ? AND master_id = 2;")
	RoomTaxData, err := ExecuteRowQuery(RoomTax.String(), bookingID)
	if err != nil {
		util.SysLogIt(err)
		return nil, false
	}

	accountDetails["room_tax"] = RoomTaxData

	Payment.WriteString(" SELECT master_name,ROUND(amount,0)*(-1) AS amount,description FROM tp_front.fd_account_detail WHERE booking_id = ? AND master_id = 4;")
	PaymentData, err := ExecuteRowQuery(Payment.String(), bookingID)
	if err != nil {
		util.SysLogIt(err)
		return nil, false
	}

	accountDetails["room_payment"] = PaymentData

	Discount.WriteString(" SELECT master_name,ROUND(amount) AS amount,description FROM tp_front.fd_account_detail WHERE booking_id = ? AND master_id = 3;")
	DiscountData, err := ExecuteRowQuery(Discount.String(), bookingID)
	if err != nil {
		util.SysLogIt(err)
		return nil, false
	}

	if len(DiscountData) > 0 {
		accountDetails["room_discount"] = DiscountData
	} else {
		accountDetails["room_discount"] = make(map[string]interface{})
	}

	bookingStuff["room_info"] = accountDetails

	return bookingStuff, true
}
