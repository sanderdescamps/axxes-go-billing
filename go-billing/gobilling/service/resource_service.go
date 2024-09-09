package service

import (
	"fmt"
	"sync"

	"github.com/sanderdescamps/go-billing-api/gobilling/log"
	"github.com/sanderdescamps/go-billing-api/gobilling/model"
)

type ResourceService struct {
	// db        *db.Driver
	billingDB *BillingDB
	dbLock    *sync.Mutex
}

func (serv *ResourceService) GetAll() ([]model.Resource, error) {
	defer serv.dbLock.Unlock()
	var resources []model.Resource
	err := serv.billingDB.GetDB().Open(model.Resource{}).Get().AsEntity(&resources)
	if err != nil {
		return nil, err
	}
	return resources, nil
}

func (serv *ResourceService) Get(id string) (*model.Resource, error) {
	defer serv.dbLock.Unlock()
	var resource model.Resource
	err := serv.billingDB.GetDB().Open(model.Resource{}).Where(model.PARAM_RESOURCE_ID, "=", id).First().AsEntity(&resource)
	if err != nil {
		return nil, err
	}
	return &resource, nil
}

func (serv *ResourceService) GetByName(name string) (*model.Resource, error) {
	defer serv.dbLock.Unlock()
	var resource model.Resource
	err := serv.billingDB.GetDB().Open(model.Resource{}).Where("name", "=", name).First().AsEntity(&resource)
	if err != nil {
		return nil, err
	}
	return &resource, nil
}

func (serv *ResourceService) GetAllForCostCenter(costCenterId string) ([]model.Resource, error) {
	defer serv.dbLock.Unlock()
	var resources []model.Resource
	err := serv.billingDB.GetDB().Open(model.Resource{}).Where("cost_center_id", "=", costCenterId).AsEntity(&resources)
	if err != nil {
		return nil, err
	}
	return resources, nil
}

func (serv *ResourceService) GetResourceCost(resource *model.Resource) model.Cost {
	totalCost := model.NewCost()
	for _, ctId := range resource.CostTypesIds {
		ct, err := serv.billingDB.CostTypes.Get(ctId)
		if err != nil {
			log.Warningf("Invalid cost type")
		}
		totalCost = totalCost.Add(ct.Cost).Multiply(resource.Value)
	}
	return *totalCost
}

func (serv *ResourceService) Exits(id string) bool {
	defer serv.dbLock.Unlock()
	return len(serv.billingDB.GetDB().Open(model.Resource{}).Where(model.PARAM_RESOURCE_ID, "=", id).Get().RawArray()) > 0
}

func (serv *ResourceService) ExitsByName(name string) bool {
	defer serv.dbLock.Unlock()
	return len(serv.billingDB.GetDB().Open(model.Resource{}).Where("name", "=", name).Get().RawArray()) > 0
}

func (serv *ResourceService) Create(name string, description string, costCenterId string, costTypeIds []string, value float64) (*model.Resource, error) {
	if serv.ExitsByName(name) {
		return nil, fmt.Errorf("resource already exist")
	}

	if !serv.billingDB.CostCenters.Exits(costCenterId) {
		return nil, fmt.Errorf("cost-center does not exist")
	}
	newResource := model.NewResource(name, description, costCenterId, value)

	for _, i := range costTypeIds {
		if serv.billingDB.CostTypes.Exits(i) && !newResource.HasCostType(i) {
			newResource.CostTypesIds = append(newResource.CostTypesIds, i)
		}
	}

	defer serv.dbLock.Unlock()
	log.Debugf("New resource: %s\n", newResource.ToString())
	err := serv.billingDB.GetDB().Insert(newResource)
	if err != nil {
		return nil, err
	}
	return newResource, nil
}

func (serv *ResourceService) AddCostType(resourceId string, costTypeId string) error {
	resource, err := serv.Get(resourceId)
	if err != nil {
		return fmt.Errorf("resource not found")
	}
	if !serv.billingDB.CostTypes.Exits(costTypeId) {
		return fmt.Errorf("invalid cost_type id")
	}
	resource.CostTypesIds = append(resource.CostTypesIds, costTypeId)
	serv.Update(resource)
	return nil
}

func (serv *ResourceService) DeleteCostType(resourceId string, costTypeId string) error {
	resource, err := serv.Get(resourceId)
	if err != nil {
		return fmt.Errorf("resource not found")
	}

	newList := []string{}
	for _, i := range resource.CostTypesIds {
		if i != costTypeId {
			newList = append(newList, i)
		}
	}
	resource.CostTypesIds = newList
	serv.Update(resource)
	return nil
}

func (serv *ResourceService) Update(r *model.Resource) error {
	defer serv.dbLock.Unlock()
	err := serv.billingDB.GetDB().Update(r)
	if err != nil {
		return err
	}
	return nil
}

func (serv *ResourceService) DeleteById(id string) error {
	defer serv.dbLock.Unlock()
	toDel := model.Resource{ResourceID: id}
	err := serv.billingDB.GetDB().Delete(toDel)
	if err != nil {
		return err
	}
	return nil
}

func (serv *ResourceService) Delete(f *model.Resource) error {
	defer serv.dbLock.Unlock()
	err := serv.billingDB.GetDB().Delete(f)
	if err != nil {
		return err
	}
	return nil
}
