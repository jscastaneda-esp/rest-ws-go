package server

import (
	"context"
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/jscastaneda-esp/rest-ws-go/database"
	"github.com/jscastaneda-esp/rest-ws-go/repository"
	"github.com/jscastaneda-esp/rest-ws-go/websocket"
	"github.com/rs/cors"
)

type Config struct {
	Port        string
	JWTSecret   string
	DatabaseUrl string
	RowsDefault string
}

type Server interface {
	Config() *Config
	Hub() *websocket.Hub
}

type Broker struct {
	config *Config
	hub    *websocket.Hub
	router *mux.Router
}

func (broker *Broker) Config() *Config {
	return broker.config
}

func (broker *Broker) Hub() *websocket.Hub {
	return broker.hub
}

func (broker *Broker) Start(binder func(s Server, r *mux.Router)) {

	broker.router = mux.NewRouter()
	binder(broker, broker.router)

	log.Println("Starting connection database")
	repo, err := database.NewPostgresRepository(broker.config.DatabaseUrl)
	if err != nil {
		log.Fatal("Database:", err)
	}
	repository.SetRepository(repo)

	log.Println("Starting websocket server")
	go broker.hub.Run()

	log.Printf("Starting server on 0.0.0.0:%s\n", broker.config.Port)
	handler := cors.Default().Handler(broker.router)
	if err := http.ListenAndServe(broker.config.Port, handler); err != nil {
		log.Fatal("ListeAndServe:", err)
	}
}

func NewServer(ctx context.Context, config *Config) (*Broker, error) {
	if config.Port == "" {
		return nil, errors.New("port is required")
	}

	if config.JWTSecret == "" {
		return nil, errors.New("jwt secret is required")
	}

	if config.DatabaseUrl == "" {
		return nil, errors.New("database url is required")
	}

	if config.RowsDefault == "" {
		return nil, errors.New("rows default is required")
	}
	if _, err := strconv.ParseUint(config.RowsDefault, 10, 64); err != nil {
		return nil, errors.New("rows default value is invalid")
	}

	broker := &Broker{
		config: config,
		hub:    websocket.NewHub(),
		router: mux.NewRouter(),
	}
	return broker, nil
}
