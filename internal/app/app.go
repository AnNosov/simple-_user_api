package app

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/AnNosov/simple_user_api/config"
	v1 "github.com/AnNosov/simple_user_api/internal/controller/http/v1"
	"github.com/AnNosov/simple_user_api/internal/usecase"
	"github.com/AnNosov/simple_user_api/internal/usecase/repo"
	"github.com/AnNosov/simple_user_api/pkg/httpserver"
	"github.com/AnNosov/simple_user_api/pkg/postgres"
	"github.com/go-chi/chi/v5"
)

func Run(cfg *config.Config) {

	pg, err := postgres.New(&cfg.Postgres)
	if err != nil {
		log.Fatal(err)
	}
	defer pg.Close()

	userActionUseCase := usecase.New(
		*repo.New(pg),
	)

	handler := chi.NewRouter()
	v1.NewRouter(handler, *userActionUseCase)
	httpServer := httpserver.New(handler, httpserver.Port(cfg.Http.Port))

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM, syscall.SIGINT) // у меня работает только с SIGINT (ctrl + c)

	select {
	case s := <-interrupt:
		log.Println("app - Run - signal: " + s.String())
	case err = <-httpServer.Notify():
		log.Println("app - Run - httpServer.Notify: ", err.Error())
	}

	err = httpServer.Shutdown()
	if err != nil {
		log.Println("app - Run - httpServer.Shutdown: ", err.Error())
	}
}
