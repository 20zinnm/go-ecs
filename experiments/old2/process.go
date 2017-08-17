package old2

type System interface {
	Execute(Entity)
}

type system struct {
	sys     System
	matcher Matcher
}
