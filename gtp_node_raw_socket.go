package gtp

//import (
//    "net"
//    "utils"
//    "time"
//    "log"
//    "golang.org/x/sys/unix"
//)
//
//
//
//type listenerUDP struct {
//    fd int
//}
//
//type connectionUDP struct {
//    fd int
//    local unix.Sockaddr
//    remote unix.Sockaddr
//}
//
//type Peer struct {
//    local *net.UDPAddr
//    remote *net.UDPAddr
//    messages map[uint32]*OutstandingMessage
//                                            /* "list" of messages map to the
//                                               sequence number. Note: GTPv1 seq
//                                               is 16bits, while GTPv2 is 24bits
//                                            */
//    connection *connectionUDP               /* concrete connection */
//    caller_channel chan<-InstanceMessage
//    action_channel chan InstanceMessage
//    timer chan *time.Time
//    teid uint32
//    log *log.Logger
//}
//
//func ListenUDP(proto string, addr *net.UDPAddr) (*listenerUDP, error) {
//    self := &listenerUDP{fd: 0}
//    var err error = nil
//
//    self.fd, err = unix.Socket(unix.AF_INET, unix.SOCK_DGRAM, unix.IPPROTO_UDP)
//    if err != nil {
//        return self, err
//    }
//    utils.CheckError(unix.SetsockoptInt(self.fd, unix.SOL_SOCKET, unix.SO_REUSEPORT, 1))
//    sa := localUDPAddrToSockaddr(addr)
//    utils.CheckError(unix.Bind(self.fd, sa))
//    return self, err
//}
//
//func localUDPAddrToSockaddr(addr *net.UDPAddr) *unix.SockaddrInet4 {
//    ipv4_addr := []byte(addr.IP.To4())
//    ipv4_addr_byte := [4]byte{0, 0, 0, 0}
//    copy(ipv4_addr_byte[:], ipv4_addr[0:3])
//    return &unix.SockaddrInet4{Addr: ipv4_addr_byte, Port: addr.Port}
//}
//
//func socketaddrToUDPAddr(sa unix.Sockaddr) *net.UDPAddr {
//    switch a := sa.(type) {
//    case *unix.SockaddrInet4:
//        ip := net.IP(a.Addr[:])
//        return &net.UDPAddr{IP: ip, Port: a.Port}
//    case *unix.SockaddrInet6:
//        ip := net.IP(a.Addr[:])
//        return &net.UDPAddr{IP: ip, Port: a.Port, Zone: string(a.ZoneId)}
//    }
//    return nil
//}
//
//func (self *listenerUDP) ReadFromUDP(buf []byte) (int, *net.UDPAddr, error){
//    n, from, err := unix.Recvfrom(self.fd, buf, 0)
//    if err != nil {
//        return n, nil, err
//    }
//    return n, socketaddrToUDPAddr(from), err
//}
//
//func NewPeer(local *net.UDPAddr, remote *net.UDPAddr,
//        parent chan InstanceMessage) {
//    self := &Peer{local: local, remote: remote,
//            connection: &connectionUDP{fd:0}}
//
//    var err error
//    self.connection, err = self.DialUDP("udp", local, remote)
//    utils.CheckError(err)
//}
//
//func (self *Peer) DialUDP(proto string, local *net.UDPAddr,
//        remote *net.UDPAddr) (*connectionUDP, error) {
//    c := &connectionUDP{fd: 0}
//
//    var err error = nil
//    c.fd, err = unix.Socket(unix.AF_INET, unix.SOCK_DGRAM, unix.IPPROTO_UDP)
//    if err != nil {
//        return c, err
//    }
//    c.local = localUDPAddrToSockaddr(local)
//    c.remote = localUDPAddrToSockaddr(remote)
//    err = unix.SetsockoptInt(c.fd, unix.SOL_SOCKET, unix.SO_REUSEPORT, 1)
//    if err != nil {
//        return c, err
//    }
//    err = unix.Bind(c.fd, c.local)
//    return c, err
//}
