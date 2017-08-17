package old3


type System interface {
	Execute(Entity)
}

type Matcher interface {
	Matches(Entity) bool
}

type systemWrapper struct {
	matcher Matcher
	system System
	mg *manager
}

func (sw *systemWrapper) execute() {
	ec := sw.mg.maxId - len(sw.mg.freeIds)

}