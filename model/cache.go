package model

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"runtime"
	"strconv"
	"strings"
	"time"
	"tp-api-common/data"
	"tp-api-common/util"
	"tp-system/config"

	"gopkg.in/mgo.v2/bson"
)

var SyncActiveProperty = func() bool {
	var Qry bytes.Buffer
	Qry.WriteString(" SELECT id, hotel_name ")
	Qry.WriteString(" FROM cf_hotel_info ")
	//Qry.WriteString(" WHERE is_approved = 1 AND status = 1 ")
	Qry.WriteString(" WHERE id = 100071932 ")
	HotelList, err := ExecuteQuery(Qry.String())
	if chkError(err) {
		util.SysLogIt("error in sync active property 1.1")
		return false
	}

	if len(HotelList) > 0 {
		for _, v := range HotelList {
			util.SysLogIt("Start Sync All Active Property ID - " + v["id"].(string))
			AddUpdatePropertyDetailInfo(v["id"].(string))
			util.SysLogIt("End Sync All Active Property ID - " + v["id"].(string))
		}
	}

	return true
}

//AddUpdatePropertyDetailInfo - Call update property detail info with error handling
func AddUpdatePropertyDetailInfo(id string) {
	defer func() {
		if err := recover(); err != nil {
			chkError(err.(error))
			util.SysLogIt("Error in add update property detail info - id - " + id)
			util.SysLogIt(err.(error))
		}
	}()
	UpdateHotelWithProperty(id)
}

func IsHotelApproved(id string) bool {
	var Qry bytes.Buffer
	Qry.WriteString(" SELECT * FROM cf_hotel_info WHERE is_approved = 1 AND id = ? ")
	Cnt, err := ExecuteQuery(Qry.String(), id)
	if chkError(err) {
		return false
	}

	if len(Cnt) == 0 {
		return false
	}

	return true
}

// UpdateSearchList -  It updates listing of hotel. As when hotel added to cache, we need to maintain it's search criteria. This API updates it's all search criteria based on operations as `List` and `Unlist`.
func UpdateSearchList(id string, operation string, callFlag ...int) bool {
	// log.Println("UpdateSearchList")
	if IsHotelApproved(id) {
		util.SysLogIt("UpdateSearchList Start")
		var Qry bytes.Buffer
		Qry.WriteString(" SELECT CHI.hotel_name, IFNULL(CFCC.name,'') AS country_name, IFNULL(CFC.name,'') AS city_name, ")
		Qry.WriteString(" IFNULL(CFS.name,'') AS state_name, IFNULL(CFL.locality,'') AS locality_name FROM cf_hotel_info AS CHI LEFT JOIN cf_city AS CFC ON CFC.id = CHI.city_id")
		Qry.WriteString(" LEFT JOIN cf_country AS CFCC ON CFCC.id = CHI.country_id")
		Qry.WriteString(" LEFT JOIN cf_states AS CFS on CFS.id = CHI.state_id")
		// Qry.WriteString(" LEFT JOIN cf_locality AS CFL ON CFL.id = CHI.locality_id WHERE CHI.is_live=1 AND CHI.id=?")
		Qry.WriteString(" LEFT JOIN cf_locality AS CFL ON CFL.id = CHI.locality_id WHERE CHI.id=?")
		hotel, err := ExecuteRowQuery(Qry.String(), id)
		if chkError(err) || hotel == nil || len(hotel) == 0 {
			return false
		}
		var sList = data.AddressUpdate{
			Operation: operation,
			HotelName: hotel["hotel_name"].(string),
			Country:   hotel["country_name"].(string),
			State:     hotel["state_name"].(string),
			City:      hotel["city_name"].(string),
			Locality:  hotel["locality_name"].(string),
		}
		// log.Println(sList)
		util.SysLogIt("sList")
		util.SysLogIt(sList)
		jsonList, _ := json.Marshal(sList)

		_, scode := SendMicroServiceRequest("POST", "updateHotelListing", string(jsonList))
		if scode != 200 {
			return false
		}
		util.SysLogIt("UpdateSearchList End")

		if operation == "List" && len(callFlag) == 0 {
			CacheChn <- CacheObj{
				Type: "updateHotelWithProperty",
				ID:   id,
			}
		}

		util.SysLogIt("UpdateSearchList For Hotel start")
		CacheChn <- CacheObj{
			Type: "changePropertyStatus",
			ID:   id,
		}
		util.SysLogIt("UpdateSearchList For Hotel end")
	}
	return true
}

// DropSearchList - Call this function when you want to remove collection : [search_list]. Mostly this api will be called when admin must reset of searchlist of all hotels
func DropSearchList() bool {
	_, scode := SendMicroServiceRequest("POST", "dropSearchList", "")
	if scode != 200 {
		return false
	}
	return true
}

func UpdatePropertyFlag(id string) bool {
	util.SysLogIt("UpdatePropertyFlag Start")
	var Qry bytes.Buffer
	Qry.WriteString("SELECT is_live,status FROM cf_hotel_info WHERE id = ? ")
	Data, err := ExecuteRowQuery(Qry.String(), id)
	if chkError(err) {
		return false
	}

	status := Data["status"].(int64)
	propertyFlag := int(Data["is_live"].(int64))
	var flag int
	switch status {
	case 1:
		flag = propertyFlag
	case 2, 4, 3:
		flag = 2
	default:
		flag = 2
	}

	var sList = data.HotelUpdateReq{
		HotelID:        id,
		PropertyStatus: flag,
	}
	// log.Println(sList)
	util.SysLogIt("sList")
	util.SysLogIt(sList)
	jsonList, _ := json.Marshal(sList)

	_, scode := SendMicroServiceRequest("POST", "updatePropertyFlag", string(jsonList))
	if scode != 200 {
		return false
	}
	util.SysLogIt("UpdatePropertyFlag End")

	return true
}

func UpdateHotelWithProperty(id string) {
	util.SysLogIt("Start UpdateHotelWithProperty - " + id)
	util.SysLogIt("Lock main list")
	UpdateHotelOnList(id)
	util.SysLogIt("unLock main list")

	util.SysLogIt("1 S")
	UpdateHotelAmenity(id)
	util.SysLogIt("1 E")

	util.SysLogIt("11 S")
	UpdateHotelDetailsAmenity(id)
	util.SysLogIt("11 E")

	util.SysLogIt("12 S")
	UpdateHotelTag(id)
	util.SysLogIt("12 E")

	util.SysLogIt("13 S")
	UpdateHotelImage(id)
	util.SysLogIt("13 E")

	util.SysLogIt("14 S")
	UpdateDetailedImage(id)
	util.SysLogIt("14 E")

	util.SysLogIt("15 S")
	UpdateHotelTax(id)
	util.SysLogIt("15 E")

	util.SysLogIt("16 S")
	UpdateAllRoomOfHotel(id)
	util.SysLogIt("16 E")

	util.SysLogIt("17 S")
	AddUpdateAllRatePlansOfHotel(id)
	util.SysLogIt("17 E")

	util.SysLogIt("18 S")
	//UpdateSearchList(id, "List", 1)
	util.SysLogIt("18 E")
	util.SysLogIt("END UpdateHotelWithProperty - " + id)
}

// UpdateHotelOnList - Update caching based on hotel deals update
func UpdateHotelOnList(id string) bool {
	//log.Println("UpdateHotelOnList")
	util.SysLogIt("UpdateHotelOnList Start")
	var Qry bytes.Buffer
	Qry.WriteString(" SELECT CHI.stripe_account, CHI.is_stripe_bank_created, CHI.is_live, CHI.id AS hotel_id, CHI.hotel_name, CHI.short_address, CHI.hotel_star, CHI.latitude, COALESCE(CHI.description, '') AS description,CHI.hotel_phone, ")
	Qry.WriteString(" CHI.longitude, CCO.name AS country_name , CCI.name AS city_name, CFS.name AS state_name,CHI.long_address,")
	Qry.WriteString(" IFNULL(CFL.locality,'') AS locality_name , CHI.property_type_id, CPT.type AS property_type, IFNULL(CHI.policy,'') AS policy, IFNULL(CHI.checkin_rules,'') AS checkin_rules, CHS.checkin_time, CHS.checkout_time")
	Qry.WriteString(" FROM cf_hotel_info AS CHI ")
	Qry.WriteString(" LEFT JOIN cf_country AS CCO On CCO.id = CHI.country_id ")
	Qry.WriteString(" LEFT JOIN cf_city AS CCI ON CCI.id = CHI.city_id")
	Qry.WriteString(" LEFT JOIN cf_states AS CFS ON CFS.id = CHI.state_id ")
	Qry.WriteString(" LEFT JOIN cf_property_type AS CPT ON CPT.id = CHI.property_type_id ")
	Qry.WriteString(" LEFT JOIN cf_hotel_settings AS CHS ON CHS.hotel_id = CHI.id ")
	Qry.WriteString(" LEFT JOIN cf_locality AS CFL ON CFL.id = CHI.locality_id WHERE CHI.status=1")
	if id != "" {
		Qry.WriteString(" AND CHI.id=? ")
	}

	if IsHotelApproved(id) {
		if id != "" {
			HotelInfo, err := ExecuteRowQuery(Qry.String(), id)
			if chkError(err) || HotelInfo == nil {
				return false
			}

			var hInfo = data.HotelUpdateReq{
				HotelID:          id,
				HotelName:        HotelInfo["hotel_name"].(string),
				Description:      HotelInfo["description"].(string),
				State:            HotelInfo["state_name"].(string),
				City:             HotelInfo["city_name"].(string),
				Locality:         HotelInfo["locality_name"].(string),
				LongAddress:      HotelInfo["long_address"].(string),
				ShortAddress:     HotelInfo["short_address"].(string),
				HotelPhone:       HotelInfo["hotel_phone"].(string),
				HotelStar:        int(HotelInfo["hotel_star"].(int64)),
				Lat:              HotelInfo["latitude"].(string),
				Long:             HotelInfo["longitude"].(string),
				PropertyType:     HotelInfo["property_type"].(string),
				PropertyTypeID:   HotelInfo["property_type_id"].(string),
				Policy:           HotelInfo["policy"].(string),
				CheckinRules:     HotelInfo["checkin_rules"].(string),
				CheckinTime:      HotelInfo["checkin_time"].(string),
				CheckoutTime:     HotelInfo["checkout_time"].(string),
				PropertyStatus:   int(HotelInfo["is_live"].(int64)),
				StripeAccount:    HotelInfo["stripe_account"].(string),
				StripeBankStatus: int(HotelInfo["is_stripe_bank_created"].(int64)),
			}
			// log.Println(hInfo)
			util.SysLogIt("hInfo")
			util.SysLogIt(hInfo)
			jsonList, _ := json.Marshal(hInfo)
			_, scode := SendMicroServiceRequest("POST", "addUpdateHotelInfo", string(jsonList))
			if scode != 200 {
				return false
			}
			// UpdateHotelAmenity(id)
			// UpdateHotelTag(id)
			// UpdateHotelImage(id)
			// UpdateDetailedImage(id)
			// UpdateHotelTax(id)
			util.SysLogIt("UpdateHotelOnList End")
		} else {
			HotelInfo, err := ExecuteQuery(Qry.String())
			if chkError(err) || HotelInfo == nil {
				return false
			}

			for _, val := range HotelInfo {
				hotelStar, _ := strconv.Atoi(val["hotel_star"].(string))
				hotelPropertyStatus, _ := strconv.Atoi(val["is_live"].(string))
				stripeBankStatus, _ := strconv.Atoi(val["is_stripe_bank_created"].(string))
				var hInfo = data.HotelUpdateReq{
					HotelID:          val["hotel_id"].(string),
					HotelName:        val["hotel_name"].(string),
					Description:      val["description"].(string),
					State:            val["state_name"].(string),
					City:             val["city_name"].(string),
					Locality:         val["locality_name"].(string),
					LongAddress:      val["long_address"].(string),
					ShortAddress:     val["short_address"].(string),
					HotelPhone:       val["hotel_phone"].(string),
					HotelStar:        hotelStar,
					Lat:              val["latitude"].(string),
					Long:             val["longitude"].(string),
					PropertyType:     val["property_type"].(string),
					PropertyTypeID:   val["property_type_id"].(string),
					Policy:           val["policy"].(string),
					CheckinRules:     val["checkin_rules"].(string),
					CheckinTime:      val["checkin_time"].(string),
					CheckoutTime:     val["checkout_time"].(string),
					PropertyStatus:   hotelPropertyStatus,
					StripeAccount:    val["stripe_account"].(string),
					StripeBankStatus: stripeBankStatus,
				}
				// log.Println(hInfo)
				util.SysLogIt("hInfo")
				util.SysLogIt(hInfo)
				jsonList, _ := json.Marshal(hInfo)
				_, scode := SendMicroServiceRequest("POST", "addUpdateHotelInfo", string(jsonList))
				if scode != 200 {
					return false
				}
			}
		}
	}
	return true
}

// UpdateHotelAmenity - Update amenity in hotel List
func UpdateHotelAmenity(id string) bool {
	util.SysLogIt("UpdateHotelAmenity Start")
	var QryAmenity bytes.Buffer
	QryAmenity.WriteString("SELECT CFHA.amenity_id, CFA.name, CFA.icon FROM cf_hotel_amenities AS CFHA INNER JOIN cf_amenity AS CFA ON CFA.id = CFHA.amenity_id ")
	QryAmenity.WriteString(" WHERE CFA.is_star_amenity=0 AND CFHA.hotel_id=?")
	hotelAmenity, err := ExecuteQuery(QryAmenity.String(), id)
	if chkError(err) {
		return false
	}
	var sAmenityArr []data.StarAmenity
	for _, val := range hotelAmenity {
		sAmenityArr = append(sAmenityArr, data.StarAmenity{ID: val["amenity_id"].(string), Name: val["name"].(string), Icon: val["icon"].(string)})
	}

	if IsHotelApproved(id) {
		if len(sAmenityArr) > 0 {
			rMap := data.HotelAmenityUpdateReq{
				HotelID: id,
				Amanity: sAmenityArr,
			}
			jsonList, _ := json.Marshal(rMap)
			_, scode := SendMicroServiceRequest("POST", "updateHotelAmenity", string(jsonList))
			if scode != 200 {
				return false
			}
		}

		//Added by meet soni for update detail amentiy in monfoDB
		if newflg := UpdateHotelDetailsAmenity(id); !newflg {
			return false
		}
	}

	util.SysLogIt("UpdateHotelAmenity End")
	return true
}

// UpdateHotelDetailsAmenity - Update details amenity in hotel list
func UpdateHotelDetailsAmenity(hotelID string) bool {
	util.SysLogIt("UpdateHotelDetailsAmenity Start")
	var Qry bytes.Buffer
	Qry.WriteString(" SELECT CAT.type, CAT.id AS amenity_type_id ")
	Qry.WriteString(" FROM cf_hotel_amenities AS CRA  ")
	Qry.WriteString(" INNER JOIN cf_amenity AS CA ON CRA.amenity_id = CA.id ")
	Qry.WriteString(" INNER JOIN cf_amenity_type AS CAT ON CAT.id = CA.amenity_type_id AND CAT.amenity_of = 1")
	Qry.WriteString(" WHERE CRA.hotel_id=? GROUP BY CAT.id")
	TypeInfo, err := ExecuteQuery(Qry.String(), hotelID)
	if chkError(err) {
		return false
	}
	var amenityMap []data.AmenityTypeMap
	for _, val := range TypeInfo {
		var AmenityQry bytes.Buffer
		AmenityQry.WriteString("SELECT CA.id AS amenity_id, CA.name AS amenity_name, CA.icon ")
		AmenityQry.WriteString("FROM cf_hotel_amenities AS CRA ")
		AmenityQry.WriteString("INNER JOIN cf_amenity AS CA ON CRA.amenity_id = CA.id AND CA.status=1 ")
		AmenityQry.WriteString("WHERE CRA.hotel_id=? AND CA.amenity_type_id = ?")
		amenityInfo, err := ExecuteQuery(AmenityQry.String(), hotelID, val["amenity_type_id"].(string))
		if chkError(err) {
			return false
		}
		amenityMap = append(amenityMap, data.AmenityTypeMap{
			TypeID:  val["amenity_type_id"].(string),
			Type:    val["type"].(string),
			Amenity: amenityInfo,
		})
	}
	if IsHotelApproved(hotelID) {
		if len(amenityMap) > 0 {
			reqMap := data.RoomAmenityUpdateReq{
				HotelID: hotelID,
				Aminity: amenityMap,
			}
			jsonList, _ := json.Marshal(reqMap)
			_, scode := SendMicroServiceRequest("POST", "updateHotelDetailAmenity", string(jsonList))
			if scode != 200 {
				return false
			}
		}
	}
	util.SysLogIt("updateHotelDetailAmenity End")
	return true
}

// UpdateHotelTag - Update hotel tag
func UpdateHotelTag(id string) bool {
	util.SysLogIt("UpdateHotelTag Start")
	var QryTag bytes.Buffer
	QryTag.WriteString("SELECT CT.id, CT.tag FROM cf_hotel_tag AS CHT INNER JOIN cf_tags AS CT ON CT.id = CHT.tag_id WHERE CHT.hotel_id=? AND CT.status=1")
	hotelTags, err := ExecuteQuery(QryTag.String(), id)
	if chkError(err) {
		return false
	}
	var hTags []data.HotelTag
	for _, val := range hotelTags {
		hTags = append(hTags, data.HotelTag{ID: val["id"].(string), Tag: val["tag"].(string)})
	}
	if IsHotelApproved(id) {
		if len(hTags) > 0 {
			rMap := data.HotelTagsUpdateReq{
				HotelID: id,
				Tags:    hTags,
			}
			jsonList, _ := json.Marshal(rMap)
			_, scode := SendMicroServiceRequest("POST", "updateHotelTags", string(jsonList))
			if scode != 200 {
				return false
			}
		}
	}
	util.SysLogIt("UpdateHotelTag End")
	return true
}

// UpdateHotelImage - Update hotel images
func UpdateHotelImage(id string) bool {
	util.SysLogIt("UpdateHotelImage Start")
	var QryImage bytes.Buffer
	QryImage.WriteString("SELECT image FROM cf_hotel_image WHERE hotel_id=? ORDER BY sortorder ASC")
	hotelImages, err := ExecuteQuery(QryImage.String(), id)
	if chkError(err) {
		return false
	}
	var hImages []string
	for _, val := range hotelImages {
		hImages = append(hImages, val["image"].(string))
	}
	if IsHotelApproved(id) {
		if len(hImages) > 0 {
			rMap := data.HotelImageUpdateReq{
				HotelID: id,
				Images:  hImages,
			}
			jsonList, _ := json.Marshal(rMap)
			_, scode := SendMicroServiceRequest("POST", "updateHotelImages", string(jsonList))
			if scode != 200 {
				return false
			}
		}
	}
	util.SysLogIt("UpdateHotelImage End")
	return true
}

// UpdateDetailedImage - Update detailed category wise images
func UpdateDetailedImage(id string) bool {
	util.SysLogIt("UpdateDetailedImage Start")
	/* c := VMongoSession.DB(config.Env.Mongo.MongoDB).C("hotel_list")
	cnt, err := c.Find(bson.M{"hotel_id": id}).Count()
	util.SysLogIt("error0")
	util.SysLogIt(err)
	util.SysLogIt("cnt")
	util.SysLogIt(cnt)
	if chkError(err) || cnt != 1 {
		return false
	}*/
	var retBody []data.HotelCategoryImage
	var QryImage bytes.Buffer
	QryImage.WriteString("SELECT CIC.id, CIC.name FROM cf_hotel_image AS CFHI  ")
	QryImage.WriteString("INNER JOIN cf_image_category AS CIC ON CIC.id = CFHI.category_id ")
	QryImage.WriteString("WHERE CFHI.hotel_id=? GROUP BY CIC.id")
	catInfo, err := ExecuteQuery(QryImage.String(), id)
	util.SysLogIt("catInfo")
	util.SysLogIt(catInfo)
	if chkError(err) {
		util.SysLogIt("error1")
		util.SysLogIt(err)
		return false
	}
	if IsHotelApproved(id) {
		if len(catInfo) > 0 {
			for _, v := range catInfo {
				var Qry bytes.Buffer
				Qry.WriteString("SELECT image FROM cf_hotel_image WHERE hotel_id = ? AND category_id = ? ORDER BY sortorder ASC")
				imageArr, err := ExecuteQuery(Qry.String(), id, v["id"].(string))
				util.SysLogIt("imageArr")
				util.SysLogIt(imageArr)
				if chkError(err) {
					util.SysLogIt("error2")
					util.SysLogIt(err)
					return false
				}
				var hImages []string
				for _, val := range imageArr {
					hImages = append(hImages, val["image"].(string))
				}
				retBody = append(retBody, data.HotelCategoryImage{ID: v["id"].(string), Category: v["name"].(string), Images: hImages})
			}
			rMap := data.HotelCategoryImageUpdateReq{
				HotelID:        id,
				CategoryImages: retBody,
			}
			// log.Println(rMap)
			util.SysLogIt("rMap")
			util.SysLogIt(rMap)
			jsonList, _ := json.Marshal(rMap)
			_, scode := SendMicroServiceRequest("POST", "updateHotelCategoryImages", string(jsonList))
			if scode != 200 {
				return false
			}
		}
	}
	util.SysLogIt("UpdateDetailedImage End")
	return true
}

// UpdateHotelTax - Set hotel taxes on bookings
func UpdateHotelTax(id string) bool {
	util.SysLogIt("UpdateHotelTax Start")

	var QryImage bytes.Buffer
	// QryImage.WriteString("SELECT id, tax, type, amount FROM cf_hotel_tax WHERE status=1 AND hotel_id=?;")
	QryImage.WriteString("SELECT id, tax, type, amount FROM cf_hotel_tax WHERE hotel_id=? AND status = 1;")
	taxes, err := ExecuteQuery(QryImage.String(), id)
	if chkError(err) {
		return false
	}
	var taxArr []data.HotelTax
	for _, val := range taxes {
		taxArr = append(taxArr, data.HotelTax{ID: val["id"].(string), Tax: val["tax"].(string), Type: val["type"].(string), Amount: val["amount"].(string)})
	}
	if IsHotelApproved(id) {
		if len(taxArr) > 0 {
			rMap := data.HotelTaxUpdateReq{
				HotelID: id,
				Tax:     taxArr,
			}
			jsonList, _ := json.Marshal(rMap)
			_, scode := SendMicroServiceRequest("POST", "updateHotelTax", string(jsonList))
			if scode != 200 {
				return false
			}
		}
	}
	util.SysLogIt("UpdateHotelTax End")
	return true
}

// DeleteHotelItem - This Function is used to remove hotel info, room types of hotel and rateplan of hotel based on your wish. you need to send true for settings of struct.
func DeleteHotelItem(deleteData data.DeleteCacheItem) bool {
	util.SysLogIt("DeleteHotelItem Start")
	jsonList, _ := json.Marshal(deleteData)
	_, scode := SendMicroServiceRequest("POST", "deleteHotelItem", string(jsonList))
	if scode != 200 {
		return false
	}
	util.SysLogIt("DeleteHotelItem End")
	return true
}

var AddUpdateAllRatePlansOfHotel = func(id string) bool {
	var Qry bytes.Buffer
	Qry.WriteString(" SELECT id,room_type_id FROM cf_rateplan WHERE status = 1 AND hotel_id = ? ")
	RatePlan, err := ExecuteQuery(Qry.String(), id)
	if chkError(err) {
		util.SysLogIt("AddUpdateAllRatePlansOfHotel 1.1 error")
		return false
	}

	if len(RatePlan) > 0 {
		for _, v := range RatePlan {
			if !AddUpdateRateplanDetails(id, v["id"].(string)) {
				util.SysLogIt("AddUpdateAllRatePlansOfHotel 1.2 error AddUpdateRateplanDetails")
				return false
			}

			if !UpdateRatePlanDeals(id, v["room_type_id"].(string), v["id"].(string)) {
				util.SysLogIt("AddUpdateAllRatePlansOfHotel 1.2 error UpdateRatePlanDeals")
				return false
			}
		}
	}

	return true
}

// AddUpdateRateplanDetails - Update Rate Plan Basic Details
func AddUpdateRateplanDetails(hid string, rpid string) bool {
	util.SysLogIt("UpdateAllRateplanDeals Start")
	if IsHotelApproved(hid) {
		var Qry bytes.Buffer
		Qry.WriteString("SELECT CFR.id AS rateplan_id, CFR.rate_plan_name, CFR.room_type_id, CFR.is_pay_at_hotel, ")
		Qry.WriteString("CASE WHEN CCP.is_non_refundable = 0 AND CCP.before_day_charge = 0 THEN 1 ELSE 0 END AS is_free_cancellation, CMT.meal_type ")
		Qry.WriteString(",CFR.cancellation_policy_id ")
		Qry.WriteString("FROM cf_rateplan AS CFR ")
		Qry.WriteString(" INNER JOIN cf_cancellation_policy AS CCP ON CCP.id = CFR.cancellation_policy_id ")
		Qry.WriteString(" INNER JOIN cf_meal_type AS CMT ON CMT.id = CFR.meal_type_id")
		// Qry.WriteString(" WHERE CFR.hotel_id=? AND CFR.status=1 AND CFR.id=?")
		Qry.WriteString(" WHERE CFR.hotel_id=? AND CFR.id=?")
		RatePlans, err := ExecuteRowQuery(Qry.String(), hid, rpid)
		if chkError(err) || RatePlans == nil {
			util.SysLogIt("UpdateAllRateplanDeals 1.1 error")
			return false
		}
		var HQry bytes.Buffer
		HQry.WriteString(" SELECT CHI.hotel_name , CCO.name AS country_name , CCI.name AS city_name, CFS.name AS state_name,")
		HQry.WriteString(" IFNULL(CFL.locality,'') AS locality_name")
		HQry.WriteString(" FROM cf_hotel_info AS CHI ")
		HQry.WriteString(" LEFT JOIN cf_country AS CCO On CCO.id = CHI.country_id ")
		HQry.WriteString(" LEFT JOIN cf_city AS CCI ON CCI.id = CHI.city_id")
		HQry.WriteString(" LEFT JOIN cf_states AS CFS ON CFS.id = CHI.state_id ")
		//HQry.WriteString(" LEFT JOIN cf_locality AS CFL ON CFL.id = CHI.locality_id WHERE CHI.id=? AND CHI.status=1")
		HQry.WriteString(" LEFT JOIN cf_locality AS CFL ON CFL.id = CHI.locality_id WHERE CHI.id=?")
		HotelInfo, err := ExecuteRowQuery(HQry.String(), hid)
		if chkError(err) || HotelInfo == nil {
			util.SysLogIt("UpdateAllRateplanDeals 1.2 error")
			return false
		}

		iph, ifc, ifb := true, true, true
		if int(RatePlans["is_pay_at_hotel"].(int64)) != 1 {
			iph = false
		}
		if int(RatePlans["is_free_cancellation"].(int64)) != 1 {
			ifc = false
		}
		if strings.Contains(RatePlans["meal_type"].(string), "breakfast") {
			ifb = true
		}

		CQry := "SELECT is_non_refundable, before_day, before_day_charge, after_day_charge FROM cf_cancellation_policy WHERE id=? AND status=1"
		CPolicy, err := ExecuteRowQuery(CQry, RatePlans["cancellation_policy_id"].(string))
		if chkError(err) || CPolicy == nil {
			util.SysLogIt("UpdateAllRateplanDeals 1.3 error")
			return false
		}

		var IncQry bytes.Buffer
		IncQry.WriteString("SELECT CFI.inclusion FROM cf_rateplan_inclusion AS CRI ")
		IncQry.WriteString("INNER JOIN cf_inclusion CFI ON CFI.id = CRI.inclusion_id ")
		IncQry.WriteString("WHERE CRI.hotel_id=? AND CRI.rateplan_id=? ORDER BY CRI.sortorder ASC")
		InclusionList, err := ExecuteQuery(IncQry.String(), hid, RatePlans["rateplan_id"].(string))
		if chkError(err) {
			util.SysLogIt("UpdateAllRateplanDeals 1.4 error")
			return false
		}
		var IncList []string
		for _, val := range InclusionList {
			IncList = append(IncList, val["inclusion"].(string))
		}

		ReqMap := data.RPUpdateReq{
			HotelName:          HotelInfo["hotel_name"].(string),
			HotelID:            hid,
			CountryName:        HotelInfo["country_name"].(string),
			CityName:           HotelInfo["city_name"].(string),
			StateName:          HotelInfo["state_name"].(string),
			LocalityName:       HotelInfo["locality_name"].(string),
			RatePlanID:         RatePlans["rateplan_id"].(string),
			RoomTypeID:         RatePlans["room_type_id"].(string),
			IsPayAtHotel:       iph,
			IsFreeCancellation: ifc,
			IsFreeBreakfast:    ifb,
			CancellationPolicy: data.CancellationPolicy{
				IsNonRefundable: int(CPolicy["is_non_refundable"].(int64)),
				BeforeDay:       int(CPolicy["before_day"].(int64)),
				BeforeDayCharge: CPolicy["before_day_charge"].(string),
				AfterDayCharge:  CPolicy["after_day_charge"].(string),
			},
			Inclusion:    IncList,
			RatePlanName: RatePlans["rate_plan_name"].(string),
		}

		jsonList, _ := json.Marshal(ReqMap)
		_, scode := SendMicroServiceRequest("POST", "updateRatePlan", string(jsonList))
		if scode != 200 {
			util.SysLogIt("UpdateAllRateplanDeals 1.5 error")
			return false
		}
	}

	util.SysLogIt("UpdateAllRateplanDeals End")
	return true
}

// UpdateRatePlanDeals - Add RatePlan Deals
func UpdateRatePlanDeals(HotelID string, RoomTypeID string, RatePlanID string) bool {
	util.SysLogIt("UpdateRatePlanDeals Start")
	if IsHotelApproved(HotelID) {
		util.SysLogIt("HotelID")
		util.SysLogIt(HotelID)
		util.SysLogIt("RoomTypeID")
		util.SysLogIt(RoomTypeID)
		util.SysLogIt("RatePlanID")
		util.SysLogIt(RatePlanID)

		currentTime, _ := time.Parse("2006-01-02", time.Now().Format("2006-01-02"))
		var Qry bytes.Buffer
		// Qry.WriteString("SELECT year, month, inv_data FROM cf_inv_data WHERE hotel_id=? AND room_id=? ORDER BY year, month ASC")
		// RoomInfo, err := ExecuteQuery(Qry.String(), HotelID, RoomTypeID)

		// 2021-04-26 - HK - START
		// Purpose : To stop getting synced old data upon activation
		currentTimeMonthYear := time.Now()
		localDate := currentTimeMonthYear.Format("2006-01-02")
		currentYr := strings.Split(localDate, "-")[0]
		currentMonth := strings.Split(localDate, "-")[1]
		Qry.WriteString("SELECT year, month, inv_data FROM cf_inv_data WHERE hotel_id=? AND room_id=? AND year >= ? AND month >= ? ORDER BY year, month ASC")
		RoomInfo, err := ExecuteQuery(Qry.String(), HotelID, RoomTypeID, currentYr, currentMonth)
		// 2021-04-26 - HK - END
		util.SysLogIt("RoomInfo")
		util.SysLogIt(RoomInfo)

		if chkError(err) {
			util.SysLogIt("err0")
			util.SysLogIt(err)
			return false
		}

		if RoomInfo == nil || len(RoomInfo) == 0 {
			util.SysLogIt("room info nill")
			util.SysLogIt(err)
			return false
		}
		var dealArr []data.Rate
		for _, val := range RoomInfo {
			var RPQry bytes.Buffer
			RPQry.WriteString("SELECT rate_rest_data FROM cf_rate_restriction_data_2 WHERE hotel_id=? AND room_id=? AND rateplan_id=? AND year=? AND month=?")
			RateInfo, err := ExecuteRowQuery(RPQry.String(), HotelID, RoomTypeID, RatePlanID, val["year"], val["month"])
			util.SysLogIt("RateInfo")
			util.SysLogIt(RateInfo)
			if chkError(err) || RateInfo == nil {
				util.SysLogIt("err of RateInfo nill")
				util.SysLogIt(err)
				return false
			}
			rateMap := make(map[string]interface{})
			err = json.Unmarshal([]byte(RateInfo["rate_rest_data"].(string)), &rateMap)
			if chkError(err) {
				util.SysLogIt("err while marshal")
				util.SysLogIt(err)
				continue
			}
			jsonMap := make(map[string]int64)
			err = json.Unmarshal([]byte(val["inv_data"].(string)), &jsonMap)
			if chkError(err) {
				util.SysLogIt("err while unmarshal")
				util.SysLogIt(err)
				continue
			}

			util.SysLogIt("jsonMap")
			util.SysLogIt(jsonMap)
			for k, v := range jsonMap {
				if _, ok := rateMap[k].(map[string]interface{}); ok {
					t, _ := time.Parse("2006-01-02", k)
					if t.Before(currentTime) {
						continue
					}
					rateinfo := rateMap[k].(map[string]interface{})
					rateMapFiltered := make(map[string]float64)
					for _, rateval := range rateinfo["rate"].([]interface{}) {
						for ratekey1, rateval1 := range rateval.(map[string]interface{}) {
							ratekey1 = strings.Replace(ratekey1, "occ_", "", -1)
							var fRate float64
							//fRate, _ = strconv.ParseFloat(rateval1.(string), 64)
							fRate = rateval1.(float64)
							rateMapFiltered[ratekey1] = fRate
						}
					}

					// var rate float64 = 500
					dealArr = append(dealArr, data.Rate{
						Date: k,
						Rate: rateMapFiltered,
						// Rate: map[string]float64{"1": rate, "2": rate},
						// Rate:      rateinfo["rate"].([]interface{}),
						Inv:       int(v),
						StopSell:  int(rateinfo["stop_sell"].(float64)),
						MinNights: int(rateinfo["min_night"].(float64)),
						Coa:       int(rateinfo["cta"].(float64)),
						Cod:       int(rateinfo["ctd"].(float64)),
					})
				}
			}
		}

		if len(dealArr) > 0 {
			rpDealsMap := data.RPDealUpdateReq{
				HotelID:    HotelID,
				RatePlanID: RatePlanID,
				Deals:      dealArr,
			}
			jsonList, _ := json.Marshal(rpDealsMap)
			_, scode := SendMicroServiceRequest("POST", "updateRatePlanDeals", string(jsonList))
			if scode != 200 {
				return false
			}
		}
	}
	util.SysLogIt("UpdateRatePlanDeals End")
	return true
}

var UpdateAllRoomOfHotel = func(id string) bool {
	util.SysLogIt("UpdateAllRoomOfHotel Start")
	if IsHotelApproved(id) {
		var Qry bytes.Buffer
		Qry.WriteString("SELECT CRT.id AS room_type_id, CRT.room_type_name, CBT.bed_type, CRT.room_size, CRT.description,")
		Qry.WriteString(" CRT.is_extra_bed, COALESCE(CEBT.extra_bed_name, '') AS extra_bed_name ")
		Qry.WriteString(" FROM cf_room_type AS CRT ")
		Qry.WriteString(" INNER JOIN cf_bed_type AS CBT ON CBT.id = CRT.bed_type_id")
		Qry.WriteString(" INNER JOIN cf_room_view AS CRV ON CRV.id = CRT.room_view_id")
		Qry.WriteString(" LEFT JOIN cf_extra_bed_type AS CEBT ON CEBT.id = CRT.extra_bed_type_id")
		Qry.WriteString("  WHERE CRT.hotel_id=?")

		util.SysLogIt("Qry")
		util.SysLogIt(Qry.String())
		util.SysLogIt("Hotel Id")
		util.SysLogIt(id)
		RoomInfo, err := ExecuteQuery(Qry.String(), id)
		util.SysLogIt("RoomInfo")
		util.SysLogIt(RoomInfo)
		if chkError(err) || RoomInfo == nil {
			util.SysLogIt("Error0")
			util.SysLogIt(err)
			return false
		}

		if len(RoomInfo) > 0 {
			for _, v := range RoomInfo {
				reqMap := data.RoomTypeReq{
					HotelID:      id,
					RoomTypeID:   v["room_type_id"].(string),
					RoomTypeName: v["room_type_name"].(string),
					BedType:      v["bed_type"].(string),
					RoomSize:     v["room_size"].(string),
					Description:  v["description"].(string),
					IsExtraBed:   int(v["is_extra_bed"].(int64)),
					ExtraBedName: v["extra_bed_name"].(string),
				}
				jsonList, _ := json.Marshal(reqMap)
				_, scode := SendMicroServiceRequest("POST", "updateRoomType", string(jsonList))
				if scode != 200 {
					util.SysLogIt("Error1 UpdateAllRoomType Syncing Data")
					util.SysLogIt(scode)
					return false
				}

				if !UpdateRoomAmenity(id, v["room_type_id"].(string)) {
					util.SysLogIt("Error11 UpdateRoomAmenity Syncing Data in updateroomtypeofhotel")
					return false
				}

				if !UpdateRoomImage(id, v["room_type_id"].(string)) {
					util.SysLogIt("Error122 UpdateRoomImage syncing data in updatetoomtypeofhotel")
					return false
				}
			}
		}
	}
	util.SysLogIt("UpdateAllRoomOfHotel End")

	return true
}

// UpdateAllRoomType - Update all rateplans or given id only
func UpdateAllRoomType(id string, roomTypeID string) bool {
	util.SysLogIt("UpdateAllRoomType Start")
	var Qry bytes.Buffer
	var err error
	var RoomInfo map[string]interface{}
	Qry.WriteString("SELECT CRT.id AS room_type_id, CRT.room_type_name, CBT.bed_type, CRT.room_size, CRT.description,")
	// Qry.WriteString(" CRT.is_extra_bed, CEBT.extra_bed_name")
	Qry.WriteString(" CRT.is_extra_bed, COALESCE(CEBT.extra_bed_name, '') AS extra_bed_name ")
	Qry.WriteString(" FROM cf_room_type AS CRT ")
	Qry.WriteString(" INNER JOIN cf_bed_type AS CBT ON CBT.id = CRT.bed_type_id")
	Qry.WriteString(" INNER JOIN cf_room_view AS CRV ON CRV.id = CRT.room_view_id")
	Qry.WriteString(" LEFT JOIN cf_extra_bed_type AS CEBT ON CEBT.id = CRT.extra_bed_type_id")
	// Qry.WriteString("  WHERE CRT.status=1 AND CRT.hotel_id=?")
	Qry.WriteString("  WHERE CRT.hotel_id=?")
	Qry.WriteString("  AND CRT.id=?")

	util.SysLogIt("Qry")
	util.SysLogIt(Qry.String())
	util.SysLogIt("Hotel Id")
	util.SysLogIt(id)
	util.SysLogIt("room id")
	util.SysLogIt(roomTypeID)
	RoomInfo, err = ExecuteRowQuery(Qry.String(), id, roomTypeID)
	util.SysLogIt("RoomInfo")
	util.SysLogIt(RoomInfo)
	if chkError(err) || RoomInfo == nil {
		util.SysLogIt("Error0")
		util.SysLogIt(err)
		return false
	}

	if IsHotelApproved(id) {
		reqMap := data.RoomTypeReq{
			HotelID:      id,
			RoomTypeID:   roomTypeID,
			RoomTypeName: RoomInfo["room_type_name"].(string),
			BedType:      RoomInfo["bed_type"].(string),
			RoomSize:     RoomInfo["room_size"].(string),
			Description:  RoomInfo["description"].(string),
			IsExtraBed:   int(RoomInfo["is_extra_bed"].(int64)),
			ExtraBedName: RoomInfo["extra_bed_name"].(string),
		}
		jsonList, _ := json.Marshal(reqMap)
		_, scode := SendMicroServiceRequest("POST", "updateRoomType", string(jsonList))
		if scode != 200 {
			util.SysLogIt("Error1 UpdateAllRoomType Syncing Data")
			util.SysLogIt(scode)
			return false
		}
	}
	util.SysLogIt("UpdateAllRoomType End")
	return true
}

// UpdateRoomAmenity - Update room type id
func UpdateRoomAmenity(hotelID string, roomTypeID string) bool {
	util.SysLogIt("UpdateRoomAmenity Start")
	var Qry bytes.Buffer
	Qry.WriteString("SELECT CAT.type, CAT.id AS amenity_type_id ")
	Qry.WriteString(" FROM cf_room_amenity AS CRA  ")
	Qry.WriteString(" INNER JOIN cf_amenity AS CA ON CRA.amenity_id = CA.id ")
	Qry.WriteString(" INNER JOIN cf_amenity_type AS CAT ON CAT.id = CA.amenity_type_id AND CAT.amenity_of = 2") // CAT.amenity_of = 1
	Qry.WriteString(" WHERE CRA.room_type_id=? GROUP BY CAT.id")
	TypeInfo, err := ExecuteQuery(Qry.String(), roomTypeID)
	if chkError(err) {
		return false
	}
	var amenityMap []data.AmenityTypeMap
	for _, val := range TypeInfo {
		var AmenityQry bytes.Buffer
		AmenityQry.WriteString("SELECT CA.id AS amenity_id, CA.name AS amenity_name, CA.icon ")
		AmenityQry.WriteString("FROM cf_room_amenity AS CRA ")
		AmenityQry.WriteString("INNER JOIN cf_amenity AS CA ON CRA.amenity_id = CA.id AND CA.status=1 ")
		AmenityQry.WriteString("WHERE CA.amenity_type_id=? AND CRA.room_type_id = ?")
		amenityInfo, err := ExecuteQuery(AmenityQry.String(), val["amenity_type_id"].(string), roomTypeID)
		if chkError(err) {
			return false
		}
		amenityMap = append(amenityMap, data.AmenityTypeMap{
			TypeID:  val["amenity_type_id"].(string),
			Type:    val["type"].(string),
			Amenity: amenityInfo,
		})
	}
	if IsHotelApproved(hotelID) {
		if len(amenityMap) > 0 {
			reqMap := data.RoomAmenityUpdateReq{
				HotelID:    hotelID,
				RoomTypeID: roomTypeID,
				Aminity:    amenityMap,
			}
			jsonList, _ := json.Marshal(reqMap)
			_, scode := SendMicroServiceRequest("POST", "updateRoomTypeAmenity", string(jsonList))
			if scode != 200 {
				return false
			}
		}
	}
	util.SysLogIt("UpdateRoomAmenity End")
	return true
}

// UpdateRoomImage - Update Room type images
func UpdateRoomImage(hotelID string, roomTypeID string) bool {
	util.SysLogIt("UpdateRoomImage Start")
	var QryImage bytes.Buffer
	QryImage.WriteString("SELECT id, image FROM cf_room_image WHERE room_type_id=? ORDER BY sort_order ASC")
	roomImages, err := ExecuteQuery(QryImage.String(), roomTypeID)
	if chkError(err) {
		return false
	}
	var rImages []string
	for _, val := range roomImages {
		rImages = append(rImages, val["image"].(string))
	}
	if IsHotelApproved(hotelID) {
		if len(rImages) > 0 {
			reqMap := data.RoomImageUpdateReq{
				HotelID:    hotelID,
				RoomTypeID: roomTypeID,
				Images:     rImages,
			}
			jsonList, _ := json.Marshal(reqMap)
			_, scode := SendMicroServiceRequest("POST", "updateRoomTypeImage", string(jsonList))
			if scode != 200 {
				return false
			}
		}
	}
	util.SysLogIt("UpdateRoomImage End")
	return true
}

func chkError(err error) bool {
	if err != nil {
		_, file, no, ok := runtime.Caller(1)
		if ok {
			content := fmt.Sprint("Exception -  on File - ", file, " Line - ", no)
			c := VMongoSession.DB(config.Env.Mongo.MongoDB).C("cache_error")
			MID := bson.NewObjectId()
			c.Insert(&data.ErrLog{ID: MID, Content: content + " /n " + err.Error()})
			go util.SendLoggerMail(content, "TP SYSTEM API Cache Error Found At ")
		}
		return true
	}
	return false
}

// SendMicroServiceRequest - it sends  micro service request to update cache on front side.
func SendMicroServiceRequest(method string, url string, data string) (map[string]interface{}, int) {
	util.SysLogIt("SendMicroServiceRequest Start")
	serviceURL := config.Env.FrontURL + url
	util.SysLogIt("serviceURL")
	util.SysLogIt(serviceURL)
	util.SysLogIt("data")
	util.SysLogIt(data)
	var token bson.ObjectId
	var payload *strings.Reader
	if data == "" {
		payload = nil
	} else {
		payload = strings.NewReader(data)
	}
	client := &http.Client{}
	token = addRequestToken(url, data)
	req, err := http.NewRequest(method, serviceURL, payload)
	if chkError(err) {
		return nil, 500
	}
	req.Header.Add("Content-Type", "application/json")
	res, err := client.Do(req)
	if chkError(err) {
		return nil, 500
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if chkError(err) {
		return nil, 500
	}
	// log.Println("res Code :: ", res.StatusCode)
	util.SysLogIt("res Code")
	util.SysLogIt(res.StatusCode)
	if res.StatusCode == 200 {
		removeRequestToken(token)
	}
	jsonMap := make(map[string]interface{})
	if string(body) != "" {
		err = json.Unmarshal(body, &jsonMap)
		if chkError(err) {
			return nil, 500
		}
	}
	util.SysLogIt("SendMicroServiceRequest End")
	return jsonMap, res.StatusCode
}

// addRequestToken - Add request Token
func addRequestToken(url string, content string) bson.ObjectId {
	c := VMongoSession.DB(config.Env.Mongo.MongoDB).C("cache_request")
	MID := bson.NewObjectId()
	c.Insert(&data.HTTPRequestToken{ID: MID, URL: url, ReqBody: content, DateTime: util.GetIsoLocalDateTime()})
	return MID
}

// removeRequestToken - Remove request token
func removeRequestToken(id bson.ObjectId) {
	c := VMongoSession.DB(config.Env.Mongo.MongoDB).C("cache_request")
	c.Remove(bson.M{"_id": id})
	return
}

// UpdateRatePlanDealsPartnerPanel - Syncs Data Of The Month Changed In Partner Panel Deals Module
func UpdateRatePlanDealsPartnerPanel(HotelID string, RoomTypeID string, RatePlanID string, year float64, month float64) bool {

	util.SysLogIt("UpdateRatePlanDealsPartnerPanel Start")
	util.SysLogIt("HotelID")
	util.SysLogIt(HotelID)
	util.SysLogIt("RoomTypeID")
	util.SysLogIt(RoomTypeID)
	util.SysLogIt("RatePlanID")
	util.SysLogIt(RatePlanID)
	util.SysLogIt("Year")
	util.SysLogIt(year)
	util.SysLogIt("Month")
	util.SysLogIt(month)

	currentTime, _ := time.Parse("2006-01-02", time.Now().Format("2006-01-02"))
	var Qry bytes.Buffer
	Qry.WriteString("SELECT inv_data FROM cf_inv_data WHERE hotel_id=? AND room_id=? AND year=? AND month=?")
	RoomInfo, err := ExecuteQuery(Qry.String(), HotelID, RoomTypeID, year, month)
	util.SysLogIt("RoomInfo")
	util.SysLogIt(RoomInfo)

	if chkError(err) {
		util.SysLogIt("err0")
		util.SysLogIt(err)
		return false
	}

	if RoomInfo == nil || len(RoomInfo) == 0 {
		util.SysLogIt("room info nill")
		util.SysLogIt(err)
		return false
	}
	var dealArr []data.Rate
	for _, val := range RoomInfo {
		var RPQry bytes.Buffer
		RPQry.WriteString("SELECT rate_rest_data FROM cf_rate_restriction_data_2 WHERE hotel_id=? AND room_id=? AND rateplan_id=? AND year=? AND month=?")
		RateInfo, err := ExecuteRowQuery(RPQry.String(), HotelID, RoomTypeID, RatePlanID, year, month)
		util.SysLogIt("RateInfo")
		util.SysLogIt(RateInfo)
		if chkError(err) || RateInfo == nil {
			util.SysLogIt("err of RateInfo nill")
			util.SysLogIt(err)
			return false
		}
		rateMap := make(map[string]interface{})
		err = json.Unmarshal([]byte(RateInfo["rate_rest_data"].(string)), &rateMap)
		if chkError(err) {
			util.SysLogIt("err while marshal")
			util.SysLogIt(err)
			continue
		}
		jsonMap := make(map[string]int64)
		err = json.Unmarshal([]byte(val["inv_data"].(string)), &jsonMap)
		if chkError(err) {
			util.SysLogIt("err while unmarshal")
			util.SysLogIt(err)
			continue
		}

		util.SysLogIt("jsonMap")
		util.SysLogIt(jsonMap)
		for k, v := range jsonMap {
			if _, ok := rateMap[k].(map[string]interface{}); ok {
				t, _ := time.Parse("2006-01-02", k)
				if t.Before(currentTime) {
					continue
				}
				rateinfo := rateMap[k].(map[string]interface{})
				rateMapFiltered := make(map[string]float64)
				for _, rateval := range rateinfo["rate"].([]interface{}) {
					for ratekey1, rateval1 := range rateval.(map[string]interface{}) {
						ratekey1 = strings.Replace(ratekey1, "occ_", "", -1)
						var fRate float64
						fRate, _ = strconv.ParseFloat(rateval1.(string), 64)
						rateMapFiltered[ratekey1] = fRate
					}
				}

				// var rate float64 = 500
				dealArr = append(dealArr, data.Rate{
					Date: k,
					Rate: rateMapFiltered,
					// Rate: map[string]float64{"1": rate, "2": rate},
					// Rate:      rateinfo["rate"].([]interface{}),
					Inv:       int(v),
					StopSell:  int(rateinfo["stop_sell"].(float64)),
					MinNights: int(rateinfo["min_night"].(float64)),
					Coa:       int(rateinfo["cta"].(float64)),
					Cod:       int(rateinfo["ctd"].(float64)),
				})
			}
		}
	}
	if IsHotelApproved(HotelID) {
		if len(dealArr) > 0 {
			rpDealsMap := data.RPDealUpdateReq{
				HotelID:    HotelID,
				RatePlanID: RatePlanID,
				Deals:      dealArr,
			}
			jsonList, _ := json.Marshal(rpDealsMap)
			_, scode := SendMicroServiceRequest("POST", "updateRatePlanDeals", string(jsonList))
			if scode != 200 {
				return false
			}
		}
	}
	util.SysLogIt("UpdateRatePlanDeals End")
	return true
}
