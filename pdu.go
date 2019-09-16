package gtp

import (
	"encoding/binary"
	"fmt"
)

// V2MessageType represents possible GTPv1 and GTPv2 message type values
type V2MessageType uint8

// GTPv2 MessageTypes
const (
	EchoRequest                                V2MessageType = 1
	EchoResponse                                             = 2
	CreateSessionRequest                                     = 32
	CreateSessionResponse                                    = 33
	ModifyBearerRequest                                      = 34
	ModifyBearerResponse                                     = 35
	DeleteSessionRequest                                     = 36
	DeleteSessionResponse                                    = 37
	RemoteUEReportNotification                               = 40
	RemoteUEReportAcknowlegement                             = 41
	ChangeNotificationRequest                                = 38
	ChangeNotificationResponse                               = 39
	ModifyBearerCommand                                      = 64
	ModifyBearerFailureIndication                            = 65
	DeleteBearerCommand                                      = 66
	DeleteBearerFailureIndication                            = 67
	BearerResourceCommand                                    = 68
	BearerResourceFailureIndication                          = 69
	DownlinkDataNotificationFailureIndication                = 70
	TraceSessionActivation                                   = 71
	TraceSessionDeactivation                                 = 72
	StopPagingIndication                                     = 73
	CreateBearerRequest                                      = 95
	CreateBearerResponse                                     = 96
	UpdateBearerRequest                                      = 97
	UpdateBearerResponse                                     = 98
	DeleteBearerRequest                                      = 99
	DeleteBearerResponse                                     = 100
	DeletePDNConnectionSetRequest                            = 101
	DeletePDNConnectionSetResponse                           = 102
	PGWDownlinkTriggeringNotification                        = 103
	PGWDownlinkTriggeringAcknowledge                         = 104
	IdentificationRequest                                    = 128
	IdentificationResponse                                   = 129
	ContextRequest                                           = 130
	ContextResponse                                          = 131
	ContextAcknowledge                                       = 132
	ForwardRelocationRequest                                 = 133
	ForwardRelocationResponse                                = 134
	ForwardRelocationCompleteNotification                    = 135
	ForwardRelocationCompleteAcknowledge                     = 136
	ForwardAccessContextNotification                         = 137
	ForwardAccessContextAcknowledge                          = 138
	RelocationCancelRequest                                  = 139
	RelocationCancelResponse                                 = 140
	ConfigurationTransferTunnel                              = 141
	DetachNotification                                       = 149
	DetachAcknowledge                                        = 150
	CSPagingIndication                                       = 151
	RANInformationRelay                                      = 152
	AlertMMENotification                                     = 153
	AlertMMEAcknowledge                                      = 154
	UEActivityNotification                                   = 155
	UEActivityAcknowledge                                    = 156
	ISRStatusIndication                                      = 157
	UERegistrationQueryRequest                               = 158
	UERegistrationQueryResponse                              = 159
	CreateForwardingTunnelRequest                            = 160
	CreateForwardingTunnelResponse                           = 161
	SuspendNotification                                      = 162
	SuspendAcknowledge                                       = 163
	ResumeNotification                                       = 164
	ResumeAcknowledge                                        = 165
	CreateIndirectDataForwardingTunnelRequest                = 166
	CreateIndirectDataForwardingTunnelResponse               = 167
	DeleteIndirectDataForwardingTunnelRequest                = 168
	DeleteIndirectDataForwardingTunnelResponse               = 169
	ReleaseAccessBearersRequest                              = 170
	ReleaseAccessBearersResponse                             = 171
	DownlinkDataNotification                                 = 176
	DownlinkDataNotificationAcknowledge                      = 177
	PGWRestartNotification                                   = 179
	PGWRestartNotificationAcknowledge                        = 180
	UpdatePDNConnectionSetRequest                            = 200
	UpdatePDNConnectionSetResponse                           = 201
)

// V2PDU represents a GTPv2 PDU.  Version field is omitted because it is always '2'.  TEID
// is set to 0 if TEIDFieldIsPresent is false.  Similarly, Priority is set to 0 if
// PriorityFieldIsPresent is false.  TotalLength includes complete header length, and body
// length, but does not include a piggybacked message length if IsCarryingPiggybackedPDU
// is true.  The SequenceNumber is actually a uint24 value.  Priority is actually a uint4 value
type V2PDU struct {
	IsCarryingPiggybackedPDU bool
	TEIDFieldIsPresent       bool
	PriorityFieldIsPresent   bool
	Type                     V2MessageType
	TotalLength              uint16
	TEID                     uint32
	SequenceNumber           uint32
	Priority                 uint8
	InformationElements      []*V2IE
}

// DecodeV2PDU decodes a stream of bytes that contain either exactly one well-formed
// GTPv2 PDU, or two GTPv2 PDUs when the piggyback flag on the first is set to true.
// Returns an error if the stream cannot be decoded into one or two PDUs.
func DecodeV2PDU(stream []byte) (pdu *V2PDU, piggybackedPdu *V2PDU, err error) {
	piggybackedPdu = nil

	if len(stream) < 8 {
		return nil, nil, fmt.Errorf("Stream length (%d) too short for a GTPv2 PDU", len(stream))
	}

	if (stream[0] >> 5) != 2 {
		return nil, nil, fmt.Errorf("GTPv2 PDU version should be 2, but in stream, it is (%d)", (stream[0] >> 5))
	}

	hasPiggybackedPdu := false
	if (stream[0] & 0x10) == 0x10 {
		hasPiggybackedPdu = true
	}

	msgLengthFieldValue := binary.BigEndian.Uint16(stream[2:4])
	totalPduLength := msgLengthFieldValue + 4

	if len(stream) < int(totalPduLength) {
		return nil, nil, fmt.Errorf("GTPv2 PDU length field is (%d), so total length should be (%d), but stream length is (%d)", msgLengthFieldValue, totalPduLength, len(stream))
	}

	if !hasPiggybackedPdu {
		if len(stream) > int(totalPduLength) {
			return nil, nil, fmt.Errorf("GTPv2 PDU length field is (%d), so total length should be (%d), but stream length is (%d)", msgLengthFieldValue, totalPduLength, len(stream))
		}
	} else {
		piggybackedPduStream := stream[totalPduLength:]

		if (piggybackedPduStream[0] & 0x10) != 0 {
			return nil, nil, fmt.Errorf("GTPv2 PDU has piggybacked PDU but the piggyback flag for that piggybacked PDU is not 0")
		}

		piggybackedPdu, _, err = DecodeV2PDU(piggybackedPduStream)

		if err != nil {
			return nil, nil, fmt.Errorf("On piggybacked PDU: %s", err)
		}
	}

	// teid := uint32(0)
	// hasTeidField := false
	// if (stream[0] & 0x08) == 0x08 {
	// 	hasTeidField = true
	// 	teid = binary.BigEndian.Uint32(stream[4:8])
	// }

	// priority := uint8(0)
	// hasPriorityField := false
	// if (stream[0] & 0x40) == 0x40 {
	// 	hasPriorityField = true
	// 	priority = (uint8(stream[11]) & 0xf0) >> 4
	// }

	return nil, nil, nil
}
