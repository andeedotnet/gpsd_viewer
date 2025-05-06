package gpsd

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"sync"
	"time"
)

var (
	DataMutex   sync.Mutex
	LastGPSData json.RawMessage
	LastSkyData json.RawMessage
	LastAISData json.RawMessage
)

// connectToGPSD tries to connect to gpsd and keeps retrying on failure
func ConnectToGPSD(gpsdAddress string) {
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
			DataMutex.Lock()
			LastGPSData = json.RawMessage(line)
			DataMutex.Unlock()

		case "SKY":
			DataMutex.Lock()
			LastSkyData = json.RawMessage(line)
			DataMutex.Unlock()

		case "AIS":
			DataMutex.Lock()
			LastAISData = json.RawMessage(line)
			DataMutex.Unlock()
		}
	}

	if err := scanner.Err(); err != nil {
		log.Printf("Scanner error: %v", err)
	}
}
