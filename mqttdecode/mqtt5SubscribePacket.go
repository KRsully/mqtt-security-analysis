package mqttdecode

import (
	"encoding/binary"

	"github.com/google/gopacket"
)

var MQTT5SubscribePacket = gopacket.RegisterLayerType(
	3902,
	gopacket.LayerTypeMetadata{Name: "MQTT 5.0 SUBSCRIBE Packet", Decoder: gopacket.DecodeFunc(DecodeMQTT5SubscribePacket)})

type mqtt5SubscribePacket struct {
	VariableHeader mqtt5SubscribeVariableHeader
	Topics         []mqtt5TopicSubscription
	Contents       []byte
}

type mqtt5TopicSubscription struct {
	TopicString string
	Options     byte
}

func (layer mqtt5SubscribePacket) LayerType() gopacket.LayerType { return MQTT5SubscribePacket }
func (layer mqtt5SubscribePacket) LayerContents() []byte         { return layer.Contents }
func (layer mqtt5SubscribePacket) LayerPayload() []byte          { return nil }

func DecodeMQTT5SubscribePacket(data []byte, packet gopacket.PacketBuilder) (err error) {
	variableHeader, err := decodeMQTT5SubscribeVariableHeader(data)
	payload, err := decodeMQTT5SubscribePayload(data[variableHeader.Length:])
	packet.AddLayer(&mqtt5SubscribePacket{variableHeader, payload, data})
	return
}

type mqtt5SubscribeVariableHeader struct {
	PacketIdentifier uint16
	Properties       []MQTT5Property
	Length           int
}

func decodeMQTT5SubscribeVariableHeader(data []byte) (header mqtt5SubscribeVariableHeader, err error) {
	header.PacketIdentifier = binary.BigEndian.Uint16(data)
	var propertiesLength int
	header.Properties, propertiesLength = extractMQTT5Properties(data[2:])
	header.Length = 2 + propertiesLength
	return
}

func decodeMQTT5SubscribePayload(data []byte) (topics []mqtt5TopicSubscription, err error) {
	pos := 0
	for pos < len(data) {
		topicString, topicLength, _ := extractUTF8String(data[pos:])
		topics = append(topics, mqtt5TopicSubscription{topicString, data[pos+2+topicLength]})
		//topicLength + 2 bytes for length value, 1 for options
		pos += topicLength + 3
	}

	return
}
