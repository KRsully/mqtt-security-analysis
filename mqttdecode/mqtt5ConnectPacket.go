package mqttdecode

import (
	"github.com/google/gopacket"
)

var MQTT5ConnectPacket = gopacket.RegisterLayerType(
	3895,
	gopacket.LayerTypeMetadata{Name: "MQTT 5.0 CONNECT Packet", Decoder: gopacket.DecodeFunc(DecodeMQTT5ConnectPacket)})

type mqtt5ConnectPacket struct {
	VariableHeader mqtt5ConnectVariableHeader
	Payload        mqtt5ConnectPayload
	Contents       []byte
}

func (layer mqtt5ConnectPacket) LayerType() gopacket.LayerType { return MQTT5ConnectPacket }
func (layer mqtt5ConnectPacket) LayerContents() []byte         { return layer.Contents }
func (layer mqtt5ConnectPacket) LayerPayload() []byte          { return nil }

func DecodeMQTT5ConnectPacket(data []byte, packet gopacket.PacketBuilder) (err error) {
	variableHeader, err := decodeMQTT5ConnectVariableHeader(data)

	payload, err := decodeMQTT5ConnectPayload(data[variableHeader.Length:], variableHeader.ConnectFlags)

	packet.AddLayer(&mqtt5ConnectPacket{variableHeader, payload, data})
	return
}

type mqtt5ConnectVariableHeader struct {
	NameLength    int
	ProtocolName  string
	ProtocolLevel uint8
	ConnectFlags  byte
	KeepAlive     []byte
	Properties    []MQTT5Property
	Length        int
}

func decodeMQTT5ConnectVariableHeader(data []byte) (header mqtt5ConnectVariableHeader, err error) {
	header.ProtocolName, header.NameLength, _ = extractUTF8String(data)
	data = data[2+header.NameLength:]
	header.ProtocolLevel = data[0]
	header.ConnectFlags = data[1]
	header.KeepAlive = data[2:4]

	data = data[4:]
	var propertiesLength int
	header.Properties, propertiesLength = extractMQTT5Properties(data)
	//2 bytes of protocol name length value, NameLength bytes of string, 4 bytes for level, flags, and keep alive,
	//	propLengthBytes of property length value and propertiesLength of property values
	header.Length = 2 + header.NameLength + 4 + propertiesLength

	return
}

type mqtt5ConnectPayload struct {
	ClientID       string
	WillProperties []MQTT5Property
	WillTopic      string
	WillPayload    string
	WillQoS        uint16
	Username       string
	Password       string
}

func decodeMQTT5ConnectPayload(data []byte, flags byte) (payload mqtt5ConnectPayload, err error) {
	//Client Identifier --> Will Properties --> Will Topic --> Will Payload --> User Name --> Password
	var stringLength int
	payload.ClientID, stringLength, _ = extractUTF8String(data)
	if flags != 0 {
		data = data[2+stringLength:]
	}

	if flags&0x4 != 0 {
		//Will Flag
		var willPropertiesLength int
		payload.WillProperties, willPropertiesLength = extractMQTT5Properties(data)
		data = data[willPropertiesLength:]
		payload.WillTopic, stringLength, _ = extractUTF8String(data)
		data = data[2+stringLength:]
		payload.WillQoS = uint16(flags & 0x18)
		// Spec says that the payload is Binary Data?
		payload.WillPayload, stringLength, _ = extractUTF8String(data)
		data = data[2+stringLength:]
	}

	if flags&0x80 != 0 {
		//Username Flag
		payload.Username, stringLength, _ = extractUTF8String(data)
		data = data[2+stringLength:]
	}
	if flags&0x40 != 0 {
		//Password Flag
		payload.Password, stringLength, _ = extractUTF8String(data)
		data = data[2+stringLength:]
	}

	if flags&0x0 != 0 {
		//Reserved Flag, must be 0

	}

	return
}
