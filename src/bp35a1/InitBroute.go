package bp35a1

import (
	"fmt"
	"strconv"

	"github.com/rs/zerolog/log"
	"github.com/yakumo-saki/b-route-reader-go/src/config"
)

// アクティブスキャン〜Bルート接続完了までを実行する
// @return スマートメーターのipv6アドレス
func InitializeBrouteConnection() (string, error) {
	isAscii, err := isAsciiMode()
	if err != nil {
		return "", err
	}
	if !isAscii {
		// WOPT 1
		log.Warn().Msg("WOPT 1 is not implemented. maybe not working.")
		return "", fmt.Errorf("command 'WOPT 1' is not implemented")
	}

	err = setupIdAndPassword()
	if err != nil {
		return "", err
	}

	sm, err := searchSmartMeter()
	if err != nil {
		return "", err
	}

	log.Info().Msgf("Found smartmeter %s", sm)
	lqi, _ := strconv.Atoi(sm.LQI)
	log.Info().Msgf("LQI: %d", lqi)
	// rssi = 0.275 * lqi - 104.27
	log.Info().Msgf("RSSI: %f", 0.275*float64(lqi)-104.27)

	ipv6, err := convertPanIdToIpv6(sm.Addr)
	if err != nil {
		return "", err
	}

	log.Info().Msgf("Smartmeter address is %s", ipv6)

	err = setBroutePanChannel(sm.Channel)
	if err != nil {
		return "", err
	}

	err = setBroutePanId(sm.PanId)
	if err != nil {
		return "", err
	}

	err = startPaCAuthentication(ipv6)
	if err != nil {
		return "", err
	}

	log.Info().Msgf("PAN authentication done.")

	return ipv6, nil
}

func setupIdAndPassword() error {
	// ID PWD
	err := setBrouteId(config.B_ROUTE_ID)
	if err != nil {
		return err
	}

	err = setBroutePassword(config.B_ROUTE_PASSWORD)
	if err != nil {
		return err
	}

	return nil
}

func searchSmartMeter() (SmartMeter, error) {

	log.Info().Msg("Active scan start. this will take some moment.")

	sm, err := activeScan()
	if err != nil {
		return sm, err
	}

	return sm, err
}
