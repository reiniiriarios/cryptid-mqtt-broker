package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	mqtt "github.com/mochi-mqtt/server/v2"
	"github.com/mochi-mqtt/server/v2/hooks/auth"
	"github.com/mochi-mqtt/server/v2/listeners"
)

const MQTT_PORT = ":1883"

const C_X = "\033[0m"
const C_R = "\033[31m"
const C_G = "\033[32m"
const C_Y = "\033[33m"
const C_B = "\033[34m"
const C_P = "\033[35m"
const C_C = "\033[36m"
const C_W = "\033[37m"

func main() {
	println(string(C_G) + "Starting MQTT Broker...")
	// Display local address
	ip := GetOutboundIP()
	println(string(C_P)+"Listening on:", ip.String()+MQTT_PORT)
	print(string(C_X))

	// Create signals channel to run server until interrupted
	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigs
		done <- true
	}()

	// Create the new MQTT Server.
	server := mqtt.New(nil)

	// Allow all connections.
	_ = server.AddHook(new(auth.AllowHook), nil)

	// Create a TCP listener on a standard port.
	tcp := listeners.NewTCP("t1", MQTT_PORT, nil)
	err := server.AddListener(tcp)
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		err := server.Serve()
		if err != nil {
			log.Fatal(err)
		}
	}()

	// Run server until interrupted
	<-done

	// Cleanup
	println()
	println(string(C_R) + "MQTT Broker Closed" + string(C_X))
}
