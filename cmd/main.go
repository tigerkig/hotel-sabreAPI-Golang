package main

import (
	"log"
	"net/http"
	"tp-system/config"
	"tp-system/model"
	"tp-system/secure"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func main() {
	var err error

	config.InitEnv()
	Route, err := secure.Initialization()
	if err != nil {
		panic(err)
	}
	model.SyncActiveProperty()
	//model.UpdateHotelDetailsAmenity("YkjA-EF6tUHT5n9ntIGCD")
	// model.UpdateSearchList("YkjA-EF6tUHT5n9ntIGCD", "List")
	//model.UpdateHotelOnList("YkjA-EF6tUHT5n9ntIGCD")
	// model.UpdateAllRoomType("YkjA-EF6tUHT5n9ntIGCD", "9yT_K5BN2bUWvzqfSmWOQ")
	// model.UpdateAllRoomType("YkjA-EF6tUHT5n9ntIGCD", "Yzjhm6Oxx0bPIHPml3zGT")
	// model.UpdateAllRateplanDeals("YkjA-EF6tUHT5n9ntIGCD", "eogJYEdKa_dD_BdF_oVPK")
	// model.UpdateAllRateplanDeals("YkjA-EF6tUHT5n9ntIGCD", "gnTGizmAT7i-4rupulPZa")
	// model.UpdateAllRateplanDeals("YkjA-EF6tUHT5n9ntIGCD", "RQrcDHaSNTqdqGQyL6YZy")
	// model.UpdateAllRateplanDeals("YkjA-EF6tUHT5n9ntIGCD", "VjT1b0BxNk6qs3zxFKmLc")
	// model.UpdateRatePlanDeals("Ze3ND7BCsdxbYzCq_qwDQ", "Q3J6MQhr9i0iBwqXFPqy7", "jXRXr4DhOOyHZXOz-7QLq")
	// model.UpdateRatePlanDeals("Ze3ND7BCsdxbYzCq_qwDQ", "Q3J6MQhr9i0iBwqXFPqy7", "AWIdiM41yNp0aF18GgDNu")
	// model.UpdateRatePlanDeals("Ze3ND7BCsdxbYzCq_qwDQ", "ciiwyjFySBPdPKp2HFW1-", "zcfVOgAb4sgDqpMmnEa-Y")
	// model.UpdateRatePlanDeals("Ze3ND7BCsdxbYzCq_qwDQ", "ciiwyjFySBPdPKp2HFW1-", "hNOpHOyqvqigBDN5sv5GT")
	// model.UpdateHotelTax("Ze3ND7BCsdxbYzCq_qwDQ")
	Run(Route)
}

// Run will start http
func Run(Route *mux.Router) {
	log.Println(config.Env.SystemPort)
	if config.Env.IsSSL == true {
		// For SSL add your chain and private key into tls.
		http.ListenAndServeTLS(config.Env.SystemPort, config.Env.StuffPath+"tls/fullchain.pem", config.Env.StuffPath+"tls/privkey.pem", handlers.CORS(handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization", "X-Auth-Token"}), handlers.AllowedMethods([]string{"GET", "POST", "PUT", "HEAD", "OPTIONS", "DELETE"}), handlers.AllowedOrigins([]string{"*"}))(Route))
	} else {
		// Without SSL
		log.Println("Without SSL")
		http.ListenAndServe(config.Env.SystemPort, handlers.CORS(handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization", "X-Auth-Token"}), handlers.AllowedMethods([]string{"GET", "POST", "PUT", "HEAD", "OPTIONS", "DELETE"}), handlers.AllowedOrigins([]string{"*"}))(Route))
	}
}
