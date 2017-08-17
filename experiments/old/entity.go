package old

const (
	AliveComponentType ComponentType = 0
)

type EntityID int

type Entity struct {
	id      EntityID
	manager *EntityManager
}

func (e *Entity) Alive() bool {

}

type AliveComponent struct {
}

func (ec AliveComponent) Type() ComponentType {
	return AliveComponentType
}
