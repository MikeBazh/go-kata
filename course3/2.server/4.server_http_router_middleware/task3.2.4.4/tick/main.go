//Задание
//Реши 3 пункта из сайта hugo:
//Обновление данных в реальном времени.
//Построение графа.
//Построение сбалансированного бинарного дерева.

package main

import (
	"fmt"
	"github.com/go-chi/chi"
	"go-kata/2.server/4.server_http_router_middleware/task3.2.4.4/tick/graph"
	"go-kata/2.server/4.server_http_router_middleware/task3.2.4.4/tick/tree"
	"log"
	"math/rand"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"
	"time"
)

func main() {
	r := chi.NewRouter()

	// обработчик, который выполняет перенаправление запроса на другой URL
	r.Use(redirectMiddleware)

	// Обработчик для /api/
	r.Get("/api/*", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello from API"))
	})
	go WorkerTree()
	go WorkerTime()
	go WorkerGraph()

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

// Middleware для перенаправления запросов
func redirectMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Если URL начинается с /api/, возвращаем текст "Hello from API"
		if strings.HasPrefix(r.URL.Path, "/api/") {
			next.ServeHTTP(w, r)
			return
		}
		if strings.HasPrefix(r.URL.Path, "/") {
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

func WorkerTime() {
	t := time.NewTicker(5 * time.Second)
	var b byte = 0
	for {
		currentTime := time.Now()
		formattedTime := currentTime.Format("2006-01-02 15:04:05")
		select {
		case <-t.C:
			err := os.WriteFile("/app/static/tasks/_index.md", []byte(fmt.Sprintf(htmlCodeIndex, formattedTime, b)), 0644)
			if err != nil {
				log.Println(err)
			}
			b++
		}
	}
}

func WorkerGraph() {
	t := time.NewTicker(5 * time.Second)
	for {
		select {
		case <-t.C:
			rand.Seed(time.Now().UnixNano())
			numNodes := rand.Intn(26) + 5
			GenGraph := graph.GenerateGraph(numNodes)
			b := graph.Mermaid(numNodes, GenGraph)
			err := os.WriteFile("/app/static/tasks/graph.md", []byte(fmt.Sprintf(htmlCodeGraph, b)), 0644)
			//err := os.WriteFile("/home/m/GolandProjects/GHcourse3/course3/2.server/4.server_http_router_middleware/task3.2.4.4/hugo/content/tasks/graph.md", []byte(fmt.Sprintf(htmlCode, b)), 0644)
			if err != nil {
				log.Println(err)
			}
		}
	}
}

func WorkerTree() {
	t := time.NewTicker(5 * time.Second)
	rand.Seed(time.Now().UnixNano())
	NewTree := tree.GenerateTree(5)
	count := 5
	var b string
	for {
		select {
		case <-t.C:
			b = NewTree.ToMermaid()
			err := os.WriteFile("/app/static/tasks/binary.md", []byte(fmt.Sprintf(htmlCodeTree, b)), 0644)
			if err != nil {
				log.Println(err)
			}
			NewTree.Insert(rand.Intn(1000))
			count++
			if count == 100 {
				NewTree = tree.GenerateTree(5)
				count = 5
			}
		}
	}
}

const htmlCodeGraph = `---
menu:
    after:
        name: graph
        weight: 1
title: Построение графа
---

## Построение графа

{{<mermaid class="text-center">}}
%s
{{< /mermaid >}}`

const htmlCodeTree = `---
menu:
    after:
        name: binary_tree
        weight: 2
title: Построение сбалансированного бинарного дерева
---

# Задача построить сбалансированное бинарное дерево
Используя AVL дерево, постройте сбалансированное бинарное дерево, на текущей странице.

Нужно написать воркер, который стартует дерево с 5 элементов, и каждые 5 секунд добавляет новый элемент в дерево.

Каждые 5 секунд на странице появляется актуальная версия, сбалансированного дерева.

При вставке нового элемента, в дерево, нужно перестраивать дерево, чтобы оно оставалось сбалансированным.

Как только дерево достигнет 100 элементов, генерируется новое дерево с 5 элементами.

{{<mermaid class="text-center">}}
%s
{{< /mermaid >}}`

const htmlCodeIndex = `---
menu:
    before:
        name: tasks
        weight: 5
title: Обновление данных в реальном времени
---

# Задача: Обновление данных в реальном времени

Напишите воркер, который будет обновлять данные в реальном времени, на текущей странице.
Текст данной задачи менять нельзя, только время и счетчик.

Файл данной страницы: /app/static/tasks/_index.md

Должен меняться счетчик и время:

Текущее время: %s

Счетчик: %d



## Критерии приемки:
- [ ] Воркер должен обновлять данные каждые 5 секунд
- [ ] Счетчик должен увеличиваться на 1 каждые 5 секунд
- [ ] Время должно обновляться каждые 5 секунд`
