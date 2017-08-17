package ecs

import (
	"reflect"
)

type process struct {
	fn   reflect.Value
	name string
	//deps []string
	past []reflect.Type
	fut  []reflect.Type
}


//func resolveDependencies(procs []*process) []*process {
//	procNames := make(map[string]*process)
//	procDeps := make(map[string]mapset.Set)
//
//	for _, p := range procs {
//		procNames[p.name] = p
//		dependencySet := mapset.NewSet()
//		for _, dep := range p.past {
//			dependencySet.Add(dep)
//		}
//		procDeps[p.name] = dependencySet
//	}
//	var resolved []*process
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
//			panic("ecs: circular dependency")
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
//	return resolved
//}

//type Option func(proc *process)

//func DependsOn(pname string) Option {
//	return func(proc *process) {
//		proc.deps = append(proc.deps, pname)
//	}
//}

//func Name(name string) Option {
//	return func(proc *process) {
//		proc.name = name
//	}
//}
