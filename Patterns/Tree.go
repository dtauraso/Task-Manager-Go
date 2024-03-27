package Patterns

import (
	"fmt"
	"reflect"
	"strconv"

	// "runtime"

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

type NextNode struct {
	id int
	// function func(x map[int]*Node3, id int, dataNodePtr *int, changes *DataChange) bool
}

type Node3 struct {
	id            int
	name          string
	parentNodeIds []int
	// nil if len(childrenNodes) > 0
	function       func(x map[int]*Node3, controlFlowNodeId int, dataNodeId *int) int
	functionPassed bool

	beforeAfterList []map[string]string
	dataNodeName    string
	// only 1 next
	nextNodeId int
	// n directions to try
	directionNodeIds []int

	value              int
	variables          map[string]int
	variableCollection []int
}

/*
graphical nodes across and their distances

a
\→ b0 ←←←←←←←←←

	   \→ c0      ↑
	      ↻\→ d0 ←↑←
	          ↓   ↑↑
	          e0 →↑↑
	   \→ c1       ↑
		  ↓\→ d1 →→↑
		  ↓	  ↓
		  ↓   e1
		  ↓\→ d3
		  ↓   ↓
		  ↓   e0
		  c2
		  ↓
		  b0
		  ↓
		  e0
		  ↓
		  c3
	   \→ e0

\→ b1
*/

const (
	levelId            = 0
	parentStateHasNext = 1
	levelIdNext        = 2
	// node maps to sequence of edge kinds horizontally from other nodes
)

const (
	characterToId                = 0
	ids1                         = 1
	ids2                         = 2
	input                        = 3
	newItemStreakCount           = 4
	useFrequencyPercentThreshold = 5
	idToCharacter                = 6
)

var reuseAttributes = map[int]*Node3{
	characterToId: {name: "characterToId", variables: map[string]int{}},
	idToCharacter: {name: "idToCharacter", variables: map[string]int{}},
	ids1:          {name: "ids1", variableCollection: []int{}},
	ids2:          {name: "ids2", variableCollection: []int{}},
	input:         {name: "input", variableCollection: []int{}},
}

// inputs
/*
111000101010
01010100110
010101101010

items saved
1, 0, 10, 01, 010101101010
*/
func loadInput(inputString string, reuseAttributes map[int]*Node3) {
	nextId := len(reuseAttributes)
	for _, char := range inputString {
		reuseAttributes[nextId] = &Node3{name: string(char)}
		reuseAttributes[input].variableCollection = append(reuseAttributes[input].variableCollection, nextId)
		nextId += 1
	}
}

const (
	a                                  = 0
	b0                                 = 1
	b1                                 = 2
	before                             = 3
	after                              = 4
	computeWaitTimeDuration            = 5
	targetTimeIsNotReached             = 6
	targetTimeIsReached                = 7
	requestFailed                      = 8
	requestSucceeded                   = 9
	NoNextNode                         = -1
	loopOverSequence                   = 0
	processItem                        = 1
	startPrediction                    = 2
	itemIsNew                          = 3
	itemIsKnown                        = 4
	itemIsUsedOften                    = 11
	predictNextItem                    = 5
	saveNewItem                        = 6
	addNewItemToOutputPrediction       = 7
	inputSequence                      = 8
	predictionSequence                 = 9
	savedOutputSequence                = 10
	currentItemStartsExistingPattern   = 12
	currentItemStartsPatternInProgress = 13
)

var reuseTree = map[int]*Node3{
	loopOverSequence:                   {name: "loopOverSequence", directionNodeIds: []int{processItem}},
	processItem:                        {name: "processItem", directionNodeIds: []int{itemIsKnown, itemIsNew}},
	itemIsKnown:                        {name: "itemIsKnown", directionNodeIds: []int{currentItemStartsExistingPattern, currentItemStartsPatternInProgress}},
	currentItemStartsExistingPattern:   {name: "currentItemStartsExistingPattern", directionNodeIds: []int{}},
	currentItemStartsPatternInProgress: {name: "currentItemStartsPatternInProgress", directionNodeIds: []int{}},

	// a:                       {name: "a", nextNodeId: b0},
	// b0:                      {name: "b0", nextNodeId: targetTimeIsNotReached, childrenNodeIds: []int{b0, b1}},
	// b1:                      {name: "b1", nextNodeId: requestFailed, childrenNodeIds: []int{before, targetTimeIsReached}},
	// before:                  {name: "before", nextNodeId: after},
	// after:                   {name: "after", nextNodeId: computeWaitTimeDuration},
	// computeWaitTimeDuration: {name: "computeWaitTimeDuration", nextNodeId: targetTimeIsNotReached},
	// targetTimeIsNotReached:  {name: "targetTimeIsNotReached", nextNodeId: before},
	// targetTimeIsReached:     {name: "targetTimeIsReached", nextNodeId: NoNextNode},
	// requestFailed:           {name: "requestFailed", nextNodeId: before},
	// requestSucceeded:        {name: "requestSucceeded", nextNodeId: NoNextNode},
}

var x1 = map[int]*Node3{
	0: {},
	1: {name: "", function: X, nextNodeId: 3},
	3: {name: "x0", function: nil, nextNodeId: 4,
		directionNodeIds: []int{5}},
	4: {name: "n0"},
	5: {name: "x1", function: X1},
	6: {name: "dataModel",
		variables: map[string]int{"varName": 7}},
	7: {name: "1"},
}

// need the control flow node to be connected to the data node ids
func SetNode3(x map[int]*Node3, controlFlowNodeId int, dataNodeId *int, after string) {
	node := x[controlFlowNodeId]
	if !structAttributeExists(node, "beforeAfterList") {
		node.beforeAfterList = []map[string]string{}
	}
	node.beforeAfterList = append(node.beforeAfterList, map[string]string{node.name: after})
	node.dataNodeName = x[*dataNodeId].name
	node.name = after
}

var X = func(x map[int]*Node3, controlFlowNodeId int, dataNodeId *int) int {

	itemId := x[*dataNodeId].variables["varName"]
	item := x[itemId]
	index, _ := strconv.Atoi(item.name)
	index += 1

	SetNode3(x, controlFlowNodeId, dataNodeId, strconv.Itoa(index))

	return 1
}

var X1 = func(x map[int]*Node3, controlFlowNodeId int, dataNodeId *int) int {
	*dataNodeId = x[*dataNodeId].variables["containerVarName"]
	return 1
}
var returnTrue = func(x map[int]*Node3, controlFlowNodeId int, dataNodeId *int) int {
	return 1
}
var returnFalse = func(x map[int]*Node3, controlFlowNodeId int, dataNodeId *int) int {
	return 0
}

func structAttributeExists(s interface{}, attributeName string) bool {
	v := reflect.ValueOf(s)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	_, found := v.Type().FieldByName(attributeName)
	return found
}

func traverseX1(
	x1 map[int]*Node3,
	id int,
	dataNodeId *int,
) int {

	functionPassCount := 1
	for functionPassCount > 0 {
		functionPassCount = 0
		item := x1[id]
		isChild := len((*item).directionNodeIds) > 0
		if !isChild {
			before := functionPassCount
			functionPassCount += item.function(x1, id, dataNodeId)
			after := functionPassCount - before
			if after == 1 {
				item.functionPassed = true
			}
		} else if isChild {
			children := item.directionNodeIds
			for key := range children {
				before := functionPassCount
				functionPassCount += traverseX1(x1, key, dataNodeId)
				after := functionPassCount - before
				if after == 1 {
					break
				}
			}
		}
		if functionPassCount == 0 {
			break
		}
		if item.nextNodeId == -1 {
			break
		}
		id = item.nextNodeId
	}
	return functionPassCount

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
	// rootNodeID := nodes2.newNode2("Root", []int{1, 2}, -1)
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
	// changes := []DataChange{}
	// dataPtr := 0
	// traverseX1(rootNodeID, 0, successChan, &dataPtr, &changes)

	// Wait for the traversal to complete
	<-successChan
}
