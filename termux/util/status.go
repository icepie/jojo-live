package util

import (
	"jojo-live/client"
	"time"

	tm "github.com/eternal-flame-AD/go-termux"
	"github.com/gorilla/websocket"
)

var (
	MainStatus          Status
	LastLightHandleTime time.Time
	LastCallTime        time.Time
	WakeTime            time.Time

	WSConnMap = make(map[string]*websocket.Conn)
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

func GetStatus() Status {

	MainStatus.IsSleep = time.Now().Before(WakeTime)

	// 东八区
	MainStatus.WakeTime = WakeTime.Add(8 * time.Hour).Format("2006-01-02 15:04:05")

	MainStatus.OnlineNum = len(WSConnMap)

	return MainStatus
}

// func GetStatusJson() []byte {

// 	sj, _ := json.Marshal(GetStatus())

// 	return sj
// }
