package main

import (
	"fmt"
	"log/slog"
	"math"
	"math/rand"
)

type Planet struct {
	id       uint64
	name     string
	universe *Universe
	star     *Star
	radius   float64
	velocity float64
	theta_dt float64
	theta    float64
	point    Point
}

func NewPlanet(universe *Universe, star *Star, radius, velocity float64) *Planet {
	id := NextId()

	center := star.Point()
	theta := rand.Float64() * 2 * math.Pi
	point := Point{
		x: radius*math.Cos(theta) + center.x,
		y: radius*math.Sin(theta) + center.y,
	}

	planet := Planet{
		id:       id,
		name:     fmt.Sprintf("Planet-%d", id),
		universe: nil,
		star:     nil,
		radius:   radius,
		velocity: velocity,
		theta_dt: velocity / radius,
		theta:    theta,
		point:    point,
	}
	slog.Info("the planet is created", slog.Uint64("planet", planet.id))
	if universe != nil {
		planet.SetUniverse(universe)
		universe.Add(&planet)
	}
	if star != nil {
		planet.SetStar(star)
		star.Add(&planet)
	}
	return &planet
}

func (planet *Planet) String() string {
	return planet.name
}

func (planet *Planet) GetId() uint64 {
	return planet.id
}

func (planet *Planet) GetUniverse() *Universe {
	return planet.universe
}

func (planet *Planet) SetUniverse(universe *Universe) {
	// if planet.universe != universe {
	// 	if planet.universe != nil {
	// 		old_universe := planet.universe
	// 		planet.universe = nil
	// 		old_universe.Del(planet)
	// 	}
	// 	if universe != nil {
	// 		planet.universe = universe
	// 		universe.Add(planet)
	// 	}
	// }
	planet.universe = universe
}

func (planet *Planet) GetStar() *Star {
	return planet.star
}

func (planet *Planet) SetStar(star *Star) {
	// if planet.star != star {
	// 	// if planet.star != nil {
	// 	// 	old_star := ship.star
	// 	// 	planet.star = nil
	// 	// 	old_star.Del(ship)
	// 	// }
	// 	if star != nil {
	// 		planet.star = star
	// 		// star.Add(ship)
	// 	}
	// }
	planet.star = star
}

func (planet *Planet) Point() Point {
	return planet.point
}

func (planet *Planet) ProcessPhysics() {
	// slog.Info("ProcessPhysics", slog.Uint64("object", object.id))
	// theta_dt := planet.velocity / planet.radius
	planet.theta += planet.theta_dt
	if planet.theta >= tau {
		planet.theta -= tau
	}
	center := planet.star.Point()
	planet.point = Point{
		x: planet.radius*math.Cos(planet.theta) + center.x,
		y: planet.radius*math.Sin(planet.theta) + center.y,
	}

	slog.Info("ProcessPhysics", "planet", planet, "point", planet.point)
}
