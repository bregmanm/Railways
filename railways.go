package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"github.com/gorilla/mux"
	"github.com/deckarep/golang-set"
)

// Global stations list
var stations = make(map[string]mapset.Set)
// Global trains list
var trains = mapset.NewSet()

// JSON structures
// Railway station
type JsonStation struct {
	Name string `json: "name"`
}
// Error response body
type ErrorResponse struct {
	Error string `json: "error"`
}
// Train
type JsonTrain struct {
	Name string `json: "name"`
}
// Stations
type JsonStations struct {
	Stations []string `json: "stations"`
}
// Trains
type JsonTrains struct {
	Trains []string `json: "trains"`
}
// TripAction
type JsonTripAction struct {
	FromStation string `json: "FromStation"`
	ToStation string `json: "ToStation"`
	Train     string `json: "Train"`
}

// Send error response to client
func SendErrorResponse(w http.ResponseWriter, status int, explanation string, params ...string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	var ErrorBody =	ErrorResponse{Error: fmt.Sprintf(explanation, params)}
	json.NewEncoder(w).Encode(ErrorBody)
}
// Get station name from URL
func GetStationNameFromUrl(r *http.Request) string {
	params := mux.Vars(r)
	return params["station"]
}
// Add the new station
func AddStation(w http.ResponseWriter, r *http.Request) {
	var station JsonStation
	_ = json.NewDecoder(r.Body).Decode(&station)
	var stationName = station.Name
	if _, ok := stations[stationName]; ok {
		// Station is already exists
		SendErrorResponse(w, http.StatusConflict, "Station %s is already exist", stationName)
	} else {
		stations[stationName] = mapset.NewSet()
	}
}
// Delete the station
func DeleteStation(w http.ResponseWriter, r *http.Request) {
	var station JsonStation
	_ = json.NewDecoder(r.Body).Decode(&station)
	var stationName = station.Name
	if _, ok := stations[stationName]; ok {
		delete(stations, stationName)
	} else {
		// Station is not exists
		SendErrorResponse(w, http.StatusNotFound, "Station %s does not exist", stationName)
	}
}
// Get stations
func GetStations(w http.ResponseWriter, r *http.Request) {
	stationsSlice := make([]string, 0, len(stations))
	for k := range stations {
		stationsSlice = append(stationsSlice, k)
	}
	var ListStations = JsonStations{Stations: stationsSlice}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ListStations)
}
// Add the new train to station
func AddTrain(w http.ResponseWriter, r *http.Request) {
	// Get station name from URL
	var stationName = GetStationNameFromUrl(r)
	// Check is the station name valid
	if stationTrains, ok := stations[stationName]; ok {
		// Get the train name
		var train JsonTrain
		_ = json.NewDecoder(r.Body).Decode(&train)
		var trainName = train.Name
		if trains.Contains(trainName) {
			// Train is already exists
			SendErrorResponse(w, http.StatusConflict, "Train %s is already exist", trainName)
		} else {
			// Add the new train to global list
			trains.Add(trainName)
			// and to given station
			stationTrains.Add(trainName)
		}
	} else {
		// Wrong station name
		SendErrorResponse(w, http.StatusNotFound, "Station %s does not exist", stationName)
	}
}
// Delete the train from station
func DeleteTrain(w http.ResponseWriter, r *http.Request) {
	// Get station name from URL
	var stationName = GetStationNameFromUrl(r)
	// Check is the station name valid
	if stationTrains, ok := stations[stationName]; ok {
		// Get the train name
		var train JsonTrain
		_ = json.NewDecoder(r.Body).Decode(&train)
		var trainName = train.Name
		if trains.Contains(trainName) {
			trains.Remove(trainName)
			stationTrains.Remove(trainName)
		} else {
			// Train does not exists
			SendErrorResponse(w, http.StatusNotFound, "Train %s does not exist", trainName)
		}
	} else {
		// Wrong station name
		SendErrorResponse(w, http.StatusNotFound, "Station %s does not exist", stationName)
	}
}
// Get trains placed on station
func GetTrains(w http.ResponseWriter, r *http.Request) {
	// Get station name from URL
	var stationName = GetStationNameFromUrl(r)
	// Check is the station name valid
	if stationTrains, ok := stations[stationName]; ok {
		trainsSlice := make([]string, 0, stationTrains.Cardinality())
		it := stationTrains.Iterator()
		for k := range it.C {
			trainsSlice = append(trainsSlice, k.(string))
		}
		var ListTrains = JsonTrains{Trains: trainsSlice}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(ListTrains)
	} else {
		// Wrong station name
		SendErrorResponse(w, http.StatusNotFound, "Station %s does not exist", stationName)
	}
}
// Move train from one station to other
func Trip(w http.ResponseWriter, r *http.Request) {
	var tripAction JsonTripAction
	_ = json.NewDecoder(r.Body).Decode(&tripAction)
	var stationNameFrom = tripAction.FromStation
	if trainsFrom, ok := stations[stationNameFrom]; ok {
		// Station of departure exists
		var trainName = tripAction.Train
		if trainsFrom.Contains(trainName) {
			// Train stays on the departure station
			var stationNameTo = tripAction.ToStation
			if trainsTo, ok := stations[stationNameTo]; ok {
				// Destination station exists
				if stationNameFrom != stationNameTo {
					// Stations are not equal - move trains
					trainsFrom.Remove(trainName)
					trainsTo.Add(trainName)
				} else {
					// Departure and arrival stations are equal
					SendErrorResponse(w, http.StatusBadRequest,
						"Station %s and %s are equal", stationNameFrom, stationNameTo)
				}
			} else {
				// Destination station is not exists
				SendErrorResponse(w, http.StatusNotFound,
					"Destination station %s does not exist", stationNameTo)
			}
		} else {
			// Train is not stay on the departure station
			SendErrorResponse(w, http.StatusBadRequest,
				"Train %s is not stay on station %s are equal",
				trainName, stationNameFrom)
		}
	} else {
		// Destination station is not exists
		SendErrorResponse(w, http.StatusNotFound,
			"Departure station %s does not exist", stationNameFrom)
	}
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/station", AddStation).Methods("POST")
	router.HandleFunc("/station", DeleteStation).Methods("DELETE")
	router.HandleFunc("/stations", GetStations).Methods("GET")
	router.HandleFunc("/{station}/train", AddTrain).Methods("POST")
	router.HandleFunc("/{station}/train", DeleteTrain).Methods("DELETE")
	router.HandleFunc("/{station}/trains", GetTrains).Methods("GET")
	router.HandleFunc("/trip", Trip).Methods("POST")
	log.Fatal(http.ListenAndServe(":9000", router))

}
