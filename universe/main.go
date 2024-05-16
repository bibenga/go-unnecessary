package main

import (
	"log"
	"log/slog"
	"os"
	"os/signal"

	u "unnecessary/universe/universe"
)

func run() {
	slog.Info("----")
	universe := u.NewUniverse(u.NewRect(0, 0, 200, 200))
	star := u.NewStar(universe, u.NewRect(0, 0, 50, 50))
	planet := u.NewPlanet(universe, star, 10, 1)
	ship := u.NewShip(universe, star, u.NewPoint(2.1234, 2.4321), 4.83)

	planet.String()
	ship.MoveToPoint(*u.NewPoint(40, 30))
	// ship.LandOn(planet)

	// go (func() {
	// 	slog.Info("11")
	// 	running := true
	// 	for running {
	// 		// select {
	// 		// case u := <-universe.SimulationTik():
	// 		// 	if u != nil {
	// 		// 		slog.Info("tik", "u", u.GetId())
	// 		// 	}
	// 		// }
	// 		u := <-universe.SimulationTik()
	// 		slog.Info("tik", "u", u.GetId())
	// 	}
	// 	slog.Info("12")
	// })()

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
