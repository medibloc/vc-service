package kms

import (
	"net/http"

	"github.com/hyperledger/aries-framework-go/pkg/framework/context"
	"github.com/labstack/echo/v4"

	log "github.com/sirupsen/logrus"
)

type service struct {
	command *command
}

func RegisterHandlers(router *echo.Echo, ariesProvider *context.Provider) {
	s := &service{
		command: newCommand(ariesProvider),
	}

	router.POST("/keys/create", s.createKeySetHandler)
}

func (s *service) createKeySetHandler(ctx echo.Context) error {
	request := &CreateKeySetRequest{}
	if err := ctx.Bind(request); err != nil {
		log.Error(err)
		return err
	}

	response, err := s.command.CreateKeySet(request)
	if err != nil {
		log.Error(err)
		return err
	}

	return ctx.JSON(http.StatusCreated, response)
}
