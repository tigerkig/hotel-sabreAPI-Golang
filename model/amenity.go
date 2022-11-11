package model

import (
	"bytes"
	"net/http"
	"strings"
	"tp-api-common/data"
	"tp-api-common/util"

	"github.com/gorilla/context"
	gonanoid "github.com/matoous/go-nanoid"
)

// AddAmenity - Add Amenity
func AddAmenity(r *http.Request, reqMap data.Amenity) bool {

	util.LogIt(r, "model - V_Amenity - AddAmenity")
	nanoid, _ := gonanoid.Nanoid()

	var Qry bytes.Buffer
	Qry.WriteString("INSERT INTO cf_amenity(id,name,amenity_type_id,is_star_amenity,icon,created_at,created_by, amenity_of) VALUES (?,?,?,?,?,?,?,?)")
	err := ExecuteNonQuery(Qry.String(), nanoid, reqMap.Name, reqMap.Type, reqMap.IsStarAmenity, reqMap.Icon, util.GetIsoLocalDateTime(), context.Get(r, "UserId"), reqMap.AmenityOf)
	if util.CheckErrorLog(r, err) {
		return false
	}

	reqStruct := util.ToMap(reqMap)
	if reqMap.IsStarAmenity == 1 {
		reqStruct["Star Amenity"] = "Yes"
	} else {
		reqStruct["Star Amenity"] = "No"
	}

	TypeValue, _ := GetModuleFieldByID(r, "AMENITY_TYPE", reqMap.Type, "type")
	reqStruct["Amenity Type"] = TypeValue

	AddLog(r, "", "AMENITY", "Create", nanoid, GetLogsValueMap(r, reqStruct, true, "IsStarAmenity,Type"))

	return true
}

// UpdateAmenity - Update Amenity
func UpdateAmenity(r *http.Request, reqMap data.Amenity) bool {
	util.LogIt(r, "model - V_Amenity - UpdateAmenity")
	var Qry bytes.Buffer

	NameValue, _ := GetModuleFieldByID(r, "AMENITY", reqMap.ID, "name")

	Qry.WriteString("UPDATE cf_amenity SET name=?,amenity_type_id=?,is_star_amenity=?,icon=?, amenity_of=? WHERE id=?")
	err := ExecuteNonQuery(Qry.String(), reqMap.Name, reqMap.Type, reqMap.IsStarAmenity, reqMap.Icon, reqMap.AmenityOf, reqMap.ID)
	if util.CheckErrorLog(r, err) {
		return false
	}

	reqStruct := util.ToMap(reqMap)
	if reqMap.IsStarAmenity == 1 {
		reqStruct["Star Amenity"] = "Yes"
	} else {
		reqStruct["Star Amenity"] = "No"
	}

	TypeValue, _ := GetModuleFieldByID(r, "AMENITY_TYPE", reqMap.Type, "type")
	reqStruct["Amenity Type"] = TypeValue

	AddLog(r, NameValue.(string), "AMENITY", "Update", reqMap.ID, GetLogsValueMap(r, reqStruct, true, "IsStarAmenity,Type"))

	return true
}

// GetAmenity -  Get Amenity Detail By ID
func GetAmenity(r *http.Request, id string) (map[string]interface{}, error) {
	util.LogIt(r, "model - V_Amenity - GetAmenity")
	var Qry bytes.Buffer

	Qry.WriteString("SELECT id,name,amenity_type_id,is_star_amenity,icon,amenity_of  FROM cf_amenity WHERE id = ?")
	RetMap, err := ExecuteRowQuery(Qry.String(), id)
	if util.CheckErrorLog(r, err) {
		return nil, err
	}

	return RetMap, nil
}

// AmenityListing - Return Datatable Listing Of Amenity
func AmenityListing(r *http.Request, reqMap data.JQueryTableUI) (map[string]interface{}, error) {
	util.LogIt(r, "model - V_Amenity - AmenityListing")
	var Qry, QryCnt, QryFilter, GroupBy, FromQry bytes.Buffer
	var status = make(map[string]interface{})

	var testColArrs [20]string
	testColArrs[0] = "CFA.id"
	testColArrs[1] = "CFA.name"
	testColArrs[2] = "CFT.type"
	testColArrs[3] = "CFA.status"
	testColArrs[4] = "created_by"

	var testArrs []map[string]string
	testArrs = append(testArrs, map[string]string{
		"key":   "amenity",
		"value": "CFA.name",
	})
	testArrs = append(testArrs, map[string]string{
		"key":   "amenity_of",
		"value": "CFA.amenity_of",
	})
	testArrs = append(testArrs, map[string]string{
		"key":   "amenity_type",
		"value": "CFT.id",
	})
	testArrs = append(testArrs, map[string]string{
		"key":   "status",
		"value": "CFA.status",
	})
	testArrs = append(testArrs, map[string]string{
		"key":   "created_at",
		"value": "DATE(from_unixtime(CFA.created_at))",
	})

	QryCnt.WriteString(" COUNT(CFA.id) AS cnt ")
	QryFilter.WriteString(" COUNT(CFA.id) AS cnt ")

	Qry.WriteString(" CFA.id, CFT.type AS amenity_type, CFA.amenity_of, CFA.name AS amenity, ST.status, CONCAT(from_unixtime(CFA.created_at),' ',CU.username) AS created_by,ST.id AS status_id ")

	FromQry.WriteString(" FROM cf_amenity AS CFA ")
	FromQry.WriteString(" INNER JOIN cf_amenity_type AS CFT ON CFT.id = CFA.amenity_type_id ")
	FromQry.WriteString(" INNER JOIN cf_user AS CU ON CU.id = CFA.created_by ")
	FromQry.WriteString(" INNER JOIN status AS ST ON ST.id = CFA.status ")
	FromQry.WriteString(" WHERE ST.id <> 3 ")
	Data, err := JQueryTable(r, reqMap, Qry, FromQry, QryCnt, QryFilter, GroupBy, testArrs, testColArrs, status, "")
	return Data, err
}

// GetAmenityList -  Get Amenity Active List
func GetAmenityList(r *http.Request) ([]map[string]interface{}, error) {
	util.LogIt(r, "model - V_Amenity - GetAmenityList")
	var Qry bytes.Buffer

	Qry.WriteString("SELECT id,name,amenity_type_id,is_star_amenity,icon FROM cf_amenity WHERE status = 1")
	RetMap, err := ExecuteQuery(Qry.String())
	if util.CheckErrorLog(r, err) {
		return nil, err
	}

	return RetMap, nil
}

// AmenityTypeWiseAmenity - Return amenity type wise amenity data
func AmenityTypeWiseAmenity(r *http.Request, HotelID string) (map[string]interface{}, error) {
	util.LogIt(r, "model - V_Amenity - AmenityTypeWiseAmenity")
	// AmenityType, err := GetAmenityTypeList(r)
	AmenityType, err := GetAmenityTypeListV1(r, "1")
	if util.CheckErrorLog(r, err) {
		return nil, err
	}

	var Qry1 bytes.Buffer
	Qry1.WriteString(" SELECT amenity_id AS id, instruction FROM cf_hotel_amenities WHERE hotel_id = ?")
	hotelAmenity, err := ExecuteQuery(Qry1.String(), HotelID)
	if util.CheckErrorLog(r, err) {
		return nil, err
	}

	values := []string{}
	for _, val := range hotelAmenity {
		values = append(values, val["id"].(string))
	}

	result := "'" + strings.Join(values, "','") + "'"

	for i := 0; i < len(AmenityType); i++ {
		var Qry bytes.Buffer

		// Qry.WriteString(" SELECT id, name, true AS is_selected FROM cf_amenity WHERE amenity_type_id = ? AND status = 1 AND id in (" + result + ") ")
		// Qry.WriteString(" UNION SELECT id,name , false AS is_selected FROM cf_amenity WHERE amenity_type_id = ? AND status = 1 AND id not in (" + result + ") ")

		Qry.WriteString(" SELECT id, name, true AS is_selected FROM cf_amenity WHERE amenity_type_id = ? AND status = 1 AND amenity_of = 1 AND id in (" + result + ") ")
		Qry.WriteString(" UNION SELECT id, name, false AS is_selected FROM cf_amenity WHERE amenity_type_id = ? AND status = 1 AND amenity_of = 1 AND id not in (" + result + ") ")

		AmenityTypeData, err := ExecuteQuery(Qry.String(), AmenityType[i]["id"], AmenityType[i]["id"])
		if util.CheckErrorLog(r, err) {
			return nil, err
		}

		for k := 0; k < len(AmenityTypeData); k++ {
			if len(hotelAmenity) == 0 {
				AmenityTypeData[k]["instruction"] = ""
			} else {
				for j := 0; j < len(hotelAmenity); j++ {
					if AmenityTypeData[k]["id"].(string) == hotelAmenity[j]["id"].(string) {
						AmenityTypeData[k]["instruction"] = hotelAmenity[j]["instruction"]
						break
					} else {
						AmenityTypeData[k]["instruction"] = ""
					}
				}
			}
		}

		if len(AmenityTypeData) > 0 {
			var dataCount int64
			for _, val := range AmenityTypeData {
				dataCount = dataCount + val["is_selected"].(int64)
			}

			if dataCount >= 1 {
				AmenityType[i]["is_selected"] = 1
			} else {
				AmenityType[i]["is_selected"] = 0
			}

			AmenityType[i]["amenity"] = AmenityTypeData
		} else {
			AmenityType[i]["is_selected"] = 0
			AmenityType[i]["amenity"] = []string{}
		}
	}

	stuff := make(map[string]interface{})
	if len(AmenityType) > 0 {
		stuff["data"] = AmenityType
	} else {
		stuff["data"] = []string{}
	}

	return stuff, nil
}

// AmenityTypeWiseAmenityAdmin - Return amenity type wise amenity data only selected view for admin panel
func AmenityTypeWiseAmenityAdmin(r *http.Request, HotelID string) (map[string]interface{}, error) {
	util.LogIt(r, "model - V_Amenity - AmenityTypeWiseAmenityAdmin")
	AmenityType, err := GetAmenityTypeList(r)
	if util.CheckErrorLog(r, err) {
		return nil, err
	}

	var Qry1 bytes.Buffer
	Qry1.WriteString(" SELECT amenity_id AS id,instruction FROM cf_hotel_amenities WHERE hotel_id = ?")
	hotelAmenity, err := ExecuteQuery(Qry1.String(), HotelID)
	if util.CheckErrorLog(r, err) {
		return nil, err
	}

	values := []string{}
	for _, val := range hotelAmenity {
		values = append(values, val["id"].(string))
	}

	result := "'" + strings.Join(values, "','") + "'"
	var NewMapArr = []map[string]interface{}{}

	for i := 0; i < len(AmenityType); i++ {
		var Qry bytes.Buffer
		Qry.WriteString(" SELECT id,name , true AS is_selected FROM cf_amenity WHERE amenity_type_id = ? AND status = 1 AND id in (" + result + ") ")
		AmenityTypeData, err := ExecuteQuery(Qry.String(), AmenityType[i]["id"])
		if util.CheckErrorLog(r, err) {
			return nil, err
		}

		for k := 0; k < len(AmenityTypeData); k++ {
			for j := 0; j < len(hotelAmenity); j++ {
				if AmenityTypeData[k]["id"] == hotelAmenity[j]["id"] {
					AmenityTypeData[k]["instruction"] = hotelAmenity[j]["instruction"]
				} else {
					AmenityTypeData[k]["instruction"] = ""
				}
			}
		}

		if len(AmenityTypeData) > 0 {
			AmenityType[i]["amenity"] = AmenityTypeData
			NewMapArr = append(NewMapArr, AmenityType[i])
		}
	}

	stuff := make(map[string]interface{})
	if len(NewMapArr) > 0 {
		stuff["data"] = NewMapArr
	}

	return stuff, nil
}

// GetAmenityListV1 -  Get Amenity Active List
func GetAmenityListV1(r *http.Request, CatgID string) ([]map[string]interface{}, error) {
	util.LogIt(r, "model - V_Amenity - GetAmenityListV1")
	var Qry bytes.Buffer

	Qry.WriteString("SELECT id, name, amenity_type_id, is_star_amenity, icon FROM cf_amenity WHERE status = 1 AND amenity_of = ?")
	RetMap, err := ExecuteQuery(Qry.String(), CatgID)
	if util.CheckErrorLog(r, err) {
		return nil, err
	}

	return RetMap, nil
}

// AmenityTypeWiseAmenityForRoom - Return amenity type wise amenity data For Room
func AmenityTypeWiseAmenityForRoom(r *http.Request, HotelID string, RoomID string) (map[string]interface{}, error) {

	util.LogIt(r, "model - V_Amenity - AmenityTypeWiseAmenityForRoom")

	AmenityType, err := GetAmenityTypeListV1(r, "2")
	if util.CheckErrorLog(r, err) {
		return nil, err
	}

	var Qry1 bytes.Buffer
	Qry1.WriteString(" SELECT amenity_id AS id, extra_detail as instruction FROM cf_room_amenity WHERE hotel_id = ? and room_type_id = ?")
	hotelRoomAmenity, err := ExecuteQuery(Qry1.String(), HotelID, RoomID)
	if util.CheckErrorLog(r, err) {
		return nil, err

	}

	values := []string{}
	for _, val := range hotelRoomAmenity {
		values = append(values, val["id"].(string))
	}

	result := "'" + strings.Join(values, "','") + "'"

	for i := 0; i < len(AmenityType); i++ {
		var Qry bytes.Buffer

		Qry.WriteString(" SELECT id, name, true AS is_selected FROM cf_amenity WHERE amenity_type_id = ? AND status = 1 AND amenity_of = 2 AND id in (" + result + ") ")
		Qry.WriteString(" UNION SELECT id, name, false AS is_selected FROM cf_amenity WHERE amenity_type_id = ? AND status = 1 AND amenity_of = 2 AND id not in (" + result + ") ")

		AmenityTypeData, err := ExecuteQuery(Qry.String(), AmenityType[i]["id"], AmenityType[i]["id"])
		if util.CheckErrorLog(r, err) {
			return nil, err
		}

		for k := 0; k < len(AmenityTypeData); k++ {

			if len(hotelRoomAmenity) == 0 {
				AmenityTypeData[k]["instruction"] = ""
			} else {
				for j := 0; j < len(hotelRoomAmenity); j++ {
					if AmenityTypeData[k]["id"].(string) == hotelRoomAmenity[j]["id"].(string) {
						AmenityTypeData[k]["instruction"] = hotelRoomAmenity[j]["instruction"]
						break
					} else {
						AmenityTypeData[k]["instruction"] = ""
					}
				}
			}
		}

		if len(AmenityTypeData) > 0 {
			var dataCount int64
			for _, val := range AmenityTypeData {
				dataCount = dataCount + val["is_selected"].(int64)
			}

			if dataCount >= 1 {
				AmenityType[i]["is_selected"] = 1
			} else {
				AmenityType[i]["is_selected"] = 0
			}

			AmenityType[i]["amenity"] = AmenityTypeData
		} else {
			AmenityType[i]["is_selected"] = 0
			AmenityType[i]["amenity"] = []string{}
		}
	}

	stuff := make(map[string]interface{})
	if len(AmenityType) > 0 {
		stuff["data"] = AmenityType
	} else {
		stuff["data"] = []string{}
	}

	return stuff, nil
}

// GetStarAmenity -  Get Star Amenity
func GetStarAmenity(r *http.Request) ([]map[string]interface{}, error) {
	util.LogIt(r, "model - V_Amenity - GetStarAmenity")
	var Qry bytes.Buffer

	Qry.WriteString("SELECT id, name, icon FROM cf_amenity WHERE status = 1 AND amenity_of = 1 AND is_star_amenity=1")
	RetMap, err := ExecuteQuery(Qry.String())
	if util.CheckErrorLog(r, err) {
		return nil, err
	}

	return RetMap, nil
}
