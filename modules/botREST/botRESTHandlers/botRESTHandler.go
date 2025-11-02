package botRESTHandlers

import (
	"errors"
	"net/http"

	"github.com/guatom999/BBBot/modules/botREST"
	"github.com/guatom999/BBBot/modules/botREST/botRESTUseCases"
	"github.com/labstack/echo/v4"
)

type (
	EchoHandlerInterface interface {
		Send(c echo.Context) error
	}

	echoHandler struct {
		echoUseCase botRESTUseCases.EchoServerInterface
	}
)

func NewEchoHandler(echoUseCase botRESTUseCases.EchoServerInterface) EchoHandlerInterface {
	return &echoHandler{echoUseCase: echoUseCase}
}

func (h *echoHandler) Send(c echo.Context) error {

	req := new(botREST.InCommingMessage)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, errors.New("invalid request body"))
	}

	if err := h.echoUseCase.Send(req); err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, "send message success")
}
