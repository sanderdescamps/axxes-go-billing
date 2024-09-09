package model

import (
	"crypto/md5"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
)

const (
	PARAM_COST_TYPE_ID = "type_id"
)

type CostType struct {
	TypeID      string `json:"type_id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Cost        Cost   `json:"cost"`
}

func (t CostType) ID() (jsonField string, value interface{}) {
	value = t.TypeID
	jsonField = PARAM_COST_TYPE_ID
	return
}

func (s *CostType) ToString() string {
	return fmt.Sprintf("{ID=%s, Name=%s}", s.TypeID, s.Name)
}

func (t *CostType) ToHash() string {
	jsonBytes, _ := json.Marshal(t)
	return fmt.Sprintf("%x", md5.Sum(jsonBytes))
}

func NewCostType(name string, description string, cost Cost) *CostType {
	id := uuid.New().String()
	s := CostType{TypeID: id, Name: name, Description: description, Cost: cost}
	return &s
}
