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

func generateUniqueId() int {
	var newId = nextFreeIdDoNotUseDirect
	nextFreeIdDoNotUseDirect += 1
	return newId
}

func initCore() {

	fmt.Print("Core: Initializing...")

	// Create mqtt broker in separate window
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" { // Windows:
		cmd = exec.Command("cmd", "/C", "start", "mosquitto", "-p", "1883", "-v")
	} else if runtime.GOOS == "darwin" { // MacOS:
		fmt.Print("MacOS...")
		cmd = exec.Command("bash", "--rcfile", "<(echo '. ~/.bashrc; some_command')")
	} else if runtime.GOOS == "linux" { // Linux:
		fmt.Print("Linux...")
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

	fmt.Println("Done.")
}
