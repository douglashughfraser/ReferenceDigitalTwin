package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	//"strconv"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

var nextFreeIdDoNotUseDirect int = 0 // Used as static
var ClientReceivers = make(map[mqtt.Client]func(string, Message))

func generateUniqueId() int {
	var newId = nextFreeIdDoNotUseDirect
	nextFreeIdDoNotUseDirect += 1
	return newId
}

func initCore() {

	// N.B. If Mosquitto doesn't launch in a seperate window, check if a mosquitto service is already running and using the port. Stop any service.

	mqtt.ERROR = log.New(os.Stdout, "[ERROR] ", 0)
	mqtt.CRITICAL = log.New(os.Stdout, "[CRIT] ", 0)
	mqtt.WARN = log.New(os.Stdout, "[WARN]  ", 0)
	//mqtt.DEBUG = log.New(os.Stdout, "[DEBUG] ", 0)

	// Start mqtt
	StartMosquitto()

	// Start MongoDB
	StartMongo()

	// Kafka needs to be manually started until I figure out how to run the sh scripts from go.
	//StartKafka()
}

// Start mosquitto mqtt broker with default settings
func StartMosquitto() {
	fmt.Print("Core: Initializing MQTT Server...")

	//Initialize client->component mapping
	switchboard = make(map[mqtt.Client]*Component)

	// Create mqtt broker in separate window
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" { // Windows:
		fmt.Print("Windows...")
		cmd = exec.Command("cmd", "/C", "start", "mosquitto", "-p", "1883", "-v")
	} else if runtime.GOOS == "darwin" { // MacOS: not working
		fmt.Print("MacOS...please run mosquitto manually")
		//cmd = exec.Command("bash"? "/bin/sh"?, "open")
	} else {
		fmt.Print("Unknown environment...please run mosquitto manually")
		//cmd = exec.Command("bash", "-c", "mosquitto", "-p", "1883", "-v")
	}

	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	} else {
		fmt.Println("Done.")
	}
}

// Start mongo instance with default settings
func StartMongo() {
	var cmd *exec.Cmd
	fmt.Print("Core: Initializing Core Mongo instance...")
	// Get current directory of file
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	// Construct core path from
	dbPath := filepath.Join(filepath.Dir(ex), "../data/db")

	if runtime.GOOS == "windows" { // Windows:
		fmt.Print("Windows...")
		cmd = exec.Command("cmd", "/C", "start", "mongod", "--port", "27017", "--dbpath", dbPath)
	} else if runtime.GOOS == "darwin" { // MacOS: not working
		fmt.Print("MacOS...please run mongod manually")
		//cmd = exec.Command("bash"? "/bin/sh"?, "open")
	} else {
		fmt.Print("Unknown environment...please run mongod manually")
	}

	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	} else {
		fmt.Println("Done.")
	}
}

/**
func StartKafka() {
	var cmd *exec.Cmd
	var kafkapath string
	if runtime.GOOS == "windows" { // Windows:
		fmt.Print("Windows...")
		kafkapath = "C:\\Users\\douglas\\Documents\\kafka_2.13-3.1.0"

		// check kafka exists
		if _, err := os.Stat(kafkapath); os.IsNotExist(err) {
			panic(err)
		}
		fmt.Print("Kafka Found...")
		cmd = exec.Command("/bin/sh",
			"C:\\Users\\douglas\\Documents\\kafka_2.13-3.1.0\\bin\\zookeeper-server-start.sh",
			"C:\\Users\\douglas\\Documents\\kafka_2.13-3.1.0\\config\\zookeeper.properties")
	} else if runtime.GOOS == "darwin" { // MacOS: not working
		fmt.Print("MacOS...please run kafka zookeeper and server manually")
		//cmd = exec.Command("bash"? "/bin/sh"?, "open")
	} else {
		fmt.Print("Unknown environment...please run kafka zookeeper and server manually")
	}

	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	} else {
		fmt.Println("Done.")
	}
}
*/
