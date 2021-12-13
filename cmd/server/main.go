package main

import (
	"log"
	"net/http"
	"time"

	"github.com/moynur/gateway/internal/helpers"
	"github.com/moynur/gateway/internal/service"
	"github.com/moynur/gateway/internal/store"
	"github.com/moynur/gateway/internal/transport/http"

	"github.com/gorilla/mux"
)

type Handler interface {
	ServeHTTP(ResponseWriter, *http.Request)
}

type ResponseWriter interface {
	Header() http.Header
	Write([]byte) (int, error)
	WriteHeader(statusCode int)
}

func main() {
	db, err := store.NewStore()
	if err != nil {
		log.Println(err)
	}
	gen := helpers.NewGenerator()
	svc := service.NewService(db, gen)
	client, _ := handler.NewHandler(svc)

	r := mux.NewRouter()
	client.ApplyRoutes(r)

	// should be as a config
	address := "0.0.0.0:8081"

	srv := &http.Server{
		Addr:              address,
		Handler:           r,
		ReadTimeout:       1 * time.Second,
		ReadHeaderTimeout: 1 * time.Second,
		WriteTimeout:      1 * time.Second,
		IdleTimeout:       1 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())

}
