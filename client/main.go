package main

import (
	"github.com/bitly/go-nsq"
	"log"
	"sync"
)

const applications_submitted = "application.submit"
const read_chan = "metrics"

func main() {
	wg := &sync.WaitGroup{}
	wg.Add(1)

	qconfig := nsq.NewConfig()

	q, _ := nsq.NewConsumer(applications_submitted, read_chan, qconfig)
	q.AddHandler(nsq.HandlerFunc(func(message *nsq.Message) error {
		log.Printf("Got a message: %v", message)
		log.Printf("ID: %v Body: '%v", message.ID, string(message.Body))

		wg.Done()
		return nil
	}))

	err := q.ConnectToNSQD("127.0.0.1:4150")
	if err != nil {
		log.Panic("Could not connect")
	}

	wg.Wait()
}
