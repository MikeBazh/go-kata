package main

import (
	"fmt"
	"github.com/go-chi/chi"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

func main() {
	r := chi.NewRouter()

	// обработчик, который выполняет перенаправление запроса на другой URL
	r.Use(redirectMiddleware)

	// Обработчик для /api/
	r.Get("/api/*", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello from API"))
	})

	err := http.ListenAndServe(":8080", r)
	if err != nil {
		fmt.Println(err)
	}
}

type ReverseProxy struct {
	host string
	port string
}

func NewReverseProxy(host, port string) *ReverseProxy {
	return &ReverseProxy{
		host: host,
		port: port,
	}
}

//const content = ``

// Middleware для перенаправления запросов
func redirectMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Если URL начинается с /api/, возвращаем текст "Hello from API"
		if strings.HasPrefix(r.URL.Path, "/api/") {
			next.ServeHTTP(w, r)
			return
		}
		if strings.HasPrefix(r.URL.Path, "/") && r.URL.Path != "/api/" {
			// перенаправляем запрос на http://hugo:1313
			proxy := NewReverseProxy("hugo", "1313")
			proxy.ServeHTTP(w, r)
			return
		}

		// Если условия не выполнены, передаем обработку следующему обработчику
		proxy := NewReverseProxy("hugo", "1313")
		proxy.ServeHTTP(w, r)
	})
}

// ServeHTTP для ReverseProxy
func (rp *ReverseProxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	targetURL := &url.URL{
		Scheme: "http",
		Host:   rp.host + ":" + rp.port,
	}

	proxy := httputil.NewSingleHostReverseProxy(targetURL)
	proxy.ServeHTTP(w, r)
}

//func WorkerTest() {
//	t := time.NewTicker(1 * time.Second)
//	var b byte = 0
//	for {
//		select {
//		case <-t.C:
//			err := os.WriteFile("/app/static/_index.md", []byte(fmt.Sprintf(content, b)), 0644)
//			if err != nil {
//				log.Println(err)
//			}
//			b++
//		}
//	}
//}
