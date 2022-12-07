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

	Mpv("https://img.tukuppt.com/newpreview_music/09/00/25/5c89106abeedd53089.mp3")

	client.GetMiLightStatus()

	if stat, err := tm.BatteryStatus(); err != nil {
		panic(err)
	} else {
		log.Println(stat.Health, stat.Percentage, stat.Status, stat.Temperature)
	}

	log.Println(client.GetMaAcIndoorTemperature())
}
