package old4

import (
	"reflect"
	"sync"
	"fmt"
)

type gamestate struct {
	components      map[reflect.Type]reflect.Value
	componentCounts map[reflect.Type][]bool
	maxId           uint
	freeIds         []uint
}

func newGameState() *gamestate {
	return &gamestate{
		components:      make(map[reflect.Type]reflect.Value),
		componentCounts: make(map[reflect.Type][]bool),
		maxId:           0,
		freeIds:         make([]uint, 0),
	}
}

type World interface {
	ExecuteFrame()
	AddProcess(proc Process, options ...Option)
	NewEntity(components ...interface{}) uint
}

// a whole
func NewWorld() World {
	return &world{
		past:    newGameState(),
		future:  newGameState(),
		procs:   make(processes, 0),
		wg:      &sync.WaitGroup{},
		RWMutex: &sync.RWMutex{},
	}
}

type world struct {
	past   *gamestate
	future *gamestate
	procs  processes
	wg     *sync.WaitGroup
	*sync.RWMutex
}

func (w *world) NewEntity(components ...interface{}) uint {
	w.Lock()
	defer w.Unlock()

	id := uint(0)
	if len(w.future.freeIds) > 0 {
		id, w.future.freeIds = w.future.freeIds[0], w.future.freeIds[1:]
	} else {
		id = w.future.maxId
		w.future.maxId++
	}

	for _, c := range components {
		typ := reflect.TypeOf(c)
		ccs, ok := w.future.componentCounts[typ]
		if !ok {
			w.future.componentCounts[typ] = make([]bool, 1)
		}
		if int(id) >= len(ccs) {
			addons := make([]bool, int(id)+1-(len(ccs)))
			for i := 0; i < len(addons); i++ {
				addons[i] = false
			}
			w.future.componentCounts[typ] = append(ccs, addons...)
		}
		w.future.componentCounts[typ][id] = true
		j := 0
		for i := 0; i < int(id); i++ {
			if w.future.componentCounts[typ][i] {
				j++
			}
		}
		// s = append(s, 0)
		n := reflect.New(typ)
		//fmt.Println(n.Type())
		//fmt.Println(n.IsNil())

		if _, ok := w.future.components[typ]; !ok {
			w.future.components[typ] = reflect.MakeSlice(reflect.SliceOf(typ), 1, 1)
		}
		w.future.components[typ] = reflect.Append(w.future.components[typ], n.Elem())
		// copy(s[i+1:], s[i:])
		length := w.future.components[typ].Len()
		reflect.Copy(w.future.components[typ].Slice(j+1, length), w.future.components[typ].Slice(j, length))
		// s[i] = x
		w.future.components[typ].Index(j).Set(reflect.ValueOf(c))
		fmt.Println(w.future.components[typ].Index(j))
	}

	return id
}

func (w *world) ExecuteFrame() {
	w.Lock()
	defer w.Unlock()

	*w.past, *w.future = *w.future, *w.past

	for _, proc := range w.procs {
		w.wg.Add(1)
		go w.executeProcess(proc)
	}
	w.wg.Wait()
}

func (w *world) executeProcess(proc *process) {
	ei := &entityIterator{
		world:       w,
		proc:        proc,
		pastIndices: make(map[reflect.Type]int),
	}
	if proc.fut != nil {
		var futi uint = 0
		ei.futureIndex = &futi
	}
	proc.fn(ei)
	w.wg.Done()
}

type Option func(proc *process)

//func DependsOn(pname string) Option {
//	return func(proc *process) {
//		proc.deps = append(proc.deps, pname)
//	}
//}

func PastComponent(c interface{}) Option {
	return func(proc *process) {
		proc.reqP[reflect.TypeOf(c)] = true
	}
}

func PastComponents(cs ...interface{}) Option {
	return func(proc *process) {
		for _, c := range cs {
			proc.reqP[reflect.TypeOf(c)] = true
		}
	}
}

func FutureComponent(c interface{}) Option {
	return func(proc *process) {
		if proc.fut != nil {
			panic("ecs: future type set multiple times")
		}
		typ := reflect.TypeOf(c)
		proc.fut = &typ
	}
}

func Name(name string) Option {
	return func(proc *process) {
		proc.name = name
	}
}

func (w *world) AddProcess(fn Process, options ...Option) {
	w.Lock()
	defer w.Unlock()

	proc := &process{fn: fn, reqP: make(map[reflect.Type]bool)}

	for _, o := range options {
		o(proc)
	}

	if proc.name == "" {
		panic("ecs: process name cannot be empty")
	}

	w.procs = append(w.procs, proc)
	//w.procs.resolveDependencies()
}

type EntityIterator interface {
	Next() bool
	Current() uint
	GetPast(c interface{}) (component interface{}, ok bool)
	SetFuture(interface{})
}

type entityIterator struct {
	world       *world
	proc        *process
	pastIndices map[reflect.Type]int
	futureIndex *uint
	current     *uint
}

func (i *entityIterator) Next() bool {
	if i.current == nil {
		var cur uint = 0
		i.current = &cur
	} else {
		*i.current++
	}
	if int(*i.current) >= int(i.world.past.maxId)-len(i.world.past.freeIds) {
		return false
	}
	matches := true
	for typ, h := range i.proc.reqP {
		if h && i.world.past.componentCounts[typ][*i.current] {
			i.pastIndices[typ]++
		} else {
			matches = false
		}
	}
	//for j := range i.world.past.componentCounts {
	//	if i.world.past.componentCounts[j][i.currentEntity] {
	//		i.pastIndices[j] += 1
	//	} else {
	//		if i.reqPast[j] {
	//			matches = false
	//		}
	//	}
	//}
	if !matches {
		return i.Next()
	}
	if i.futureIndex != nil {
		*i.futureIndex++
	}
	return true
}

func (i *entityIterator) Current() uint {
	if i.current == nil {
		panic("ecs: current called before first iteration")
	}
	return *i.current
}

func (i *entityIterator) GetPast(c interface{}) (component interface{}, ok bool) {
	typ := reflect.TypeOf(c)
	if i.current == nil {
		panic("ecs: past called before first iteration")
	}
	if i.proc.reqP[typ] {
		cs, ok := i.world.past.components[typ]
		if ok {
			fmt.Println(reflect.TypeOf(cs.Index(i.pastIndices[typ]).Interface()))
			fmt.Println(cs.Index(i.pastIndices[typ]))
			return cs.Index(i.pastIndices[typ]).Interface(), true
		}
	}
	return nil, false
}

func (i *entityIterator) SetFuture(c interface{}) {
	if i.current == nil {
		panic("ecs: future called before first iteration")
	}
	if i.proc.fut == nil {
		panic("ecs: future called by a process that does not modify future state")
	}
	i.world.future.components[reflect.TypeOf(c)].Index(int(*i.futureIndex)).Set(reflect.ValueOf(c))
}

type Process func(iterator EntityIterator)
