package main

import (
	"encoding/json"
	"github.com/esiqveland/queuetutor/models"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
	"github.com/bitly/go-nsq"
	"github.com/codegangsta/negroni"
)

type appContext struct {
	producer *nsq.Producer
}
func main() {
	qconfig := nsq.NewConfig()
	w, err := nsq.NewProducer("127.0.0.1:4150", qconfig)
	if err != nil {
		panic(err)
	}

	context := &appContext{w}

	mw := negroni.Classic()
	router := mux.NewRouter()
    router.
		HandleFunc("/applications/submit", context.SubmitApplication).
		Methods("POST")

	mw.UseHandler(router)
	mw.Run(":8000")
}

func (c *appContext) SubmitApplication(w http.ResponseWriter, r *http.Request) {
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("error reading request: %v", err)
		fail(w, "invalid payload")
		return
	}
	app := models.Application{}
	err = json.Unmarshal(data, &app)
	if err != nil {
		log.Printf("invalid unmarshal data: %v", err)
		fail(w, "invalid payload")
		return
	}
	if !app.Valid() {
		log.Printf("invalid application filled in")
		fail(w, "invalid application data")
		return
	}
	err = c.producer.Publish("application.submit", data)
	if err != nil {
		log.Printf("error publishing application %v", err)
		failS(w, 500, "error saving application")
		return
	}

	w.WriteHeader(201)
}
func failS(w http.ResponseWriter, status int, msg string) {
	w.WriteHeader(status)
	w.Write([]byte(msg))
}
func fail(w http.ResponseWriter, msg string) {
	failS(w, 400, msg)
}
