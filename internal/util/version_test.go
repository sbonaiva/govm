package util_test

import (
	"fmt"
	"runtime"
	"testing"

	"github.com/sbonaiva/govm/internal/util"
	"github.com/stretchr/testify/assert"
)

func TestVersion(t *testing.T) {
	assert.Equal(
		t,
		fmt.Sprintf("0.0.1 %s/%s", runtime.GOOS, runtime.GOARCH),
		util.GoVersionManager,
	)
}
