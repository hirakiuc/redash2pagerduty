package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/hirakiuc/redash2pagerduty/internal/server"
)

// NumOfWorkers describe the number of workers.
const NumOfWorkers int = 3

func main() {
	srv, err := server.NewAPIServer(":8080", NumOfWorkers)
	if err != nil {
		fmt.Printf("Error %v\n", err)
		os.Exit(-1)
	}

	defer srv.Close()
	fmt.Println("Server started at", srv.Addr())

	sigCh := make(chan os.Signal, 1)
	defer close(sigCh)

	signal.Notify(sigCh, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGINT)
	go func() {
		fmt.Println("listening signals.")
		<-sigCh
		fmt.Println("Signal received.")
		srv.Close()
	}()

	srv.Listen()
}
