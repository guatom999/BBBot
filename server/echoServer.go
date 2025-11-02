package server

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/guatom999/BBBot/config"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"golang.org/x/net/context"
)

type (
	Server interface {
		Start(pctx context.Context)
		GracefulShutdown(pctx context.Context, close <-chan os.Signal)
	}

	server struct {
		dg  *discordgo.Session
		app *echo.Echo
		cfg *config.Config
	}
)

func NewEchoServer(dg *discordgo.Session, cfg *config.Config) Server {
	return &server{
		dg:  dg,
		app: echo.New(),
		cfg: cfg,
	}
}

func (s *server) GracefulShutdown(pctx context.Context, close <-chan os.Signal) {
	<-close

	ctx, cancel := context.WithTimeout(pctx, time.Second*10)
	defer cancel()

	if err := s.app.Shutdown(ctx); err != nil {
		log.Fatalf("Failed to shutdown server...")
		panic(err)
	}

	log.Println("Shutting Down Server....")

}

func (s *server) Start(pctx context.Context) {

	s.app.Use(middleware.TimeoutWithConfig(middleware.TimeoutConfig{
		Skipper:      middleware.DefaultSkipper,
		ErrorMessage: "Error: Request Timeout",
		Timeout:      time.Second * 10,
	}))

	//Cors
	s.app.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		Skipper:      middleware.DefaultSkipper,
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.GET, echo.POST, echo.PUT, echo.DELETE, echo.PATCH},
	}))

	s.app.Use(middleware.Logger())

	close := make(chan os.Signal, 1)
	signal.Notify(close, syscall.SIGINT, syscall.SIGTERM)

	go s.GracefulShutdown(pctx, close)

	s.BotRestModules()

	if err := s.app.Start(s.cfg.EchoApp.Port); err != nil {
		log.Fatalf("Failed to shutdown:%v", err)
	}

}
