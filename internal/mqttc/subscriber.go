package mqttc

import (
	"encoding/json"
	"fmt"
	"gestanestle/aizm-server/internal/db"
	"gestanestle/aizm-server/internal/models"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/ncruces/go-strftime"
)

const topic = "aizm/conditions"

func Subscribe() {
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
}

var messagePubHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	//fmt.Printf("Received message: %s from topic: %s\n", msg.Payload(), msg.Topic())
	var event models.Conditions
	err := json.Unmarshal([]byte(msg.Payload()), &event)
	if err != nil {
		fmt.Println(err.Error())
	}

	fmt.Printf("Machine ID: %v ", event.ID)

	time := strftime.Format("%Y-%m-%d %H:%M:%S+00", time.Now())

	event.Time = time

	fmt.Println(event.Time)
	d := db.Dao{}
	d.PersistEvent(event)
}

var connectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
	fmt.Println("Connected")
}

var connectLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
	fmt.Println(err)
	sub(client)
}

func sub(client mqtt.Client) {
	token := client.Subscribe(topic, 1, nil)
	token.Wait()
	fmt.Printf("Subscribed to topic %s", topic)
	fmt.Println()
}
