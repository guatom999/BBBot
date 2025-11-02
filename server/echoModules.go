package server

import (
	echoHandlers "github.com/guatom999/BBBot/modules/botREST/botRESTHandlers"
	echoUseCases "github.com/guatom999/BBBot/modules/botREST/botRESTUseCases"
)

func (s *server) BotRestModules() {
	botModulesUsecase := echoUseCases.NewEchoUserCase(s.dg, s.cfg)
	botModulesHanlder := echoHandlers.NewEchoHandler(botModulesUsecase)

	router := s.app.Group("/bot")

	router.POST("/send", botModulesHanlder.Send)

}
