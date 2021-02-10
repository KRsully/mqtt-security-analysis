package mqttdecode

import (
	"errors"

	"github.com/google/gopacket"
)

var mqttFixedHeaderLayer = gopacket.RegisterLayerType(
	3883,
	gopacket.LayerTypeMetadata{Name: "MQTT Fixed Header", Decoder: gopacket.DecodeFunc(DecodeMQTTFixedHeader)})

type mqttFixedHeader struct {
	ControlPacketType string
	Flags             byte
	RemainingLength   int
	Contents          []byte
	Payload           []byte
}

func (layer mqttFixedHeader) LayerType() gopacket.LayerType { return MQTT3ConnectPacket }
func (layer mqttFixedHeader) LayerContents() []byte         { return layer.Contents }
func (layer mqttFixedHeader) LayerPayload() []byte          { return nil }

func decodeControlPacketType(header byte) controlPacketType {
	//MQTT control packet type is determined by the value of the 4 highest bits of the packet's first byte
	return controlPacketType((header & 0xF0) >> 4)
}

func DecodeMQTTFixedHeader(data []byte, packet gopacket.PacketBuilder) (err error) {
	remainingLength, err := decodeVariableByteInteger(data[1:])
	ctlPacketType := decodeControlPacketType(data[0])

	if err != nil {
		return err
	}

	if ctlPacketType < 1 || ctlPacketType > 15 {
		return errors.New("Invalid Control Type Error")
	}

	packet.AddLayer(&mqttFixedHeader{decodeControlPacketType(data[0]).String(),
		data[0] & 0xF, remainingLength, data[0:(remainingLength/128 + 1)], data[(remainingLength/128 + 2):]})

	//If the packet is a PUBLISH with QoS>0, or of type PUBACK through UNSUBACK, it has a variable header
	// if (ctlPacketType == PUBLISH && data[0]&0x3 != 0) || (PUBLISH < ctlPacketType && ctlPacketType < PINGREQ) {
	// 	return packet.NextDecoder(MQTTVariableHeaderLayer)
	// }
	switch ctlPacketType {
	case 1:
		DecodeMQTT3ConnectPacket(data[2:], packet)
	}

	return nil

}
