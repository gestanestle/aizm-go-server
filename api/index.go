package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/gorilla/mux"
)

var messagePubHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	fmt.Printf("Received message: %s from topic: %s\n", msg.Payload(), msg.Topic())
}

var connectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
	fmt.Println("Connected")
}

var connectLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
	fmt.Println(err)
	sub(client)
}

func main() {

	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tcp://%s:%d", "broker.emqx.io", 1883))
	opts.SetClientID("aizm-server")
	opts.SetUsername("admin")
	opts.SetPassword("public")
	opts.SetDefaultPublishHandler(messagePubHandler)
	opts.SetCleanSession(true)
	opts.OnConnect = connectHandler
	opts.OnConnectionLost = connectLostHandler
	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	sub(client)

	r := mux.NewRouter()
	r.HandleFunc("/", Hello).Methods("GET")
	http.ListenAndServe(":4000", r)
}

func sub(client mqtt.Client) {
	topic := "aizm/conditions"
	token := client.Subscribe(topic, 1, nil)
	token.Wait()
	fmt.Printf("Subscribed to topic %s", topic)
}

func Hello(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	res := Response{
		Status:  "OK",
		Message: "Hello from AIZM Server",
	}
	json.NewEncoder(w).Encode(res)
}

type Response struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}
