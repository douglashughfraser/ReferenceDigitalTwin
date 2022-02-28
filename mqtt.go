package main

import (
	"encoding/json"
	"fmt"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

// Gives access to component profiles, needed within message handling callbacks.
var switchboard map[mqtt.Client]*Component

func (c *Component) MQTTMakeClient(options *mqtt.ClientOptions) *mqtt.Token {
	client := mqtt.NewClient(options)

	token := client.Connect()
	if token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	switchboard[client] = c

	return &token
}

func MQTTGetTemplateClientOptions(clientId string) *mqtt.ClientOptions {
	options := mqtt.NewClientOptions()
	// broker IP and port
	options.AddBroker("tcp://127.0.0.1:1883")
	options.SetClientID(clientId)
	options.SetDefaultPublishHandler(func(client mqtt.Client, msg mqtt.Message) {
		fmt.Printf("Default handler called for msg with payload: %v\n", msg.Payload())
	})
	options.OnConnect = func(client mqtt.Client) {
		options := client.OptionsReader()
		id := options.ClientID()
		fmt.Printf("OnConnect handler called for Client: %v\n", id)
	}
	options.OnConnectionLost = func(client mqtt.Client, err error) {
		options := client.OptionsReader()
		id := options.ClientID()
		fmt.Printf("ConnectionLost han dler called for Client: %v, error: %v\n", id, err.Error())
	}

	return options
}

func (c *Component) MQTTSubscribe(subscriptions []string) {
	// Map of channels containing recieved messages, split by topic.
	//SubMessages := make(map[string]chan Message)

	for _, topic := range subscriptions {
		// Make channel of mqtt.message type, messages to MQTT Client are sent to these channels to be read
		// buffered with size 2: Don't want channel to block, front entry will be removed by handler if limit is reached, not 1 to give leeway
		//fmt.Println(SubMessages)
		//SubMessages[topic] = make(chan Message, 5)

		callback := func(client mqtt.Client, msg mqtt.Message) {

			// Decode JSON to Message interface
			var decodedMessage Message
			if err := json.Unmarshal(msg.Payload(), &decodedMessage); err != nil {
				panic(err)
			}

			fmt.Printf("%v received MQTT message\n", switchboard[client].id)
		}
		subscription := c.mqtt.Subscribe(topic, 1, callback)
		// Check for errors
		if subscription.Wait() && subscription.Error() != nil {
			panic(subscription.Error())
		}
	}
}

func (c *Component) MQTTPublish(topic string, qos byte, retained bool, msg Message) {
	// Encode message to JSON
	jsonMessage, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}

	// Publish message to server
	token := c.mqtt.Publish(topic, qos, retained, jsonMessage)

	// Verify sent
	if token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
}
