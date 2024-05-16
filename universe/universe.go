package main

import (
	"fmt"
	"log/slog"
	"sync/atomic"
)

var id atomic.Uint64

func NextId() uint64 {
	return id.Add(1)
}

type IObject interface {
	String() string
	GetId() uint64
	SetUniverse(*Universe)
	GetUniverse() *Universe
	ProcessPhysics()
	// Point() Point
}

// type IUniverseObject interface {
// 	SetUniverse(IUniverse)
// 	GetUniverse() IUniverse
// }

// type IUniverse interface {
// 	String() string
// 	GetId() uint64
// 	Rect() *Rect
// 	Add(IObject)
// 	Del(IObject)
// 	ProcessPhysics()
// }

// type Object struct {
// 	id       uint64
// 	name     string
// 	universe IUniverse
// }

// func NewObject(universe IUniverse) IObject {
// 	id := NextId()
// 	obj := Object{
// 		id:       id,
// 		name:     fmt.Sprintf("Object-%d", id),
// 		universe: nil,
// 	}
// 	// slog.Info("the object is created", slog.Uint64("object", obj.id))
// 	slog.Info("the object is created", "object", obj.id)
// 	obj.SetUniverse(universe)
// 	return &obj
// }

// func (object *Object) String() string {
// 	return object.name
// }

// func (object *Object) GetId() uint64 {
// 	return object.id
// }

// func (object *Object) GetUniverse() IUniverse {
// 	return object.universe
// }

// func (object *Object) SetUniverse(universe IUniverse) {
// 	if object.universe != universe {
// 		if object.universe != nil {
// 			old_universe := object.universe
// 			object.universe = nil
// 			old_universe.Del(object)
// 		}
// 		if universe != nil {
// 			object.universe = universe
// 			universe.Add(object)
// 		}
// 	}
// }

// func (object *Object) ProcessPhysics() {
// 	// slog.Info("ProcessPhysics", slog.Uint64("object", object.id))
// 	slog.Info("ProcessPhysics", "object", object)
// }

type Universe struct {
	id            uint64
	name          string
	rect          Rect
	objects       []IObject
	tik           int
	simulationTik chan *Universe
}

func NewUniverse(rect Rect) *Universe {
	id := NextId()
	universe := Universe{
		id:            id,
		name:          fmt.Sprintf("Universe-%d", id),
		rect:          rect,
		objects:       []IObject{},
		simulationTik: make(chan *Universe),
	}
	slog.Info("the universe is created", slog.Uint64("universe", universe.id), "rect", rect)
	return &universe
}

func (universe *Universe) String() string {
	return universe.name
}

func (universe *Universe) GetId() uint64 {
	return universe.id
}

func (universe *Universe) Rect() *Rect {
	return &universe.rect
}

func (universe *Universe) SimulationTik() chan *Universe {
	return universe.simulationTik
}

func (universe *Universe) Add(obj IObject) {
	if obj.GetUniverse() != universe {
		if obj.GetUniverse() != nil {
			obj.GetUniverse().Del(obj)
		}
		obj.SetUniverse(universe)
	}
	// if _, ok := universe.objects[obj.GetId()]; !ok {
	// 	universe.objects[obj.GetId()] = obj
	// 	slog.Info("object is added", slog.Uint64("universe", universe.id), slog.Uint64("object", obj.GetId()))
	// }
	universe.objects = append(universe.objects, obj)
}

func (universe *Universe) Del(obj IObject) {
	if obj.GetUniverse() == universe {
		obj.SetUniverse(nil)
	}
	// if _, ok := universe.objects[obj.GetId()]; ok {
	// 	delete(universe.objects, obj.GetId())
	// 	slog.Info("object is removed", slog.Uint64("universe", universe.id), slog.Uint64("object", obj.GetId()))
	// }
	objects := []IObject{}
	for _, o := range universe.objects {
		if o != obj {
			objects = append(objects, o)
		}
	}
	universe.objects = objects

	panic(1)
}

func (universe *Universe) ProcessPhysics() {
	slog.Info("=========================")
	universe.tik += 1
	universe.simulationTik <- universe
	slog.Info("ProcessPhysics", "universe", universe.id, "tik", universe.tik)
	for _, obj := range universe.objects {
		obj.ProcessPhysics()
	}
}
