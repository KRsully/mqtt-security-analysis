/*
* Eventually this will drive our MQTT3.1.1 and MQTT5.0 clients, but for now is just a copy of https://github.com/eclipse/paho.mqtt.golang/blob/master/cmd/simple/main.go
 */

package main

import (
	"flag"

	"github.com/KRsully/mqtt-security-analysis/mqttclients"
)

func main() {
	port := flag.String("p", "1883", "port to capture MQTT packets on")
	brokerAddress := flag.String("b", "127.0.0.1", "IP address or FQDN of the target broker")
	certFilePath := flag.String("cafile", "certs/ca.crt", "filepath to the certificate of ")
	basicTest := flag.Bool("test-basic", false, "Run basic pub/sub test")
	basicCertTest := flag.Bool("test-cert", false, "Run basic pub/sub test using a certificate")
	mqttVersion := flag.Int("-V", 4, "MQTT version 4 for 3.1.1, or 5 for 5.0")

	flag.Parse()

	if *mqttVersion == 4 {
		// Run a test from the 3.1.1 tests table
		if *basicTest {
			mqttclients.Paho3BasicPubSubTest(*brokerAddress, *port)
		} else if *basicCertTest {
			mqttclients.Paho3BasicCertTest(*brokerAddress, *port, *certFilePath)
		}
	} else if *mqttVersion == 5 {
		// Run a test from the 5.0 tests table
		if *basicTest {
			mqttclients.Paho5BasicPubSubTest(*brokerAddress, *port)
		} else if *basicCertTest {
			//mqttclients.Paho5BasicCertTest(*brokerAddress, *port, *certFilePath)
		}
	} else {
		// Either return an error or use a default version?
	}

}
