package service

import (
	"fmt"
	"strings"
	"sync"

	"github.com/sanderdescamps/go-billing-api/gobilling/log"
	"github.com/sanderdescamps/go-billing-api/gobilling/model"
)

type UserService struct {
	billingDB *BillingDB
	dbLock    *sync.Mutex
}

// func (serv *UserService) GetDB() *db.Driver {
// 	serv.dbLock.Lock()
// 	return serv.db
// }

func (serv *UserService) Get(id string) (*model.User, error) {
	var user model.User
	err := serv.billingDB.GetDB().Open(model.User{}).Where(model.PARAM_USER_ID, "=", id).First().AsEntity(&user)
	serv.dbLock.Unlock()
	if err != nil {
		return nil, err
	}
	return &user, nil
}
func (serv *UserService) GetByUsername(username string) (*model.User, error) {
	defer serv.dbLock.Unlock()
	var user model.User
	err := serv.billingDB.GetDB().Open(model.User{}).Where("username", "=", username).First().AsEntity(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (serv *UserService) Exits(id string) bool {
	defer serv.dbLock.Unlock()
	l := len(serv.billingDB.GetDB().Open(model.User{}).Where(model.PARAM_USER_ID, "=", id).Get().RawArray()) > 0
	return l
}

func (serv *UserService) ExitsUsername(username string) bool {
	defer serv.dbLock.Unlock()
	l := len(serv.billingDB.GetDB().Open(model.User{}).Where("username", "=", username).Get().RawArray()) > 0
	return l
}

func (serv *UserService) Create(username string, password string, roleId string) (*model.User, error) {
	if serv.ExitsUsername(username) {
		return nil, fmt.Errorf("user already exist")
	}

	if !serv.billingDB.Roles.Exits(roleId) {
		return nil, fmt.Errorf("invalid role ID")
	}

	user := model.NewUser(username, password, roleId)

	defer serv.dbLock.Unlock()
	log.Debugf("New User: %s\n", user.ToString())
	err := serv.billingDB.GetDB().Insert(user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (serv *UserService) Update(u *model.User) error {
	defer serv.dbLock.Unlock()
	err := serv.billingDB.GetDB().Update(u)
	if err != nil {
		return err
	}
	return nil
}

func (serv *UserService) HasPermissions(userId string, p []string) (bool, error) {
	user, err := serv.Get(userId)
	if err != nil {
		return false, err
	}

	var userRole *model.Role
	if user.RoleId != "" {
		userRole, err = serv.billingDB.Roles.Get(user.RoleId)
		if err != nil {
			return false, err
		}
	} else {
		return false, fmt.Errorf("user has no role")
	}

	if strings.ToLower(userRole.Name) == "admin" {
		log.Errorf("it is the admin")
		return true, nil
	} else {
		log.Errorf("it is the %s", userRole.Name)
	}
	for _, i := range p {
		if !userRole.IsAllowed(i) {
			return false, nil
		}
	}
	return true, nil
}

func (serv *UserService) SetNewApiToken(id string) (string, error) {
	user, err := serv.Get(id)
	if err != nil {
		return "", err
	}

	token, err := user.SetNewApiToken()
	if err != nil {
		return "", err
	}

	err = serv.Update(user)
	if err != nil {
		return "", err
	}
	return token, nil
}
