package mqttdecode

import (
	"encoding/binary"

	"github.com/google/gopacket"
)

var MQTT5UnsubscribePacket = gopacket.RegisterLayerType(
	3904,
	gopacket.LayerTypeMetadata{Name: "MQTT 5.0 UNSUBSCRIBE Packet", Decoder: gopacket.DecodeFunc(DecodeMQTT5UnsubscribePacket)})

type topic struct {
	topicString string
	topicLength int
}

type mqtt5UnsubscribePacket struct {
	VariableHeader mqtt5UnsubscribeVariableHeader
	TopicFilters   []string
	Contents       []byte
}

func (layer mqtt5UnsubscribePacket) LayerType() gopacket.LayerType { return MQTT5UnsubscribePacket }
func (layer mqtt5UnsubscribePacket) LayerContents() []byte         { return layer.Contents }
func (layer mqtt5UnsubscribePacket) LayerPayload() []byte          { return nil }

func DecodeMQTT5UnsubscribePacket(data []byte, packet gopacket.PacketBuilder) (err error) {
	variableHeader, err := decodeMQTT5UnsubscribeVariableHeader(data)
	topics, err := decodeMQTT5UnsubscribePayload(data[variableHeader.Length:])
	packet.AddLayer(&mqtt5UnsubscribePacket{variableHeader, topics, data})
	return
}

type mqtt5UnsubscribeVariableHeader struct {
	PacketIdentifier uint16
	Properties       []MQTT5Property
	Length           int
}

func decodeMQTT5UnsubscribeVariableHeader(data []byte) (header mqtt5UnsubscribeVariableHeader, err error) {
	header.PacketIdentifier = binary.BigEndian.Uint16(data)
	var propertiesLength int
	header.Properties, propertiesLength = extractMQTT5Properties(data[2:])
	header.Length = 2 + propertiesLength
	return
}

func decodeMQTT5UnsubscribePayload(data []byte) (topics []string, err error) {
	pos := 0

	for pos < len(data) {
		topicString, topicLength, _ := extractUTF8String(data[pos:])
		topics = append(topics, topicString)
		//topicLength + 2 bytes for length value
		pos += topicLength + 2
	}

	return
}
