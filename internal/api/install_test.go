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

type installCmdSuite struct {
	suite.Suite
	ctx     context.Context
	handler *handler.InstallHandlerMock
	cmd     *cobra.Command
}

func TestInstallCmd(t *testing.T) {
	suite.Run(t, new(installCmdSuite))
}

func (r *installCmdSuite) SetupTest() {
	r.ctx = context.Background()
	r.handler = new(handler.InstallHandlerMock)
	r.cmd = api.NewInstallCmd(r.ctx, r.handler)
}

func (r *installCmdSuite) TearDownTest() {
	r.handler.AssertExpectations(r.T())
}

func (r *installCmdSuite) TestSuccess() {
	// Arrange
	r.handler.On("Handle", r.ctx, &domain.Action{Version: "1.15.0"}).Return(nil)

	// Act
	output, _ := test.CaptureOutput(func() error {
		r.cmd.Run(r.cmd, []string{"1.15.0"})
		return nil
	})

	// Assert
	r.Equal("Go version \"1.15.0\" installed successfully!\nPlease, reopen your terminal to start using new version.\n", output)
}

func (r *installCmdSuite) TestErrorHandling() {
	// Arrange
	r.handler.On("Handle", r.ctx, &domain.Action{Version: "1.24.0"}).Return(errors.New("install error"))

	// Act
	output, _ := test.CaptureOutput(func() error {
		r.cmd.Run(r.cmd, []string{"1.24.0"})
		return nil
	})

	// Assert
	r.Equal("install error\n", output)
}

func (r *installCmdSuite) TestInvalidArguments() {
	// Arrange
	r.cmd.SetArgs([]string{})

	// Act
	err := r.cmd.Execute()

	// Assert
	r.Error(err)
	r.EqualError(err, "accepts 1 arg(s), received 0")
}
