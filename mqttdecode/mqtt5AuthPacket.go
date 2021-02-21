package mqttdecode

import (
	"github.com/google/gopacket"
)

var MQTT5AuthPacket = gopacket.RegisterLayerType(
	3907,
	gopacket.LayerTypeMetadata{Name: "MQTT 5.0 AUTH Packet", Decoder: gopacket.DecodeFunc(DecodeMQTT5AuthPacket)})

type mqtt5AuthPacket struct {
	VariableHeader mqtt5AuthVariableHeader
	Contents       []byte
}

func (layer mqtt5AuthPacket) LayerType() gopacket.LayerType { return MQTT5AuthPacket }
func (layer mqtt5AuthPacket) LayerContents() []byte         { return layer.Contents }
func (layer mqtt5AuthPacket) LayerPayload() []byte          { return nil }

func DecodeMQTT5AuthPacket(data []byte, packet gopacket.PacketBuilder) (err error) {
	variableHeader, err := decodeMQTT5AuthVariableHeader(data)

	packet.AddLayer(&mqtt5AuthPacket{variableHeader, data})
	return
}

type mqtt5AuthVariableHeader struct {
	ReasonCode     byte
	ReasonCodeName string
	Properties     []MQTT5Property
}

func decodeMQTT5AuthVariableHeader(data []byte) (header mqtt5AuthVariableHeader, err error) {
	header.ReasonCode = data[0]
	header.ReasonCodeName = resolveReasonCode(header.ReasonCode)
	header.Properties, _ = extractMQTT5Properties(data[1:])
	return
}
