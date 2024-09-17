package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
)

var (
	version   string
	buildTime string
)

func echoHandler(w http.ResponseWriter, r *http.Request) {
	// Print method and URL
	fmt.Fprintf(w, "%s %s %s\n", r.Method, r.URL.Path, r.Proto)

	// Print headers
	for name, values := range r.Header {
		for _, value := range values {
			fmt.Fprintf(w, "%s: %s\n", name, value)
		}
	}

	// Print remote ip address
	fmt.Fprintf(w, "x-goparrot-RemoteAddr: %s\n", r.RemoteAddr)

	// Print my ip address
	fmt.Fprintf(w, "x-goparrot-LocalAddr: %s\n", r.Context().Value(http.LocalAddrContextKey))

	// Print body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Unable to read body", http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, "\n%s", body)
}

func main() {
	programName := os.Args[0]

	helpFlag := flag.Bool("help", false, "print help")
	portFlag := flag.Int("port", 8080, "port to listen on")
	keyFileFlag := flag.String("key", "", "tls key file")
	certFileFlag := flag.String("cert", "", "tls certificate file")
	versionFlag := flag.Bool("version", false, "print version information")
	flag.Parse()

	if *helpFlag {
		fmt.Printf("Usage: %s [options]\n", programName)
		flag.PrintDefaults()
		return
	}

	if *versionFlag {
		fmt.Printf("%s version %s, built at %s\n", programName, version, buildTime)
		return
	}

	http.HandleFunc("/", echoHandler)

	address := fmt.Sprintf(":%d", *portFlag)
	if *certFileFlag != "" && *keyFileFlag != "" {
		fmt.Printf("%s: Listening on port %d with TLS keyfile=%s certfile=%s\n", programName, *portFlag, *certFileFlag, *keyFileFlag)
		if err := http.ListenAndServeTLS(address, *certFileFlag, *keyFileFlag, nil); err != nil {
			fmt.Printf("%s: Error starting server with TLS: %s\n", programName, err)
		}
	} else {
		fmt.Printf("%s: Listening on port %d\n", programName, *portFlag)
		if err := http.ListenAndServe(address, nil); err != nil {
			fmt.Printf("%s: Error starting server: %s\n", programName, err)
		}
	}
}
