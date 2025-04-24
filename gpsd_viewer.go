package main

import (
	"bufio"
	"embed"
	"encoding/json"
	"flag"
	"fmt"
	"io/fs"
	"log"
	"net"
	"net/http"
	"os"
	"sync"
	"time"
)

// GPSData represents the structure for a single GPS reading
type GPSData struct {
	Time    time.Time `json:"time"`
	Lat     float64   `json:"lat"`
	Lon     float64   `json:"lon"`
	Alt     float64   `json:"alt"`
	Speed   float64   `json:"speed"`
	RawJSON string    `json:"raw_json"`
}

var (
	dataMutex     sync.Mutex
	dataStore     []GPSData
	lastSkyData   json.RawMessage
	lastAISData   json.RawMessage
	gpsdAddress   string
	listenAddress string
)

// content holds our static web server content.
//
//go:embed web/*
var content embed.FS

func main() {
	// Configure logging
	log.SetOutput(os.Stdout)
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// CLI flag to specify gpsd address (default: localhost:2947)
	flag.StringVar(&gpsdAddress, "gpsd", "localhost:2947", "Address of gpsd server (e.g. localhost:2947)")
	flag.StringVar(&listenAddress, "p", "9000", "Address of web server (e.g. 9000)")
	flag.Parse()

	// Start GPSD connection handler in its own goroutine
	go connectToGPSD()

	// Set up HTTP routes
	http.HandleFunc("/json", gpsHandler) // JSON output route
	//http.Handle("/", http.StripPrefix("/", http.FileServer(http.Dir("./web")))) // Static file server
	contentStatic, _ := fs.Sub(content, "web")
	http.Handle("/", http.FileServer(http.FS(contentStatic)))

	log.Println("Web server running at port " + listenAddress)
	if err := http.ListenAndServe(":"+listenAddress, nil); err != nil {
		log.Fatalf("Web server error: %v", err)
	}
}

// connectToGPSD tries to connect to gpsd and keeps retrying on failure
func connectToGPSD() {
	for {
		log.Printf("Attempting to connect to gpsd at %s...", gpsdAddress)
		conn, err := net.Dial("tcp", gpsdAddress)
		if err != nil {
			log.Printf("Connection failed: %v. Retrying in 5 seconds...", err)
			time.Sleep(5 * time.Second)
			continue
		}
		log.Println("Connected to gpsd.")
		handleGPSD(conn)
		conn.Close()
		log.Println("Disconnected from gpsd. Retrying in 5 seconds...")
		time.Sleep(5 * time.Second)
	}
}

// handleGPSD reads and processes incoming data from gpsd
func handleGPSD(conn net.Conn) {
	fmt.Fprintf(conn, "?WATCH={\"enable\":true,\"TPV\":true,\"SKY\":true,\"AIS\":true,\"json\":true}\n")
	scanner := bufio.NewScanner(conn)

	for scanner.Scan() {
		line := scanner.Text()
		now := time.Now()

		var parsed map[string]interface{}
		if err := json.Unmarshal([]byte(line), &parsed); err != nil {
			log.Printf("JSON parse error: %v, line: %s", err, line)
			continue
		}

		class, _ := parsed["class"].(string)

		switch class {
		case "TPV":
			// Only process data with mode 3 (3D fix)
			mode, _ := parsed["mode"].(float64)
			if mode != 3 {
				continue
			}

			lat, _ := parsed["lat"].(float64)
			lon, _ := parsed["lon"].(float64)
			alt, _ := parsed["alt"].(float64)
			speed, _ := parsed["speed"].(float64)

			entry := GPSData{
				Time:    now,
				Lat:     lat,
				Lon:     lon,
				Alt:     alt,
				Speed:   speed,
				RawJSON: line,
			}

			dataMutex.Lock()
			dataStore = append(dataStore, entry)
			if len(dataStore) > 10 {
				dataStore = dataStore[len(dataStore)-10:]
			}
			dataMutex.Unlock()

		case "SKY":
			dataMutex.Lock()
			lastSkyData = json.RawMessage(line)
			dataMutex.Unlock()

		case "AIS":
			dataMutex.Lock()
			lastAISData = json.RawMessage(line)
			dataMutex.Unlock()
		}
	}

	if err := scanner.Err(); err != nil {
		log.Printf("Scanner error: %v", err)
	}
}

// gpsHandler handles requests to /json and returns the current GPS data
func gpsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	dataMutex.Lock()
	defer dataMutex.Unlock()

	response := map[string]interface{}{
		"tpv": dataStore,
		"sky": json.RawMessage(`null`),
		"ais": json.RawMessage(`null`),
	}

	if lastSkyData != nil {
		response["sky"] = lastSkyData
		response["ais"] = lastAISData
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, fmt.Sprintf("Error encoding JSON response: %v", err), http.StatusInternalServerError)
		log.Printf("HTTP JSON error: %v", err)
	}
}
