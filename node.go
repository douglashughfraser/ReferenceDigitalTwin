package main

import (
	"fmt"

	//"strconv"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"go.mongodb.org/mongo-driver/mongo"
)

type Message struct {
	Id   string  `json:"id,omitempty"`
	Data float64 `json:"data,omitempty"`
	Str  string  `json:"str,omitempty"`
}

type Component struct {
	behaviour string
	id        string
	mqtt      mqtt.Client
	mongo     mongo.Database
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

	profile := Component{
		behaviour: componentType,
		id:        strid,
	}

	//Add subscription to unique topic
	if subscriptions == nil {
		subscriptions = []string{"component/" + strid}
	} else {
		subscriptions = append(subscriptions, "component/"+strid)
	}

	profile.MQTTMakeClient(MQTTGetTemplateClientOptions(strid))
	profile.MQTTSubscribe(subscriptions)

	// Connect to core database and add client to profile.
	profile.ConnectDB()

	return &profile
}
