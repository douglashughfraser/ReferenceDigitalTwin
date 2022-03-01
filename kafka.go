package main

import (
	"context"
	"fmt"
	kafka "github.com/segmentio/kafka-go"
	"net"
)

// Initialize kafka writers for each topic and store them in a map within the component profile
func (c *Component) KafkaMakeWriters(topics []string) {
	KafkaWriters := make(map[string]*kafka.Writer)
	for _, topic := range topics {
		KafkaWriters[topic] = KafkaMakeWriter(topic)
	}

	c.KafkaWriters = KafkaWriters
}

// Intialize a kafka writer for a given topic using the default broker address
func KafkaMakeWriter(topic string) *kafka.Writer {
	addr, err := net.ResolveTCPAddr("tcp", "localhost:9092")
	if err != nil {
		panic(err)
	}

	return &kafka.Writer{
		Addr:  addr,
		Topic: topic,
		Async: true,
		Completion: func(messages []kafka.Message, err error) {
			if err != nil {
				for _, message := range messages {
					fmt.Println("Message %v failed\n", message.Value)
				}
			}
		},
	}
}

func (c *Component) KafkaPublish(topic string, ctx context.Context) {
	/**
	for {
		// each kafka message has a key and value. The key is used
		// to decide which partition (and consequently, which broker)
		// the message gets published on
		err := c.KafkaWriters[topic].WriteMessages(ctx, kafka.Message{
			Key: []byte(strconv.Itoa(i)),
			// create an arbitrary message payload for the value
			Value: []byte("this is message" + strconv.Itoa(i)),
		})
		if err != nil {
			panic("could not write message " + err.Error())
		}

		// log a confirmation once the message is written
		fmt.Println("writes:", i)
		i++
		// sleep for a second
		time.Sleep(time.Second)
	}**/
}

func consume(topic string, ctx context.Context) {
	// initialize a new reader with the brokers and topic
	// the groupID identifies the consumer and prevents
	// it from receiving duplicate messages
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{"localhost:9092"}, // default broker address
		Topic:   topic,
		GroupID: "my-group",
	})
	for {
		// the `ReadMessage` method blocks until we receive the next event
		msg, err := r.ReadMessage(ctx)
		if err != nil {
			panic("could not read message " + err.Error())
		}
		// after receiving the message, log its value
		fmt.Println("received: ", string(msg.Value))
	}
}

func makeKafkaReader(topic string) *kafka.Reader {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{"localhost:9092"}, // default broker address
		Topic:   topic,
		GroupID: "my-group",
	})

	return r
}
