package gtpv2

import (
	"encoding/binary"
	"fmt"
)

// MessageType represents possible GTPv2 message type values
type MessageType uint8

// GTPv2 MessageTypes
const (
	EchoRequest                                MessageType = 1
	EchoResponse                               MessageType = 2
	CreateSessionRequest                       MessageType = 32
	CreateSessionResponse                      MessageType = 33
	ModifyBearerRequest                        MessageType = 34
	ModifyBearerResponse                       MessageType = 35
	DeleteSessionRequest                       MessageType = 36
	DeleteSessionResponse                      MessageType = 37
	RemoteUEReportNotification                 MessageType = 40
	RemoteUEReportAcknowlegement               MessageType = 41
	ChangeNotificationRequest                  MessageType = 38
	ChangeNotificationResponse                 MessageType = 39
	ModifyBearerCommand                        MessageType = 64
	ModifyBearerFailureIndication              MessageType = 65
	DeleteBearerCommand                        MessageType = 66
	DeleteBearerFailureIndication              MessageType = 67
	BearerResourceCommand                      MessageType = 68
	BearerResourceFailureIndication            MessageType = 69
	DownlinkDataNotificationFailureIndication  MessageType = 70
	TraceSessionActivation                     MessageType = 71
	TraceSessionDeactivation                   MessageType = 72
	StopPagingIndication                       MessageType = 73
	CreateBearerRequest                        MessageType = 95
	CreateBearerResponse                       MessageType = 96
	UpdateBearerRequest                        MessageType = 97
	UpdateBearerResponse                       MessageType = 98
	DeleteBearerRequest                        MessageType = 99
	DeleteBearerResponse                       MessageType = 100
	DeletePDNConnectionSetRequest              MessageType = 101
	DeletePDNConnectionSetResponse             MessageType = 102
	PGWDownlinkTriggeringNotification          MessageType = 103
	PGWDownlinkTriggeringAcknowledge           MessageType = 104
	IdentificationRequest                      MessageType = 128
	IdentificationResponse                     MessageType = 129
	ContextRequest                             MessageType = 130
	ContextResponse                            MessageType = 131
	ContextAcknowledge                         MessageType = 132
	ForwardRelocationRequest                   MessageType = 133
	ForwardRelocationResponse                  MessageType = 134
	ForwardRelocationCompleteNotification      MessageType = 135
	ForwardRelocationCompleteAcknowledge       MessageType = 136
	ForwardAccessContextNotification           MessageType = 137
	ForwardAccessContextAcknowledge            MessageType = 138
	RelocationCancelRequest                    MessageType = 139
	RelocationCancelResponse                   MessageType = 140
	ConfigurationTransferTunnel                MessageType = 141
	DetachNotification                         MessageType = 149
	DetachAcknowledge                          MessageType = 150
	CSPagingIndication                         MessageType = 151
	RANInformationRelay                        MessageType = 152
	AlertMMENotification                       MessageType = 153
	AlertMMEAcknowledge                        MessageType = 154
	UEActivityNotification                     MessageType = 155
	UEActivityAcknowledge                      MessageType = 156
	ISRStatusIndication                        MessageType = 157
	UERegistrationQueryRequest                 MessageType = 158
	UERegistrationQueryResponse                MessageType = 159
	CreateForwardingTunnelRequest              MessageType = 160
	CreateForwardingTunnelResponse             MessageType = 161
	SuspendNotification                        MessageType = 162
	SuspendAcknowledge                         MessageType = 163
	ResumeNotification                         MessageType = 164
	ResumeAcknowledge                          MessageType = 165
	CreateIndirectDataForwardingTunnelRequest  MessageType = 166
	CreateIndirectDataForwardingTunnelResponse MessageType = 167
	DeleteIndirectDataForwardingTunnelRequest  MessageType = 168
	DeleteIndirectDataForwardingTunnelResponse MessageType = 169
	ReleaseAccessBearersRequest                MessageType = 170
	ReleaseAccessBearersResponse               MessageType = 171
	DownlinkDataNotification                   MessageType = 176
	DownlinkDataNotificationAcknowledge        MessageType = 177
	PGWRestartNotification                     MessageType = 179
	PGWRestartNotificationAcknowledge          MessageType = 180
	UpdatePDNConnectionSetRequest              MessageType = 200
	UpdatePDNConnectionSetResponse             MessageType = 201
)

var messageNames = []string{
	"Reserved", "Echo Request", "Echo Response", "Version Not Supported Indication", "Reserved",
	"Reserved", "Reserved", "Reserved", "Reserved", "Reserved",
	"Reserved", "Reserved", "Reserved", "Reserved", "Reserved",
	"Reserved", "Reserved", "Reserved", "Reserved", "Reserved", // 19
	"Reserved", "Reserved", "Reserved", "Reserved", "Reserved",
	"Reserved", "Reserved", "Reserved", "Reserved", "Reserved",
	"Reserved", "Reserved", "Create Session Request", "Create Session Response", "Modify Bearer Request",
	"Modify Bearer Response", "Delete Session Request", "Delete Session Response", "Change Notification Request", "Change Notification Response", // 39
	"Remote UE Report Notification", "Remote UE Report Acknowledge", "Reserved", "Reserved", "Reserved",
	"Reserved", "Reserved", "Reserved", "Reserved", "Reserved",
	"Reserved", "Reserved", "Reserved", "Reserved", "Reserved",
	"Reserved", "Reserved", "Reserved", "Reserved", "Reserved", // 59
	"Reserved", "Reserved", "Reserved", "Reserved", "Modify Bearer Command",
	"Modify Bearer Failure Indication", "Delete Bearer Command", "Delete Bearer Failure Indication", "Bearer Resource Command", "Bearer Resource Failure Indication",
	"Downlink Data Notification Failure Indication", "Trace Session Activation", "Trace Session Deactivation", "Stop Paging Indication", "Reserved",
	"Reserved", "Reserved", "Reserved", "Reserved", "Reserved", // 79
	"Reserved", "Reserved", "Reserved", "Reserved", "Reserved",
	"Reserved", "Reserved", "Reserved", "Reserved", "Reserved",
	"Reserved", "Reserved", "Reserved", "Reserved", "Reserved",
	"Create Bearer Request", "Create Bearer Response", "Update Bearer Request", "Update Bearer Response", "Delete Bearer Request", // 99
	"Delete Bearer Response", "Delete PDN Connection Set Request", "Delete PDN Connection Set Response", "PGW Downlink Triggering Notification", "PGW Downlink Triggering Acknowledge",
	"Reserved", "Reserved", "Reserved", "Reserved", "Reserved",
	"Reserved", "Reserved", "Reserved", "Reserved", "Reserved",
	"Reserved", "Reserved", "Reserved", "Reserved", "Reserved", // 119
	"Reserved", "Reserved", "Reserved", "Reserved", "Reserved",
	"Reserved", "Reserved", "Reserved", "Identification Request", "Identification Response",
	"Context Request", "Context Response", "Context Acknowledge", "Forward Relocation Request", "Forward Relocation Response",
	"Forward Relocation Complete Notification", "Forward Relocation Complete Acknowledge", "Forward Access Context Notification", "Forward Access Context Acknowledge", "Relocation Cancel Request", // 139
	"Relocation Cancel Response", "Configuration Transfer Tunnel", "Reserved", "Reserved", "Reserved",
	"Reserved", "Reserved", "Reserved", "Reserved", "Detach Notification",
	"Detach Acknowledge", "CS Paging Indication", "RAN Information Relay", "Alert MME Notification", "Alert MME Acknowledge",
	"UE Activity Notification", "UE Activity Acknowledge", "ISR Status Indication", "UE Registration Query Request", "UE Registration Query Response", // 159
	"Create Forwarding Tunnel Request", "Create Forwarding Tunnel Response", "Suspend Notification", "Suspend Acknowledge", "Resume Notification",
	"Resume Acknowledge", "Create Indirect Data Forwarding Tunnel Request", "Create Indirect Data Forwarding Tunnel Response", "Delete Indirect Data Forwarding Tunnel Request", "Delete Indirect Data Forwarding Tunnel Response",
	"Release Access Bearers Request", "Release Access Bearers Response", "Reserved", "Reserved", "Reserved",
	"Reserved", "Downlink Data Notification", "Downlink Data Notification Acknowledge", "Reserved", "PGW Restart Notification", // 179
	"PGW Restart Notification Acknowledge", "Reserved", "Reserved", "Reserved", "Reserved",
	"Reserved", "Reserved", "Reserved", "Reserved", "Reserved",
	"Reserved", "Reserved", "Reserved", "Reserved", "Reserved",
	"Reserved", "Reserved", "Reserved", "Reserved", "Reserved", // 199
	"Update PDN Connection Set Request", "Update PDN Connection Set Response", "Reserved", "Reserved", "Reserved",
	"Reserved", "Reserved", "Reserved", "Reserved", "Reserved",
	"Reserved", "Modify Access Bearers Request", "Modify Access Bearers Response", "Reserved", "Reserved",
	"Reserved", "Reserved", "Reserved", "Reserved", "Reserved", // 219
	"Reserved", "Reserved", "Reserved", "Reserved", "Reserved",
	"Reserved", "Reserved", "Reserved", "Reserved", "Reserved",
	"Reserved", "MBMS Session Start Request", "MBMS Session Start Response", "MBMS Session Update Request", "MBMS Session Update Response",
	"MBMS Session Stop Request", "MBMS Session Stop Response", "Reserved", "Reserved", "Reserved", // 239
	"Reserved", "Reserved", "Reserved", "Reserved", "Reserved",
	"Reserved", "Reserved", "Reserved", "Reserved", "Reserved",
	"Reserved", "Reserved", "Reserved", "Reserved", "Reserved",
	"Reserved", // 255
}

// NameOfMessageForType returns a string identifier (from TS 29.274 section 8.1) for
// a GTPv2 IE based on the type integer value
func NameOfMessageForType(msgType MessageType) string {
	return messageNames[int(msgType)]
}

// PDU represents a GTPv2 PDU.  Version field is omitted because it is always '2'.  TEID
// should be 0 if TEIDFieldIsPresent is false.  Similarly, Priority should be 0 if
// PriorityFieldIsPresent is false.  TotalLength includes complete header length, and body
// length, but does not include a piggybacked message length if IsCarryingPiggybackedPDU
// is true.  The SequenceNumber is actually a uint24 value.  Priority is actually a uint4 value.
// For these two, upper bits beyond the actual encode size are ignored and should be zero.
type PDU struct {
	IsCarryingPiggybackedPDU bool
	TEIDFieldIsPresent       bool
	PriorityFieldIsPresent   bool
	Type                     MessageType
	TotalLength              uint16
	TEID                     uint32
	SequenceNumber           uint32
	Priority                 uint8
	InformationElements      []*IE
}

// NewPDU constructs a new base GTPv2 PDU.  It uses a builder pattern to
// add non-mandatory elements, including a TEID and a priority.  A piggybacked
// PDU is added at the time of encoding and revealed on decoding.  If you change
// struct values after construction, Encode() may not operate as expected and may
// even panic, so the struct values should usually be treated as read-only.
// This version of the constructor will panic if the length of the IEs exceeds
// the maximum PDU length.  If you want to be able to catch this condition,
// construct the PDU struct manually.
func NewPDU(pduType MessageType, sequenceNumber uint32, ies []*IE) *PDU {
	pduLength := uint32(8)

	for _, ie := range ies {
		// compute of IE length is data length + 4 bytes for IE header
		pduLength += uint32(len(ie.Data) + 4)
	}

	if pduLength > 0xffff {
		panic("Combined IE lengths exceed maximum PDU length")
	}

	return &PDU{
		IsCarryingPiggybackedPDU: false,
		TEIDFieldIsPresent:       false,
		PriorityFieldIsPresent:   false,
		Type:                     pduType,
		TEID:                     0,
		SequenceNumber:           sequenceNumber,
		Priority:                 0,
		InformationElements:      ies,
		TotalLength:              uint16(pduLength),
	}
}

// AddTEID sets the TEID field and the teid presence flag
func (pdu *PDU) AddTEID(teid uint32) *PDU {
	pdu.TEIDFieldIsPresent = true
	pdu.TEID = teid
	pdu.TotalLength += 4

	return pdu
}

// AddPriority sets the priority field and the priority presence flag
func (pdu *PDU) AddPriority(priority uint8) *PDU {
	pdu.PriorityFieldIsPresent = true
	pdu.Priority = priority & 0x0f
	return pdu
}

// Encode encodes the GTPv2 PDU as a byte stream in network byte order,
// suitable for trasmission.
func (pdu *PDU) Encode() []byte {
	encoded := make([]byte, pdu.TotalLength)

	encoded[0] = 0x40
	encoded[1] = uint8(pdu.Type)
	binary.BigEndian.PutUint16(encoded[2:4], pdu.TotalLength-4)

	ieOffsetByteIndex := 0

	if pdu.TEIDFieldIsPresent {
		encoded[0] |= 0x08
		binary.BigEndian.PutUint32(encoded[4:8], pdu.TEID)
		binary.BigEndian.PutUint32(encoded[8:12], pdu.SequenceNumber<<8)

		if pdu.PriorityFieldIsPresent {
			encoded[0] |= 0x04
			encoded[11] = pdu.Priority << 4
		}
		ieOffsetByteIndex = 12
	} else {
		binary.BigEndian.PutUint32(encoded[4:8], pdu.SequenceNumber<<8)
		ieOffsetByteIndex = 8
	}

	for _, ie := range pdu.InformationElements {
		encodedIE := ie.Encode()
		offsetForEndOfIE := ieOffsetByteIndex + len(encodedIE)

		copy(encoded[ieOffsetByteIndex:offsetForEndOfIE], encodedIE)

		ieOffsetByteIndex = offsetForEndOfIE
	}

	return encoded
}

// DecodePDU decodes a stream of bytes that contain either exactly one well-formed
// GTPv2 PDU, or two GTPv2 PDUs when the piggyback flag on the first is set to true.
// Returns an error if the stream cannot be decoded into one or two PDUs.
func DecodePDU(stream []byte) (pdu *PDU, piggybackedPdu *PDU, err error) {
	piggybackedPdu = nil

	if len(stream) < 8 {
		return nil, nil, fmt.Errorf("stream length (%d) too short for a GTPv2 PDU", len(stream))
	}

	if (stream[0] >> 5) != 2 {
		return nil, nil, fmt.Errorf("GTPv2 PDU version should be 2, but in stream, it is (%d)", (stream[0] >> 5))
	}

	hasPiggybackedPdu := (stream[0] & 0x10) == 0x10

	msgLengthFieldValue := binary.BigEndian.Uint16(stream[2:4])
	totalPduLength := msgLengthFieldValue + 4

	if len(stream) < int(totalPduLength) {
		return nil, nil, fmt.Errorf("GTPv2 PDU length field is (%d), so total length should be (%d), but stream length is (%d)", msgLengthFieldValue, totalPduLength, len(stream))
	}

	if !hasPiggybackedPdu {
		if len(stream) != int(totalPduLength) {
			return nil, nil, fmt.Errorf("GTPv2 PDU length field is (%d), so total length should be (%d), but stream length is (%d)", msgLengthFieldValue, totalPduLength, len(stream))
		}
	} else {
		piggybackedPduStream := stream[totalPduLength:]

		if (piggybackedPduStream[0] & 0x10) != 0 {
			return nil, nil, fmt.Errorf("GTPv2 PDU has piggybacked PDU but the piggyback flag for that piggybacked PDU is not 0")
		}

		piggybackedPdu, _, err = DecodePDU(piggybackedPduStream)

		if err != nil {
			return nil, nil, fmt.Errorf("on piggybacked PDU: %s", err)
		}

		if len(stream) != int(totalPduLength)+int(piggybackedPdu.TotalLength) {
			return nil, nil, fmt.Errorf("stream contains more than single PDU and piggybacked PDU")
		}
	}

	teid := uint32(0)
	sequenceNumber := uint32(0)
	var headerLength int
	hasTeidField := false

	if (stream[0] & 0x08) == 0x08 {
		hasTeidField = true
		teid = binary.BigEndian.Uint32(stream[4:8])
		sequenceNumber = binary.BigEndian.Uint32(stream[8:12]) >> 8
		headerLength = 12
	} else {
		sequenceNumber = binary.BigEndian.Uint32(stream[4:8]) >> 8
		headerLength = 8
	}

	hasPriorityField := (stream[0] & 0x04) == 0x04

	priority := uint8(0)
	if hasPriorityField {
		priority = (uint8(stream[11]) & 0xf0) >> 4
	}

	pdu = &PDU{
		IsCarryingPiggybackedPDU: hasPiggybackedPdu,
		TEIDFieldIsPresent:       hasTeidField,
		PriorityFieldIsPresent:   hasPriorityField,
		TEID:                     teid,
		Priority:                 priority,
		SequenceNumber:           sequenceNumber,
		TotalLength:              totalPduLength,
		Type:                     MessageType(stream[1]),
	}

	ieSet := make([]*IE, 0, 10)

	for i := headerLength; i < int(totalPduLength); {
		nextIEInStream, err := DecodeIE(stream[i:])

		if err != nil {
			return nil, nil, err
		}

		ieSet = append(ieSet, nextIEInStream)

		i += int(nextIEInStream.TotalLength)
	}

	pdu.InformationElements = ieSet

	return pdu, piggybackedPdu, nil
}
