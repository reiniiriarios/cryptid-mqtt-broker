package main

import (
	"fmt"
	"strconv"
	"time"

	"github.com/fatih/color"
	mqtt "github.com/mochi-mqtt/server/v2"
	"github.com/mochi-mqtt/server/v2/packets"
)

const DISPLAY_HOT_THRESHOLD = 25 // degrees celcius
const DISPLAY_WET_THRESHOLD = 75 // percentage relative humidity

func receiveDisplayTemperature(cl *mqtt.Client, sub packets.Subscription, pk packets.Packet) {
	temp, err := strconv.ParseFloat(string(pk.Payload), 32)
	if err != nil {
		println(
			time.Now().Format(TIMESTAMP_FORMAT),
			color.RedString("Unable to read Display temperature:"), pk.Payload,
		)
	}
	if temp > DISPLAY_HOT_THRESHOLD {
		println(
			time.Now().Format(TIMESTAMP_FORMAT),
			color.RedString("🔥 Display is HOT! "+fmt.Sprintf("%.2f", temp)+"°C"),
		)
	}
}

func receiveDisplayHumidity(cl *mqtt.Client, sub packets.Subscription, pk packets.Packet) {
	humidity, err := strconv.ParseFloat(string(pk.Payload), 32)
	if err != nil {
		println(
			time.Now().Format(TIMESTAMP_FORMAT),
			color.RedString("Unable to read Display humidity:"), pk.Payload,
		)
	}
	if humidity > DISPLAY_WET_THRESHOLD {
		println(
			time.Now().Format(TIMESTAMP_FORMAT),
			color.RedString("💦 Display is WET! "+fmt.Sprintf("%.2f", humidity)+"% rH"),
		)
	}
}
