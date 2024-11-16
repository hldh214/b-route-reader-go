package main

import (
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/yakumo-saki/b-route-reader-go/src/echonet"
	"github.com/yakumo-saki/b-route-reader-go/src/ha"
	"os"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/yakumo-saki/b-route-reader-go/src/bp35a1"
	"github.com/yakumo-saki/b-route-reader-go/src/config"
	"github.com/yakumo-saki/b-route-reader-go/src/global"
	"github.com/yakumo-saki/b-route-reader-go/src/logger"
)

var exitcode = 0

func main() {
	ret := run()
	os.Exit(ret)
}

func run() int {
	config.Initialize()
	logger.Initiallize()

	log.Info().Msgf("Version %s. Build %s", global.Version, global.GitBuild)
	log.Info().Msgf("%s", global.Url)
	err := bp35a1.Connect()
	if err != nil {
		log.Err(err).Msg("Serial port open error. Exiting.")
		exitcode = 1
		goto EXIT
	}

	err = runWithSerialPort()
	if err != nil {
		exitcode = 1
		log.Err(err).Msg("ERR")
	}

EXIT:
	err = bp35a1.Close()
	if err != nil {
		log.Err(err).Msg("Error occured in close connection. do nothing.")
	}

	if exitcode == 0 {
		log.Info().Msg("Normal end.")
	}

	return exitcode
}

func runWithSerialPort() error {
	var err error
	var mqttClient mqtt.Client

	// MQTT for HA
	opts := mqtt.NewClientOptions()
	opts.AddBroker(config.MQTT_BROKER)
	opts.SetClientID("b-route-reader")
	opts.SetUsername(config.MQTT_USERNAME)
	opts.SetPassword(config.MQTT_PASSWORD)
	opts.SetAutoReconnect(true)
	opts.SetOnConnectHandler(func(client mqtt.Client) {
		log.Info().Msg("MQTT connected")
		ha.InitTopic(client)
	})
	mqttClient = mqtt.NewClient(opts)
	if token := mqttClient.Connect(); token.Wait() && token.Error() != nil {
		return fmt.Errorf("error occured while connecting to mqtt broker: %w", token.Error())
	}

	mqttClient.Publish(ha.StatusTopic, 0, true, ha.StatusConnecting)
	err = bp35a1.StartConnection()
	if err != nil {
		mqttClient.Publish(ha.StatusTopic, 0, true, ha.StatusError)
		return fmt.Errorf("test connection failed: %w", err)
	}

	initResult, err := bp35a1.InitializeBrouteConnection()
	if err != nil {
		mqttClient.Publish(ha.StatusTopic, 0, true, ha.StatusError)
		log.Err(err).Msg(". Exiting.")
		return fmt.Errorf("cannot initialize B-route connection: %w", err)
	}
	mqttClient.Publish(ha.RSSITopic, 0, true, fmt.Sprintf("%d", initResult.Rssi))
	mqttClient.Publish(ha.StatusTopic, 0, true, ha.StatusConnected)

	// echonet start
	err = bp35a1.GetSmartMeterInitialData(initResult.Ipv6)
	if err != nil {
		mqttClient.Publish(ha.StatusTopic, 0, true, ha.StatusError)
		return fmt.Errorf("error occured while initializing echonet lite: %w", err)
	}

	log.Info().Msg("Starting main loop")

	// TODO シグナルハンドリング

	nowTimer := time.NewTimer(config.NOW_CONSUMPTION_WAIT)
	totalTimer := time.NewTimer(config.TOTAL_CONSUMPTION_WAIT)

	for {
		select {
		case <-nowTimer.C:
			ret, err := bp35a1.GetNowConsumptionData(initResult.Ipv6)
			if err != nil {
				mqttClient.Publish(ha.StatusTopic, 0, true, ha.StatusError)
				return fmt.Errorf("error occured while getting consumption: %w", err)
			}

			nowTimer = time.NewTimer(config.NOW_CONSUMPTION_WAIT)

			log.Info().Msgf("Smartmeter Response: %v", ret)

			mqttClient.Publish(ha.StatusTopic, 0, true, ha.StatusRunning)
			mqttClient.Publish(ha.InstantaneousElectricPowerTopic, 0, true, ret[fmt.Sprintf("%02X", echonet.P_NOW_DENRYOKU)].String())
			mqttClient.Publish(ha.InstantaneousCurrentTopic, 0, true, ret[fmt.Sprintf("%02X", echonet.P_NOW_DENRYUU)].String())
			mqttClient.Publish(ha.InstantaneousCurrentRPhaseTopic, 0, true, ret[fmt.Sprintf("%02X_Rphase", echonet.P_NOW_DENRYUU)].String())
			mqttClient.Publish(ha.InstantaneousCurrentTPhaseTopic, 0, true, ret[fmt.Sprintf("%02X_Tphase", echonet.P_NOW_DENRYUU)].String())
		case <-totalTimer.C:
			ret, err := bp35a1.GetDeltaConsumptionData(initResult.Ipv6)
			if err != nil {
				mqttClient.Publish(ha.StatusTopic, 0, true, ha.StatusError)
				return fmt.Errorf("error occured while getting delta consumption: %w", err)
			}

			totalTimer = time.NewTimer(config.TOTAL_CONSUMPTION_WAIT)

			log.Info().Msgf("Smartmeter Response: %v", ret)

			nd := ret[fmt.Sprintf("%02X", echonet.P_DELTA_DENRYOKU)]
			rd := ret[fmt.Sprintf("%02X", echonet.P_DELTA_DENRYOKU_R)]

			if nd.String() == "0" || rd.String() == "0" {
				log.Info().Msg("No data. Skip.")
				continue
			}

			mqttClient.Publish(ha.StatusTopic, 0, true, ha.StatusRunning)
			mqttClient.Publish(ha.NormalDirectionCumulativeElectricEnergyTopic, 0, true, nd.String())
			mqttClient.Publish(ha.ReverseDirectionCumulativeElectricEnergyTopic, 0, true, rd.String())
		}
	}
}
