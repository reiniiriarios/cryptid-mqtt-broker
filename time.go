package main

import (
	"strconv"
	"time"

	mqtt "github.com/mochi-mqtt/server/v2"
)

const UPDATE_TIME_INTERVAL = 60
const TIMEZONE = "America/Los_Angeles"

// Publishes a unix timestamp _in the local timezone_.
// Assume here that this will be digested by embedded
// systems that cannot as easily convert from UTC.
func timePublish(s *mqtt.Server) {
	for range time.Tick(time.Second * UPDATE_TIME_INTERVAL) {
		location, err := time.LoadLocation(TIMEZONE)
		if err != nil {
			s.Log.Error("timePublish error", err)
			continue
		}
		_, tzOffset := time.Now().In(location).Zone()
		ts := time.Now().Unix() + int64(tzOffset)

		payload := []byte(strconv.FormatInt(ts, 10))
		err = s.Publish("current_time", payload, false, 0)
		if err != nil {
			s.Log.Error("timePublish error", err)
			return
		}
	}
}
