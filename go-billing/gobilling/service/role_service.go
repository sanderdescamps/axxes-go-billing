package service

import (
	"fmt"
	"sync"

	"github.com/sanderdescamps/go-billing-api/gobilling/log"
	"github.com/sanderdescamps/go-billing-api/gobilling/model"
)

type RoleService struct {
	billingDB *BillingDB
	dbLock    *sync.Mutex
}

func (serv *RoleService) Get(id string) (*model.Role, error) {
	var role model.Role
	err := serv.billingDB.GetDB().Open(model.Role{}).Where(model.PARAM_ROLE_ID, "=", id).First().AsEntity(&role)
	serv.dbLock.Unlock()
	if err != nil {
		return nil, err
	}
	return &role, nil
}

func (serv *RoleService) GetByName(name string) (*model.Role, error) {
	var role model.Role
	err := serv.billingDB.GetDB().Open(model.Role{}).Where("name", "=", name).First().AsEntity(&role)
	serv.dbLock.Unlock()
	if err != nil {
		return nil, err
	}
	return &role, nil
}

func (serv *RoleService) Exits(id string) bool {
	defer serv.dbLock.Unlock()
	l := len(serv.billingDB.GetDB().Open(model.Role{}).Where(model.PARAM_ROLE_ID, "=", id).Get().RawArray()) > 0
	return l
}

func (serv *RoleService) ExitsByName(name string) bool {
	defer serv.dbLock.Unlock()
	l := len(serv.billingDB.GetDB().Open(model.Role{}).Where("name", "=", name).Get().RawArray()) > 0
	return l
}

func (serv *RoleService) Create(name string, permissions []string) (*model.Role, error) {
	defer serv.dbLock.Unlock()
	if serv.Exits(name) {
		return nil, fmt.Errorf("role already exist")
	}
	role := model.NewRole(name, permissions)

	log.Debugf("New Role: %s\n", role.ToString())
	err := serv.billingDB.GetDB().Insert(role)
	if err != nil {
		return nil, err
	}
	return role, nil
}
