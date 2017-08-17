package old2

type EntityId uint

type Entity interface {
	Id() EntityId
	Has(componentType ComponentType) bool
	Get(componentType ComponentType) interface{}
	Set(componentType ComponentType, component interface{})
	Del(componentType ComponentType)
	// Kill deletes all component data for an entity. Note, however, that it does not remove the entity ID; thus, it is not sufficient to determine whether an entity is alive. This should be a separate component.
	Kill()
}

type entity struct {
	id      int
	// version corresponds with the last gametick this entity saw
	version int
	manager *Manager
}

func (e *entity) Has(componentType ComponentType) bool {
	e.manager.mut.RLock()
	defer e.manager.mut.RUnlock()

	return e.manager[componentType][e.id] > 0
}

func (e *entity) Get(componentType ComponentType) interface{} {
	e.manager.mut.RLock()
	defer e.manager.mut.RUnlock()

	return *e.components[componentType]
}

func (e *entity) Set(componentType ComponentType, component interface{}) {
	e.manager.mut.Lock()
	had := e.manager.componentCountsF[componentType][e.id] > 0
	if !had {
		e.manager.componentCountsF[componentType][e.id] = 1
		e.manager.componentsF = append(e.manager.componentsF, nil)
		i := uint(0)
		for j := 0; j < e.id; j++ {
			i += e.manager.componentCountsF[componentType][i]
		}
		copy(e.manager.componentsF[componentType][i+1:], e.manager.componentsF[componentType][i:])
		e.manager.componentsF[componentType][i] = component
	} else {
		i := uint(0)
		for j := 0; j < e.id; j++ {
			i += e.manager.componentCountsF[componentType][i]
		}
		e.manager.componentsF[componentType][i] = component
	}
	e.manager.mut.Unlock()
}

func (e *entity) Del(componentType ComponentType) {
	e.manager.mut.Lock()
	e.manager.componentCountsF[componentType][e.id] = 0
	i := uint(0)
	for j := 0; j < e.id; j++ {
		i += e.manager.componentCountsF[componentType][i]
	}
	e.manager.componentsF[componentType] = append(e.manager.componentsF[componentType][:i], e.manager.componentsF[componentType][i+1:]...)
	e.manager.mut.Unlock()
}
