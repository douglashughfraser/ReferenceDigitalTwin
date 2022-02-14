package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type DataStructure struct {
	id   int
	data []float64
}

func (asset *PhysicalAsset) run() {
	for {
		data := rand.NormFloat64()*15 + 50
		// Quality of Service, 0: At most once, 1: At least once, 2: Exactly once
		asset.SendMessage("Sensors", 1, false, Message{
			Id:   asset.id,
			Data: data,
			Str:  "physical reading"})
		fmt.Printf("Sensor %v: \t Sent data reading: %v\n", asset.id, data)
		time.Sleep(time.Duration(10000+rand.Intn(500)-250) * time.Millisecond)
	}
}

func (twin *DigitalTwin) run() {
	// Handle received messages
	for _, topicMessages := range twin.subMessages {
		select {
		case msg := <-topicMessages:
			fmt.Printf("\tDT Trying to send message\n")
			twin.SendMessage("fmt", 1, true, Message{
				Id:   twin.id,
				Data: 0.0,
				Str:  fmt.Sprintf("Digital Twin (%v) Recieved: %v\n", twin.id, msg)})
		}
	}
}

func listener(subscriptions []string) {
	// Handle received messages

	MQTTClient := mqtt.NewClient(GetMQTTClientOptions("listener"))
	// Connect component to mqttbroker
	connection := MQTTClient.Connect()
	if connection.Wait() && connection.Error() != nil {
		panic(connection.Error())
	}

	// Map of channels containing recieved messages, split by topic.
	SubMessages := make(map[string]chan Message)

	for _, topic := range subscriptions {
		// Make channel of mqtt.message type, messages to MQTT Client are sent to these channels to be read
		// buffered with size 2: Don't want channel to block, front entry will be removed by handler if limit is reached, not 1 to give leeway
		fmt.Println(SubMessages)
		SubMessages[topic] = make(chan Message, 5)

		callback := func(client mqtt.Client, msg mqtt.Message) {

			// Decode JSON to Message interface
			var decodedMessage Message
			if err := json.Unmarshal(msg.Payload(), &decodedMessage); err != nil {
				panic(err)
			}

			fmt.Printf("\tListener: Topic: %v, Message: %v\n", msg.Topic(), decodedMessage.Str)
		}

		subscription := MQTTClient.Subscribe(topic, 1, callback)
		// Check for errors
		if subscription.Wait() && subscription.Error() != nil {
			panic(subscription.Error())
		}
	}
}
