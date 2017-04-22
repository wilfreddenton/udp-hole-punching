# 👊 UDP Hole Punching 👊

Yusuke punches a hole right through your NAT

![spirit punch](http://i.imgur.com/ZwBpD0a.jpg)

**Disclaimer**: This is not a production ready chat application. While it does create AES encrypted connections `client <-> client` and `client <-> server`, this code has not been audited or tested by any security specialists. This is simply an exercise for me to learn more about P2P networking and technologies as well as provide some examples of the technologies in use for others who are interested in learning. Additionally, the udp client does not implement a protocol that ensures the successful delivery of messages and so some will be lost over spotty connections.

## Install

```
go get github.com/wilfreddenton/udp-hole-punching
```

## Usage

### 1. Setup rendezvous server

The main package is the rendezvous server. Find a VPS or something to host it on. You can run everything locally but it won't really be testing whether or not hole punching works because it's on the same machine. Make sure that the server has TCP and UDP ports open to incoming traffic from 0-65535.

### 2. Adjust UI settings

There are two UIs that you can use `gui` which is a web UI and `term-ui` which is a terminal UI.

Before you use one you should open the `main.go` file and switch the `serverTCPIP` and `serverUDPIP` constants to the IP address of your rendezvous server (no port).

To run the web UI

1. `cd gui/ui`
2. `npm install`
3. `npm run build`
4. `cd ..`
5. `go build`
6. `gui` or `./gui`
7. point your browser to `localhost:8000`

To disconnect and start a new chat simply refresh.

To run the terminal UI

1. `cd term-ui`
2. `go build`
3. `term-ui` or `./term-ui`

To disconnect and start a new chat `ctrl-c` to exit the program and run it again.

### 3. Find a friend

If not a friend then get access to a computer behind a different router and set up a client on there.

### 4. Test it out

Run the clients and provide the PeerID of one client to the other client and if the network topology permits hole punching then you will establish an encrypted connection between the clients.

## Architecture

![spirit punch architecture](http://i.imgur.com/dZNEhpw.png)

1. Both clients register themselves using their ID with the rendezvous server
2. Client A makes an "establish" request to the rendezvous server sending the `ID` of the peer it would like to being communicating with
3. Upon receiving the "establish" request from client A and verifying that both client A and the requested peer, client B, have registered, the server sends an "establish" response back to client A as well as client B informing the peers of each other's information.
4. The peers can now send requests directly to each other with the information they've received from the rendezvous server. They create this connection using the hole-punching algorithm described in reference 1.

## Simplification of the algorithm

To make the implementation of hole punching a little simpler, the clients to not attempt to connect to each other's private IP addresses. Clients that are behind the same NAT will still be able to connect but they will do so with public IP addresses and not private ones. The routers I have tested seemed to understand that the peers were on it's local network and facilitated the connection without going to the outside internet.

## References

1. [Peer-to-Peer Communication Across Network Address Translators](https://www.usenix.org/legacy/event/usenix05/tech/general/full_papers/ford/ford.pdf)
2. [Diffie-Hellman Key Exchange: A Non-mathematician’s explanation](http://academic.regis.edu/cias/ia/palmgren_-_diffie-hellman_key_exchange.pdf)
3. powered by [Curve25519](https://cr.yp.to/ecdh.html)
