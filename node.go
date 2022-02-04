package main

import (
	"encoding/json"
	"fmt"

	//"strconv"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

/*
Constructs a ComponentProfile for components:
* Generates a unique id
* Creates a new MQTTClient for the component and subscribes the client to a slice of *subscriptions* topics.
* Implicitly subscribes the component to it's own unique topic -- Format: "dt/components/<id>"
	- This should be used to communicate with the component. */
func newComponentProfile(strid string, componentType string, subscriptions []string) *ComponentProfile {

	//Generate id for component
	//strid := strconv.Itoa(generateUniqueId())
	fmt.Printf("--- ID: %v; Creating %v component --- \n", strid, componentType)

	//Add subscription to unique topic
	if subscriptions == nil {
		subscriptions = []string{"component/" + strid}
	} else {
		subscriptions = append(subscriptions, "component/"+strid)
	}

	MQTTClient := mqtt.NewClient(GetMQTTClientOptions(strid))
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

			fmt.Printf("\tSub Callback received: Topic: %v, Message: %v\n", msg.Topic(), decodedMessage.Str)
			fmt.Printf("\tSubMessages length: %v\n", len(SubMessages[msg.Topic()]))

			// Use global client->receive() mapping to forward message to component
			ClientReceivers[client](msg.Topic(), decodedMessage)
		}

		subscription := MQTTClient.Subscribe(topic, 1, callback)
		// Check for errors
		if subscription.Wait() && subscription.Error() != nil {
			panic(subscription.Error())
		}
	}

	sendMessage := func(topic string, qos byte, retained bool, msg Message) {
		// Encode message to JSON
		jsonMessage, err := json.Marshal(msg)
		if err != nil {
			panic(err)
		}

		// Publish message to server
		token := MQTTClient.Publish(topic, qos, retained, jsonMessage)

		// Verify sent
		if token.Wait() && token.Error() != nil {
			panic(token.Error())
		}
	}

	receiveMessage := func(topic string, msg Message) {
		fmt.Printf("\tReceived %v message: %v\n", topic, msg.Data)
		// What to do when a message from this topic is recieved
		msgs := SubMessages[topic]
		// If channel at capacity, remove oldest message
		if len(msgs) == cap(msgs) {
			<-msgs
		}
		//Add new message
		msgs <- msg
	}

	// Add entry to global mapping of clients to receivers
	ClientReceivers[MQTTClient] = receiveMessage

	// Bring it all together to create and return profile
	profile := ComponentProfile{
		id:             strid,
		componentType:  componentType,
		mqttClient:     MQTTClient,
		subMessages:    SubMessages,
		SendMessage:    sendMessage,
		ReceiveMessage: receiveMessage,
	}
	return &profile
}
