/*
* Eventually this will be a Paho MQTT3.1.1 client library used for packet generation, but for now is just a copy of https://github.com/eclipse/paho.mqtt.golang/blob/master/cmd/simple/main.go
 */
package mqttclients

import (
	"fmt"
	"log"
	"os"
	"time"

	paho "github.com/eclipse/paho.mqtt.golang"
)

var msgHandler paho.MessageHandler = func(client paho.Client, msg paho.Message) {
	fmt.Printf("TOPIC: %s\n", msg.Topic())
	fmt.Printf("MSG: %s\n", msg.Payload())
}

func makePaho3(keepAliveTime time.Duration, pingTimeout time.Duration, brokerAddress string, port string) paho.Client {
	paho.DEBUG = log.New(os.Stdout, "", 0)
	paho.ERROR = log.New(os.Stdout, "", 0)
	options := paho.NewClientOptions().AddBroker(brokerAddress + ":" + port).SetClientID("paho3")
	options.SetKeepAlive(keepAliveTime)
	options.SetDefaultPublishHandler(msgHandler)
	options.SetPingTimeout(pingTimeout)

	return paho.NewClient(options)
}

func main() {
	paho.DEBUG = log.New(os.Stdout, "", 0)
	paho.ERROR = log.New(os.Stdout, "", 0)
	opts := paho.NewClientOptions().AddBroker("tcp://iot.eclipse.org:1883").SetClientID("gotrivial")
	opts.SetKeepAlive(2 * time.Second)
	opts.SetDefaultPublishHandler(msgHandler)
	opts.SetPingTimeout(1 * time.Second)

	c := paho.NewClient(opts)
	if token := c.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	if token := c.Subscribe("go-mqtt/sample", 0, nil); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
		os.Exit(1)
	}

	for i := 0; i < 5; i++ {
		text := fmt.Sprintf("this is msg #%d!", i)
		token := c.Publish("go-mqtt/sample", 0, false, text)
		token.Wait()
	}

	time.Sleep(6 * time.Second)

	if token := c.Unsubscribe("go-mqtt/sample"); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
		os.Exit(1)
	}

	c.Disconnect(250)

	time.Sleep(1 * time.Second)
}
