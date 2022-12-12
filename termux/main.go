package main

import (
	"encoding/json"
	"jojo-live/client"
	"jojo-live/util"
	"jojo-live/ws"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"github.com/juju/ratelimit"
)

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

func GetStatus() util.Status {

	util.MainStatus.IsSleep = time.Now().Before(util.WakeTime)

	// 东八区
	util.MainStatus.WakeTime = util.WakeTime.Add(8 * time.Hour).Format("2006-01-02 15:04:05")

	util.MainStatus.OnlineNum = ws.WsHub.Len()

	return util.MainStatus
}

func Broadcast() {
	for {
		data, _ := json.Marshal(ws.WsMessage{
			Type: "status",
			Data: GetStatus(),
		})
		ws.WsHub.Broadcast(data)
		time.Sleep(3 * time.Second)
	}
}

func main() {

	// go util.UpdateIndoorTemperature()
	go util.UpdateOtherStatus()

	go Broadcast()

	// hub := ws.NewHub()

	// gin
	go ws.WsHub.Run()

	r := gin.Default()

	r.GET("/ws", gin.WrapF(func(w http.ResponseWriter, r *http.Request) {
		ws.ServeWs(ws.WsHub, w, r)
	}))

	// CORS middleware
	r.Use(cors.Default())

	r.Use(RateLimitMiddleware(1*time.Second, 1200, 1))

	r.GET("/status", func(c *gin.Context) {

		c.JSON(200, GetStatus())
	})

	r.GET("/light/on", func(c *gin.Context) {

		// 判断是否到睡醒时间
		if time.Now().Before(util.WakeTime) {
			c.JSON(403, "JOJO正(要)睡觉哦!")
			return
		}

		if time.Since(util.LastLightHandleTime) < 5*time.Second {
			c.JSON(403, "操作太频繁了")
			return
		}

		util.LastLightHandleTime = time.Now()

		err := client.SetMiLightPower(true)
		if err != nil {
			c.JSON(500, err.Error())
			return
		}

		c.JSON(200, "开灯成功")
	})

	r.GET("sleep", func(c *gin.Context) {

		// 判断是否到睡醒时间
		if time.Now().Before(util.WakeTime) {
			c.JSON(200, "已经是睡眠模式啦")
			return
		}

		util.WakeTime = time.Now()

		hour := c.Query("hour")

		if hour != "" {
			h, err := strconv.Atoi(hour)
			if err != nil {
				c.JSON(500, err.Error())
				return
			}
			util.WakeTime = time.Now().Add(time.Duration(h) * time.Hour)
			if h > 12 && h < 1 {
				c.JSON(500, "睡觉时长不合法")
			}

		} else {
			util.WakeTime = time.Now().Add(1 * time.Hour)
		}

		c.JSON(200, "准备睡")
	})

	r.GET("/light/off", func(c *gin.Context) {

		// 判断是否到睡醒时间
		if time.Now().Before(util.WakeTime) {
			c.JSON(403, "JOJO正(要)睡觉哦!")
			return
		}

		if time.Since(util.LastLightHandleTime) < 5*time.Second {
			c.JSON(403, "操作太频繁了")
			return
		}

		util.LastLightHandleTime = time.Now()

		err := client.SetMiLightPower(false)
		if err != nil {
			c.JSON(500, err.Error())
			return
		}

		c.JSON(200, "关灯成功")
	})

	r.GET("/call", func(c *gin.Context) {

		// 判断是否到睡醒时间
		if time.Now().Before(util.WakeTime) {
			c.JSON(403, "JOJO正(要)睡觉哦!")
			return
		}

		// 判断是否在 23:00 ~ 8:00
		now := time.Now().Add(8 * time.Hour)
		if now.Hour() > 23 || now.Hour() < 8 {
			c.JSON(403, "深夜不能呼叫JOJO!")
			return
		}

		if time.Since(util.LastCallTime) < 5*time.Second {
			c.JSON(403, "呼叫太频繁了")
			return
		}

		util.LastCallTime = time.Now()

		err := util.Mpv("https://img.tukuppt.com/newpreview_music/09/00/25/5c89106abeedd53089.mp3")
		if err != nil {
			c.JSON(500, err.Error())
			return
		}

		c.JSON(200, "已呼叫")
	})

	r.Run(":8080")
}
