package control

type Subject interface {
	ControlType() string

	Control(cmds []interface{})
}
