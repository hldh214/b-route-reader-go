package ha

import (
	"encoding/json"
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/yakumo-saki/b-route-reader-go/src/config"
	"strings"
)

var DEVICE = map[string]string{
	"identifiers":  config.MQTT_TOPIC_DEVICE_NAME,
	"name":         "SA-M0 B-route Reader",
	"manufacturer": "Atmark Techno, Inc.",
	"model":        "Armadillo-Box WS1",
	"sw_version":   "Linux abws1-0 3.14.36-at11 #1 PREEMPT Sat Mar 31 02:35:12 JST 2018 armv5tejl GNU/Linux",
}

const (
	StatusConnecting = "Connecting"
	StatusConnected  = "Connected"
	StatusRunning    = "Running"
	StatusError      = "Error"
)

func InitTopic(client mqtt.Client) {
	topicPrefix := fmt.Sprintf("homeassistant/sensor/%s", config.MQTT_TOPIC_DEVICE_NAME)

	statusTopic := fmt.Sprintf("%s/%s_sensor_%s/config", topicPrefix, config.MQTT_TOPIC_DEVICE_NAME, Status)
	statusJson, _ := json.Marshal(map[string]interface{}{
		"state_topic":     StatusTopic,
		"icon":            "mdi:information",
		"device_class":    "enum",
		"attributes":      map[string]interface{}{"options": []string{StatusConnecting, StatusConnected, StatusRunning, StatusError}},
		"state_class":     "measurement",
		"entity_category": "diagnostic",
		"name":            "Status",
		"object_id":       strings.Replace(StatusTopic, "/", "_", -1),
		"unique_id":       strings.Replace(StatusTopic, "/", "_", -1),
		"device":          DEVICE,
	})
	client.Publish(statusTopic, 0, true, statusJson)

	rssiTopic := fmt.Sprintf("%s/%s_sensor_%s/config", topicPrefix, config.MQTT_TOPIC_DEVICE_NAME, RSSI)
	rssiJson, _ := json.Marshal(map[string]interface{}{
		"state_topic":         RSSITopic,
		"icon":                "mdi:signal",
		"device_class":        "signal_strength",
		"state_class":         "measurement",
		"unit_of_measurement": "dBm",
		"entity_category":     "diagnostic",
		"name":                "RSSI",
		"object_id":           strings.Replace(RSSITopic, "/", "_", -1),
		"unique_id":           strings.Replace(RSSITopic, "/", "_", -1),
		"device":              DEVICE,
	})
	client.Publish(rssiTopic, 0, true, rssiJson)

	ndceeTopic := fmt.Sprintf("%s/%s_sensor_%s/config", topicPrefix, config.MQTT_TOPIC_DEVICE_NAME, NormalDirectionCumulativeElectricEnergy)
	ndceeJson, _ := json.Marshal(map[string]interface{}{
		"state_topic":         NormalDirectionCumulativeElectricEnergyTopic,
		"icon":                "mdi:flash",
		"device_class":        "energy",
		"state_class":         "total_increasing",
		"unit_of_measurement": "kWh",
		"entity_category":     "diagnostic",
		"name":                "Normal Direction Cumulative Electric Energy",
		"object_id":           strings.Replace(NormalDirectionCumulativeElectricEnergyTopic, "/", "_", -1),
		"unique_id":           strings.Replace(NormalDirectionCumulativeElectricEnergyTopic, "/", "_", -1),
		"device":              DEVICE,
	})
	client.Publish(ndceeTopic, 0, true, ndceeJson)

	rdceeTopic := fmt.Sprintf("%s/%s_sensor_%s/config", topicPrefix, config.MQTT_TOPIC_DEVICE_NAME, ReverseDirectionCumulativeElectricEnergy)
	rdceeJson, _ := json.Marshal(map[string]interface{}{
		"state_topic":         ReverseDirectionCumulativeElectricEnergyTopic,
		"icon":                "mdi:flash",
		"device_class":        "energy",
		"state_class":         "total_increasing",
		"unit_of_measurement": "kWh",
		"entity_category":     "diagnostic",
		"name":                "Reverse Direction Cumulative Electric Energy",
		"object_id":           strings.Replace(ReverseDirectionCumulativeElectricEnergyTopic, "/", "_", -1),
		"unique_id":           strings.Replace(ReverseDirectionCumulativeElectricEnergyTopic, "/", "_", -1),
		"device":              DEVICE,
	})
	client.Publish(rdceeTopic, 0, true, rdceeJson)

	iepTopic := fmt.Sprintf("%s/%s_sensor_%s/config", topicPrefix, config.MQTT_TOPIC_DEVICE_NAME, InstantaneousElectricPower)
	iepJson, _ := json.Marshal(map[string]interface{}{
		"state_topic":         InstantaneousElectricPowerTopic,
		"icon":                "mdi:flash",
		"device_class":        "power",
		"state_class":         "measurement",
		"unit_of_measurement": "W",
		"entity_category":     "diagnostic",
		"name":                "Instantaneous Electric Power",
		"object_id":           strings.Replace(InstantaneousElectricPowerTopic, "/", "_", -1),
		"unique_id":           strings.Replace(InstantaneousElectricPowerTopic, "/", "_", -1),
		"device":              DEVICE,
	})
	client.Publish(iepTopic, 0, true, iepJson)

	icTopic := fmt.Sprintf("%s/%s_sensor_%s/config", topicPrefix, config.MQTT_TOPIC_DEVICE_NAME, InstantaneousCurrent)
	icJson, _ := json.Marshal(map[string]interface{}{
		"state_topic":         InstantaneousCurrentTopic,
		"icon":                "mdi:current-ac",
		"device_class":        "current",
		"state_class":         "measurement",
		"unit_of_measurement": "A",
		"entity_category":     "diagnostic",
		"name":                "Instantaneous Current",
		"object_id":           strings.Replace(InstantaneousCurrentTopic, "/", "_", -1),
		"unique_id":           strings.Replace(InstantaneousCurrentTopic, "/", "_", -1),
		"device":              DEVICE,
	})
	client.Publish(icTopic, 0, true, icJson)

	icrTopic := fmt.Sprintf("%s/%s_sensor_%s/config", topicPrefix, config.MQTT_TOPIC_DEVICE_NAME, InstantaneousCurrentRPhase)
	icrJson, _ := json.Marshal(map[string]interface{}{
		"state_topic":         InstantaneousCurrentRPhaseTopic,
		"icon":                "mdi:current-ac",
		"device_class":        "current",
		"state_class":         "measurement",
		"unit_of_measurement": "A",
		"entity_category":     "diagnostic",
		"name":                "Instantaneous Current R Phase",
		"object_id":           strings.Replace(InstantaneousCurrentRPhaseTopic, "/", "_", -1),
		"unique_id":           strings.Replace(InstantaneousCurrentRPhaseTopic, "/", "_", -1),
		"device":              DEVICE,
	})
	client.Publish(icrTopic, 0, true, icrJson)

	ictTopic := fmt.Sprintf("%s/%s_sensor_%s/config", topicPrefix, config.MQTT_TOPIC_DEVICE_NAME, InstantaneousCurrentTPhase)
	ictJson, _ := json.Marshal(map[string]interface{}{
		"state_topic":         InstantaneousCurrentTPhaseTopic,
		"icon":                "mdi:current-ac",
		"device_class":        "current",
		"state_class":         "measurement",
		"unit_of_measurement": "A",
		"entity_category":     "diagnostic",
		"name":                "Instantaneous Current T Phase",
		"object_id":           strings.Replace(InstantaneousCurrentTPhaseTopic, "/", "_", -1),
		"unique_id":           strings.Replace(InstantaneousCurrentTPhaseTopic, "/", "_", -1),
		"device":              DEVICE,
	})
	client.Publish(ictTopic, 0, true, ictJson)
}
