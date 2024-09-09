package service_test

import (
	"testing"

	"github.com/sanderdescamps/go-billing-api/gobilling/service"
)

func TestRandemString(t *testing.T) {
	randString := service.RandomString(16)
	t.Logf("Random String: %s", randString)
}
