/*
* Eventually this will be a Paho MQTT5.0 client library used for packet generation
* For now, it is sad and empty
 */
package mqttclients

import (
	"context"
	"fmt"
	"log"
	"net"

	paho "github.com/eclipse/paho.golang/paho"
)

func Paho5BasicPubSubTest(brokerAddress string, port string) {

	fullBrokerAddress := brokerAddress + ":" + port

	connexion, err := net.Dial("tcp", fullBrokerAddress)
	if err != nil {
		log.Fatalf("Failed to connect to %s: %s", fullBrokerAddress, err)
	}

	c := paho.NewClient()
	c.Conn = connexion

	connectPacket := &paho.Connect{
		ClientID:     "",
		CleanStart:   true,
		Username:     "paho5user",
		UsernameFlag: true,
	}

	// The Client*.Connect() function requires a Connect* to be passed as an argument, and returns a Connack*
	connackPacket, err := c.Connect(context.Background(), connectPacket)
	if err != nil {
		log.Fatalln(err)
	}
	if connackPacket.ReasonCode != 0 {
		log.Fatalf("Failed to connect to %s : %d - %s", fullBrokerAddress, connackPacket.ReasonCode, connackPacket.Properties.ReasonString)
	}

	fmt.Printf("Connected to %s\n", fullBrokerAddress)

}
