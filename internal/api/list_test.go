package api_test

import (
	"context"
	"errors"
	"testing"

	"github.com/sbonaiva/govm/internal/api"
	"github.com/sbonaiva/govm/internal/handler"
	"github.com/sbonaiva/govm/internal/test"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/suite"
)

type listCmdSuite struct {
	suite.Suite
	ctx     context.Context
	handler *handler.ListHandlerMock
	cmd     *cobra.Command
}

func TestListCmd(t *testing.T) {
	suite.Run(t, new(listCmdSuite))
}

func (r *listCmdSuite) SetupTest() {
	r.ctx = context.Background()
	r.handler = new(handler.ListHandlerMock)
	r.cmd = api.NewListCmd(r.ctx, r.handler)
}

func (r *listCmdSuite) TearDownTest() {
	r.handler.AssertExpectations(r.T())
}

func (r *listCmdSuite) TestSuccess() {
	// Arrange
	r.handler.On("Handle", r.ctx).Return(nil)
	r.cmd.SetArgs([]string{})

	// Act
	output, _ := test.CaptureOutput(func() error {
		r.cmd.Run(r.cmd, []string{})
		return nil
	})

	// Assert
	r.Empty(output)
}

func (r *listCmdSuite) TestErrorHandling() {
	// Arrange
	r.handler.On("Handle", r.ctx).Return(errors.New("list error"))
	r.cmd.SetArgs([]string{})

	// Act
	output, _ := test.CaptureOutput(func() error {
		r.cmd.Run(r.cmd, []string{})
		return nil
	})

	// Assert
	r.Equal("Error: list error\n", output)
}
