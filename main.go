package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gorilla/mux"
)

type Message struct {
	Id      string `json:"id"`
	Balance uint64 `json:"balance"`
}

func (msg *Message) IsValid() bool {
	return len(msg.Id) > 0
}

func handleNotify(d *Dispatcher, ctx context.Context) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		msg, err := decodeJson(w, r)
		if err != nil {
			return
		}

		// Enqueue task
		d.Work(ctx, func(ctx context.Context) {
			fmt.Fprintf(os.Stdout, "Worked ! %s:%d\n", msg.Id, msg.Balance)
		})

		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "{\"result\": \"Received\"}")
	}
}

func decodeJson(w http.ResponseWriter, r *http.Request) (*Message, error) {
	if r.Body == nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "{\"result\":\"require json\"}")
		return nil, errors.New("Require json")
	}

	msg := Message{}
	err := json.NewDecoder(r.Body).Decode(&msg)
	if err != nil {
		fmt.Fprintf(w, "{\"result\":\"decode failed\"}")
		return nil, err
	}
	if msg.IsValid() == false {
		fmt.Fprintf(w, "{\"result\":\"unexpected json\"}")
		return nil, errors.New("Invalid json")
	}

	return &msg, nil
}

func routes(ctx context.Context, d *Dispatcher) *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/", handleNotify(d, ctx)).
		Methods("POST")

	return r
}

func server(ctx context.Context, addr string, d *Dispatcher) (listener net.Listener, ch chan error) {
	ch = make(chan error)

	listener, err := net.Listen("tcp", addr)
	if err != nil {
		panic(err)
	}

	go func() {
		mux := routes(ctx, d)
		ch <- http.Serve(listener, mux)
	}()

	return listener, ch
}

func main() {
	d := NewDispatcher(3)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	listener, ch := server(ctx, ":8080", d)
	defer listener.Close()
	fmt.Println("Server started at", listener.Addr())

	sigCh := make(chan os.Signal, 1)
	defer close(sigCh)

	signal.Notify(sigCh, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGINT)
	go func() {
		fmt.Println("listening signals.")
		<-sigCh
		fmt.Println("Signal received.")
		d.Wait()
		fmt.Println("Dispatcher wait end.")
		cancel()
		fmt.Println("cancel.")
		listener.Close()
		fmt.Println("Listener close.")
	}()

	<-ch
}
