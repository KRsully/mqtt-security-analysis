package mqttdecode

import (
	"encoding/binary"
	"log"
)

type propertyIdentifier int

const (
	payloadFormatIndicator = propertyIdentifier(iota + 1)
	messageExpiryInterval
	contentType
	_
	_
	_
	_
	_
	responseTopic
	correlationData
	_
	subscriptionIdentifier
	_
	_
	_
	_
	_
	sessionExpiryIterval
	assignedClientIdentifier
	serverKeepAlive
	_
	authenticationMethod
	authenticationData
	requestProblemInformation
	willDelayInterval
	requestResponseInformation
	ResponseInformation
	_
	serverReference
	_
	_
	reasonString
	_
	receiveMaximum
	topicAliasMaximum
	topicAlias
	maximumQoS
	retainAvailable
	userProperty
	maximumPacketSize
	wildcardSubscriptionAvailable
	subscriptionIdentifierAvailable
	sharedSubscriptionAvailable
)

func (pi propertyIdentifier) String() string {
	return [...]string{"", "Payload Format Indicator", "Message Expiry Iterval", "Content Type", "", "", "", "",
		"Response Topic", "Correlation Data", "", "Subscription Identifier", "", "", "", "",
		"", "Session Expiry Interval", "Assigned Client Identifier", "Server Keep Alive", "", "Authentication Method", "Authentication Data", "Request Problem Information",
		"Will Delay Interval", "Request Response Information", "Response Information", "", "Server Reference", "", "", "Reason String",
		"", "Receive Maximum", "Topic Alias Maximum", "Maximum QoS", "Retain Available", "User Property", "Maximum Packet Size",
		"Wildcard Subscription Available", "Subscription Identifier Available", "Shared Subscription Available"}[pi]
}

func (pi propertyIdentifier) resolveDecoder() func(data []byte) string {
	decoder := func(data []byte) string { return string(data) }
	switch pi {
	case 1, 23, 25, 36, 37, 40, 41, 42:
		//Property is a single Byte
	case 2, 17, 24, 39:
		//Property is a Four-Byte Integer
		decoder = func(data []byte) string { return string(binary.BigEndian.Uint32(data)) }
	case 3, 8, 18, 21, 26, 28, 31:
		//Property is a single UTF-8 Encoded String
		decoder = func(data []byte) string {
			propString, _, _ := extractUTF8String(data)
			return propString
		}
	case 9, 22:
		//Property is Binary Data
	case 11:
		//Property is a Variable Byte Integer
		decoder = func(data []byte) string {
			value, _, _ := decodeVariableByteInteger(data)
			return string(value)
		}
	case 19, 33, 34, 35:
		//Property is a Two-Byte Integer
		decoder = func(data []byte) string { return string(binary.BigEndian.Uint16(data)) }
	case 38:
		//Property is a pair of UTF-8 Encoded Strings
		decoder = func(data []byte) string {
			nameString, strLength, _ := extractUTF8String(data)
			valueString, _, _ := extractUTF8String(data[2+strLength:])
			return nameString + ":" + valueString
		}
	}
	return decoder
}

func (pi propertyIdentifier) resolveContentSize(data []byte) (size int) {
	size = 0
	switch pi {
	case 1, 23, 25, 36, 37, 40, 41, 42:
		//Property is a single Byte
		size = 1
	case 2, 17, 24, 39:
		//Property is a Four-Byte Integer
		size = 4
	case 3, 8, 18, 21, 26, 28, 31:
		//Property is a single UTF-8 Encoded String
		_, size, _ = extractUTF8String(data)
	case 9, 22:
		//Property is Binary Data
		// What size is it????
		size = 1
	case 11:
		//Property is a Variable Byte Integer
		_, size, _ = decodeVariableByteInteger(data)
	case 19, 33, 34, 35:
		//Property is a Two-Byte Integer
		size = 2
	case 38:
		//Property is a pair of UTF-8 Encoded Strings
		_, strLength1, _ := extractUTF8String(data)
		_, strLength2, _ := extractUTF8String(data[2+strLength1:])
		size = strLength1 + strLength2
	}

	return size
}

type MQTT5Property struct {
	Identifier   int
	PropertyName string
	Contents     []byte
	//ContentDecoder func(data []byte) string
	Length int
}

// func (property MQTT5Property) String() string {
// 	return property.ContentDecoder(property.Contents)
// }

func extractMQTT5Properties(data []byte) (properties []MQTT5Property, totalLength int) {
	propertiesLength, propertiesLengthInBytes, err := decodeVariableByteInteger(data)
	if err != nil {
		//Malformed properties length value
	}
	totalLength = propertiesLength + propertiesLengthInBytes
	data = data[propertiesLengthInBytes:]
	propertiesLength -= propertiesLengthInBytes

	for propertiesLength > 0 {
		propInteger, propIntegerLengthInBytes, _ := decodeVariableByteInteger(data)
		data = data[propIntegerLengthInBytes:]
		propIdentifier := propertyIdentifier(propInteger)
		log.Println(propIdentifier)
		propLength := propIdentifier.resolveContentSize(data)
		properties = append(properties, MQTT5Property{
			propInteger,
			propIdentifier.String(),
			data[:propLength],
			//propIdentifier.resolveDecoder(),
			propLength,
		})
		// Number of bytes to code for length and property length
		propertiesLength -= propIntegerLengthInBytes + propLength
	}

	return properties, totalLength
}
