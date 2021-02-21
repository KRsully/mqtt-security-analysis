package mqttdecode

import (
	"encoding/binary"

	"github.com/google/gopacket"
)

var MQTT5PublishPacket = gopacket.RegisterLayerType(
	3897,
	gopacket.LayerTypeMetadata{Name: "MQTT 5.0 PUBLISH Packet", Decoder: gopacket.DecodeFunc(DecodeMQTT5PublishPacket)})

type mqtt5PublishPacket struct {
	VariableHeader mqtt5PublishVariableHeader
	Contents       []byte
	Payload        string
}

func (layer mqtt5PublishPacket) LayerType() gopacket.LayerType { return MQTT5PublishPacket }
func (layer mqtt5PublishPacket) LayerContents() []byte         { return layer.Contents }
func (layer mqtt5PublishPacket) LayerPayload() []byte          { return nil }

func DecodeMQTT5PublishPacket(data []byte, packet gopacket.PacketBuilder) (err error) {
	variableHeader, err := decodeMQTT5PublishVariableHeader(data)

	packet.AddLayer(&mqtt5PublishPacket{variableHeader, data, string(data[variableHeader.Length:])})
	return
}

type mqtt5PublishVariableHeader struct {
	TopicLength      int
	TopicString      string
	PacketIdentifier uint16
	Properties       []MQTT5Property
	Length           int
}

func decodeMQTT5PublishVariableHeader(data []byte) (header mqtt5PublishVariableHeader, err error) {
	flags := data[0]
	// Flags are 4 bits: 	0 0 0 0 DUP QoS RETAIN
	//						7 6 5 4  3  2,1   0
	QoS := flags & 0x3
	data = data[1:]
	var propertiesStartIndex int
	header.TopicString, header.TopicLength, _ = extractUTF8String(data)
	if QoS == 2 || QoS == 1 {
		header.PacketIdentifier = binary.BigEndian.Uint16(data[2+header.TopicLength:])
		propertiesStartIndex = 2 + header.TopicLength + 2
	} else {
		propertiesStartIndex = 2 + header.TopicLength
	}
	var propertiesLength int
	header.Properties, propertiesLength = extractMQTT5Properties(data[propertiesStartIndex:])

	header.Length = propertiesStartIndex + propertiesLength
	return
}
