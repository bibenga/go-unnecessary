package universe

import (
	"fmt"
	"log/slog"
)

type Ship struct {
	id          uint64
	log         *slog.Logger
	name        string
	universe    *Universe
	star        *Star
	point       Point
	fuel        int64
	maxVelocity float64
	velocity    float64
	destination *Point
	planet      *Planet
	landed      bool
}

var _ IObject = &Ship{}
var _ fmt.Stringer = &Ship{}

func NewShip(universe *Universe, star *Star, point *Point, velocity float64) *Ship {
	id := NextId()
	if point == nil {
		point = universe.Rect().RendomPoint()
	}
	ship := Ship{
		id:          id,
		log:         slog.Default().With("universe", universe.GetId(), "star", star.GetId(), "ship", id),
		name:        fmt.Sprintf("Ship-%d", id),
		universe:    nil,
		star:        nil,
		point:       *point,
		fuel:        1000,
		landed:      false,
		maxVelocity: velocity,
	}
	ship.log.Info("the ship is created")
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
	ship.log.Info("move to a point", "order", "move", "point", ship.point)
}

func (ship *Ship) LandOn(planet *Planet) {
	if ship.landed {
		return
	}
	ship.velocity = ship.maxVelocity
	ship.destination = nil
	ship.planet = planet
	ship.landed = false
	ship.log.Info("move to a planet", "order", "land", "planet", planet)
}

func (ship *Ship) ProcessPhysics() {
	if ship.planet != nil {
		if ship.landed {
			ship.fuel = 1000
			ship.point = ship.planet.Point()
			ship.log.Info("relax on the planet", "planet", ship.planet, "point", ship.point, "fuel", ship.fuel)
		} else {
			ship.fuel -= 5
			p1, p2 := ship.point, ship.planet.Point()
			d := Distance(&p1, &p2)
			if d <= ship.velocity {
				ship.point = p2
				ship.landed = true
				ship.log.Info("land on the planet", "planet", ship.planet, "point", ship.point, "fuel", ship.fuel)
			} else {
				t := ship.velocity / d
				ship.point = Point{
					x: (1-t)*p1.x + t*p2.x,
					y: (1-t)*p1.y + t*p2.y,
				}
				ship.log.Info("move to the planet", "planet", ship.planet, "point", ship.point, "fuel", ship.fuel)
			}
		}
	} else if ship.destination != nil {
		ship.fuel -= 5
		p1, p2 := &ship.point, ship.destination
		d := Distance(p1, p2)
		if d <= ship.velocity {
			ship.velocity = 0
			ship.point = *ship.destination
			ship.destination = nil
			ship.log.Info("arrived to the point", "point", ship.point, "fuel", ship.fuel)
		} else {
			t := ship.velocity / d
			ship.point = Point{
				x: (1-t)*p1.x + t*p2.x,
				y: (1-t)*p1.y + t*p2.y,
			}
			ship.log.Info("move to the point", "point", ship.point, "fuel", ship.fuel)
		}
	} else {
		ship.fuel -= 1
		ship.log.Info("relax at the point", "point", ship.point, "fuel", ship.fuel)
	}
}
