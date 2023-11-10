package main

import (
	"time"

	mqtt "github.com/mochi-mqtt/server/v2"
)

func testPublish(s *mqtt.Server) {
	for range time.Tick(time.Second * 30) {
		err := s.Publish("test", []byte("test message"), false, 0)
		if err != nil {
			s.Log.Error("server.Publish", "error", err)
			return
		}
		s.Log.Info("issued test message")
	}
}
