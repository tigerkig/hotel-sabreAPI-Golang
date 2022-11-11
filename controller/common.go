package controller

import (
	"bytes"
	b64 "encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"log"
	"math/rand"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
	"tp-api-common/data"
	"tp-api-common/util"
	"tp-system/config"
	"tp-system/model"
	"tp-system/model/partner"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awsutil"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/disintegration/imaging"
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
)

// GetTimeZoneList - Returns Time Zone list
func GetTimeZoneList(w http.ResponseWriter, r *http.Request) {
	util.LogIt(r, "common - common - GetTimeZoneList")
	defer util.CommonDeferred(w, r, "common", "common", "GetTimeZoneList")
	util.RespondData(r, w, util.StaticArrTimeZone, 200)
}

// GetDateFormatList - Returns date format list
func GetDateFormatList(w http.ResponseWriter, r *http.Request) {
	util.LogIt(r, "common - common - GetDateFormatList")
	defer util.CommonDeferred(w, r, "common", "common", "GetDateFormatList")
	util.RespondData(r, w, util.StaticArrDateFormat, 200)
}

// GetTimeFormatList - Returns time format list
func GetTimeFormatList(w http.ResponseWriter, r *http.Request) {
	util.LogIt(r, "common - common - GetTimeFormatList")
	defer util.CommonDeferred(w, r, "common", "common", "GetTimeFormatList")
	util.RespondData(r, w, util.StaticArrTimeFormat, 200)
}

// StatusList - Returns List of status
func StatusList(w http.ResponseWriter, r *http.Request) {
	util.LogIt(r, "common - common - StatusList")
	defer util.CommonDeferred(w, r, "common", "common", "StatusList")
	util.RespondData(r, w, model.StatusList(r), 200)
}

// IconList - Returns Icon List For Amenity
func IconList(w http.ResponseWriter, r *http.Request) {
	util.LogIt(r, "common - common - IconList")
	defer util.CommonDeferred(w, r, "common", "common", "IconList")
	util.RespondData(r, w, util.IconList, 200)
}

// BankList - Returns Bank List
func BankList(w http.ResponseWriter, r *http.Request) {
	util.LogIt(r, "common - common - BankList")
	defer util.CommonDeferred(w, r, "common", "common", "BankList")
	util.RespondData(r, w, util.BankList, 200)
}

// GetResponseCodeMsg - It Will Get Response Code And Send Accoding Message
func GetResponseCodeMsg(w http.ResponseWriter, r *http.Request) {
	util.LogIt(r, "common - common - GetResponseCodeMsg")
	defer util.CommonDeferred(w, r, "common", "common", "GetResponseCodeMsg")
	var reqMap = make(map[string]interface{})
	code := r.URL.Query().Get("code")
	if code == "" {
		reqMap["message"] = util.RespCode
		util.RespondData(r, w, reqMap, 200)
		return
	} else {
		if util.RespCode[code] == "" {
			util.Respond(r, w, nil, 400, "10009")
			return
		}
		reqMap["message"] = util.RespCode[code]
		util.RespondData(r, w, reqMap, 200)
	}
}

// UpdateMultipleStatusModuleWise - Update Bulk Status
func UpdateMultipleStatusModuleWise(w http.ResponseWriter, r *http.Request) {
	util.LogIt(r, "common - common - UpdateMultipleStatusModuleWise")
	defer util.CommonDeferred(w, r, "common", "common", "UpdateMultipleStatusModuleWise")
	var reqMap data.StatusChange

	err := json.NewDecoder(r.Body).Decode(&reqMap)
	if err != nil {
		util.RespondBadRequest(r, w)
		return
	}

	for _, v := range reqMap.Status {
		ValidateString := ValidateNotNullStructString(v.ID)
		ValidateFloat := ValidateNotNullStructFloat(v.Status)
		if ValidateString == 0 || ValidateFloat == 0 {
			util.RespondBadRequest(r, w)
			return
		}

		flag, httpStatus, err := model.UpdateStatusModuleWise(r, reqMap.Module, int(v.Status), v.ID)
		if flag == 0 || err != nil {
			// util.RespondWithError(r, w, "500")
			util.Respond(r, w, nil, httpStatus, "100014")
			return
		}
	}
	util.Respond(r, w, nil, 204, "")
}

// DeleteMultipleStatusModuleWise - Delete Bulk Data
func DeleteMultipleStatusModuleWise(w http.ResponseWriter, r *http.Request) {
	util.LogIt(r, "common - common - DeleteMultipleStatusModuleWise")
	defer util.CommonDeferred(w, r, "common", "common", "DeleteMultipleStatusModuleWise")
	var reqMap data.StatusChange

	err := json.NewDecoder(r.Body).Decode(&reqMap)
	if err != nil {
		util.RespondBadRequest(r, w)
		return
	}

	for _, v := range reqMap.Status {
		ValidateString := ValidateNotNullStructString(v.ID)
		if ValidateString == 0 {
			util.RespondBadRequest(r, w)
			return
		}

		flag, httpStatus, err := model.UpdateStatusModuleWise(r, reqMap.Module, 3, v.ID)
		if flag == 0 || err != nil {
			// util.RespondWithError(r, w, "500")
			util.Respond(r, w, nil, httpStatus, "100014")
			return
		}
	}
	util.Respond(r, w, nil, 204, "")
}

// UpdateStatusModuleWise - Update Status
func UpdateStatusModuleWise(w http.ResponseWriter, r *http.Request) {
	util.LogIt(r, "common - common - UpdateStatusModuleWise")
	defer util.CommonDeferred(w, r, "common", "common", "UpdateStatusModuleWise")
	var reqMap data.StatusForSingle
	err := json.NewDecoder(r.Body).Decode(&reqMap)
	if err != nil {
		util.RespondBadRequest(r, w)
		return
	}
	if reqMap.HotelID != "" {
		context.Set(r, "HotelId", reqMap.HotelID)
	}

	vars := mux.Vars(r)
	id := vars["id"]
	reqMap.ID = id

	apiModule := vars["module"]
	if _, ok := util.StaticApiModules[apiModule]; !ok {
		util.Respond(r, w, nil, 406, "100011")
		return
	}

	var status = int(reqMap.Status)

	// 2020-05-14 - HK - START
	// Purpose : Dependent Record Checking Added Before Deativating Record
	if _, ok := util.DependentTables[apiModule]; ok {
		if status == 2 || status == 3 {
			// cnt, err := model.CheckDependentRecord(r, id, apiModule)
			var cnt int64
			var err error

			// If condition added here for special check for property_tag module because for all dependent
			// table moules, all have status column exists in their dependent table. But in the case of
			// property_tag we do have a slight different mechanism, where we remove tags directly upon
			// update operation and hence status column was not taken into cf_hotel_tag table. - HK - 2021-04-19
			if apiModule == "property_tag" || apiModule == "amenity" {
				cnt, err = model.CheckDependentRecordWithoutStatus(r, id, apiModule)
			} else {
				cnt, err = model.CheckDependentRecord(r, id, apiModule)
			}

			if err != nil || cnt > 0 {
				util.LogIt(r, "common - common - Dependent Record Found In Dependent Table")
				util.Respond(r, w, nil, 406, "")
				return
			}
		}
	}
	// 2020-05-14 - HK - END

	if status == 1 {
		if apiModule == "room_type" {
			// Image Count Validation
			ImageCount, err := partner.GetHotelRoomImageCount(r, reqMap.ID)
			if err != nil {
				util.LogIt(r, "common - common - Error getting room image count")
				util.RespondBadRequest(r, w)
				return
			}
			if ImageCount == 0 {
				util.LogIt(r, "common - common - No Images Are Binded To Room")
				util.Respond(r, w, nil, 406, "No Images Are Binded To Room")
				return
			}

			// Amenity Count Validation
			AmenityCount, err := partner.GetHotelRoomAmenityCount(r, reqMap.ID)
			if err != nil {
				util.LogIt(r, "common - common - Error getting room amenity count")
				util.RespondBadRequest(r, w)
				return
			}
			if AmenityCount == 0 {
				util.LogIt(r, "common - common - No Amenities Are Binded To Room")
				util.Respond(r, w, nil, 406, "No Amenities Are Binded To Room")
				return
			}
		}
	}

	Vmodule := util.StaticApiModules[apiModule].(map[string]string)
	flag, httpStatus, err := model.UpdateStatusModuleWise(r, Vmodule["api"], status, reqMap.ID)
	if flag == 0 || err != nil {
		// util.RespondWithError(r, w, "500")
		util.Respond(r, w, nil, httpStatus, "100014")
		return
	}

	util.Respond(r, w, nil, 204, "")
}

// DeleteStatusModuleWise - Delete Status
func DeleteStatusModuleWise(w http.ResponseWriter, r *http.Request) {
	util.LogIt(r, "common - common - DeleteStatusModuleWise")
	defer util.CommonDeferred(w, r, "common", "common", "DeleteStatusModuleWise")
	if r.URL.Query().Get("hotelid") != "" {
		context.Set(r, "HotelId", r.URL.Query().Get("hotelid"))
	}

	vars := mux.Vars(r)
	id := vars["id"]

	apiModule := vars["module"]
	if _, ok := util.StaticApiModules[apiModule]; !ok {
		util.Respond(r, w, nil, 406, "100011")
		return
	}

	// 2020-06-20 - HK - START
	// Purpose : Dependent Record Checking Added Before Deleting Record
	if _, ok := util.DependentTables[apiModule]; ok {
		// cnt, err := model.CheckDependentRecord(r, id, apiModule)
		var cnt int64
		var err error

		// If condition added here for special check for property_tag module because for all dependent
		// table moules, all have status column exists in their dependent table. But in the case of
		// property_tag we do have a slight different mechanism, where we remove tags directly upon
		// update operation and hence status column was not taken into cf_hotel_tag table. - HK - 2021-04-19
		if apiModule == "property_tag" || apiModule == "amenity" {
			cnt, err = model.CheckDependentRecordWithoutStatus(r, id, apiModule)
		} else {
			cnt, err = model.CheckDependentRecord(r, id, apiModule)
		}
		if err != nil || cnt > 0 {
			util.LogIt(r, "common - common - Dependent Record Found In Dependent Table")
			util.Respond(r, w, nil, 406, "")
			return
		}
	}
	// 2020-06-20 - HK - END

	Vmodule := util.StaticApiModules[apiModule].(map[string]string)

	flag, httpStatus, err := model.UpdateStatusModuleWise(r, Vmodule["api"], 3, id)
	if flag == 0 || err != nil {
		// util.RespondWithError(r, w, "500")
		util.Respond(r, w, nil, httpStatus, "100014")
		return
	}

	util.Respond(r, w, nil, 204, "")
}

// UploadB64Image - Base 64 Image File Upload And Return Name Of An Image
func UploadB64Image(r *http.Request, base64ImageArr string, width int, height int, x int, y int, folderName string) (string, error) {
	util.LogIt(r, fmt.Sprint("Common - Common - UploadImage"))
	queryParm := strings.Split(base64ImageArr, ",")
	imgExt := queryParm[0]
	base64Image := queryParm[1]

	sDec, err := b64.StdEncoding.DecodeString(base64Image)
	if util.CheckErrorLog(r, err) {
		return "", err
	}

	var ranmonNumber = strconv.Itoa(rand.Intn(100))
	r1 := bytes.NewReader(sDec)
	imageName := ranmonNumber + "-" + time.Now().Local().Format("20060102150405")
	var im image.Image
	if imgExt == "data:image/jpeg;base64" || imgExt == "data:image/jpg;base64" || imgExt == "data:image/JPG;base64" || imgExt == "data:image/JPEG;base64" {
		im, err = jpeg.Decode(r1)
		imageName = imageName + ".jpg"
	} else if imgExt == "data:image/png;base64" || imgExt == "data:image/PNG;base64" {
		im, err = png.Decode(r1)
		imageName = imageName + ".png"
	} else if imgExt == "data:image/gif;base64" || imgExt == "data:image/GIF;base64" {
		im, err = gif.Decode(r1)
		imageName = imageName + ".gif"
	} else {
		return "", errors.New("No Image found")
	}

	if util.CheckErrorLog(r, err) {
		return "", err
	}
	CheckFolderExists(folderName)
	f, err := os.OpenFile(config.Env.StuffPath+folderName+"/"+imageName, os.O_WRONLY|os.O_CREATE, 0777)
	if util.CheckErrorLog(r, err) {
		return "", err
	}

	png.Encode(f, im)
	src, err := imaging.Open(config.Env.StuffPath + folderName + "/" + imageName)
	if util.CheckErrorLog(r, err) {
		return "", err
	}
	b := src.Bounds()
	imgWidth := b.Max.X
	imgHeight := b.Max.Y
	dstimg := imaging.Resize(src, imgWidth, imgHeight, imaging.Lanczos)
	err = imaging.Save(dstimg, config.Env.StuffPath+folderName+"/"+imageName)
	if util.CheckErrorLog(r, err) {
		return "", err
	}

	return imageName, nil
}

// UploadImageFormData -  Image Upload with form data
func UploadImageFormData(r *http.Request, folderName string, fileObj *multipart.FileHeader) (string, error) {
	util.LogIt(r, fmt.Sprint("Functions - Functions - UploadImageFormData"))
	var errors error
	var ImageName string
	file, err := fileObj.Open()
	if err != nil {
		return "", err
	}
	defer file.Close()

	fileInfo := fileObj.Filename

	var ranmonNumber = strconv.Itoa(rand.Intn(100))
	FileExt := strings.Split(fileInfo, ".")[1]
	ImageNameGenerate := ranmonNumber + "-" + time.Now().Local().Format("20060102150405") + "." + FileExt
	CheckFolderExists(folderName)
	f, err := os.OpenFile(config.Env.StuffPath+folderName+"/"+ImageNameGenerate, os.O_WRONLY|os.O_CREATE, 0777)
	if err != nil {
		errors = err
	}
	ImageName = ImageNameGenerate
	defer f.Close()
	io.Copy(f, file)

	_, err = UploadImageToS3(r, ImageName, folderName)
	if err != nil {
		return "", err
	}

	file.Close()
	f.Close()

	err = os.Remove(config.Env.StuffPath + folderName + "/" + ImageNameGenerate)
	if err != nil {
		return "", err
	}

	return ImageName, errors
}

// GetActivityWiseLog - Returns Activity Logs  By Passing Id Of Object, Module
func GetActivityWiseLog(w http.ResponseWriter, r *http.Request) {
	util.LogIt(r, "Common - Common - GetActivityWiseLog")
	defer util.CommonDeferred(w, r, "Common", "Common", "GetActivityWiseLog")
	var reqMap data.JQueryTableUI

	err := json.NewDecoder(r.Body).Decode(&reqMap)
	if err != nil {
		util.RespondBadRequest(r, w)
		return
	}

	var Limit = int(reqMap.Limit)
	var Offset = int(reqMap.Offset)
	var Sort = reqMap.Order[0]
	var Sort_Field = Sort.Field
	var direction = Sort.Direction

	vars := mux.Vars(r)
	apiModule := vars["module"]
	if _, ok := util.StaticApiModules[apiModule]; !ok {
		util.Respond(r, w, nil, 406, "100011")
		return
	}

	Vmodule := util.StaticApiModules[apiModule].(map[string]string)
	reqMap.External = Vmodule["api"]

	Data, err := model.GetActivityWiseLog(r, reqMap.External, reqMap.ID, Offset, Limit, direction, int(Sort_Field), reqMap.Search)
	if err != nil {
		util.RespondWithError(r, w, "500")
		return
	}

	util.Respond(r, w, Data, 200, "")

}

// GetCountry - Returns array of country
func GetCountry(w http.ResponseWriter, r *http.Request) {
	util.LogIt(r, "common - common - GetCountry")
	defer util.CommonDeferred(w, r, "common", "common", "GetCountry")

	flag, err := model.GetModuleData(r, "country", "")
	if err != nil {
		util.RespondWithError(r, w, "500")
		return
	}

	util.Respond(r, w, flag, 200, "")
}

// GetCity - Returns City array by passing state id
func GetCity(w http.ResponseWriter, r *http.Request) {
	util.LogIt(r, "common - common - GetCity")
	defer util.CommonDeferred(w, r, "common", "common", "GetCity")

	vars := mux.Vars(r)
	id := vars["id"]

	flag, err := model.GetModuleData(r, "city", id)
	if err != nil {
		util.RespondWithError(r, w, "500")
		return
	}

	util.Respond(r, w, flag, 200, "")
}

// GetState - Returns state array by passing country id
func GetState(w http.ResponseWriter, r *http.Request) {
	util.LogIt(r, "common - common - GetState")
	defer util.CommonDeferred(w, r, "common", "common", "GetState")

	vars := mux.Vars(r)
	id := vars["id"]

	flag, err := model.GetModuleData(r, "state", id)
	if err != nil {
		util.RespondWithError(r, w, "500")
		return
	}

	util.Respond(r, w, flag, 200, "")
}

// ChangeSortOrderOfModule - Update Status
func ChangeSortOrderOfModule(w http.ResponseWriter, r *http.Request) {
	util.LogIt(r, "common - common - ChangeSortOrderOfModule")
	defer util.CommonDeferred(w, r, "common", "common", "ChangeSortOrderOfModule")
	var reqMap data.SortOrder

	err := json.NewDecoder(r.Body).Decode(&reqMap)
	if err != nil {
		util.RespondBadRequest(r, w)
		return
	}

	vars := mux.Vars(r)
	id := vars["id"]
	reqMap.ID = id

	var SortOrder = int(reqMap.SortOrder)

	flag, err := model.ChangeSortOrderOfModule(r, "HOTEL_IMAGE", SortOrder, reqMap.ID)
	if flag == 0 || err != nil {
		util.RespondWithError(r, w, "500")
		return
	}

	util.Respond(r, w, nil, 204, "")
}

//UploadImageToS3 - Upload Image To S3
func UploadImageToS3(r *http.Request, fileName string, folderName string) (map[string]interface{}, error) {
	util.LogIt(r, fmt.Sprint("Common - Common - UploadImageToS3"))
	stuff := make(map[string]interface{})

	AWSKey := config.Env.AwsKey
	AWSSecretKey := config.Env.AwsSecret
	token := ""

	creds := credentials.NewStaticCredentials(AWSKey, AWSSecretKey, token)
	_, err := creds.Get()
	if util.CheckErrorLog(r, err) {
		return nil, err
	}

	cfg := aws.NewConfig().WithRegion(config.Env.AwsRegion).WithCredentials(creds)
	svc := s3.New(session.New(), cfg)
	file, err := os.Open(config.Env.StuffPath + folderName + "/" + fileName)
	if util.CheckErrorLog(r, err) {
		util.LogIt(r, fmt.Sprintf("Error - %s", err))
	}
	defer file.Close()

	fileInfo, _ := file.Stat()
	size := fileInfo.Size()
	buffer := make([]byte, size) // read file content to buffer
	file.Read(buffer)
	fileBytes := bytes.NewReader(buffer)
	fileType := http.DetectContentType(buffer)
	path := "/" + folderName + "/" + fileInfo.Name()

	params := &s3.PutObjectInput{
		Bucket:        aws.String(config.Env.AwsBucket),
		Key:           aws.String(path),
		Body:          fileBytes,
		ContentLength: aws.Int64(size),
		ContentType:   aws.String(fileType),
		ACL:           aws.String("public-read"),
	}

	resp, err := svc.PutObject(params)
	if util.CheckErrorLog(r, err) {
		return nil, err
		fmt.Sprintf("response %s", awsutil.StringValue(resp))
	}

	defer file.Close()

	stuff["file"] = "" + config.Env.StuffPath + folderName + "/" + fileName + ""
	return stuff, err
}

//DeleteImageFromS3 - Delete Image From S3
func DeleteImageFromS3(r *http.Request, folder string, image string) bool {
	util.LogIt(r, fmt.Sprint("Common - Common - DeleteImageFromS3"))
	token := ""
	creds := credentials.NewStaticCredentials(config.Env.AwsKey, config.Env.AwsSecret, token)
	_, err := creds.Get()
	if util.CheckErrorLog(r, err) {
		return false
	}

	FileName := folder + "/" + image
	cfg := aws.NewConfig().WithRegion(config.Env.AwsRegion).WithCredentials(creds)
	svc := s3.New(session.New(), cfg)

	params := &s3.DeleteObjectInput{
		Bucket: aws.String(config.Env.AwsBucket),
		Key:    aws.String(FileName),
	}

	_, err = svc.DeleteObject(params)
	if util.CheckErrorLog(r, err) {
		return false
	}

	return true
}

//CheckFolderExists - Check Folder Exists
func CheckFolderExists(folderName string) {
	if _, err := os.Stat(config.Env.StuffPath + folderName); os.IsExist(err) {
		// never triggers, because err is nil if file actually exists
	} else {
		if err := os.Mkdir(config.Env.StuffPath+folderName, 0755); os.IsExist(err) {
			// triggers if dir already exists
		}
	}
}

// AmenityCategoryList - Returns Categor List For Amenity
func AmenityCategoryList(w http.ResponseWriter, r *http.Request) {

	util.LogIt(r, "common - common - AmenityCategoryList")
	defer util.CommonDeferred(w, r, "common", "common", "AmenityCategoryList")

	util.RespondData(r, w, util.AmenityCtegory, 200)
}

// GetSelectAllData - Get Select All Data Created By Meet Soni 2Nd May 2020
func GetSelectAllData(w http.ResponseWriter, r *http.Request) {
	util.LogIt(r, "common - common - GetSelectAllData")
	defer util.CommonDeferred(w, r, "common", "common", "GetSelectAllData")
	id := util.GetRequestData(r, "id")
	apiModule := util.GetRequestData(r, "module")

	if _, ok := util.StaticApiModules[apiModule]; !ok {
		util.Respond(r, w, nil, 406, "100011")
		return
	}

	var status = "1"
	if id != "" {
		status = "1,2"
	}

	Vmodule := util.StaticApiModules[apiModule].(map[string]string)
	flag, err := model.GetSelectData(r, Vmodule["api"], id, true, status)
	if flag == 0 || err != nil {
		util.RespondWithError(r, w, "500")
		return
	}

	util.Respond(r, w, flag, 200, "")
}

// CommonChangeSortOrderOfModule - Update Status
func CommonChangeSortOrderOfModule(w http.ResponseWriter, r *http.Request) {
	util.LogIt(r, "common - common - CommonChangeSortOrderOfModule")
	defer util.CommonDeferred(w, r, "common", "common", "CommonChangeSortOrderOfModule")
	var reqMap data.SortOrder

	err := json.NewDecoder(r.Body).Decode(&reqMap)
	if err != nil {
		util.RespondBadRequest(r, w)
		return
	}

	vars := mux.Vars(r)
	id := vars["id"]
	reqMap.ID = id

	var SortOrder = int(reqMap.SortOrder)
	log.Println(SortOrder)

	apiModule := vars["module"]
	if _, ok := util.DependentStaticAPIModules[apiModule]; !ok {
		util.Respond(r, w, nil, 406, "100011")
		return
	}
	log.Println(apiModule)

	Vmodule := util.DependentStaticAPIModules[apiModule].(map[string]string)
	log.Println(Vmodule)

	log.Println(Vmodule["dep_module"])

	SortOrderChangeModule := Vmodule["dep_module"]

	flag, err := model.CommonChangeSortOrderOfModule(r, SortOrderChangeModule, SortOrder, reqMap.ID, apiModule)
	if flag == 0 || err != nil {
		util.RespondWithError(r, w, "500")
		return
	}

	util.Respond(r, w, nil, 204, "")
}

// GetLocalityList - Returns Locality by passing city id
func GetLocalityList(w http.ResponseWriter, r *http.Request) {
	util.LogIt(r, "common - common - GetLocalityList")
	defer util.CommonDeferred(w, r, "common", "common", "GetLocalityList")

	vars := mux.Vars(r)
	id := vars["id"]

	flag, err := model.GetModuleData(r, "locality", id)
	if err != nil {
		util.RespondWithError(r, w, "500")
		return
	}

	util.Respond(r, w, flag, 200, "")
}

// TaxCategoryList - Returns Categry List For Tax
func TaxCategoryList(w http.ResponseWriter, r *http.Request) {

	util.LogIt(r, "common - common - TaxCategoryList")
	defer util.CommonDeferred(w, r, "common", "common", "TaxCategoryList")

	util.RespondData(r, w, util.TaxInit, 200)
}

// UploadSingleImageFormData -  Single Image Upload with form data
func UploadSingleImageFormData(r *http.Request, folderName string, imageKey string) (string, error) {
	util.LogIt(r, fmt.Sprint("Functions - Functions - UploadSingleImageFormData"))

	var errors error
	file, fileInfo, err := r.FormFile(imageKey)
	if util.CheckErrorLog(r, err) {
		return "", err
	}

	s := strings.Split(fileInfo.Filename, ".")
	FileExt := s[len(s)-1]

	var ImageName string
	var ranmonNumber = strconv.Itoa(rand.Intn(100))
	ImageNameGenerate := ranmonNumber + "-" + time.Now().Local().Format("20060102150405") + "." + FileExt

	f, err := os.OpenFile(config.Env.StuffPath+folderName+"/"+ImageNameGenerate, os.O_WRONLY|os.O_CREATE, 0777)
	if err != nil {
		errors = err
	}
	ImageName = ImageNameGenerate
	defer f.Close()
	io.Copy(f, file)

	_, err = UploadImageToS3(r, ImageName, folderName)
	if err != nil {
		return "", err
	}

	file.Close()
	f.Close()

	err = os.Remove(config.Env.StuffPath + folderName + "/" + ImageNameGenerate)
	if err != nil {
		return "", err
	}

	return ImageName, errors
}

//WrapHandler - Wrap Handler
func WrapHandler(handler func(w http.ResponseWriter, r *http.Request), privileges string) func(w http.ResponseWriter, r *http.Request) {
	h := func(w http.ResponseWriter, r *http.Request) {
		if !CheckPrivilege(r, privileges) {
			util.Respond(r, w, nil, 403, "")
			return
		}
		handler(w, r) //handler
	}
	return h
}

// CheckPrivilege checks privileges and return if there is not authorize request - 2021-05-04 - HK
func CheckPrivilege(r *http.Request, privilege string) bool {
	util.LogIt(r, fmt.Sprint("Common - Common - CheckPrivilege"))

	flag := model.CheckPrivilege(r, privilege)
	return flag
}
