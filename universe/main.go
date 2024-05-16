package main

import (
	"log"
	"log/slog"
	"os"
	"os/signal"
)

// func run() {
// 	slog.Info("----")
// 	universe := NewUniverse(NewRect(0, 0, 200, 200))
// 	star := NewStar(universe, NewRect(0, 0, 50, 50))
// 	planet := NewPlanet(universe, star, 10, 1)
// 	ship := NewShip(universe, star, &Point{2.1234, 2.4321}, 4.83)
// 	// ship.MoveToPoint(Point{40, 30})
// 	ship.LandOn(planet)

// 	ctx, cancel := context.WithCancel(context.Background())
// 	defer cancel()

// 	go (func() {
// 		slog.Info("11")
// 		running := true
// 		for running {
// 			select {
// 			case <-ctx.Done():
// 				running = false
// 			case u := <-universe.SimulationTik():
// 				if u != nil {
// 					slog.Info("tik", "u", u.GetId())
// 				}
// 			}
// 		}
// 		slog.Info("12")
// 	})()

// 	// align time
// 	a := time.Now().Truncate(time.Second).Add(1 * time.Second)
// 	slog.Info("alignment", "a", a)
// 	aligin := time.NewTimer(time.Until(a))
// 	<-aligin.C

// 	tiker := time.NewTicker(1 * time.Second)
// 	defer tiker.Stop()

// 	osSignal := make(chan os.Signal, 1)
// 	signal.Notify(osSignal, os.Interrupt)

// 	running := true
// 	// out:
// 	for running {
// 		universe.ProcessPhysics()
// 		select {
// 		case s := <-osSignal:
// 			slog.Info("os signal received", "signal", s)
// 			running = false
// 			// break out
// 		case <-ctx.Done():
// 			slog.Info("ctx.done")
// 			running = false
// 			// break out
// 		case <-tiker.C:
// 		}
// 	}

// 	cancel()
// 	tiker.Stop()
// }

func run() {
	slog.Info("----")
	universe := NewUniverse(NewRect(0, 0, 200, 200))
	star := NewStar(universe, NewRect(0, 0, 50, 50))
	planet := NewPlanet(universe, star, 10, 1)
	ship := NewShip(universe, star, &Point{2.1234, 2.4321}, 4.83)
	// ship.MoveToPoint(Point{40, 30})
	ship.LandOn(planet)

	go (func() {
		slog.Info("11")
		running := true
		for running {
			// select {
			// case u := <-universe.SimulationTik():
			// 	if u != nil {
			// 		slog.Info("tik", "u", u.GetId())
			// 	}
			// }
			u := <-universe.SimulationTik()
			slog.Info("tik", "u", u.GetId())
		}
		slog.Info("12")
	})()

	osSignal := make(chan os.Signal, 1)
	signal.Notify(osSignal, os.Interrupt)

	go universe.Run()

	s := <-osSignal
	slog.Info("os signal received", "signal", s)

	universe.Stop()
}

func main() {
	// log.SetFlags(log.LstdFlags | log.Lmicroseconds | log.Lshortfile)
	log.SetFlags(log.Lmicroseconds)

	// logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
	// 	AddSource: false,
	// }))
	// slog.SetDefault(logger)

	slog.Info("started")
	run()
	slog.Info("stopped")
}
