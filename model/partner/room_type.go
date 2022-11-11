package partner

import (
	"bytes"
	"net/http"
	"tp-api-common/data"
	"tp-api-common/util"
	"tp-system/config"
	"tp-system/model"

	"github.com/gorilla/context"
	gonanoid "github.com/matoous/go-nanoid"
)

// AddRoomType - Add Room For Hotel
func AddRoomType(r *http.Request, reqMap data.RoomType) (map[string]interface{}, error) {

	util.LogIt(r, "model - V_Room_Type - AddRoomType")

	var Qry bytes.Buffer
	nanoid, _ := gonanoid.Nanoid()

	Qry.WriteString("INSERT INTO cf_room_type(id, room_type_name, bed_type_id, room_size, description, room_view_id, is_extra_bed, extra_bed_type_id, max_occupancy, inventory, sort_order, created_at, created_by, hotel_id) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?)")
	// err := model.ExecuteNonQuery(Qry.String(), nanoid, reqMap.Name, reqMap.BedType, reqMap.RoomSize, reqMap.Description, reqMap.RoomView, reqMap.IsExtraBed, reqMap.ExtraBed, reqMap.MaxOccupancy, reqMap.Inventory, reqMap.SortOrder, util.GetIsoLocalDateTime(), context.Get(r, "UserId"), context.Get(r, "HotelId"))
	err := model.ExecuteNonQuery(Qry.String(), nanoid, reqMap.Name, reqMap.BedType, reqMap.RoomSize, reqMap.Description, reqMap.RoomView, reqMap.IsExtraBed, reqMap.ExtraBed, reqMap.MaxOccupancy, reqMap.Inventory, reqMap.SortOrder, util.GetIsoLocalDateTime(), context.Get(r, "UserId"), reqMap.HotelID)
	if util.CheckErrorLog(r, err) {
		return nil, err
	}

	// model.UpdateAllRoomType(context.Get(r, "HotelId").(string), nanoid) // 2020-06-24 - HK - Room Add Sync With Mongo Added - Partner Panel

	model.AddLog(r, "", "ROOM_TYPE", "Create", nanoid, model.GetLogsValueMap(r, util.ToMap(reqMap), false, ""))

	var retMap = make(map[string]interface{})
	retMap["id"] = nanoid

	model.CacheChn <- model.CacheObj{
		Type: "roomType",
		// ID:         context.Get(r, "HotelId").(string),
		ID:         reqMap.HotelID,
		Additional: nanoid,
	}

	// invFilled := model.CheckAndFillInvdataOnActiveStatus(r, context.Get(r, "HotelId").(string), nanoid)
	invFilled := model.CheckAndFillInvdataOnActiveStatus(r, reqMap.HotelID, nanoid)
	if !invFilled {
		util.LogIt(r, "Issue on sync data on add room type - "+nanoid)
	}

	return retMap, nil
}

// UpdateRoomType - Updates Room Info
func UpdateRoomType(r *http.Request, reqMap data.RoomType) bool {
	util.LogIt(r, "model - V_Room_Type - UpdateRoomType")

	var Qry bytes.Buffer

	Qry.WriteString("UPDATE cf_room_type SET room_type_name=?, bed_type_id=?, room_size=?, description=?, room_view_id=?, is_extra_bed=?, extra_bed_type_id=?, max_occupancy=?, inventory=?, sort_order=?, created_at=?, created_by=? WHERE id=? AND hotel_id=?")
	// err := model.ExecuteNonQuery(Qry.String(), reqMap.Name, reqMap.BedType, reqMap.RoomSize, reqMap.Description, reqMap.RoomView, reqMap.IsExtraBed, reqMap.ExtraBed, reqMap.MaxOccupancy, reqMap.Inventory, reqMap.SortOrder, util.GetIsoLocalDateTime(), context.Get(r, "UserId"), reqMap.ID, context.Get(r, "HotelId"))
	err := model.ExecuteNonQuery(Qry.String(), reqMap.Name, reqMap.BedType, reqMap.RoomSize, reqMap.Description, reqMap.RoomView, reqMap.IsExtraBed, reqMap.ExtraBed, reqMap.MaxOccupancy, reqMap.Inventory, reqMap.SortOrder, util.GetIsoLocalDateTime(), context.Get(r, "UserId"), reqMap.ID, reqMap.HotelID)
	if util.CheckErrorLog(r, err) {
		return false
	}

	// model.UpdateAllRoomType(context.Get(r, "HotelId").(string), reqMap.ID) // 2020-06-24 - HK - Room Add Sync With Mongo Added - Partner Panel

	model.AddLog(r, "", "ROOM_TYPE", "UPDATE", reqMap.ID, model.GetLogsValueMap(r, util.ToMap(reqMap), false, ""))

	model.CacheChn <- model.CacheObj{
		Type: "roomType",
		// ID:         context.Get(r, "HotelId").(string),
		ID:         reqMap.HotelID,
		Additional: reqMap.ID,
	}

	// invFilled := model.CheckAndFillInvdataOnActiveStatus(r, context.Get(r, "HotelId").(string), reqMap.ID)
	invFilled := model.CheckAndFillInvdataOnActiveStatus(r, reqMap.HotelID, reqMap.ID)
	if !invFilled {
		util.LogIt(r, "Issue on sync data on update room type - "+reqMap.ID)
	}

	return true
}

// RoomTypeListing - Return Datatable Listing Of Cancel Policy
func RoomTypeListing(r *http.Request, reqMap data.JQueryTableUI) (map[string]interface{}, error) {

	util.LogIt(r, "model - V_Room_Type - RoomTypeListing")

	var Qry, QryCnt, QryFilter, GroupBy, FromQry bytes.Buffer
	var status = make(map[string]interface{})

	var testColArrs [20]string
	testColArrs[0] = "CRT.id"
	testColArrs[1] = "CRT.room_type_name"
	testColArrs[2] = "ST.status"
	testColArrs[3] = "CBT.bed_type"
	testColArrs[4] = "CRV.room_view_name"
	testColArrs[5] = "CRT.room_size"
	testColArrs[6] = "CRT.max_occupancy"
	testColArrs[7] = "CRT.sort_order"
	testColArrs[8] = "created_by"

	var testArrs []map[string]string
	testArrs = append(testArrs, map[string]string{
		"key":   "room_type_name",
		"value": "CRT.room_type_name",
	})
	testArrs = append(testArrs, map[string]string{
		"key":   "bed_type",
		"value": "CBT.bed_type",
	})
	testArrs = append(testArrs, map[string]string{
		"key":   "room_view",
		"value": "CRV.id",
	})
	testArrs = append(testArrs, map[string]string{
		"key":   "status",
		"value": "CRT.status",
	})
	testArrs = append(testArrs, map[string]string{
		"key":   "created_at",
		"value": "DATE(from_unixtime(CRT.created_at))",
	})

	QryCnt.WriteString(" COUNT(CRT.id) AS cnt ")
	QryFilter.WriteString(" COUNT(CRT.id) AS cnt ")

	Qry.WriteString(" CRT.id, CRT.room_type_name, CRT.room_size, CRT.description, ST.status, CONCAT(from_unixtime(CRT.created_at),' ',CHC.username) AS created_by, ST.id AS status_id, CRT.sort_order, ")
	Qry.WriteString(" CRT.bed_type_id, CBT.bed_type, ")
	Qry.WriteString(" CRT.room_view_id, CRV.room_view_name, ")
	Qry.WriteString(" CRT.is_extra_bed, CRT.extra_bed_type_id, ")
	Qry.WriteString(" CASE WHEN CRT.is_extra_bed = 1 AND CRT.extra_bed_type_id != '' THEN CEBT.extra_bed_name ELSE '' END AS extra_bed_name, ")
	Qry.WriteString(" CRT.max_occupancy, CRT.inventory ")

	FromQry.WriteString(" FROM cf_room_type AS CRT ")

	FromQry.WriteString(" INNER JOIN cf_bed_type AS CBT ON CBT.id = CRT.bed_type_id ")
	FromQry.WriteString(" INNER JOIN cf_room_view AS CRV ON CRV.id = CRT.room_view_id ")
	FromQry.WriteString(" LEFT JOIN cf_extra_bed_type AS CEBT ON CEBT.id = CRT.extra_bed_type_id ")

	FromQry.WriteString(" INNER JOIN cf_hotel_client AS CHC ON CHC.id = CRT.created_by ")
	FromQry.WriteString(" INNER JOIN status AS ST ON ST.id = CRT.status ")
	FromQry.WriteString(" WHERE ST.id <> 3 AND CRT.hotel_id = '")
	// FromQry.WriteString(context.Get(r, "HotelId").(string))
	FromQry.WriteString(reqMap.HotelID)
	FromQry.WriteString("'")
	Data, err := model.JQueryTable(r, reqMap, Qry, FromQry, QryCnt, QryFilter, GroupBy, testArrs, testColArrs, status, "")
	return Data, err
}

// GetRoomType - Get Room Type Info
func GetRoomType(r *http.Request, roomID string) (map[string]interface{}, error) {

	util.LogIt(r, "model - V_Room_Type - GetRoomType")

	var Qry bytes.Buffer

	Qry.WriteString("SELECT id, room_type_name, bed_type_id, room_size, description, room_view_id, is_extra_bed, extra_bed_type_id, max_occupancy, inventory, sort_order, created_at, created_by, hotel_id FROM cf_room_type WHERE id = ? AND hotel_id = ?")
	RetMap, err := model.ExecuteRowQuery(Qry.String(), roomID, context.Get(r, "HotelId"))
	if util.CheckErrorLog(r, err) {
		return nil, err
	}

	return RetMap, nil
}

// UpdateRoomAmenity - Update Amenity Of Room
func UpdateRoomAmenity(r *http.Request, reqMap data.RoomAmenity) bool {

	util.LogIt(r, "model - V_Partner_Room - UpdateRoomAmenity")

	var DelQry bytes.Buffer
	var AmenityArr = reqMap.Amenity

	DelQry.WriteString("DELETE FROM cf_room_amenity WHERE hotel_id = ? and room_type_id = ?")
	err := model.ExecuteNonQuery(DelQry.String(), reqMap.HotelID, reqMap.RoomID)
	if util.CheckErrorLog(r, err) {
		return false
	}

	for _, val := range AmenityArr {
		var Qry bytes.Buffer
		nanoid, _ := gonanoid.Nanoid()
		Qry.WriteString("INSERT INTO cf_room_amenity SET id=?, amenity_id=?, extra_detail = ?, room_type_id = ?, hotel_id = ?")
		err = model.ExecuteNonQuery(Qry.String(), nanoid, val.ID, val.Description, reqMap.RoomID, reqMap.HotelID)
		if util.CheckErrorLog(r, err) {
			return false
		}
	}

	// model.UpdateRoomAmenity(reqMap.HotelID, reqMap.RoomID) // 2020-06-24 - HK - Room Amenity Sync With Mongo Added - Partner Panel

	model.AddLog(r, "", "ROOM_TYPE", "Update Room Amenities", reqMap.RoomID, map[string]interface{}{})

	model.CacheChn <- model.CacheObj{
		Type:       "roomTypeAmenity",
		ID:         reqMap.HotelID,
		Additional: reqMap.RoomID,
	}

	return true
}

// GetHotelRoomImageCount - Get Hotel Room Image Count
func GetHotelRoomImageCount(r *http.Request, RoomID string) (int64, error) {

	util.LogIt(r, "model - V_Partner_Hotel - GetHotelRoomImageCount")
	var Qry bytes.Buffer

	Qry.WriteString("SELECT COUNT(id) AS id FROM cf_room_image WHERE room_type_id = ? AND hotel_id = ?")
	Data, err := model.ExecuteRowQuery(Qry.String(), RoomID, context.Get(r, "HotelId"))
	if util.CheckErrorLog(r, err) {
		return 0, err
	}

	return Data["id"].(int64), nil
}

// UploadHotelRoomImage - Upload Hotel Room Image Max 5 At One Time
func UploadHotelRoomImage(r *http.Request, reqMap map[string]interface{}) bool {
	util.LogIt(r, "model - V_Partner_RoomType - UploadHotelRoomImage")

	var SQLGet bytes.Buffer
	SQLGet.WriteString(" SELECT CASE WHEN MAX(sort_order) IS NULL THEN 0 ELSE MAX(sort_order) END AS sortorder FROM cf_room_image WHERE room_type_id = ? AND hotel_id = ? ")
	SortOrder, err := model.ExecuteRowQuery(SQLGet.String(), reqMap["RoomID"], reqMap["HotelID"])
	if util.CheckErrorLog(r, err) {
		return false
	}

	var latestSortOrder = int(SortOrder["sortorder"].(int64))

	for i, val := range reqMap["Image"].([]string) {
		var Qry bytes.Buffer
		nanoid, _ := gonanoid.Nanoid()
		Qry.WriteString("INSERT INTO cf_room_image SET id = ?, room_type_id = ?, image = ?, sort_order = ?, hotel_id = ?")
		err := model.ExecuteNonQuery(Qry.String(), nanoid, reqMap["RoomID"], val, latestSortOrder+i+1, reqMap["HotelID"])
		if util.CheckErrorLog(r, err) {
			return false
		}
	}

	// model.UpdateRoomImage(reqMap["HotelID"].(string), reqMap["RoomID"].(string)) // 2020-06-24 - HK - Room Photo Sync With Mongo Added - Partner Panel

	BeforeUpdate, _ := model.GetModuleFieldByID(r, "ROOM_TYPE", reqMap["RoomID"].(string), "room_type_name")

	model.AddLog(r, "", "ROOM_TYPE", "Upload Hotel Room Images", reqMap["RoomID"].(string), map[string]interface{}{"Room Name": BeforeUpdate, "Total Image": len(reqMap["Image"].([]string))})

	model.CacheChn <- model.CacheObj{
		Type:       "roomTypeImages",
		ID:         reqMap["HotelID"].(string),
		Additional: reqMap["RoomID"].(string),
	}

	return true
}

// GetHotelRoomImageName - Get Hotel Room Image Name For Delete Function
func GetHotelRoomImageName(r *http.Request, reqMap data.RoomImageName) (string, error) {
	util.LogIt(r, "model - V_Partner_RoomType - GetHotelRoomImageName")

	var Qry bytes.Buffer

	Qry.WriteString("SELECT image FROM cf_room_image WHERE id = ? AND room_type_id = ? AND hotel_id = ? ")
	// Data, err := model.ExecuteRowQuery(Qry.String(), reqMap.Image, context.Get(r, "HotelId"))
	Data, err := model.ExecuteRowQuery(Qry.String(), reqMap.Image, reqMap.RoomID, reqMap.HotelID)
	if util.CheckErrorLog(r, err) {
		return "", err
	}

	return Data["image"].(string), nil
}

// DeleteHotelRoomImage - Delete Hotel Room Image
func DeleteHotelRoomImage(r *http.Request, reqMap data.RoomImageName) bool {
	util.LogIt(r, "model - V_Partner_RoomType - DeleteHotelRoomImage")
	var Qry bytes.Buffer

	/*
		var Qry1 bytes.Buffer
		Qry1.WriteString("SELECT room_type_id FROM cf_room_image WHERE id = ? AND hotel_id = ? ")
		Data, err := model.ExecuteRowQuery(Qry1.String(), reqMap.Image, context.Get(r, "HotelId"))
		if util.CheckErrorLog(r, err) {
			return false
		}

		BeforeUpdate, _ := model.GetModuleFieldByID(r, "ROOM_TYPE", Data["room_type_id"].(string), "room_type_name")
	*/

	BeforeUpdate, _ := model.GetModuleFieldByID(r, "ROOM_TYPE", reqMap.RoomID, "room_type_name")
	Qry.WriteString("DELETE FROM cf_room_image WHERE id=? AND hotel_id = ?")
	// err := model.ExecuteNonQuery(Qry.String(), reqMap.Image, context.Get(r, "HotelId"))
	err := model.ExecuteNonQuery(Qry.String(), reqMap.Image, reqMap.HotelID)
	if util.CheckErrorLog(r, err) {
		return false
	}

	// model.UpdateRoomImage(context.Get(r, "HotelId").(string), Data["room_type_id"].(string)) // 2020-06-24 - HK - Room Photo Sync With Mongo Added - Partner Panel

	// model.AddLog(r, "", "ROOM_TYPE", "Delete Hotel Room Image", Data["room_type_id"].(string), map[string]interface{}{"Room Name": BeforeUpdate.(string), "Image": reqMap.Image})
	model.AddLog(r, "", "ROOM_TYPE", "Delete Hotel Room Image", reqMap.RoomID, map[string]interface{}{"Room Name": BeforeUpdate.(string), "Image": reqMap.Image})

	model.CacheChn <- model.CacheObj{
		Type: "roomTypeImages",
		// ID:   context.Get(r, "HotelId").(string),
		// Additional: Data["room_type_id"].(string),
		ID:         reqMap.HotelID,
		Additional: reqMap.RoomID,
	}

	return true
}

// GetRoomImageList - Get Room Image List
func GetRoomImageList(r *http.Request, roomID string, HotelID string) ([]map[string]interface{}, error) {
	util.LogIt(r, "model - V_Room_Type - GetRoomImageList")

	var Qry bytes.Buffer

	Qry.WriteString("SELECT id, CONCAT('" + config.Env.AwsBucketURL + "room/" + "',image) AS image, sort_order FROM cf_room_image WHERE room_type_id = ? AND hotel_id = ? ORDER BY sort_order")
	// RetMap, err := model.ExecuteQuery(Qry.String(), roomID, context.Get(r, "HotelId"))
	RetMap, err := model.ExecuteQuery(Qry.String(), roomID, HotelID)
	if util.CheckErrorLog(r, err) {
		return nil, err
	}

	return RetMap, nil
}

// GetRoomTypeList - Get Room Type List For Other Module
func GetRoomTypeList(r *http.Request, HotelID string) ([]map[string]interface{}, error) {

	util.LogIt(r, "model - V_Room_Type - GetRoomTypeList")

	var Qry bytes.Buffer

	Qry.WriteString("SELECT id, room_type_name FROM cf_room_type WHERE status = 1 AND hotel_id = ?")
	RetMap, err := model.ExecuteQuery(Qry.String(), HotelID)
	if util.CheckErrorLog(r, err) {
		return nil, err
	}

	return RetMap, nil
}

// GetRoomTypeWithHotel - Get Room Type Info Of Passed Property - 2021-04-27 - HK
func GetRoomTypeWithHotel(r *http.Request, roomID string, HotelID string) (map[string]interface{}, error) {
	util.LogIt(r, "model - V_Room_Type - GetRoomTypeWithHotel")

	var Qry bytes.Buffer
	Qry.WriteString("SELECT id, room_type_name, bed_type_id, room_size, description, room_view_id, is_extra_bed, extra_bed_type_id, max_occupancy, inventory, sort_order, created_at, created_by, hotel_id FROM cf_room_type WHERE id = ? AND hotel_id = ?")
	RetMap, err := model.ExecuteRowQuery(Qry.String(), roomID, HotelID)
	if util.CheckErrorLog(r, err) {
		return nil, err
	}

	return RetMap, nil
}

// SortRoomImage - Change Sort Order Of Room Image - 2021-04-27 - HK
func SortRoomImage(r *http.Request, reqMap map[string]interface{}) (int, error) {
	util.LogIt(r, "Model - V_Room_Type - SortRoomImage")

	var SQLQry, SiQry, SQLUpdate, SQLNewUpdate bytes.Buffer

	SiQry.WriteString(" SELECT sort_order FROM cf_room_image WHERE id = ? AND hotel_id = ? AND room_type_id = ?")
	SortOrderOfGivenID, err := model.ExecuteRowQuery(SiQry.String(), reqMap["id"], reqMap["hotel_id"], reqMap["room_id"])
	if util.CheckErrorLog(r, err) {
		return 0, err
	}

	SQLQry.WriteString(" SELECT id FROM cf_room_image WHERE sort_order = ? AND hotel_id = ? AND room_type_id = ?")
	IDOfChangeSortOrder, err := model.ExecuteRowQuery(SQLQry.String(), reqMap["sortorder"], reqMap["hotel_id"], reqMap["room_id"])
	if util.CheckErrorLog(r, err) {
		return 0, err
	}

	SQLUpdate.WriteString("UPDATE cf_room_image SET sort_order = ? WHERE id = ? AND hotel_id = ? AND room_type_id = ?")
	err = model.ExecuteNonQuery(SQLUpdate.String(), reqMap["sortorder"], reqMap["id"], reqMap["hotel_id"], reqMap["room_id"])
	if util.CheckErrorLog(r, err) {
		return 0, err
	}

	SQLNewUpdate.WriteString("UPDATE cf_room_image SET sort_order = ? WHERE id = ? AND hotel_id = ? AND room_type_id = ?;")
	err = model.ExecuteNonQuery(SQLNewUpdate.String(), SortOrderOfGivenID["sort_order"], IDOfChangeSortOrder["id"], reqMap["hotel_id"], reqMap["room_id"])
	if util.CheckErrorLog(r, err) {
		return 0, err
	}

	// AddLog(r, "", LogModule, "Update Sort Order Of Image", mainID.(string), map[string]interface{}{"Object Name": DepColumnValue["log_column"].(string)})
	// AddLog(r, "", "HOTEL", "Update Sort Order Of Image", context.Get(r, "HotelId").(string), map[string]interface{}{})

	return 1, nil
}

// GetHotelRoomAmenityCount - Get Hotel Room Amenity Count - 2021-05-14 - HK
func GetHotelRoomAmenityCount(r *http.Request, RoomID string) (int64, error) {
	util.LogIt(r, "model - V_Partner_Hotel - GetHotelRoomAmenityCount")
	var Qry bytes.Buffer

	Qry.WriteString("SELECT COUNT(id) AS id FROM cf_room_amenity WHERE room_type_id = ? AND hotel_id = ?")
	Data, err := model.ExecuteRowQuery(Qry.String(), RoomID, context.Get(r, "HotelId"))
	if util.CheckErrorLog(r, err) {
		return 0, err
	}

	return Data["id"].(int64), nil
}
