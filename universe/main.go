package main

import (
	"log"
	"log/slog"
	"os"
	"os/signal"
)

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
