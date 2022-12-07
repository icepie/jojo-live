package main

import (
	"log"
	"os/exec"

	"jojo-live/client"

	tm "github.com/eternal-flame-AD/go-termux"
)

func Mpv(url string) error {
	mpv := exec.Command("mpv", url)
	err := mpv.Start()
	return err
}

func main() {

	client.SetMiLightPower(true)

	client.SetMiLightPower(false)

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

	log.Println(client.MideaAc.IndoorTemperature(), client.MideaAc.Outdoortemperature(), client.MideaAc.PowerState(), client.MideaAc.SwingMode(), client.MideaAc.FanSpeed(), client.MideaAc.TurboMode())
}
