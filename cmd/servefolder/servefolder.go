/*

This app serves a folder via HTTP.
Useful for quick sharing. Not suitable for public hosting over the internet.

*/
package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
)

const (
	appName    = "servefolder"
	appVersion = "v0.1.0"
	appAuthor  = "Andras Belicza"
	appHome    = "https://github.com/icza/toolbox"
)

var (
	version = flag.Bool("version", false, "print version info and exit")
	addr    = flag.String("addr", ":8080", "address to start the server on")
)

func main() {
	flag.Usage = printUsage
	flag.Parse()

	if *version {
		printVersion()
		return
	}

	args := flag.Args()

	path := ""
	if len(args) > 0 {
		path = args[0]
	}

	path, err := filepath.Abs(path)
	if err != nil {
		log.Printf("Failed to resolve %s", path)
		os.Exit(1)
	}

	log.Printf("Serving folder: %s", path)

	// Find out and print which addresses we're listening on:
	host, port, err := net.SplitHostPort(*addr)
	if err != nil {
		log.Printf("Failed to split addr: %v", err)
		os.Exit(2)
	}
	if host != "" {
		// Host is explicit:
		log.Printf("Listening on http://%s/", *addr)
	} else {
		// Host is missing, we'll listen on all available interfaces:
		printLocalInterfaces(port)
	}

	http.Handle("/", http.FileServer(http.Dir(path)))
	log.Print(http.ListenAndServe(*addr, nil))
}

func printLocalInterfaces(port string) {
	ifaces, err := net.Interfaces()
	if err != nil {
		log.Printf("Failed to get interfaces: %v", err)
		os.Exit(11)
	}

	for _, iface := range ifaces {
		addrs, err := iface.Addrs()
		if err != nil {
			log.Printf("Failed to get interface addresses: %v", err)
			os.Exit(12)
		}

		for _, addr := range addrs {
			if ipNet, ok := addr.(*net.IPNet); ok {
				if ipv4 := ipNet.IP.To4(); ipv4 != nil {
					log.Printf("Listening on http://%s:%s/", ipv4, port)
				}
			}
		}
	}
}

func printVersion() {
	fmt.Println(appName, "version:", appVersion)
	fmt.Println("Platform:", runtime.GOOS, runtime.GOARCH)
	fmt.Println("Built with:", runtime.Version())
	fmt.Println("Author:", appAuthor)
	fmt.Println("Home page:", appHome)
}

func printUsage() {
	fmt.Println("Usage:")
	name := os.Args[0]
	fmt.Printf("%s [FLAGS] [folder-to-serve]\n", name)
	fmt.Println("(The current working directory is served if not specified.)")
	fmt.Println()
	fmt.Println("Flags:")

	flag.PrintDefaults()
}
