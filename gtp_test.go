package gtp

import (
    "testing"
    "fmt"
    "net"
    "time"
)

func testEq(a, b []byte) bool {
    if a == nil && b == nil {
        return true;

    }

    if a == nil || b == nil {
        return false;
    }
    if len(a) != len(b) {
        fmt.Println("length: ", len(a), len(b))
        return false
    }
    for i := range a {
        if a[i] != b[i] {
            fmt.Printf("at: %d 0x%x|0x%x\n", i, a[i], b[i])
            return false
        }
    }

    return true
}

func TestVersion1Tunnel(t *testing.T) {
    raw := []byte{  0x32,0xff,0x00,0x58,0x00,0x00,0x00,0x01,0x28,0xdb,0x00,
                    0x00,0x45,0x00,0x00,0x54,0x00,0x00,0x40,0x00,0x40,0x01,
                    0x5e,0xa5,0xca,0x0b,0x28,0x9e,0xc0,0xa8,0x28,0xb2,0x08,
                    0x00,0xbe,0xe7,0x00,0x00,0x28,0x7b,0x04,0x11,0x20,0x4b,
                    0xf4,0x3d,0x0d,0x00,0x08,0x09,0x0a,0x0b,0x0c,0x0d,0x0e,
                    0x0f,0x10,0x11,0x12,0x13,0x14,0x15,0x16,0x17,0x18,0x19,
                    0x1a,0x1b,0x1c,0x1d,0x1e,0x1f,0x20,0x21,0x22,0x23,0x24,
                    0x25,0x26,0x27,0x28,0x29,0x2a,0x2b,0x2c,0x2d,0x2e,0x2f,
                    0x30,0x31,0x32,0x33,0x34,0x35,0x36,0x37}

    p, ok := NewGTP(raw).(*GTPv1Packet)
    if !ok {
        t.Error("Failed to convert to GTPMessageV1")
    }

    if p.version != 1 {
        t.Error("Version should be 1, but it is", p.version)
    }
    // if h.hasProtol != true {
    //     t.Error("f_protocol_type should be true, but it is",
    //             h.f_protocol_type)
    // }
    if p.hasSeq() != true {
        t.Error("f_seq should be true, but it is",
                p.hasSeq())
    }
    if p.msg_type != 0xff {
        t.Error("msg_type should be 0xff, but it is",
                p.msg_type)
    }
    if p.msg_length != 88 {
        t.Error("msg_length should be 88, but it is",
                p.msg_length)
    }
    if p.teid != 0x0001 {
        t.Error("teid should be 0x0001, but it is",
                p.teid)
    }
    if p.seq != 0x28db {
        t.Error("seq should be 0x28db, but it is",
                p.seq)
    }
    if len(p.data) != 84 {
        t.Error("data length should be 80, but it is",
                len(p.data))
    }

    raw_data := p.Encode()
    if testEq(raw, raw_data) != true {
        t.Error("payload should be ", raw, ", but it is",
                raw_data)
    }

    seq := int(p.seq)
    if !p.hasSeq() {
        seq = -1
    }
    p1 := NewGTPv1Packet(false, seq, p.n_pdu, p.ext_header, p.msg_type, p.teid, p.ies, p.data)
    if p1.msg_length != 88 {
        t.Error("msg_length should be 88, but it is",
                p1.msg_length)
    }
    raw_data = p1.Encode()
    if testEq(raw, raw_data) != true {
        t.Errorf("payload should be \n[% x]\n[% x]\n",
                raw, raw_data)
    }
}

func TestVersion1IE(t *testing.T) {
    raw := []byte{ 0x32, 0x1a, 0x00, 0x10, 0x00, 0x00, 0x00, 0x00,
                   0x00, 0x00, 0x00, 0x00, 0x10, 0xe4, 0x03, 0xfb,
                   0x94, 0x85, 0x00, 0x04, 0xac, 0x13, 0x01, 0xc6}

    p, ok := NewGTP(raw).(*GTPv1Packet)
    if !ok {
        t.Error("Failed to convert to GTPMessageV1")
    }

    if p.version != 1 {
        t.Error("Version should be 1, but it is", p.version)
    }
    // if h.hasProtol != true {
    //     t.Error("f_protocol_type should be true, but it is",
    //             h.f_protocol_type)
    // }
    if p.hasSeq() != true {
        t.Error("f_seq should be true, but it is",
                p.hasSeq())
    }
    if p.msg_type != 0x1a {
        t.Error("msg_type should be 0x1a, but it is",
                p.msg_type)
    }
    if p.msg_length != 16 {
        t.Error("msg_length should be 16, but it is",
                p.msg_length)
    }
    if p.teid != 0x0 {
        t.Error("teid should be 0x0, but it is",
                p.teid)
    }
    if p.seq != 0x0 {
        t.Error("seq should be 0x28db, but it is",
                p.seq)
    }
    if len(p.data) != 12 {
        t.Error("data length should be 12, but it is",
                len(p.data))
    }

    raw_data := p.Encode()
    if testEq(raw, raw_data) != true {
        t.Error("payload should be ", raw, ", but it is",
                raw_data)
    }

    if len(p.ies) != 2 {
        t.Error("number of ies should be 2, but it is",
                len(p.ies))
    }
    seq := int(p.seq)
    if !p.hasSeq() {
        seq = -1
    }
    empty_bytes := make([]byte, 0)
    p1 := NewGTPv1Packet(false, seq, p.n_pdu, p.ext_header, p.msg_type, p.teid, p.ies, empty_bytes)
    if p1.msg_length != p.msg_length {
        t.Error("msg_length should be 88, but it is",
                p1.msg_length)
    }
    raw_data = p1.Encode()
    if testEq(raw, raw_data) != true {
        t.Errorf("payload should be \n[% x]\n[% x]\n",
                raw, raw_data)
    }
}

func TestVersion2(t *testing.T) {
    raw := []byte{  0x48,0x23,0x00,0x53,0x39,0xf0,0x00,0x05,0x00,0x1a,0xcc,0x00,
                    0x02,0x00,0x02,0x00,0x10,0x00,0x5d,0x00,0x30,0x00,0x49,0x00,
                    0x01,0x00,0x05,0x02,0x00,0x02,0x00,0x10,0x00,0x57,0x00,0x19,
                    0x00,0xc1,0x05,0x40,0x3b,0x30,0x9b,0xa5,0x26,0x65,0x26,0x06,
                    0xae,0x00,0x20,0x01,0x0b,0x00,0x00,0x00,0x00,0x00,0x00,0x00,
                    0x00,0x06,0x5e,0x00,0x04,0x00,0x05,0x00,0x00,0x0c,0x03,0x00,
                    0x01,0x00,0x38,0x48,0x00,0x08,0x00,0x00,0x00,0x61,0xa8,0x00,
                    0x01,0x11,0x70 }

    p, ok := NewGTP(raw).(*GTPv2Packet)
    if !ok {
        t.Error("Failed to convert to GTPv1Packet")
    }

    if p.version != 2 {
        t.Error("Version should be 1, but it is", p.version)
    }

    if !p.hasTEID() {
        t.Error("hasTEID should be true, but it is", p.hasTEID())
    }

    if p.msg_type != 35 {
        t.Error("msg type should be 35, but it is", p.msg_type)
    }
    if p.msg_length != 83 {
        t.Error("msg type should be 83, but it is", p.msg_length)
    }
    if p.teid != 972029957 {
        t.Error("ted should be 972029957, but it is", p.teid)
    }
    if p.seq != 6860 {
        t.Error("seq should be 6860, but it is",
                p.seq)
    }
    if len(p.data) != 75 {
        t.Error("data length should be 75, but it is",
                len(p.data))
    }
    if len(p.ies) != 4 {
        t.Error("ie length should be 4, but it is",
                len(p.ies))
    }
    raw_data := p.Encode()
    if testEq(raw, raw_data) != true {
        fmt.Printf("%x\n%x", raw, raw_data)
        t.Errorf("payload should be %x, but it is %x\n",
                raw, raw_data)
    }


    p2 := NewGTPv2Packet(TypeToGTPMessage[GTPv2][p.msg_type],
            p.seq, p.teid, p.ies)
    if p2.version != 2 {
        t.Error("Version should be 2, but it is", p2.version)
    }

    if !p2.hasTEID() {
        t.Error("hasTEID should be true, but it is", p2.hasTEID())
    }

    if p2.msg_type != 35 {
        t.Error("msg type should be 35, but it is", p2.msg_type)
    }
    if p2.msg_length != 83 {
        t.Error("msg type should be 83, but it is", p2.msg_length)
    }
    if p2.teid != 972029957 {
        t.Error("ted should be 972029957, but it is", p2.teid)
    }
    if p2.seq != 6860 {
        t.Error("seq should be 6860, but it is",
                p2.seq)
    }
    if len(p2.data) != 75 {
        t.Error("data length should be 75, but it is",
                len(p2.data))
    }
    if len(p2.ies) != 4 {
        t.Error("ie length should be 4, but it is",
                len(p2.ies))
    }
    raw_data = p2.Encode()

    if testEq(raw, raw_data) != true {
        t.Errorf("payload should be \n[% x]\n[% x]\n",
                raw, raw_data)
    }
 }

func TestTunnelManager(t *testing.T) {
    address_string := fmt.Sprintf("%s:%d", "127.0.0.1", 2123)
    listen_addr, err := net.ResolveUDPAddr("udp", address_string)
    if err != nil {
        panic(err)
    }
    from_child := make(chan InstanceMessage)
    to_child := make(chan InstanceMessage)
    go NewNode(listen_addr, from_child, to_child)

    address_string = fmt.Sprintf("%s:%d", "127.0.0.1", 5000)
    receiver_addr, err := net.ResolveUDPAddr("udp", address_string)
    if err != nil {
        panic(err)
    }
    l := make(chan InstanceMessage)
    to_child<-NewInstanceMessageSend(receiver_addr,
            GTPV2_MSG_ECHO_REQUEST,
            make([]*GTPPacketIEData, 0))
    var reply_count = 0

WAIT:
    for {
        select {
        case msg := <-from_child:
            switch (msg.request) {
            case SEND_MESSAGE_REPLY:
                reply_count++
                fmt.Printf("[%d] reply", reply_count)
                if (reply_count == 1) {
                    if msg.status != STATUS_TIMEOUT {
                        t.Errorf("[%d] expected %s, but got %s",
                                reply_count, STATUS_ERROR, msg.status)
                    }
                    // setup listener here, and try again
                    go NewNode(receiver_addr, from_child, to_child)
                    to_child<-NewInstanceMessageSend(receiver_addr,
                            GTPV2_MSG_ECHO_REQUEST,
                            make([]*GTPPacketIEData, 0))
                } else if (reply_count > 1) {
                    if msg.status != STATUS_OK {
                        t.Errorf("[%d] expected timeout, but got %s",
                                reply_count, msg.status)
                    }
                    break WAIT
                }
            default:
                t.Error("unexpected reply received: ", msg.request, msg.status.String())
            }
            fmt.Println("Reply received: ", msg.status)
        case msg := <-l:
            fmt.Println("Second message: ", msg.request, msg.status)
        case <-time.After(50000 * time.Millisecond):
            fmt.Println("TIMEOUT")
            t.Error("timeout after 500ms")
            break WAIT
            // blocks
        }
    }
}

// func TestGTPListenerPeer(t *testing.T) {
//     address_string := fmt.Sprintf("%s:%d", "127.0.0.1", 2123)
//     listen_addr, err := net.ResolveUDPAddr("udp", address_string)
//     if err != nil {
//         panic(err)
//     }
//     c := make(chan InstanceMessage)
//     go NewListener(listen_addr, c)
//
//     address_string = fmt.Sprintf("%s:%d", "127.0.0.1", 5000)
//     receiver_addr, err := net.ResolveUDPAddr("udp", address_string)
//     if err != nil {
//         panic(err)
//     }
//     l := make(chan InstanceMessage)
    // NewPeer(listen_addr, receiver_addr, l)

// WAIT:
//     for {
//         select {
//         case msg := <-c:
//             switch (msg.request) {
//             case SEND_MESSAGE_REPLY:
//                 reply_count++
//                 fmt.Printf("[%d] reply", reply_count)
//                 if (reply_count == 1) {
//                     if msg.status != STATUS_TIMEOUT {
//                         t.Errorf("[%d] expected %s, but got %s",
//                                 reply_count, STATUS_ERROR, msg.status)
//                     }
//                     // setup listener here, and try again
//                     go NewNode(receiver_addr, l)
//                     c<-NewInstanceMessageSend(receiver_addr,
//                             GTPV2_MSG_ECHO_REQUEST,
//                             make([]*GTPPacketIEData, 0))
//                 } else if (reply_count > 1) {
//                     if msg.status != STATUS_OK {
//                         t.Errorf("[%d] expected timeout, but got %s",
//                                 reply_count, msg.status)
//                     }
//                     break WAIT
//                 }
//             default:
//                 t.Error("unexpected reply received: ", msg.request, msg.status.String())
//             }
//             fmt.Println("Reply received: ", msg.status)
//         case msg := <-l:
//             fmt.Println("Second message: ", msg.request, msg.status)
//         case <-time.After(50000 * time.Millisecond):
//             fmt.Println("TIMEOUT")
//             t.Error("timeout after 500ms")
//             break WAIT
//             // blocks
//         }
//     }
//}
