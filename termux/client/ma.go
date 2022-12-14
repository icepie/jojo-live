package client

import (
	"errors"
	ma "jojo-live/midea-ac"
)

func GetMaAcIndoorTemperature() (temperature float64, err error) {

	MideaAcDevice := ma.NewDevice("192.168.2.217", 162727724063545, "6444")

	isAuth := MideaAcDevice.AuthenticateV3("0952ccb99ed54264bb5774e68b126acb24f7f61aa7c143b5a06af9601057b7ac", "6EBCEACB480D0F68395A05F2D9E96299379477CF986950E770103FADB29232E806619956B3BCA2C9660EAA86C5DB7C0891587F7613CDD0A12B14A86B5E252BF2")
	if !isAuth {
		panic("Authentication failed")
	}

	MideaAc := ma.NewAirConditioningDevice(MideaAcDevice)

	MideaAc.Refresh()

	temperature = MideaAc.IndoorTemperature()
	if temperature == 0 {
		err = errors.New("indoor temperature is 0")
	}

	return
}
