package main

import (
	"fmt"
	"log/slog"
)

type IStar interface {
	IObject
	Add(IStarObject)
}

type IStarObject interface {
	IObject
	SetStar(IStar)
	GetStar() IStar
}

type Star struct {
	id       uint64
	name     string
	universe *Universe
	objects  []IObject
	rect     Rect
	point    Point
}

func NewStar(universe *Universe, rect Rect) *Star {
	id := NextId()
	star := Star{
		id:       id,
		name:     fmt.Sprintf("Star-%d", id),
		universe: nil,
		objects:  []IObject{},
		rect:     rect,
		point:    rect.Center(),
	}
	slog.Info("the star is created", slog.Uint64("star", star.id))
	if universe != nil {
		star.SetUniverse(universe)
		universe.Add(&star)
	}
	return &star
}

func (star *Star) String() string {
	return star.name
}

func (star *Star) GetId() uint64 {
	return star.id
}

func (star *Star) GetUniverse() *Universe {
	return star.universe
}

func (star *Star) SetUniverse(universe *Universe) {
	// if star.universe != universe {
	// 	if star.universe != nil {
	// 		old_universe := star.universe
	// 		star.universe = nil
	// 		old_universe.Del(star)
	// 	}
	// 	if universe != nil {
	// 		star.universe = universe
	// 		universe.Add(star)
	// 	}
	// }
	star.universe = universe
}

func (star *Star) Add(obj IObject) {
	// if planet.GetStar() != star {
	// 	planet.SetStar(star)
	// }
	// if _, ok := universe.objects[obj.GetId()]; !ok {
	// 	universe.objects[obj.GetId()] = obj
	// 	slog.Info("object is added", slog.Uint64("universe", universe.id), slog.Uint64("object", obj.GetId()))
	// }
	star.objects = append(star.objects, obj)
}

func (star *Star) Point() Point {
	return star.point
}

func (star *Star) ProcessPhysics() {
	// slog.Info("ProcessPhysics", slog.Uint64("object", object.id))
	slog.Info("ProcessPhysics", "star", star, "point", star.point)
}
