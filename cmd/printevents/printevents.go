package main

import (
	"github.com/crenz/locationengine"
	log "github.com/Sirupsen/logrus"
	"github.com/mattn/go-colorable"
)

const API_KEY = "your_api_key"
const PLACE = "your_place_uuid"

func telemetryCallback(receiver string, item[] locationengine.Item) {
	log.Infof("[Telemetry] %s: %+v\n", receiver, item)
}

func eventCallback(evt locationengine.Event, receiver string, item locationengine.Item) {
	log.Infof("[Event %i] %s: %+v\n", evt, receiver, item)
}

func init() {
	// Set to log.DebugLevel to see each individual MQTT message received
	log.SetLevel(log.InfoLevel)
	log.SetOutput(colorable.NewColorableStdout())
}

func main() {
	log.Infoln("Connecting to Kontakt.io broker")
	l := locationengine.New()
	l.RegisterEventCallback(locationengine.EvtItemAppeared, eventCallback)
	l.RegisterEventCallback(locationengine.EvtItemDisappeared, eventCallback)
	l.RegisterEventCallback(locationengine.EvtItemProximityChange, eventCallback)
	l.RegisterEventCallback(locationengine.EvtItemRSSIChange, eventCallback)

	if connected, _ := l.Connect(API_KEY); !connected {
		log.Fatalln("Failed to connect!")
	}
	log.Infoln("Connected successfully")

	if err := l.Subscribe([]string{PLACE}); err != nil {
		log.Fatalln("Failed to subscribe!")
	}
	log.Infoln("Subscribed successfully")

	defer l.Disconnect()

	log.Infoln("Listening for incoming messages")

	for {
		l.ReceiveMessage()
	}


}
