package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"time"

	"github.com/google/gopacket/pcap"
)

const mqttPort = 1883
const mqttSecPort = 8883

type pcapFlags struct {
	port        *int
	promiscuity *bool
	snaplen     *int
	timeout     *time.Duration
	tail        []string
}

//Use pcap to print the name, descriptor, and IP addresses of the machine's network interfaces to the console.
func listAvailableInterfaces() {
	devices, err := pcap.FindAllDevs()
	if err != nil {
		log.Panic(err)
	}

	for _, device := range devices {
		fmt.Printf("%s (%s):\n", device.Name, device.Description)
		for _, address := range device.Addresses {
			fmt.Printf("%12s %s \n", "-", address.IP)
		}
	}
}

//use nmap to determine devices on the network that have port 1883 (mqtt) and/or port 8883 (mqtt-ssl) open.
func listNetworkedMqttListeners() {

}

//Given the
func resolveIPToName(targetIP net.IP) string {
	devices, err := pcap.FindAllDevs()
	if err != nil {
		log.Panic(err)
	}

	for _, device := range devices {
		for _, address := range device.Addresses {
			if address.IP.Equal(targetIP) {
				return device.Name
			}
		}
	}

	return ""
}

//Create an inactive pcap handle given the user-supplied parameters or the defaults.
//	deviceName is the name of the interface to create the handle on
//	Returns a pointer to the InactiveHandle object
func createInactiveHandle(deviceName string, flags *pcapFlags) *pcap.InactiveHandle {
	inactive, err := pcap.NewInactiveHandle(deviceName)
	if err != nil {
		log.Fatal(err)
	}
	defer inactive.CleanUp()

	if err = inactive.SetPromisc(*flags.promiscuity); err != nil {
		log.Fatal(err)
	} else if err = inactive.SetSnapLen(*flags.snaplen); err != nil {
		log.Fatal(err)
	} else if err = inactive.SetTimeout(*flags.timeout); err != nil {
		log.Fatal(err)
	}

	log.Printf("Created inactive pcap handle to capture traffic on device : %s.\n%30s %t\n%30s %d\n%30s %v\n",
		deviceName, "Promiscuous:", *flags.promiscuity, "Snapshot Length:", *flags.snaplen, "Timeout:", *flags.timeout)

	return inactive
}

//Parse the flags provided through the command line
//	Returns a pointer to a struct with the values of the flags
func parseFlags() *pcapFlags {
	flags := pcapFlags{
		port:        flag.Int("p", mqttPort, "port to capture MQTT packets on"),
		promiscuity: flag.Bool("pro", false, "promiscuous mode on"),
		snaplen:     flag.Int("sl", 512, "maximum bytes per packet capture"), //Is MQTT packet max 512 bytes?
		timeout:     flag.Duration("t", pcap.BlockForever, "time to wait for additional packets after packet capture before returning"),
	}

	flag.Parse()

	flags.tail = flag.Args()

	return &flags
}

//Driver for the packet capturing
func captureMQTTPackets() {
	flags := parseFlags()

	if len(flags.tail) < 1 {
		log.Fatal("No interface name or IP address supplied.")
	}

	deviceName := flags.tail[0]
	targetIP := net.ParseIP(deviceName)
	if targetIP != nil {
		deviceName = resolveIPToName(targetIP)
	}
	inactive := createInactiveHandle(deviceName, flags)
	log.Printf("Inactive Handle %v Port %d", inactive, *flags.port) //temp
}

func main() {
	if len(os.Args) < 2 {
		listAvailableInterfaces()
	} else {
		captureMQTTPackets()
	}
}
