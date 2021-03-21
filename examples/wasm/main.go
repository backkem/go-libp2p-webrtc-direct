package main

import (
	"context"
	"fmt"
	"io/ioutil"

	smux "github.com/libp2p/go-libp2p-core/mux"
	tpt "github.com/libp2p/go-libp2p-core/transport"
	mplex "github.com/libp2p/go-libp2p-mplex"
	direct "github.com/libp2p/go-libp2p-webrtc-direct"
	ma "github.com/multiformats/go-multiaddr"
	"github.com/pion/webrtc/v3"
)

func main() {
	maddr, err := ma.NewMultiaddr("/ip4/127.0.0.1/tcp/9090/http/p2p-webrtc-direct")
	check(err)

	transport := direct.NewTransport(
		webrtc.Configuration{},
		new(mplex.Transport),
	)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	c, err := transport.Dial(ctx, maddr, "peerA")
	check(err)
	defer c.Close()
	fmt.Println("[dialer] Opened connection")

	s, err := c.OpenStream(context.Background())
	check(err)
	fmt.Println("[dialer] Opened stream")

	_, err = s.Write([]byte("hey, how is it going. I am the dialer"))
	check(err)

	err = s.Close()
	check(err)
}

func handleConn(c tpt.CapableConn) {
	for {
		s, err := c.AcceptStream()
		if err != nil {
			return
		}

		fmt.Println("[listener] Got stream")
		go handleStream(s)
	}
}
func handleStream(s smux.MuxedStream) {
	b, err := ioutil.ReadAll(s)
	check(err)
	fmt.Println("[listener] Received:")
	fmt.Println(string(b))
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}
