package main

import (
	"jojo-live/client"
	"time"

	tm "github.com/eternal-flame-AD/go-termux"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var (
	indoorTemperature float64
)

type Battery struct {
	BatteryPercentage  int
	BatterISCharging   bool
	BatteryHealth      string
	BatteryTemperature float64
}

type Status struct {
	Battery           Battery
	LightPower        bool
	IndoorTemperature float64
}

func updateIndoorTemperature() {
	for {
		indoorTemperature, _ = client.GetMaAcIndoorTemperature()
		time.Sleep(10 * time.Second)
	}
}

func main() {

	go updateIndoorTemperature()

	// gin

	r := gin.Default()

	// CORS middleware
	r.Use(cors.Default())

	r.GET("/status", func(c *gin.Context) {
		var status Status

		lightStatus, err := client.GetMiLightStatus()
		if err == nil {
			status.LightPower = lightStatus.Result[0].(string) == "on"
		}

		status.IndoorTemperature = indoorTemperature

		if stat, err := tm.BatteryStatus(); err == nil {
			status.Battery.BatteryPercentage = stat.Percentage
			status.Battery.BatterISCharging = stat.Status == "Charging"
			status.Battery.BatteryHealth = stat.Health
			status.Battery.BatteryTemperature = stat.Temperature
		}

		c.JSON(200, status)
	})

	r.GET("/light/on", func(c *gin.Context) {
		client.SetMiLightPower(true)
		c.JSON(200, "ok")
	})

	r.GET("/light/off", func(c *gin.Context) {
		client.SetMiLightPower(false)
		c.JSON(200, "ok")
	})

	r.GET("/call", func(c *gin.Context) {
		Mpv("https://img.tukuppt.com/newpreview_music/09/00/25/5c89106abeedd53089.mp3")
		c.JSON(200, "ok")
	})

	r.Run(":8080")
}
