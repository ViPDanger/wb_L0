package handler

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type HttpClient struct {
	httpServer *http.Server
}

// Запуск http Клиента
func (s *HttpClient) Run(addres string, port string, natsHandler NatsHanlder) error {
	// Отстройка параметров
	ctx_HTTP_Shutdown, cancel_on_http := context.WithCancel(context.Background())
	ctx, stop := signal.NotifyContext(ctx_HTTP_Shutdown, os.Interrupt, os.Kill, syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	mux := http.NewServeMux()
	SetupHttpHandlers(mux, cancel_on_http, natsHandler)
	s.httpServer = &http.Server{
		Addr:           addres + ":" + port,
		Handler:        mux,
		MaxHeaderBytes: 1 << 20,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   20 * time.Second,
	}

	// Запуск сервера
	go func() {
		if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Fatal Error: %v", err)
		}
	}()
	time.Sleep(1 * time.Second)
	log.Print("Sever starter on addres: ", s.httpServer.Addr)

	// Выключение сервера
	<-ctx.Done()
	log.Println("Shutting down")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	s.httpServer.Close()
	err := s.httpServer.Shutdown(shutdownCtx)
	if err != nil {
		log.Println("Shutting down the server: ", err)
		return err
	}

	return nil
}
