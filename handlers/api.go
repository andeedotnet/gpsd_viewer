package handlers

import (
	"encoding/json"
	"fmt"
	gpsd "gpsd_viewer/models"
	"log"
	"net/http"
)

// gpsHandler handles requests to /json and returns the current GPS data
func APIHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	gpsd.DataMutex.Lock()
	defer gpsd.DataMutex.Unlock()

	response := map[string]interface{}{
		"tpv": json.RawMessage(`null`),
		"sky": json.RawMessage(`null`),
		"ais": json.RawMessage(`null`),
	}

	if gpsd.LastGPSData != nil {
		response["tpv"] = gpsd.LastGPSData
	}
	if gpsd.LastSkyData != nil {
		response["sky"] = gpsd.LastSkyData
	}

	if gpsd.LastAISData != nil {
		response["ais"] = gpsd.LastAISData
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, fmt.Sprintf("Error encoding JSON response: %v", err), http.StatusInternalServerError)
		log.Printf("HTTP JSON error: %v", err)
	}
}
