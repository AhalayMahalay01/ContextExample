package main

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"time"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	go handleSignals(cancel)

	if err := startServer(ctx); err != nil {
		log.Fatal(err)

	}
}
func handleSignals(cancel context.CancelFunc) {
	sigCh := make(chan os.Signal)
	signal.Notify(sigCh, os.Interrupt)
	for {
		sig := <-sigCh
		switch sig {
		case os.Interrupt:
			cancel()
			return
		}
	}
}

func startServer(ctx context.Context) error {
	ladder, err := net.ResolveTCPAddr("tcp", ":8084")
	if err != nil {
		return err
	}

	l, err := net.ListenTCP("tcp", ladder)
	if err != nil {
		return err
	}
	defer l.Close()

	for {
		select {
		case <-ctx.Done():
			log.Println("Server STOP!")
			return nil
		default:
			if err := l.SetDeadline(time.Now().Add(time.Second)); err != nil {
				return err
			}

			_, err := l.Accept()
			if err != nil {
				if os.IsTimeout(err) {
					continue
				}
				return err
			}
			log.Println("New client connected!")
		}
	}
}
