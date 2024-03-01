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

type Nodes2 struct {
	nodes []*Node2
}

func (n *Nodes2) newNode2(
	value string,
	children []int,
	next int,
) int {

	newNodeId := len(n.nodes)
	n.nodes = append(n.nodes, &Node2{
		id:       newNodeId,
		children: children,
		next:     next,
		value:    value})

	return newNodeId
}

func (n *Nodes2) dft(nodeId int, level int) {
	if nodeId == -1 {
		return
	}
	node := n.nodes[nodeId]
	indent := strings.Repeat(" ", level*2)
	fmt.Printf("%s%s\n", indent, node.value)
	for _, childId := range node.children {
		n.dft(childId, level+1)

	}
	n.dft(node.next, level)
}

func (n *Nodes2) addTask(task Task, views []string, lists map[string]int) (taskViews map[string]int) {

	return map[string]int{}
}

type Task struct {
	title string
	tags  []string
}

const (
	title = "title"
	tags  = "tags"
)

func MakeTree() {

	nodes2 := Nodes2{}
	views := []string{"today"}
	lists := map[string]int{"todayStart": -1, "todayEnd": -1}
	taskViews := nodes2.addTask(Task{}, views, lists)
	fmt.Printf("%v\n", taskViews)
	taskTitleId := nodes2.newNode2("task title", nil, -1)
	taskTitleAttributeId := nodes2.newNode2(title, nil, taskTitleId)
	taskTags2Id := nodes2.newNode2("task tag 2", nil, -1)
	taskTagsId := nodes2.newNode2("task tag", nil, taskTags2Id)
	taskTagsAttributeId := nodes2.newNode2(tags, nil, taskTagsId)
	taskTitleRootId := nodes2.newNode2("title field", []int{taskTitleAttributeId}, -1)
	taskTagsRootId := nodes2.newNode2("tags field", []int{taskTagsAttributeId}, -1)
	rootId := nodes2.newNode2(
		"0",
		[]int{taskTitleRootId, taskTagsRootId},
		-1,
	)
	previewId := nodes2.newNode2("0", []int{taskTitleId, taskTagsId}, -1)
	nodes2.dft(rootId, 0)
	fmt.Printf("\n")
	nodes2.dft(previewId, 0)

}
