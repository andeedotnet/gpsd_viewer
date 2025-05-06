package main

import (
	"embed"
	"flag"
	api "gpsd_viewer/handlers"
	gpsd "gpsd_viewer/models"
	"io/fs"
	"log"
	"net/http"
	"os"
)

var (
	gpsdAddress       string
	listenAddress     string
	useEmbeddedWebapp string
)

// content holds our static web server content.
//
//go:embed static/*
var content embed.FS

func main() {
	// Configure logging
	log.SetOutput(os.Stdout)
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// Read CLI flags
	flag.StringVar(&gpsdAddress, "gpsd", "localhost:2947", "Address of gpsd server (e.g. localhost:2947)")
	flag.StringVar(&listenAddress, "p", "9000", "Address of web server (e.g. 9000)")
	flag.StringVar(&useEmbeddedWebapp, "e", "true", "Use embedded webapp")
	flag.Parse()

	// Start GPSD client in its own goroutine
	go gpsd.ConnectToGPSD(gpsdAddress)

	// Set up HTTP routes
	http.HandleFunc("/api", api.APIHandler)

	if useEmbeddedWebapp == "true" {
		contentStatic, _ := fs.Sub(content, "static")
		http.Handle("/", http.FileServer(http.FS(contentStatic)))
	} else {
		http.Handle("/", http.StripPrefix("/", http.FileServer(http.Dir("./static")))) // Static file server
	}

	log.Println("Web server running at port " + listenAddress)
	if err := http.ListenAndServe(":"+listenAddress, nil); err != nil {
		log.Fatalf("Web server error: %v", err)
	}
}
