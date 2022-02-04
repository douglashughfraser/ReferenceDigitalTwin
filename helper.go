package main

import (
	"fmt"
	"math/rand"
	"time"
)

type DataStructure struct {
	id   int
	data []float64
}

func getBehaviour(componentType string) ComponentBehaviour {
	var behaviour ComponentBehaviour
	if componentType == "normDist" {
		behaviour = func(profile *ComponentProfile) {
			// Publish init message to subscribers, retained (true) by broker for future subscribers.
			profile.mqttClient.Publish("fmt", 1, true, Message{
				Id:   profile.id,
				Data: 0.0,
				Str:  "Physical asset operational"})
			// Produce and publish data, forever
			for {
				data := rand.NormFloat64()*15 + 50
				// Quality of Service, 0: At most once, 1: At least once, 2: Exactly once
				profile.SendMessage("Sensors", 1, false, Message{
					Id:   profile.id,
					Data: data,
					Str:  "physical reading"})
				fmt.Printf("Sensor %v: \t Sent data reading: %v\n", profile.id, data)
				time.Sleep(time.Duration(10000+rand.Intn(500)-250) * time.Millisecond)
			}
		}
	} else if componentType == "dt" {
		behaviour = func(profile *ComponentProfile) {
			// Handle received messages
			for _, topicMessages := range profile.subMessages {
				select {
				case msg := <-topicMessages:
					fmt.Printf("\tDT Trying to send message\n")
					profile.SendMessage("fmt", 1, true, Message{
						Id:   profile.id,
						Data: 0.0,
						Str:  fmt.Sprintf("Digital Twin (%v) Recieved: %v\n", profile.id, msg)})
				}
			}
		}
	} else if componentType == "fmtListener" {
		behaviour = func(profile *ComponentProfile) {
			// Handle received messages
			for _, topicMessages := range profile.subMessages {
				select {
				case msg := <-topicMessages:
					fmt.Printf("Logger recieved: %v\n", msg.Str)
				}
			}
		}
	}
	return behaviour
}
