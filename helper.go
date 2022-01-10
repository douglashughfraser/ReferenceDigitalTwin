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

/**
Database integration

func (record DataStructure)GetID(id int){
	// not sure what goes in here
}

func (record DataStructure)SetID() int{
	return record.id
}

func AfterFind(db *Database) error {
	// does anything need to be done after find?
	return nil
}**/
func getBehaviour(componentType string) ComponentBehaviour {
	var behaviour ComponentBehaviour
	if componentType == "normDist" {
		behaviour = func(profile *ComponentProfile) {
			// Publish init message to subscribers, retained (true) by broker for future subscribers.
			profile.mqttClient.Publish(profile.id, 1, true, Message{
				id:   profile.id,
				data: 0.0,
				str:  "Physical asset operational"})
			for {
				select {
					case msg := recieved message:
					fmt.Printf("Asset %v recieved! Sender: %v, Data: %v\n", profile.id, msg.id, msg.str)
					profile.mqttClient.Publish(profile.id, 1, true, Message{
						id:   profile.id,
						data: 0.0,
						str:  "Physical asset terminating"})
					if msg.data == 0.00 {
						break
					}
				default:
					fmt.Printf("Sensor %v ready to send!\n", profile.id)
					// Quality of Service, 0: At most once, 1: At least once, 2: Exactly once
					profile.mqttClient.Publish("Sensors", 1, false, Message{
						id:   profile.id,
						data: rand.NormFloat64()*15 + 50,
						str:  "physical reading"})
					profile.mqttClient.Publish(profile.id, 0, false, "super message")
					fmt.Printf("Sensor %v sent!\n", profile.id)
					time.Sleep(time.Duration(500+rand.Intn(500)-250) * time.Millisecond)
				}
			}
		}
	} else if componentType == "monitor" {
		behaviour = func(profile *ComponentProfile) {
			messagesRecieved := 0
			for {
				fmt.Printf("Model waiting to receive.\n")
				select {
				case msg := <-profile.input:
					fmt.Printf("DB Recieved! Sender: %v, Data: %v\n", msg.id, msg.data)
					messagesRecieved += 1
					if messagesRecieved > 30 {
						profile.mqttClient.Publish(profile.id, 1, false, Message{
							id:   profile.id,
							data: 1.00,
							str:  "Be Prosperous"})
					} else {
						profile.mqttClient.Publish(profile.id, 1, false, Message{
							id:   profile.id,
							data: 0.00,
							str:  "Maintain"})
					}
				}
			}
		}
	} else if componentType == "listener" {
		behaviour = func(profile *ComponentProfile) {
			for {
				fmt.Printf("Actuator waiting to receive.\n")
				select {
				case msg := <-profile.input:
					if msg.data == 1.00 {
						fmt.Printf("Actuator %v does something\n", profile.id)
					} else if msg.data == 0.00 {
						fmt.Printf("Actuator %v does nothing\n", profile.id)
					} else {
						fmt.Printf("Actuator %v is confused\n", profile.id)
					}
				}
			}
		}
	}
	return behaviour
}

/**
// Database integration
//Make database structure
store, err := disk.New("./data", ".json")
if (err != nil){
	print("on no\n")
}

db, err := hare.New(store)

defer profile.waitgroup.Done()
for {
	select{
	case msg := <- profile.input:

		// Assumes all messages contain data for storage i.e. can't find it? Store it.
		record := DataStructure{}

		if err := db.Find("components", msg.id, &record); err != nil{
			err := db.Update("components", msg.id)
		}else{
			// If err, add entry to database
			recID, err := db.Insert("components", msg.id)
		}
	}
}**/

/**
func initModel() ([]int, *chan Message) { //operation func(Message)
	modelChan := make(chan Message)
	go listen("Subscriber 1", modelChan)

	// generate every id
	var subIds []int
	for i := 0; i < numSensors; i++ {
		subIds = append(subIds, i)
	}

	return subIds, &modelChan
}
**/
