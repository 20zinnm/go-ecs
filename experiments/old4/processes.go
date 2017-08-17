package old4

import (
	"reflect"
)

type process struct {
	fn   Process
	name string
	//deps []string
	fut  *reflect.Type
	reqP map[reflect.Type]bool
}

type processes []*process

//func (processes processes) resolveDependencies() {
//	procNames := make(map[string]*process)
//	procDeps := make(map[string]mapset.Set)
//
//	for _, p := range processes {
//		procNames[p.name] = p
//		dependencySet := mapset.NewSet()
//		for _, dep := range p.deps {
//			dependencySet.Add(dep)
//		}
//		procDeps[p.name] = dependencySet
//	}
//	var resolved processes
//	for len(procDeps) != 0 {
//		// Get all nodes from the graph which have no dependencies
//		readySet := mapset.NewSet()
//		for name, deps := range procDeps {
//			if deps.Cardinality() == 0 {
//				readySet.Add(name)
//			}
//		}
//
//		// If there aren't any ready nodes, then we have a cicular dependency
//		if readySet.Cardinality() == 0 {
//			panic("circular dependency")
//		}
//
//		// Remove the ready nodes and add them to the resolved graph
//		for name := range readySet.Iter() {
//			delete(procDeps, name.(string))
//			resolved = append(resolved, procNames[name.(string)])
//		}
//
//		// Also make sure to remove the ready nodes from the
//		// remaining node dependencies as well
//		for name, deps := range procDeps {
//			diff := deps.Difference(readySet)
//			procDeps[name] = diff
//		}
//	}
//	processes = resolved
//}
