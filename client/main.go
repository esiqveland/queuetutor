package main

import (
	"github.com/bitly/go-nsq"
	"log"
	"sync"
	"encoding/json"
	"github.com/esiqveland/queuetutor/models"
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

		app := models.Application{}
		err := json.Unmarshal(message.Body, &app)
		if err != nil { // maybe do something else here? put it on a failed application queue for manual inspection?
			return err
		}

		wg.Done()
		return nil
	}))

	err := q.ConnectToNSQD("127.0.0.1:4150")
	if err != nil {
		log.Panic("Could not connect")
	}

	wg.Wait()
}

func noopHandler() {
	// app = &application{} // *application (pointer)

	return // notice pointer did not escape scope
	// this is taken advantage of by the compiler and not add pointer to garbage collection,
	// but cleans it up from stack when leaving the scope
}
