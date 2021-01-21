package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
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

/* func resolveIPToName(targetIP net.IP) string {
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
} */

//Create an inactive pcap handle given the user-supplied parameters or the defaults.
//	deviceName is the name of the interface to create the handle on
//	Returns a pointer to the InactiveHandle object
func createInactiveHandle(deviceName string, flags *pcapFlags) *pcap.InactiveHandle {
	inactive, err := pcap.NewInactiveHandle(deviceName)
	if err != nil {
		log.Fatal(err)
	}

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
		snaplen:     flag.Int("sl", 65535, "maximum bytes per packet capture"), //Is MQTT packet max 512 bytes?
		timeout:     flag.Duration("t", pcap.BlockForever, "time to wait for additional packets after packet capture before returning (min. nanoseconds)"),
	}

	flag.Parse()

	flags.tail = flag.Args()

	return &flags
}

//Parse the given interface string to determine if pcap can find a valid interface
//	Receives a string which may be the IP address or the name of the network interface
//	Returns a string with the name of one of the machine's interfaces
func parseInterface(deviceIdentifier string) string {

	devices, err := pcap.FindAllDevs()
	if err != nil {
		log.Panic(err)
	}

	for _, device := range devices {
		if device.Name == deviceIdentifier {
			return device.Name
		}
		for _, address := range device.Addresses {
			targetIP := net.ParseIP(deviceIdentifier)
			if address.IP.Equal(targetIP) {
				return device.Name
			}
		}
	}

	log.Fatal("Cannot resolve interface name")

	return ""
}

//Driver for the packet capturing
func captureMQTTPackets() {
	flags := parseFlags()

	if len(flags.tail) < 1 {
		log.Fatal("No interface name or IP address supplied.")
	}

	deviceName := parseInterface(flags.tail[0])

	inactive := createInactiveHandle(deviceName, flags)

	pcapHandle, err := inactive.Activate()
	if err != nil {
		log.Panic(err)
	}

	inactive.CleanUp()

	defer pcapHandle.Close()

	//pcapHandle.SetBPFFilter("tcp and port " + strconv.Itoa(*flags.port))
	pcapHandle.SetBPFFilter("tcp and port 1883 or port 8883")

	packetSource := gopacket.NewPacketSource(pcapHandle, layers.LayerTypeEthernet)
	for packet := range packetSource.Packets() {
		fmt.Println(packet)
	}
}

func main() {
	if len(os.Args) < 2 {
		listAvailableInterfaces()
	} else {
		captureMQTTPackets()
	}
}
