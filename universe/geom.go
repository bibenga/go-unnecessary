package main

import (
	"fmt"
	"math"
	"math/rand"
)

const threshold = 1e-9
const tau = 2 * math.Pi

type Point struct {
	x, y float64
}

func NewPoint(x, y float64) *Point {
	return &Point{x: x, y: y}
}

// func (point *Point) String() string {
// 	return fmt.Sprintf("(%.2f,%.2f)", point.x, point.y)
// }

func (point Point) String() string {
	return fmt.Sprintf("{%.2f,%.2f}", point.x, point.y)
}

func (point *Point) X() float64 {
	return point.x
}

func (point *Point) Y() float64 {
	return point.y
}

func (point *Point) NearlyEqual(p2 *Point) bool {
	return math.Abs(point.x-p2.x) < threshold && math.Abs(point.y-p2.y) < threshold
}

// func (point *Point) Add(p2 *Point) *Point {
// 	return &Point{
// 		x: point.x + p2.x,
// 		y: point.y + p2.y,
// 	}
// }

// func (point *Point) Sub(p2 *Point) *Point {
// 	return &Point{
// 		x: point.x - p2.x,
// 		y: point.y - p2.y,
// 	}
// }

func (point *Point) DistanceTo(p2 *Point) float64 {
	return Distance(point, p2)
}

func Distance(p1, p2 *Point) float64 {
	dx := (p1.x - p2.x)
	dy := (p1.y - p2.y)
	d := math.Pow(math.Pow(dx, 2)+math.Pow(dy, 2), 0.5)
	return d
}

func NextPointOnLine(p1, p2 *Point, velocity float64) *Point {
	d := Distance(p1, p2)
	if d < velocity {
		return p2
	}
	t := velocity / d
	return &Point{
		x: (1-t)*p1.x + t*p2.x,
		y: (1-t)*p1.y + t*p2.y,
	}
}

func NextTheta(theta, delta float64) float64 {
	theta += delta
	if delta > 0 {
		if theta > tau {
			theta -= tau
		}
	} else {
		if theta < 0 {
			theta += tau
		}
	}
	return theta
}

func NextPointOnCircle(center *Point, radius, theta float64) *Point {
	return &Point{
		x: radius*math.Cos(theta) + center.x,
		y: radius*math.Sin(theta) + center.y,
	}
}

type Rect struct {
	p1, p2 Point
}

func NewRect(x1, y1, x2, y2 float64) Rect {
	return Rect{
		p1: Point{x: x1, y: y1},
		p2: Point{x: x2, y: y2},
	}
}

// func (rect *Rect) String() string {
// 	return fmt.Sprintf("((%.2f,%.2f),(%.2f,%.2f))", rect.p1.x, rect.p1.y, rect.p2.x, rect.p2.y)
// }

func (rect Rect) String() string {
	// return fmt.Sprintf("{{%.2f,%.2f},{%.2f,%.2f}}", rect.p1.x, rect.p1.y, rect.p2.x, rect.p2.y)
	return fmt.Sprintf("{%v,%v}", rect.p1, rect.p2)
}

func (rect *Rect) P1() Point {
	return rect.p1
}

func (rect *Rect) P2() Point {
	return rect.p2
}

func (rect *Rect) Center() Point {
	p1, p2 := rect.p1, rect.p2
	return Point{
		x: p1.x + (p2.x-p1.x)/2,
		y: p1.y + (p2.y-p1.y)/2,
	}
}

func (rect *Rect) ContainsPoint(p Point) bool {
	p1, p2 := rect.p1, rect.p2
	return (p1.x <= p.x && p.x <= p2.x) && (p1.y <= p.y && p.y <= p2.y)
}

func (rect *Rect) ContainsRect(r Rect) bool {
	return rect.ContainsPoint(r.p1) && rect.ContainsPoint(r.p2)
}

func (rect *Rect) RendomPoint() *Point {
	p1, p2 := rect.p1, rect.p2
	return &Point{
		x: p1.x + rand.Float64()*(p2.x-p1.x),
		y: p1.y + rand.Float64()*(p2.y-p1.y),
	}
}
