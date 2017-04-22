package shared

import (
	"bufio"
	"crypto/cipher"
	"encoding/base64"
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"sync"
)

type Conn interface {
	Send(*Message) error
	Protocol() string
	GetAddr() net.Addr
	GetSecret() ([32]byte, error)
	SetSecret([32]byte)
}

type Client interface {
	WasKeySent() bool
	SetKeySent(bool)
	WasKeyReceived() bool
	SetKeyReceived(bool)
	GetServer() Server
	GetLog() *log.Logger
	GetSelf() *Peer
	GetPeer() *Peer
	SetPeer(*Peer)
	GetPeerConn() Conn
	SetPeerConn(Conn)
	GetServerConn() Conn
	SetServerConn(Conn)
	Connect()
	Stop()
	Start() error
	RegisteredCallback(Client)
	ConnectingCallback(Client)
	ConnectedCallback(Client)
	MessageCallback(Client, string)
	OnRegistered(func(Client))
	OnConnecting(func(Client))
	OnConnected(func(Client))
	OnMessage(func(Client, string))
}

type Server interface {
	Stop()
	Listen()
	CreateConn(net.Addr) (Conn, error)
	OnMessage(f func(Conns, Conn, *Message))
}

type UDPPayload struct {
	Bytes []byte
	Addr  *net.UDPAddr
}

type UDPConn struct {
	send   chan *UDPPayload
	addr   *net.UDPAddr
	secret string
}

func convertSecret(secretText string) ([32]byte, error) {
	// ensure secret has been set
	var secret [32]byte
	if secretText == "" {
		return secret, errors.New("secret has not been set")
	}

	// decode to byte slice
	bs, err := base64.StdEncoding.DecodeString(secretText)
	if err != nil {
		return secret, errors.New("could not decode secret")
	}

	// copy byte slice into byte array
	copy(secret[:], bs)
	return secret, nil
}

func (c *UDPConn) Send(m *Message) error {
	b, err := MessageOut(c, m)
	if err != nil {
		return err
	}

	c.send <- &UDPPayload{Bytes: b, Addr: c.addr}
	return err
}

func (c *UDPConn) Protocol() string {
	return "UDP"
}

func (c *UDPConn) GetAddr() net.Addr {
	return c.addr
}

func (c *UDPConn) GetSecret() ([32]byte, error) {
	return convertSecret(c.secret)
}

func (c *UDPConn) SetSecret(secret [32]byte) {
	c.secret = base64.StdEncoding.EncodeToString(secret[:])
}

func NewUDPConn(send chan *UDPPayload, addr *net.UDPAddr) *UDPConn {
	return &UDPConn{
		send: send,
		addr: addr,
	}
}

type TCPConn struct {
	C      *net.TCPConn
	secret string
}

func (c *TCPConn) Send(m *Message) error {
	b, err := MessageOut(c, m)
	if err != nil {
		return err
	}

	c.C.Write(b)
	return nil
}

func (c *TCPConn) Protocol() string {
	return "TCP"
}

func (c *TCPConn) GetAddr() net.Addr {
	return c.C.RemoteAddr()
}

func (c *TCPConn) GetSecret() ([32]byte, error) {
	return convertSecret(c.secret)
}

func (c *TCPConn) SetSecret(secret [32]byte) {
	c.secret = base64.StdEncoding.EncodeToString(secret[:])
}

func NewTCPConn(c *net.TCPConn) *TCPConn {
	return &TCPConn{C: c}
}

type Conns map[string]Conn

type Endpoint struct {
	IP   string `json:"ip"`
	Port int    `json:"port"`
}

func (e Endpoint) String() string {
	return e.IP + ":" + strconv.Itoa(e.Port)
}

type Blocks map[string]cipher.Block

type Registration struct {
	Username  string `json:"username"`
	PublicKey string `json:"publicKey"`
}

type History struct {
	w *bufio.Writer
	m *sync.Mutex
}

func (h *History) Add(text string) {
	h.m.Lock()
	defer h.m.Unlock()
	fmt.Fprintln(h.w, text)
	h.w.Flush()
}

func NewHistory(f *os.File) *History {
	return &History{
		w: bufio.NewWriter(f),
		m: &sync.Mutex{},
	}
}

type Message struct {
	Type    string      `json:"type"`
	PeerID  string      `json:"peerID,omitempty"`
	Error   string      `json:"error,omitempty"`
	Content interface{} `json:"data,omitempty"`
	Encrypt bool        `json:"-"`
	addr    *net.UDPAddr
}

func (m *Message) GetAddr() *net.UDPAddr {
	return m.addr
}

func (m *Message) SetAddr(addr *net.UDPAddr) *Message {
	m.addr = addr
	return m
}

type Peer struct {
	ID         string       `json:"id,omitempty"`
	Username   string       `json:"username,omitempty"`
	Endpoint   Endpoint     `json:"endpoint,omitempty"`
	PublicKey  string       `json:"publicKey,omitempty"`
	PrivateKey [32]byte     `json:"-"`
	Addr       *net.UDPAddr `json:"-"`
}

func (p *Peer) GetPublicKey() ([32]byte, error) {
	var key [32]byte
	bs, err := base64.StdEncoding.DecodeString(p.PublicKey)
	if err != nil {
		return key, err
	}
	copy(key[:], bs)
	return key, nil
}

func (p *Peer) SetPublicKey(key [32]byte) {
	p.PublicKey = base64.StdEncoding.EncodeToString(key[:])
}

type Peers map[string]*Peer
