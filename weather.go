package main

import (
	"fmt"
	"strconv"
	"time"

	"github.com/fatih/color"
	mqtt "github.com/mochi-mqtt/server/v2"
	"github.com/mochi-mqtt/server/v2/packets"
)

func receiveWeatherTemperature(cl *mqtt.Client, sub packets.Subscription, pk packets.Packet) {
	temp, err := strconv.ParseFloat(string(pk.Payload), 32)
	if err != nil {
		println(
			time.Now().Format(TIMESTAMP_FORMAT),
			color.RedString("Unable to read weather temp."), pk.Payload,
		)
	}
	println(
		time.Now().Format(TIMESTAMP_FORMAT),
		"It's currently "+fmt.Sprintf("%.2f", temp)+"°C",
	)
}
