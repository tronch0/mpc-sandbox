package clihook

import (
	"flag"
	"fmt"
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func RunCliHook() (int, string) {
	sPort := flag.Int("sp", -16351, "port to be listen on (should be between 1 and 65535")
	dest := flag.String("d", "", "Destination multiaddr string")

	flag.Parse()

	sourcePort := *sPort
	destAdd := *dest

	if isValidPort(sourcePort) == false {
		validNodePort := generateRandomNodePort()
		fmt.Printf("generated valid port to be use (new port value=%d)\n", validNodePort)
		sourcePort = validNodePort
	}

	fmt.Printf("Input-Parameters: \n "+
		"source-port: %d,\n"+
		" dest: %s\n", sourcePort, destAdd)

	return sourcePort, destAdd
}

func generateRandomNodePort() int {
	// need to be equal or lower than 65535
	var (
		min = 1
		max = 65535
	)

	return min + rand.Intn(max-min)
}

func isValidPort(port int) bool {
	if port == -16351 {
		fmt.Println("the input parameter SourcePort was not supplied")
		return false
	}

	if port > 65535 || port < 1 {
		fmt.Printf("the input parameter SourcePort was supplied with invalid value (sourcePort=%d)\n", port)
		return false
	}

	return true
}
