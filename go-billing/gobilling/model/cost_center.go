package model

import (
	"crypto/md5"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
)

const (
	PARAM_COSTCENTER_ID = "cost_center_id"
)

type CostCenter struct {
	Id          string `json:"cost_center_id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

func (r CostCenter) ID() (jsonField string, value interface{}) {
	value = r.Id
	jsonField = PARAM_COSTCENTER_ID
	return
}

func (t CostCenter) ToHash() string {
	jsonBytes, _ := json.Marshal(t)
	return fmt.Sprintf("%x", md5.Sum(jsonBytes))
}

func (r *CostCenter) ToString() string {
	return fmt.Sprintf("{ID=%s, Name=%s}", r.Id, r.Name)
}

func NewCostCenter(name string, description string) *CostCenter {
	id := uuid.New().String()
	c := CostCenter{Id: id, Name: name, Description: description}
	return &c
}
