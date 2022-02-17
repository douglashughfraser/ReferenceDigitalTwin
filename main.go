package main

import "fmt"

//"fmt"

func main() {

	initCore()

	// Slice storing profiles of every component
	var components []Component = make([]Component, 0)

	// Create component profiles and connect to MQTT brok
	components = append(components, *newComponentProfile("PhysicalAsset", "NormDistAsset", nil))
	components = append(components, *newComponentProfile("DigitalTwin", "DigitalTwin", []string{"Sensors"}))
	listener([]string{"Sensors"})

	// Iterate through component profiles, in order, calling the appropriate behaviour for that component
	// and running it as a goroutine.
	for i, _ := range components {
		components[i].BehaviourDispatch()
	}

	fmt.Printf("\n-------------------- SET UP COMPLETE --------------------\n\n")

	fmt.Scanln()
}
