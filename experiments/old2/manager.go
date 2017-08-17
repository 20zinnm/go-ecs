package old2

import "sync"

type Manager struct {
	past    *gameState
	future  *gameState
	mut     sync.RWMutex
	systems []system
}

func (m *Manager) NewEntity() EntityId {

}

func (m *Manager) RegisterSystem(matcher Matcher, sys System) {
	m.mut.Lock()
	m.systems = append(m.systems, system{sys, matcher})
	m.mut.Unlock()
}

func (m *Manager) Tick() {
	m.mut.Lock()
	defer m.mut.Unlock()

	for i, sys := range m.systems {
		newFuture := m.past
		newPast := m.future


	}
}
