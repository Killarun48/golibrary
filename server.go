package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "golibrary/docs"
	"golibrary/internal/modules"

	httpSwagger "github.com/swaggo/http-swagger/v2"
)

var pathDB = "./users.db"

type Server struct {
	srv     *http.Server
	users   map[string]string
	sigChan chan os.Signal
}

func NewServer(addr string) *Server {

	server := &Server{
		sigChan: make(chan os.Signal, 1),
		users:   make(map[string]string),
	}

	signal.Notify(server.sigChan, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)

	// инициализируем фасад "библиотеки"
	library := modules.NewLibraryFacade(pathDB)

	// инициализируем маршруты
	r := library.GetRoutes()

	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("doc.json"),
	))

	srv := &http.Server{
		Addr:    addr,
		Handler: r,
	}

	server.srv = srv

	//time.Sleep(1 * time.Second)

	return server
}

func (s *Server) Serve() {
	go func() {
		log.Println("Starting server...")
		if err := s.srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	<-s.sigChan

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := s.srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}

	log.Println("Server stopped gracefully")
}

func (s *Server) Stop() {
	s.sigChan <- syscall.Signal(1)
}
