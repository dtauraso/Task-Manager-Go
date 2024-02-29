package Patterns

import (
	"fmt"
	"strings"
)

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

func dft(nodeId int, level int) {
	if nodeId == -1 {
		return
	}
	node := nodes2[nodeId]
	indent := strings.Repeat(" ", level*2)
	fmt.Printf("%s%s\n", indent, node.value)
	for _, childId := range node.children {
		dft(childId, level+1)

	}
	dft(node.next, level)
}

func MakeTree() {
	taskTitleId := newNode2(&nodes2, "task title", nil, -1)
	taskTitleAttributeId := newNode2(&nodes2, title, nil, taskTitleId)
	taskTags2Id := newNode2(&nodes2, "task tag 2", nil, -1)
	taskTagsId := newNode2(&nodes2, "task tag", nil, taskTags2Id)
	taskTagsAttributeId := newNode2(&nodes2, tags, nil, taskTagsId)
	taskTitleRootId := newNode2(&nodes2, "title field", []int{taskTitleAttributeId}, -1)
	taskTagsRootId := newNode2(&nodes2, "tags field", []int{taskTagsAttributeId}, -1)
	rootId := newNode2(
		&nodes2,
		"task",
		[]int{taskTitleRootId, taskTagsRootId},
		-1,
	)
	previewId := newNode2(&nodes2, "preview", []int{taskTitleId, taskTagsId}, -1)
	dft(rootId, 0)
	fmt.Printf("\n")
	dft(previewId, 0)

}
