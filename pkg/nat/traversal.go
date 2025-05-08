package nat

import (
	"fmt"
	"net"
	"polaris/log"
	"time"

	"github.com/pion/stun/v3"
)

const (
	udp           = "udp4"
	pingMsg       = "ping"
	pongMsg       = "pong"
	timeoutMillis = 500
)

func NewNatTraversal(addrCallback func(stun.XORMappedAddress) error, targetHost string) (*NatTraversal, error) {
	conn, err := net.ListenUDP(udp, nil)
	if err != nil {
		return nil, fmt.Errorf("listen: %w", err)
	}

	log.Infof("Listening on %s", conn.LocalAddr())

	messageChan := listen(conn)
	s := &NatTraversal{
		conn:         conn,
		messageChan:  messageChan,
		cancel:       make(chan struct{}),
		addrChan:     make(chan stun.XORMappedAddress),
		addrCallback: addrCallback,
		targetHost:   targetHost,
	}

	go s.updateNatAddr()

	return s, nil
}

type NatTraversal struct {
	//peerAddr    *net.UDPAddr
	conn        *net.UDPConn
	messageChan <-chan []byte
	addrChan    chan stun.XORMappedAddress
	cancel      chan struct{}

	stunAddr     *stun.XORMappedAddress
	addrCallback func(stun.XORMappedAddress) error
	targetHost   string
	targetPort   int
}

func (s *NatTraversal) Cancel() {

	close(s.cancel)
	s.conn.Close()
}

func (s *NatTraversal) updateNatAddr() {
	for addr := range s.addrChan {
		if s.stunAddr == nil || s.stunAddr.String() != addr.String() { //new address
			log.Warnf("My public address: %s\n", addr)
			if s.addrCallback != nil { //execute callback
				if err := s.addrCallback(addr); err != nil {
					log.Warnf("callback error: %v", err)
				}
			}

			s.targetPort = addr.Port
			log.Infof("now proxy to target host: %s:%d", s.targetHost, s.targetPort)
			s.stunAddr = &addr
		}
	}
}

func (s *NatTraversal) sendStunServerBindingMsg() error {
	for _, srv := range getStunServers() {
		log.Debugf("try to connect to stun server: %s", srv)
		srvAddr, err := net.ResolveUDPAddr(udp, srv)
		if err != nil {
			log.Warnf("Failed to resolve server addr: %s", err)
			continue
		}
		err = sendBindingRequest(s.conn, srvAddr)
		if err != nil {
			log.Warnf("send binding request: %w", err)
			continue
		}
		return nil
	}
	return fmt.Errorf("failed to get STUN address")
}

func (s *NatTraversal) getNatAddr(msg []byte) (*stun.XORMappedAddress, error) {
	if !stun.IsMessage(msg) {
		return nil, fmt.Errorf("not a stun message")
	}

	m := new(stun.Message)
	m.Raw = msg
	decErr := m.Decode()
	if decErr != nil {
		return nil, fmt.Errorf("decode: %w", decErr)
	}
	var xorAddr stun.XORMappedAddress
	if getErr := xorAddr.GetFrom(m); getErr != nil {
		return nil, fmt.Errorf("getFrom: %w", getErr)
	}
	s.addrChan <- xorAddr

	return &xorAddr, nil

}

func (s *NatTraversal) StartProxy() {

	tick := time.NewTicker(10 * time.Second)

	go func() { //tcker message to check public ip and port
		defer tick.Stop()
		for {
			select {
			case <-s.cancel:
				log.Infof("stun nat proxy cancelled")
				return
			case <-tick.C:
				err := s.sendStunServerBindingMsg()
				if err != nil {
					log.Warnf("send stun server binding msg: %w", err)
				}
			}
		}
	}()

	for {
		select {
		case <-s.cancel:
			log.Infof("stun nat proxy cancelled")
			return
		case m := <-s.messageChan:
			if stun.IsMessage(m) {
				s.getNatAddr(m)
			} else {
				peerAddr, err := net.ResolveUDPAddr(udp, fmt.Sprintf("%s:%d", s.targetHost, s.targetPort))
				if err != nil {
					log.Errorf("resolve peeraddr: %w", err)
					continue
				}

				send(m, s.conn, peerAddr)
			}
		}

	}
}

func listen(conn *net.UDPConn) <-chan []byte {
	messages := make(chan []byte)
	go func() {
		for {
			buf := make([]byte, 10240)

			n, _, err := conn.ReadFromUDP(buf)
			if err != nil {
				close(messages)

				return
			}
			buf = buf[:n]

			messages <- buf
		}
	}()

	return messages
}

func sendBindingRequest(conn *net.UDPConn, addr *net.UDPAddr) error {
	m := stun.MustBuild(stun.TransactionID, stun.BindingRequest)

	err := send(m.Raw, conn, addr)
	if err != nil {
		return fmt.Errorf("binding: %w", err)
	}

	return nil
}

func send(msg []byte, conn *net.UDPConn, addr *net.UDPAddr) error {
	_, err := conn.WriteToUDP(msg, addr)
	if err != nil {
		return fmt.Errorf("send: %w", err)
	}

	return nil
}

func sendStr(msg string, conn *net.UDPConn, addr *net.UDPAddr) error {
	return send([]byte(msg), conn, addr)
}
