package main

import (
	"github.com/bitly/go-nsq"
	"log"
	"sync"
	"encoding/json"
	"github.com/esiqveland/queuetutor/models"
	"math/rand"
	"errors"
	"fmt"
)

const applications_submitted = "application.submit"
const read_chan = "email"

func main() {
	// this is only here to make sure the program does not start and quit (finish) immediately.
	// we can do this waitgroup or wait for something else while working
	wg := &sync.WaitGroup{}
	wg.Add(1) // Add one to the waitgroup. This will wait for 1 task completed.

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
		err = sendSketchyEmailCopy(app, 0.1)
		if err != nil {
			return err
		}

		wg.Done() // Mark one task in the waitgroup as done
		return nil
	}))

	err := q.ConnectToNSQD("127.0.0.1:4150")
	if err != nil {
		log.Panic("Could not connect")
	}

	wg.Wait() // wait for all in waitgroup to finish.
}


// sendEmailCopy is a no-op email handler for models.Application
func sendEmailCopy(app models.Application) (error) {
	// app = &application{} // *application (pointer)
	log.Printf("Sending email receipt to: %v", app.Email)

	return nil // notice pointer did not escape scope
	// this is taken advantage of by the compiler and not add pointer to garbage collection,
	// but cleans it up from stack when leaving the scope
}

// sendEmailCopy is a no-op email handler for models.Application
// with a `prob` chance of failing (return with error)
func sendSketchyEmailCopy(app models.Application, prob float64) (error) {
	r := rand.Float64()
	log.Printf("Sending email receipt to: %v", app.Email)

	if prob >= r {
		return errors.New(fmt.Sprintf("Error sending email for application to '%v'", app.Email))
	}

	// app = &application{} // *application (pointer)

	return nil // notice pointer did not escape scope
	// this is taken advantage of by the compiler and not add pointer to garbage collection,
	// but cleans it up from stack when leaving the scope
}
