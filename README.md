# mqtt-security-analysis
# Dependencies
* `github.com/google/gopacket`
* On Linux - `libpcap-dev`

# Operation
From within the repository:
* `go run main.go`: Display the network interfaces on the system
* `go run main.go <interface-name-or-IP>`: Begin capturing packets on the designated interface. If the interface IP is provided, `main.go` will attempt to resolve it to the iterface name.

## Running from Raspberry Pi
To list the local network interfaces: `sudo -E go run main.go `

To begin packet capture: `sudo -E go run main.go <interface-name-or-IP>`
## Flags:
* `-pro`: turn promiscuous mode on
* `-p`: Port to capture on (default 1883 and 8883)
* `-sl`: Set the maximum snapshot length (default 65535)
* `-t`: Set the timeout (default nonblocking)
