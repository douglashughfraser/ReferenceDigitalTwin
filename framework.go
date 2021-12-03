package main

import (
	"fmt"
	"math/rand"
	"reflect"
	"sync"
	"time"
)

const numSensors = 10

func initSensors(wg *sync.WaitGroup, chans []chan float64) {
	for i := 0; i < numSensors; i++ {
		wg.Add(1) // add sensor to wait group
		chans[i] = make(chan float64)

		go func(ch chan float64, interval int) {
			defer wg.Done()
			for i := 0.00; i < 100.00; i++ {
				ch <- rand.Float64() * 100 // add probabilitic distribution?
				time.Sleep(time.Duration(interval+rand.Intn(500)-250) * time.Millisecond)
			}
		}(chans[i], 500) // Randomize update interval
	}
}

func main() {

	var wg sync.WaitGroup // used to synchronize sensors to start only once receiver created.
	wg.Add(1)             // Add listener to wait group first before creating sensors

	// create channels for data exchange
	// currently only being made for sensors
	chans := make([]chan float64, numSensors)

	initSensors(&wg, chans) // create sensors for each channel, adding each sensor to the waitgroup (passed using pointer ref)

	// initialize a receiver for all channels
	cases := make([]reflect.SelectCase, len(chans)) // make an array of SelectCases (templates for receiving/sending)
	for i, ch := range chans {                      // for each channel
		cases[i] = reflect.SelectCase{ // make a recieve case for the type of that channel
			Dir:  reflect.SelectRecv,
			Chan: reflect.ValueOf(ch)}
	}

	for i := 0; i < len(chans); i++ {
		fmt.Print(len(chans[i]))
	}

	// listen on all channels
	go func([]chan float64) {
		//var data [][]float64 // variable array to collect data in
		wg.Done()
		for {
			ch, v, _ := reflect.Select(cases)
			// ok will be true if the channel has not been closed.
			fmt.Printf("Sensor %v: %v \n", ch, v)
			//add new data from channel to array, bit of conversion needed to handle reflection
			//data[ch] = append(data[ch], v.Interface().(float64))

			// actuators
		}
	}(chans)

	fmt.Scanln()
}
