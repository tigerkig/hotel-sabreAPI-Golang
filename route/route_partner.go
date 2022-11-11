package route

import (
	"tp-system/controller"
	partner "tp-system/controller/partner"
	"tp-system/secure/middleware"

	"github.com/gorilla/mux"
)

// InitializePartnerRouters - Defines all partner panel routes of application
func InitializePartnerRouters(Route *mux.Router) {

	/**************************************************** [ FORM DATA API START] ***********************************************************/
	FormData := Route.PathPrefix("/form-cpi"+"/"+version).HeadersRegexp("Content-Type", "application/json|multipart/form-data").Subrouter().StrictSlash(true)
	FormData.HandleFunc("/hotel/image", partner.UploadHotelImage).Methods("POST")
	FormData.HandleFunc("/room_type/image", partner.UploadHotelRoomImage).Methods("POST")
	FormData.Use(middleware.FormDataPartnerAPIMiddleware)
	/**************************************************** [ FORM DATA API END ] ***********************************************************/

	/* Auth Router Start */
	web := Route.PathPrefix("/auth" + "/" + version).Subrouter().StrictSlash(true)
	web.HandleFunc("/login", controller.PartnerLogin).Methods("POST") //MS - 2020-04-20

	web.HandleFunc("/booking/inv_update", partner.UpdateHotelRoomRateDataForBooking).Methods("POST") //HK - 2020-08-01
	web.Use(middleware.WebMiddleware)

	api := Route.PathPrefix("/cpi" + "/" + version).Subrouter().StrictSlash(true)

	//Auth Details Of Particular Token
	api.HandleFunc("/authdetail/{token}", controller.GetConsoleAuthDetails).Methods("GET") //MS - 2020-04-20
	api.HandleFunc("/logout", controller.Logout).Methods("POST")                           //MS - 2020-04-20
	api.HandleFunc("/resetPassword", controller.ResetPassword).Methods("PUT")
	//Amenity Type
	api.HandleFunc("/amenity_type", controller.GetAmenityTypeList).Methods("GET") //MS - 2020-04-20

	//Amenity
	api.HandleFunc("/amenity", controller.GetAmenityList).Methods("GET")                                              //MS - 2020-04-20 - New Created
	api.HandleFunc("/amenity/icon/list", controller.IconList).Methods("GET")                                          //MS - 2020-04-20
	api.HandleFunc("/amenity/list", controller.AmenityTypeWiseAmenity).Methods("GET").Queries("hotelid", "{hotelid}") //MS - 2020-04-20

	// Common Amenity Use
	api.HandleFunc("/amenity_type/{id}", controller.GetAmenityTypeListV1).Methods("GET") //HK - 2020-05-06
	api.HandleFunc("/amenity/{id}", controller.GetAmenityListV1).Methods("GET")          //HK - 2020-05-06

	//Property Tags
	api.HandleFunc("/property_tag", controller.GetPropPertyTagList).Methods("GET") //MS - 2020-04-20

	//Bed Type
	api.HandleFunc("/bed_type", controller.GetBedTypeList).Methods("GET") //MS - 2020-04-20

	//Inclusion
	api.HandleFunc("/inclusion", controller.GetInclusionList).Methods("GET") //MS - 2020-04-20

	//Image Category
	api.HandleFunc("/image_category", controller.GetImageCategoryList).Methods("GET") //MS - 2020-04-20

	//Room View
	api.HandleFunc("/room_view", controller.GetRoomViewList).Methods("GET") //MS - 2020-04-20

	//Bed Type
	api.HandleFunc("/extra_bed_type", controller.GetExtraBedTypeList).Methods("GET") //MS - 2020-04-20

	//Room Type - Common
	api.HandleFunc("/room_type", partner.GetRoomTypeList).Methods("GET") //HK - 2020-05-09

	//Rate Plan - Common
	api.HandleFunc("/rate_plan/room_type/{id}", partner.GetRatePlanList).Methods("GET") //HK - 2020-05-30

	//Cancel Policy - Common
	api.HandleFunc("/can_policy/{hotelid}/list", partner.GetCancelPolicyList).Methods("GET") //HK - 2020-05-09

	//Cancel Policy - Common
	api.HandleFunc("/meal_type", controller.GetMealTypeList).Methods("GET") //HK - 2020-05-19

	// Property Type Common
	api.HandleFunc("/property_type", controller.GetPropertyTypeList).Methods("GET") //HK - 2020-06-24

	//Hotel
	api.HandleFunc("/hotel", partner.AddHotelByHotelier).Methods("POST") //MS - 2021-05-04
	api.HandleFunc("/hotel", partner.HotelListOfHotelier).Methods("GET") //MS - 2021-05-04
	// api.HandleFunc("/hotel/{id}", partner.UpdateHotelBasicInfo).Methods("PUT")                        //MS - 2020-04-20 //Changed at 2021-04-26
	api.HandleFunc("/hotel/bank/{id}", controller.UpdateBankDetails).Methods("PUT")                   //MS - 2020-04-20 //Changed at 2021-04-26
	api.HandleFunc("/hotel/location/{id}", partner.UpdateLocation).Methods("PUT")                     //MS - 2020-04-20 12:03 AM
	api.HandleFunc("/hotel/amenity/{id}", partner.UpdateAmenity).Methods("PUT")                       //MS - 2020-04-20 16:16 PM
	api.HandleFunc("/hotel/basic/{id}", partner.UpdateHotelBasicInfo).Methods("PUT")                  //MS - 2020-04-20 16:16 PM
	api.HandleFunc("/hotel/policy/{id}", partner.UpdatePolicyRules).Methods("PUT")                    //MS - 2020-04-21 11:02 AM
	api.HandleFunc("/hotel/image/sort_order/{id}", controller.ChangeSortOrderOfModule).Methods("PUT") //MS - 2020-04-21 02:20 PM
	api.HandleFunc("/hotel/image/{id}/delete/{hotelid}", partner.DeleteHotelImage).Methods("DELETE")  //MS - 2020-04-21 04:27 PM
	api.HandleFunc("/hotel/{id}", controller.ViewHotelInfo).Methods("GET")                            //MS - 2020-04-21
	// api.HandleFunc("/hotel/view", controller.ViewHotelInfo).Methods("GET")                            //MS - 2020-05-02
	api.HandleFunc("/hotel/image/listing", partner.GetHotelImageList).Methods("GET").Queries("hotelid", "{hotelid}") //HK - 2020-05-09

	api.HandleFunc("/hotel/dumpdata", partner.FillInvRateData).Methods("POST")                 //HK - 2020-05-14
	api.HandleFunc("/hotel/invdata", partner.GetHotelRoomRateData).Methods("POST")             //HK - 2020-05-15
	api.HandleFunc("/hotel/updateinvratedata", partner.UpdateHotelRoomRateData).Methods("PUT") //HK - 2020-06-03
	api.HandleFunc("/hotel/updatelog", partner.GetUpdateLogsOfProperty).Methods("POST")        //HK - 2020-08-08
	api.HandleFunc("/hotel/review_check/{id}", partner.GetReviewOfHotel).Methods("GET")        //HP - 2021-05-20
	api.HandleFunc("/hotel/verify_submit/{id}", partner.VerifyHotel).Methods("POST")           //HP - 2021-05-20

	//Common Api
	api.HandleFunc("/response", controller.GetResponseCodeMsg).Methods("GET")       //MS - 2020-04-20
	api.HandleFunc("/status", controller.StatusList).Methods("GET")                 //MS - 2020-04-20
	api.HandleFunc("/bank_list", controller.BankList).Methods("GET")                //MS - 2020-04-20
	api.HandleFunc("/logs/{module}", controller.GetActivityWiseLog).Methods("POST") //MS - 2020-04-20
	api.HandleFunc("/country", controller.GetCountry).Methods("GET")                //MS - 2020-04-20
	api.HandleFunc("/city/{id}", controller.GetCity).Methods("GET")                 //MS - 2020-04-20
	api.HandleFunc("/state/{id}", controller.GetState).Methods("GET")               //MS - 2020-04-20
	api.HandleFunc("/locality/{id}", controller.GetLocalityList).Methods("GET")     //HK - 2020-05-13

	//Hotel Cancellation Policy
	api.HandleFunc("/can_policy", partner.AddCancelPolicy).Methods("POST")                                         //HK - 2020-05-01
	api.HandleFunc("/can_policy/{id}", partner.UpdateCancelPolicy).Methods("PUT")                                  //HK - 2020-05-01
	api.HandleFunc("/can_policy/listing", partner.CancelPolicyListing).Methods("POST")                             //HK - 2020-05-01
	api.HandleFunc("/can_policy/{id}", partner.GetCancelPolicyInfo).Methods("GET").Queries("hotelid", "{hotelid}") //HK - 2020-06-24

	//Hotel Room Type
	api.HandleFunc("/room_type", partner.AddRoomType).Methods("POST")             //HK - 2020-05-06
	api.HandleFunc("/room_type/{id}", partner.UpdateRoomType).Methods("PUT")      //HK - 2020-05-07
	api.HandleFunc("/room_type/{id}", partner.GetRoomType).Methods("GET")         //HK - 2020-05-08
	api.HandleFunc("/room_type/listing", partner.RoomTypeListing).Methods("POST") //HK - 2020-05-08

	api.HandleFunc("/room_type/amenity/{id}", partner.UpdateRoomAmenity).Methods("PUT")                     //HK - 2020-05-08
	api.HandleFunc("/room_type/amenity/list/{id}", controller.AmenityTypeWiseAmenityForRoom).Methods("GET") //HK - 2020-05-08

	api.HandleFunc("/room_type/image/{id}/delete", partner.DeleteHotelRoomImage).Methods("DELETE") //HK - 2020-05-08
	// api.HandleFunc("/{module}/image/sort_order/{id}", controller.CommonChangeSortOrderOfModule).Methods("PUT") //HK - 2020-05-08
	api.HandleFunc("/{module}/image/sort_order/{id}", partner.SortRoomImage).Methods("PUT")  //HK - 2021-04-27
	api.HandleFunc("/room_type/image/{id}/listing", partner.GetRoomImageList).Methods("GET") //HK - 2020-05-08

	//Hotel Rate Plan Room Wise
	api.HandleFunc("/rate_plan", partner.AddRatePlan).Methods("POST")             //HK - 2020-05-09
	api.HandleFunc("/rate_plan/{id}", partner.UpdateRatePlan).Methods("PUT")      //HK - 2020-05-09
	api.HandleFunc("/rate_plan/{id}", partner.GetRatePlan).Methods("GET")         //HK - 2020-05-09
	api.HandleFunc("/rate_plan/listing", partner.RatePlanListing).Methods("POST") //HK - 2020-05-09

	//Hotel Tax
	api.HandleFunc("/tax", partner.AddTax).Methods("POST")                          //HK - 2020-06-04
	api.HandleFunc("/tax/{id}", partner.UpdateTax).Methods("PUT")                   //HK - 2020-06-05
	api.HandleFunc("/tax/{id}", partner.GetTaxInfo).Methods("GET")                  //HK - 2020-06-05
	api.HandleFunc("/tax/listing", partner.TaxListing).Methods("POST")              //HK - 2020-06-05
	api.HandleFunc("/tax/category/list", controller.TaxCategoryList).Methods("GET") //HK - 2020-06-05

	//Display Setting
	api.HandleFunc("/display_setting", controller.GetDisplaySettings).Methods("GET") //MS - 2020-07-28

	//Update & Delete Status
	api.HandleFunc("/{module}/{id}/status", controller.UpdateStatusModuleWise).Methods("PUT")    //HK - 2020-05-01
	api.HandleFunc("/{module}/{id}/delete", controller.DeleteStatusModuleWise).Methods("DELETE") //HK - 2020-05-01

	// Sync Hotel Inv
	api.HandleFunc("/sync_inv", partner.AddInv).Methods("GET")

	api.Use(middleware.PartnerMiddleware)
}
