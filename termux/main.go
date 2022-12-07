package main

import (
	"log"
	"os/exec"

	"jojo-live/client"
	ma "jojo-live/midea-ac"

	tm "github.com/eternal-flame-AD/go-termux"
)

func Mpv(url string) error {
	mpv := exec.Command("mpv", url)
	err := mpv.Start()
	return err
}

func main() {

	// c := xiaomiio.NewXiaoMiio("", "")
	// err := c.Login()
	// if err != nil {
	// 	panic(err)
	// }

	// deviceList, err := c.GetDevices()
	// if err != nil {
	// 	panic(err)
	// }

	// log.Println(deviceList)

	// for _, device := range deviceList {

	// 	if device.Model == "yeelink.light.lamp1" {
	// 		log.Println(device.Did)
	// 	}
	// }

	Mpv("https://img.tukuppt.com/newpreview_music/09/00/25/5c89106abeedd53089.mp3")

	client.GetMiLightStatus()

	if stat, err := tm.BatteryStatus(); err != nil {
		panic(err)
	} else {
		log.Println(stat.Health, stat.Percentage, stat.Status, stat.Temperature)
	}

	device := ma.NewDevice("192.168.2.217", 162727724063545, "6444")

	isAuth := device.AuthenticateV3("0952ccb99ed54264bb5774e68b126acb24f7f61aa7c143b5a06af9601057b7ac", "6EBCEACB480D0F68395A05F2D9E96299379477CF986950E770103FADB29232E806619956B3BCA2C9660EAA86C5DB7C0891587F7613CDD0A12B14A86B5E252BF2")
	if !isAuth {
		panic("Authentication failed")
	}

	ac := ma.NewAirConditioningDevice(device)

	log.Println(ac.IndoorTemperature(), ac.Outdoortemperature(), ac.PowerState(), ac.SwingMode(), ac.FanSpeed(), ac.TurboMode())
}
