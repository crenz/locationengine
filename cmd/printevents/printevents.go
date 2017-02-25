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

	connected, _ := l.Connect(API_KEY)

	if (connected != l.IsConnected()) {
		log.Errorln("Return value does not equal connection state")
	}

	if (!l.IsConnected()) {
		log.Errorln("Failed to connect!")
	}

	log.Print("Connected successfully")

	err := l.Subscribe([]string{PLACE})
	if err != nil {
		log.Errorln("Failed to subscribe!")
	}
	log.Print("Subscribed successfully")

	defer l.Disconnect()

	log.Infoln("Listening for incoming messages")

	for {
		l.ReceiveMessage()
	}


}
