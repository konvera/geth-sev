package main

import (
	"crypto/tls"
	"io"
	"log"
	"net/http"

	"github.com/konvera/geth-sev/constellation/atls"
	"github.com/konvera/geth-sev/constellation/attestation/azure/snp"
)

func helloHandler(w http.ResponseWriter, r *http.Request) {
	// Write "Hello, world!" to the response body
	io.WriteString(w, "Hello, world!\n")
}

func main() {
	// Set up a /hello resource handler
	http.HandleFunc("/hello", helloHandler)

	issuer := snp.NewIssuer(nil)
	tlsConfig, err := atls.CreateAttestationServerTLSConfig(issuer, nil)
	if (err != nil) {
		log.Fatal(err)
	}

	listener, err := tls.Listen("tcp", ":8443", tlsConfig)
	if (err != nil) {
		log.Fatal(err)
	}

	// Listen to HTTPS connections with the server certificate and wait
	log.Fatal(http.Serve(listener, nil))
}
