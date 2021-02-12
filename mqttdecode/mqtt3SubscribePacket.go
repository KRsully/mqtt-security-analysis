package mqttdecode

import (
	"encoding/binary"

	"github.com/google/gopacket"
)

var MQTT3SubscribePacket = gopacket.RegisterLayerType(
	3891,
	gopacket.LayerTypeMetadata{Name: "MQTT 3.1.1 SUBSCRIBE Packet", Decoder: gopacket.DecodeFunc(DecodeMQTT3SubscribePacket)})

type mqtt3SubscribePacket struct {
	VariableHeader mqtt3SubscribeVariableHeader
	Topics         []topicSubscription
	Contents       []byte
}

type topicSubscription struct {
	TopicString string
	TopicLength uint16
	QoS         int
}

func (layer mqtt3SubscribePacket) LayerType() gopacket.LayerType { return MQTT3SubscribePacket }
func (layer mqtt3SubscribePacket) LayerContents() []byte         { return layer.Contents }
func (layer mqtt3SubscribePacket) LayerPayload() []byte          { return nil }

func DecodeMQTT3SubscribePacket(data []byte, packet gopacket.PacketBuilder) (err error) {
	variableHeader, err := decodeMQTT3SubscribeVariableHeader(data)
	payload, err := decodeMQTT3SubscribePayload(data[2:])
	packet.AddLayer(&mqtt3SubscribePacket{variableHeader, payload, data})
	return
}

type mqtt3SubscribeVariableHeader struct {
	PacketIdentifier uint16
}

func decodeMQTT3SubscribeVariableHeader(data []byte) (header mqtt3SubscribeVariableHeader, err error) {
	header.PacketIdentifier = binary.BigEndian.Uint16(data)

	return
}

func decodeMQTT3SubscribePayload(data []byte) (topics []topicSubscription, err error) {
	pos := 0

	for pos < len(data) {
		topicString, topicLength, _ := extractUTF8String(data[pos:])
		topics = append(topics, topicSubscription{topicString, topicLength, int(data[pos+2+int(topicLength)])})
		//topicLength + 2 bytes for length value, 1 for QoS value
		pos += int(topicLength) + 3
	}

	return
}
