package server

import (
	"net/http"

	echoHandlers "github.com/guatom999/BBBot/modules/botREST/botRESTHandlers"
	echoUseCases "github.com/guatom999/BBBot/modules/botREST/botRESTUseCases"
	"github.com/labstack/echo/v4"
)

func (s *server) BotRestModules() {
	botModulesUsecase := echoUseCases.NewEchoUserCase(s.dg, s.cfg)
	botModulesHanlder := echoHandlers.NewEchoHandler(botModulesUsecase)

	// Health check endpoint
	s.app.GET("/", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{
			"status":  "ok",
			"service": "BBBot",
		})
	})

	router := s.app.Group("/bot")

	router.POST("/send", botModulesHanlder.Send)

}
