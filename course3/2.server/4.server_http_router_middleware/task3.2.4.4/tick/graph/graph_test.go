package graph

import (
	"reflect"
	"testing"
)

func TestGenerateGraph(t *testing.T) {
	type args struct {
		numNodes int
	}
	Node1 := Node{1, "n", "rect", nil}
	Node2 := Node{1, "n", "circle", nil}
	Node3 := Node{1, "n", "square", nil}
	Node1.Links = []*Node{&Node2, &Node3}
	Node2.Links = []*Node{&Node3}
	wantGraph := []*Node{&Node1, &Node2, &Node3}
	tests := []struct {
		name string
		args args
		want []*Node
	}{
		{"tes1_3Nodes", args{3}, wantGraph}, // TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GenerateGraph(tt.args.numNodes); reflect.TypeOf(got) != reflect.TypeOf(tt.want) {
				t.Errorf("GenerateGraph() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMermaid(t *testing.T) {
	type args struct {
		numNodes int
		graph    []*Node
	}
	Node1 := Node{1, "n", "rect", nil}
	Node2 := Node{2, "n", "circle", nil}
	Node3 := Node{3, "n", "square", nil}
	Node1.Links = []*Node{&Node2, &Node3}
	Node2.Links = []*Node{&Node3}
	Graph := []*Node{&Node2}
	wantCode2 := "graph LR\n" + "Node2((circle))\n" + "Node2 --> Node3\n"
	tests := []struct {
		name string
		args args
		want string
	}{
		{"test1", args{3, Graph}, wantCode2}, // TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Mermaid(tt.args.numNodes, tt.args.graph); got != tt.want {
				t.Errorf("Mermaid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getRandomForm(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{"test1", "randomform"}, // TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getRandomForm(); reflect.TypeOf(got) != reflect.TypeOf(tt.want) {
				t.Errorf("getRandomForm() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_isLinked(t *testing.T) {
	type args struct {
		node1 *Node
		node2 *Node
	}
	Node1 := Node{1, "n", "rect", nil}
	Node2 := Node{2, "n", "circle", nil}
	Node3 := Node{3, "n", "square", nil}
	Node1.Links = []*Node{&Node2, &Node3}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"test1", args{&Node1, &Node2}, true},
		{"test2", args{&Node2, &Node3}, false}, // TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isLinked(tt.args.node1, tt.args.node2); got != tt.want {
				t.Errorf("isLinked() = %v, want %v", got, tt.want)
			}
		})
	}
}
