package ha

import (
	"fmt"
	"github.com/yakumo-saki/b-route-reader-go/src/config"
)

const RSSI = "RSSI"
const NormalDirectionCumulativeElectricEnergy = "NormalDirectionCumulativeElectricEnergy"
const ReverseDirectionCumulativeElectricEnergy = "ReverseDirectionCumulativeElectricEnergy"
const InstantaneousElectricPower = "InstantaneousElectricPower"
const InstantaneousCurrent = "InstantaneousCurrent"
const InstantaneousCurrentRPhase = "InstantaneousCurrentRPhase"
const InstantaneousCurrentTPhase = "InstantaneousCurrentTPhase"

var RSSITopic = fmt.Sprintf("%s/%s/%s", config.MQTT_TOPIC_PREFIX, config.MQTT_TOPIC_DEVICE_NAME, RSSI)
var NormalDirectionCumulativeElectricEnergyTopic = fmt.Sprintf("%s/%s/%s", config.MQTT_TOPIC_PREFIX, config.MQTT_TOPIC_DEVICE_NAME, NormalDirectionCumulativeElectricEnergy)
var ReverseDirectionCumulativeElectricEnergyTopic = fmt.Sprintf("%s/%s/%s", config.MQTT_TOPIC_PREFIX, config.MQTT_TOPIC_DEVICE_NAME, ReverseDirectionCumulativeElectricEnergy)
var InstantaneousElectricPowerTopic = fmt.Sprintf("%s/%s/%s", config.MQTT_TOPIC_PREFIX, config.MQTT_TOPIC_DEVICE_NAME, InstantaneousElectricPower)
var InstantaneousCurrentTopic = fmt.Sprintf("%s/%s/%s", config.MQTT_TOPIC_PREFIX, config.MQTT_TOPIC_DEVICE_NAME, InstantaneousCurrent)
var InstantaneousCurrentRPhaseTopic = fmt.Sprintf("%s/%s/%s", config.MQTT_TOPIC_PREFIX, config.MQTT_TOPIC_DEVICE_NAME, InstantaneousCurrentRPhase)
var InstantaneousCurrentTPhaseTopic = fmt.Sprintf("%s/%s/%s", config.MQTT_TOPIC_PREFIX, config.MQTT_TOPIC_DEVICE_NAME, InstantaneousCurrentTPhase)
