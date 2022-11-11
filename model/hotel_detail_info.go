package model

import (
	"bytes"
	"net/http"
	"tp-api-common/util"
	"tp-system/config"
)

// GetHotelDetailInfo - Get Hotel Detail Info
func GetHotelDetailInfo(r *http.Request, id string) (map[string]interface{}, error) {
	util.LogIt(r, "model - Model - GetHotelDetailInfo")
	resMap := make(map[string]interface{})
	var Select, Qry, SettingQry, ImageCategoryQry bytes.Buffer

	Select.WriteString(" SELECT ")
	Qry.WriteString(" CHI.id, hotel_name, short_address, long_address, hotel_phone, hotel_star, description, ")
	Qry.WriteString(" latitude, longitude, policy, checkin_rules,  ")
	Qry.WriteString(" CFC.name AS city, CFS.name AS state, CFY.name AS country,  ")
	Qry.WriteString(" CFL.locality,CFU.username AS account_manager, CPT.type AS property_type, group_concat(CFT.tag) AS tags  ")
	Qry.WriteString(" FROM cf_hotel_info AS CHI ")
	Qry.WriteString(" INNER JOIN cf_property_type AS CPT ON CPT.id = CHI.property_type_id ")
	Qry.WriteString(" LEFT JOIN cf_city AS CFC ON CFC.id = CHI.city_id ")
	Qry.WriteString(" LEFT JOIN cf_states AS CFS ON CFS.id = CHI.state_id ")
	Qry.WriteString(" LEFT JOIN cf_country AS CFY ON CFY.id = CHI.country_id ")
	Qry.WriteString(" LEFT JOIN cf_locality AS CFL ON CFL.id = CHI.locality_id ")
	Qry.WriteString(" LEFT JOIN cf_user AS CFU ON CFU.id = CHI.account_manager ")
	Qry.WriteString(" LEFT JOIN cf_hotel_tag AS CHT ON CHT.hotel_id = CHI.id ")
	Qry.WriteString(" LEFT JOIN cf_tags AS CFT ON CFT.id = CHT.tag_id ")
	Qry.WriteString(" WHERE CHI.status <> 3 AND CHI.id = ? ")
	Qry.WriteString(" GROUP BY CHI.id ")
	HotelInfoMap, err := ExecuteRowQuery(Select.String()+Qry.String(), id)
	if util.CheckErrorLog(r, err) {
		return nil, err
	}

	SettingQry.WriteString(" checkin_time,checkout_time,commission_amount,account_holder_name,account_number,swift_code,bank_name ")
	SettingQry.WriteString(" FROM cf_hotel_settings ")
	SettingQry.WriteString(" WHERE hotel_id = ? ")
	SettingMap, err := ExecuteRowQuery(Select.String()+SettingQry.String(), id)
	if util.CheckErrorLog(r, err) {
		return nil, err
	}

	ImageCategoryQry.WriteString("SELECT id,name AS category FROM cf_image_category WHERE status = 1")
	CategoryListing, err := ExecuteQuery(ImageCategoryQry.String())
	if util.CheckErrorLog(r, err) {
		return nil, err
	}

	if len(CategoryListing) > 0 {
		for _, v := range CategoryListing {
			var ImageQry bytes.Buffer
			ImageQry.WriteString("SELECT CONCAT('" + config.Env.AwsBucketURL + "hotel/" + "',image) AS image FROM cf_hotel_image WHERE category_id = ? AND hotel_id = ?")
			ImgMap, err := ExecuteQuery(ImageQry.String(), v["id"], id)
			if util.CheckErrorLog(r, err) {
				return nil, err
			}
			if len(ImgMap) > 0 {
				v["image"] = ImgMap
			} else {
				v["image"] = []string{}
			}

		}
	}

	var AmenityArrData = make(map[string]interface{})
	AmenityArrData, err = AmenityTypeWiseAmenityAdmin(r, id)
	if util.CheckErrorLog(r, err) {
		return nil, err
	}

	resMap["hotel_amenity"] = AmenityArrData

	var CancelPolicyQry bytes.Buffer
	CancelPolicyQry.WriteString("SELECT * FROM cf_cancellation_policy WHERE status = 1 AND hotel_id = ?")
	RetMap, err := ExecuteQuery(CancelPolicyQry.String(), id)
	if util.CheckErrorLog(r, err) {
		return nil, err
	}

	var TaxQry bytes.Buffer
	TaxQry.WriteString(" SELECT ")
	TaxQry.WriteString(" CHT.id, CHT.tax, CHT.description, CHT.type, CHT.amount ")
	TaxQry.WriteString(" FROM cf_hotel_tax AS CHT ")
	TaxQry.WriteString(" WHERE CHT.hotel_id = ? ")
	TaxMap, err := ExecuteQuery(TaxQry.String(), id)
	if util.CheckErrorLog(r, err) {
		return nil, err
	}

	resMap["hotel"] = HotelInfoMap
	resMap["setting"] = SettingMap
	resMap["image"] = CategoryListing
	resMap["cancellation_policy"] = RetMap
	resMap["room_list"], _ = GetRoomTypeList(r, id)
	resMap["tax_list"] = TaxMap

	return resMap, nil
}
