package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	re "github.com/wilfreddenton/reDo"
	"github.com/wilfreddenton/udp-hole-punching/shared"
	"github.com/wilfreddenton/udp-hole-punching/udp_client"
)

const (
	serverTCPPort = ":7001"
	serverUDPPort = ":9001"
)

var (
	serverTCPIP = "0.0.0.0"
	serverUDPIP = "127.0.0.1"
	serverIP    = flag.String("serverIP", "", "IP address of rendezvous server")
)

func main() {
	flag.Parse()

	if *serverIP != "" {
		serverTCPIP = *serverIP
		serverUDPIP = *serverIP
	}

	fmt.Print("\n  UDP Hole Punching v0.0.1 👊\n\n")

	log.SetFlags(log.LstdFlags | log.Lshortfile)
	var err error

	// get username from user
	var username string
	for username == "" || len(username) > 32 {
		fmt.Println("  Username (<= 32 chars)")
		fmt.Print("  > ")
		_, err = fmt.Scanln(&username)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Print("\n")
	}

	// create message history
	wd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	fName := fmt.Sprintf("%s/history-%s.txt", wd, username)
	hf, err := os.Create(fName)
	if err != nil {
		log.Fatal(err)
	}
	defer hf.Close()

	h := shared.NewHistory(hf)

	var c shared.Client
	var sAddr *net.UDPAddr
	var addr *net.UDPAddr
	sAddr, err = net.ResolveUDPAddr("udp", serverUDPIP+serverUDPPort)
	if err != nil {
		log.Fatal(err)
	}

	err = re.Do(5, func() error {
		addr, err = net.ResolveUDPAddr("udp", shared.GenPort())
		if err != nil {
			log.Fatal(err)
		}

		c, err = udp_client.New(username, addr, sAddr)
		return err
	})
	if err != nil {
		log.Fatal(err)
	}

	c.OnRegistered(registeredCallback)
	c.OnConnecting(connectingCallback)
	c.OnConnected(createConnectedCallback(h))
	c.OnMessage(createMessageCallback(h))

	fmt.Println("  ID")
	fmt.Printf("  > %s\n\n", c.GetSelf().ID)

	err = c.Start()
	if err != nil {
		log.Fatal(err)
	}

	exit := make(chan os.Signal)
	signal.Notify(exit, syscall.SIGINT, syscall.SIGTERM)
	log.Println(<-exit)

	c.Stop()
}
