# MQTT Packet Capture
# Dependencies
* `github.com/google/gopacket`
* On Linux - `libpcap-dev`

# Operation
From within the repository:
* `go run main.go`: Display the network interfaces on the system
* `go run main.go [flags] <interface-name-or-IP>`: Begin capturing packets on the designated interface. If the interface IP is provided, `main.go` will attempt to resolve it to the iterface name.

## Running from Raspberry Pi
To list the local network interfaces: `sudo -E go run main.go `

To begin packet capture: `sudo -E go run main.go [flags] <interface-name-or-IP>`
## Flags:
* `-pro`: turn promiscuous mode on
* `-p <port-number>`: Port to capture on (default 1883 and 8883)
* `-sl <snapshot-length>`: Set the maximum snapshot length (default 65535)
* `-t <timeout-length>`: Set the timeout (default nonblocking)