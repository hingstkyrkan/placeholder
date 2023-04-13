/*
 * Simple "Hello World"-style Placeholder HTTP Server for testing
 */

package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
)

// Want to read it in the request handler
var listener net.Listener
var dumpHeaders bool
var dumpBody bool

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, I am %s\n", listener.Addr())

	log.Println(r.Method, r.URL.Path, r.RemoteAddr, r.UserAgent())

	if dumpHeaders {
		fmt.Fprintf(w, "\nDumping headers:\n\n")
		for header, values := range r.Header {
			for _, value := range values {
				fmt.Fprintf(w, "%s: %s\n", header, value)
			}
		}
	}

	if dumpBody {
		fmt.Fprintf(w, "\nDumping body:\n\n")
		body, _ := io.ReadAll(r.Body)
		fmt.Fprintf(w, "%s\n", body)
	}
}
func main() {
	// Listen on a random port by default, unless overridden
	listen := flag.String("listen", "tcp@:0", "tcp@[ip]:<port>, unix@<path>, systemd")
	wanttls := flag.Bool("tls", false, "enable tls on listener")
	wantheaders := flag.Bool("headers", false, "dump http headers in response")
	wantbody := flag.Bool("body", false, "dump http body in response")

	flag.Parse()
	listen_parts := strings.Split(*listen, "@")

	if *wantheaders {
		dumpHeaders = true
	}

	if *wantbody {
		dumpBody = true
	}

	var err error

	switch listen_parts[0] {
	case "tcp":
		fallthrough

	case "unix":
		listener, err = net.Listen(listen_parts[0], listen_parts[1])
		if err != nil {
			log.Fatal(err)
		}

	case "systemd":
		// From systemd/sd-daemon.h, first passed filedescriptor is 3
		const SD_LISTEN_FDS_START = 3

		if os.Getenv("LISTEN_FDS") == "1" {
			listener, err = net.FileListener(os.NewFile(SD_LISTEN_FDS_START, "systemd socket"))
		} else {
			log.Fatal("Couldn't find (exactly 1) filedescriptor passed")
		}

		if err != nil {
			log.Fatal(err)
		}

	default:
		log.Panic("Unsupported listener: ", listen_parts[0])
	}

	if *wanttls {
		cert, err := tls.LoadX509KeyPair("cert.pem", "key.pem")
		if err != nil {
			log.Fatal(err)
		}

		tlsconf := &tls.Config{Certificates: []tls.Certificate{cert}}

		// Replace listener with a TLS-wrapped listener
		listener = tls.NewListener(listener, tlsconf)
	}

	// With multiple listener options, some random, print where we ended up listening
	log.Println("Listening on:", listener.Addr())

	http.HandleFunc("/", hello)
	err = http.Serve(listener, nil)
	log.Fatal(err)
}
