package mqttdecode

import (
	"log"

	"github.com/google/gopacket"
)

var MQTT3ConnectPacket = gopacket.RegisterLayerType(
	3884,
	gopacket.LayerTypeMetadata{Name: "MQTT 3.1.1 CONNECT Packet", Decoder: gopacket.DecodeFunc(DecodeMQTT3ConnectPacket)})

type mqtt3ConnectPacket struct {
	VariableHeader mqtt3ConnectVariableHeader
	Payload        mqtt3ConnectPayload
	Contents       []byte
}

func (layer mqtt3ConnectPacket) LayerType() gopacket.LayerType { return MQTT3ConnectPacket }
func (layer mqtt3ConnectPacket) LayerContents() []byte         { return layer.Contents }
func (layer mqtt3ConnectPacket) LayerPayload() []byte          { return nil }

func DecodeMQTT3ConnectPacket(data []byte, packet gopacket.PacketBuilder) (err error) {
	log.Printf("Connect Data: %v", data)
	variableHeader, err := decodeMQTT3ConnectVariableHeader(data)
	payload, err := decodeMQTT3ConnectPayload(data[variableHeader.Length:], variableHeader.ConnectFlags)

	packet.AddLayer(&mqtt3ConnectPacket{variableHeader, payload, data})
	return
}

type mqtt3ConnectVariableHeader struct {
	NameLength    uint16
	ProtocolName  string
	ProtocolLevel uint8
	ConnectFlags  byte
	KeepAlive     []byte
	Length        uint16
}

func decodeMQTT3ConnectVariableHeader(data []byte) (header mqtt3ConnectVariableHeader, err error) {
	header.ProtocolName, header.NameLength, _ = extractUTF8String(data)
	header.ProtocolLevel = data[header.NameLength]
	header.ConnectFlags = data[header.NameLength+1]
	header.KeepAlive = data[header.NameLength+2 : header.NameLength+4]
	//2 bytes of protocol name length, NameLength bytes of string, and 4 bytes for level, flags, and keep alive
	header.Length = 2 + header.NameLength + 4

	return
}

type mqtt3ConnectPayload struct {
	ClientID    string
	WillTopic   string
	WillMessage string
	WillQoS     uint16
	Username    string
	Password    string
}

func decodeMQTT3ConnectPayload(data []byte, flags byte) (payload mqtt3ConnectPayload, err error) {
	//Client Identifier --> Will Retain --> Will Message --> User Name --> Password
	log.Printf("Payload: %v", data)
	var stringLength uint16
	payload.ClientID, stringLength, _ = extractUTF8String(data)
	log.Printf("ClientID: %s", payload.ClientID)
	if flags != 0 {
		data = data[stringLength+1:]
	}

	if flags&0x4 == 1 {
		//Will Flag
		payload.WillTopic, stringLength, _ = extractUTF8String(data)
		data = data[stringLength+1:]
		payload.WillQoS = uint16(flags & 18)
		payload.WillMessage, stringLength, _ = extractUTF8String(data)
		data = data[stringLength+1:]
		log.Printf("will")
	}

	if flags&0x80 == 1 {
		//Username Flag
		payload.Username, stringLength, _ = extractUTF8String(data)
		data = data[stringLength+1:]
		log.Printf("uname")
	}
	if flags&0x40 == 1 {
		//Password Flag
		payload.Password, stringLength, _ = extractUTF8String(data)
		data = data[stringLength+1:]
		log.Printf("passwd")
	}

	if flags&0x0 != 0 {
		//Reserved Flag, must be 0

	}

	return
}
