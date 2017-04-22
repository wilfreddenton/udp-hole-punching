package udp_client

import (
	"encoding/base64"
	"net"
	"time"

	"github.com/wilfreddenton/udp-hole-punching/base_client"
	"github.com/wilfreddenton/udp-hole-punching/shared"
	"github.com/wilfreddenton/udp-hole-punching/udp_server"
)

type Client struct {
	*base_client.Client
	sAddr *net.UDPAddr
}

func (c *Client) Connect() {
	l := c.GetLog()
	self := c.GetSelf()
	peer := c.GetPeer()
	pConn := c.GetPeerConn()

	for i := 0; i < 5; i += 1 {
		if c.WasKeyReceived() {
			// tell user that client connected to peer
			l.Printf("connected to peer %s", peer.Username)
			c.ConnectedCallback(c)
			return
		}

		l.Printf("punching through to peer %s at %s", peer.Username, pConn.GetAddr())
		pConn.Send(&shared.Message{
			Type:   "connect",
			PeerID: self.ID,
		})
		time.Sleep(3 * time.Second)
	}

	l.Printf("could not connect to peer %s at %s", peer.Username, pConn.GetAddr())
}

func (c *Client) Start() error {
	s := c.GetServer()

	// add rendezvous server connection
	sConn, err := s.CreateConn(c.sAddr)
	if err != nil {
		return err
	}

	c.SetServerConn(sConn)

	// get public key
	pubKey, err := c.GetSelf().GetPublicKey()
	if err != nil {
		return err
	}

	// start server
	go s.Listen()

	// send greeting message to server
	sConn.Send(&shared.Message{
		Type:    "greeting",
		Content: base64.StdEncoding.EncodeToString(pubKey[:]),
	})

	return nil
}

func New(username string, addr *net.UDPAddr, sAddr *net.UDPAddr) (*Client, error) {
	// create udp server
	s, err := udp_server.New(addr)
	if err != nil {
		return nil, err
	}

	bc, err := base_client.New(username, s)
	if err != nil {
		return nil, err
	}

	c := &Client{
		Client: bc,
		sAddr:  sAddr,
	}

	s.OnMessage(shared.CreateMessageCallback(c))

	return c, nil
}
