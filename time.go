package main

import (
	"strconv"
	"time"

	mqtt "github.com/mochi-mqtt/server/v2"
)

func timePublish(s *mqtt.Server) {
	var ts []byte
	for range time.Tick(time.Second * 60) {
		ts = []byte(strconv.FormatInt(time.Now().Unix(), 10))
		err := s.Publish("current_time", ts, false, 0)
		if err != nil {
			s.Log.Error("timePublish error", err)
			return
		}
	}
}
