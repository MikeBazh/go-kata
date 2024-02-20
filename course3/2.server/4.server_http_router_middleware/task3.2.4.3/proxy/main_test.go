package main

import (
	"fmt"
	"io"
	"net/http"
	"reflect"
	"testing"
)

func TestNewReverseProxy(t *testing.T) {
	type args struct {
		host string
		port string
	}
	want := &ReverseProxy{
		host: "host1",
		port: "port1",
	}
	tests := []struct {
		name string
		args args
		want *ReverseProxy
	}{
		{"test1", args{"host1", "port1"}, want}, // TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewReverseProxy(tt.args.host, tt.args.port); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewReverseProxy() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_Handle(t *testing.T) {
	type args struct {
		link           string
		wantMessage    string
		wantStatusCode int
	}
	tests := []struct {
		name string
		args args
	}{
		{"test1", args{"http://localhost:8080/api/", "Hello from API", http.StatusOK}},
		{"test2", args{"http://localhost:8080/api/2", "Hello from API", http.StatusOK}},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := http.Get(tt.args.link)
			if err != nil {
				t.Error(err)
			}
			if res.StatusCode != tt.args.wantStatusCode {
				t.Errorf("status want %d, got %d", tt.args.wantStatusCode, res.StatusCode)
			}
			//fmt.Println(res.Request)
			body, _ := io.ReadAll(res.Body)
			result := string(body)
			if tt.args.wantMessage != result {
				t.Errorf("message want %s, got %s", tt.args.wantMessage, result)
			}
		})
	}
}

func Test_redirectMiddleware(t *testing.T) {
	type args struct {
		link           string
		wantMessage    string
		wantStatusCode int
	}
	tests := []struct {
		name string
		args args
	}{

		{"test3", args{"http://localhost:8080/", "<!DOCTYPE html>", http.StatusOK}},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			go main()
			res, err := http.Get(tt.args.link)
			if err != nil {
				t.Error(err)
			}
			if res.StatusCode != tt.args.wantStatusCode {
				t.Errorf("status want %d, got %d", tt.args.wantStatusCode, res.StatusCode)
			}
			//fmt.Println(res.Header)
			for key, values := range res.Header {
				for _, value := range values {
					fmt.Printf("%s: %s\n", key, value)
				}
			}

			body, _ := io.ReadAll(res.Body)
			result := string(body[:15])
			if tt.args.wantMessage != result {
				t.Errorf("message want %s, got %s", tt.args.wantMessage, result)
			}
		})
	}
}
