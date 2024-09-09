package service

import (
	"fmt"
	"sync"

	"github.com/sanderdescamps/go-billing-api/gobilling/log"
	"github.com/sanderdescamps/go-billing-api/gobilling/model"
	db "github.com/sonyarouje/simdb"
)

// CostType

type CostTypeService struct {
	// db        *db.Driver
	billingDB *BillingDB
	dbLock    *sync.Mutex
}

func (serv *CostTypeService) GetDB() *db.Driver {
	serv.billingDB.dbLock.Lock()
	return &(serv.billingDB).db
}

func (serv *CostTypeService) Get(id string) (*model.CostType, error) {
	var CostType model.CostType
	err := serv.GetDB().Open(model.CostType{}).Where(model.PARAM_COST_TYPE_ID, "=", id).First().AsEntity(&CostType)
	serv.dbLock.Unlock()
	if err != nil {
		return nil, err
	}
	return &CostType, nil
}

func (serv *CostTypeService) GetAll() ([]model.CostType, error) {
	var CostType []model.CostType
	err := serv.GetDB().Open(model.CostType{}).Get().AsEntity(&CostType)
	serv.dbLock.Unlock()
	if err != nil {
		return nil, err
	}
	return CostType, nil
}

func (serv *CostTypeService) GetByName(name string) (*model.CostType, error) {
	var CostType model.CostType
	err := serv.GetDB().Open(model.CostType{}).Where("name", "=", name).First().AsEntity(&CostType)
	serv.dbLock.Unlock()
	if err != nil {
		return nil, err
	}
	return &CostType, nil
}

func (serv *CostTypeService) ExitsByName(name string) bool {
	defer serv.dbLock.Unlock()
	return len(serv.GetDB().Open(model.CostType{}).Where("name", "=", name).Get().RawArray()) > 0
}

func (serv *CostTypeService) Exits(id string) bool {
	defer serv.dbLock.Unlock()
	return len(serv.GetDB().Open(model.CostType{}).Where(model.PARAM_COST_TYPE_ID, "=", id).Get().RawArray()) > 0
}

func (serv *CostTypeService) Create(s *model.CostType) error {
	if serv.ExitsByName(s.Name) || serv.Exits(s.TypeID) {
		return fmt.Errorf("CostType already exist")
	}

	log.Debugf("New costType: %s\n", s.ToString())
	err := serv.GetDB().Insert(s)
	serv.dbLock.Unlock()
	if err != nil {
		return err
	}
	return nil
}

func (serv *CostTypeService) Update(t *model.CostType) error {
	err := serv.GetDB().Update(t)
	serv.dbLock.Unlock()
	if err != nil {
		return err
	}
	return nil
}

func (serv *CostTypeService) DeleteById(id string) error {
	toDel := model.CostType{
		TypeID: id,
	}
	err := serv.GetDB().Delete(toDel)
	serv.dbLock.Unlock()
	if err != nil {
		return err
	}
	return nil
}

func (serv *CostTypeService) Delete(t *model.CostType) error {
	err := serv.GetDB().Delete(t)
	serv.dbLock.Unlock()
	if err != nil {
		return err
	}
	return nil
}

func (serv *CostTypeService) Save(CostType ...*model.CostType) error {
	for _, s := range CostType {
		if ok := serv.ExitsByName(s.Name); ok {
			err := serv.Update(s)
			if err != nil {
				return err
			}
		} else {
			err := serv.Create(s)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
