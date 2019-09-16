package gtp

import (
    "net"
)
/*
A struct for keeping communication internal to the node
*/

type InstanceMessageCode uint8
type InstanceStatusCode uint8

//go:generate stringer -type=InstanceMessageCode
const (
    SEND_MESSAGE        InstanceMessageCode = 1 + iota
    SEND_MESSAGE_REPLY
    NODE_MESSAGE_RECEIVED
    NODE_MESSAGE_SEND
    NODE_MESSAGE_SEND_REPLY
    NODE_MESSAGE_SETUP_RETRY
    NODE_MESSAGE_TIMEOUT
)

//go:generate stringer -type=InstanceStatusCode
const (
    STATUS_OK           InstanceStatusCode = 1 + iota
    STATUS_TIMEOUT
    STATUS_ERROR
)

func NewInstanceMessageSend(dest *net.UDPAddr, msg *GTPMessage, ies[]*GTPPacketIEData) InstanceMessage {
    return InstanceMessage{addr: dest, status: STATUS_OK, request: SEND_MESSAGE, msg: msg, ies: ies}
}

func (self *InstanceMessage) ToOutstandingMessage(seq uint32, teid uint32,
        dest* Destination) *OutstandingMessage {
    var bytes []byte

    switch (self.msg.version) {
    case GTPv1:
        packet := NewGTPv1Packet(false, int(seq >> 16), 0, 0, self.msg.id, teid, self.ies,
                make([]byte, 0))
        bytes = packet.Encode()
    case GTPv2:
        packet := NewGTPv2Packet(self.msg, seq, teid, self.ies)
        bytes = packet.Encode()
    }
    var timeout int = self.timeout
    var retries int = self.retry

    if seq != 0 {
        if timeout <= 0 {
            timeout = GTP_MESSAGE_TIMEOUT_DEFAULT_MS
        }
        if retries <= 0 {
            retries = GTP_MESSAGE_RETRY_DEFAULT
        }
    }
    return &OutstandingMessage{dest: dest, msg: self.msg, packet: bytes,
            sequence: seq, retries: retries, timeout: timeout, instance: self }
}

func (self *InstanceMessage) ReplyExpected() bool {
    return self.msg != nil && self.msg.response != nil
}

func (self *InstanceMessage) ShouldForward(msg *TunnelManager) bool {
    return true
}

/*
Determine whether we should reply to this message, that requires a reply.
*/
func (self *InstanceMessage) ShouldReply(node *TunnelManager) bool {
    return true
}

/*
Generate a response message from the provided InstanceMessage, and node.
*/
func (self *InstanceMessage) GenerateResponse(node *TunnelManager) *InstanceMessage {
    resp := &InstanceMessage{
        addr: self.addr,
        msg: self.msg.response,
        ies: make([]*GTPPacketIEData, 0),
        sequence: self.sequence}
    resp.fillOutMandatoryIE(node)

    var packet interface{}
    empty_tunnel_data := make([]byte, 0)
    switch (self.msg.version) {
    case GTPv1:
        packet = NewGTPv1Packet(false, int(self.sequence >> 16), 0, 0, self.msg.id,
                node.TEID(nil), resp.ies, empty_tunnel_data)
    case GTPv2:
        packet = NewGTPv2Packet(resp.msg, self.sequence,
                node.TEID(&Destination{}), resp.ies)
    }
    resp.packet = packet
    return resp
}

func (self *InstanceMessage) fillOutMandatoryIE(node *TunnelManager) {

}
