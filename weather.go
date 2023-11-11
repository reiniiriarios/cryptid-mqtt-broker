package main

import (
	"fmt"
	"strconv"
	"time"

	"github.com/fatih/color"
	mqtt "github.com/mochi-mqtt/server/v2"
	"github.com/mochi-mqtt/server/v2/packets"
)

func receiveError(topic string, payload []byte) {
	println(
		time.Now().Format(TIMESTAMP_FORMAT),
		color.RedString("Unable to read `"+topic+"`"),
		payload,
	)
}

func logReceipt(s string) {
	println(time.Now().Format(TIMESTAMP_FORMAT), s)
}

func receiveWeatherTemperature(cl *mqtt.Client, sub packets.Subscription, pk packets.Packet) {
	temp, err := strconv.ParseFloat(string(pk.Payload), 32)
	if err != nil {
		receiveError(pk.TopicName, pk.Payload)
		return
	}
	logReceipt("It's currently " + fmt.Sprintf("%.2f", temp) + "°C")
}

func receiveWeatherHumidity(cl *mqtt.Client, sub packets.Subscription, pk packets.Packet) {
	humidity, err := strconv.ParseInt(string(pk.Payload), 10, 8)
	if err != nil {
		receiveError(pk.TopicName, pk.Payload)
		return
	}
	logReceipt("The humidity is " + fmt.Sprint(humidity) + "% rH")
}

func receiveWeatherFeelsLike(cl *mqtt.Client, sub packets.Subscription, pk packets.Packet) {
	temp, err := strconv.ParseFloat(string(pk.Payload), 32)
	if err != nil {
		receiveError(pk.TopicName, pk.Payload)
		return
	}
	logReceipt("It feels like " + fmt.Sprintf("%.2f", temp) + "°C")
}

func receiveWeatherCondition(cl *mqtt.Client, sub packets.Subscription, pk packets.Packet) {
	logReceipt("The current condition is: " + string(pk.Payload))
}
