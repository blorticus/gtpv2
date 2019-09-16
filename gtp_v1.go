package gtp

import (
    "fmt"
    "encoding/binary"
    "bytes"
    "io"
    "utils"
)

func NewGTPv1PacketFromBuffer(flags uint8, buf io.Reader) *GTPv1Packet {
    self := &GTPv1Packet{
        GTPPacket: GTPPacket{version: GTPv1, flags: flags}, seq: 0, n_pdu: 0, ext_header: 0}
    self.read_buffer(buf)
    return self
}
// the sequence int is so we can destinguish between prescent (>=0), vs not
// prescent (-1)
func NewGTPv1Packet(prime bool, sequence int, n_pdu uint8, ext_header uint8,
        msg_type uint8, teid uint32, ies []*GTPPacketIEData, tunnel_data []byte) *GTPv1Packet {
    self := &GTPv1Packet{
        GTPPacket: GTPPacket{version: GTPv1, flags: 0x30,
                             msg_type: msg_type, ies: ies} }
    self.msg_length = GTPV1_MANDATORY_LENGTH
    if !prime {
        self.flags |= GTPV1_F_PROTOCOL_TYPE
    }
    if sequence != -1 {
        self.flags |= GTPV1_F_SEQ
        self.seq = uint16(sequence)
    }
    if n_pdu != 0 {
        self.flags |= GTPV1_F_NPUD
        self.n_pdu = n_pdu
    }
    if ext_header != 0 {
        self.flags |= GTPV1_F_EXT_HEADER
        self.ext_header = ext_header
    }
    if (self.hasExtHeader() || self.hasSeq() || self.hasNPDU()) {
        self.msg_length += 4
    }

    self.teid = teid
    self.data = make([]byte, 0)
    if self.msg_type == GTPV1_TUNNEL || len(tunnel_data) > 0 {
        self.data = append(self.data, tunnel_data...)
    } else {
        for _, ie := range(self.ies) {
            ie_bytes := ie.Encode()
            self.data = append(self.data, ie_bytes...)
        }
    }
    self.msg_length += uint16(len(self.data))

    return self
}

func (self *GTPv1Packet) read_buffer(buf io.Reader) {
    utils.CheckError(binary.Read(buf, binary.BigEndian, &self.msg_type))
    utils.CheckError(binary.Read(buf, binary.BigEndian, &self.msg_length))
    data_length := self.msg_length
    utils.CheckError(binary.Read(buf, binary.BigEndian, &self.teid))
    if (self.hasExtHeader() || self.hasSeq() || self.hasNPDU()) {
        data_length -= 4
        var data uint32
        utils.CheckError(binary.Read(buf, binary.BigEndian, &data))
        if self.hasSeq() {
            self.seq = uint16(data >> 16);
        }
        if self.hasNPDU() {
            self.n_pdu = uint8(data & GTPV1_NPUD_MASK >> 8)
        }
        if self.hasExtHeader() {
            self.ext_header = uint8(data & GTPV1_EXTENSION_HEADER_MASK)
        }
    }
    self.data = make([]byte, data_length)
    if _, err := io.ReadFull(buf, self.data); err != nil {
        panic(err)
	}
    if self.msg_type != GTPV1_TUNNEL && len(self.data) > 0 {
        self.read_buffer_for_ie()
    }
}

func (self *GTPv1Packet) hasSeq() bool {
    return (self.flags & GTPV1_F_SEQ) != NONE
}

func (self *GTPv1Packet) hasExtHeader() bool {
    return (self.flags & GTPV1_F_EXT_HEADER) != NONE
}

func (self *GTPv1Packet) hasNPDU() bool {
    return (self.flags & GTPV1_F_EXT_HEADER) != NONE
}

func (self *GTPv1Packet) DataLength() uint16 {
    length := self.msg_length
    if (self.hasExtHeader() || self.hasSeq() || self.hasNPDU()) {
        length -= 4
    }
    return length
}

func (self *GTPv1Packet) read_buffer_for_ie() {
    length := self.DataLength()
    if length <= 0 {
        panic(fmt.Errorf("unexpected data length: %d", length))
    }
    ie_buf := bytes.NewReader(self.data)
    for length > 0 {
        var ie_id uint8
        utils.CheckError(binary.Read(ie_buf, binary.BigEndian, &ie_id))
        length -= 1
        if (ie_id & GTPV1_IE_FORMAT) == 0x0 {
            ie, ok := IDToGTPv1IE[ie_id]
            if !ok {
                fmt.Errorf("unknow GTPv1 IE[%d]", ie_id)
            }
            ie_size := ie.size
            length -= (ie_size)
            payload := make([]byte, ie_size)
            if _, err := io.ReadFull(ie_buf, payload); err != nil {
                panic(err)
            }
            self.ies = append(self.ies, NewGTPPacketIEData(ie_id, ie_size, payload, ie))
        } else {
            // TLV
            var ie_size uint16
            utils.CheckError(binary.Read(ie_buf, binary.BigEndian, &ie_size))
            length -= (ie_size + 2)
            payload := make([]byte, ie_size)
            if _, err := io.ReadFull(ie_buf, payload); err != nil {
                panic(err)
            }
            self.ies = append(self.ies, NewGTPPacketIEData(ie_id, ie_size, payload, IDToGTPv1IE[ie_id]))
        }
    }
}

func (self *GTPv1Packet) Encode() []byte {
    buf := new(bytes.Buffer)
    utils.CheckError(binary.Write(buf, binary.LittleEndian, self.flags))
    utils.CheckError(binary.Write(buf, binary.LittleEndian, self.msg_type))
    utils.CheckError(binary.Write(buf, binary.BigEndian, self.msg_length))
    utils.CheckError(binary.Write(buf, binary.BigEndian, self.teid))
    if (self.hasExtHeader() || self.hasSeq() || self.hasNPDU()) {
        data := uint32(0)
        if self.hasSeq() {
            data |= uint32(self.seq) << 16
        }
        if self.hasNPDU() {
            data |= uint32(self.n_pdu) << 8
        }
        if self.hasExtHeader() {
            data |= uint32(self.ext_header)
        }
        utils.CheckError(binary.Write(buf, binary.BigEndian, data))
    }
    utils.CheckError(binary.Write(buf, binary.LittleEndian, self.data))

    return buf.Bytes()
}

func NewGTPv1Message(id uint8, name string, description string,
        mandatory []*GTPIE, response *GTPMessage) *GTPMessage {
    d := &GTPMessage{version: GTPv1, id: id, name: name,
            description: description, mandatory: mandatory,
            response: response}
    TypeToGTPMessage[GTPv1][d.id] = d
    NameToGTPMessage[d.name] = d
    return d
}

func NewGTPv1IE(name string, id uint8, size uint16,
        description string) *GTPIE {

    d := &GTPIE{version: GTPv1, name: name, id: id, ie_type: OctetString,
            size: size, description: description}
    if size != 0 {
        d.format = TV
    } else {
        d.format = TLV
    }
    NameToGTPIE[d.name] = d
    IDToGTPv1IE[d.id] = d
    return d
}
