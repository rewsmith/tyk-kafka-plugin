package main

import (
	"log"
	//"gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
	"encoding/json"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"net/http"
)

const (
	host  = "host.docker.internal"
	topic = "test"
	group = "myGroup"
)

// Message struct
type Message struct {
	Data json.RawMessage `json:"data"`
}

// SendToKafka does what it says on the tin
func SendToKafka(rw http.ResponseWriter, r *http.Request) {

	log.Println("SendToKafka plugin starting")

	p, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": host})
	if err != nil {
		panic(err)
	}

	defer p.Close()

	// Delivery report handler for produced messages
	go func() {
		for e := range p.Events() {
			switch ev := e.(type) {
			case *kafka.Message:
				if ev.TopicPartition.Error != nil {
					log.Println("Delivery failed", ev.TopicPartition)
				} else {
					log.Println("Delivered message to ", ev.TopicPartition)
				}
			}
		}
	}()

	// Produce messages to topic (asynchronously)
	top := topic

	// Declare a new Message struct
	var m Message

	// Try to decode the request body into the struct. If there is an error,
	// respond to the client with the error message and a 400 status code.
	err = json.NewDecoder(r.Body).Decode(&m)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	p.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &top, Partition: kafka.PartitionAny},
		Value:          []byte(m.Data),
	}, nil)

	// Wait for message deliveries before shutting down
	p.Flush(15 * 1000)

}

// ConsumeFromKafka does what it says on the tin
func ConsumeFromKafka(rw http.ResponseWriter, r *http.Request) {

	log.Println("ConsumeFromKafka plugin starting")

	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": host,
		"group.id":          group,
		"auto.offset.reset": "earliest",
	})

	if err != nil {
		panic(err)
	}

	c.SubscribeTopics([]string{topic, "^aRegex.*[Tt]opic"}, nil)

	msg, err := c.ReadMessage(-1)
	if err == nil {
		log.Println("Message on ", msg.TopicPartition, string(msg.Value))
	} else {
		// The client will automatically try to recover from all errors.
		log.Println("Consumer error: ", err, msg)
	}

	c.Close()

	var m Message
	json.Unmarshal(msg.Value, &m)

	/* jsonData, err := json.Marshal(string(msg.Value))
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		return
	} */

	// send HTTP response from Golang plugin
	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	//rw.Write(jsonData)
	rw.Write(m.Data)

}

func main() {}
