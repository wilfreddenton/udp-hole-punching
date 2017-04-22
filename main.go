package main

import (
	"fmt"
	"log"
	"net"

	"github.com/wilfreddenton/crypto"
	"github.com/wilfreddenton/udp-hole-punching/shared"
	"github.com/wilfreddenton/udp-hole-punching/udp_server"
)

var (
	pubKey [32]byte
	priKey [32]byte
)

func route(peers shared.Peers, conns shared.Conns, conn shared.Conn, m *shared.Message) (*shared.Message, error) {
	switch m.Type {
	case "greeting":
		return greetingHandler(conn, m)
	case "register":
		return registerHandler(peers, conn, m)
	case "establish":
		return establishHandler(peers, conns, m)
	default:
		return notFoundHandler(m)
	}
}

func createMessageCallback(peers shared.Peers) func(cs shared.Conns, c shared.Conn, m *shared.Message) {
	return func(cs shared.Conns, c shared.Conn, m *shared.Message) {
		// log request
		log.Printf("Request from client at %s over %s with type %s", c.GetAddr(), c.Protocol(), m.Type)

		// route request to a handler
		res, err := route(peers, cs, c, m)

		// respond with error if there was one
		if err != nil {
			c.Send(&shared.Message{
				Type:  m.Type,
				Error: err.Error(),
			})
			return
		}

		// respond
		err = c.Send(res)
		if err != nil {
			log.Print(err)
		}
	}
}

func main() {
	fmt.Println("UDP Hole Punching Rendezvous Server v0.0.1")

	var err error
	priKey, pubKey, err = crypto.GenKeyPair()
	if err != nil {
		log.Fatal(err)
	}

	udpAddr, err := net.ResolveUDPAddr("udp", "0.0.0.0:9001")
	if err != nil {
		log.Fatal(err)
	}

	udpS, err := udp_server.New(udpAddr)
	if err != nil {
		log.Fatal(err)
	}

	udpPeers := make(shared.Peers)
	udpS.OnMessage(createMessageCallback(udpPeers))
	udpS.Listen()
}
