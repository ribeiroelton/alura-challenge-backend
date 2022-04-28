package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/ribeiroelton/alura-challenge-backend/config"
	"github.com/ribeiroelton/alura-challenge-backend/internal/core/domain/usecase"
	"github.com/ribeiroelton/alura-challenge-backend/pkg/logger"
	"github.com/ribeiroelton/alura-challenge-backend/pkg/mailer"
	"github.com/ribeiroelton/alura-challenge-backend/pkg/validator"
	"github.com/ribeiroelton/alura-challenge-backend/repository"
	"github.com/ribeiroelton/alura-challenge-backend/web/ui"
)

//setupServer defines all components required to create an application
func setupServer() *ui.Server {
	//TODO Adopt Viper
	//TODO Analyze Google Wire DI

	//Config
	c := &config.Config{
		ServerAddress: ":8080",
		ConnString:    "mongodb://admin:admin@localhost:27017",
		DatabaseName:  "challenge",
	}

	//Logger
	zap, err := logger.NewZapLogger()
	if err != nil {
		log.Fatalf("error while creating logger, details %v", err)
	}

	//Mailer
	mailer := mailer.NewMailer()

	//Validator
	validator := validator.NewGOValidator()

	//Repositories
	tr, err := repository.NewTransactionRepository(c)
	if err != nil {
		log.Fatalf("error while creating transaction repository, details %v \n", err)
	}

	ir, err := repository.NewImportRepository(c)
	if err != nil {
		log.Fatalf("error while creating import repository, details %v \n", err)
	}

	ur, err := repository.NewUserRepository(c)
	if err != nil {
		log.Fatalf("error while creating user repository, details %v \n", err)
	}

	//Services
	tsConfig := usecase.TransactionServiceConfig{
		Log:           zap,
		TransactionDB: tr,
		ImportDB:      ir,
		Validator:     validator,
	}
	ts := usecase.NewTransactionService(tsConfig)

	usConfig := usecase.UserServiceConfig{
		Log:    zap,
		DB:     ur,
		Mailer: mailer,
	}
	us := usecase.NewUserService(usConfig)

	//Handlers
	thConfig := &ui.TransactionsHandlerConfig{
		Service: ts,
		Log:     zap,
	}
	th := ui.NewTransactionsHandler(thConfig)

	uhConfig := &ui.UserHandlerConfig{
		Service: us,
		Log:     zap,
	}
	uh := ui.NewUserHandler(uhConfig)

	//Server
	serverConfig := ui.ServerConfig{
		Config: c,
		Log:    zap,
	}
	s := ui.NewServer(serverConfig)

	s.Srv.GET("/upload", th.GetUpload)
	s.Srv.POST("/upload", th.PostUpload)
	s.Srv.GET("/users", uh.GetUsers)
	s.Srv.GET("/users-edit", uh.GetUsersEdit)
	s.Srv.POST("/users-edit", uh.PostUsersEdit)

	return s
}

func main() {
	s := setupServer()

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
