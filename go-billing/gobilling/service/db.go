package service

import (
	"sync"

	db "github.com/sonyarouje/simdb"
)

type BillingDB struct {
	db          db.Driver
	dbLock      *sync.Mutex
	CostTypes   *CostTypeService
	CostCenters *CostCenterService
	Resources   *ResourceService
	Users       *UserService
	Roles       *RoleService
}

func (db *BillingDB) GetDB() *db.Driver {
	db.dbLock.Lock()
	return &db.db
}

func NewTreeDB(path string) *BillingDB {
	driver, err := db.New(path)
	mu := sync.Mutex{}
	if err != nil {
		return nil
	}

	billingDB := BillingDB{
		db:     *driver,
		dbLock: &mu,
	}

	billingDB.Resources = &ResourceService{billingDB: &billingDB, dbLock: &mu}
	billingDB.CostTypes = &CostTypeService{billingDB: &billingDB, dbLock: &mu}
	billingDB.CostCenters = &CostCenterService{billingDB: &billingDB, dbLock: &mu}
	billingDB.Users = &UserService{billingDB: &billingDB, dbLock: &mu}
	billingDB.Roles = &RoleService{billingDB: &billingDB, dbLock: &mu}
	return &billingDB
}
