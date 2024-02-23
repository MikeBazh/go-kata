package temp1

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

// Node представляет узел в графе
type Node struct {
	ID    int
	Name  string
	Form  string // "circle", "rect", "square", "ellipse", "round-rect", "rhombus"
	Links []*Node
}

// generateGraph генерирует случайный граф с numNodes узлами
// generateGraph генерирует случайный граф с numNodes узлами
// generateGraph генерирует случайный граф с numNodes узлами
func generateGraph(numNodes int) []*Node {
	rand.Seed(time.Now().UnixNano())

	// Создаем узлы графа
	nodes := make([]*Node, numNodes)
	for i := 0; i < numNodes; i++ {
		nodes[i] = &Node{
			ID:   i,
			Name: fmt.Sprintf("Node %d", i),
			Form: getRandomForm(),
		}
	}

	// Создаем связи между узлами графа
	for i := 0; i < numNodes; i++ {
		for j := i + 1; j < numNodes; j++ { // Начинаем с j = i + 1, чтобы избежать дублирования связей
			// Исключаем связь на самого себя и обратные связи
			if i != j && !isLinked(nodes[i], nodes[j]) {
				nodes[i].Links = append(nodes[i].Links, nodes[j])
			}
		}
	}

	return nodes
}

// getRandomForm возвращает случайную форму для узла
func getRandomForm() string {
	forms := []string{"circle", "rect", "round-rect", "rhombus"}
	return forms[rand.Intn(len(forms))]
}

// isLinked проверяет, есть ли связь между двумя узлами
func isLinked(node1, node2 *Node) bool {
	for _, link := range node1.Links {
		if link == node2 {
			return true
		}
	}
	return false
}

func MCode() string {
	//numNodes := rand.Intn(26) + 5
	numNodes := 5
	graph := generateGraph(numNodes)

	// Генерация кода Mermaid и запись его в строку
	var mermaidCode strings.Builder
	mermaidCode.WriteString("graph LR\n")
	for _, node := range graph {
		nodeID := fmt.Sprintf("Node%d", node.ID)

		// Форматируем имя узла в зависимости от его формы
		if node.Form == "round-rect" {
			mermaidCode.WriteString(fmt.Sprintf("  %s(%s%s)\n", nodeID, node.Form, node.Name))
		} else if node.Form == "circle" {
			mermaidCode.WriteString(fmt.Sprintf("  %s((%s%s))\n", nodeID, node.Form, node.Name))
			//} else if node.Form == "square" {
			//	mermaidCode.WriteString(fmt.Sprintf("  %s[%s]\n", nodeID, node.Form))
		} else if node.Form == "rect" {
			mermaidCode.WriteString(fmt.Sprintf("  %s[%s%s]\n", nodeID, node.Form, node.Name))
		} else if node.Form == "rhombus" {
			mermaidCode.WriteString(fmt.Sprintf("  %s{%s%s}\n", nodeID, node.Form, node.Name))
		}

		for _, link := range node.Links {
			linkID := fmt.Sprintf("Node%d", link.ID)
			mermaidCode.WriteString(fmt.Sprintf("  %s --> %s\n", nodeID, linkID))
		}
	}

	// Вывод кода Mermaid
	fmt.Println(mermaidCode.String())
	return mermaidCode.String()
}

//func main() {
//	fmt.Println(MCode())
//}
