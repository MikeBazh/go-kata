package main

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"net/http"
)

type IPResponse struct {
	IP string `json:"ip"`
}

type GeoResponse struct {
	Lat float64 `json:"lat"`
	Lon float64 `json:"lon"`
}

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Get("/geo", getGeo)

	http.ListenAndServe(":8081", r)
}

func getGeo(w http.ResponseWriter, r *http.Request) {
	ip, err := getExternalIP()
	fmt.Println(ip)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	geo, err := getGeoLocation(ip)
	fmt.Println(geo)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(geo)
}

func getExternalIP() (string, error) {
	resp, err := http.Get("https://api.ipify.org?format=json")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var ipResponse IPResponse
	if err := json.NewDecoder(resp.Body).Decode(&ipResponse); err != nil {
		return "", err
	}

	return ipResponse.IP, nil
}

func getGeoLocation(ip string) (*GeoResponse, error) {
	resp, err := http.Get(fmt.Sprintf("http://ip-api.com/json/%s", ip))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var geoResponse GeoResponse
	if err := json.NewDecoder(resp.Body).Decode(&geoResponse); err != nil {
		return nil, err
	}

	return &geoResponse, nil
}
