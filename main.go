package main

import (
	"context"
	"crypto/rand"
	"fmt"
	"github.com/libp2p/go-libp2p"
	crypto "github.com/libp2p/go-libp2p-core/crypto"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/network"
	"github.com/libp2p/go-libp2p-core/peerstore"
	ma "github.com/multiformats/go-multiaddr"
	"log"
)

const (
	firstPeerRawAddress  = "/ip4/0.0.0.0/tcp/12300" //  /ip4/0.0.0.0/tcp/%d "0.0.0.0:12300"
	secondPeerRawAddress = "/ip4/0.0.0.0/tcp/12351" //  /ip4/0.0.0.0/tcp/%d "0.0.0.0:12300"
	protocol             = "/echo/1.0.0"
)

func main() {

	// create first peer
	firstHost, firstPeerCancelFn := createHost(firstPeerRawAddress)
	defer firstPeerCancelFn()

	// create second peer
	secondHost, secondPeerCancelFn := createHost(secondPeerRawAddress)
	defer secondPeerCancelFn()

	// first peer: register handler for incoming traffic
	firstHost.SetStreamHandler(protocol, func(s network.Stream) {
		log.Println("first-peer: got new stream!!")
	})
	firstPeerFullAddr := getFullListingAddress(firstHost)
	log.Printf("first-peer: I am %s\n", firstPeerFullAddr)

	// second peer: register handler for incoming traffic
	secondHost.SetStreamHandler(protocol, func(s network.Stream) {
		log.Println("second-peer: got new stream!!")
	})
	secondPeerFullAddr := getFullListingAddress(secondHost)
	log.Printf("second-peer: I am %s\n", secondPeerFullAddr)

	connectHostAToHostB(firstHost, secondHost)
	connectHostAToHostB(secondHost, firstHost)

	streamToSecondHost, err := firstHost.NewStream(context.Background(), secondHost.ID(), protocol)
	if err != nil {
		panic(err)
	}

	streamToFirstHost, err := secondHost.NewStream(context.Background(), firstHost.ID(), protocol)
	if err != nil {
		panic(err)
	}

	_, err = streamToSecondHost.Write([]byte("hiiii its me!! (first send to second)"))
	if err != nil {
		panic(err)
	}
	streamToSecondHost.Close()

	_, err = streamToFirstHost.Write([]byte("hi its me!! (second send to first)"))
	if err != nil {
		panic(err)
	}

	streamToFirstHost.Close()
}
func connectHostAToHostB(hostA host.Host, hostB host.Host) {
	hostA.Peerstore().AddAddr(hostB.ID(), getFullListingAddress(hostB), peerstore.PermanentAddrTTL)
}
func createHost(listenAdd string) (host.Host, context.CancelFunc) {
	log.Print("creating lib2p2 host...")

	ctx, cancelFn := context.WithCancel(context.Background())

	h, err := libp2p.New(ctx,
		libp2p.Identity(generateIdentity()),
		libp2p.ListenAddrStrings(listenAdd),
		libp2p.DisableRelay(),
		libp2p.NoSecurity,
	)
	if err != nil {
		panic(err)
	}

	log.Printf("lib2p2 host created successfully (id: %s, address: %s)", h.ID(), h.Addrs()[0])

	return h, cancelFn
}

func getFullListingAddress(h host.Host) ma.Multiaddr {
	hostAddr, _ := ma.NewMultiaddr(fmt.Sprintf("/ipfs/%s", h.ID().Pretty()))
	addr := h.Addrs()[0]

	return addr.Encapsulate(hostAddr)
}

func generateIdentity() crypto.PrivKey {
	sk, _, err := crypto.GenerateEd25519Key(rand.Reader)
	if err != nil {
		panic(err)
	}

	return sk
}
