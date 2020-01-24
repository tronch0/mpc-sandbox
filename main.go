package main

import (
	"context"
	"fmt"
	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p-core/peerstore"
	"github.com/multiformats/go-multiaddr"
	"github.com/tronch0/mpc-sandbox/clihook"
	"log"
)

const (
	protocol = "/mpc-sandbox/1.0.0"
)

func main() {
	sourcePort, destAdd := clihook.RunCliHook()

	rawStringAddress := fmt.Sprintf("/ip4/0.0.0.0/tcp/%d", sourcePort)
	h, hostPeerCancelFn := createHost(rawStringAddress)
	defer hostPeerCancelFn()
	printNodeAddresses(h)
	if destAdd == "" {
		listenAndIdle(h)
	} else {
		connectAndListen(h, destAdd)
	}
}

func connectAndListen(h host.Host, destAdd string) {
	addressInfo := connect(h, destAdd)

	s, err := h.NewStream(context.Background(), addressInfo.ID, protocol)
	if err != nil {
		panic(err)
	}

	StreamHandler(s)

	select {}

}
func listenAndIdle(h host.Host) {
	h.SetStreamHandler(protocol, StreamHandler)
	printConnectionDetails(h)

	<-make(chan struct{})
}

func connect(h host.Host, destAdd string) *peer.AddrInfo {
	addressInfo := parseAddress(destAdd)

	h.Peerstore().AddAddrs(addressInfo.ID, addressInfo.Addrs, peerstore.PermanentAddrTTL)

	return addressInfo
}

func printNodeAddresses(h host.Host) {
	fmt.Println("This node's multiaddresses:")
	for _, la := range h.Addrs() {
		fmt.Printf(" - %v\n", la)
	}
	fmt.Println()
}

func parseAddress(destAdd string) *peer.AddrInfo {
	maddr, err := multiaddr.NewMultiaddr(destAdd)
	if err != nil {
		log.Fatalln(err)
	}

	info, err := peer.AddrInfoFromP2pAddr(maddr)
	if err != nil {
		log.Fatalln(err)
	}

	return info
}

func printConnectionDetails(h host.Host) {
	var port string

	for _, la := range h.Network().ListenAddresses() {
		if p, err := la.ValueForProtocol(multiaddr.P_TCP); err == nil {
			port = p
			break
		}
	}

	if port == "" {
		panic("error finding host local listening-port")
	}

	fmt.Printf("Run 'go run . -d /ip4/127.0.0.1/tcp/%v/p2p/%s' on another console. (same dir)\n", port, h.ID().Pretty())
	fmt.Printf("\nlistening...\n\n")
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
