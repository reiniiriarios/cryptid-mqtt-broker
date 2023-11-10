package main

import (
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

	authRules := &auth.Ledger{
		Auth: auth.AuthRules{ // Auth disallows all by default
			// @todo add actual auth
			{Username: "cryptid", Password: "public", Allow: true},
			{Remote: "127.0.0.1:*", Allow: true},
			{Remote: "localhost:*", Allow: true},
			{Remote: "172.16.0.*:*", Allow: true},
		},
		ACL: auth.ACLRules{ // ACL allows all by default
			{Remote: "127.0.0.1:*"}, // local superuser allow all
			{
				Username: "cryptid", Filters: auth.Filters{
					"test":      auth.ReadWrite,
					"hello":     auth.ReadWrite,
					"cryptid/#": auth.ReadWrite,
					"updates/#": auth.WriteOnly,
					"status/#":  auth.ReadOnly,
				},
			},
			{
				// Otherwise, no clients have publishing permissions
				Filters: auth.Filters{
					"#":         auth.Deny,
					"updates/#": auth.Deny,
					"status/#":  auth.ReadOnly,
				},
			},
		},
	}

	// Create the new MQTT Server.
	server := mqtt.New(nil)

	// Set logging.
	level := new(slog.LevelVar)
	server.Log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: level,
	}))
	level.Set(slog.LevelDebug)

	// Add hook.
	err := server.AddHook(new(auth.Hook), &auth.Options{
		Ledger: authRules,
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
	_ = server.Close()
	println(string(C_R) + "MQTT Broker Closed" + string(C_X))
}
