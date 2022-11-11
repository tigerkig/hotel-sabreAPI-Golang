package config

import (
	"bytes"
	"log"
	"os"

	"github.com/BurntSushi/toml"
	"gopkg.in/mgo.v2"
)

// Env Represents All Structs From Toml
var Env struct {
	SystemPort         string       `toml:"system_port"`
	StuffPath          string       `toml:"system_path"`
	LogPath            string       `toml:"log_path"`
	EnableDebug        int          `toml:"enable_debug"`
	IsSSL              bool         `toml:"is_ssl_enable"`
	MSQL               *MSQL        `toml:"mysql"`
	Redis              *Redis       `toml:"redis"`
	Mongo              *Mongo       `toml:"mongo"`
	Stripe             *Stripe      `toml:"stripe"`
	AppKey             string       `toml:"app_key"`
	StuffURL           string       `toml:"stuff_url"`
	AwsKey             string       `toml:"awd_key_id"`
	AwsSecret          string       `toml:"aws_secret_key"`
	AwsRegion          string       `toml:"aws_region"`
	AwsBucket          string       `toml:"aws_bucket"`
	AwsBucketURL       string       `toml:"aws_bucket_url"`
	RoomFolder         string       `toml:"room_folder"`
	HotelFolder        string       `toml:"hotel_folder"`
	ProfileFolder      string       `toml:"profile_folder"`
	FrontURL           string       `toml:"front_url"`
	PropertyTypeFolder string       `toml:"property_type_folder"`
	PopularCityFolder  string       `toml:"popular_city_folder"`
	InvAuthKey         string       `toml:"invAuthKey"`
	FrontWebURL        string       `toml:"front_web_url"`
	PartnerWebURL      string       `toml:"partner_url"`
	MicroServiceType   string       `toml:"microservice"`
	MGOSession         *mgo.Session `toml:"mgo_session"`
}

// MSQL - Mysql Server Detail
type MSQL struct {
	ServerHost string `toml:"host"`
	Port       int    `toml:"port"`
	DBName     string `toml:"database_name"`
	User       string `toml:"user"`
	Password   string `toml:"password"`
}

// Redis - Redis Database
type Redis struct {
	RedisHost     string `toml:"redis_host"`
	RedisPassword string `toml:"redis_password"`
	RedisDB       int    `toml:"redis_db"`
}

type Stripe struct {
	Public string `toml:"StripeAppKey"`
	Secret string `toml:"StripeSecretKey"`
	APIURL string `toml:"StripeURL"`
}

// Mongo - Mongo Db Details
type Mongo struct {
	MongoHost     string `toml:"mongo_host"`
	MongoUser     string `toml:"mongo_user"`
	MongoPassword string `toml:"mongo_password"`
	MongoDB       string `toml:"mongo_db"`
}

// InitEnv - Initialize environment of system
func InitEnv() {
	var env bytes.Buffer
	env.WriteString(os.Getenv("TP_BACKEND"))
	log.Println(os.Getenv("TP_BACKEND"))
	env.WriteString("config.toml")
	if _, err := toml.DecodeFile(env.String(), &Env); err != nil {
		panic(err)
	}
}
