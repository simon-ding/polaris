package nat

import (
	"fmt"
	"net"
	"polaris/log"

	"github.com/pion/stun/v3"
)

const (
	udp           = "udp4"
	pingMsg       = "ping"
	pongMsg       = "pong"
	timeoutMillis = 500
)

func NewNatTraversal() (*NatTraversal, error) {
	conn, err := net.ListenUDP(udp, nil)
	if err != nil {
		return nil, fmt.Errorf("listen: %w", err)
	}

	log.Infof("Listening on %s", conn.LocalAddr())

	messageChan := listen(conn)

	return &NatTraversal{
		conn:        conn,
		messageChan: messageChan,
		cancel:      make(chan struct{}),
	}, nil
}

type NatTraversal struct {
	//peerAddr    *net.UDPAddr
	conn        *net.UDPConn
	messageChan <-chan []byte
	cancel      chan struct{}

	stunAddr *stun.XORMappedAddress
}

func (s *NatTraversal) Cancel() {
	
	close(s.cancel)
	s.conn.Close()
}


func (s *NatTraversal) StunAddr() (*stun.XORMappedAddress, error) {
	for _, srv := range getStunServers() {
		log.Debugf("try to connect to stun server: %s", srv)
		srvAddr, err := net.ResolveUDPAddr(udp, srv)
		if err != nil {
			log.Warnf("Failed to resolve server addr: %s", err)
			continue
		}
		err = sendBindingRequest(s.conn, srvAddr)
		if err != nil {
			return nil, fmt.Errorf("send binding request: %w", err)
		}
		select {
		case message, ok := <-s.messageChan:
			if !ok {
				continue
			}
			if stun.IsMessage(message) {
				m := new(stun.Message)
				m.Raw = message
				decErr := m.Decode()
				if decErr != nil {
					log.Warnf("decode:", decErr)

					break
				}
				var xorAddr stun.XORMappedAddress
				if getErr := xorAddr.GetFrom(m); getErr != nil {
					log.Warnf("getFrom:", getErr)

					continue
				}
				if s.stunAddr == nil || s.stunAddr.String() != xorAddr.String() {
					log.Warnf("My public address: %s\n", xorAddr)
					s.stunAddr = &xorAddr
				}
				return &xorAddr, nil

			}
		}

	}
	return nil, fmt.Errorf("failed to get STUN address")
}

func (s *NatTraversal) StartProxy(targetAddr string) error {
	log.Infof("Starting NAT traversal proxy to %s", targetAddr)
	peerAddr, err := net.ResolveUDPAddr(udp, targetAddr)
	if err != nil {
		return fmt.Errorf("resolve peeraddr: %w", err)
	}

	if s.stunAddr == nil {
		addr, err := s.StunAddr()
		if err != nil {
			return fmt.Errorf("get STUN address: %w", err)
		}
		log.Infof("STUN address: %s", addr)
	}
	for {
		select {
		case <-s.cancel:
			log.Infof("cancelled")
			return nil
		case m := <-s.messageChan:
			//log.Infof("Received message: %d", len(m))
			send(m, s.conn, peerAddr)
		}
	
	}
}

func listen(conn *net.UDPConn) <-chan []byte {
	messages := make(chan []byte)
	go func() {
		for {
			buf := make([]byte, 10240)

			n, addr, err := conn.ReadFromUDP(buf)
			if err != nil {
				close(messages)

				return
			}
			log.Debugf("Received message from %s: %d", addr, n)
			buf = buf[:n]
			log.Debugf("recevied message %s", string(buf))

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
