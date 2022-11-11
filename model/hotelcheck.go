package model

import (
	"bytes"
	"tp-api-common/util"
)

// CheckHotelVerificationEligibility - Checks If Hotel's Whole Details Are Fulfilled Or Not
// If not doesn't allow to make verification happed - HK - 2021-05-19
func CheckHotelVerificationEligibility(hotelID string) (bool, string) {
	util.SysLogIt("CheckHotelVerificationEligibility Start")

	flgHotelBasic := CheckHotelBasicInfoExistance(hotelID)
	if !flgHotelBasic {
		util.SysLogIt("CheckHotelVerificationEligibility - Hotel basic info not fulfilled")
		return false, "Hotel basic info not fulfilled"
	}
	util.SysLogIt("CheckHotelVerificationEligibility - Hotel basic info checking passed")

	flgHotelAmenity := CheckHotelAmenityInfoExistance(hotelID)
	if !flgHotelAmenity {
		util.SysLogIt("CheckHotelVerificationEligibility - Hotel amenity info not fulfilled")
		return false, "Hotel amenity info not fulfilled"
	}

	flgHotelImage := CheckHotelImageInfoExistance(hotelID)
	if !flgHotelImage {
		util.SysLogIt("CheckHotelVerificationEligibility - Hotel image info not fulfilled")
		return false, "Hotel image info not fulfilled"
	}

	flgHotelBank := CheckHotelBankInfoExistance(hotelID)
	if !flgHotelBank {
		util.SysLogIt("CheckHotelVerificationEligibility - Hotel bank info not fulfilled")
		return false, "Hotel bank info not fulfilled"
	}

	flgHotelCanPolicy := CheckHotelCancelPolicyInfoExistance(hotelID, false)
	if !flgHotelCanPolicy {
		util.SysLogIt("CheckHotelVerificationEligibility - Hotel cancel policy info not fulfilled")
		return false, "Hotel cancel policy info not fulfilled"
	}

	flgHotelRoom := CheckHotelRoomInfoExistance(hotelID, false)
	if !flgHotelRoom {
		util.SysLogIt("CheckHotelVerificationEligibility - Hotel room info not fulfilled")
		return false, "Hotel room info not fulfilled"
	}

	util.SysLogIt("CheckHotelVerificationEligibility End")
	return true, ""
}

// CheckHotelBasicInfoExistance - Checks If Hotel's Basic Info Fulfilled Or Not - HK - 2021-05-17
func CheckHotelBasicInfoExistance(hotelID string) bool {
	util.SysLogIt("CheckHotelBasicInfoExistance Start")

	var HotelBasicQry, HotelTagQry bytes.Buffer
	HotelBasicQry.WriteString(" SELECT hotel_name, hotel_star, description, hotel_phone, property_type_id, ")
	HotelBasicQry.WriteString(" short_address, long_address, latitude, longitude, locality_id, city_id, state_id, country_id, ")
	HotelBasicQry.WriteString(" policy, checkin_rules")
	HotelBasicQry.WriteString(" FROM cf_hotel_info WHERE id = ?")
	HotelBasicInfo, err := ExecuteRowQuery(HotelBasicQry.String(), hotelID)
	if err != nil {
		util.SysLogIt("Error getting hotel basic info")
		return false
	}

	validateFlagStr := util.ValidateNotNullAndString(HotelBasicInfo, []string{"hotel_name", "description", "hotel_phone", "property_type_id", "short_address", "long_address", "locality_id", "policy", "checkin_rules", "latitude", "longitude"})
	validateFlagInt := util.ValidateNotNullAndInt(HotelBasicInfo, []string{"hotel_star", "city_id", "state_id", "country_id"})
	if validateFlagStr == 0 || validateFlagInt == 0 {
		util.SysLogIt("Error hotel basic info missing")
		return false
	}

	HotelTagQry.WriteString("SELECT COUNT(id) AS id FROM cf_hotel_tag WHERE hotel_id = ?")
	Data, err := ExecuteRowQuery(HotelTagQry.String(), hotelID)
	if err != nil {
		util.SysLogIt("Error getting hotel tag count")
		return false
	}

	if Data["id"].(int64) == 0 {
		util.SysLogIt("Error getting hotel tag count zero found")
		return false
	}

	util.SysLogIt("CheckHotelBasicInfoExistance End")
	return true
}

// CheckHotelAmenityInfoExistance - Checks If Hotel's Amenities Info Fulfilled Or Not - HK - 2021-05-18
func CheckHotelAmenityInfoExistance(hotelID string) bool {
	util.SysLogIt("CheckHotelAmenityInfoExistance Start")

	var HotelAmenityQry bytes.Buffer
	HotelAmenityQry.WriteString("SELECT COUNT(id) AS id FROM cf_hotel_amenities WHERE hotel_id = ?")
	Data, err := ExecuteRowQuery(HotelAmenityQry.String(), hotelID)
	if err != nil {
		util.SysLogIt("Error getting hotel amenity count")
		return false
	}

	if Data["id"].(int64) == 0 {
		util.SysLogIt("Error - hotel amenity count zero found")
		return false
	}

	util.SysLogIt("CheckHotelAmenityInfoExistance End")
	return true
}

// CheckHotelImageInfoExistance - Checks If Hotel's Images Info Fulfilled Or Not - HK - 2021-05-18
func CheckHotelImageInfoExistance(hotelID string) bool {
	util.SysLogIt("CheckHotelImageInfoExistance Start")

	var HotelImageQry bytes.Buffer
	HotelImageQry.WriteString("SELECT COUNT(id) AS id FROM cf_hotel_image WHERE hotel_id = ?")
	Data, err := ExecuteRowQuery(HotelImageQry.String(), hotelID)
	if err != nil {
		util.SysLogIt("CheckHotelImageInfoExistance - Error - getting hotel image count")
		return false
	}

	if Data["id"].(int64) == 0 {
		util.SysLogIt("CheckHotelImageInfoExistance - Error - hotel image count zero found")
		return false
	}

	util.SysLogIt("CheckHotelImageInfoExistance End")
	return true
}

// CheckHotelBankInfoExistance - Checks If Hotel's Bank Info Fulfilled Or Not - HK - 2021-05-18
func CheckHotelBankInfoExistance(hotelID string) bool {
	util.SysLogIt("CheckHotelBankInfoExistance Start")

	var HotelBankQry bytes.Buffer
	//HotelBankQry.WriteString(" SELECT account_holder_name, account_number, swift_code, bank_name, ")
	HotelBankQry.WriteString(" SELECT  ")
	HotelBankQry.WriteString(" checkin_time, checkout_time ")
	HotelBankQry.WriteString(" FROM cf_hotel_settings WHERE hotel_id = ?")
	HotelBankInfo, err := ExecuteRowQuery(HotelBankQry.String(), hotelID)
	if err != nil {
		util.SysLogIt("CheckHotelBankInfoExistance - Error - getting hotel bank info")
		return false
	}

	//validateFlagStr := util.ValidateNotNullAndString(HotelBankInfo, []string{"account_holder_name", "account_number", "swift_code", "bank_name", "checkin_time", "checkout_time"})
	validateFlagStr := util.ValidateNotNullAndString(HotelBankInfo, []string{"checkin_time", "checkout_time"})
	if validateFlagStr == 0 {
		util.SysLogIt("CheckHotelBankInfoExistance - Error - hotel bank info missing")
		return false
	}

	util.SysLogIt("CheckHotelBankInfoExistance End")
	return true
}

// CheckHotelCancelPolicyInfoExistance - Checks If Hotel's Cancel Policy Info Fulfilled Or Not - HK - 2021-05-19
func CheckHotelCancelPolicyInfoExistance(hotelID string, statusChk bool) bool {
	util.SysLogIt("CheckHotelCancelPolicyInfoExistance Start")

	var HotelCanPolicyQry bytes.Buffer
	HotelCanPolicyQry.WriteString("SELECT COUNT(id) AS id FROM cf_cancellation_policy WHERE hotel_id = ? ")
	if statusChk {
		HotelCanPolicyQry.WriteString(" AND status = 1 ")
	}
	Data, err := ExecuteRowQuery(HotelCanPolicyQry.String(), hotelID)
	if err != nil {
		util.SysLogIt("CheckHotelCancelPolicyInfoExistance - Error - getting hotel cancel policy count")
		return false
	}

	if Data["id"].(int64) == 0 {
		util.SysLogIt("CheckHotelCancelPolicyInfoExistance - Error - hotel cancel policy count zero found")
		return false
	}

	util.SysLogIt("CheckHotelCancelPolicyInfoExistance End")
	return true
}

// CheckHotelRoomInfoExistance - Checks If Hotel's Room Info Fulfilled Or Not - HK - 2021-05-19
func CheckHotelRoomInfoExistance(hotelID string, statusChk bool) bool {
	util.SysLogIt("CheckHotelRoomInfoExistance Start")

	var HotelRoomQry bytes.Buffer
	HotelRoomQry.WriteString("SELECT id, room_type_name, max_occupancy, inventory FROM cf_room_type WHERE hotel_id = ? ")
	if statusChk {
		HotelRoomQry.WriteString(" AND status = 1 ")
	}
	Data, err := ExecuteQuery(HotelRoomQry.String(), hotelID)
	if err != nil {
		util.SysLogIt("CheckHotelRoomInfoExistance - Error - getting hotel room count for " + hotelID)
		return false
	}

	if len(Data) == 0 {
		util.SysLogIt("CheckHotelRoomInfoExistance - Error - hotel room count zero found for " + hotelID)
		return false
	}

	if len(Data) > 0 {

		for _, j := range Data {

			roomID := j["id"].(string)
			roomName := j["room_type_name"].(string)

			// Amenity Count Checking Added
			AmenityDataCount, err := CheckHotelRoomAmenityInfoExistance(hotelID, roomID)
			if err != nil {
				util.SysLogIt("CheckHotelRoomInfoExistance - Error - getting room amenity count for " + hotelID + " - " + roomID + " - " + roomName)
				return false
			}
			if AmenityDataCount == 0 {
				util.SysLogIt("CheckHotelRoomInfoExistance - Error - room amenity count zero found for " + hotelID + " - " + roomID + " - " + roomName)
				return false
			}

			// Room Image Count Checking Added
			RoomImageCnt, err := CheckHotelRoomImageInfoExistance(hotelID, roomID)
			if err != nil {
				util.SysLogIt("CheckHotelRoomInfoExistance - Error - getting room image count for " + hotelID + " - " + roomID + " - " + roomName)
				return false
			}
			if RoomImageCnt == 0 {
				util.SysLogIt("CheckHotelRoomInfoExistance - Error - room image count zero found for " + hotelID + " - " + roomID + " - " + roomName)
				return false
			}

			// Room Rate Plan Count Checking Added
			RoomRatePlanCnt, err := CheckHotelRoomRatePlanInfoExistance(hotelID, roomID, false)
			if err != nil {
				util.SysLogIt("CheckHotelRoomInfoExistance - Error - getting rate plan count for " + hotelID + " - " + roomID + " - " + roomName)
				return false
			}
			if RoomRatePlanCnt == 0 {
				util.SysLogIt("CheckHotelRoomInfoExistance - Error - rate plan count zero found for " + hotelID + " - " + roomID + " - " + roomName)
				return false
			}
		}
	}

	util.SysLogIt("CheckHotelRoomInfoExistance End")
	return true
}

// CheckHotelRoomAmenityInfoExistance - Checks If Room Has Amenity Info Exists Or Not - 2021-05-19 - HK
func CheckHotelRoomAmenityInfoExistance(HotelID string, RoomID string) (int64, error) {
	util.SysLogIt("CheckHotelRoomAmenityhInfoExistance Start " + HotelID + " - " + RoomID)

	var Qry bytes.Buffer
	Qry.WriteString("SELECT COUNT(id) AS id FROM cf_room_amenity WHERE hotel_id = ? and room_type_id = ?")
	Data, err := ExecuteRowQuery(Qry.String(), HotelID, RoomID)
	if err != nil {
		util.SysLogIt("CheckHotelRoomAmenityInfoExistance - Error - getting room amenity count for " + HotelID + " - " + RoomID)
		return -1, err
	}

	util.SysLogIt("CheckHotelRoomAmenityInfoExistance End")
	return Data["id"].(int64), nil
}

// CheckHotelRoomImageInfoExistance - Checks If Room Has Image Info Exists Or Not - 2021-05-19 - HK
func CheckHotelRoomImageInfoExistance(HotelID string, RoomID string) (int64, error) {
	util.SysLogIt("CheckHotelRoomImageInfoExistance Start " + HotelID + " - " + RoomID)

	var Qry bytes.Buffer
	Qry.WriteString("SELECT COUNT(id) AS id FROM cf_room_image WHERE hotel_id = ? and room_type_id = ?")
	Data, err := ExecuteRowQuery(Qry.String(), HotelID, RoomID)
	if err != nil {
		util.SysLogIt("CheckHotelRoomImageInfoExistance - Error - getting room image count for " + HotelID + " - " + RoomID)
		return -1, err
	}

	util.SysLogIt("CheckHotelRoomImageInfoExistance End")
	return Data["id"].(int64), nil
}

// CheckHotelRoomRatePlanInfoExistance - Checks If Room Has Image Info Exists Or Not - 2021-05-19 - HK
func CheckHotelRoomRatePlanInfoExistance(HotelID string, RoomID string, statusChk bool) (int64, error) {
	util.SysLogIt("CheckHotelRoomRatePlanInfoExistance Start " + HotelID + " - " + RoomID)

	var Qry bytes.Buffer
	Qry.WriteString("SELECT COUNT(id) AS id FROM cf_rateplan WHERE hotel_id = ? and room_type_id = ?")
	if statusChk {
		Qry.WriteString(" AND status = 1 ")
	}
	Data, err := ExecuteRowQuery(Qry.String(), HotelID, RoomID)
	if err != nil {
		util.SysLogIt("CheckHotelRoomRatePlanInfoExistance - Error - getting rate plan count for " + HotelID + " - " + RoomID)
		return -1, err
	}

	util.SysLogIt("CheckHotelRoomRatePlanInfoExistance End")
	return Data["id"].(int64), nil
}
