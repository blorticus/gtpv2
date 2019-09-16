// Package gtp provides encoding and decoding of GTP packets
package gtp

import (
    "fmt"
    "encoding/binary"
    "bytes"
    "utils"
)

const (
    NONE = 0x0
    GTP_VERSION_MASK = 0xE0
    GTPV1_NPUD_MASK = 0xFF00
    GTPV1_EXTENSION_HEADER_MASK = 0x00FF
    GTPV1_VERSION_UINT32_MASK = 0xFF00
    GTPV2_VERSION_UINT32_MASK = 0xFFF0
    GTPV2_F_PIGGYBACK = 0x10
    GTPV2_F_TEID = 0x08
    GTPV2_TEID_MASK = 0xFFFFFF00
    GTPV1_F_PROTOCOL_TYPE = 0x20
    GTPV1_F_EXT_HEADER = 0x04
    GTPV1_F_SEQ = 0x02
    GTPV1_F_NPUD = 0x01
    GTPV1_TUNNEL = 0xff
    GTPV1_MANDATORY_LENGTH = 0              // TEID is after heading
    GTPV2_MANDATORY_LENGTH = 4              // SEQ + SPARE
    GTPV2_IE_MANDATORY_FIELD_LENGTH = 4
    GTPV2_IE_SPARE = uint8(0)
    GTPV1_IE_FORMAT = 0x80
)

type GTPVersion uint8

const (
    GTPv0 GTPVersion = 0
    GTPv1 GTPVersion = 1
    GTPv2 GTPVersion = 2
)

type GTPPacket struct {
    version GTPVersion
    flags uint8
    msg_type uint8
    msg_length uint16
    teid uint32
    raw []byte
    data []byte
    ies []*GTPPacketIEData
}

type GTPv1Packet struct {
    GTPPacket
    seq uint16
    n_pdu uint8
    ext_header uint8
}

type GTPv2Packet struct {
    GTPPacket
    seq uint32
}

type GTPPacketInterface interface {
    Sequence() uint32
}

func (self *GTPPacket) GetType() uint8 {
    return self.msg_type
}

func packetSequenceFromInterface(packet interface{}) uint32 {
    sequence := uint32(0)
    switch p := packet.(type) {
    case *GTPv1Packet:
        sequence = uint32(p.seq)
    case *GTPv2Packet:
        sequence = p.seq
    }
    return sequence
}

func GTPMessageFromInterface(packet interface{}) *GTPMessage {
    switch p := packet.(type) {
    case *GTPv1Packet:
        return TypeToGTPMessage[GTPv1][p.msg_type]
    case *GTPv2Packet:
        return TypeToGTPMessage[GTPv2][p.msg_type]
    }
    return nil
}

func NewGTP(raw_data []byte) interface{} {
    var flags uint8

    buf := bytes.NewReader(raw_data)
    err := binary.Read(buf, binary.BigEndian, &flags)
    if err != nil {
        panic(err)
    }
    version := GTPVersion((flags & GTP_VERSION_MASK) >> 5)
    switch (version) {
    case GTPv1, GTPv0:
        return NewGTPv1PacketFromBuffer(flags, buf)
    case GTPv2:
        return NewGTPv2PacketFromBuffer(flags, buf)
    }
    panic(fmt.Errorf("unknown version: ", version))
}

func getIEFromId(version GTPVersion, ie_id uint8) *GTPIE {
    switch(version) {
    case GTPv2:
        ie, ok := IDToGTPv2IE[ie_id]
        if !ok {
            return &GTPIE{version: GTPv2, format: TLV}
        }
        return ie
    case GTPv1:
        ie, ok := IDToGTPv1IE[ie_id]
        if !ok {
            return &GTPIE{version: GTPv1, format: TLV}
        }
        return ie
    }
    return nil
}

type GTPMessage struct {
    id uint8
    version GTPVersion
    name string
    description string
    response *GTPMessage
    mandatory []*GTPIE
}


var TypeToGTPMessage = map[GTPVersion]map[uint8]*GTPMessage{
    GTPv1: map[uint8]*GTPMessage{},
    GTPv2: map[uint8]*GTPMessage{}}
var NameToGTPMessage = map[string]*GTPMessage {}
var IDToGTPv1IE = map[uint8]*GTPIE {}
var IDToGTPv2IE = map[uint8]*GTPIE {}
var NameToGTPIE = map[string]*GTPIE {}

type GTPIEType uint8
type GTPFormat uint8

const (
    Unsigned8 GTPIEType = 1 + iota
    Unsigned32
    Unsigned64
    Enumerated
    UTF8String
    OctetString
    Time
    Address
)

const (
    TV GTPFormat = 0
    TLV GTPFormat = 1
)

type GTPIE struct {
    id uint8
    version GTPVersion
    ie_type GTPIEType
    size uint16
    format GTPFormat
    name string
    description string
}

func (self *GTPMessage) ToPacket(seq uint32, teid uint32,
        ies []*GTPPacketIEData) *GTPv2Packet {
    return NewGTPv2Packet(self, seq, teid, ies)
}

func (self *GTPMessage) GetType() uint8 {
    return self.id
}

func (self *GTPMessage) GetName() string {
    return self.name
}

func ConstructGTPMessage(msg_string string) []byte {
    msg, ok := NameToGTPMessage[msg_string]
    if !ok {
        panic(fmt.Errorf("message not found [%s]", msg_string))
    }
    ies := make([]*GTPPacketIEData, 0)
    packet := NewGTPv2Packet(msg, 0, 0, ies)
    return packet.Encode()
}

type GTPPacketIEData struct {
    id uint8
    length uint16
    data []byte
    format GTPFormat
    version GTPVersion
}

func NewGTPPacketIEData(id uint8, length uint16, payload []byte, ie *GTPIE) *GTPPacketIEData {
    d := &GTPPacketIEData{id: id, length: length, data: payload, format: TLV}
    if ie != nil {
        d.format = ie.format
        d.version = ie.version
    }
    return d
}

// for GTPv1 - when the most significant bit of (8-bit) type field is 1, it's TLV,
// TV other wise.  For GTPv2 - it's always TLV encoding.
func (self *GTPPacketIEData) Encode() []byte {
    buf := new(bytes.Buffer)

    utils.CheckError(binary.Write(buf, binary.LittleEndian, self.id))
    switch self.version {
    case GTPv2:
        utils.CheckError(binary.Write(buf, binary.BigEndian, self.length))
        utils.CheckError(binary.Write(buf, binary.LittleEndian, GTPV2_IE_SPARE))
    default:
        if self.format == TLV {
            utils.CheckError(binary.Write(buf, binary.BigEndian, self.length))
        }
    }
    utils.CheckError(binary.Write(buf, binary.LittleEndian, self.data))

    return buf.Bytes()
}
