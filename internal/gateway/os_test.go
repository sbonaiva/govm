package gateway_test

import (
	"os"
	"testing"

	"github.com/sbonaiva/govm/internal/gateway"
	"github.com/stretchr/testify/suite"
)

type osGatewaySuite struct {
	suite.Suite
	gateway gateway.OsGateway
}

func TestOsGateway(t *testing.T) {
	suite.Run(t, new(osGatewaySuite))
}

func (r *osGatewaySuite) SetupTest() {
	r.gateway = gateway.NewOsGateway()
}

func (r *osGatewaySuite) TestGetUserHomeDir() {
	usrHomeDir, err := r.gateway.GetUserHomeDir()
	r.NoError(err)
	r.NotEmpty(usrHomeDir)
}

func (r *osGatewaySuite) TestStat() {
	fi, err := r.gateway.Stat(".")
	r.NoError(err)
	r.NotNil(fi)
}

func (r *osGatewaySuite) TestCreateAndRemoveDir() {
	err := r.gateway.CreateDir("xpto", 0755)
	r.NoError(err)
	err = r.gateway.RemoveDir("xpto")
	r.NoError(err)
}

func (r *osGatewaySuite) TestCreateAndRemoveFile() {
	_, err := r.gateway.CreateFile("xpto")
	r.NoError(err)
	err = r.gateway.RemoveFile("xpto")
	r.NoError(err)
}

func (r *osGatewaySuite) TestOpenFile() {
	file, err := r.gateway.OpenFile("os_test.go")
	r.NoError(err)
	r.NotNil(file)
}

func (r *osGatewaySuite) TestReadFile() {
	data, err := r.gateway.ReadFile("os_test.go")
	r.NoError(err)
	r.NotEmpty(data)
}

func (r *osGatewaySuite) TestWriteFile() {
	err := r.gateway.WriteFile("xpto", []byte("test"), 0755)
	r.NoError(err)
	err = r.gateway.RemoveFile("xpto")
	r.NoError(err)
}

func (r *osGatewaySuite) TestGetEnv() {
	os.Setenv("XPTO", "test")
	env := r.gateway.GetEnv("XPTO")
	r.Equal("test", env)
}
