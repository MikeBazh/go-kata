package tree

import (
	"strings"
	"testing"
)

func TestAVLTree_ToMermaid(t1 *testing.T) {
	type fields struct {
		Root *Node
	}
	node3 := &Node{3, 1, nil, nil}
	node2 := &Node{1, 1, nil, nil}
	node := &Node{2, 2, node2, node3}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{"test1", fields{node}, "graph TD;\n  2\n  2 --> 1\n  1\n  2 --> 3\n  3\n"}, // TODO: Add test cases.
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := &AVLTree{
				Root: tt.fields.Root,
			}
			if got := t.ToMermaid(); got != tt.want {
				t1.Errorf("ToMermaid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_generateMermaid(t *testing.T) {
	type args struct {
		node    *Node
		builder *strings.Builder
	}
	node3 := &Node{3, 1, nil, nil}
	node2 := &Node{1, 1, nil, nil}
	node := &Node{2, 2, node2, node3}
	//var builder *strings.Builder
	builder := &strings.Builder{}
	tests := []struct {
		name string
		args args
	}{
		{"test1", args{node, builder}}, // TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			generateMermaid(tt.args.node, tt.args.builder)
		})
	}
}
