package old2

type Matcher interface {
	Matches(components ...ComponentType) bool
}