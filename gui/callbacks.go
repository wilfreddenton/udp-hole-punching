package main

import (
	"fmt"

	"github.com/googollee/go-socket.io"
	"github.com/wilfreddenton/udp-hole-punching/shared"
)

func createRegisteredCallback(so socketio.Socket) func(shared.Client) {
	return func(c shared.Client) {
		so.Emit("enter", c.GetSelf().ID)
	}
}

func createConnectingCallback(so socketio.Socket) func(shared.Client) {
	return func(c shared.Client) {
		peer := c.GetPeer()
		pConn := c.GetPeerConn()

		fmt.Println("connecting")
		so.Emit("connecting", fmt.Sprintf(`{
			"username": "%s",
			"id": "%s",
			"addr": "%s"
		}`, peer.Username, peer.ID, pConn.GetAddr()))
	}
}

func createConnectedCallback(so socketio.Socket) func(shared.Client) {
	return func(c shared.Client) {
		so.Emit("connected")
	}
}

func createMessageCallback(so socketio.Socket) func(shared.Client, string) {
	return func(c shared.Client, text string) {
		so.Emit("message", text)
	}
}
