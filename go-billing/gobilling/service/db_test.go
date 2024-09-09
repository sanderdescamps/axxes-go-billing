package service_test

import (
	"github.com/sanderdescamps/go-billing-api/gobilling/service"
)

var billingDB *service.BillingDB

func init() {
	billingDB = service.NewTreeDB("./data")
}
