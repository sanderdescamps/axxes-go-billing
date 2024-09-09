package service_test

import (
	"fmt"
	"testing"

	"github.com/sanderdescamps/go-billing-api/gobilling/service"
)

func TestResource(t *testing.T) {

	testId := service.RandomString(6)

	t.Run("SingleResource", func(t *testing.T) {
		name := fmt.Sprintf("name-single-%s", testId)
		t.Cleanup(func() {
			DeleteResourceIfExists(t, name)
		})
		_, err := billingDB.Resources.Create(name, "this is a test", nil, 1)
		if err != nil {
			t.Error(err)
		}

	})

	t.Run("ResourceWithCostType", func(t *testing.T) {
		cType1Name := fmt.Sprintf("type-1-%s", testId)
		cType2Name := fmt.Sprintf("type-2-%s", testId)
		rName := fmt.Sprintf("withCostType-%s", testId)

		t.Cleanup(func() {
			DeleteResourceIfExists(t, rName)
			DeleteCostTypeIfExists(t, cType1Name)
			DeleteCostTypeIfExists(t, cType2Name)
		})

		newResource, err := billingDB.Resources.Create(rName, "this is a test", nil, 1)
		if err != nil {
			t.Error(err)
		}

		cType1 := CreateTestCostType(t, cType1Name, 1.2, 0.0)
		cType2 := CreateTestCostType(t, cType2Name, 0.0, 3.2)

		billingDB.Resources.AddCostType(newResource.ResourceID, cType1.TypeID)
		billingDB.Resources.AddCostType(newResource.ResourceID, cType2.TypeID)

	})
}

//---------------------------------------------------------------------

func DeleteResourceIfExists(t *testing.T, name string) {
	if resource, err := billingDB.Resources.GetByName(name); err == nil {
		err := billingDB.Resources.Delete(resource)
		if err != nil {
			t.Errorf("failed to delete resource (name=%s): %v", name, err)
		}
	} else {
		t.Errorf("resource (name=%s) does not exist", name)
	}
}
