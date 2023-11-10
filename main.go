package main

import (
	"flag"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	mqtt "github.com/mochi-mqtt/server/v2"
	"github.com/mochi-mqtt/server/v2/hooks/auth"
	"github.com/mochi-mqtt/server/v2/listeners"
)

const MQTT_PORT = ":1883"

func main() {
	// Create signals channel to run server until interrupted
	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigs
		done <- true
	}()

	// Create the new MQTT Server.
	server := mqtt.New(&mqtt.Options{
		InlineClient: true,
	})

	// Set logging.
	server.Log = slog.New(NewPrettyLogHandler(os.Stdout, PrettyLogHandlerOptions{
		SlogOpts: slog.HandlerOptions{
			Level: slog.LevelDebug,
		},
	}))

	server.Log.Info("Starting MQTT Broker...")
	server.Log.Info("Listening on: " + GetOutboundIP().String() + MQTT_PORT)

	// Get ledger from yaml file
	authFilePath := flag.String("path", "auth.yaml", "path to data auth file")
	flag.Parse()
	authData, err := os.ReadFile(*authFilePath)
	if err != nil {
		log.Fatal(err)
	}

	// Add hook with auth options.
	err = server.AddHook(new(auth.Hook), &auth.Options{
		Data: authData,
	})
	if err != nil {
		log.Fatal(err)
	}

	// Create a TCP listener on a standard port.
	tcp := listeners.NewTCP("t1", MQTT_PORT, nil)
	err = server.AddListener(tcp)
	if err != nil {
		log.Fatal(err)
	}

	// Run.
	go func() {
		err := server.Serve()
		if err != nil {
			log.Fatal(err)
		}
	}()

	// Test direct interactions.
	go testPublish(server)

	// Display
	_ = server.Subscribe("display/temperature", 100, receiveDisplayTemperature)
	_ = server.Subscribe("display/humidity", 101, receiveDisplayHumidity)

	// Weather
	_ = server.Subscribe("weather/temperature", 200, receiveWeatherTemperature)

	// Run server until interrupted
	<-done

	// Cleanup
	server.Log.Warn("Closing broker...")
	_ = server.Close()
	server.Log.Warn("MQTT Broker Closed")
}
