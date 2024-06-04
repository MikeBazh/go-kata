package main

import (
	"io"
	"net/http"
)

func main() {
	http.HandleFunc("/geo", geoHandler)
	http.HandleFunc("/weather", weatherHandler)
	http.ListenAndServe(":8080", nil)
}

func geoHandler(w http.ResponseWriter, r *http.Request) {
	resp, err := http.Get("http://localhost:8081/geo")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	w.Header().Set("Content-Type", resp.Header.Get("Content-Type"))
	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
}

func weatherHandler(w http.ResponseWriter, r *http.Request) {
	resp, err := http.Get("http://localhost:8082/weather")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	w.Header().Set("Content-Type", resp.Header.Get("Content-Type"))
	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
}
