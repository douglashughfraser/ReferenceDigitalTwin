package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime"

	//"strconv"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

var nextFreeIdDoNotUseDirect int = 0 // Used as static
var ClientReceivers = make(map[mqtt.Client]func(string, Message))

type Message struct {
	Id   string  `json:"id,omitempty"`
	Data float64 `json:"data,omitempty"`
	Str  string  `json:"str,omitempty"`
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
	subMessages       map[string]chan Message
	SendMessage       func(topic string, qos byte, retained bool, msg Message)
	ReceiveMessage    func(topic string, msg Message)
}

func generateUniqueId() int {
	var newId = nextFreeIdDoNotUseDirect
	nextFreeIdDoNotUseDirect += 1
	return newId
}

func main() {

	// Create mqtt broker in separate window
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" { // Windows:
		cmd = exec.Command("cmd", "/C", "start", "mosquitto", "-p", "1883", "-v")
	} else if runtime.GOOS == "darwin" { // MacOS:
		cmd = exec.Command("bash", "-c", "mosquitto", "-p", "1883", "-v")
	} else if runtime.GOOS == "linux" { // Linux:
		cmd = exec.Command("bash", "-c", "mosquitto", "-p", "1883", "-v")
	} else {
		panic(runtime.GOOS)
	}

	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}

	mqtt.ERROR = log.New(os.Stdout, "[ERROR] ", 0)
	mqtt.CRITICAL = log.New(os.Stdout, "[CRIT] ", 0)
	mqtt.WARN = log.New(os.Stdout, "[WARN]  ", 0)
	//mqtt.DEBUG = log.New(os.Stdout, "[DEBUG] ", 0)

	// Slice storing profiles of every component
	var components []ComponentProfile = make([]ComponentProfile, 0)

	// Create component profiles and connect to MQTT broker
	components = append(components, *newComponentProfile("asset", "normDist", nil))
	components = append(components, *newComponentProfile("dt", "dt", []string{"Sensors"}))
	components = append(components, *newComponentProfile("listener", "fmtListener", []string{"fmt"}))

	// Iterate through component profiles, in order, loading the behaviour for that component
	// and running it as a goroutine.
	for i := range components {
		doTheThings := getBehaviour(components[i].componentType)
		go doTheThings(&components[i])
	}

	fmt.Printf("\n-------------------- SET UP COMPLETE --------------------\n\n")

	fmt.Scanln()

}
