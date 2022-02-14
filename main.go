package main

//"fmt"

func main() {

	initCore()
	/**
		// Slice storing profiles of every component
		var components []Component = make([]Component, 0)

		// Create component profiles and connect to MQTT broker
		components = append(components, *newComponentProfile("PhysicalAsset", "normDist", nil))
		components = append(components, *newComponentProfile("DigitalTwin", "dt", []string{"Sensors"}))
		listener([]string{"Sensors"})

		// Iterate through component profiles, in order, loading the behaviour for that component
		// and running it as a goroutine.
		for _, component := range components {
			component.run()
		}

		fmt.Printf("\n-------------------- SET UP COMPLETE --------------------\n\n")

		fmt.Scanln()
	**/
}
