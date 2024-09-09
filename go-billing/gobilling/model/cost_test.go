package model_test

import (
	"testing"

	"github.com/sanderdescamps/go-billing-api/gobilling/model"
)

func TestCost(t *testing.T) {
	t.Run("Add", func(t *testing.T) {
		a := model.Cost{Fixed: 2.6, CostPerSec: 5.0}
		b := model.Cost{Fixed: 2.4, CostPerSec: 5.0}
		a.Add(b)
		if !a.Equal(model.Cost{Fixed: 5.0, CostPerSec: 10.0}) {
			t.Error("Adding costs incorrect")
		}
	})
	t.Run("Multiply", func(t *testing.T) {
		a := model.Cost{Fixed: 2.6, CostPerSec: 5.0}
		b := 5.5
		a.Multiply(b)
		if !a.Equal(model.Cost{Fixed: 2.6 * 5.5, CostPerSec: 5.0 * 5.5}) {
			t.Error("Multiply costs incorrect")
		}
	})
}
