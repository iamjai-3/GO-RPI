package function

import (
	"fmt"
	"os"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/stianeikeland/go-rpio"
)

var topic = "topic1"

var (
	pin = rpio.Pin(4)
)

var messagePubHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	recievedData := string(msg.Payload())
	if recievedData == "on" {
		if err := rpio.Open(); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		defer rpio.Close()

		pin.Output()

		for x := 0; x < 20; x++ {
			pin.Toggle()
			time.Sleep(time.Second / 5)
		}
		return
	}
}

var connectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
	fmt.Println("Connected")

}

var connectLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
	fmt.Printf("Connect lost: %v", err)
}

func publish(client mqtt.Client) {

	token := client.Publish(topic, 0, false, "on")
	token.Wait()
	time.Sleep(time.Second)

}

func sub(client mqtt.Client) {
	token := client.Subscribe(topic, 0, nil)
	token.Wait()
	fmt.Printf("Subscribed to topic: %s \n", topic)
}

func Mqtt() {
	// mqtt.ERROR = log.New(os.Stdout, "[ERROR] ", 0)
	// mqtt.CRITICAL = log.New(os.Stdout, "[CRIT] ", 0)
	// mqtt.WARN = log.New(os.Stdout, "[WARN]  ", 0)
	// mqtt.DEBUG = log.New(os.Stdout, "[DEBUG] ", 0)

	var broker = "broker.emqx.io"
	var port = 1883
	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tcp://%s:%d", broker, port))
	opts.SetClientID("iamjai")
	opts.SetUsername("iotr-admin")
	opts.SetPassword("12345")
	opts.SetDefaultPublishHandler(messagePubHandler)
	opts.SetAutoReconnect(true)
	opts.SetMaxReconnectInterval(10 * time.Second)

	opts.OnConnect = connectHandler
	opts.OnConnectionLost = connectLostHandler
	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
	sub(client)
	// publish(client)

	// client.Disconnect(250)
}
