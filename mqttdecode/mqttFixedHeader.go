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

func appendFlags(payload []byte, flags byte) []byte {
	/*
	*	Some packet types need to be able to check the flags in the fixed header
	*  	For now, we'll prepend the flags to the payload before we call the packet-specific decoder
	*	Yeah it's not great.
	 */
	payload = append(payload, 0)
	copy(payload[1:], payload)
	payload[0] = flags
	return payload
}

func decodeMQTTFixedHeader(data []byte, packet gopacket.PacketBuilder) (err error) {
	remainingLength, _, err := decodeVariableByteInteger(data[1:])
	ctlPacketType := decodeControlPacketType(data[0])
	flags := data[0] & 0xF
	if err != nil {
		return err
	}

	if ctlPacketType < 1 {
		return errors.New("Invalid Control Type Error")
	}

	payloadStartIndex := len(data) - remainingLength

	packet.AddLayer(&mqttFixedHeader{decodeControlPacketType(data[0]).String(),
		flags, remainingLength, data[0:payloadStartIndex], data[payloadStartIndex:]})
	//log.Printf("packet payload: %v\n", data[payloadStartIndex:])
	switch ctlPacketType {
	/*If a decoder function runs into any error, gopacket panics
	* Since there doesn't see to be a uniform way of determining if the packet is 3.1.1 or 5.0,
	* we'll just have a deferred error/panic check for every single decoder...
	 */
	case 1:
		defer func() {
			if err := recover(); err != nil {
				err = DecodeMQTT3ConnectPacket(data[payloadStartIndex:], packet)
			}
		}()
		err = DecodeMQTT5ConnectPacket(data[payloadStartIndex:], packet)
	case 2:
		defer func() {
			if err := recover(); err != nil {
				err = DecodeMQTT3ConnAckPacket(data[payloadStartIndex:], packet)
			}
		}()
		err = DecodeMQTT5ConnAckPacket(data[payloadStartIndex:], packet)
	case 3:
		defer func() {
			if err := recover(); err != nil {
				err = DecodeMQTT3PublishPacket(appendFlags(data[payloadStartIndex:], flags), packet)
			}
		}()
		err = DecodeMQTT5PublishPacket(appendFlags(data[payloadStartIndex:], flags), packet)
	case 4:
		defer func() {
			if err := recover(); err != nil {
				err = DecodeMQTT3PubAckPacket(data[payloadStartIndex:], packet)
			}
		}()
		err = DecodeMQTT5PubAckPacket(data[payloadStartIndex:], packet)
	case 5:
		defer func() {
			if err := recover(); err != nil {
				err = DecodeMQTT3PubRecPacket(data[payloadStartIndex:], packet)
			}
		}()
		err = DecodeMQTT5PubRecPacket(data[payloadStartIndex:], packet)
	case 6:
		defer func() {
			if err := recover(); err != nil {
				err = DecodeMQTT3PubRelPacket(data[payloadStartIndex:], packet)
			}
		}()
		err = DecodeMQTT5PubRelPacket(data[payloadStartIndex:], packet)
	case 7:
		defer func() {
			if err := recover(); err != nil {
				err = DecodeMQTT3PubCompPacket(data[payloadStartIndex:], packet)
			}
		}()
		err = DecodeMQTT5PubCompPacket(data[payloadStartIndex:], packet)
	case 8:
		defer func() {
			if err := recover(); err != nil {
				err = DecodeMQTT3SubscribePacket(data[payloadStartIndex:], packet)
			}
		}()
		err = DecodeMQTT5SubscribePacket(data[payloadStartIndex:], packet)
	case 9:
		defer func() {
			if err := recover(); err != nil {
				err = DecodeMQTT3SubAckPacket(data[payloadStartIndex:], packet)
			}
		}()
		err = DecodeMQTT5SubAckPacket(data[payloadStartIndex:], packet)
	case 10:
		defer func() {
			if err := recover(); err != nil {
				err = DecodeMQTT3UnsubscribePacket(data[payloadStartIndex:], packet)
			}
		}()
		err = DecodeMQTT5UnsubscribePacket(data[payloadStartIndex:], packet)
	case 11:
		defer func() {
			if err := recover(); err != nil {
				err = DecodeMQTT3UnsubAckPacket(data[payloadStartIndex:], packet)
			}
		}()
		err = DecodeMQTT5UnsubAckPacket(data[payloadStartIndex:], packet)
	case 12:
		//PINGREQ Packet
	case 13:
		//PINGRESP Packet
	case 14:
		//DISCONNECT Packet
		if len(data) > 2 {
			//Sometimes, MQTT5 disconnect packets have no variable header?
			err = DecodeMQTT5DisconnectPacket(data[payloadStartIndex:], packet)
		}
	case 15:
		//AUTH Packet (MQTT 5.0 only)
		err = DecodeMQTT5AuthPacket(data[payloadStartIndex:], packet)
	default:
		//?
	}

	return nil

}
