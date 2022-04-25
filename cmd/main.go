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

//configureServer defines all components required to create an application
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

	serverConfig := ui.ServerConfig{
		Config: c,
		Log:    zap,
	}
	s := ui.NewServer(serverConfig)

	tr, err := repository.NewTransactionRepository(c)
	if err != nil {
		log.Fatalf("error while creating transaction repository, details %v \n", err)
	}

	tsConfig := usecase.TransactionServiceConfig{
		Log: zap,
		DB:  tr,
	}
	ts := usecase.NewTransactionService(tsConfig)

	thConfig := &ui.TransactionsHandlerConfig{
		Service: ts,
		Log:     zap,
		Srv:     s.Srv,
	}
	ui.NewTransactionsHandler(thConfig)

	ur, err := repository.NewUserRepository(c)
	if err != nil {
		log.Fatalf("error while creating user repository, details %v \n", err)
	}

	usConfig := usecase.UserServiceConfig{
		Log: zap,
		DB:  ur,
	}

	us := usecase.NewUserService(usConfig)

	uhConfig := &ui.UserHandlerConfig{
		Service: us,
		Log:     zap,
		Srv:     s.Srv,
	}
	ui.NewUserHandler(uhConfig)

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
