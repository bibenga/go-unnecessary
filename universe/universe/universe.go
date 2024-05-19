package universe

import (
	"fmt"
	"log/slog"
	"sync/atomic"
	"time"
)

var id atomic.Uint64

func NextId() uint64 {
	return id.Add(1)
}

type IObject interface {
	GetId() uint64
	SetUniverse(*Universe)
	GetUniverse() *Universe
	ProcessPhysics()
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
	id      uint64
	log     *slog.Logger
	name    string
	rect    Rect
	objects []IObject
	running atomic.Bool
	tik     uint64
	stop    chan int
	stopped chan int
}

var _ fmt.Stringer = &Universe{}

func NewUniverse(rect Rect) *Universe {
	id := NextId()
	universe := Universe{
		id:      id,
		log:     slog.Default().With("universe", id),
		name:    fmt.Sprintf("Universe-%d", id),
		rect:    rect,
		objects: []IObject{},
		stop:    make(chan int),
		stopped: make(chan int),
	}
	universe.running.Store(false)
	universe.log.Info("the universe is created", "rect", rect)
	return &universe
}

func (universe *Universe) String() string {
	return universe.name
}

func (universe *Universe) GetId() uint64 {
	return universe.id
}

func (universe *Universe) Log() *slog.Logger {
	return universe.log
}

func (universe *Universe) Rect() *Rect {
	return &universe.rect
}

func (universe *Universe) Add(obj IObject) {
	if obj.GetUniverse() != universe {
		if obj.GetUniverse() != nil {
			obj.GetUniverse().Del(obj)
		}
		obj.SetUniverse(universe)
	}
	universe.objects = append(universe.objects, obj)
}

func (universe *Universe) Del(obj IObject) {
	if obj.GetUniverse() == universe {
		obj.SetUniverse(nil)
	}
	for i, o := range universe.objects {
		if o == obj {
			objects := universe.objects
			if len(objects) == 1 {
				universe.objects = []IObject{}
			} else {
				universe.objects = append(objects[:i], objects[i+1:]...)
			}
		}
	}
}

func (universe *Universe) ProcessPhysics() {
	slog.Debug("=========================")
	universe.tik++
	universe.log.Info("The Universe plays with gravity", "tik", universe.tik)
	for _, obj := range universe.objects {
		obj.ProcessPhysics()
	}
}

func (universe *Universe) Run() {
	tiker := time.NewTicker(1 * time.Second)
	defer tiker.Stop()

	// universe.ctx, universe.cancel = context.WithCancel(context.Background())
	// defer universe.cancel()
	// defer universe.stopWait.Done()

	defer func() {
		universe.running.Store(false)
		universe.stopped <- 1
	}()

	universe.log.Info("simulation started")
	universe.running.Store(true)
out:
	for {
		select {
		// case <-universe.ctx.Done():
		// slog.Info("stopSignal received")
		// 	break out
		case <-universe.stop:
			universe.log.Info("someone wants to stop The Universe")
			break out
		case <-tiker.C:
			universe.ProcessPhysics()
		}
	}

	universe.log.Info("simulation stopped")
}

func (universe *Universe) Stop() {
	// can be called once
	if universe.running.Load() {
		universe.stop <- 1
		<-universe.stopped
		universe.log.Info("The Universe is saved or not")
	}
}
