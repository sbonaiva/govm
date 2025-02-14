package api_test

import (
	"context"
	"errors"
	"os"
	"testing"

	"github.com/sbonaiva/govm/internal/api"
	"github.com/sbonaiva/govm/internal/domain"
	"github.com/sbonaiva/govm/internal/handler"
	"github.com/sbonaiva/govm/internal/test"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/suite"
)

type uninstallCmdSuite struct {
	suite.Suite
	ctx     context.Context
	handler *handler.UninstallHandlerMock
	cmd     *cobra.Command
}

func TestUninstallCmd(t *testing.T) {
	suite.Run(t, new(uninstallCmdSuite))
}

func (r *uninstallCmdSuite) SetupTest() {
	r.ctx = context.Background()
	r.handler = new(handler.UninstallHandlerMock)
	r.cmd = api.NewUninstallCmd(r.ctx, r.handler)
}

func (r *uninstallCmdSuite) TearDownTest() {
	r.handler.AssertExpectations(r.T())
}

func (r *uninstallCmdSuite) TestSuccess() {
	// Arrange
	r.handler.On("Handle", r.ctx, &domain.Uninstall{}).Return(nil)

	input := []byte("y")
	rdr, wtr, err := os.Pipe()
	if err != nil {
		r.T().Fatal(err)
	}

	_, err = wtr.Write(input)
	if err != nil {
		r.T().Fatal(err)
	}
	wtr.Close()

	defer func(v *os.File) { os.Stdin = v }(os.Stdin)
	os.Stdin = rdr

	// Act
	output, err := test.CaptureOutput(func() error {
		r.cmd.Run(r.cmd, []string{})
		return nil
	})

	// Assert
	r.NoError(err)
	r.Equal("Confirm uninstall current Go version? (y/n): Go uninstalled successfully!\n", output)
}

func (r *uninstallCmdSuite) TestError() {
	// Arrange
	r.handler.On("Handle", r.ctx, &domain.Uninstall{}).Return(errors.New("uninstall error"))

	input := []byte("y")
	rdr, wtr, err := os.Pipe()
	if err != nil {
		r.T().Fatal(err)
	}

	_, err = wtr.Write(input)
	if err != nil {
		r.T().Fatal(err)
	}
	wtr.Close()

	defer func(v *os.File) { os.Stdin = v }(os.Stdin)
	os.Stdin = rdr

	// Act
	output, err := test.CaptureOutput(func() error {
		r.cmd.Run(r.cmd, []string{})
		return nil
	})

	// Assert
	r.NoError(err)
	r.Equal("Confirm uninstall current Go version? (y/n): Error: uninstall error\n", output)
}

func (r *uninstallCmdSuite) TestAbort() {
	// Arrange
	input := []byte("n")
	rdr, wtr, err := os.Pipe()
	if err != nil {
		r.T().Fatal(err)
	}

	_, err = wtr.Write(input)
	if err != nil {
		r.T().Fatal(err)
	}
	wtr.Close()

	defer func(v *os.File) { os.Stdin = v }(os.Stdin)
	os.Stdin = rdr

	// Act
	output, err := test.CaptureOutput(func() error {
		r.cmd.Run(r.cmd, []string{})
		return nil
	})

	// Assert
	r.NoError(err)
	r.Equal("Confirm uninstall current Go version? (y/n): Uninstall aborted by user\n", output)
}
