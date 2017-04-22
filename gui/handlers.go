package main

import (
	"fmt"
	"log"
	"net"

	"github.com/googollee/go-socket.io"
	"github.com/wilfreddenton/udp-hole-punching/udp_client"
)

func enterHandler(so socketio.Socket, s *state, u *user) {
	s.username = u.Username

	var err error
	var sAddr *net.UDPAddr
	var addr *net.UDPAddr
	sAddr, err = net.ResolveUDPAddr("udp", serverUDPIP+serverUDPPort)
	if err != nil {
		log.Fatal(err)
	}

	// create self address
	addr, err = net.ResolveUDPAddr("udp", ":9002")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("set client")
	s.client, err = udp_client.New(s.username, addr, sAddr)
	err = s.client.Start()
	if err != nil {
		log.Print(err)
		so.Emit("error", err.Error())
		return
	}

	s.client.OnRegistered(createRegisteredCallback(so))
	s.client.OnConnecting(createConnectingCallback(so))
	s.client.OnConnected(createConnectedCallback(so))
	s.client.OnMessage(createMessageCallback(so))

	s.id = s.client.GetSelf().ID
}
