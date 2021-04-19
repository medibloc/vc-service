package verifiable

import (
	"net/http"

	log "github.com/sirupsen/logrus"

	"github.com/labstack/echo/v4"

	"github.com/hyperledger/aries-framework-go/pkg/framework/context"
)

type service struct {
	command *command
}

func RegisterHandlers(router *echo.Echo, ariesProvider *context.Provider) {
	s := &service{
		command: newCommand(ariesProvider),
	}

	router.POST("/credentials/issue", s.issueCredentialHandler)
	router.POST("/credentials/verify", s.verifyCredentialHandler)
	router.POST("/credentials/derive", s.deriveCredentialHandler)
	router.POST("/presentations/prove", s.provePresentationHandler)
	router.POST("/presentations/verify", s.verifyPresentationHandler)
}

func (s *service) issueCredentialHandler(ctx echo.Context) error {
	request := &IssueCredentialRequest{}
	if err := ctx.Bind(request); err != nil {
		log.Error(err)
		return err
	}

	vc, err := s.command.IssueCredential(request)
	if err != nil {
		log.Error(err)
		return err
	}

	return ctx.JSON(http.StatusCreated, vc)
}

func (s *service) verifyCredentialHandler(ctx echo.Context) error {
	request := &VerifyCredentialRequest{}
	if err := ctx.Bind(request); err != nil {
		log.Error(err)
		return err
	}

	if err := s.command.VerifyCredential(request); err != nil {
		log.Error(err)
		return err
	}

	return ctx.JSON(http.StatusCreated, "done") //TODO: use a proper response type
}

func (s *service) deriveCredentialHandler(ctx echo.Context) error {
	request := &DeriveCredentialRequest{}
	if err := ctx.Bind(request); err != nil {
		log.Error(err)
		return err
	}

	revealVC, err := s.command.DeriveCredential(request)
	if err != nil {
		log.Error(err)
		return err
	}

	return ctx.JSON(http.StatusCreated, revealVC)
}

func (s *service) provePresentationHandler(ctx echo.Context) error {
	request := &ProvePresentationRequest{}
	if err := ctx.Bind(request); err != nil {
		log.Error(err)
		return err
	}

	vp, err := s.command.ProvePresentation(request)
	if err != nil {
		log.Error(err)
		return err
	}

	return ctx.JSON(http.StatusCreated, vp)
}

func (s *service) verifyPresentationHandler(ctx echo.Context) error {
	request := &VerifyPresentationRequest{}
	if err := ctx.Bind(request); err != nil {
		log.Error(err)
		return err
	}

	if err := s.command.VerifyPresentation(request); err != nil {
		log.Error(err)
		return err
	}

	return ctx.JSON(http.StatusCreated, "done") //TODO: use a proper response type
}
