package service

import (
	"fmt"
	"sync"

	"github.com/sanderdescamps/go-billing-api/gobilling/log"
	"github.com/sanderdescamps/go-billing-api/gobilling/model"
)

type CostCenterService struct {
	// db        *db.Driver
	billingDB *BillingDB
	dbLock    *sync.Mutex
}

func (serv *CostCenterService) GetAll() ([]model.CostCenter, error) {
	defer serv.dbLock.Unlock()
	var costCenter []model.CostCenter
	err := serv.billingDB.GetDB().Open(model.CostCenter{}).Get().AsEntity(&costCenter)
	if err != nil {
		return nil, err
	}
	return costCenter, nil
}

func (serv *CostCenterService) Get(id string) (*model.CostCenter, error) {
	defer serv.dbLock.Unlock()
	var costCenter model.CostCenter
	err := serv.billingDB.GetDB().Open(model.CostCenter{}).Where(model.PARAM_COSTCENTER_ID, "=", id).First().AsEntity(&costCenter)
	if err != nil {
		return nil, err
	}
	return &costCenter, nil
}

func (serv *CostCenterService) GetByName(name string) (*model.CostCenter, error) {
	defer serv.dbLock.Unlock()
	var costCenter model.CostCenter
	err := serv.billingDB.GetDB().Open(model.CostCenter{}).Where("name", "=", name).First().AsEntity(&costCenter)
	if err != nil {
		return nil, err
	}
	return &costCenter, nil
}

func (serv *CostCenterService) Exits(id string) bool {
	defer serv.dbLock.Unlock()
	return len(serv.billingDB.GetDB().Open(model.CostCenter{}).Where(model.PARAM_COSTCENTER_ID, "=", id).Get().RawArray()) > 0
}

func (serv *CostCenterService) ExitsByName(name string) bool {
	defer serv.dbLock.Unlock()
	return len(serv.billingDB.GetDB().Open(model.CostCenter{}).Where("name", "=", name).Get().RawArray()) > 0
}

func (serv *CostCenterService) Create(name string, description string) (*model.CostCenter, error) {
	if serv.ExitsByName(name) {
		return nil, fmt.Errorf("cost-center already exist")
	}
	newCostCenter := model.NewCostCenter(name, description)

	defer serv.dbLock.Unlock()
	log.Debugf("New cost-center: %s\n", newCostCenter.ToString())
	err := serv.billingDB.GetDB().Insert(newCostCenter)
	if err != nil {
		return nil, err
	}
	return newCostCenter, nil
}

func (serv *CostCenterService) Update(r *model.CostCenter) error {
	defer serv.dbLock.Unlock()
	err := serv.billingDB.GetDB().Update(r)
	if err != nil {
		return err
	}
	return nil
}

func (serv *CostCenterService) DeleteById(id string) error {
	toDel, err := serv.Get(id)
	if err != nil {
		return err
	}
	return serv.Delete(toDel)
}

func (serv *CostCenterService) Delete(f *model.CostCenter) error {
	defer serv.dbLock.Unlock()
	err := serv.billingDB.GetDB().Delete(f)
	if err != nil {
		return err
	}
	return nil
}

func (serv *CostCenterService) GetTotalCost(costCenterId string) (*model.Cost, error) {
	resources, err := serv.billingDB.Resources.GetAllForCostCenter(costCenterId)
	if err != nil {
		return nil, err
	}
	totalCost := model.NewCost()
	for _, r := range resources {
		resourceCost := serv.billingDB.Resources.GetResourceCost(&r)
		totalCost.Add(resourceCost)
	}
	return totalCost, nil
}
