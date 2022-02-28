package main

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

type DataStructure struct {
	id   int
	data []float64
}

func (component *Component) Run() {
	if component.behaviour == "NormDistAsset" {
		go component.NormDistAsset()
	} else if component.behaviour == "DigitalTwin" {
		go component.DigitalTwin()
	} else {
		panic(fmt.Sprintf("No appropriate behaviour %v found.", component.behaviour))
	}
}

func (asset *Component) NormDistAsset() {
	fmt.Printf("%v asset running NormDistBehaviour\n", asset.id)
	for {
		data := rand.NormFloat64()*15 + 50

		// Publish Event
		// Quality of Service, 0: At most once, 1: At least once, 2: Exactly once
		asset.MQTTPublish("Sensors", 1, false, Message{
			Id:   asset.id,
			Data: data,
			Str:  "physical reading"})

		// Insert Data into DB
		doc := bson.D{{"Reading", data}}
		result, err := asset.mongo.Collection("Readings").InsertOne(context.TODO(), doc)
		if err != nil {
			panic(err)
		}
		fmt.Printf("Sensor %v: \t Sent and stored data reading: %v\n\t StorageID: %v\n", asset.id, data, result)
		time.Sleep(time.Duration(10000+rand.Intn(500)-250) * time.Millisecond)
	}
}

func (twin *Component) DigitalTwin() {
	fmt.Printf("%v asset running DigitalTwin behaviour\n", twin.id)
	for {
		data := rand.NormFloat64()*15 + 50
		// Quality of Service, 0: At most once, 1: At least once, 2: Exactly once
		twin.MQTTPublish("Sensors", 1, false, Message{
			Id:   twin.id,
			Data: data,
			Str:  "physical reading"})
		fmt.Printf("Sensor %v: \t Sent data reading: %v\n", twin.id, data)
		time.Sleep(time.Duration(10000+rand.Intn(500)-250) * time.Millisecond)
	}
}
