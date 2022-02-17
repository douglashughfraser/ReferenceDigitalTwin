package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	//"strconv"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Message struct {
	Id   string  `json:"id,omitempty"`
	Data float64 `json:"data,omitempty"`
	Str  string  `json:"str,omitempty"`
}

type Component struct {
	behaviour  string
	id         string
	mqttClient mqtt.Client
	db         mongo.Database
}

/*
Constructs a ComponentProfile for components:
* Generates a unique id
* Creates a new MQTTClient for the component and subscribes the client to a slice of *subscriptions* topics.
* Implicitly subscribes the component to it's own unique topic -- Format: "dt/components/<id>"
	- This should be used to communicate with the component. */
func newComponentProfile(strid string, componentType string, subscriptions []string) *Component {

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

		}

		subscription := MQTTClient.Subscribe(topic, 1, callback)
		// Check for errors
		if subscription.Wait() && subscription.Error() != nil {
			panic(subscription.Error())
		}
	}

	profile := Component{
		behaviour:  componentType,
		id:         strid,
		mqttClient: MQTTClient,
	}

	// Connect to core database and add client to profile.
	profile.ConnectDB()

	return &profile
}

func (c *Component) PublishMQTT(topic string, qos byte, retained bool, msg Message) {
	// Encode message to JSON
	jsonMessage, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}

	// Publish message to server
	token := c.mqttClient.Publish(topic, qos, retained, jsonMessage)

	// Verify sent
	if token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
}

func (c *Component) ConnectDB() {

	// Username and password can be added here
	uri := fmt.Sprintf("mongodb://localhost/27017")

	// Connect to MongoDB instance
	client, err := mongo.NewClient(options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatal(err)
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}

	c.db = *client.Database(c.id)
}
