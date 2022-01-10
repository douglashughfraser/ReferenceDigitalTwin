package main

import (
	"fmt"
	"strconv"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

const numSensors = 1
const numModels = 1
const numActuators = 1

type Message struct {
	id   string
	data float64
	str  string
}

type ComponentBehaviour = func(profile *ComponentProfile)

// Struct containing all the data about a component.
// Used but ComponentBehaviour function and to communicate with the component.
type ComponentProfile struct {
	id                string
	componentType     string
	input             chan Message
	mqttClient        mqtt.Client
	mqttSubscriptions []string
}

var nextFreeIdDoNotUseDirect int = 0

func generateUniqueId() int {
	var newId = nextFreeIdDoNotUseDirect
	nextFreeIdDoNotUseDirect += 1
	return newId
}

/*
Constructs a ComponentProfile for components:
* Generates a unique id
* Creates a new MQTTClient for the component and subscribes the client to a slice of *subscriptions* topics.
* Implicitly subscribes the component to it's own unique topic -- Format: "iot/components/<id>"
	- This should be used to communicate with the component. */
func newComponentProfile(componentType string, subscriptions []string) *ComponentProfile {

	//Generate id for component
	strid := strconv.Itoa(generateUniqueId())

	//Add subscription to unique topic
	if subscriptions == nil {
		subscriptions = []string{"iot/components/" + strid}
	} else {
		subscriptions = append(subscriptions, "iot/components/"+strid)
	}

	MQTTClient := mqtt.NewClient(GetMQTTClientOptions(strid))

	// Connect component to mqttbroker
	connection := MQTTClient.Connect()
	if connection.Wait() && connection.Error() != nil {
		panic(connection.Error())
	}

	for _, topic := range subscriptions {
		subscription := MQTTClient.Subscribe(topic, 1, nil)
		if subscription.Wait() && subscription.Error() != nil {
			panic(subscription.Error())
		} else {
			fmt.Println("Subscribed to %v", topic)
		}
	}

	profile := ComponentProfile{
		id:            strid,
		componentType: componentType,
		mqttClient:    MQTTClient,
	}

	return &profile
}

func main() {

	// Slice storing profiles of every component
	var components []ComponentProfile = make([]ComponentProfile, 0)

	// Create component profiles and connect to MQTT broker
	components = append(components, *newComponentProfile("normDistSensor", nil))
	components = append(components, *newComponentProfile("monitor", []string{"iot/type/Sensors"}))

	// Iterate through component profiles, in order, loading the behaviour for that component
	// and running it as a goroutine.

	fmt.Println(len(components))
	for i, component := range components {
		doTheThings := getBehaviour(component.componentType)
		fmt.Println(i)
		go doTheThings(&component)
	}

	fmt.Printf("Initialized\n")

	fmt.Scanln()

}
