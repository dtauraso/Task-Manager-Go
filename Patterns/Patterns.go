package Patterns

import (
	"fmt"
)

const (
	x = "x"
	y = "y"
	z = "z"
)

func R3Test(v *Variables) bool {

	x, okX := v.State[x]

	y, okY := v.State[y]

	z, okZ := v.State[z]

	if !okX {
		return false
	}
	if _, okXInt := x.(int); !okXInt {
		return false
	}
	if !okY {
		return false
	}
	if _, okYInt := y.(int); !okYInt {
		return false
	}
	if !okZ {
		return false
	}
	if _, okZInt := z.(int); !okZInt {
		return false
	}
	return true
}

func add1(x int) int {
	return x + 1
}

func subtract1(x int) int {
	return x - 1
}

func move1Unit(v *Variables, c *Caretaker, dimensionName string, direction func(int) int) {

	c.UpdateMemento(v.CreateMemento())
	dimension := v.State[dimensionName].(int)
	dimension = direction(dimension)
	v.State[dimensionName] = dimension

}
func moveForward1UnitX(v *Variables, c *Caretaker)  { move1Unit(v, c, x, add1) }
func moveForward1UnitY(v *Variables, c *Caretaker)  { move1Unit(v, c, y, add1) }
func moveForward1UnitZ(v *Variables, c *Caretaker)  { move1Unit(v, c, z, add1) }
func moveBackward1UnitX(v *Variables, c *Caretaker) { move1Unit(v, c, x, subtract1) }
func moveBackward1UnitY(v *Variables, c *Caretaker) { move1Unit(v, c, y, subtract1) }
func moveBackward1UnitZ(v *Variables, c *Caretaker) { move1Unit(v, c, z, subtract1) }

type Operation struct {
	Id           int
	VariableName string
	FunctionName string
	TypeName     string
}
type Node1 struct {
	Id int
	// OperationId   int
	VariableName  string
	FunctionName  string
	TypeName      string
	Edges         map[string][]int
	ParentChildId int
}

func (n1 *Node1) GetLastEdge(edgeName string) int {
	return n1.Edges[edgeName][len(n1.Edges[edgeName])-1]
}

type Storage struct {
	Id           int
	Node1Id      int
	StreakLength int
}
type Variables struct {
	State             map[string]interface{}
	IfConditionResult bool
}

func (v *Variables) CreateMemento() Memento {
	memento := map[string]interface{}{}

	for key, value := range v.State {
		memento[key] = value
	}
	return Memento{State: memento}
}
func (v *Variables) SetMemento(m Memento) {
	v.State = m.State
}

type Memento struct {
	State map[string]interface{}
}

type Caretaker struct {
	memento Memento
}

// {sequenceVarName: Memento} of mementos for each sequence to process
func (c *Caretaker) UpdateMemento(m Memento) {
	c.memento = m

}

func (c *Caretaker) GetMemento() Memento {
	return c.memento
}

const (
	mF1UX = "moveForward1UnitX"
	mF1UY = "moveForward1UnitY"
	mF1UZ = "moveForward1UnitZ"
	mB1UX = "moveBackward1UnitX"
	mB1UY = "moveBackward1UnitY"
	mB1UZ = "moveBackward1UnitZ"
)

var functions = map[string]interface{}{
	mF1UX: moveForward1UnitX,
	mF1UY: moveForward1UnitY,
	mF1UZ: moveForward1UnitZ,
	mB1UX: moveBackward1UnitX,
	mB1UY: moveBackward1UnitY,
	mB1UZ: moveBackward1UnitZ,
}

func equal(a1, a2 interface{}) bool {
	return a1 == a2
}

func notEqual(a1, a2 interface{}) bool {
	return a1 != a2
}

type UniqueOrderedSet struct {
	itemsMap   map[int]struct{}
	itemsSlice []int
}

func NewUniqueOrderedSet() *UniqueOrderedSet {
	return &UniqueOrderedSet{
		itemsMap:   make(map[int]struct{}),
		itemsSlice: []int{},
	}
}

func (s *UniqueOrderedSet) Add(item int) {
	if _, exists := s.itemsMap[item]; !exists {
		s.itemsMap[item] = struct{}{}
		s.itemsSlice = append(s.itemsSlice, item)
	}
}

type SequenceHierarchy struct {
	Sequences                    *[]*Node1
	FunctionNameToNodeIds        map[string]*map[int]int
	FirstNodeIdLastSequenceAdded int
	NodeIdsLastSequenceAdded     map[int]struct{}
}

func (sh *SequenceHierarchy) CreateSequenceOfOperationChangeNames(
	v *Variables,
	c *Caretaker,
	sequence []string) {

	// when the command changes
	// note what variable values changed
	// record the changes as a sequence of operation change names

	head := -1
	prev := head
	lastOperationName := ""
	sh.NodeIdsLastSequenceAdded = map[int]struct{}{}
	functionNameOccurrenceCounts := map[string]int{}
	// {operation name: {node id(s) of occurrence}}
	sh.FirstNodeIdLastSequenceAdded = len(*sh.Sequences)
	for _, functionName := range sequence {
		functions[functionName].(func(v *Variables, c *Caretaker))(v, c)
		if functionName != lastOperationName {
			if _, ok := functionNameOccurrenceCounts[functionName]; !ok {
				functionNameOccurrenceCounts[functionName] = 1
			} else {
				functionNameOccurrenceCounts[functionName] += 1
			}
			changedVariableName := ""
			typeName := ""
			// likely to be O(1) due to each operation only changing 1 variable at a time
			for variableName, value := range v.State {
				prevValue := c.GetMemento().State[variableName]
				if value != prevValue {
					changedVariableName = variableName
					typeName = fmt.Sprintf("%T", value)
				}
			}

			newNodeId := len(*sh.Sequences)
			sh.NodeIdsLastSequenceAdded[newNodeId] = struct{}{}
			if pointer := sh.FunctionNameToNodeIds[functionName]; pointer == nil {
				sh.FunctionNameToNodeIds[functionName] = &map[int]int{}
			}
			(*sh.FunctionNameToNodeIds[functionName])[newNodeId] = functionNameOccurrenceCounts[functionName]

			temp := Node1{
				Id:           newNodeId,
				VariableName: changedVariableName,
				FunctionName: functionName,
				TypeName:     typeName,
				Edges:        map[string][]int{"prev": {prev}, "next": {-1}}}
			if prev >= 0 {
				newEdges := (*sh.Sequences)[prev].Edges
				newEdges["next"] = []int{temp.Id}
				(*sh.Sequences)[prev].Edges = newEdges
			}
			*sh.Sequences = append(*sh.Sequences, &temp)
			prev = temp.Id
		}
		lastOperationName = functionName
	}
}

var operations = map[int]Operation{
	0: {VariableName: "x", FunctionName: mF1UX, TypeName: "int"},
	1: {VariableName: "x", FunctionName: mB1UX, TypeName: "int"},
	2: {VariableName: "y", FunctionName: mF1UY, TypeName: "int"},
	3: {VariableName: "y", FunctionName: mB1UY, TypeName: "int"},
	4: {VariableName: "z", FunctionName: mF1UZ, TypeName: "int"},
	5: {VariableName: "z", FunctionName: mB1UZ, TypeName: "int"}}
var sequences = []Node1{}

type SequencePair struct {
	A1 int
	A2 int
}

type CategoryTracker struct {
	IsVisited                bool
	TotalSequenceLengthFound int
	SequencePairMatches      SequencePair
}

var catagoryTracker = map[int]CategoryTracker{}

// {sequence node id: {isVisited, totalSequenceLengthFound, sequenceIndexPairMatches}
// sequenceIndexPairMatches is list of node id pairs from stored sequence and new sequence that match
// use isVisited to remove duplicate sequence trackers if they are discovered
// 1 entry per sequence after algorithm is done {last sequence node id: {isVisited: true, totalSequenceLengthFound: n}

// 1) make sequence
// 2) find the connections between new sequece and already existing sequence

func (sh *SequenceHierarchy) Categorize() {

	// trackingDict := map[int]CategoryTracker{}
	nodeIdMatches := map[int]int{}
	functionNameCurrentOccurrenceCount := map[string]int{}
	newSequenceIdTracker := sh.FirstNodeIdLastSequenceAdded
	sizeOfNewSequence := 0
	for ; newSequenceIdTracker != -1; newSequenceIdTracker = (*sh.Sequences)[newSequenceIdTracker].GetLastEdge("next") {

		functionNameNewSequence := (*sh.Sequences)[newSequenceIdTracker].FunctionName
		// prevents 1 occurrence from connecting to more than 1 occurrences in previously saved sequences
		if _, isOccurrenceRecord := functionNameCurrentOccurrenceCount[functionNameNewSequence]; !isOccurrenceRecord {
			functionNameCurrentOccurrenceCount[functionNameNewSequence] = 1
		} else {
			functionNameCurrentOccurrenceCount[functionNameNewSequence] += 1
		}
		nodeIds := sh.FunctionNameToNodeIds[functionNameNewSequence]
		for nodeId, occurrenceCount := range *nodeIds {
			if _, isNodeIdInNewSequence := sh.NodeIdsLastSequenceAdded[nodeId]; isNodeIdInNewSequence {
				continue
			}
			if occurrenceCount != functionNameCurrentOccurrenceCount[functionNameNewSequence] {
				continue
			}
			// nodeId node matches with newSequenceIdTracker node
			nodeIdMatches[nodeId] = newSequenceIdTracker
		}
		sizeOfNewSequence += 1
	}
	fmt.Printf("%v\n", nodeIdMatches)
	x := map[int]struct{}{}
	for key := range nodeIdMatches {
		x[key] = struct{}{}
	}
	visited := x

	advancedCount := 1
	for advancedCount > 0 {
		advancedCount = 0
		y := map[int]struct{}{}
		for nodeId := range x {
			nextNodeId := (*sh.Sequences)[nodeId].GetLastEdge("next")
			// skip over if at end of sequence
			if nextNodeId == -1 {
				continue
			}
			// skip over if node id has already been visited
			if _, ok := visited[nextNodeId]; ok {
				continue
			}
			visited[nextNodeId] = struct{}{}
			advancedCount += 1
			y[nextNodeId] = struct{}{}
		}

		if advancedCount > 0 {
			x = y
		}
	}
	fmt.Printf("%v\n", x)

	nodeIdSize := map[int]int{}
	for nodeId, _ := range x {
		nodeIdSize[nodeId] = 1
	}
	atBeginingCount := 0
	nodeIdSizeLength := len(nodeIdSize)
	for atBeginingCount < nodeIdSizeLength {
		for nodeId := range nodeIdSize {
			prevNodeId := (*sh.Sequences)[nodeId].GetLastEdge("prev")
			if prevNodeId == -1 {
				atBeginingCount += 1
				continue
			}
			nodeIdSize[prevNodeId] = nodeIdSize[nodeId] + 1
			delete(nodeIdSize, nodeId)
			// fmt.Printf("%v %v,%v\n", atBeginingCount, nodeIdSizeLength, prevNodeId)
		}

	}
	fmt.Printf("node id size %v\n", nodeIdSize)
	fmt.Printf("sequence size %v\n", sizeOfNewSequence)
	// todo: The new sequence is a copy of 1 of the previous sequences.
	// todo: The new and previous sequences have the same length, but number of matches < length of each sequence.
	for nodeId := range nodeIdSize {
		previousSequenceId := nodeId
		for ; previousSequenceId != -1; previousSequenceId = (*sh.Sequences)[previousSequenceId].GetLastEdge("next") {
			if newSequenceId, match := nodeIdMatches[previousSequenceId]; match {

				newSequenceEdges := (*sh.Sequences)[newSequenceId].Edges
				previousSequenceEdges := (*sh.Sequences)[previousSequenceId].Edges

				newSequenceSize := sizeOfNewSequence
				previousSequenceSize := nodeIdSize[nodeId]
				if newSequenceSize > previousSequenceSize {
					newSequenceEdges["parent"] = []int{previousSequenceId}
					if _, isChildKeyPresent := previousSequenceEdges["child"]; !isChildKeyPresent {
						previousSequenceEdges["child"] = []int{newSequenceId}
					} else {
						previousSequenceEdges["child"] = append(previousSequenceEdges["child"], newSequenceId)
					}
				} else if newSequenceSize < previousSequenceSize {
					if _, isParentKeyPresent := previousSequenceEdges["parent"]; !isParentKeyPresent {
						previousSequenceEdges["parent"] = []int{newSequenceId}
					} else {
						previousSequenceEdges["parent"] = append(newSequenceEdges["parent"], newSequenceId)
					}
					if _, isChildKeyPresent := newSequenceEdges["child"]; !isChildKeyPresent {
						newSequenceEdges["child"] = []int{previousSequenceId}
					} else {
						newSequenceEdges["child"] = append(newSequenceEdges["child"], previousSequenceId)
					}
				}
			}
		}
	}

	// update operationNameToNodes
	// add nodes to sequences and connect them with node ids found in operationNameToNodes
	// find the sequences the nodes that match with the input are part of
	// shorter sequences are above longer sequences
	// sequences of same length are lower than the sequence of items they have in common in relative order
	// the lower sequence doesn't match all its items with the higher sequence.
	// lower sequence was entered before the higher sequences and
	// only has 1 match with n different higher sequences and only 1 of those
	// sequences has been entered
	// count number of nodes in linked list

}
func Pattern() {

	item1 := Variables{State: map[string]interface{}{x: 0, y: 0, z: 0},
		IfConditionResult: true}
	if !R3Test(&item1) {
		return
	}
	caretaker1 := Caretaker{}
	itemSequence1 := []string{
		mF1UY,
		mF1UY,
		mB1UX,
		mB1UX,
		mB1UY,
		mB1UY,
		mF1UX,
		mF1UX,
		mF1UZ,
		mF1UZ}
	// nodes := []*Node1{}
	sh := SequenceHierarchy{Sequences: &[]*Node1{}, FunctionNameToNodeIds: map[string]*map[int]int{
		mF1UX: nil,
		mF1UY: nil,
		mF1UZ: nil,
		mB1UX: nil,
		mB1UY: nil,
		mB1UZ: nil,
	}}
	sh.CreateSequenceOfOperationChangeNames(&item1, &caretaker1, itemSequence1)
	// fmt.Printf("here\n")
	// for _, item := range nodes {
	// 	fmt.Printf("%v\n", item)
	// }

	// fmt.Printf("\n\n")

	item2 := Variables{State: map[string]interface{}{x: 0, y: 0, z: 0},
		IfConditionResult: true}

	caretaker2 := Caretaker{}
	// mF1UY, mB1UX, mB1UY, mF1UX, mF1UZ
	itemSequence2 := []string{mB1UZ, mB1UX, mF1UY, mB1UX, mF1UX, mB1UY}
	sh.CreateSequenceOfOperationChangeNames(&item2, &caretaker2, itemSequence2)
	// fmt.Printf("\n\n")

	for _, item := range *sh.Sequences {
		fmt.Printf("%v\n", item)
	}
	fmt.Printf("\n")

	for operationName, item := range sh.FunctionNameToNodeIds {
		fmt.Printf("%v: %v\n", operationName, item)
	}
	fmt.Printf("%v\n", sh.FirstNodeIdLastSequenceAdded)
	for nodeIds, _ := range sh.NodeIdsLastSequenceAdded {
		fmt.Printf("%v, ", nodeIds)
	}
	fmt.Printf("\n")

	sh.Categorize()
	for _, item := range *sh.Sequences {
		fmt.Printf("%v\n", item)
	}
	fmt.Printf("\n")

	item3 := Variables{State: map[string]interface{}{x: 0, y: 0, z: 0},
		IfConditionResult: true}

	caretaker3 := Caretaker{}
	itemSequence3 := []string{mF1UX}
	sh.CreateSequenceOfOperationChangeNames(&item3, &caretaker3, itemSequence3)

	for _, item := range *sh.Sequences {
		fmt.Printf("%v\n", item)
	}
	fmt.Printf("\n")

	for operationName, item := range sh.FunctionNameToNodeIds {
		fmt.Printf("%v: %v\n", operationName, item)
	}
	fmt.Printf("%v\n", sh.FirstNodeIdLastSequenceAdded)
	for nodeIds, _ := range sh.NodeIdsLastSequenceAdded {
		fmt.Printf("%v, ", nodeIds)
	}
	fmt.Printf("\n")

	sh.Categorize()
	for _, item := range *sh.Sequences {
		fmt.Printf("%v\n", item)
	}
	fmt.Printf("\n")

	// checkFunctions := map[int][]string{}
	// for _, item := range itemSequence1 {
	// 	// fmt.Printf("%v. %v, %v\n", item, "check", strings.Contains(item, "check"))
	// 	// fmt.Printf("item1 %v\n", item1)

	// 	if strings.Contains(item, "check") {
	// 		// fmt.Printf("%v, %v\n", item, i)
	// 		// functions[item].(func(v *Variables, c *Caretaker))(&item1, &caretaker1)
	// 		// fmt.Printf("%v, %v\n", item, functions[item].(func(v *Variables, c *Caretaker) bool)(&item1, &caretaker1))
	// 		if !functions[item].(func(v *Variables, c *Caretaker) bool)(&item1, &caretaker1) {
	// 			continue
	// 		}
	// 		if entry := checkFunctions[0]; len(entry) >= 1 {
	// 			checkFunctions[0] = append(checkFunctions[0], item)
	// 		} else {
	// 			checkFunctions[0] = []string{item}
	// 		}
	// 	} else {
	// 		functions[item].(func(v *Variables, c *Caretaker))(&item1, &caretaker1)
	// 	}

	// }
	// fmt.Printf("%v", checkFunctions)
	// myBlocks := Blocks{Blocks: map[string]Block{}, MaxInt: 0}

	// myBlocks.Blocks["leftY"] = Block{Id: "leftY", FunctionName: "leftY"}
	// myBlocks.Blocks["forward"] = Block{Id: "forward", FunctionName: "forward"}
	// myBlocks.Blocks["checkLeftX"] = Block{Id: "checkLeftX", FunctionName: "checkLeftX"}
	// myBlocks.Blocks["path"] = Block{Id: "path",
	// 	NestedBlock: map[string]Block{
	// 		"0": {Id: "0",
	// 			Variables: map[string]Variable{
	// 				"x": {Value: Atom{IntValue: 0, TypeValueSet: "IntValue"}},
	// 				"y": {Value: Atom{IntValue: 0, TypeValueSet: "IntValue"}},
	// 				"z": {Value: Atom{IntValue: 0, TypeValueSet: "IntValue"}},
	// 			},
	// 			Sequence: []Link{
	// 				Link{Ids: []string{"forward"}},
	// 				Link{Ids: []string{"checkLeftX"}}},
	// 		}}}
	// inputs := map[string][]string{
	// 	"leftY":      []string{"leftY"},
	// 	"forward":    []string{"forward"},
	// 	"checkLeftX": []string{"checkLeftX"},
	// }
	// functionNameFunction := map[string]func(blocks Blocks, path []string, sequencePos int) bool{
	// "leftY":      leftY,
	// "checkLeftX": checkLeftX,
	// }
	// myBlocks.Blocks["cond"] = Block{Id: "cond",
	// 	NestedBlock: map[string]Block{
	// 		"instances": {Id: "instances", Sequence: LinkedList{}},
	// 		"0":         {Id: "0", FunctionName: "condFunction"}}}
	// myBlocks.Blocks["if"] = Block{Id: "if",
	// 	NestedBlock: map[string]Block{
	// 		"instances": {Id: "instances", Sequence: LinkedList{}},
	// 		"0":         {Id: "0", Sequence: LinkedList{LinkedList: []LinkedNode{{Data: Link{Ids: []string{"cond", "0"}}}}, FirstNode: 0, LastNode: 0, CurrentNode: 0}}}}
	// sequence of blocks for different directions
	// all spirals have to be larger than 1 unit spiral

	// detect repeating
	// 1 small spiral
	// 1 large spiral
	// detect repeating for small spiral
	// detect parts of small spiral as part of large spiral
	// update small spiral to have parts of large spiral
	// 1 small large spiral
	// detect part of small large spiral using current spiral template
	// 1 wierd spiral (70% or less match with spiral detector)
	// detect spiral parts and generate spiral using the spiral parts it has detected
	// simplify saved sequenes by deleting nodes that don't match to spiral traits
	// remove example spiral sequences so there is idealy 1 unit spiral to detect all future spirals

}

// matching
// needs to match different sequences at different times per sequence
// ith input != ith position is existing pattern
