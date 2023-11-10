package main

import (
	"fmt"
	"time"

	mqtt "github.com/mochi-mqtt/server/v2"
	"github.com/mochi-mqtt/server/v2/packets"
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

func testReceive(cl *mqtt.Client, sub packets.Subscription, pk packets.Packet) {
	println(
		"Inline receive: client="+cl.ID,
		"sub="+fmt.Sprint(sub.Identifier),
		"topic="+pk.TopicName,
		"payload="+string(pk.Payload),
	)
}
