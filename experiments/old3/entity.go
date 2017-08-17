package old3

type EntityId uint

type Entity interface {
	Id() EntityId
	// Has returns whether an Entity has a component of the given type in the past gamestate.
	Has(componentType ComponentType) bool
	// Get returns the instance of a component related to an entity in the past gamestate.
	Get(componentType ComponentType) interface{}
	Del(componentType ComponentType)
	// Set updates the instance of a component related to an entity in the future gamestate.
	// The component that is set is determined when the system is initially registered. If this is called on a system that does not alter the future state, the engine will panic (programmer error).
	// This may alter which entities a System gets in the next iteration.
	Set(component interface{})
	// Kill deletes all component data for an entity in the future gamestate and frees the ID for future allocation.
	Kill()
}

type entity struct {
	m           *manager
	id          EntityId
	pastIndices []uint
	futureComponents []bool
}

func (e *entity) Has(componentType ComponentType) bool {
	return e.m.componentIndices[componentType][e.id]
}

func (e *entity) Get(componentType ComponentType) interface{} {
	return e.m.components[componentType][e.pastIndices[componentType]]
}

func (e *entity) Set(component interface{}) {

}

func (e *entity) Id() EntityId {
	return e.id
}

//
//func (e *entity) Has(componentType ComponentType) bool {
//	e.mut.RLock()
//	defer e.mut.RUnlock()
//	return e.world.componentIndices[componentType][e.id]
//}
//
//func (e *entity) Get(componentType ComponentType) interface{} {
//	e.mut.RLock()
//	defer e.mut.RUnlock()
//	i := 0
//	for j := 0; j < int(e.id); j++ {
//		if e.componentIndices[componentType][j] {
//			i++
//		}
//	}
//	return e.components[componentType][i]
//}
//
//func (e *entity) Set(componentType ComponentType, component interface{}) {
//	panic("implement me")
//}
//
//func (e *entity) Del(componentType ComponentType) {
//	panic("implement me")
//}

func (e *entity) Kill() {
	panic("implement me")
}
