package Patterns

import "fmt"

type Node2 struct {
	id       int
	children []int
	next     int
	value    string
}

var nodes2 = []*Node2{}

func newNode2(
	nodes *[]*Node2,
	value string,
	children []int,
	next int,
) int {

	newNodeId := len(*nodes)
	(*nodes) = append(*nodes, &Node2{
		id:       newNodeId,
		children: children,
		next:     next,
		value:    value})

	return newNodeId
}

const (
	title = "title"
	tags  = "tags"
)

func MakeTree() {
	taskTitleId := newNode2(&nodes2, "task title", nil, -1)
	taskTitleAttributeId := newNode2(&nodes2, title, nil, taskTitleId)
	taskTagsId := newNode2(&nodes2, "task tag", nil, taskTitleId)
	taskTagsAttributeId := newNode2(&nodes2, tags, nil, taskTagsId)
	rootId := newNode2(
		&nodes2,
		"title field",
		[]int{taskTitleAttributeId, taskTagsAttributeId},
		-1,
	)
	for _, node := range nodes2 {
		fmt.Printf("Node ID: %d, Value: %s, Children: %v, Next: %d\n", node.id, node.value, node.children, node.next)
	}

	fmt.Printf("%v\n", rootId)
}
