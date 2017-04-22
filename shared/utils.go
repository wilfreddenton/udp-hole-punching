package shared

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"time"

	"github.com/wilfreddenton/crypto"
)

func init() {
	rand.Seed(time.Now().Unix())
}

func MessageIn(c Conn, b []byte) (*Message, error) {
	m := &Message{}
	err := json.Unmarshal(b, m)

	// if there is an error, check if message is encrypted, if so, decrypt and unmarshal
	if err != nil {
		var secret [32]byte
		secret, err = c.GetSecret()
		if err == nil {
			// decrypt
			b, err = crypto.Decrypt(b, secret)
			if err != nil {
				return m, err
			}

			// unmarshal into Request struct
			err = json.Unmarshal(b, m)
		}

		// if there was an error unmarshalling initially and either the message wasn't encrypted or unmarshaling the unencrypted message failed
		if err != nil {
			log.Print(err)
			return m, err
		}
	}

	return m, nil
}

func MessageOut(c Conn, m *Message) ([]byte, error) {
	b, err := json.Marshal(m)
	if err != nil {
		return b, err
	}

	if m.Encrypt {
		var s [32]byte
		s, err = c.GetSecret()
		if err != nil {
			return b, fmt.Errorf("cannot encrypt with an empty secret")
		}
		// encrypt message content
		b, err = crypto.Encrypt(b, s)
		if err != nil {
			return b, err
		}
	}

	return b, nil
}

func route(client Client, cs Conns, c Conn, m *Message) (*Message, error) {
	switch m.Type {
	case "greeting":
		return greetingHandler(client, c, m)
	case "register":
		return registerHandler(client, c, m)
	case "establish":
		return establishHandler(client, c, m)
	case "connect":
		return connectHandler(client, c, m)
	case "key":
		return keyHandler(client, c, m)
	case "message":
		return messageHandler(client, c, m)
	}
	return nil, nil
}

func CreateMessageCallback(client Client) func(Conns, Conn, *Message) {
	return func(cs Conns, c Conn, m *Message) {
		// ensure there was no error during registration
		res, err := route(client, cs, c, m)
		if err != nil {
			fmt.Println(err)
			client.GetLog().Fatal(err)
		}

		if res != nil {
			c.Send(res)
		}
	}
}

func GenPort() string {
	return ":" + strconv.Itoa(rand.Intn(65535-10000)+10000)
}
