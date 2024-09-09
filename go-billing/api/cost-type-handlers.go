package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/sanderdescamps/go-billing-api/gobilling/log"
	"github.com/sanderdescamps/go-billing-api/gobilling/model"
)

// CostType

func getAllCostTypes(w http.ResponseWriter, r *http.Request) {
	result := NewResult[model.CostType]()

	if r.URL.Query().Has("name") {
		name := r.URL.Query().Get("name")
		t, err := billingDB.CostTypes.GetByName(name)
		if err != nil {
			result.Error(err.Error(), http.StatusBadRequest).Send(w)
			log.Debugf("%v", err)
			return
		}
		log.Debugf("CostType found ID=%s", t.TypeID)
		result.AddResult(*t)
	} else {
		var err error
		costTypes, err := billingDB.CostTypes.GetAll()
		if err != nil {
			result.Error("failed to get all costTypes", http.StatusBadGateway).Send(w)
			log.Errorf(err.Error())
			return
		}
		result.AddResult(costTypes...)
	}
	result.Send(w)
}

func getCostType(w http.ResponseWriter, r *http.Request) {
	result := NewResult[model.CostType]()

	params := mux.Vars(r)
	ctId := params["id"]

	if !billingDB.CostTypes.Exits(ctId) {
		result.Error("CostType not found", http.StatusNotFound).Send(w)
		return
	}

	costType, err := billingDB.CostTypes.Get(ctId)
	if err != nil {
		result.Error("failed to get costType by ID", http.StatusBadGateway)
		return
	}
	result.AddResult(*costType)
	result.Send(w)
}

type CreateCostTypeRequestBody struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	CostPerSec  *float64 `json:"cost_per_sec"`
	FixedCost   *float64 `json:"cost_fixed"`
}

func createCostType(w http.ResponseWriter, r *http.Request) {
	result := NewResult[model.CostType]()

	var body CreateCostTypeRequestBody
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		result.Error("Invalid request body", http.StatusBadRequest).Send(w)
		return
	}

	if billingDB.CostTypes.ExitsByName(body.Name) {
		result.Error("CostType already exists", http.StatusBadRequest).Send(w)
		return
	}

	newCostType := model.NewCostType(body.Name, body.Description, model.Cost{Fixed: *body.FixedCost, CostPerSec: *body.CostPerSec})
	err = billingDB.CostTypes.Create(newCostType)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to create CostType: %s", err), http.StatusBadRequest)
		return
	}
	log.Debugf("New CostType created ID=%s", newCostType.TypeID)
	result.Changed().Status(http.StatusCreated).AddResult(*newCostType)
	result.Send(w)
}

func updateCostType(w http.ResponseWriter, r *http.Request) {
	result := NewResult[model.CostType]()

	params := mux.Vars(r)
	ctId := params["id"]

	toChangeCostType, err := billingDB.CostTypes.Get(ctId)
	if err != nil {
		result.Error("CostType not found", http.StatusNotFound).Send(w)
		return
	}
	origHash := toChangeCostType.ToHash()

	var body CreateCostTypeRequestBody
	err = json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		result.Error("Invalid request body", http.StatusBadRequest).Send(w)
		return
	}

	if body.Description != "" {
		toChangeCostType.Description = body.Description
	}
	if body.Name != "" {
		log.Warningf("TODO: check for name conflicts")
		toChangeCostType.Name = body.Name
	}
	if body.FixedCost != nil {
		toChangeCostType.Cost.Fixed = *body.FixedCost
	}
	if body.CostPerSec != nil {
		toChangeCostType.Cost.CostPerSec = *body.CostPerSec
	}

	if origHash != toChangeCostType.ToHash() {
		err := billingDB.CostTypes.Update(toChangeCostType)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to update costType: %v", err), http.StatusBadGateway)
			return
		}
		result.Changed()
	}
	result.AddResult(*toChangeCostType).Send(w)
}

func deleteCostType(w http.ResponseWriter, r *http.Request) {
	result := NewResult[model.CostType]()

	params := mux.Vars(r)
	ctId := params["id"]

	if !billingDB.CostTypes.Exits(ctId) {
		result.Error("CostType not found", http.StatusNotFound).Send(w)
		return
	}

	err := billingDB.CostTypes.DeleteById(ctId)
	if err != nil {
		result.Error("Failed to delete CostType", http.StatusBadGateway).Send(w)
		log.Errorf("Failed to delete CostType: %s", err.Error())
		return
	}
	result.Msg("Delete costType ID=%s", ctId)
	log.Debugf("Delete costType ID=%s", ctId)
	result.Changed().Send(w)
}
