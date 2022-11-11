package model

import (
	"bytes"
	"net/http"
	"tp-api-common/util"
	"tp-system/config"
)

// GetFrontSettings - Returns settings front
func GetFrontSettings(r *http.Request) (map[string]interface{}, error) {
	util.LogIt(r, "Model - front - GetFrontSettings")

	var Qry bytes.Buffer
	Qry.WriteString("SELECT id, `key`, par_value, description FROM cf_filteration_settings;")
	CntData, err := ExecuteQuery(Qry.String())
	if util.CheckErrorLog(r, err) {
		return nil, err
	}
	resBody := make(map[string]interface{})
	for _, v := range CntData {
		key := v["key"].(string)
		value := v["par_value"].(string)
		resBody[key] = value
		if key == "show_popular_filtration" && value == "yes" {
			resBody["popular_filtration"] = util.PopularFiltrations
		}
		if key == "show_customer_ratings" && value == "yes" {
			resBody["customer_ratings"] = util.CustomerRating
		}
		if key == "show_hotel_ratings" && value == "yes" {
			resBody["hotel_ratings"] = util.HotelRating
		}
		if key == "show_property_type" && value == "yes" {
			ptype, flg := GetPropertyType(r)
			if flg == false {
				return nil, err
			}
			resBody["property_type"] = ptype
		}
		if key == "show_amenities" && value == "yes" {
			amenity, err := GetStarAmenity(r)
			if err != nil {
				return nil, err
			}
			resBody["amenities"] = amenity
		}
		resBody["sorting"] = util.HListSort
	}
	resBody["bucket_url"] = config.Env.AwsBucketURL
	resBody["room_folder"] = config.Env.RoomFolder
	resBody["hotel_folder"] = config.Env.HotelFolder
	resBody["profile_folder"] = config.Env.ProfileFolder
	resBody["property_type_folder"] = config.Env.PropertyTypeFolder
	resBody["popular_city_folder"] = config.Env.PopularCityFolder

	return resBody, nil
}

// GetHomePageData - Returns home page data
func GetHomePageData(r *http.Request) (map[string]interface{}, error) {
	util.LogIt(r, "Model - front - GetHomePageData")

	// Property List Data
	var Qry bytes.Buffer
	retBody := make(map[string]interface{})
	Qry.WriteString("SELECT CPT.id, CPT.type, CPT.image, count(CHI.id) AS cnt FROM cf_property_type AS CPT ")
	Qry.WriteString("LEFT JOIN cf_hotel_info AS CHI ON CHI.property_type_id = CPT.id  WHERE CPT.status=1 ")
	Qry.WriteString("GROUP BY CPT.id;")
	PropertyData, err := ExecuteQuery(Qry.String())
	if util.CheckErrorLog(r, err) {
		return nil, err
	}
	retBody["propertyList"] = PropertyData

	// City List Data
	var CityQry bytes.Buffer
	CityQry.WriteString("SELECT id,city_id, city_name, image, description, sort_order FROM cms_popular_city WHERE status = 1 ORDER BY sort_order ASC;")
	CityData, err := ExecuteQuery(CityQry.String())
	if util.CheckErrorLog(r, err) {
		return nil, err
	}
	retBody["cityList"] = CityData

	// Recommended For You Data
	var PropertyQry bytes.Buffer
	PropertyQry.WriteString("SELECT CRP.hotel_id, CRP.sort_order, CHI.hotel_name,CHI.short_address, CCO.name AS country_name, CCI.name AS city_name, CFS.name AS state_name, ")
	PropertyQry.WriteString(" CPT.type AS property_type, IFNULL(CFL.locality,'') AS locality_name, CHI.hotel_star,")
	PropertyQry.WriteString("(SELECT image FROM cf_hotel_image WHERE hotel_id= CRP.hotel_id ORDER BY sortorder ASC LIMIT 1) AS image ")
	PropertyQry.WriteString(" FROM cms_recommended_property AS CRP ")
	PropertyQry.WriteString(" INNER JOIN cf_hotel_info AS CHI ON CHI.id = CRP.hotel_id ")
	PropertyQry.WriteString(" INNER JOIN cf_country AS CCO On CCO.id = CHI.country_id")
	PropertyQry.WriteString(" INNER JOIN cf_city AS CCI ON CCI.id = CHI.city_id")
	PropertyQry.WriteString(" INNER JOIN cf_states AS CFS ON CFS.id = CHI.state_id")
	PropertyQry.WriteString(" INNER JOIN cf_property_type AS CPT ON CPT.id = CHI.property_type_id ")
	PropertyQry.WriteString(" LEFT JOIN cf_locality AS CFL ON CFL.id = CHI.locality_id WHERE CHI.status=1 AND CRP.status=1 ORDER BY CRP.sort_order ASC")
	HotelData, err := ExecuteQuery(PropertyQry.String())
	if util.CheckErrorLog(r, err) {
		return nil, err
	}
	retBody["recommendedForYou"] = HotelData

	// CMS Data
	var CmsQry bytes.Buffer
	CmsQry.WriteString("SELECT page, slug FROM cf_cms;")
	CmsData, err := ExecuteQuery(CmsQry.String())
	if util.CheckErrorLog(r, err) {
		return nil, err
	}
	retBody["cmsList"] = CmsData

	return retBody, nil
}

// GetRatingQuestions - Returns rating questions
func GetRatingQuestions(r *http.Request) ([]map[string]interface{}, error) {
	util.LogIt(r, "Model - front - GetRatingQuestions")
	var Qry bytes.Buffer
	Qry.WriteString("SELECT id, question FROM cf_feedback_questions WHERE status=1")
	PropertyData, err := ExecuteQuery(Qry.String())
	if util.CheckErrorLog(r, err) {
		return nil, err
	}
	return PropertyData, nil
}

// GetCmsData - Returns Cms Data
func GetCmsData(r *http.Request, slug string) (map[string]interface{}, error) {
	util.LogIt(r, "Model - front - GetCmsData")

	var Qry bytes.Buffer
	Qry.WriteString("SELECT page, content FROM cf_cms WHERE slug = ?")
	CmsData, err := ExecuteRowQuery(Qry.String(), slug)
	if util.CheckErrorLog(r, err) {
		return nil, err
	}
	return CmsData, nil
}

// GetCmsListData - Returns Cms list Data
func GetCmsListData(r *http.Request) ([]map[string]interface{}, error) {
	util.LogIt(r, "Model - front - GetCmsListData")

	var Qry bytes.Buffer
	Qry.WriteString("SELECT page, slug, content FROM cf_cms")
	CmsListData, err := ExecuteQuery(Qry.String())
	if util.CheckErrorLog(r, err) {
		return nil, err
	}
	return CmsListData, nil
}
