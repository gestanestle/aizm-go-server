package main

import (
	"fmt"
	"net/http"
	"strconv"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/gofor-little/env"
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

	if err := env.Load(".env"); err != nil {
		panic(err)
	}

	broker, _ := env.MustGet("MQTT_HOST")
	portEnv, _ := env.MustGet("MQTT_PORT")
	port, _ := strconv.Atoi(portEnv)
	clientId, _ := env.MustGet("CLIENT_ID")
	mqttUser, _ := env.MustGet("MQTT_USER")
	mqttPass, _ := env.MustGet("MQTT_PASS")

	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tcp://%s:%d", broker, port))
	opts.SetClientID(clientId)
	opts.SetUsername(mqttUser)
	opts.SetPassword(mqttPass)
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
	http.ListenAndServe(":4000", r)
}

func sub(client mqtt.Client) {
	topic := "aizm/conditions"
	token := client.Subscribe(topic, 1, nil)
	token.Wait()
	fmt.Printf("Subscribed to topic %s", topic)
}
