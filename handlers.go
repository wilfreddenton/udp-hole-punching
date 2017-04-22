package main

import (
	"encoding/base64"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/mitchellh/mapstructure"
	"github.com/wilfreddenton/crypto"
	"github.com/wilfreddenton/udp-hole-punching/shared"
)

func greetingHandler(conn shared.Conn, m *shared.Message) (*shared.Message, error) {
	// ensure that public key was sent in greeting request
	str, ok := m.Content.(string)
	if !ok {
		return nil, errors.New("greeting request must contain client's public key")
	}

	// get public key contained in content
	bs, err := base64.StdEncoding.DecodeString(str)
	if err != nil {
		return nil, err
	}

	// create shared secret from private key and peer public key
	var clientPubKey [32]byte
	copy(clientPubKey[:], bs[:])
	conn.SetSecret(crypto.GenSharedSecret(priKey, clientPubKey))

	// send greeting response
	return &shared.Message{
		Type:    "greeting",
		Content: base64.StdEncoding.EncodeToString(pubKey[:]),
	}, nil
}

// register the requesting peer in the server
func registerHandler(peers shared.Peers, c shared.Conn, m *shared.Message) (*shared.Message, error) {
	// map -> structure the content
	var registration shared.Registration
	err := mapstructure.Decode(m.Content, &registration)
	if err != nil {
		return nil, err
	}

	// register peer
	endpoint := strings.Split(c.GetAddr().String(), ":")
	if len(endpoint) != 2 {
		return nil, errors.New("address is not valid")
	}

	port, err := strconv.Atoi(endpoint[1])
	if err != nil {
		return nil, err
	}

	peers[m.PeerID] = &shared.Peer{
		ID:       m.PeerID,
		Username: registration.Username,
		Endpoint: shared.Endpoint{
			IP:   endpoint[0],
			Port: port,
		},
	}
	log.Printf("Registered peer: %s at addr %s", m.PeerID, c.GetAddr().String())

	// confirm registry to peer
	return &shared.Message{
		Type:    "register",
		Encrypt: true,
	}, nil
}

// facilitate in the establishing of the p2p connection
func establishHandler(peers shared.Peers, conns shared.Conns, m *shared.Message) (*shared.Message, error) {
	// make sure requesting peer has registered with server
	rp, ok := peers[m.PeerID]
	if !ok {
		return nil, errors.New("client is not registered with this server")
	}

	// make sure that a valid payload was sent
	id, ok := m.Content.(string)
	if !ok {
		return nil, errors.New("request content is malformed")
	}

	// make sure the other peer has registered with the server
	op, ok := peers[id]
	if !ok {
		return nil, fmt.Errorf("The peer: %s has not registered with the server.", id)
	}

	// get conn for other peer
	conn, ok := conns[op.Endpoint.String()]
	if !ok {
		return nil, fmt.Errorf("Could not resolve the peer: %s's conn", id)
	}

	// send requesting peer's endpoint to other peer
	conn.Send(&shared.Message{
		Type:    "establish",
		Content: rp,
		Encrypt: true,
	})

	// send requesting peer other peer's endpoint
	return &shared.Message{
		Type:    "establish",
		Content: op,
		Encrypt: true,
	}, nil
}

func notFoundHandler(m *shared.Message) (*shared.Message, error) {
	return nil, fmt.Errorf("Request type %s undefined", m.Type)
}
