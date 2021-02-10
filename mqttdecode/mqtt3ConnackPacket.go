package mqttdecode

import (
	"github.com/google/gopacket"
)

var MQTT3ConnAckPacket = gopacket.RegisterLayerType(
	3885,
	gopacket.LayerTypeMetadata{Name: "MQTT 3.1.1 CONNACK Packet", Decoder: gopacket.DecodeFunc(decodeMQTT3ConnAckPacket)})

type mqtt3ConnAckPacket struct {
	VariableHeader mqtt3ConnAckVariableHeader
	Contents       []byte
}

func (layer mqtt3ConnAckPacket) LayerType() gopacket.LayerType { return MQTT3ConnAckPacket }
func (layer mqtt3ConnAckPacket) LayerContents() []byte         { return layer.Contents }
func (layer mqtt3ConnAckPacket) LayerPayload() []byte          { return nil }

func decodeMQTT3ConnAckPacket(data []byte, packet gopacket.PacketBuilder) (err error) {
	variableHeader, err := decodeMQTT3ConnectVariableHeader(data)
	payload, err := decodeMQTT3ConnectPayload(data[variableHeader.Length+1:], variableHeader.ConnectFlags)

	packet.AddLayer(&mqtt3ConnectPacket{variableHeader, payload, data})
	return
}

type mqtt3ConnAckVariableHeader struct {
	ConnectAckFlag     uint8
	ConnectReturnCode  uint8
	ReturnCodeResponse string
}

func decodeMQTT3ConnAckVariableHeader(data []byte) (header mqtt3ConnAckVariableHeader, err error) {
	header.ConnectAckFlag = data[0]
	header.ConnectReturnCode = data[1]
	switch header.ConnectReturnCode {
	case 0:
		header.ReturnCodeResponse = "Connection Accepted"
	case 1:
		header.ReturnCodeResponse = "Connection Refused - unacceptable protocol version"
	case 2:
		header.ReturnCodeResponse = "Connection Refused - identifier rejected"
	case 3:
		header.ReturnCodeResponse = "Connection Refused - server unavailable"
	case 4:
		header.ReturnCodeResponse = "Connection Refused - bad username or password"
	case 5:
		header.ReturnCodeResponse = "Connection Refused - not authorized"
	default:
		header.ReturnCodeResponse = "Unknown"
	}

	return
}
