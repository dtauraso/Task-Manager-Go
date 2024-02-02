package Patterns

// generic object struct
type Variables struct {
	State              map[string]interface{}
	StructInstanceName string
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
	memento map[string][]Memento
}

func (c *Caretaker) InitMemento(variableName string) {
	if c.memento == nil {
		c.memento = map[string][]Memento{}
	}
	c.memento[variableName] = []Memento{}
}

func (c *Caretaker) UpdateMemento(variableName string, m Memento) {

	if _, ok := c.memento[variableName]; !ok {
		c.InitMemento(variableName)
	}
	c.memento[variableName] = append(c.memento[variableName], m)
}

func (c *Caretaker) GetLastMemento(variableName string) Memento {
	if _, ok := c.memento[variableName]; !ok {
		return Memento{}
	}
	if len(c.memento[variableName]) == 0 {
		return Memento{}

	}
	return c.memento[variableName][len(c.memento[variableName])-1]
}
