package Patterns

import "reflect"

func checkAddSubtractChange(v *Variables, c *Caretaker) bool {

	if _, ok := v.State[variableName].(int); !ok {
		return false
	}
	var variablePrev int
	memento := c.GetLastMemento(variableName)
	if reflect.DeepEqual(memento, Memento{}) {
		variablePrev = 0
	} else {
		variablePrev = memento.State[variableName].(int)
	}
	variable := v.State[variableName].(int)
	return variable == arithmaticFunction(variablePrev)

}
