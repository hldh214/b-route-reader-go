package config

import (
	"github.com/joho/godotenv"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var LOG_LEVEL = ""

var B_ROUTE_ID = ""
var B_ROUTE_PASSWORD = ""
var SERIAL = ""

var MQTT_BROKER = "tcp://localhost:1883"
var MQTT_USERNAME = ""
var MQTT_PASSWORD = ""

var MQTT_TOPIC_PREFIX = "brr"
var MQTT_TOPIC_DEVICE_NAME = "smartmeter"

// SKSCAN を繰り返す回数
var ACTIVE_SCAN_COUNT = 10

// 瞬間消費電力値を取得する間隔（秒）
var NOW_CONSUMPTION_WAIT = 20 * time.Second

// 積算消費電力値を取得する間隔（秒）
var TOTAL_CONSUMPTION_WAIT = 180 * time.Second

// EchonetLiteのGET時にリトライする回数
var MAX_ECHONET_GET_RETRY = 3

// ログに日時を出力するか（systemdで動かす場合はfalse）
var LOG_NO_DATETIME = true

// 環境変数からconfigをセット
func Initialize() {
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	exPath := filepath.Dir(ex)
	err = godotenv.Load(filepath.Join(exPath, ".env"))
	if err != nil {
		panic("Error loading .env file")
	}

	LOG_LEVEL = os.Getenv("LOG_LEVEL")
	B_ROUTE_ID = os.Getenv("B_ROUTE_ID")
	B_ROUTE_PASSWORD = os.Getenv("B_ROUTE_PASSWORD")
	SERIAL = os.Getenv("SERIAL")
	MQTT_BROKER = os.Getenv("MQTT_BROKER")
	MQTT_USERNAME = os.Getenv("MQTT_USERNAME")
	MQTT_PASSWORD = os.Getenv("MQTT_PASSWORD")

	if LOG_LEVEL == "" {
		LOG_LEVEL = "INFO"
	} else {
		LOG_LEVEL = strings.ToUpper(LOG_LEVEL)
		switch LOG_LEVEL {
		case "DEBUG":
		case "INFO":
		case "WARN":
		case "ERROR":
		default:
			panic("LOG_LEVEL is not in [DEBUG|INFO|WARN|ERROR]")
		}
	}

	if B_ROUTE_ID == "" {
		panic("B_ROUTE_ID env value is not set")
	}
	if B_ROUTE_PASSWORD == "" {
		panic("B_ROUTE_PASSWORD env value is not set")
	}
	if SERIAL == "" {
		panic("SERIAL env value is not set")
	}
	if MQTT_BROKER == "" {
		panic("MQTT_BROKER env value is not set")
	}

	noDatetime := os.Getenv("LOG_NO_DATETIME")
	LOG_NO_DATETIME = (strings.ToLower(noDatetime) == "true")

}
