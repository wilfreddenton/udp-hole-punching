package shared

import (
	"encoding/base64"
	"errors"
	"fmt"
	"net"

	"github.com/mitchellh/mapstructure"
	"github.com/wilfreddenton/crypto"
)

func greetingHandler(c Client, serverConn Conn, m *Message) (*Message, error) {
	l := c.GetLog()
	self := c.GetSelf()
	// quit the client if greeting fails
	if m.Error != "" {
		l.Fatal(m.Error)
		return nil, errors.New(m.Error)
	}

	// ensure that server sent back a public key string
	s, ok := m.Content.(string)
	if !ok {
		return nil, errors.New("expected to receive public key with greeting")
	}

	// get server public key
	bs, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return nil, err
	}
	var pubKey [32]byte
	copy(pubKey[:], bs)

	// get self public keySent
	sPubKey, err := self.GetPublicKey()
	if err != nil {
		return nil, err
	}

	// create and store secret
	serverConn.SetSecret(crypto.GenSharedSecret(self.PrivateKey, pubKey))

	// send register message to server
	return &Message{
		Type:   "register",
		PeerID: self.ID,
		Content: Registration{
			Username:  self.Username,
			PublicKey: base64.StdEncoding.EncodeToString(sPubKey[:]),
		},
	}, nil
}

func registerHandler(c Client, serverConn Conn, m *Message) (*Message, error) {
	// quit the client if registration fails
	if m.Error != "" {
		return nil, errors.New(m.Error)
	}

	c.RegisteredCallback(c)
	return nil, nil
}

func establishHandler(c Client, serverConn Conn, m *Message) (*Message, error) {
	l := c.GetLog()
	l.Print("establish request from server")
	if m.Error != "" {
		return nil, errors.New(m.Error)
	}

	var p Peer
	err := mapstructure.Decode(m.Content, &p)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	c.SetPeer(&Peer{
		ID:       p.ID,
		Username: p.Username,
	})

	var addr net.Addr
	switch serverConn.Protocol() {
	case "UDP":
		addr, err = net.ResolveUDPAddr("udp", p.Endpoint.String())
	case "TCP":
		addr, err = net.ResolveTCPAddr("tcp", p.Endpoint.String())
	default:
		addr, err = nil, fmt.Errorf("unknown Conn protocol %s", serverConn.Protocol())
	}
	if err != nil {
		return nil, err
	}

	if c.GetPeerConn() != nil && c.GetPeerConn().GetAddr().String() != addr.String() {
		l.Print("ignoring establish request because the client is already connected to a peer")
		return nil, nil
	}

	go func() {
		pConn, err := c.GetServer().CreateConn(addr)
		if err != nil {
			return
		}

		c.SetPeerConn(pConn)

		go c.Connect()

		c.ConnectingCallback(c)
	}()
	return nil, nil
}

func connectHandler(c Client, peerConn Conn, m *Message) (*Message, error) {
	self := c.GetSelf()
	l := c.GetLog()

	pConn := c.GetPeerConn()
	if pConn == nil {
		return nil, nil
	}

	if pConn != peerConn {
		// if addresses are the same then this is the correct peer but the listener has picked up the message and created a new conn
		if pConn.GetAddr().String() == peerConn.GetAddr().String() {
			pConn = peerConn
			c.SetPeerConn(pConn)
		}
		return nil, errors.New("received connect message from unknown peer")
	}

	l.Printf("connection mirror request from peer %s at %s, sending mirror...", self.Username, pConn.GetAddr())

	pubKey, err := self.GetPublicKey()
	if err != nil {
		return nil, err
	}

	// confirm public key was sent to peer
	defer c.SetKeySent(true)

	return &Message{
		Type:    "key",
		PeerID:  self.ID,
		Content: base64.StdEncoding.EncodeToString(pubKey[:]),
	}, nil
}

func keyHandler(c Client, peerConn Conn, m *Message) (*Message, error) {
	l := c.GetLog()
	pConn := c.GetPeerConn()
	if pConn != peerConn {
		return nil, errors.New("received key message from unknown peer")
	}

	// ensure that public key was sent with message
	s, ok := m.Content.(string)
	if !ok {
		return nil, errors.New("no public key was sent with key message")
	}

	// decode and store the sent public key
	bs, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return nil, err
	}

	var pubKey [32]byte
	copy(pubKey[:], bs)

	// create and store cipher with other peer's public key
	pConn.SetSecret(crypto.GenSharedSecret(c.GetSelf().PrivateKey, pubKey))

	// confirm peer's public key was received
	c.SetKeyReceived(true)

	l.Printf("received communication mirror from peer %s at %s", c.GetPeer().Username, pConn.GetAddr())
	return nil, nil
}

func messageHandler(c Client, peerConn Conn, m *Message) (*Message, error) {
	pConn := c.GetPeerConn()
	if pConn != peerConn {
		return nil, errors.New("received message message from unknown peer")
	}
	text, ok := m.Content.(string)
	if !ok {
		return nil, errors.New("message message must send some text in content field")
	}

	c.MessageCallback(c, text)
	// c.messageHook(c, s)
	return nil, nil
}
