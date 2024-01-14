package Patterns

func move1Unit(v *Variables, c *Caretaker, dimensionName string, direction func(int) int) {

	c.UpdateMemento(v.StructInstanceName, v.CreateMemento())

	if _, ok := v.State[dimensionName].(int); !ok {
		return
	}
	dimension := v.State[dimensionName].(int)
	dimension = direction(dimension)
	v.State[dimensionName] = dimension

}
func moveForward1UnitX(v *Variables, c *Caretaker)  { move1Unit(v, c, x, add) }
func moveForward1UnitY(v *Variables, c *Caretaker)  { move1Unit(v, c, y, add) }
func moveForward1UnitZ(v *Variables, c *Caretaker)  { move1Unit(v, c, z, add) }
func moveBackward1UnitX(v *Variables, c *Caretaker) { move1Unit(v, c, x, subtract) }
func moveBackward1UnitY(v *Variables, c *Caretaker) { move1Unit(v, c, y, subtract) }
func moveBackward1UnitZ(v *Variables, c *Caretaker) { move1Unit(v, c, z, subtract) }
