/*
* A little CLI to send sets of packets to a designated interface
* Right now the two tests run depending on the flags provided, but the idea is to have a bunch of tests
*	that can be chosen from or loaded from a file or something of the like.
* Will we get that far?
 */

package main

import (
	"flag"

	"github.com/KRsully/mqtt-security-analysis/mqttclients"
)

func main() {
	port := flag.String("p", "1883", "port to capture MQTT packets on")
	brokerAddress := flag.String("b", "127.0.0.1", "IP address or FQDN of the target broker")
	certFilePath := flag.String("cafile", "certs/ca.crt", "filepath to the certificate for verifying the server")
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
