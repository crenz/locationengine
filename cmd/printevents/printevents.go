package main

import (
	"github.com/crenz/locationengine"
	log "github.com/Sirupsen/logrus"
	"github.com/mattn/go-colorable"
)

const API_KEY = "your_api_key"
const PLACE = "your_place_uuid"

func callback(evt locationengine.Event, receiver string, item locationengine.Item) {
	log.Infof("Callback [%i] %s %+v\n", evt, receiver, item)
}

func init() {
	// Set to log.DebugLevel to see each individual MQTT message received
	log.SetLevel(log.InfoLevel)
	log.SetOutput(colorable.NewColorableStdout())
}

func main() {
	log.Infoln("Connecting to Kontakt.io broker")
	l := locationengine.New()
	l.RegisterCallback(locationengine.EvtItemAppeared, callback)
	l.RegisterCallback(locationengine.EvtItemDisappeared, callback)
	l.RegisterCallback(locationengine.EvtItemProximityChange, callback)
	l.RegisterCallback(locationengine.EvtItemRSSIChange, callback)

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
