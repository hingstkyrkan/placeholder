/*
 * Simple "Hello World"-style Placeholder HTTP Server for testing
 */

package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"strings"

	"github.com/coreos/go-systemd/v22/activation"
)


// Want to read it in the request handler 
var listener net.Listener

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, I am %s\n", listener.Addr())
}

func main() {
	// Listen on a random port by default, unless overridden
	listen := flag.String("listen", "tcp@:0", "tcp@[ip]:<port>, unix@<path>, systemd")
	flag.Parse()
	listen_parts := strings.Split(*listen, "@")

	var err error

	switch listen_parts[0] {
	case "tcp", "udp":
		fallthrough

	case "unix":
		listener, err = net.Listen(listen_parts[0], listen_parts[1])
		if err != nil {
			log.Fatal(err)
		}

	case "systemd":
		listeners, err := activation.Listeners()
		if err != nil {
			log.Fatal(err)
		}

		listener = listeners[0]

	default:
		log.Panic("Unsupported listener: ", listen_parts[0])
	}

	// With multiple listener options, some random, print where we ended up listening
	log.Println("Listening on:", listener.Addr())

	http.HandleFunc("/", hello)
	err = http.Serve(listener, nil)
	log.Fatal(err)
}
