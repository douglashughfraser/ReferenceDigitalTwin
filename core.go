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

	fmt.Print("Core: Initializing MQTT Server...")

	// Create mqtt broker in separate window
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" { // Windows:
		fmt.Print("Windows...")
		cmd = exec.Command("cmd", "/C", "start", "mosquitto", "-p", "1883", "-v")
	} else if runtime.GOOS == "darwin" { // MacOS: not working
		fmt.Print("MacOS...")
		cmd = exec.Command("bash", "--rcfile", "<(echo '. ~/.bashrc; some_command')")
	} else if runtime.GOOS == "linux" { // Linux: not working
		fmt.Print("Linux...")
		cmd = exec.Command("bash", "-c", "mosquitto", "-p", "1883", "-v")
	} else {
		panic(runtime.GOOS)
	}

	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	} else {
		fmt.Println("Done.")
	}

	mqtt.ERROR = log.New(os.Stdout, "[ERROR] ", 0)
	mqtt.CRITICAL = log.New(os.Stdout, "[CRIT] ", 0)
	mqtt.WARN = log.New(os.Stdout, "[WARN]  ", 0)
	//mqtt.DEBUG = log.New(os.Stdout, "[DEBUG] ", 0)

	// --- Create Core DB Catalogue ---

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
		fmt.Print("MacOS...")
		cmd = exec.Command("bash", "open")
	} else {
		panic(runtime.GOOS)
	}

	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	} else {
		fmt.Println("Done.")
	}
}
