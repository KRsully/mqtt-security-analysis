/*
* This will eventually have a bunch of different tests/some sort of interpreter that runs tests with the Paho 3 library
* Modified from https://github.com/eclipse/paho.mqtt.golang/blob/master/cmd/simple/main.go
 */
package mqttclients

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	paho "github.com/eclipse/paho.mqtt.golang"
)

var msgHandler paho.MessageHandler = func(client paho.Client, msg paho.Message) {
	fmt.Printf("[%s]: %s\n", msg.Topic(), msg.Payload())
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

func Paho3BasicPubSubTest(brokerAddress string, port string) {
	//paho.DEBUG = log.New(os.Stdout, "", 0)
	//paho.ERROR = log.New(os.Stdout, "", 0)
	opts := paho.NewClientOptions().AddBroker("mqtt://" + brokerAddress + ":" + port).SetClientID("paho3BasicTestCID")
	opts.SetKeepAlive(2 * time.Second)
	opts.SetDefaultPublishHandler(msgHandler)
	opts.SetPingTimeout(1 * time.Second)
	opts.SetUsername("paho3BasicTestUser")

	c := paho.NewClient(opts)

	if token := c.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	if token := c.Subscribe("paho3BasicTestTopic", 0, nil); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
		os.Exit(1)
	}

	for i := 0; i < 5; i++ {
		text := fmt.Sprintf("This is msg #%d!", i)
		token := c.Publish("paho3BasicTestTopic", 0, false, text)
		token.Wait()
	}

	time.Sleep(2 * time.Second)

	if token := c.Unsubscribe("paho3BasicTestTopic"); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
		os.Exit(1)
	}

	c.Disconnect(250)
}

func Paho3BasicCertTest(brokerAddress string, port string, certFilePath string) {
	//paho.DEBUG = log.New(os.Stdout, "", 0)
	//paho.ERROR = log.New(os.Stdout, "", 0)
	opts := paho.NewClientOptions().AddBroker("mqtts://" + brokerAddress + ":" + port).SetClientID("paho3BasicCertTest")
	opts.SetKeepAlive(2 * time.Second)
	opts.SetDefaultPublishHandler(msgHandler)
	opts.SetPingTimeout(1 * time.Second)

	certFile, err := ioutil.ReadFile(certFilePath)
	if err != nil {
		log.Printf("Error reading file %s: %s", certFile, err.Error())
	}
	certPool := x509.NewCertPool()
	if err != nil {
		certPool.AppendCertsFromPEM(certFile)
	}

	// var cert tls.Certificate
	// cert.Leaf, err = x509.ParseCertificate(certFile)

	opts.SetTLSConfig(&tls.Config{
		RootCAs:            certPool,
		ClientAuth:         tls.NoClientCert,
		Certificates:       []tls.Certificate{},
		InsecureSkipVerify: true,
	})
	opts.SetUsername("paho3BasicCertTest")

	c := paho.NewClient(opts)

	if token := c.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	if token := c.Subscribe("paho3BasicCertTest", 0, nil); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
		os.Exit(1)
	}

	for i := 0; i < 5; i++ {
		text := fmt.Sprintf("This is msg #%d!", i)
		token := c.Publish("paho3BasicCertTest", 0, false, text)
		token.Wait()
	}

	time.Sleep(2 * time.Second)

	if token := c.Unsubscribe("paho3BasicCertTest"); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
		os.Exit(1)
	}

	c.Disconnect(250)
}
