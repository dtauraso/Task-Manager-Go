package Patterns

import (
	"fmt"
	// "runtime"
	"strings"
)

type Node2 struct {
	id       int
	children []int
	next     int
	value    string
	// prevParentPredictsParent map[string]int                 // prev tracker's nextPredictions is parent if 1 nonzero length string("parent") is generated for the parent key
	// nextPredictions          map[string]map[int]interface{} // eventually will reach a node with an empty nextPredictions map (higher level version can map to "parent" key)
	// using trackers of different subtrees to travel unique path for processing when 1 certain key is reached
	// allow all the data to be stored as paths of Node2 items but allow unique paths to be visited for displaying to the user
	// childPredictions map[string][]int
	// prev parent  tracker's next predictions although successfull need to verify they are pointing to the right current child sequence
	// an active tag is active under a certain list and it's parent tag that is also active under the same list must be able to reach the child tag and the specific activation node for the same list
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
		value:    value,
	})

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

	for i := 0; i < len(views); i++ {
		switch views[i] {
		case "today":
			// preview of task
			// tag bar
			// 	activate node in tag tree under "today list"
			// tasks in today list
		default:

		}
	}

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

func (n *Nodes2) dft2(nodeId int, level int, successChan chan struct{}) {
	if nodeId == -1 {
		return
	}
	node := n.nodes[nodeId]
	indent := strings.Repeat(" ", level*2)
	fmt.Printf("%s%s\n", indent, node.value)

	// Create a channel to signal successful visitation of a child node
	childSuccess := make(chan struct{})

	for _, childId := range node.children {
		go func(childID int) {
			n.dft2(childID, level+1, childSuccess)
		}(childId)
	}

	// Wait for child visitation signals
	for range node.children {
		select {
		case <-childSuccess:
			// If one child's visit function is successful, cancel sibling goroutines
			close(successChan)
			return
		}
	}

	// Signal successful visitation of this node
	successChan <- struct{}{}
}

// // Call the dft function with a root node ID to start the traversal
// rootNodeID := 0 // Example root node ID
// successChan := make(chan struct{})
// go nodes2.dft(rootNodeID, 0, successChan)

// // Wait for the traversal to complete
// <-successChan

type Node3 struct {
	next     map[string]map[int]func(x interface{}, id int, dataNodePtr *int) bool
	children map[string]map[int]func(x interface{}, id int, dataNodePtr *int) bool
}

func Init(name string, nextNodeId int, function func(x interface{}, id int, dataNodePtr *int) bool) map[string]map[int]func(x interface{}, id int, dataNodePtr *int) bool {
	return map[string]map[int]func(x interface{}, id int, dataNodePtr *int) bool{name: {nextNodeId: function}}
}

var x1 = map[int]*Node3{
	0: {},
	1: {next: Init("x0", 3, X)},
	3: {next: Init("n0", 4, X),
		children: Init("x1", 5, X1)},
}

var X = func(x interface{}, id int, dataNodePtr *int) bool {

	// item := x.(map[int]map[string]map[int]func(x interface{}, id int) bool)[id]
	return true
}
var X1 = func(x interface{}, id int, dataNodePtr *int) bool {
	return true
}

func traverseX1(x1 map[int]map[string]map[int]func(x interface{}, id int, dataNodePtr *int) bool, id int, dataNodePtr *int) {
	// if len(x1[id]) == 0 {
	// return true
	// }
	/*
	   record of data changes
	   id of item
	   	name of next node "x0"
	   	name of function
	   	direction "next" or "child"
	   	collection of (before, after) strings from function
	*/
	item := x1[id]
	var pass bool
	child, hasChild := item["child"]
	// var nextNodeId int
	if hasChild {
		pass := false
		for key := range child {
			// pc, _, _, _ := runtime.Caller(0)
			// functionName := runtime.FuncForPC(pc).Name()
			if child[key](x1, key, dataNodePtr) {
				pass = true
				traverseX1(x1, key, dataNodePtr)
			}
		}
		if !pass {
			// return false
		}
	}
	var functionPass bool
	// for key := range child {
	// 	functionPass = child[key](x1, key)
	// }
	if !functionPass {
		// return false
	}
	// next, hasNext := item["next"]
	// if !hasNext {
	// return true
	// }
	if !hasChild {

	}
	if !hasChild || pass {

	}
	for key, _ := range x1[id] {
		// pass := false
		if key == "child" {
			// for singleKey := range value {
			// 	// if traverseX1(x1, singleKey) {
			// 	// 	pass = true
			// 	// }
			// }
		}

		// for nextId, function := range value {
		// 	if function(x1, id) {
		// 		// if
		// 		// /*value(id)*/ {
		// 		// if traverseX1(x1, key) {
		// 		// 	return true
		// 		// }
		// 	}
		// }

	}

	// return false
}
func MakeTree() {

	// nodes2 := Nodes2{}
	// views := []string{"today"}
	// lists := map[string]int{
	// 	"todayStart": nodes2.newNode2("todayStart", nil, -1),
	// 	"todayEnd":   nodes2.newNode2("todayEnd", nil, -1)}
	// taskViews := nodes2.addTask(Task{}, views, lists)
	// fmt.Printf("%v\n", taskViews)
	// taskTitleId := nodes2.newNode2("task title", nil, -1)
	// taskTitleAttributeId := nodes2.newNode2(title, nil, taskTitleId)
	// taskTags2Id := nodes2.newNode2("task tag 2", nil, -1)
	// taskTagsId := nodes2.newNode2("task tag", nil, taskTags2Id)
	// taskTagsAttributeId := nodes2.newNode2(tags, nil, taskTagsId)
	// taskTitleRootId := nodes2.newNode2("title field", []int{taskTitleAttributeId}, -1)
	// taskTagsRootId := nodes2.newNode2("tags field", []int{taskTagsAttributeId}, -1)
	// rootId := nodes2.newNode2(
	// 	"0",
	// 	[]int{taskTitleRootId, taskTagsRootId},
	// 	-1,
	// )
	// previewId := nodes2.newNode2("0", []int{taskTitleId, taskTagsId}, -1)
	// nodes2.dft(rootId, 0)
	// fmt.Printf("\n")
	// nodes2.dft(previewId, 0)
	// Create a sample tree structure
	nodes2 := Nodes2{}
	rootNodeID := nodes2.newNode2("Root", []int{1, 2}, -1)
	node1ID := nodes2.newNode2("Node 1", []int{3, 4}, -1)
	node2ID := nodes2.newNode2("Node 2", []int{5}, -1)
	node3ID := nodes2.newNode2("Node 3", nil, -1)
	node4ID := nodes2.newNode2("Node 4", nil, -1)
	node5ID := nodes2.newNode2("Node 5", nil, -1)

	fmt.Println(node1ID)
	fmt.Println(node2ID)
	fmt.Println(node3ID)
	fmt.Println(node4ID)
	fmt.Println(node5ID)

	// Run the depth-first traversal on the tree starting from the root node
	successChan := make(chan struct{})
	go nodes2.dft2(rootNodeID, 0, successChan)

	// Wait for the traversal to complete
	<-successChan
}
