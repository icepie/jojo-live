package main

import (
	"jojo-live/client"
	"net/http"
	"strconv"
	"time"

	tm "github.com/eternal-flame-AD/go-termux"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"github.com/juju/ratelimit"
)

var (
	status              Status
	lastLightHandleTime time.Time
	lastCallTime        time.Time
	wakeTime            time.Time
)

type Battery struct {
	BatteryPercentage  int
	BatterISCharging   bool
	BatteryHealth      string
	BatteryTemperature float64
}

type Status struct {
	Battery           Battery
	IsSleep           bool
	WakeTime          string
	LightPower        bool
	IndoorTemperature float64
}

func updateIndoorTemperature() {
	for {
		status.IndoorTemperature, _ = client.GetMaAcIndoorTemperature()
		time.Sleep(10 * time.Second)
	}
}

func updateOtherStatus() {
	for {
		// var status Status

		lightStatus, err := client.GetMiLightStatus()
		if err == nil {
			status.LightPower = lightStatus.Result[0].(string) == "on"
		}

		if stat, err := tm.BatteryStatus(); err == nil {
			status.Battery.BatteryPercentage = stat.Percentage
			status.Battery.BatterISCharging = stat.Status != "DISCHARGING"
			status.Battery.BatteryHealth = stat.Health
			status.Battery.BatteryTemperature = stat.Temperature
		}

		time.Sleep(5 * time.Second)
	}
}

func RateLimitMiddleware(fillInterval time.Duration, cap, quantum int64) gin.HandlerFunc {
	bucket := ratelimit.NewBucketWithQuantum(fillInterval, cap, quantum)
	return func(c *gin.Context) {
		if bucket.TakeAvailable(1) < 1 {
			c.String(http.StatusForbidden, "操作太频繁了...")
			c.Abort()
			return
		}
		c.Next()
	}
}

func main() {

	go updateIndoorTemperature()
	go updateOtherStatus()

	// gin

	r := gin.Default()

	// CORS middleware
	r.Use(cors.Default())

	r.Use(RateLimitMiddleware(1*time.Second, 1200, 1))

	r.GET("/status", func(c *gin.Context) {

		status.IsSleep = time.Now().Before(wakeTime)

		// 东八区
		status.WakeTime = wakeTime.Add(8 * time.Hour).Format("2006-01-02 15:04:05")

		c.JSON(200, status)
	})

	r.GET("/light/on", func(c *gin.Context) {

		// 判断是否到睡醒时间
		if time.Now().Before(wakeTime) {
			c.JSON(403, "JOJO正(要)睡觉哦!")
			return
		}

		if time.Since(lastLightHandleTime) < 5*time.Second {
			c.JSON(403, "操作太频繁了")
			return
		}

		lastLightHandleTime = time.Now()

		err := client.SetMiLightPower(true)
		if err != nil {
			c.JSON(500, err.Error())
			return
		}

		c.JSON(200, "开灯成功")
	})

	r.GET("sleep", func(c *gin.Context) {

		// 判断是否到睡醒时间
		if time.Now().Before(wakeTime) {
			c.JSON(200, "已经是睡眠模式啦")
			return
		}

		wakeTime = time.Now()

		hour := c.Query("hour")

		if hour != "" {
			h, err := strconv.Atoi(hour)
			if err != nil {
				c.JSON(500, err.Error())
				return
			}
			wakeTime = time.Now().Add(time.Duration(h) * time.Hour)
			if h > 12 && h < 1 {
				c.JSON(500, "睡觉时长不合法")
			}

		} else {
			wakeTime = time.Now().Add(1 * time.Hour)
		}

		c.JSON(200, "准备睡")
	})

	r.GET("/light/off", func(c *gin.Context) {

		// 判断是否到睡醒时间
		if time.Now().Before(wakeTime) {
			c.JSON(403, "JOJO")
			return
		}

		if time.Since(lastLightHandleTime) < 5*time.Second {
			c.JSON(403, "操作太频繁了")
			return
		}

		lastLightHandleTime = time.Now()

		err := client.SetMiLightPower(false)
		if err != nil {
			c.JSON(500, err.Error())
			return
		}

		c.JSON(200, "关灯成功")
	})

	r.GET("/call", func(c *gin.Context) {

		// 判断是否到睡醒时间
		if time.Now().Before(wakeTime) {
			c.JSON(403, "JOJO正(要)睡觉哦!")
			return
		}

		if time.Since(lastCallTime) < 5*time.Second {
			c.JSON(403, "呼叫太频繁了")
			return
		}

		lastCallTime = time.Now()

		err := Mpv("https://img.tukuppt.com/newpreview_music/09/00/25/5c89106abeedd53089.mp3")
		if err != nil {
			c.JSON(500, err.Error())
			return
		}

		c.JSON(200, "已呼叫")
	})

	r.Run(":8080")
}
