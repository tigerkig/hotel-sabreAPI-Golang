package route

import (
	"net/http"
	"tp-system/config"
	"tp-system/controller"
	"tp-system/controller/front"
	"tp-system/secure/middleware"

	"github.com/gorilla/mux"
)

const version = "v1"

// InitializeRouters - Defines all routes of application
func InitializeRouters(Route *mux.Router) {

	StripeAPI := Route.PathPrefix("/stripe").Subrouter().StrictSlash(true)
	StripeAPI.HandleFunc("/reauth/{id}", controller.ReCreateAccountLinks).Methods("GET") //MS - 2021-05-20
	StripeAPI.HandleFunc("/return/{id}", controller.ReturnAccountURL).Methods("GET")     //MS - 2021-05-21
	StripeAPI.Use(middleware.WebMiddleware)
	/* Auth Router Start */

	CrossAPI := Route.PathPrefix("/cross" + "/" + version).Subrouter().StrictSlash(true)

	CrossAPI.HandleFunc("/sms/{id}", controller.GetSMSGatewayAndTemplate).Methods("GET")              //MS - 2020-08-20
	CrossAPI.HandleFunc("/email_template/{id}", controller.GetEmailTemplateDetailInfo).Methods("GET") //MS - 2020-07-21
	CrossAPI.HandleFunc("/payment_gateway", controller.GetActivatePaymentGateway).Methods("GET")      //MS - 2020-07-29
	CrossAPI.HandleFunc("/create_card/{id}", controller.CreateCard).Methods("GET")                    //HP - 2021-10-31
	CrossAPI.HandleFunc("/virtual_card_info/{id}", controller.VirtualCardInfo).Methods("GET")         //HP - 2021-10-31

	CrossAPI.Use(middleware.CrossMiddleware)

	/**************************************************** [ FORM DATA API START] ***********************************************************/
	FormData := Route.PathPrefix("/form-api"+"/"+version).HeadersRegexp("Content-Type", "application/json|multipart/form-data").Subrouter().StrictSlash(true)

	FormData.HandleFunc("/property_type", controller.WrapHandler(controller.AddPropertyTypeNew, "24")).Methods("POST")
	FormData.HandleFunc("/property_type/{id}", controller.WrapHandler(controller.UpdatePropertyTypeNew, "24")).Methods("PUT")

	// Popular City
	FormData.HandleFunc("/popular_city", controller.WrapHandler(controller.AddPopularCity, "32")).Methods("POST")        //HK - 2020-06-09
	FormData.HandleFunc("/popular_city/{id}", controller.WrapHandler(controller.UpdatePopularCity, "32")).Methods("PUT") //HK - 2020-06-09

	// Company Info
	FormData.HandleFunc("/company_info/{id}", controller.WrapHandler(controller.UpdateCompanyInfo, "34")).Methods("PUT") //HK - 2020-06-10

	FormData.Use(middleware.FormDataAPIMiddleware)

	web := Route.PathPrefix("/web" + "/" + version).Subrouter().StrictSlash(true)
	web.HandleFunc("/login", controller.Login).Methods("POST")              //MS - 2020-04-11
	web.HandleFunc("/settings", controller.GetFrontSettings).Methods("GET") //Umesh - 2020-05-31
	web.HandleFunc("/static", controller.GetFrontStaticData).Methods("GET") //Meet - 2020-08-18

	web.HandleFunc("/home_page_data", controller.GetHomePageData).Methods("GET")      //Umesh - 2020-05-31
	web.HandleFunc("/rating_questions", controller.GetRatingQuestions).Methods("GET") //Umesh - 2020-05-31
	web.HandleFunc("/special_request", controller.GetSpecialRequest).Methods("GET")   //Umesh - 2020-05-31
	web.HandleFunc("/cms/{slug}", controller.GetCmsData).Methods("GET")               //HP - 2021-05-29
	web.HandleFunc("/cms", controller.GetCmsListData).Methods("GET")                  //HP - 2021-05-29
	// On Boarding Inquiry
	web.HandleFunc("/hotel_inquiry", controller.AddListingInquiry).Methods("POST") //HK - 2020-06-17
	web.HandleFunc("/country", controller.GetCountry).Methods("GET")               //HK - 2020-06-17
	web.HandleFunc("/city/{id}", controller.GetCity).Methods("GET")                //HK - 2020-06-17
	web.HandleFunc("/state/{id}", controller.GetState).Methods("GET")              //HK - 2020-06-17

	web.HandleFunc("/company_info/{id}", controller.GetCompanyInfo).Methods("GET")           //MS - 2020-07-24
	web.HandleFunc("/list_your_property", front.AddListProperty).Methods("POST")             //MS - 2021-04-29
	web.HandleFunc("/verify_token/{token}", front.VerifyPartnerActivateToken).Methods("GET") //MS - 2021-04-30
	web.Use(middleware.WebMiddleware)

	api := Route.PathPrefix("/api" + "/" + version).Subrouter().StrictSlash(true)
	//Auth Details Of Particular Token
	api.HandleFunc("/authdetail/{token}", controller.GetAuthDetails).Methods("GET") //MS - 2020-04-11
	api.HandleFunc("/logout", controller.Logout).Methods("POST")                    //MS - 2020-04-11
	api.HandleFunc("/resetPassword", controller.ResetPassword).Methods("PUT")
	//Update & Delete Status
	api.HandleFunc("/{module}/{id}/status", controller.UpdateStatusModuleWise).Methods("PUT")                  //MS - 2020-04-13
	api.HandleFunc("/{module}/{id}/delete", controller.DeleteStatusModuleWise).Methods("DELETE")               //MS - 2020-04-13
	api.HandleFunc("/hotel/{id}/propertyStatus/{propertyStatus}", controller.UpdateHotelStatus).Methods("PUT") //MS - 2021-04-19

	//Get Privileges List
	api.HandleFunc("/privileges", controller.GetPrivilegeList).Methods("GET") //MS - 2020-04-13

	//User Role
	api.HandleFunc("/user_role", controller.WrapHandler(controller.AddUserRole, "4")).Methods("POST")                                   //MS - 2020-04-13
	api.HandleFunc("/user_role/{id}", controller.WrapHandler(controller.GetUserRole, "3")).Methods("GET")                               //MS - 2020-04-13
	api.HandleFunc("/user_role/{id}", controller.WrapHandler(controller.UpdateUserRole, "4")).Methods("PUT")                            //MS - 2020-04-13
	api.HandleFunc("/user_role/listing", controller.WrapHandler(controller.UserRoleListing, "3")).Methods("POST")                       //MS - 2020-04-13
	api.HandleFunc("/user_role/listing", controller.WrapHandler(controller.UserRoleListing, "3")).Methods("POST")                       //MS - 2020-04-13
	api.HandleFunc("/user_role", controller.GetUserRoleList).Methods("GET")                                                             //MS - 2020-04-13
	api.HandleFunc("/{module:user_role}/{id}/status", controller.WrapHandler(controller.UpdateStatusModuleWise, "4")).Methods("PUT")    //HK - 2021-05-05
	api.HandleFunc("/{module:user_role}/{id}/delete", controller.WrapHandler(controller.DeleteStatusModuleWise, "4")).Methods("DELETE") //HK - 2021-05-05

	//User
	api.HandleFunc("/user", controller.WrapHandler(controller.AddUser, "6")).Methods("POST")                                       //MS - 2020-04-13
	api.HandleFunc("/user/{id}", controller.WrapHandler(controller.UpdateUser, "6")).Methods("PUT")                                //MS - 2020-04-13
	api.HandleFunc("/user/{id}", controller.WrapHandler(controller.GetUser, "5")).Methods("GET")                                   //MS - 2020-04-13
	api.HandleFunc("/user/listing", controller.WrapHandler(controller.UserListing, "5")).Methods("POST")                           //MS - 2020-04-13
	api.HandleFunc("/user/reset/{id}", controller.WrapHandler(controller.ResetPwdAndPrivileges, "6")).Methods("PUT")               //MS - 2020-04-13
	api.HandleFunc("/user", controller.GetUserList).Methods("GET")                                                                 //MS - 2020-04-15
	api.HandleFunc("/{module:user}/{id}/status", controller.WrapHandler(controller.UpdateStatusModuleWise, "6")).Methods("PUT")    //HK - 2021-05-05
	api.HandleFunc("/{module:user}/{id}/delete", controller.WrapHandler(controller.DeleteStatusModuleWise, "6")).Methods("DELETE") //HK - 2021-05-05

	//Amenity Type
	api.HandleFunc("/amenity_type", controller.WrapHandler(controller.AddAmenityType, "10")).Methods("POST")            //MS - 2020-04-14
	api.HandleFunc("/amenity_type/{id}", controller.WrapHandler(controller.UpdateAmenityType, "10")).Methods("PUT")     //MS - 2020-04-14
	api.HandleFunc("/amenity_type/listing", controller.WrapHandler(controller.AmenityTypeListing, "9")).Methods("POST") //MS - 2020-04-14
	// api.HandleFunc("/amenity_type", controller.GetAmenityTypeList).Methods("GET")                //MS - 2020-04-14
	api.HandleFunc("/amenity_type/{id}", controller.WrapHandler(controller.GetAmenityType, "9")).Methods("GET")                             //HK - 2020-05-04
	api.HandleFunc("/amenity_type/catg_list/{id}", controller.GetAmenityTypeListCatgWise).Methods("GET")                                    //HK - 2020-05-02
	api.HandleFunc("/amenity_type/category/list", controller.AmenityCategoryList).Methods("GET")                                            //HK - 2020-05-01
	api.HandleFunc("/amenity_type", controller.GetAmenityTypeList).Methods("GET")                                                           //MS - 2020-04-14
	api.HandleFunc("/{module:amenity_type}/{id}/status", controller.WrapHandler(controller.UpdateStatusModuleWise, "10")).Methods("PUT")    //HK - 2021-05-05
	api.HandleFunc("/{module:amenity_type}/{id}/delete", controller.WrapHandler(controller.DeleteStatusModuleWise, "10")).Methods("DELETE") //HK - 2021-05-05

	//Amenity
	api.HandleFunc("/amenity", controller.WrapHandler(controller.AddAmenity, "12")).Methods("POST")                                    //MS - 2020-04-14
	api.HandleFunc("/amenity/{id}", controller.WrapHandler(controller.UpdateAmenity, "12")).Methods("PUT")                             //MS - 2020-04-14
	api.HandleFunc("/amenity/{id}", controller.WrapHandler(controller.GetAmenity, "11")).Methods("GET")                                //MS - 2020-04-14
	api.HandleFunc("/amenity/listing", controller.WrapHandler(controller.AmenityListing, "11")).Methods("POST")                        //MS - 2020-04-14
	api.HandleFunc("/amenity/icon/list", controller.IconList).Methods("GET")                                                           //MS - 2020-04-17
	api.HandleFunc("/{module:amenity}/{id}/status", controller.WrapHandler(controller.UpdateStatusModuleWise, "12")).Methods("PUT")    //HK - 2021-05-05
	api.HandleFunc("/{module:amenity}/{id}/delete", controller.WrapHandler(controller.DeleteStatusModuleWise, "12")).Methods("DELETE") //HK - 2021-05-05

	//Property Tags
	api.HandleFunc("/property_tag", controller.WrapHandler(controller.AddPropPertyTag, "18")).Methods("POST")                               //MS - 2020-04-15
	api.HandleFunc("/property_tag/{id}", controller.WrapHandler(controller.UpdatePropPertyTag, "18")).Methods("PUT")                        //MS - 2020-04-15
	api.HandleFunc("/property_tag/listing", controller.WrapHandler(controller.PropPertyTagListing, "17")).Methods("POST")                   //MS - 2020-04-15
	api.HandleFunc("/property_tag", controller.GetPropPertyTagList).Methods("GET")                                                          //MS - 2020-04-15
	api.HandleFunc("/{module:property_tag}/{id}/status", controller.WrapHandler(controller.UpdateStatusModuleWise, "18")).Methods("PUT")    //HK - 2021-05-05
	api.HandleFunc("/{module:property_tag}/{id}/delete", controller.WrapHandler(controller.DeleteStatusModuleWise, "18")).Methods("DELETE") //HK - 2021-05-05

	//Bed Type
	api.HandleFunc("/bed_type", controller.WrapHandler(controller.AddBedType, "14")).Methods("POST")                                    //MS - 2020-04-15
	api.HandleFunc("/bed_type/{id}", controller.WrapHandler(controller.UpdateBedType, "14")).Methods("PUT")                             //MS - 2020-04-15
	api.HandleFunc("/bed_type/listing", controller.WrapHandler(controller.BedTypeListing, "13")).Methods("POST")                        //MS - 2020-04-15
	api.HandleFunc("/bed_type", controller.GetBedTypeList).Methods("GET")                                                               //MS - 2020-04-15
	api.HandleFunc("/bed_type/{id}", controller.WrapHandler(controller.GetBedType, "13")).Methods("GET")                                //HK - 2021-04-21
	api.HandleFunc("/{module:bed_type}/{id}/status", controller.WrapHandler(controller.UpdateStatusModuleWise, "14")).Methods("PUT")    //HK - 2021-05-05
	api.HandleFunc("/{module:bed_type}/{id}/delete", controller.WrapHandler(controller.DeleteStatusModuleWise, "14")).Methods("DELETE") //HK - 2021-05-05

	//Inclusion
	api.HandleFunc("/inclusion", controller.WrapHandler(controller.AddInclusion, "22")).Methods("POST")                                  //MS - 2020-04-15
	api.HandleFunc("/inclusion/{id}", controller.WrapHandler(controller.UpdateInclusion, "22")).Methods("PUT")                           //MS - 2020-04-15
	api.HandleFunc("/inclusion/listing", controller.WrapHandler(controller.InclusionListing, "21")).Methods("POST")                      //MS - 2020-04-15
	api.HandleFunc("/inclusion", controller.GetInclusionList).Methods("GET")                                                             //MS - 2020-04-15
	api.HandleFunc("/inclusion/{id}", controller.WrapHandler(controller.GetInclusion, "21")).Methods("GET")                              //HK - 2021-04-20
	api.HandleFunc("/{module:inclusion}/{id}/status", controller.WrapHandler(controller.UpdateStatusModuleWise, "22")).Methods("PUT")    //HK - 2021-05-05
	api.HandleFunc("/{module:inclusion}/{id}/delete", controller.WrapHandler(controller.DeleteStatusModuleWise, "22")).Methods("DELETE") //HK - 2021-05-05

	//Image Category
	api.HandleFunc("/image_category", controller.WrapHandler(controller.AddImageCategory, "20")).Methods("POST")                              //MS - 2020-04-16
	api.HandleFunc("/image_category/listing", controller.WrapHandler(controller.ImageCategoryListing, "19")).Methods("POST")                  //MS - 2020-04-16
	api.HandleFunc("/image_category/{id}", controller.WrapHandler(controller.UpdateImageCategory, "20")).Methods("PUT")                       //MS - 2020-04-16
	api.HandleFunc("/image_category", controller.GetImageCategoryList).Methods("GET")                                                         //MS - 2020-04-16
	api.HandleFunc("/{module:image_category}/{id}/status", controller.WrapHandler(controller.UpdateStatusModuleWise, "20")).Methods("PUT")    //HK - 2021-05-05
	api.HandleFunc("/{module:image_category}/{id}/delete", controller.WrapHandler(controller.DeleteStatusModuleWise, "20")).Methods("DELETE") //HK - 2021-05-05

	//Hotel 51 52
	api.HandleFunc("/hotel", controller.WrapHandler(controller.AddHotelInfo, "52")).Methods("POST")                           //MS - 2020-04-16
	api.HandleFunc("/hotel/listing", controller.WrapHandler(controller.HotelListing, "51")).Methods("POST")                   //MS - 2020-04-16
	api.HandleFunc("/hotel/commission/{id}", controller.WrapHandler(controller.UpdateCommissionSetting, "52")).Methods("PUT") //MS - 2020-04-16
	api.HandleFunc("/hotel/account/{id}", controller.WrapHandler(controller.UpdateAccountManager, "52")).Methods("PUT")       //MS - 2020-04-16 - Finished At 17Th
	api.HandleFunc("/hotel/{id}", controller.WrapHandler(controller.UpdateHotelInfo, "52")).Methods("PUT")                    //MS - 2020-04-17
	api.HandleFunc("/hotel/reset/{id}", controller.WrapHandler(controller.ResetHotelUserPwd, "52")).Methods("PUT")            //MS - 2020-04-17
	api.HandleFunc("/hotel/view/{id}", controller.WrapHandler(controller.ViewHotelInfo, "51")).Methods("GET")                 //MS - 2020-04-17
	api.HandleFunc("/hotel/{id}", controller.WrapHandler(controller.GetHotelInfo, "51")).Methods("GET")                       //MS - 2020-04-17
	api.HandleFunc("/hotel/bank/{id}", controller.WrapHandler(controller.UpdateBankDetails, "52")).Methods("PUT")             //MS - 2020-04-17
	api.HandleFunc("/hotel", controller.GetHotelListForOtherModule).Methods("GET")                                            //HK - 2020-06-18
	api.HandleFunc("/hotel/live/{id}", controller.WrapHandler(controller.UpdateHotelStatusToLive, "52")).Methods("PUT")       //HK - 2020-06-26
	api.HandleFunc("/hotel/client/listing", controller.HotelierListing).Methods("POST")                                       //MS - 2021-05-04
	api.HandleFunc("/hotel/approved/{id}", controller.ApprovedHotel).Methods("PUT")                                           //MS - 2021-05-04
	api.HandleFunc("/booking/{id}", controller.UpdateBooking).Methods("PUT")                                                  //HP - 2022-04-29

	//Room View
	api.HandleFunc("/room_view", controller.WrapHandler(controller.AddRoomView, "8")).Methods("POST")                                   //MS - 2020-04-17
	api.HandleFunc("/room_view/{id}", controller.WrapHandler(controller.UpdateRoomView, "8")).Methods("PUT")                            //MS - 2020-04-17
	api.HandleFunc("/room_view/listing", controller.WrapHandler(controller.RoomViewListing, "7")).Methods("POST")                       //MS - 2020-04-17
	api.HandleFunc("/room_view", controller.GetRoomViewList).Methods("GET")                                                             //MS - 2020-04-17
	api.HandleFunc("/room_view/{id}", controller.WrapHandler(controller.GetRoomView, "7")).Methods("GET")                               //HK - 2021-04-21
	api.HandleFunc("/{module:room_view}/{id}/status", controller.WrapHandler(controller.UpdateStatusModuleWise, "8")).Methods("PUT")    //HK - 2021-05-05
	api.HandleFunc("/{module:room_view}/{id}/delete", controller.WrapHandler(controller.DeleteStatusModuleWise, "8")).Methods("DELETE") //HK - 2021-05-05

	//Extra Bed Type
	api.HandleFunc("/extra_bed_type", controller.WrapHandler(controller.AddExtraBedType, "16")).Methods("POST")                               //MS - 2020-04-18
	api.HandleFunc("/extra_bed_type/{id}", controller.WrapHandler(controller.UpdateExtraBedType, "16")).Methods("PUT")                        //MS - 2020-04-18
	api.HandleFunc("/extra_bed_type/listing", controller.WrapHandler(controller.ExtraBedTypeListing, "15")).Methods("POST")                   //MS - 2020-04-18
	api.HandleFunc("/extra_bed_type", controller.GetExtraBedTypeList).Methods("GET")                                                          //MS - 2020-04-18
	api.HandleFunc("/{module:extra_bed_type}/{id}/status", controller.WrapHandler(controller.UpdateStatusModuleWise, "16")).Methods("PUT")    //HK - 2021-05-05
	api.HandleFunc("/{module:extra_bed_type}/{id}/delete", controller.WrapHandler(controller.DeleteStatusModuleWise, "16")).Methods("DELETE") //HK - 2021-05-05

	//Special Request
	api.HandleFunc("/special_request", controller.WrapHandler(controller.AddSpecialRequest, "26")).Methods("POST")                             //MS - 2020-04-22
	api.HandleFunc("/special_request/{id}", controller.WrapHandler(controller.UpdateSpecialRequest, "26")).Methods("PUT")                      //HK - 2020-04-28
	api.HandleFunc("/special_request/listing", controller.WrapHandler(controller.SpecialRequestListing, "25")).Methods("POST")                 //HK - 2020-04-28
	api.HandleFunc("/{module:special_request}/{id}/status", controller.WrapHandler(controller.UpdateStatusModuleWise, "26")).Methods("PUT")    //HK - 2021-05-05
	api.HandleFunc("/{module:special_request}/{id}/delete", controller.WrapHandler(controller.DeleteStatusModuleWise, "26")).Methods("DELETE") //HK - 2021-05-05

	//City,State,Country API
	api.HandleFunc("/country", controller.GetCountry).Methods("GET")                 //MS - 2020-04-16
	api.HandleFunc("/city/{id}", controller.GetCity).Methods("GET")                  //MS - 2020-04-16
	api.HandleFunc("/state/{id}", controller.GetState).Methods("GET")                //MS - 2020-04-16
	api.HandleFunc("/city/locality/{id}", controller.GetLocalityList).Methods("GET") //HK - 2020-05-13

	//Response Code Msg
	api.HandleFunc("/response", controller.GetResponseCodeMsg).Methods("GET")       //MS - 2020-04-13
	api.HandleFunc("/status", controller.StatusList).Methods("GET")                 //MS - 2020-04-14
	api.HandleFunc("/bank_list", controller.BankList).Methods("GET")                //MS - 2020-04-17
	api.HandleFunc("/logs/{module}", controller.GetActivityWiseLog).Methods("POST") //MS - 2020-04-14 - Finished at 15Th April

	//Property Type
	api.HandleFunc("/property_type", controller.WrapHandler(controller.AddPropertyTypeNew, "24")).Methods("POST")                            //HK - 2020-04-27
	api.HandleFunc("/property_type/{id}", controller.WrapHandler(controller.UpdatePropertyTypeNew, "24")).Methods("PUT")                     //HK - 2020-04-27
	api.HandleFunc("/property_type/{id}", controller.WrapHandler(controller.GetPropertyTypeInfo, "23")).Methods("GET")                       //HK - 2020-06-05
	api.HandleFunc("/property_type/listing", controller.WrapHandler(controller.PropertyTypeListing, "23")).Methods("POST")                   //HK - 2020-04-27
	api.HandleFunc("/property_type", controller.GetPropertyTypeList).Methods("GET")                                                          //HK - 2020-06-18
	api.HandleFunc("/{module:property_type}/{id}/status", controller.WrapHandler(controller.UpdateStatusModuleWise, "24")).Methods("PUT")    //HK - 2021-05-05
	api.HandleFunc("/{module:property_type}/{id}/delete", controller.WrapHandler(controller.DeleteStatusModuleWise, "24")).Methods("DELETE") //HK - 2021-05-05

	//Meal Type
	api.HandleFunc("/meal_type", controller.WrapHandler(controller.AddMealType, "28")).Methods("POST")                                   //HK - 2020-05-04
	api.HandleFunc("/meal_type/{id}", controller.WrapHandler(controller.UpdateMealType, "28")).Methods("PUT")                            //HK - 2020-05-04
	api.HandleFunc("/meal_type/listing", controller.WrapHandler(controller.MealTypeListing, "27")).Methods("POST")                       //HK - 2020-05-04
	api.HandleFunc("/{module:meal_type}/{id}/status", controller.WrapHandler(controller.UpdateStatusModuleWise, "28")).Methods("PUT")    //HK - 2021-05-05
	api.HandleFunc("/{module:meal_type}/{id}/delete", controller.WrapHandler(controller.DeleteStatusModuleWise, "28")).Methods("DELETE") //HK - 2021-05-05

	//Locality
	api.HandleFunc("/locality", controller.WrapHandler(controller.AddLocality, "30")).Methods("POST")                                   //HK - 2020-05-16
	api.HandleFunc("/locality/{id}", controller.WrapHandler(controller.UpdateLocality, "30")).Methods("PUT")                            //HK - 2020-05-16
	api.HandleFunc("/locality/{id}", controller.WrapHandler(controller.GetLocality, "29")).Methods("GET")                               //HK - 2020-05-16
	api.HandleFunc("/locality/listing", controller.WrapHandler(controller.LocalityListing, "29")).Methods("POST")                       //HK - 2020-05-16
	api.HandleFunc("/{module:locality}/{id}/status", controller.WrapHandler(controller.UpdateStatusModuleWise, "30")).Methods("PUT")    //HK - 2021-05-05
	api.HandleFunc("/{module:locality}/{id}/delete", controller.WrapHandler(controller.DeleteStatusModuleWise, "30")).Methods("DELETE") //HK - 2021-05-05

	//Popular City
	// api.HandleFunc("/popular_city", controller.AddPopularCity).Methods("POST")             //HK - 2020-06-09
	// api.HandleFunc("/popular_city/{id}", controller.UpdatePopularCity).Methods("PUT")      //HK - 2020-06-09
	api.HandleFunc("/popular_city/{id}", controller.WrapHandler(controller.GetPopularCityinfo, "31")).Methods("GET")                        //HK - 2020-06-09
	api.HandleFunc("/popular_city/listing", controller.WrapHandler(controller.PopularCityListing, "31")).Methods("POST")                    //HK - 2020-06-09
	api.HandleFunc("/popular_city/{id}/statuschange", controller.WrapHandler(controller.UpdatePopularCityStatus, "32")).Methods("PUT")      //HK - 2021-04-26
	api.HandleFunc("/{module:popular_city}/{id}/status", controller.WrapHandler(controller.UpdateStatusModuleWise, "32")).Methods("PUT")    //HK - 2021-05-05
	api.HandleFunc("/{module:popular_city}/{id}/delete", controller.WrapHandler(controller.DeleteStatusModuleWise, "32")).Methods("DELETE") //HK - 2021-05-05

	//Feedback/Question
	api.HandleFunc("/question", controller.WrapHandler(controller.AddQuestion, "42")).Methods("POST")                                   //HK - 2020-06-09
	api.HandleFunc("/question/{id}", controller.WrapHandler(controller.UpdateQuestion, "42")).Methods("PUT")                            //HK - 2020-06-09
	api.HandleFunc("/question/listing", controller.WrapHandler(controller.QuestionListing, "41")).Methods("POST")                       //HK - 2020-06-09
	api.HandleFunc("/{module:question}/{id}/status", controller.WrapHandler(controller.UpdateStatusModuleWise, "42")).Methods("PUT")    //HK - 2021-05-05
	api.HandleFunc("/{module:question}/{id}/delete", controller.WrapHandler(controller.DeleteStatusModuleWise, "42")).Methods("DELETE") //HK - 2021-05-05

	//Company Info
	api.HandleFunc("/company_info", controller.WrapHandler(controller.CompanyInfo, "33")).Methods("GET") //HK - 2020-06-10
	// api.HandleFunc("/company_info/{id}", controller.UpdateCompanyInfo).Methods("PUT") //HK - 2020-06-10
	api.HandleFunc("/company_info/{id}", controller.WrapHandler(controller.GetCompanyInfo, "33")).Methods("GET") //HK - 2020-06-10

	// Email Configuration
	api.HandleFunc("/email", controller.WrapHandler(controller.AddEmailConfig, "44")).Methods("POST")                                //HK - 2020-06-15
	api.HandleFunc("/email/{id}", controller.WrapHandler(controller.UpdateEmailConfig, "44")).Methods("PUT")                         //HK - 2020-06-16
	api.HandleFunc("/email/{id}", controller.WrapHandler(controller.GetEmailConfigInfo, "43")).Methods("GET")                        //HK - 2020-06-16
	api.HandleFunc("/email/listing", controller.WrapHandler(controller.EmailListing, "43")).Methods("POST")                          //HK - 2020-06-16
	api.HandleFunc("/{module:email}/{id}/status", controller.WrapHandler(controller.UpdateStatusModuleWise, "44")).Methods("PUT")    //HK - 2021-05-05
	api.HandleFunc("/{module:email}/{id}/delete", controller.WrapHandler(controller.DeleteStatusModuleWise, "44")).Methods("DELETE") //HK - 2021-05-05

	// On Board Inquiry Listing
	api.HandleFunc("/hotel_inquiry/listing", controller.WrapHandler(controller.OnBoardInquiryListing, "53")).Methods("POST")              //HK - 2020-06-17
	api.HandleFunc("/hotel_inquiry/{id}", controller.WrapHandler(controller.GetInquiryDetailInfo, "53")).Methods("GET")                   //HK - 2020-06-16
	api.HandleFunc("/{module:hotel_inquiry}/{id}/status", controller.WrapHandler(controller.UpdateStatusModuleWise, "54")).Methods("PUT") //HK - 2021-05-05

	// Add Property To Recommended List
	api.HandleFunc("/recommended", controller.WrapHandler(controller.AddHotelToRecommendedList, "40")).Methods("POST")                     //HK - 2020-06-18
	api.HandleFunc("/recommended/{id}", controller.WrapHandler(controller.UpdateHotelToRecommendedList, "40")).Methods("PUT")              //HK - 2020-06-18
	api.HandleFunc("/recommended/listing", controller.WrapHandler(controller.RecommendedHotelList, "39")).Methods("POST")                  //HK - 2020-06-18
	api.HandleFunc("/recommended/{id}", controller.WrapHandler(controller.GetRecommendedHotelInfo, "39")).Methods("GET")                   //HK - 2020-06-16
	api.HandleFunc("/{module:recommended}/{id}/status", controller.WrapHandler(controller.UpdateStatusModuleWise, "40")).Methods("PUT")    //HK - 2021-05-06
	api.HandleFunc("/{module:recommended}/{id}/delete", controller.WrapHandler(controller.DeleteStatusModuleWise, "40")).Methods("DELETE") //HK - 2021-05-06

	// General Setting
	api.HandleFunc("/gen_setting", controller.WrapHandler(controller.GetWebDefaultSettings, "35")).Methods("GET")    //HK - 2020-06-20
	api.HandleFunc("/gen_setting", controller.WrapHandler(controller.UpdateWebDefaultSettings, "36")).Methods("PUT") //HK - 2020-06-20

	//Email Template
	api.HandleFunc("/email_template/{id}", controller.WrapHandler(controller.UpdateEmailTemplate, "46")).Methods("PUT")      //MS - 2020-07-21
	api.HandleFunc("/email_template/{id}", controller.WrapHandler(controller.GetEmailTemplate, "45")).Methods("GET")         //MS - 2020-07-21
	api.HandleFunc("/email_template/listing", controller.WrapHandler(controller.EmailTemplateListing, "45")).Methods("POST") //MS - 2020-07-21
	api.HandleFunc("/init/email_template", controller.GetEmailList).Methods("GET")                                           //MS - 2020-07-24

	//Display Setting
	api.HandleFunc("/display_setting", controller.WrapHandler(controller.SetDisplaySetting, "38")).Methods("PUT")  //MS - 2020-07-28
	api.HandleFunc("/display_setting", controller.WrapHandler(controller.GetDisplaySettings, "37")).Methods("GET") //MS - 2020-07-28
	api.HandleFunc("/init/display_setting", controller.DisplaySettingInit).Methods("GET")                          //MS - 2020-07-28

	//Payment Gateway Configuration
	api.HandleFunc("/payment_gateway/{id}", controller.WrapHandler(controller.UpdatePaymentConfiguration, "50")).Methods("PUT")      //MS - 2020-07-29
	api.HandleFunc("/payment_gateway/{id}", controller.WrapHandler(controller.GetPaymentConfigDetail, "49")).Methods("GET")          //MS - 2020-07-29
	api.HandleFunc("/payment_gateway/activate/{id}", controller.WrapHandler(controller.ActivatePaymentGateway, "50")).Methods("PUT") //MS - 2020-07-29
	api.HandleFunc("/payment_gateway/listing", controller.WrapHandler(controller.PaymentGatewayListing, "49")).Methods("POST")       //MS - 2020-07-29
	api.HandleFunc("/payment_gateway", controller.GetPaymentGatewayList).Methods("GET")                                              //MS - 2020-07-29
	api.HandleFunc("/activate/payment_gateway", controller.GetActivatePaymentGateway).Methods("GET")                                 //MS - 2020-07-29

	//Hotel Configuration Review
	api.HandleFunc("/hotel_configuration/room_type", controller.GetRoomTypeList).Methods("GET")     //MS - 2020-08-14
	api.HandleFunc("/hotel_configuration/rate_plan", controller.GetRatePlanFromRoom).Methods("GET") //MS - 2020-08-14
	api.HandleFunc("/hotel_configuration/invdata", controller.GetHotelRoomRateData).Methods("GET")  //MS - 2020-08-14

	//SMS Template
	api.HandleFunc("/sms/{id}", controller.WrapHandler(controller.UpdateSMSTemplate, "48")).Methods("PUT")      //MS - 2020-08-20
	api.HandleFunc("/sms_gateway", controller.WrapHandler(controller.UpdateSMSGateway, "48")).Methods("PUT")    //MS - 2020-08-20
	api.HandleFunc("/sms/{id}", controller.WrapHandler(controller.GetSMSTemplate, "47")).Methods("GET")         //MS - 2020-08-20
	api.HandleFunc("/sms_gateway", controller.GetSmsGatewayDetail).Methods("GET")                               //MS - 2020-08-20
	api.HandleFunc("/sms/listing", controller.WrapHandler(controller.SMSTemplateListing, "47")).Methods("POST") //MS - 2020-08-20

	//Hotel Detailed Info
	api.HandleFunc("/hotel_detailed_info/{id}", controller.GetHotelDetailInfo) //MS - 2021-05-05

	//ListYourProperty Listing
	api.HandleFunc("/list_your_property/listing", front.ListYourPropertyListing).Methods("POST") //MS - 2021-05-11

	api.HandleFunc("/cms/{id}", controller.WrapHandler(controller.GetCms, "70")).Methods("GET")         //HP - 2021-05-29
	api.HandleFunc("/cms", controller.WrapHandler(controller.AddCms, "71")).Methods("POST")             //HP - 2021-05-29
	api.HandleFunc("/cms/{id}", controller.WrapHandler(controller.UpdateCms, "71")).Methods("PUT")      //HP - 2021-05-29
	api.HandleFunc("/cms/listing", controller.WrapHandler(controller.CmsListing, "70")).Methods("POST") //HP - 2021-05-29

	// Virtual Card
	api.HandleFunc("/virtual_card_active/{id}", controller.VirtualCardActive).Methods("GET") //HP - 2021-11-02
	api.HandleFunc("/virtual_card", controller.VirtualCardList).Methods("POST")              //HP - 2022-05-03
	api.HandleFunc("/virtual_card/{id}", controller.VirtualCardDetail).Methods("GET")        //HP - 2022-05-03

	api.Use(middleware.APIMiddleware)

	//Static File Handler
	fileHandler := http.StripPrefix("/stuff", (http.FileServer(http.Dir(config.Env.StuffPath))))
	Route.PathPrefix("/stuff/").Handler(fileHandler)
}
