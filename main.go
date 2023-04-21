package main

import (
	"context"
	"errors"
	"log"
	"net"
	"os"
	"os/signal"

	"github.com/arianvp/systemd-creds/server"
	"github.com/coreos/go-systemd/v22/activation"
	"golang.org/x/sys/unix"
)

func main() {
	os.Exit(run())
}

func run() int {
	ctx, cancelFunc := signal.NotifyContext(context.Background(), os.Interrupt, unix.SIGTERM)
	defer cancelFunc()

	listeners, err := activation.Listeners()
	if err != nil {
		log.Print(err)
		return 1
	}

	if len(listeners) == 0 {
		l, err := net.Listen("unix", "/tmp/test.sock")
		if err != nil {
			log.Print(err)
			return 1
		}
		listeners = []net.Listener{l}
	}
	for _, l := range listeners {
		defer l.Close()
	}

	connChan := make(chan net.Conn, 100)
	defer close(connChan)

	for _, l := range listeners {
		// This is a bit confusing still
		go func(l net.Listener) {
			for {
				conn, err := l.Accept()
				if err != nil {
					// The listener is closed when ctx is cancelled
					if errors.Is(err, net.ErrClosed) {
						return
					}
					log.Print(err)
					cancelFunc()
				}
				connChan <- conn
			}
		}(l)
	}

	s := server.Server{
		Store: nil, // TODO: Store
	}

	go s.Start(ctx, connChan)
	<-ctx.Done()
	return 0
}
