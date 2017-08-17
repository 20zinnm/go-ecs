package main

import (
	"github.com/20zinnm/go-ecs"
	"fmt"
)

type Position struct {
	X, Y float32
}

type Velocity struct {
	X, Y float32
}

func main() {
	world := ecs.NewWorld()
	world.AddProcess(func(vel Velocity, pos Position) (nextPos Position) {
		nextPos.X = pos.X + vel.X
		nextPos.Y = pos.Y + vel.Y
		return
	})
	var tick = 0
	world.AddProcess(func(pos Position) {
		tick++
		fmt.Printf("Tick %d: Entity at (%f,%f)\n", tick, pos.X, pos.Y)
	})
	world.NewEntity(Position{10, 10}, Velocity{1, 1})
	for i := 0; i < 10; i++ {
		world.ExecuteFrame()
	}
	//fmt.Println(id)
}

//
//func PhysicsProc(vel Velocity, pos Position) (nextPos Position) {
//	nextPos.X = pos.X + vel.X
//	nextPos.Y += pos.Y + vel.Y
//	return
//}