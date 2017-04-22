package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/wilfreddenton/udp-hole-punching/shared"
)

func registeredCallback(c shared.Client) {
	fmt.Println("  (1) Connect with a Peer")
	fmt.Println("  (2) Wait for a peer to connect")
	fmt.Println("  (3) Exit")
	var n int
	for {
		fmt.Print("  > ")
		fmt.Scanln(&n)
		fmt.Print("\n")

		switch n {
		case 1:
			fmt.Println("  PeerID")
			fmt.Print("  > ")
			var id string
			for id == "" {
				fmt.Scanln(&id)
			}
			fmt.Print("\n")
			c.GetServerConn().Send(&shared.Message{
				Type:    "establish",
				PeerID:  c.GetSelf().ID,
				Content: id,
			})
			return
		case 2:
			fmt.Print("  waiting...\n\n")
			return
		case 3:
			fmt.Print("~ bye ~\n")
			os.Exit(0)
		default:
			continue
		}
	}
}

func connectingCallback(c shared.Client) {
	peer := c.GetPeer()
	pConn := c.GetPeerConn()
	fmt.Println("  connecting to peer...")
	fmt.Printf("    Username: %s\n", peer.Username)
	fmt.Printf("    ID: %s\n", peer.ID)
	fmt.Printf("    Address: %s\n\n", pConn.GetAddr())
}

func spacing(s1, s2 string) string {
	dif := len(s1) - len(s2)
	var spacing string
	if dif < 0 {
		for i := 0; i > dif; i -= 1 {
			spacing += " "
		}
	}
	return spacing
}

func createConnectedCallback(h *shared.History) func(c shared.Client) {
	return func(c shared.Client) {
		self := c.GetSelf()
		peer := c.GetPeer()

		fmt.Printf("  Connected to %s over an encrypted channel\n", peer.Username)
		// start chat process
		go func() {
			for {
				fmt.Printf("  %s > ", self.Username)
				r := bufio.NewReader(os.Stdin)
				bytes, _, _ := r.ReadLine()
				text := string(bytes)
				for text == "" {
					fmt.Println("  No empty messages allowed")
					fmt.Print("  > ")
					r = bufio.NewReader(os.Stdin)
					bytes, _, _ := r.ReadLine()
					text = string(bytes)
				}

				spacing := spacing(self.Username, peer.Username)
				h.Add(fmt.Sprintf("%s%s >         %s", self.Username, spacing, text))
				c.GetPeerConn().Send(&shared.Message{
					Type:    "message",
					PeerID:  self.ID,
					Content: text,
				})
			}
		}()
	}
}

func createMessageCallback(h *shared.History) func(c shared.Client, text string) {
	return func(c shared.Client, text string) {
		pUsername := c.GetPeer().Username
		spacing := spacing(pUsername, c.GetSelf().Username)
		h.Add(fmt.Sprintf("%s%s < %s", pUsername, spacing, text))
	}
}
