package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/ribeiroelton/alura-challenge-backend/config"
	"github.com/ribeiroelton/alura-challenge-backend/internal/core/domain/usecase"
	"github.com/ribeiroelton/alura-challenge-backend/pkg/logger"
	"github.com/ribeiroelton/alura-challenge-backend/repository"
	"github.com/ribeiroelton/alura-challenge-backend/web/ui"
)

func configureServer() *ui.Server {
	//TODO Adopt Viper
	c := &config.Config{
		ServerAddress: ":8080",
		ConnString:    "mongodb://admin:admin@localhost:27017",
		DatabaseName:  "challenge",
	}

	zap, err := logger.NewZapLogger()
	if err != nil {
		log.Fatalf("error while creating logger, details %v", err)
	}

	db, err := repository.NewMongoDB(c)
	if err != nil {
		log.Fatalf("error while creating database, details %v \n", err)
	}

	svcConfig := usecase.TransactionServiceConfig{
		Log: zap,
		DB:  db,
	}
	svc := usecase.NewTransactionService(svcConfig)

	serverConfig := ui.ServerConfig{
		Config: c,
		Log:    zap,
	}
	s := ui.NewServer(serverConfig)

	thConfig := ui.TransactionsHandlerConfig{
		Service: svc,
		Log:     zap,
		Srv:     s.Srv,
	}
	ui.NewTransactionsHandler(thConfig)

	return s
}

func main() {
	s := configureServer()

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)

	done := make(chan bool, 1)

	go func() {
		if err := s.StartServer(); err != nil {
			log.Fatalf("error while starting UI server, details %v, \n", err)
		}
	}()

	go func() {
		log.Println("Press CTRL+C to stop")
		<-ch
		log.Println("Stopped")
		done <- true
	}()

	<-done

}
