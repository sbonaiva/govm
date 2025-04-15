package api_test

import (
	"context"
	"errors"
	"testing"

	"github.com/sbonaiva/govm/internal/api"
	"github.com/sbonaiva/govm/internal/domain"
	"github.com/sbonaiva/govm/internal/handler"
	"github.com/sbonaiva/govm/internal/test"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/suite"
)

type updateCmdSuite struct {
	suite.Suite
	ctx     context.Context
	handler *handler.UpdateHandlerMock
	cmd     *cobra.Command
}

func TestUpdateCmd(t *testing.T) {
	suite.Run(t, new(updateCmdSuite))
}

func (r *updateCmdSuite) SetupTest() {
	r.ctx = context.Background()
	r.handler = new(handler.UpdateHandlerMock)
	r.cmd = api.NewUpdateCmd(r.ctx, r.handler)
}

func (r *updateCmdSuite) TearDownTest() {
	r.handler.AssertExpectations(r.T())
}

func (r *updateCmdSuite) TestSuccess() {
	// Arrange
	r.handler.On("Handle", r.ctx, &domain.Action{UpdateStrategy: domain.PatchStrategy}).Return("1.15.1", nil)

	// Act
	output, _ := test.CaptureOutput(func() error {
		r.cmd.Run(r.cmd, []string{"1.15.0"})
		return nil
	})

	// Assert
	r.Equal("Go updated to version \"1.15.1\" successfully!\nPlease, reopen your terminal to start using new version.\n", output)
}

func (r *updateCmdSuite) TestErrorHandling() {
	// Arrange
	r.handler.On("Handle", r.ctx, &domain.Action{UpdateStrategy: domain.PatchStrategy}).Return("", errors.New("update error"))

	// Act
	output, _ := test.CaptureOutput(func() error {
		r.cmd.Run(r.cmd, []string{"1.15.0"})
		return nil
	})

	// Assert
	r.Equal("update error\n", output)
}
