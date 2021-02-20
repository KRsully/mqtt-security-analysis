package mqttdecode

import (
	"errors"

	"github.com/google/gopacket"
)

var MQTTFixedHeaderLayer = gopacket.RegisterLayerType(
	3883,
	gopacket.LayerTypeMetadata{Name: "MQTT Fixed Header", Decoder: gopacket.DecodeFunc(decodeMQTTFixedHeader)})

type mqttFixedHeader struct {
	ControlPacketType string
	Flags             byte
	RemainingLength   int
	Contents          []byte
	Payload           []byte
}

func (layer mqttFixedHeader) LayerType() gopacket.LayerType { return MQTTFixedHeaderLayer }
func (layer mqttFixedHeader) LayerContents() []byte         { return layer.Contents }
func (layer mqttFixedHeader) LayerPayload() []byte          { return nil }

func decodeControlPacketType(header byte) controlPacketType {
	//MQTT control packet type is determined by the value of the 4 highest bits of the packet's first byte
	return controlPacketType((header & 0xF0) >> 4)
}

func decodeMQTTFixedHeader(data []byte, packet gopacket.PacketBuilder) (err error) {
	remainingLength, _, err := decodeVariableByteInteger(data[1:])
	ctlPacketType := decodeControlPacketType(data[0])

	if err != nil {
		return err
	}

	if ctlPacketType < 1 {
		return errors.New("Invalid Control Type Error")
	}

	payloadStartIndex := len(data) - remainingLength

	packet.AddLayer(&mqttFixedHeader{decodeControlPacketType(data[0]).String(),
		data[0] & 0xF, remainingLength, data[0:payloadStartIndex], data[payloadStartIndex:]})
	//log.Printf("packet payload: %v\n", data[payloadStartIndex:])
	switch ctlPacketType {
	/*If a decoder function runs into any error, gopacket panics
	* Since there doesn't see to be a uniform way of determining if the packet is 3.1.1 or 5.0,
	* we'll just have a deferred error/panic check for every single decoder...
	 */
	case 1:
		defer func() {
			if err := recover(); err != nil {
				err = DecodeMQTT5ConnectPacket(data[payloadStartIndex:], packet)
			}
		}()
		err = DecodeMQTT3ConnectPacket(data[payloadStartIndex:], packet)
	case 2:
		err = DecodeMQTT3ConnAckPacket(data[payloadStartIndex:], packet)
		if err != nil {
			//DecodeMQTT5
		}
	case 3:
		err = DecodeMQTT3PublishPacket(data[payloadStartIndex:], packet)
		if err != nil {
			//DecodeMQTT5
		}
	case 4:
		err = DecodeMQTT3PubAckPacket(data[payloadStartIndex:], packet)
		if err != nil {
			//DecodeMQTT5
		}
	case 5:
		err = DecodeMQTT3PubRecPacket(data[payloadStartIndex:], packet)
		if err != nil {
			//DecodeMQTT5
		}
	case 6:
		err = DecodeMQTT3PubRelPacket(data[payloadStartIndex:], packet)
		if err != nil {
			//DecodeMQTT5
		}
	case 7:
		err = DecodeMQTT3PubCompPacket(data[payloadStartIndex:], packet)
		if err != nil {
			//DecodeMQTT5
		}
	case 8:
		err = DecodeMQTT3SubscribePacket(data[payloadStartIndex:], packet)
		if err != nil {
			//DecodeMQTT5
		}
	case 9:
		err = DecodeMQTT3SubAckPacket(data[payloadStartIndex:], packet)
		if err != nil {
			//DecodeMQTT5
		}
	case 10:
		err = DecodeMQTT3UnsubscribePacket(data[payloadStartIndex:], packet)
		if err != nil {
			//DecodeMQTT5
		}
	case 11:
		err = DecodeMQTT3UnsubAckPacket(data[payloadStartIndex:], packet)
		if err != nil {
			//DecodeMQTT5
		}
	case 12:
		//PINGREQ Packet
	case 13:
		//PINGRESP Packet
	case 14:
		//DISCONNECT Packet
	case 15:
		//AUTH Packet (MQTT 5.0 only)
	default:
	}

	return nil

}
