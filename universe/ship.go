package main

import (
	"fmt"
	"log/slog"
)

// type IShip interface {
// 	IObject
// 	MoveToPoint(point Point)
// }

type Ship struct {
	id          uint64
	name        string
	universe    *Universe
	star        *Star
	point       Point
	maxVelocity float64
	velocity    float64
	destination *Point
	planet      *Planet
	landed      bool
}

func NewShip(universe *Universe, star *Star, point *Point, velocity float64) *Ship {
	id := NextId()
	if point == nil {
		point = universe.Rect().RendomPoint()
	}
	ship := Ship{
		id:          id,
		name:        fmt.Sprintf("Ship-%d", id),
		universe:    nil,
		star:        nil,
		point:       *point,
		landed:      false,
		maxVelocity: velocity,
	}
	slog.Info("the ship is created", slog.Uint64("ship", ship.id))
	if universe != nil {
		ship.SetUniverse(universe)
		universe.Add(&ship)
	}
	if star != nil {
		ship.SetStar(star)
		star.Add(&ship)
	}

	return &ship
}

func (ship Ship) String() string {
	return ship.name
}

func (ship *Ship) GetId() uint64 {
	return ship.id
}

func (ship *Ship) GetUniverse() *Universe {
	return ship.universe
}

func (ship *Ship) SetUniverse(universe *Universe) {
	// if ship.universe != universe {
	// 	if ship.universe != nil {
	// 		old_universe := ship.universe
	// 		ship.universe = nil
	// 		old_universe.Del(ship)
	// 	}
	// 	if universe != nil {
	// 		ship.universe = universe
	// 		universe.Add(ship)
	// 	}
	// }
	ship.universe = universe
}

func (ship *Ship) GetStar() *Star {
	return ship.star
}

func (ship *Ship) SetStar(star *Star) {
	// if ship.star != star {
	// 	if ship.star != nil {
	// 		// old_star := ship.star
	// 		ship.star = nil
	// 		// old_star.Del(ship)
	// 	}
	// 	if star != nil {
	// 		ship.star = star
	// 		// star.Add(ship)
	// 	}
	// }
	ship.star = star
}

func (ship *Ship) MoveToPoint(point Point) {
	if ship.landed {
		return
	}
	ship.destination = &point
	ship.velocity = ship.maxVelocity
	slog.Info("order", "ship", ship, "type", "move", "point", ship.point)
}

func (ship *Ship) LandOn(planet *Planet) {
	if ship.landed {
		return
	}
	ship.velocity = ship.maxVelocity
	ship.destination = nil
	ship.planet = planet
	ship.landed = false
	slog.Info("order", "ship", ship, "type", "land", "planet", planet)
}

func (ship *Ship) ProcessPhysics() {
	if ship.planet != nil {
		if ship.landed {
			ship.point = ship.planet.Point()
			slog.Info("relax", "ship", ship, "planet", ship.planet, "point", ship.point)
		} else {
			p1, p2 := ship.point, ship.planet.Point()
			d := Distance(&p1, &p2)
			if d <= ship.velocity {
				ship.point = p2
				ship.landed = true
				slog.Info("landed", "ship", ship, "planet", ship.planet, "point", ship.point)
			} else {
				t := ship.velocity / d
				ship.point = Point{
					x: (1-t)*p1.x + t*p2.x,
					y: (1-t)*p1.y + t*p2.y,
				}
				slog.Info("moved", "ship", ship, "planet", ship.planet, "point", ship.point)
			}
		}
	} else if ship.destination != nil {
		p1, p2 := &ship.point, ship.destination
		d := Distance(p1, p2)
		if d <= ship.velocity {
			ship.velocity = 0
			ship.point = *ship.destination
			ship.destination = nil
			slog.Info("arrived", "ship", ship, "point", ship.point)
		} else {
			t := ship.velocity / d
			ship.point = Point{
				x: (1-t)*p1.x + t*p2.x,
				y: (1-t)*p1.y + t*p2.y,
			}
			slog.Info("moved", "ship", ship, "point", ship.point)
		}
	} else {
		slog.Info("relax", "ship", ship, "point", ship.point)
	}
}

// async def process_physics(self) -> None:
//     # super().process_physics()
//     # if self._universe and self._universe.tik > 1:
//     #     self.destroy()
//     if self._planet is not None:
//         if self._is_landed:
//             if self._should_take_off:
//                 planet = self._planet
//                 self._point = planet.point
//                 self._planet = None
//                 self._is_landed = False
//                 self._should_take_off = False
//                 _l.info('%s: take off from %s: %s', self, planet, self._point)
//                 await ship_took_off.send_async(self, planet=planet)
//             else:
//                 self._point = self._planet.point
//                 _l.info('%s: relaxes on %s', self, self._planet)
//         else:
//             p1, p2 = self._point, self._planet.point
//             d = p1.distance(p2)
//             if d <= self._velocity:
//                 self._point = self._planet.point
//                 self._is_landed = True
//                 _l.info('%s: landed on %s', self, self._planet)
//                 await ship_landed.send_async(self, planet=self._planet)
//             else:
//                 t = self._velocity / d
//                 self._point = Point(
//                     x=(1 - t) * p1.x + t * p2.x,
//                     y=(1 - t) * p1.y + t * p2.y
//                 )
//                 _l.info('%s: flying to %s: %.2f; %s -> %s',
//                         self, self._planet, self._point.distance(self._planet.point),
//                         self._point, self._planet.point)

//     elif self._dst_point is not None:
//         p1, p2 = self._point, self._dst_point
//         d = p1.distance(p2)
//         if d <= self._velocity:
//             self._velocity = 0.0
//             self._point = self._dst_point
//             self._src_point = None
//             self._dst_point = None
//             if not self._universe.rect.contains(self._point):
//                 self.destroy()
//                 await ship_exploded.send_async(self)
//             else:
//                 _l.info('%s: is arrived to %s', self, self.point)
//                 await ship_arrived.send_async(self)
//         else:
//             t = self._velocity / d
//             self._point = Point(
//                 x=(1 - t) * p1.x + t * p2.x,
//                 y=(1 - t) * p1.y + t * p2.y
//             )
//             if not self._universe.rect.contains(self._point):
//                 self.destroy()
//                 await ship_exploded.send_async(self)
//             else:
//                 _l.info('%s: flying to %s: %.2f; %s -> %s',
//                         self, self._dst_point, self._point.distance(self._dst_point),
//                         self._point, self._dst_point)

//     else:
//         _l.info('%s: relaxes at %s', self, self._point)
