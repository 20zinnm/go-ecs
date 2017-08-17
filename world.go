package ecs

import (
	"reflect"
	"sync"
	"github.com/deckarep/golang-set"
)

type World interface {
	ExecuteFrame()
	NewEntity(components ...interface{}) uint
	AddProcess(fn interface{})
}
type gamestate struct {
	components      map[reflect.Type]reflect.Value
	componentCounts map[reflect.Type][]uint
}

func newGamestate() *gamestate {
	return &gamestate{
		components:      make(map[reflect.Type]reflect.Value),
		componentCounts: make(map[reflect.Type][]uint),
	}
}

func NewWorld() World {
	return &world{
		past:        newGamestate(),
		future:      newGamestate(),
		entityCount: 0,
		free:        nil,
		procs:       nil,
		writtenTo:   mapset.NewThreadUnsafeSet(),
		wg:          &sync.WaitGroup{},
		Mutex:       &sync.Mutex{},
	}
}

type world struct {
	past        *gamestate
	future      *gamestate
	entityCount uint
	free        []uint
	procs       []*process
	writtenTo   mapset.Set
	wg          *sync.WaitGroup
	*sync.Mutex
}

func (w *world) ExecuteFrame() {
	w.Lock()
	defer w.Unlock()

	*w.past = *w.future
	for _, proc := range w.procs {
		w.wg.Add(1)
		go func(proc *process) {
			w.executeProcess(proc);
			w.wg.Done()
		}(proc)
	}
	w.wg.Wait()
}

func (w *world) AddProcess(fn interface{}) {
	w.Lock()
	defer w.Unlock()

	var proc process

	// Input--past components process requires
	vf := reflect.ValueOf(fn)
	if vf.Kind() != reflect.Func {
		panic("ecs: AddProcess called with an argument that is not a function")
	}
	proc.fn = vf

	for i := 0; i < vf.Type().NumIn(); i++ {
		proc.past = append(proc.past, vf.Type().In(i))
	}

	// Output--future components process writes to
	for i := 0; i < vf.Type().NumOut(); i++ {
		proc.fut = append(proc.fut, vf.Type().Out(i))
		added := w.writtenTo.Add(vf.Type().Out(i))
		if !added {
			panic("ecs: multiple processes will write to the same future component")
		}
	}

	// Add to procs list
	w.procs = append(w.procs, &proc)
	//w.procs = resolveDependencies(w.procs)
}

func (w *world) executeProcess(proc *process) {
	pastIndices := make([]uint, len(proc.past))
	futureIndices := make([]uint, len(proc.fut))
	for id := uint(0); id < w.entityCount; id++ {
		matches := true
		for j, typ := range proc.past {
			cc := w.past.componentCounts[typ]
			if id == 0 {
				pastIndices[j] = 0
			} else {
				pastIndices[j] += cc[id-1]
			}
			if cc[id] <= 0 {
				matches = false
			}
		}
		if !matches {
			continue
		}
		past := make([]reflect.Value, len(proc.past))
		for j, typ := range proc.past {
			past[j] = w.past.components[typ].Index(int(pastIndices[j]))
		}
		res := proc.fn.Call(past)
		for j, val := range res {
			if id == 0 {
				futureIndices[j] = 0
			} else {
				futureIndices[j]++
			}
			w.future.components[proc.fut[j]].Index(int(futureIndices[j])).Set(val)
		}
	}
}

func (w *world) NewEntity(components ...interface{}) uint {
	w.Lock()
	defer w.Unlock()
	id := uint(0)
	if len(w.free) > 0 {
		id, w.free = w.free[0], w.free[1:]
	} else {
		id = w.entityCount
		w.entityCount++
	}
	for _, c := range components {
		typ := reflect.TypeOf(c)

		ccs, ok := w.future.componentCounts[typ]
		if !ok {
			w.future.componentCounts[typ] = make([]uint, w.entityCount)
		}
		if int(id) >= len(ccs) {
			w.future.componentCounts[typ] = append(ccs, make([]uint, int(w.entityCount)-len(ccs))...)
		}
		w.future.componentCounts[typ][id] = 1
		j := 0
		for i := 0; i < int(id); i++ {
			j += int(w.future.componentCounts[typ][i])
		}
		if _, ok := w.future.components[typ]; !ok {
			w.future.components[typ] = reflect.MakeSlice(reflect.SliceOf(typ), 1, 8)
		}
		// s = append(s, 0)
		n := reflect.New(typ)
		w.future.components[typ] = reflect.Append(w.future.components[typ], n.Elem())
		// copy(s[i+1:], s[i:])
		length := w.future.components[typ].Len()
		reflect.Copy(w.future.components[typ].Slice(j+1, length), w.future.components[typ].Slice(j, length))
		// s[i] = x
		w.future.components[typ].Index(j).Set(reflect.ValueOf(c))
	}
	return id
}
