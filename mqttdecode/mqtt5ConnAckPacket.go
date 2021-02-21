package mqttdecode

import (
	"github.com/google/gopacket"
)

var MQTT5ConnAckPacket = gopacket.RegisterLayerType(
	3896,
	gopacket.LayerTypeMetadata{Name: "MQTT 5.0 CONNACK Packet", Decoder: gopacket.DecodeFunc(DecodeMQTT5ConnAckPacket)})

type mqtt5ConnAckPacket struct {
	VariableHeader mqtt5ConnAckVariableHeader
	Contents       []byte
}

func (layer mqtt5ConnAckPacket) LayerType() gopacket.LayerType { return MQTT5ConnAckPacket }
func (layer mqtt5ConnAckPacket) LayerContents() []byte         { return layer.Contents }
func (layer mqtt5ConnAckPacket) LayerPayload() []byte          { return nil }

func DecodeMQTT5ConnAckPacket(data []byte, packet gopacket.PacketBuilder) (err error) {
	variableHeader, err := decodeMQTT5ConnAckVariableHeader(data)

	packet.AddLayer(&mqtt5ConnAckPacket{variableHeader, data})
	return
}

type mqtt5ConnAckVariableHeader struct {
	ConnectAckFlag    uint8
	ConnectReturnCode uint8
	ReturnCodeName    string
	Properties        []MQTT5Property
}

func decodeMQTT5ConnAckVariableHeader(data []byte) (header mqtt5ConnAckVariableHeader, err error) {
	header.ConnectAckFlag = data[0]
	header.ConnectReturnCode = data[1]
	header.ReturnCodeName = resolveReasonCode(header.ConnectReturnCode)
	header.Properties, _ = extractMQTT5Properties(data[2:])
	return
}
