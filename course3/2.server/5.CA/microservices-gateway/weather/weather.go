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

type WeatherResponse struct {
	Weather []struct {
		Description string `json:"description"`
	} `json:"weather"`
	Main struct {
		Temp float64 `json:"temp"`
	} `json:"main"`
}

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Get("/weather", getWeather)

	http.ListenAndServe(":8082", r)
}

func getWeather(w http.ResponseWriter, r *http.Request) {
	ip, err := getExternalIP()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	geo, err := getGeoLocation(ip)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	weather, err := getWeatherByCoordinates(geo.Lat, geo.Lon)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(weather)
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
	//fmt.Println(fmt.Sprintf("http://ip-api.com/json/%s", ip))
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

func getWeatherByCoordinates(lat, lon float64) (*WeatherResponse, error) {
	apiKey := "bd5e378503939ddaee76f12ad7a97608"

	resp, err := http.Get(fmt.Sprintf("https://api.openweathermap.org/data/2.5/weather?lat=%f&lon=%f&units=metric&appid=%s", lat, lon, apiKey))
	//fmt.Println(fmt.Sprintf("https://api.openweathermap.org/data/2.5/weather?lat=%f&lon=%f&units=metric&appid=%s", lat, lon, apiKey))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var weatherResponse WeatherResponse
	if err := json.NewDecoder(resp.Body).Decode(&weatherResponse); err != nil {
		return nil, err
	}

	return &weatherResponse, nil
}
