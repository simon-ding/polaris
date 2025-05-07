package nat

import (
	"fmt"
	"log"
	"net"
	"time"

	"github.com/pion/stun/v3"
)

const (
	udp           = "udp4"
	pingMsg       = "ping"
	pongMsg       = "pong"
	timeoutMillis = 500
)

type natTraversal struct {
	peerAddr *net.UDPAddr
	cancel   chan struct{}
	port     <-chan int
}

func (s *natTraversal) Port() int {
	return <-s.port
}

func (s *natTraversal) Cancel() {
	s.cancel <- struct{}{}
}

func NatTraversal(targetAddr string) (*natTraversal, error) { //nolint:gocognit,cyclop

	srvAddr, err := net.ResolveUDPAddr(udp, getStunServers()[0])
	if err != nil {
		log.Fatalf("Failed to resolve server addr: %s", err)
	}

	conn, err := net.ListenUDP(udp, nil)
	if err != nil {
		return nil, fmt.Errorf("listen: %w", err)
	}

	log.Printf("Listening on %s", conn.LocalAddr())

	peerAddr, err := net.ResolveUDPAddr(udp, targetAddr)
	if err != nil {
		return nil, fmt.Errorf("resolve peeraddr: %w", err)
	}
	err = sendBindingRequest(conn, srvAddr)
	if err != nil {
		return nil, fmt.Errorf("send binding request: %w", err)
	}
	nt := &natTraversal{
		peerAddr: peerAddr,
		cancel:   make(chan struct{}),
		port:     make(chan int),
	}
	go func() {
		err := doTraversal(conn, peerAddr, nt.cancel)
		if err != nil {
			log.Println("nat traversal error:", err)
		}
	}()
	return nt, nil

}

func doTraversal(conn *net.UDPConn, peerAddr *net.UDPAddr, quit <-chan struct{}) error {
	defer func() {
		_ = conn.Close()
	}()

	var publicAddr stun.XORMappedAddress

	messageChan := listen(conn)
	//var peerAddrChan <-chan string

	keepalive := time.Tick(timeoutMillis * time.Millisecond)
	keepaliveMsg := pingMsg

	gotPong := false
	sentPong := false

	for {
		select {
		case message, ok := <-messageChan:
			if !ok {
				return nil
			}

			switch {
			case string(message) == pingMsg:
				keepaliveMsg = pongMsg

			case string(message) == pongMsg:
				if !gotPong {
					log.Println("Received pong message.")
				}

				// One client may skip sending ping if it receives
				// a ping message before knowning the peer address.
				keepaliveMsg = pongMsg

				gotPong = true

			case stun.IsMessage(message):
				m := new(stun.Message)
				m.Raw = message
				decErr := m.Decode()
				if decErr != nil {
					log.Println("decode:", decErr)

					break
				}
				var xorAddr stun.XORMappedAddress
				if getErr := xorAddr.GetFrom(m); getErr != nil {
					log.Println("getFrom:", getErr)

					break
				}

				if publicAddr.String() != xorAddr.String() {
					log.Printf("My public address: %s\n", xorAddr)
					publicAddr = xorAddr

					//peerAddrChan = getPeerAddr()
				}

			default:
				send(message, conn, peerAddr)
			}

		case <-keepalive:
			// Keep NAT binding alive using STUN server or the peer once it's known
			err := sendStr(keepaliveMsg, conn, peerAddr)
			if keepaliveMsg == pongMsg {
				sentPong = true
			}
			_ = sentPong

			if err != nil {
				log.Panicln("keepalive:", err)
			}

		case <-quit:
			_ = conn.Close()
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
