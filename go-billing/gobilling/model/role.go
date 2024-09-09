package model

import (
	"fmt"
	"strings"

	"github.com/google/uuid"
)

const (
	PARAM_ROLE_ID = "role_id"
)

type Role struct {
	RoleId      string   `json:"role_id"`
	Name        string   `json:"name"`
	Permissions []string `json:"permissions"`
}

func NewRole(name string, permissions []string) *Role {
	id := uuid.New().String()
	return &Role{
		RoleId:      id,
		Name:        name,
		Permissions: permissions,
	}
}

func (r Role) ID() (jsonField string, value interface{}) {
	value = r.Name
	jsonField = PARAM_ROLE_ID
	return
}

func (r *Role) Add(p string) {
	r.Permissions = append(r.Permissions, p)
}

func (r *Role) Remove(p string) {
	newPermissions := []string{}
	for _, perm := range r.Permissions {
		if perm != p {
			newPermissions = append(newPermissions, p)
		}
	}
	r.Permissions = newPermissions
}

func (r *Role) ToString() string {
	return fmt.Sprintf("{ID=%s, Name=%s}", r.RoleId, r.Name)
}

func (r *Role) IsAdmin() bool {
	return strings.ToLower(r.Name) == "admin"
}

func (r *Role) IsAllowed(p string) bool {
	if r.IsAdmin() {
		return true
	}
	for _, perm := range r.Permissions {
		if perm == p {
			return true
		}
	}
	return false
}
