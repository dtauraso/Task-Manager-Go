package Patterns

import "fmt"

type Node2 struct {
	id       int
	children []int
	next     int
	value    string
}

var nodes2 = map[int]*Node2{}

func newNode2(
	nodes *map[int]*Node2,
	id int,
	children []int,
	next int,
	value string) int {

	newNodeId := len(*nodes)
	(*nodes)[newNodeId] = &Node2{
		id:       newNodeId,
		children: children,
		next:     next,
		value:    value}

	return newNodeId
}

func MakeTree() {
	rootId := newNode2(
		&nodes2,
		0,
		[]int{newNode2(&nodes2, 1, nil, 2, "child1"),
			newNode2(&nodes2, 2, nil, 3, "child2")},
		3,
		"root")
	for _, node := range nodes2 {
		fmt.Printf("Node ID: %d, Value: %s, Children: %v, Next: %d\n", node.id, node.value, node.children, node.next)
	}

	fmt.Printf("%v\n", rootId)
}
