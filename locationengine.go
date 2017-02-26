package locationengine

import (
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/satori/go.uuid"
	log "github.com/Sirupsen/logrus"
	"encoding/json"
	"regexp"
	"crypto/tls"
)

const topicPrefix = "/presence/stream/"

type Event int

const (
	EvtItemAppeared Event = iota
	EvtItemDisappeared
	EvtItemRSSIChange
	EvtItemProximityChange
)

type Callback func(evt Event, receiver string, item Item)

type Item struct {
	Timestamp uint
	SourceId string
	TrackingId string
	Rssi int
	Proximity string
	Scantype string
	DeviceAddress string
}

type LocationEngine interface {

	Connect(apiKey string) (bool, error)
	IsConnected() bool
	Disconnect()
	Subscribe(id[] string) error
	ReceiveMessage()
	RegisterCallback(evt Event, callback Callback)
}

type locationEngine struct {
	clientID string
	broker string
	user string
	apiKey string

	url string

	mqttClient mqtt.Client

	mqttHandler mqtt.MessageHandler
	messages chan mqtt.Message
	subscriptions map[string]bool

	knownItems map[string]map[string]Item

	topicRegex * regexp.Regexp
	callbacks map[Event]Callback
}


func New() LocationEngine {
	l := &locationEngine{}

	l.broker = "ssl://ovs.kontakt.io:8083"
	l.user = "user"
	l.clientID = uuid.NewV4().String()
	l.messages = make(chan mqtt.Message)
	l.mqttHandler = func(client mqtt.Client, msg mqtt.Message) {
		l.messages <- msg
	}
	l.subscriptions = make(map[string]bool)
	l.knownItems = make(map[string]map[string]Item)

	l.topicRegex = regexp.MustCompile(`^/presence/stream/(.+)$`)
	l.callbacks = make(map[Event]Callback)

	return l
}

func (l * locationEngine) Connect(apiKey string) (bool, error) {
	if l.IsConnected() {
		return true, nil
	}
	l.apiKey = apiKey

	opts := mqtt.NewClientOptions()
	opts.SetClientID(l.clientID)
	opts.AddBroker(l.broker)
	opts.SetUsername(l.user)
	opts.SetPassword(l.apiKey)

	tlsCfg := &tls.Config {
		InsecureSkipVerify: true,
	}
	opts.SetTLSConfig(tlsCfg)
	log.WithFields(log.Fields{
		"component": "LocationEngine",
	}).Warn("WARNING - Not checking TLS certificate of broker!")

	pahoOnConnectHandler := func(cm mqtt.Client) {
		log.WithFields(log.Fields{
			"component": "LocationEngine",
		}).Debug("Established connection to broker")
		l.resubscribe()
	}
	opts.SetOnConnectHandler(pahoOnConnectHandler)
	pahoConnectionLostHandler := func(cm mqtt.Client, e error) {
		log.WithFields(log.Fields{
			"component": "LocationEngine",
			"error": e,
		}).Debug("Lost connection to broker")
	}
	opts.SetConnectionLostHandler(pahoConnectionLostHandler)

	l.mqttClient = mqtt.NewClient(opts)
	if token := l.mqttClient.Connect(); token.Wait() && token.Error() != nil {
		log.WithFields(log.Fields{
			"component": "LocationEngine",
			"error": token.Error(),
		}).Debug("Unable to connect to broker")
		return false, token.Error()
	}
	return true, nil
}


func (l * locationEngine) IsConnected () bool {
	return l.mqttClient != nil && l.mqttClient.IsConnected()
}

func (l * locationEngine) Disconnect() {
	if (l.IsConnected()) {
		l.mqttClient.Disconnect(250)
	}
}

func (l * locationEngine) resubscribe() {
	var s []string
	for k := range l.subscriptions {
		s = append(s, k)
	}

	l.Subscribe(s)
}

func (l * locationEngine) Subscribe(id[] string) error {
	subscriptions := make(map[string]byte)

	for _, value := range id {
		topic := topicPrefix + value
		log.WithFields(log.Fields{
			"component": "LocationEngine",
			"topic": topic,
		}).Debug("Adding subscription")
		l.subscriptions[value] = true
		subscriptions[topic] = 2
	}

	if token := l.mqttClient.SubscribeMultiple(subscriptions, l.mqttHandler); token.Wait() && token.Error() != nil {
		log.WithFields(log.Fields{
			"component": "LocationEngine",
			"error":     token.Error(),
		}).Error("Error adding subscriptions")
		return token.Error()
	}
	return nil
}


func logItem(i Item, msg string) {
	log.WithFields(log.Fields{
		"component": "LocationEngine",
		"source": i.SourceId,
		"trackingId": i.TrackingId,
		"address": i.DeviceAddress,
		"proximity": i.Proximity,
		"rssi": i.Rssi,
	}).Debug(msg)
}

func (l * locationEngine) invokeCallback(evt Event, receiver string, item Item) {
	c, exists := l.callbacks[evt]
	if exists {
		c(evt, receiver, item)
	}
}

func (l * locationEngine) ReceiveMessage() {
	msg := <-l.messages
	var itemsSeen []Item

	matches := l.topicRegex.FindStringSubmatch(string(msg.Topic()))
	if len(matches) != 2 {
		log.WithFields(log.Fields{
			"component": "LocationEngine",
			"topic": string(msg.Topic()),
			"payload": string(msg.Payload()),
		}).Debugf("Ignoring unrecognized MQTT topic: %v", matches)
		return
	}
	receiver := matches[1]

	if err := json.Unmarshal(msg.Payload(), &itemsSeen); err != nil {
		log.WithFields(log.Fields{
			"component": "LocationEngine",
			"error":     err,
		}).Error("Error parsing payload")
	}


	log.WithFields(log.Fields{
		"component": "LocationEngine",
		"receiver": receiver,
		"deviceList": itemsSeen,
	}).Debug("Incoming device list")

	knownItems, exists := l.knownItems[receiver]
	if (!exists) {
		knownItems = make(map[string]Item)
		l.knownItems[receiver] = knownItems
	}

	addressMap := make(map[string]bool)
	for _, i := range itemsSeen {
		addressMap[i.DeviceAddress] = true
		prevValue, exists := knownItems[i.DeviceAddress]
		if (!exists) {
			logItem(i, "New device in range")
			l.knownItems[receiver][i.DeviceAddress] = i
			l.invokeCallback(EvtItemAppeared, receiver, i)
		} else {
			if i.Proximity != prevValue.Proximity {
				logItem(i, "Proximity change")
				prevValue.Proximity = i.Proximity
				l.knownItems[receiver][i.DeviceAddress] = prevValue
				l.invokeCallback(EvtItemProximityChange, receiver, i)
			}
			if i.Rssi != prevValue.Rssi {
				logItem(i, "RSSI change")
				prevValue.Rssi = i.Rssi
				l.knownItems[receiver][i.DeviceAddress] = prevValue
				l.invokeCallback(EvtItemRSSIChange, receiver, i)
			}
		}
	}
	// Check whether existing devices all appear in list
	for _, i := range knownItems {
		_, exists := addressMap[i.DeviceAddress]
		if !exists {
			logItem(i, "Device out of range")
			delete(l.knownItems[receiver], i.DeviceAddress)
			l.invokeCallback(EvtItemDisappeared, receiver, i)
		}
	}
}

func (l * locationEngine) RegisterCallback(evt Event, callback Callback) {
	l.callbacks[evt] = callback
}
