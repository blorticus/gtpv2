package gtp

import (
    "fmt"
    "encoding/binary"
    "bytes"
    "io"
    "utils"
)

func NewGTPv2Packet(msg *GTPMessage,
        sequence uint32, teid uint32, ies []*GTPPacketIEData) *GTPv2Packet {

    self := &GTPv2Packet{
        GTPPacket: GTPPacket{version: GTPv2, flags: 0x40,
                             msg_type: msg.id} }
    self.msg_length = GTPV2_MANDATORY_LENGTH
    if teid != 0 {
        self.msg_length += 4
        self.flags |= GTPV2_F_TEID
        self.teid = teid
    }
    self.seq = sequence
    self.ies = ies
    self.data = make([]byte, 0)
    for _, ie := range(self.ies) {
        self.data = append(self.data, ie.Encode()...)
    }
    self.msg_length += uint16(len(self.data))

    return self
}


func NewGTPv2PacketFromBuffer(flags uint8, buf io.Reader) *GTPv2Packet {
    self := &GTPv2Packet{
        GTPPacket: GTPPacket{version: GTPv2, flags: flags, teid: 0},
        seq: 0}
    self.read_buffer(buf)
    return self
}

func (self *GTPv2Packet) hasTEID() bool {
    return (self.flags & GTPV2_F_TEID) != NONE
}

func (self *GTPv2Packet) hasPiggyBack() bool {
    return (self.flags & GTPV2_F_PIGGYBACK) != NONE
}

func (self *GTPv2Packet) read_buffer(buf io.Reader) {
    utils.CheckError(binary.Read(buf, binary.BigEndian, &self.msg_type))
    utils.CheckError(binary.Read(buf, binary.BigEndian, &self.msg_length))
    // payload length is msg-length minus teid size
    if self.hasTEID() {
        utils.CheckError(binary.Read(buf, binary.BigEndian, &self.teid))
    }
    utils.CheckError(binary.Read(buf, binary.BigEndian, &self.seq))
    self.seq &= GTPV2_TEID_MASK
    self.seq = self.seq >> 8

    self.data = make([]byte, self.DataLength())
    if _, err := io.ReadFull(buf, self.data); err != nil {
        panic(err)
    }

    if len(self.data) > 0 {
        self.read_buffer_for_ie()
    }
}

func (self *GTPv2Packet) DataLength() uint16 {
    length := self.msg_length - GTPV2_IE_MANDATORY_FIELD_LENGTH
    if self.hasTEID() {
        length -= 4
    }
    return length
}

func (self *GTPv2Packet) read_buffer_for_ie() {

    length := self.DataLength()
    if length <= 0 {
        panic(fmt.Errorf("unexpected data length: %d", length))
    }
    ie_buf := bytes.NewReader(self.data)
    for length > 0 {
        var ie_id, ie_spare uint8
        var ie_length uint16
        // IE length is mandafory fields (4 bytes) + data
        utils.CheckError(binary.Read(ie_buf, binary.BigEndian, &ie_id))
        utils.CheckError(binary.Read(ie_buf, binary.BigEndian, &ie_length))
        utils.CheckError(binary.Read(ie_buf, binary.BigEndian, &ie_spare))

        length -= GTPV2_IE_MANDATORY_FIELD_LENGTH + ie_length
        payload := make([]byte, ie_length)
        if _, err := io.ReadFull(ie_buf, payload); err != nil {
            panic(err)
        }

        self.ies = append(self.ies, NewGTPPacketIEData(ie_id, ie_length, payload,
                getIEFromId(GTPv2, ie_id)))
    }
}

func (self *GTPv2Packet) Encode() []byte {
    buf := new(bytes.Buffer)
    utils.CheckError(binary.Write(buf, binary.LittleEndian, self.flags))
    utils.CheckError(binary.Write(buf, binary.LittleEndian, self.msg_type))
    utils.CheckError(binary.Write(buf, binary.BigEndian, self.msg_length))
    if self.hasTEID() {
        utils.CheckError(binary.Write(buf, binary.BigEndian, self.teid))
    }
    data := self.seq << 8
    utils.CheckError(binary.Write(buf, binary.BigEndian, data))
    utils.CheckError(binary.Write(buf, binary.LittleEndian, self.data))

    return buf.Bytes()
}

func NewGTPv2Message(id uint8, name string, description string,
        mandatory []*GTPIE, response *GTPMessage) *GTPMessage {
    d := &GTPMessage{version: GTPv2, id: id, name: name,
            description: description, mandatory: mandatory,
            response: response}
    TypeToGTPMessage[GTPv2][d.id] = d
    NameToGTPMessage[d.name] = d
    return d
}

func NewGTPv2IE(name string, id uint8, size uint16,
        description string) *GTPIE {

    d := &GTPIE{version: GTPv2, name: name, id: id, ie_type: OctetString,
            size: size, description: description, format: TLV}
    NameToGTPIE[d.name] = d
    IDToGTPv2IE[d.id] = d
    return d
}
