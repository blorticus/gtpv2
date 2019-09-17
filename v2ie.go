package gtp

import (
	"encoding/binary"
	"fmt"
)

// V2IEType represents the various IE types for GTPv2
type V2IEType uint8

// These represent possible GTPv2 PDU types.  In some cases, includes the
// full name and its abbreviation (e.g., for IMSI)
const (
	InternationalMobileSubscriberIdentity                = 1
	IMSI                                                 = 1
	Cause                                                = 2
	RecoveryRestartCounter                               = 3
	STNSR                                                = 51
	AccessPointName                                      = 71
	APN                                                  = 71
	AggregateMaximumBitRate                              = 72
	AMBR                                                 = 72
	EPSBearerID                                          = 73
	EBI                                                  = 73
	IPAddress                                            = 74
	MobileEquipmentIdentity                              = 75
	MEI                                                  = 75
	MSISDN                                               = 76
	Indication                                           = 77
	ProtocolConfigurationOptions                         = 78
	PCI                                                  = 78
	PDNAddressAllocation                                 = 79
	PAA                                                  = 79
	BearerLevelQualityofService                          = 80
	BearerQoS                                            = 80
	FlowQualityofService                                 = 81
	FlowQoS                                              = 81
	RATType                                              = 82
	ServingNetwork                                       = 83
	EPSBearerLevelTrafficFlowTemplate                    = 84
	BearerTFT                                            = 84
	TrafficAggregationDescription                        = 85
	TAD                                                  = 85
	UserLocationInformation                              = 86
	ULI                                                  = 86
	FullyQualifiedTunnelEndpointIdentifier               = 87
	FTEID                                                = 87
	TMSI                                                 = 88
	GlobalCNId                                           = 89
	S103PDNDataForwardingInfo                            = 90
	S103PDF                                              = 90
	S1UDataForwardingInfo                                = 91
	S1UDF                                                = 91
	DelayValue                                           = 92
	BearerContext                                        = 93
	ChargingID                                           = 94
	ChargingCharacteristics                              = 95
	TraceInformation                                     = 96
	BearerFlags                                          = 97
	PDNType                                              = 99
	ProcedureTransactionID                               = 100
	MMContextGSMKeyandTriplets                           = 103
	MMContextUMTSKeyUsedCipherandQuintuplets             = 104
	MMContextGSMKeyUsedCipherandQuintuplets              = 105
	MMContextUMTSKeyandQuintuplets                       = 106
	MMContextEPSSecurityContextQuadrupletsandQuintuplets = 107
	MMContextUMTSKeyQuadrupletsandQuintuplets            = 108
	PDNConnection                                        = 109
	PDUNumbers                                           = 110
	PTMSI                                                = 111
	PTMSISignature                                       = 112
	HopCounter                                           = 113
	UETimeZone                                           = 114
	TraceReference                                       = 115
	CompleteRequestMessage                               = 116
	GUTI                                                 = 117
	FContainer                                           = 118
	FCause                                               = 119
	PLMNID                                               = 120
	TargetIdentification                                 = 121
	PacketFlowID                                         = 123
	RABContext                                           = 124
	SourceRNCPDCPContextInfo                             = 125
	PortNumber                                           = 126
	APNRestriction                                       = 127
	SelectionMode                                        = 128
	SourceIdentification                                 = 129
	ChangeReportingAction                                = 131
	FullyQualifiedPDNConnectionSetIdentifier             = 132
	FQCSID                                               = 132
	Channelneeded                                        = 133
	eMLPPPriority                                        = 134
	NodeType                                             = 135
	FullyQualifiedDomainName                             = 136
	FQDN                                                 = 136
	TransactionIdentifier                                = 137
	TI                                                   = 137
	MBMSSessionDuration                                  = 138
	MBMSServiceArea                                      = 139
	MBMSSessionIdentifier                                = 140
	MBMSFlowIdentifier                                   = 141
	MBMSIPMulticastDistribution                          = 142
	MBMSDistributionAcknowledge                          = 143
	RFSPIndex                                            = 144
	UserCSGInformation                                   = 145
	UCI                                                  = 145
	CSGInformationReportingAction                        = 146
	CSGID                                                = 147
	CSGMembershipIndication                              = 148
	CMI                                                  = 148
	Serviceindicator                                     = 149
	DetachType                                           = 150
	LocalDistiguishedName                                = 151
	LDN                                                  = 151
	NodeFeatures                                         = 152
	MBMSTimetoDataTransfer                               = 153
	Throttling                                           = 154
	AllocationRetentionPriority                          = 155
	ARP                                                  = 155
	EPCTimer                                             = 156
	SignallingPriorityIndication                         = 157
	TemporaryMobileGroupIdentity                         = 158
	TMGI                                                 = 158
	AdditionalMMcontextforSRVCC                          = 159
	AdditionalflagsforSRVCC                              = 160
	MDTConfiguration                                     = 162
	AdditionalProtocolConfigurationOptions               = 163
	APCO                                                 = 163
	AbsoluteTimeofMBMSDataTransfer                       = 164
	HeNBInformationReporting                             = 165
	IPv4ConfigurationParameters                          = 166
	IP4CP                                                = 166
	ChangetoReportFlags                                  = 167
	ActionIndication                                     = 168
	TWANIdentifier                                       = 169
	ULITimestamp                                         = 170
	MBMSFlags                                            = 171
	RANNASCause                                          = 172
	CNOperatorSelectionEntity                            = 173
	TrustedWLANModeIndication                            = 174
	NodeNumber                                           = 175
	NodeIdentifier                                       = 176
	PresenceReportingAreaAction                          = 177
	PresenceReportingAreaInformation                     = 178
	TWANIdentifierTimestamp                              = 179
	OverloadControlInformation                           = 180
	LoadControlInformation                               = 181
	Metric                                               = 182
	SequenceNumber                                       = 183
	APNandRelativeCapacity                               = 184
	WLANOffloadabilityIndication                         = 185
	PagingandServiceInformation                          = 186
	IntegerNumber                                        = 187
	MillisecondTimeStamp                                 = 188
	MonitoringEventInformation                           = 189
	ECGIList                                             = 190
	RemoteUEContext                                      = 191
	RemoteUserID                                         = 192
	RemoteUEIPinformation                                = 193
	CIoTOptimizationsSupportIndication                   = 194
	SCEFPDNConnection                                    = 195
	HeaderCompressionConfiguration                       = 196
	ExtendedProtocolConfigurationOptions                 = 197
	ePCO                                                 = 197
	ServingPLMNRateControl                               = 198
	Counter                                              = 199
	MappedUEUsageType                                    = 200
	SecondaryRATUsageDataReport                          = 201
	UPFunctionSelectionIndicationFlags                   = 202
	ExtensionType                                        = 254
	PrivateExtension                                     = 255
)

var v2IENames = []string{
	"Reserved", "International Mobile Subscriber Identity (IMSI)", "Cause",
	"Recovery (Restart Counter)", "Reserved", "Reserved",
	"Reserved", "Reserved", "Reserved", "Reserved", "Reserved",
	"Reserved", "Reserved", "Reserved", "Reserved", "Reserved",
	"Reserved", "Reserved", "Reserved", "Reserved", "Reserved",
	"Reserved", "Reserved", "Reserved", "Reserved", "Reserved",
	"Reserved", "Reserved", "Reserved", "Reserved", "Reserved",
	"Reserved", "Reserved", "Reserved", "Reserved", "Reserved",
	"Reserved", "Reserved", "Reserved", "Reserved", "Reserved",
	"Reserved", "Reserved", "Reserved", "Reserved", "Reserved",
	"Reserved", "Reserved", "Reserved", "Reserved", "Reserved",
	"STN-SR", "Reserved", "Reserved", "Reserved", "Reserved",
	"Reserved", "Reserved", "Reserved", "Reserved", "Reserved",
	"Reserved", "Reserved", "Reserved", "Reserved", "Reserved",
	"Reserved", "Reserved", "Reserved", "Reserved", "Reserved",
	"Access Point Name (APN)", "Aggregate Maximum Bit Rate (AMBR)", "EPS Bearer ID (EBI)",
	"IP Address", "Mobile Equipment Identity (MEI)", "MSISDN", "Indication",
	"Protocol Configuration Options (PCO)", "PDN Address Allocation (PAA)",
	"Bearer Level Quality of Service (Bearer QoS)", "Flow Quality of Service (Flow QoS)",
	"RAT Type", "Serving Network", "EPS Bearer Level Traffic Flow Template (Bearer TFT)",
	"Traffic Aggregation Description (TAD)", "User Location Information (ULI)",
	"Fully Qualified Tunnel Endpoint Identifier (F-TEID)", "TMSI", "Global CN-Id",
	"S103 PDN Data Forwarding Info (S103PDF)", "S1-U Data Forwarding Info (S1UDF)", "Delay Value",
	"Bearer Context", "Charging ID", "Charging Characteristics", "Trace Information",
	"Bearer Flags", "Reserved", "PDN Type", "Procedure Transaction ID",
	"Reserved", "Reserved",
	"MM Context (GSM Key and Triplets)", "MM Context (UMTS Key, Used Cipher and Quintuplets)",
	"MM Context (GSM Key, Used Cipher and Quintuplets)", "MM Context (UMTS Key and Quintuplets)",
	"MM Context (EPS Security Context, Quadruplets and Quintuplets)",
	"MM Context (UMTS Key, Quadruplets and Quintuplets)",
	"PDN Connection", "PDU Numbers", "P-TMSI", "P-TMSI Signature", "Hop Counter",
	"UE Time Zone", "Trace Reference", "Complete Request Message", "GUTI",
	"F-Container", "F-Cause", "PLMN ID", "Target Identification", "Reserved",
	"Packet Flow ID", "RAB Context", "Source RNC PDCP Context Info", "Port Number",
	"APN Restriction", "Selection Mode", "Source Identification",
	"Reserved", "Change Reporting Action", "Fully Qualified PDN Connection Set Identifier (FQ-CSID)",
	"Channel needed", "eMLPP Priority", "Node Type", "Fully Qualified Domain Name (FQDN)",
	"Transaction Identifier (TI)", "MBMS Session Duration", "MBMS Service Area",
	"MBMS Session Identifier", "MBMS Flow Identifier", "MBMS IP Multicast Distribution",
	"MBMS Distribution Acknowledge", "RFSP Index", "User CSG Information (UCI)",
	"CSG Information Reporting Action", "CSG ID", "CSG Membership Indication (CMI)",
	"Service indicator", "Detach Type", "Local Distiguished Name (LDN)",
	"Node Features", "MBMS Time to Data Transfer", "Throttling", "Allocation/Retention Priority (ARP)",
	"EPC Timer", "Signalling Priority Indication", "Temporary Mobile Group Identity (TMGI)",
	"Additional MM context for SRVCC", "Additional flags for SRVCC", "Reserved",
	"MDT Configuration", "Additional Protocol Configuration Options (APCO)",
	"Absolute Time of MBMS Data Transfer", "H(e)NB Information Reporting ",
	"IPv4 Configuration Parameters (IP4CP)", "Change to Report Flags",
	"Action Indication", "TWAN Identifier", "ULI Timestamp", "MBMS Flags", "RAN/NAS Cause",
	"CN Operator Selection Entity", "Trusted WLAN Mode Indication", "Node Number",
	"Node Identifier", "Presence Reporting Area Action", "Presence Reporting Area Information",
	"TWAN Identifier Timestamp", "Overload Control Information", "Load Control Information",
	"Metric", "Sequence Number", "APN and Relative Capacity", "WLAN Offloadability Indication",
	"Paging and Service Information", "Integer Number", "Millisecond Time Stamp",
	"Monitoring Event Information", "ECGI List", "Remote UE Context", "Remote User ID",
	"Remote UE IP information", "CIoT Optimizations Support Indication", "SCEF PDN Connection",
	"Header Compression Configuration", "Extended Protocol Configuration Options (ePCO)",
	"Serving PLMN Rate Control", "Counter", "Mapped UE Usage Type", "Secondary RAT Usage Data Report",
	"UP Function Selection Indication Flags", "Reserved", "Reserved", "Reserved",
	"Reserved", "Reserved", "Reserved", "Reserved", "Reserved", "Reserved",
	"Reserved", "Reserved", "Reserved", "Reserved", "Reserved", "Reserved",
	"Reserved", "Reserved", "Reserved", "Reserved", "Reserved", "Reserved",
	"Reserved", "Reserved", "Reserved", "Reserved", "Reserved", "Reserved",
	"Reserved", "Reserved", "Reserved", "Reserved", "Reserved", "Reserved",
	"Reserved", "Reserved", "Reserved", "Reserved", "Reserved", "Reserved",
	"Reserved", "Reserved", "Reserved", "Reserved", "Reserved", "Reserved",
	"Reserved", "Reserved", "Reserved", "Reserved", "Reserved", "Reserved",
	"IE Extension", "Private Extension",
}

// NameOfV2IEForType returns a string identifier (from TS 29.274 section 8.1) for
// a GTPv2 IE based on the type integer
func NameOfV2IEForType(ieType V2IEType) string {
	return v2IENames[int(ieType)]
}

// V2IE is a GTPv2 Information Element.  DataLength is the length of just
// the contained data, in bytes.  TotalLength is the DataLength plus the
// header length.  InstanceNumber is actually uint4.  Data is the BigEndian
// data bytes.
type V2IE struct {
	Type           V2IEType
	DataLength     uint16
	TotalLength    uint16
	InstanceNumber uint8
	Data           []byte
}

// DecodeV2IE consumes bytes from the start of stream to produce a V2IE.
// The TotalLength field of the resulting V2IE provides the count of bytes
// from stream that are consumed to produce this IE.  Return an error if
// decoding fails.
func DecodeV2IE(stream []byte) (*V2IE, error) {
	if len(stream) < 4 {
		return nil, fmt.Errorf("Insufficient octets in stream for a complete GTPv2 IE")
	}

	ie := &V2IE{
		Type:           V2IEType(stream[0]),
		DataLength:     binary.BigEndian.Uint16(stream[1:3]),
		InstanceNumber: uint8(stream[3]) & 0x0f,
	}

	ie.TotalLength = ie.DataLength + 4

	if len(stream) < int(ie.TotalLength) {
		return nil, fmt.Errorf("Next IE length field is (%d), which requires (%d) bytes in stream, but there are only (%d) bytes", ie.DataLength, ie.TotalLength, len(stream))
	}

	ie.Data = make([]byte, ie.DataLength)
	copy(ie.Data, stream[4:ie.DataLength+4])

	return ie, nil
}
