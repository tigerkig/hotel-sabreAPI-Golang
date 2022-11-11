package secure

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"tp-system/config"
	"tp-system/model"

	//
	_ "github.com/go-sql-driver/mysql"
	"github.com/stripe/stripe-go"
	"gopkg.in/mgo.v2"

	"tp-system/route"

	"github.com/go-redis/redis"
	"github.com/gorilla/mux"
)

// Initialization - It will initialize app and maintain mysql, Redis Database connection and initialize Routes
func Initialization() (*mux.Router, error) {
	log.Println("Initialization DB")

	// psqlInfo := fmt.Sprintf(config.Env.MSQL.User + ":" + config.Env.MSQL.Password + "@/" + config.Env.MSQL.DBName)
	psqlInfo := fmt.Sprintf(config.Env.MSQL.User + ":" + config.Env.MSQL.Password + "@(" + config.Env.MSQL.ServerHost + ":" + strconv.Itoa(config.Env.MSQL.Port) + ")/" + config.Env.MSQL.DBName)
	dba, err := sql.Open("mysql", psqlInfo)
	if err != nil {
		panic(err)
	}

	client := redis.NewClient(&redis.Options{
		Addr:     config.Env.Redis.RedisHost,
		Password: config.Env.Redis.RedisPassword,
		DB:       config.Env.Redis.RedisDB,
	})

	mSession, err := mgo.Dial(config.Env.Mongo.MongoHost)
	if err != nil {
		panic(err)
	}

	model.VDB = dba
	model.RedisClient = client
	model.VMongoSession = mSession
	config.Env.MGOSession = mSession
	var Router = mux.NewRouter()
	//Init stripe key while starting program
	stripe.Key = config.Env.Stripe.Secret

	route.InitializeRouters(Router)
	route.InitializePartnerRouters(Router)
	go InitLogRotate()
	go model.HandleCaching()
	go model.HandleSendMailObject()

	// go CronRoutine() // 2021-05-11 - HK - Daily Cron
	return Router, nil
}
