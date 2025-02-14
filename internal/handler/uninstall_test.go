package handler_test

import (
	"context"
	"testing"

	"github.com/sbonaiva/govm/internal/handler"
	"github.com/stretchr/testify/suite"
)

type uninstallHandlerSuite struct {
	suite.Suite
	ctx     context.Context
	handler handler.UninstallHandler
}

func TestUninstallHandler(t *testing.T) {
	suite.Run(t, new(uninstallHandlerSuite))
}

func (s *uninstallHandlerSuite) SetupTest() {
	s.ctx = context.Background()
	s.handler = handler.NewUninstall()
}
