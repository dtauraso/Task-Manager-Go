package Patterns

import (
	"fmt"
)

type Node1 struct {
	Id             int
	SequenceLength int
	VariableName   string
	FunctionName   string
	TypeName       string
	Edges          map[string][]int
	ParentChildId  int
}

func (n1 *Node1) GetLastEdge(edgeName string) int {
	return n1.Edges[edgeName][len(n1.Edges[edgeName])-1]
}

const (
	mF1UX = "moveForward1UnitX"
	mF1UY = "moveForward1UnitY"
	mF1UZ = "moveForward1UnitZ"
	mB1UX = "moveBackward1UnitX"
	mB1UY = "moveBackward1UnitY"
	mB1UZ = "moveBackward1UnitZ"
	cAC   = "checkAddChange"
	cSC   = "checkSubtractChange"
)

var functions = map[string]interface{}{
	mF1UX: moveForward1UnitX,
	mF1UY: moveForward1UnitY,
	mF1UZ: moveForward1UnitZ,
	mB1UX: moveBackward1UnitX,
	mB1UY: moveBackward1UnitY,
	mB1UZ: moveBackward1UnitZ,
}
var functionNameMapCheckFunctionName = map[string]string{
	mF1UX: cAC,
	mF1UY: cAC,
	mF1UZ: cAC,
	mB1UX: cSC,
	mB1UY: cSC,
	mB1UZ: cSC,
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
	Sequences              *[]*Node1
	FunctionNameToNodeIds  map[string]*map[int]int
	FunctionNameToNodeIds2 map[string]*map[int]struct{}

	FirstNodeIdLastSequenceAdded int
	NodeIdsLastSequenceAdded     map[int]struct{}
}

func (sh *SequenceHierarchy) CreateSequenceOfCheckFunctionNames(
	v *Variables,
	c *Caretaker,
	sequence []string) *[]*Node1 {

	head := -1
	prev := head
	lastOperationName := ""
	Sequence := &[]*Node1{}
	sh.FirstNodeIdLastSequenceAdded = len(*sh.Sequences)
	for _, functionName := range sequence {
		functions[functionName].(func(v *Variables, c *Caretaker))(v, c)
		if functionName != lastOperationName {

			changedVariableName := ""
			typeName := ""
			// likely to be O(1) due to each operation only changing 1 variable at a time
			for variableName, value := range v.State {
				prevValue := c.GetLastMemento(v.StructInstanceName).State[variableName]
				if value != prevValue {
					changedVariableName = variableName
					typeName = fmt.Sprintf("%T", value)
				}
			}

			newNodeId := len(*Sequence)

			temp := Node1{
				Id:             newNodeId,
				VariableName:   changedVariableName,
				SequenceLength: len(sequence),
				FunctionName:   functionNameMapCheckFunctionName[functionName],
				TypeName:       typeName,
				Edges:          map[string][]int{"prev": {prev}, "next": {-1}}}
			if prev >= 0 {
				newEdges := (*Sequence)[prev].Edges
				newEdges["next"] = []int{temp.Id}
				(*Sequence)[prev].Edges = newEdges
			}
			*Sequence = append(*Sequence, &temp)
			prev = temp.Id
		}
		lastOperationName = functionName
	}

	return Sequence
}

func (sh *SequenceHierarchy) CreateSequenceOfOperationChangeNames(
	v *Variables,
	c *Caretaker,
	sequence []string) {

	// when the command changes
	// note what variable values changed
	// record the changes as a sequence of operation change names
	// insert comparison function that verifies the operation change happened.
	// record check functions
	// need to keep sequence separate till know how it connects to saved sequences

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
				prevValue := c.GetLastMemento(v.StructInstanceName).State[variableName]
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

func (sh *SequenceHierarchy) Categorize2(newSequence *[]*Node1) {

	// 1 to many matches from newSequence to sh.Sequences
	// 1 to many matches where 1 item from newSequence matches with many different sequences
	// long sequence(higher) that uses shorter sequences(lower) many times
	// meauring complexity
	// measuring specificy
	// measuring reusability
	// patterns as a hierarchy of sequences with limited lengths
	// using different levels in hierarchy that represent new patterns when connected as a new pattern
	// recognize the pattern by building a hierarchy
	// saving changing sequence of patterns of length 1 is saving change of changes.
	// build pattern tree with pattern of arbitrary length using sequences of length 2
	// pattern tree needs to be height balanced
	// reading the input sequence, replacing subsequences with parent check function name when they match partitions of length 2
	// if the match sequence exists in sh.Sequence
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

type Node struct {
	Id                  int
	Edges               map[string][]int
	IsActive            bool
	ActiveChildrenCount int
}

var nodes []*Node

func updateBottom(Bottom *map[string][]int, item rune, newNodeId int) *map[string][]int {
	bottomEdges := (*Bottom)[string(item)]
	bottomEdges = append(bottomEdges, newNodeId)
	(*Bottom)[string(item)] = bottomEdges
	return Bottom
}
func addToDoublyLinkedList(nodes *[]*Node, prev, newNodeId, parentNodeId int) (tempId int) {
	temp := Node{
		Id:    newNodeId,
		Edges: map[string][]int{"prev": {prev}, "next": {-1}, "parents": {parentNodeId}}}
	if prev >= 0 {
		newEdges := (*nodes)[prev].Edges
		newEdges["next"] = []int{temp.Id}
		(*nodes)[prev].Edges = newEdges
	}
	*nodes = append(*nodes, &temp)
	return temp.Id
}
func makeDoublyLinkedListForArrays(nodes *[]*Node, Bottom *map[string][]int, childIds []int, prev int, item interface{}) (tempId int) {
	parentNodeId := doublyLinkSequence(nodes, Bottom, item)
	newNodeId := len(*nodes)
	childIds = append(childIds, newNodeId)
	tempId = addToDoublyLinkedList(nodes, prev, newNodeId, parentNodeId)
	return tempId
}
func doublyLinkSequence(nodes *[]*Node, Bottom *map[string][]int, sequence interface{}) (parentNodeId int) {

	head := -1
	prev := head
	childIds := []int{}
	fmt.Printf("%s\n", fmt.Sprintf("%T", sequence))
	// tree entries off slightly
	// make specific parent node for each different sequence visiting same subsequence
	// make 1 node connecting to the specific parent nodes for sequences with variation that are the same
	// at 1 point in time
	// change where the read write heads are (repeat visitting same items
	// write down the autocompleted sequence)
	// autocomplete stops when sequence finishes last node
	// bias in favor of patterns
	// measure connection strength
	// stronger connections dictate the order
	switch sequence.(type) {
	case []interface{}:
		parentNodeId = len(*nodes) + len(sequence.([]interface{}))

		for _, item := range sequence.([]interface{}) {
			parentNodeId := doublyLinkSequence(nodes, Bottom, item)
			newNodeId := len(*nodes)
			childIds = append(childIds, newNodeId)
			tempId := addToDoublyLinkedList(nodes, prev, newNodeId, parentNodeId)
			// tempId := makeDoublyLinkedListForArrays(nodes, Bottom, childIds, prev, item)
			prev = tempId
		}
	case string:
		parentNodeId = len(*nodes) + len(sequence.(string))
		for _, item := range sequence.(string) {
			newNodeId := len(*nodes)
			childIds = append(childIds, newNodeId)
			Bottom = updateBottom(Bottom, item, newNodeId)
			tempId := addToDoublyLinkedList(nodes, prev, newNodeId, parentNodeId)
			prev = tempId
		}
	default:
		return -1
	}
	*nodes = append(*nodes, &Node{Id: parentNodeId, Edges: map[string][]int{"children": childIds}})
	return parentNodeId
}

// doubly link new nodes
// if the same sequence of length > 1 is revsited from a different starting point
//		replace sequence with 1 node for each starting point
// increase occurrence count for new nodes(1) and revisiting nodes(x + 1)
// if first and last node have a frequency ratio too far apart
//		make sorted list of nodes from sequence prev and current in assending order using occurrence count
//		remove nodes on the list where frquency ratio is out of bounds
//		take remaining nodes on the list and make 1 parent node they all connect it (typical sequence)
// if parent of nodes exists
// 		activated nodes move to parent
// 		1 goroutine for dft 1 level up from activated nodes (autocomplete)
// 			dft the next link unless next has no links
// 			if next has n links
//				there is no if statement as 1 of the options is already predetermined to be run by the begining of the sequence
//  			1 goroutine per link
// 		repeat increase ocurrence count step

func Hierarchy() {

	// _ = doublyLinkSequence(&nodes, &Bottom, "title")
	// _ = doublyLinkSequence(&nodes, &Bottom, "tag")
	// _ = doublyLinkSequence(&nodes, &Bottom, []interface{}{[]interface{}{"title", "tag"}})
	// make sure the typical sequence happens more often than the random data
	_ = doublyLinkSequence(&nodes, &Bottom, []interface{}{
		[]interface{}{"x", "t", "i"},
		[]interface{}{"y", "t", "i"},
		[]interface{}{"t", "i", "x"}})

	/*
		typical sequence construction
			matching items in sequence get stronger connections
			non matching items in sequence get weaker connections
			sort by increasing frquency
			sequence can be extended while frequencey ratio from least frequent item to the most frequent item is < 10%
			when item at pos 0 is >= 10% weaker than the item at pos len(sequence) - 1, then item at pos 0 is removed
			remaining items winn when pos 0 is < 10% weaker than pos len(sequence) - 1
			the typical sequence strengths become permanent
	*/
	// get a first match wth input
	// make list of candidates for possible match
	// first candidate sequence wins
	text0 := "title" // "tag"

	winningNodeIds := []int{}
	bottomTrackers := map[int]struct{}{}
	for i := 0; i < len(text0); i++ {
		letter := string(text0[i])
		fmt.Printf("%s\n", letter)

		tempBottomTrackers1 := map[int]struct{}{}

		parents := map[int]struct{}{}
		// assumes leftmost index of sequence in Bottom[letter] is the start of the sequence
		for _, nodeId := range Bottom[letter] {
			if nodeId == -1 {
				continue
			}
			if nodes[nodeId].IsActive {
				continue
			}

			// need to make sure nodeId's aren't part of the same sequence
			if _, ok := parents[nodes[nodeId].Edges["parents"][0]]; !ok {
				parents[nodes[nodeId].Edges["parents"][0]] = struct{}{}
				tempBottomTrackers1[nodeId] = struct{}{}
			}

		}
		fmt.Printf("parents %v\n", parents)
		bottomTrackers = tempBottomTrackers1

		tempBottomTrackers := map[int]struct{}{}
		fmt.Printf("bottom trackers loop start %v\n", bottomTrackers)
		for nodeId := range bottomTrackers {

			if nodeId == -1 {
				continue
			}
			nodes[nodeId].IsActive = true

			tempBottomTrackers[nodes[nodeId].Edges["next"][0]] = struct{}{}
			parentId := nodes[nodeId].Edges["parents"][0]

			nodes[parentId].ActiveChildrenCount += 1
			if nodes[parentId].ActiveChildrenCount == len(nodes[parentId].Edges["children"]) {
				nodes[parentId].IsActive = true
				winningNodeIds = append(winningNodeIds, parentId)
			}

		}
		bottomTrackers = tempBottomTrackers
		fmt.Printf("bottom trackers loop end %v\n", bottomTrackers)
		for _, nodeRef := range nodes {
			fmt.Printf("Node %v\n", *nodeRef)
		}

	}
	// todo: travel and activate nodes for the remainder of the sequence
	// stop activating nodes if the next step has more than 1 possibility
	// autocomplete is depth first traveral
	// autocomplete stops at tracker n if tracker n is on a "user input" node
	fmt.Printf("%v\n", winningNodeIds)
	for _, nodeRef := range nodes {
		fmt.Printf("Node %v\n", *nodeRef)
	}
	fmt.Printf("\n")

}

func Pattern() {

	item1 := Variables{State: map[string]interface{}{x: 0, y: 0, z: 0},
		StructInstanceName: "item1"}

	caretaker := Caretaker{}
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
	Sequence := sh.CreateSequenceOfCheckFunctionNames(&item1, &caretaker, itemSequence1)
	fmt.Printf("check function sequence\n")
	for _, item := range *Sequence {
		fmt.Printf("%v\n", item)
	}
	fmt.Printf("\n")
	return
	sh.CreateSequenceOfOperationChangeNames(&item1, &caretaker, itemSequence1)
	// fmt.Printf("here\n")
	// for _, item := range nodes {
	// 	fmt.Printf("%v\n", item)
	// }

	// fmt.Printf("\n\n")

	item2 := Variables{State: map[string]interface{}{x: 0, y: 0, z: 0},
		StructInstanceName: "item2"}

	// mF1UY, mB1UX, mB1UY, mF1UX, mF1UZ
	itemSequence2 := []string{mB1UZ, mB1UX, mB1UX, mF1UY, mB1UX, mF1UX, mB1UY}
	sh.CreateSequenceOfOperationChangeNames(&item2, &caretaker, itemSequence2)
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
		StructInstanceName: "item3"}

	itemSequence3 := []string{mF1UX, mF1UX, mF1UX}
	sh.CreateSequenceOfOperationChangeNames(&item3, &caretaker, itemSequence3)

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

	fmt.Printf("%v\n", caretaker.memento)
	// fmt.Printf("%v\n", caretaker)
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
