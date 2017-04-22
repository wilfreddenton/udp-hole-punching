package udp_server

import (
	"errors"
	"log"
	"net"
	"sync"
	"time"

	"github.com/wilfreddenton/udp-hole-punching/shared"
)

type Server struct {
	c               *net.UDPConn
	conns           shared.Conns
	send            chan *shared.UDPPayload
	messageCallback func(shared.Conns, shared.Conn, *shared.Message)
	exit            chan bool
	wg              *sync.WaitGroup
}

func (s *Server) sender() {
	s.wg.Add(1)
	defer s.wg.Done()

	for {
		select {
		case <-s.exit:
			log.Print("exiting UDP sender")
			return
		case p := <-s.send:
			_, err := s.c.WriteToUDP(p.Bytes, p.Addr)
			if err != nil {
				log.Print(err)
			}
		}
	}
}

func (s *Server) serve(b []byte, c shared.Conn) {
	defer s.wg.Done()
	m, err := shared.MessageIn(c, b)
	if err != nil {
		c.Send(&shared.Message{
			Error: "Malformed payload was sent",
		})
		return
	}

	go s.messageCallback(s.conns, c, m)
}

func (s *Server) receiver() {
	s.wg.Add(1)
	defer s.wg.Done()

	for {
		select {
		case <-s.exit:
			log.Print("exiting UDP receiver")
			s.c.Close()
			return
		default:
		}

		buf := make([]byte, 2048)
		s.c.SetDeadline(time.Now().Add(time.Second))
		n, addr, err := s.c.ReadFromUDP(buf)
		if err != nil {
			if opErr, ok := err.(*net.OpError); ok && opErr.Timeout() {
				continue
			}

			delete(s.conns, addr.String())
			log.Print(err)
			return
		}

		c, ok := s.conns[addr.String()]
		if !ok {
			c = shared.NewUDPConn(s.send, addr)
			s.conns[addr.String()] = c
		}

		// process message
		s.wg.Add(1)
		go s.serve(buf[:n], c)
	}
}

func (s *Server) CreateConn(addr net.Addr) (shared.Conn, error) {
	if addr == nil {
		return nil, errors.New("Conns addr must not be nil")
	}

	udpAddr, ok := addr.(*net.UDPAddr)
	if !ok {
		return nil, errors.New("could not assert net.Addr to *net.UDPAddr")
	}

	c := shared.NewUDPConn(s.send, udpAddr)
	s.conns[addr.String()] = c
	return c, nil
}

func (s *Server) OnMessage(f func(cs shared.Conns, c shared.Conn, m *shared.Message)) {
	s.messageCallback = f
}

func (s *Server) Stop() {
	close(s.exit)
	s.wg.Wait()
	log.Print("UDP server exited")
}

func (s *Server) Listen() {
	go s.sender()

	s.receiver()
}

func New(addr *net.UDPAddr) (*Server, error) {
	// create udp conn
	c, err := net.ListenUDP("udp", addr)
	if err != nil {
		return nil, err
	}

	return &Server{
		c:               c,
		conns:           make(shared.Conns),
		send:            make(chan *shared.UDPPayload, 100),
		messageCallback: func(cs shared.Conns, c shared.Conn, m *shared.Message) {},
		exit:            make(chan bool),
		wg:              &sync.WaitGroup{},
	}, nil
}
