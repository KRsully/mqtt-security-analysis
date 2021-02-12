package mqttdecode

import (
	"encoding/binary"

	"github.com/google/gopacket"
)

var MQTT3PublishPacket = gopacket.RegisterLayerType(
	3886,
	gopacket.LayerTypeMetadata{Name: "MQTT 3.1.1 PUBLISH Packet", Decoder: gopacket.DecodeFunc(DecodeMQTT3PublishPacket)})

type mqtt3PublishPacket struct {
	VariableHeader mqtt3PublishVariableHeader
	Contents       []byte
	Payload        string
}

func (layer mqtt3PublishPacket) LayerType() gopacket.LayerType { return MQTT3SubAckPacket }
func (layer mqtt3PublishPacket) LayerContents() []byte         { return layer.Contents }
func (layer mqtt3PublishPacket) LayerPayload() []byte          { return nil }

func DecodeMQTT3PublishPacket(data []byte, packet gopacket.PacketBuilder) (err error) {
	variableHeader, err := decodeMQTT3PublishVariableHeader(data)

	packet.AddLayer(&mqtt3PublishPacket{variableHeader, data, string(data[2+variableHeader.TopicLength:])})
	return
}

type mqtt3PublishVariableHeader struct {
	TopicLength      uint16
	TopicString      string
	PacketIdentifier uint16
}

func decodeMQTT3PublishVariableHeader(data []byte) (header mqtt3PublishVariableHeader, err error) {
	header.TopicString, header.TopicLength, _ = extractUTF8String(data)
	header.PacketIdentifier = binary.BigEndian.Uint16(data[int(header.TopicLength):])
	return
}
