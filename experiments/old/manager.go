package old

import (
	"reflect"
	"sync"
	"context"
)

const (
	CtxEntityId int = iota
)

func NewManager() *EntityManager {
	return &EntityManager{
		componentsP:      make([][]Component, 0),
		componentCountsP: make([][]uint, 0),
		componentsF:      make([][]Component, 0),
		componentCountsF: make([][]uint, 0),
		procs:            make([]Process, 0),
		emut:             sync.RWMutex{},
	}
}

type EntityManager struct {
	past *gameState

	entityCount      uint
	procs            []Process
	emut             sync.RWMutex
}

func (e *EntityManager) ExecuteFrame(ctx context.Context) {
	e.emut.Lock()
	defer e.emut.Unlock()
}

func (e *EntityManager) execProcess(proc Process) {
	pastIndices := make([]uint, len(proc.pastComponents))
	for i := 0; i < len(pastIndices); i++ {
		pastIndices[i] = 0
	}
	futureIndex := 0
	for eid := uint(0); eid < e.entityCount; eid++ {
		satisfies := true
		for n := range proc.pastComponents {
			if e.componentCountsP[n][eid] == 0 {
				satisfies = false
				break
			}
		}
		if !satisfies {
			break
		}
		for n := range proc.pastComponents {
			if n > 0 {
				pastIndices[n] += e.componentCountsP[proc.pastComponents[n-1]][eid]
			} else {
				pastIndices[0] += e.componentCountsP[proc.pastComponents[0]][eid]
			}
		}
		if proc.ctx {
			ctx := context.WithValue(context.Background(), CtxEntityId, eid)

		}
	}
}

func (e *EntityManager) AddProcess(proc interface{}) {
	e.emut.Lock()
	defer e.emut.Unlock()

	var pr Process

	// Input
	vf := reflect.ValueOf(proc)
	if vf.Kind() != reflect.Func {
		panic("AddProcess requires a function argument.")
	}
	if vf.Type().NumIn() < 1 {
		panic("malformed Process")
	}
	ps := 0
	if vf.Type().In(0) == reflect.TypeOf(context.Background()) {
		pr.ctx = true
		ps = 1
	}
	for i := ps; i < vf.Type().NumIn(); i++ {
		n := reflect.New(vf.Type().In(i)).Elem().Interface()
		nc, ok := n.(Component)
		if !ok {
			panic("malformed Process")
		}
		pr.pastComponents = append(pr.pastComponents, nc.(Component).Type())
	}

	// Output
	switch vf.Type().NumOut() {
	case 0:
		pr.futureComponent = nil
	case 1:
		n := reflect.New(vf.Type().Out(0)).Elem().Interface()
		nc, ok := n.(Component)
		if !ok {
			panic("malformed Process")
		}
		pr.futureComponent = &nc.(Component).Type()
	default:
		panic("malformed Process")
	}

	pr.proc = vf

	// Add to procs list
	e.procs = append(e.procs, pr)
}
