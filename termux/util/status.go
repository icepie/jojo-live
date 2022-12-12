package util

import (
	"jojo-live/client"
	"time"

	tm "github.com/eternal-flame-AD/go-termux"
)

var (
	MainStatus          Status
	LastLightHandleTime time.Time
	LastCallTime        time.Time
	WakeTime            time.Time
)

type Battery struct {
	BatteryPercentage  int
	BatterISCharging   bool
	BatteryHealth      string
	BatteryTemperature float64
}

type Status struct {
	Battery           Battery
	OnlineNum         int
	IsSleep           bool
	WakeTime          string
	LightPower        bool
	IndoorTemperature float64
}

func UpdateIndoorTemperature() {
	for {
		MainStatus.IndoorTemperature, _ = client.GetMaAcIndoorTemperature()
		time.Sleep(10 * time.Second)
	}
}

func UpdateOtherStatus() {
	for {
		// var status Status

		lightStatus, err := client.GetMiLightStatus()
		if err == nil {
			MainStatus.LightPower = lightStatus.Result[0].(string) == "on"
		}

		if stat, err := tm.BatteryStatus(); err == nil {
			MainStatus.Battery.BatteryPercentage = stat.Percentage
			MainStatus.Battery.BatterISCharging = stat.Status != "DISCHARGING"
			MainStatus.Battery.BatteryHealth = stat.Health
			MainStatus.Battery.BatteryTemperature = stat.Temperature
		}

		time.Sleep(5 * time.Second)
	}
}
