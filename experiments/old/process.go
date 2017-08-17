package old

import "reflect"

type Process struct {
	ctx             bool
	futureComponent *ComponentType
	pastComponents  []ComponentType
	proc            reflect.Value
}
