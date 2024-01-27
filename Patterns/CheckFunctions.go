package Patterns

import (
	"reflect"
)

func checkAddSubtractChange(v *Variables, c *Caretaker, arithmaticFunction func(int) int) bool {

	// assume 1 variable changed by +1 or -1

	changedVariableName := ""

	for variableName, value := range v.State {
		prevValue := c.GetLastMemento(v.StructInstanceName).State[variableName]
		if value != prevValue {
			changedVariableName = variableName
		}
	}
	if _, ok := v.State[changedVariableName].(int); !ok {
		return false
	}
	var variablePrev int
	memento := c.GetLastMemento(changedVariableName)
	if reflect.DeepEqual(memento, Memento{}) {
		variablePrev = 0
	} else {
		variablePrev = memento.State[changedVariableName].(int)
	}
	variable := v.State[changedVariableName].(int)
	return variable == arithmaticFunction(variablePrev)

}
func checkAddChange(v *Variables, c *Caretaker) bool {
	return checkAddSubtractChange(v, c, add)
}
func checkSubtractChange(v *Variables, c *Caretaker) bool {
	return checkAddSubtractChange(v, c, subtract)
}
