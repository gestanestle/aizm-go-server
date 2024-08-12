package mqttc

import (
	"encoding/json"
	"fmt"
	"gestanestle/aizm-server/internal/db"
	"gestanestle/aizm-server/internal/models"
	"math/rand"
	"os"
	"strconv"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/ncruces/go-strftime"
)

var topic = os.Getenv("MQTT_TOPIC")
var host = os.Getenv("MQTT_HOST")
var port, _ = strconv.Atoi(os.Getenv("MQTT_PORT"))
var user = os.Getenv("MQTT_USER")
var pass = os.Getenv("MQTT_PASS")

func randSeq(n int) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func Subscribe() {

	id := randSeq(10)

	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tcp://%s:%d", host, port))
	opts.SetClientID(id)
	opts.SetUsername(user)
	opts.SetPassword(pass)
	opts.SetDefaultPublishHandler(messagePubHandler)
	opts.SetCleanSession(true)
	opts.OnConnect = connectHandler
	opts.OnConnectionLost = connectLostHandler
	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	fmt.Printf("MQTT Client ID: %s", id)
	fmt.Println()

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

	loc, _ := time.LoadLocation("Asia/Shanghai")
	time := strftime.Format("%Y-%m-%d %H:%M:%S+00", time.Now().In(loc))

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
