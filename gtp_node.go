package gtp

import (
    "os"
    "net"
    "log"
    "time"
    "fmt"
    "bytes"
    "utils"
    "encoding/binary"
    "crypto/rand"
)

/*
29274 - e00
7.6 Reliable Delivery of Signalling Messages

Retransmission requirements in the current subclause do not apply to the Initial
messages that do not have Triggered messages. Reliable delivery in GTPv2
messages is accomplished by retransmission of these messages. A message shall be
retransmitted if and only if a reply is expected for that message and the reply
has not yet been received. There may be limits placed on the total number of
retransmissions to avoid network overload.

Initial messages and their Triggered messages, as well as Triggered messages and
their Triggered Reply messages are matched based on the Sequence Number and the
IP address and port rules in subclause 4.2 "Protocol stack". Therefore, an
Initial message and its Triggered message, as well as a Triggered message and
its Triggered Reply message shall have exactly the same Sequence Number value. A
retransmitted GTPv2 message (an Initial or a Triggered) has the exact same GTPv2
message content, including the GTP header, UDP ports, source and destination IP
addresses as the originally transmitted GTPv2 message.

For each triplet of local IP address, local UDP port and remote peer's IP
address a GTP entity maintains a sending queue with signalling messages to be
sent to that peer. The message at the front of the queue shall be sent with a
Sequence Number, and if the message has an expected reply, it shall be held in a
list until a reply is received or until the GTP entity has ceased retransmission
of that message. The Sequence Number shall be unique for each outstanding
Initial message sourced from the same IP/UDP endpoint. A node running GTP may
have several outstanding messages waiting for replies. Not counting
retransmissions, a single GTP message with an expected reply shall be answered
with a single GTP reply, regardless whether it is per UE, per APN, or per bearer
*/

/* the TunnelManager shall keep track of the outstanding messages, and handle
retransmissions, as well as validating any reply (as indicated in Section 11) */
type TunnelManager struct {
    addr *net.UDPAddr
    destinations map[*net.UDPAddr]*Destination
    messages map[uint32]*OutstandingMessage
                                            /* "list" of messages map to the
                                               sequence number. Note: GTPv1 seq
                                               is 16bits, while GTPv2 is 24bits
                                            */
    server *net.UDPConn                     /* listening server conn */
    caller_channel chan<-InstanceMessage
    action_channel chan InstanceMessage
    timer chan *time.Time
    log *log.Logger
}

// each destination have it's own unique teid
type Destination struct {
    addr *net.UDPAddr
    connection *net.UDPConn
    teid uint32
}

type OutstandingMessage struct {
    dest *Destination
    msg *GTPMessage
    packet []byte
    sequence uint32
    timeout int
    retries int
    timer *time.Timer
    instance *InstanceMessage
}


const (
    GTP_MESSAGE_TIMEOUT_DEFAULT_MS  = int(10)
    GTP_MESSAGE_RETRY_DEFAULT       = int(5)
)

/* Use for messaging between nodes */
type InstanceMessage struct {
    request InstanceMessageCode
    status InstanceStatusCode
    addr *net.UDPAddr            /* message destination */
    msg *GTPMessage             /* GTP message to send */
    ies []*GTPPacketIEData      /* GTP IE data */
    packet interface{}          /* GTPPacket */
    timeout int                 /* send message timeout */
    retry int                   /* send message retry */
    err error
    sequence uint32             /* sequence number */
}

func randomBytes(n int) ([]byte, error) {
    b := make([]byte, n)
    _, err := rand.Read(b)
    if err != nil {
        return nil, err
    }

    return b, nil
}

/* Generate random sequence number base on GTPVersion */
func randomSequence(version GTPVersion) (uint32, error) {
    b, err := randomBytes(4)
    if err != nil {
        return 0, err
    }
    var val uint32
	buf := bytes.NewReader(b)
	err = binary.Read(buf, binary.LittleEndian, &val)
	if err != nil {
        return 0, nil
	}
    return val, nil
}

func NewNode(local *net.UDPAddr,
        to_parent chan<-InstanceMessage,
        from_parent <-chan InstanceMessage) *TunnelManager {
    self := &TunnelManager{addr: local,
            messages: make(map[uint32]*OutstandingMessage),
            destinations: make(map[*net.UDPAddr]*Destination),
            caller_channel: to_parent,
            action_channel: make(chan InstanceMessage),
            timer: nil}
    perfix := fmt.Sprintf("C[%s] ", self.addr)
    self.log = log.New(os.Stdout, perfix,
            log.Ldate|log.Ltime|log.Lshortfile)

    incoming := make(chan InstanceMessage)
    go self.listen(incoming)
    for {
        select {
        case msg := <-from_parent:
            self.log.Printf("parent message: %s", msg.request)
            go self.handle_parent_message(&msg)
        case msg := <-incoming:
            self.log.Printf("incoming: %s %s", msg.request, msg.msg.name)
            // handle all incoming messages
            self.handle_incoming_message(&msg)
        case msg := <-self.action_channel:
            self.log.Printf("self: %s", msg.request)
            go self.handle_self_message(&msg)
        }
    }
    return self
}

/*
Set up listening for UDP packet.
*/
func (self *TunnelManager) listen(outgoing chan<-InstanceMessage) {
    var err error = nil
    self.server, err = net.ListenUDP("udp", self.addr)
    utils.CheckError(err)
    self.log.Printf("listening")
    buf := make([]byte, GTP_NODE_READ_BUFFER_SIZE)
    for {
        _, addr, err := self.server.ReadFromUDP(buf)
        // XXX need to check n == GTP_NODE_READ_BUFFER_SIZE, and handle
        // multiple reads
        if err != nil {
            panic(err)
        }
        outgoing <- bufferToInstanceMessage(buf, addr)
    }
}


/*
Set up listening for UDP packet.
*/
const GTP_NODE_READ_BUFFER_SIZE = int(1024)
//
// func (self *TunnelManager) listen(outgoing chan<-InstanceMessage) {
//     var err error = nil
//     self.server, err = net.ListenUDP("udp", self.addr)
//     if err != nil {
//         panic(err)
//     }
//     self.log.Printf("Listening")
//     defer self.server.Close()
//     buf := make([]byte, GTP_NODE_READ_BUFFER_SIZE)
//     for {
//         _, addr, err := self.server.ReadFromUDP(buf)
//         // XXX need to check n == GTP_NODE_READ_BUFFER_SIZE, and handle
//         // multiple reads
//         utils.CheckError(err)
//         outgoing <- bufferToInstanceMessage(buf, addr)
//     }
// }

func bufferToInstanceMessage(buf []byte,
        addr *net.UDPAddr) InstanceMessage {
    packet := NewGTP(buf)
    return InstanceMessage{request: NODE_MESSAGE_RECEIVED,
            addr: addr, packet: packet,
            sequence: packetSequenceFromInterface(packet),
            msg: GTPMessageFromInterface(packet)}
}


func (self *TunnelManager) handle_self_message(msg *InstanceMessage) {
    switch (msg.request) {
    case NODE_MESSAGE_SEND:
        // given a sequence number, send the OutstandingMessage
        out_msg, ok := self.messages[msg.sequence]
        if !ok {
            fmt.Println("Error sequence [%d] does not exist for SEND_MESSAGE_NODE",
                    msg.sequence)
        }
        self.Send(out_msg)
    case NODE_MESSAGE_SETUP_RETRY:
        // given a sequence number, send the OutstandingMessage
        out_msg, ok := self.messages[msg.sequence]
        if !ok {
            fmt.Println("Error sequence [%d] does not exist for SEND_MESSAGE_NODE",
                    msg.sequence)
            return
        }
        out_msg.timer = time.NewTimer(time.Duration(out_msg.timeout) * time.Millisecond)
        go func() {
            <- out_msg.timer.C
            out_msg.retries--
            self.Send(out_msg)
        }()
    case NODE_MESSAGE_TIMEOUT:
        out_msg, ok := self.messages[msg.sequence]
        if !ok {
            fmt.Println("Error sequence [%d] does not exist for SEND_MESSAGE_NODE",
                    msg.sequence)
            return
        }
        out_msg.instance.request = SEND_MESSAGE_REPLY
        out_msg.instance.status = STATUS_TIMEOUT
        self.caller_channel<-*out_msg.instance
    }
}

func (self *TunnelManager) handle_incoming_message(msg *InstanceMessage) {
    self.log.Printf("message sequence [%d]", msg.sequence)
    if msg.sequence != 0 {
        out_msg, ok := self.messages[msg.sequence]
        if ok {
            // stop retries
            out_msg.timer.Stop()
            delete(self.messages, msg.sequence)
            self.log.Printf("delete outstanding message [%d]", msg.sequence)
            // forward it back to parent
            if msg.ShouldForward(self) {
                self.log.Printf("forwarding message [%d]: %s",
                        msg.sequence, msg.msg.name)
                msg.request = SEND_MESSAGE_REPLY
                msg.status = STATUS_OK
                self.caller_channel<-*msg
            }
        } else {
            if msg.ReplyExpected() && msg.ShouldReply(self) {
                reply := msg.GenerateResponse(self)
                self.SendReply(reply)
            } else {
                self.log.Printf("NOT generating reply")
            }
        }
    }
}

func (self *TunnelManager) handle_parent_message(msg *InstanceMessage) {
    dest := self.findOrMakeDestination(msg.addr)
    switch (msg.request) {
    case SEND_MESSAGE:
        var seq uint32 = 0
        if msg.ReplyExpected() {
            // XXX ignores the same message type on the same destination
            // for now.  At one point, we'll have to figure out what to do in
            // that instance.
            ok := bool(true)
            for ok {
                // generate unique sequence number for this node
                seq, _ = randomSequence(GTPv2)
                switch (msg.msg.version) {
                case GTPv1:
                    seq &= GTPV1_VERSION_UINT32_MASK
                case GTPv2:
                    seq &= GTPV2_VERSION_UINT32_MASK
                }
                _, ok = self.messages[seq]
            }
        }
        teid := self.TEID(dest)
        out_msg := msg.ToOutstandingMessage(seq, teid, dest)
        if seq != 0 {
            self.messages[seq] = out_msg
        }
        self.Send(out_msg)
    }
}

func (self *TunnelManager) TEID(dest *Destination) uint32 {
    // XXX this should be dyanmic
    // see 29274-5.5.3
    return dest.teid
}

func (self *TunnelManager) findOrMakeDestination(addr *net.UDPAddr) *Destination {
    dest, ok := self.destinations[addr]
    if !ok {
        // XXX ignore error for now
        teid, _ := randomSequence(GTPv2)
        dest = &Destination{addr: addr, teid: teid}
        self.destinations[addr] = dest
    }
    return dest
}

func (self *TunnelManager) reply_status(msg_code InstanceMessageCode,
        msg_status InstanceStatusCode, err error) {
    fmt.Println("reply status: ", msg_code, msg_status, err)
    self.caller_channel <- InstanceMessage{err: err, request: msg_code, status: msg_status}
}

/* Send a GTP message to the addr, the node will take care of retry, and */
func (self *TunnelManager) Send(msg *OutstandingMessage) {
    self.log.Printf("sending message type [%s] {teid:%d seq:%d}",
        msg.msg.name, msg.dest.teid, msg.sequence)
    _, err := self.server.WriteTo(msg.packet, msg.dest.addr)
    if err != nil {
        self.reply_status(SEND_MESSAGE_REPLY, STATUS_ERROR, err)
        return
    }
    if msg.sequence != 0 {
        if msg.retries > 0 {
            self.action_channel <- InstanceMessage{request: NODE_MESSAGE_SETUP_RETRY,
                    sequence: msg.sequence}
        } else {
            self.action_channel <- InstanceMessage{request: NODE_MESSAGE_TIMEOUT,
                    sequence: msg.sequence}
        }
    } else {
        self.action_channel <- InstanceMessage{request: NODE_MESSAGE_SEND_REPLY,
                sequence: msg.sequence}
    }
}

/* Send a GTP message to the addr, the node will take care of retry, and */
func (self *TunnelManager) SendReply(msg *InstanceMessage) {
    var packet []byte

    switch p := msg.packet.(type) {
    case *GTPv1Packet:
        packet = p.Encode()
    case *GTPv2Packet:
        packet = p.Encode()
    }
    _, err := self.server.WriteTo(packet, msg.addr)
    if err != nil {
        self.reply_status(SEND_MESSAGE_REPLY, STATUS_ERROR, err)
    } else {
        self.reply_status(SEND_MESSAGE_REPLY, STATUS_OK, nil)
    }
}
