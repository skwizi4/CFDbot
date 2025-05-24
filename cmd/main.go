package main

import (
	"TradeBot/internal/app"
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	var testRange string
	a := app.New("MusicApp")
	a.Run(gracefulShutDown(), testRange)
}

func gracefulShutDown() context.Context {
	ctx, cancel := context.WithCancel(context.Background())
	c := make(chan os.Signal)

	signal.Notify(c, syscall.SIGTERM, os.Interrupt)
	go func() {
		<-c
		log.Print("services stopped by graceful	ShutDown")
		cancel()
		os.Exit(0)

	}()

	return ctx
}
