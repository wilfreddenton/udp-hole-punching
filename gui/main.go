package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/googollee/go-socket.io"
	"github.com/rs/cors"
	"github.com/wilfreddenton/udp-hole-punching/shared"
)

type user struct {
	Username string `json:"username"`
	Protocol string `json:"protocol"`
}

const (
	serverUDPPort = ":9001"
	serverTCPPort = ":7001"
)

var (
	serverTCPIP = "0.0.0.0"
	serverUDPIP = "127.0.0.1"
	useCors     = flag.Bool("cors", false, "Use CORS or not")
	serverIP    = flag.String("serverIP", "", "IP address of rendezvous server")
	STATE       = &state{}
)

type state struct {
	username string
	protocol string
	id       string
	client   shared.Client
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./ui/dist/")
}

func main() {
	flag.Parse()

	if *serverIP != "" {
		serverTCPIP = *serverIP
		serverUDPIP = *serverIP
	}

	server, err := socketio.NewServer(nil)
	if err != nil {
		log.Fatal(err)
	}

	server.On("connection", func(so socketio.Socket) {
		fmt.Println("connection")
		// when user submits username
		so.On("enter", func(msg string) {
			u := &user{}
			err := json.Unmarshal([]byte(msg), u)
			if err != nil {
				log.Print(err)
				so.Emit("error", err.Error())
			}

			enterHandler(so, STATE, u)
		})
		// when user enters peerID
		so.On("establish", func(peerID string) {
			fmt.Println("establish")
			STATE.client.GetServerConn().Send(&shared.Message{
				Type:    "establish",
				PeerID:  STATE.id,
				Content: peerID,
			})
		})
		// when user sends a message
		so.On("message", func(text string) {
			fmt.Println("message:", text)
			STATE.client.GetPeerConn().Send(&shared.Message{
				Type:    "message",
				PeerID:  STATE.id,
				Content: text,
			})
		})
		// when user resets chat
		so.On("reset", func(text string) {
			fmt.Println("reset")
			STATE.client.Stop()
		})
		// socket.io disconnection event
		so.On("disconnection", func() {
			fmt.Println("disconnection")
			if STATE.client != nil {
				STATE.client.Stop()
			}
		})
	})
	server.On("error", func(so socketio.Socket, err error) {
		log.Println("error:", err)
	})

	mux := http.NewServeMux()
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./ui/dist/static"))))
	mux.Handle("/socket.io/", server)
	mux.HandleFunc("/", indexHandler)

	handler := http.Handler(mux)
	if *useCors {
		fmt.Println("CORS enabled")
		handler = cors.New(cors.Options{
			AllowedOrigins:   []string{"*"},
			AllowCredentials: true,
		}).Handler(mux)
	}

	go http.ListenAndServe(":8000", handler)

	exit := make(chan os.Signal)
	signal.Notify(exit, syscall.SIGINT, syscall.SIGTERM)
	log.Println(<-exit)

	if STATE.client != nil {
		STATE.client.Stop()
	}
}
