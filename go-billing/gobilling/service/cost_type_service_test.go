package service_test

import (
	"fmt"
	"testing"

	"github.com/sanderdescamps/go-billing-api/gobilling/model"
	"github.com/sanderdescamps/go-billing-api/gobilling/service"
)

func TestCostType(t *testing.T) {

	testId := service.RandomString(6)

	t.Run("SingleCostType", func(t *testing.T) {
		name := fmt.Sprintf("single-%s", testId)
		t.Cleanup(func() {
			DeleteCostTypeIfExists(t, name)
		})
		CreateTestCostType(t, name, 0.0, 0.0)
	})

}

//---------------------------------------------------------------------

func CreateTestCostType(t *testing.T, name string, costPerSecond float64, fixCost float64) *model.CostType {
	newCostType := model.NewCostType(name, "this is a test", model.Cost{CostPerSec: costPerSecond, Fixed: fixCost})
	err := billingDB.CostTypes.Create(newCostType)
	if err != nil {
		t.Error(err)
	}
	return newCostType
}

func DeleteCostTypeIfExists(t *testing.T, name string) {
	if costType, err := billingDB.CostTypes.GetByName(name); err == nil {
		err := billingDB.CostTypes.Delete(costType)
		if err != nil {
			t.Errorf("failed to delete CostType (name=%s): %v", name, err)
		}
	} else {
		t.Errorf("CostType (name=%s) does not exist", name)
	}
}
