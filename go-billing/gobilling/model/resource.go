package model

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
)

const (
	PARAM_RESOURCE_ID = "resource_id"
)

type Resource struct {
	ResourceID   string   `json:"resource_id"`
	Name         string   `json:"name"`
	Description  string   `json:"description"`
	CreationTime string   `json:"creation_time"`
	CostTypesIds []string `json:"cost_type_ids"`
	CostCenterId string   `json:"cost_center_id"`
	Value        float64  `json:"value"`
}

func (r Resource) ID() (jsonField string, value interface{}) {
	value = r.ResourceID
	jsonField = PARAM_RESOURCE_ID
	return
}

func (r Resource) HasCostType(costTypeId string) bool {
	for _, i := range r.CostTypesIds {
		if i == costTypeId {
			return true
		}
	}
	return false
}

func (t Resource) ToHash() string {
	jsonBytes, _ := json.Marshal(t)
	return fmt.Sprintf("%x", md5.Sum(jsonBytes))
}

func (r *Resource) ToString() string {
	return fmt.Sprintf("{ID=%s, Name=%s}", r.ResourceID, r.Name)
}

func NewResource(name string, description string, costCenterId string, value float64) *Resource {
	id := uuid.New().String()
	now := time.Now()
	creationTime := now.Format(time.RFC822)
	r := Resource{
		ResourceID:   id,
		Name:         name,
		Description:  description,
		CostCenterId: costCenterId,
		CreationTime: creationTime,
		Value:        value,
	}
	return &r
}
