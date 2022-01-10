package main

import (
	"fmt"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func GetMQTTClientOptions(clientId string) *mqtt.ClientOptions {
	options := mqtt.NewClientOptions()
	// broker IP and port
	options.AddBroker("tcp://127.0.0.1:1883")
	options.SetClientID(clientId)
	options.SetDefaultPublishHandler(func(client mqtt.Client, msg mqtt.Message) {
		fmt.Printf("Default handler got message %v\n", msg.Payload())
	})
	options.OnConnect = func(client mqtt.Client) {
		fmt.Printf("OnConnect handler called for a client\n")
	}
	options.OnConnectionLost = func(client mqtt.Client, err error) {
		fmt.Printf("ConnectionLost handler called for a client, error: %v\n", err.Error())
	}

	return options
}

func CreateMQTTClient(options *mqtt.ClientOptions) *mqtt.Token {
	client := mqtt.NewClient(options)

	token := client.Connect()
	if token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	return &token
}
