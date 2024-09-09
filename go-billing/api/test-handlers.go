package api

import (
	"fmt"
	"net/http"
)

func test1(w http.ResponseWriter, r *http.Request) {
	resource_vm1, _ := billingDB.Resources.GetByName("vm1")

	fmt.Printf("resource vm1: %s\n", resource_vm1.ToString())

	cost := billingDB.Resources.GetResourceCost(resource_vm1)

	fmt.Printf("resource vm1 cost: %s\n", cost.ToString())

}
