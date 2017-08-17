package main

import (
	"github.com/20zinnm/go-ecs"
	"fmt"
	"github.com/20zinnm/go-ecs/src/components"
)

func main() {
	world := ecs.NewWorld()
	world.AddProcess(func(ei ecs.EntityIterator) {
		for ei.Next() {
			posn, ok := ei.GetPast(components.Position{})
			if !ok {
				panic("got an entity that does not fulfill the past component requirements")
			}
			pos, ok := posn.(components.Position)
			if !ok {
				panic("component retrieved is of the wrong type")
			}
			fmt.Printf("Processing entity %d located at %f, %f.", ei.Current(), pos.X, pos.Y)
		}
	}, ecs.Name("printer"), ecs.PastComponent(components.Position{}))
	myEntity := world.NewEntity(components.Position{10.10, 10.11})
	world.ExecuteFrame()
	fmt.Println("Yay! I made an entity!")
	fmt.Printf("Its ID is %d.", myEntity)
}
