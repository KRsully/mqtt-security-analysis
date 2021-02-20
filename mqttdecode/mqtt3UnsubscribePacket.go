package mqttdecode

import (
	"encoding/binary"

	"github.com/google/gopacket"
)

var MQTT3UnsubscribePacket = gopacket.RegisterLayerType(
	3893,
	gopacket.LayerTypeMetadata{Name: "MQTT 3.1.1 UNSUBSCRIBE Packet", Decoder: gopacket.DecodeFunc(DecodeMQTT3UnsubscribePacket)})

type topic struct {
	topicString string
	topicLength int
}

type mqtt3UnsubscribePacket struct {
	VariableHeader mqtt3SubAckVariableHeader
	TopicFilters   []string
	Contents       []byte
}

func (layer mqtt3UnsubscribePacket) LayerType() gopacket.LayerType { return MQTT3UnsubscribePacket }
func (layer mqtt3UnsubscribePacket) LayerContents() []byte         { return layer.Contents }
func (layer mqtt3UnsubscribePacket) LayerPayload() []byte          { return nil }

func DecodeMQTT3UnsubscribePacket(data []byte, packet gopacket.PacketBuilder) (err error) {
	variableHeader, err := decodeMQTT3UnsubscribeVariableHeader(data)

	packet.AddLayer(&mqtt3SubAckPacket{variableHeader, data[2:], data})
	return
}

type mqtt3UnsubscribeVariableHeader struct {
	PacketIdentifier uint16
}

func decodeMQTT3UnsubscribeVariableHeader(data []byte) (header mqtt3SubAckVariableHeader, err error) {
	header.PacketIdentifier = binary.BigEndian.Uint16(data)

	return
}

func decodeMQTT3UnsubscribePayload(data []byte) (topics []topic, err error) {
	pos := 0

	for pos < len(data) {
		topicString, topicLength, _ := extractUTF8String(data[pos:])
		topics = append(topics, topic{topicString, topicLength})
		//topicLength + 2 bytes for length value
		pos += int(topicLength) + 2
	}

	return
}
