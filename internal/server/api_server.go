package server

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/hirakiuc/redash2pagerduty/internal/redash"
)

// APIServer describe the API server.
type APIServer struct {
	listener   net.Listener
	errCh      chan error
	dispatcher *Dispatcher

	// Context
	ctx    context.Context
	cancel context.CancelFunc
}

func writeResponse(w http.ResponseWriter, code int, status string, result string) {
	w.WriteHeader(code)
	fmt.Fprintf(w, "{\"status\":\"%s\", \"result\":\"%s\"}\n", status, result)
}

func handleNotify(d *Dispatcher, ctx context.Context) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		event, err := decodeJSON(w, r)
		if err != nil {
			fmt.Printf("  => Failed decode %v\n", err)
			return
		}

		// Enqueue task
		d.Work(ctx, func(ctx context.Context) {
			fmt.Fprintf(os.Stdout, "Worked ! %s\n", event.Event)
		})

		writeResponse(w, http.StatusOK, "ok", "Received")
	}
}

func decodeJSON(w http.ResponseWriter, r *http.Request) (*redash.WebhookEvent, error) {
	if r.Body == nil {
		writeResponse(w, http.StatusBadRequest, "failure", "require json")
		return nil, errors.New("Require json")
	}

	// TODO validate with JsonSchema
	event, err := redash.Parse(r.Body)
	if err != nil {
		writeResponse(w, http.StatusBadRequest, "failure", "Invalid json")
		return nil, err
	}

	return event, nil
}

func routes(ctx context.Context, d *Dispatcher) *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/", handleNotify(d, ctx)).
		Methods("POST")

	return r
}

// NewAPIServer return a API server instance.
func NewAPIServer(addr string, numOfworkers int) (*APIServer, error) {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, err
	}

	d := NewDispatcher(numOfworkers)
	ctx, cancel := context.WithCancel(context.Background())
	ch := make(chan error)
	go func() {
		mux := routes(ctx, d)
		ch <- http.Serve(listener, mux)
	}()

	return &APIServer{
		listener:   listener,
		errCh:      ch,
		dispatcher: d,

		ctx:    ctx,
		cancel: cancel,
	}, nil
}

// Listen start listening API request.
func (srv *APIServer) Listen() {
	<-srv.errCh
	defer close(srv.errCh)
}

// Addr return this API server address.
func (srv *APIServer) Addr() net.Addr {
	return srv.listener.Addr()
}

// Close cleanup this API server.
func (srv *APIServer) Close() {
	(srv.dispatcher).Wait()
	fmt.Println("Dispatcher wait end.")

	srv.cancel()
	fmt.Println("cancel.")

	(srv.listener).Close()
	fmt.Println("Listener close.")
}
