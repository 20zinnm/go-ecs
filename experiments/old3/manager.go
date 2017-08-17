package old3

import "sync"

type manager struct {
	components [][]interface{}
	// entities stores an array for each ComponentType such that entities[type]_id represents whether an entity with a given id has a given component.
	// This is used to calculate offsets for entities during iterations.
	componentIndices [][]bool
	maxId            uint
	freeIds          []uint
	mut    sync.RWMutex
}
