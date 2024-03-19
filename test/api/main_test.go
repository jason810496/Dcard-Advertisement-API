package api_test

import (
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	SetupPackage()
	code := m.Run()
	TeardownPackage()
	os.Exit(code)
}
