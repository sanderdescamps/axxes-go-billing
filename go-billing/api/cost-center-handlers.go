package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/sanderdescamps/go-billing-api/gobilling/log"
	"github.com/sanderdescamps/go-billing-api/gobilling/model"
)

// CostCenter

func getAllCostCenters(w http.ResponseWriter, r *http.Request) {
	result := NewResult[model.CostCenter]()

	if r.URL.Query().Has("name") {
		name := r.URL.Query().Get("name")
		t, err := billingDB.CostCenters.GetByName(name)
		if err != nil {
			result.Error(err.Error(), http.StatusBadRequest).Send(w)
			log.Debugf("%v", err)
			return
		}
		log.Debugf("CostCenter found ID=%s", t.Id)
		result.AddResult(*t)
	} else {
		var err error
		costCenters, err := billingDB.CostCenters.GetAll()
		if err != nil {
			result.Error("failed to get all costTypes", http.StatusBadGateway).Send(w)
			log.Errorf(err.Error())
			return
		}
		result.AddResult(costCenters...)
	}
	result.Send(w)
}

func getCostCenter(w http.ResponseWriter, r *http.Request) {
	result := NewResult[model.CostCenter]()

	params := mux.Vars(r)
	ccId := params["id"]

	if !billingDB.CostCenters.Exits(ccId) {
		result.Error("CostCenter not found", http.StatusNotFound).Send(w)
		return
	}

	costCenter, err := billingDB.CostCenters.Get(ccId)
	if err != nil {
		result.Error("failed to get costCenter by ID", http.StatusBadGateway)
		return
	}
	result.AddResult(*costCenter).Send(w)
}

func getTotalCostCenterCost(w http.ResponseWriter, r *http.Request) {
	result := NewResult[model.Cost]()
	params := mux.Vars(r)
	ccId := params["id"]

	total, err := billingDB.CostCenters.GetTotalCost(ccId)
	if err != nil {
		result.Error("failed to get total cost", http.StatusBadGateway)
		log.Errorf(err.Error())
		return
	}
	result.AddResult(*total).Send(w)
}

type CreateCostCenterRequestBody struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

func createCostCenter(w http.ResponseWriter, r *http.Request) {
	result := NewResult[model.CostCenter]()

	var body CreateCostCenterRequestBody
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		result.Error("Invalid request body", http.StatusBadRequest).Send(w)
		return
	}

	if billingDB.CostCenters.ExitsByName(body.Name) {
		result.Error("CostCenter already exists", http.StatusBadRequest).Send(w)
		return
	}

	newCostCenter, err := billingDB.CostCenters.Create(body.Name, body.Description)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to create CostCenter: %s", err), http.StatusBadRequest)
		return
	}
	log.Debugf("New CostCenter created ID=%s", newCostCenter.Id)
	result.Changed().Status(http.StatusCreated).AddResult(*newCostCenter)
	result.Send(w)
}

func updateCostCenter(w http.ResponseWriter, r *http.Request) {
	result := NewResult[model.CostCenter]()

	params := mux.Vars(r)
	ccId := params["id"]

	toChangeCostCenter, err := billingDB.CostCenters.Get(ccId)
	if err != nil {
		result.Error("CostCenter not found", http.StatusNotFound).Send(w)
		return
	}
	origHash := toChangeCostCenter.ToHash()

	var body CreateCostCenterRequestBody
	err = json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		result.Error("Invalid request body", http.StatusBadRequest).Send(w)
		return
	}

	if body.Description != "" {
		toChangeCostCenter.Description = body.Description
	}
	if body.Name != "" {
		log.Warningf("TODO: check for name conflicts")
		toChangeCostCenter.Name = body.Name
	}

	if origHash != toChangeCostCenter.ToHash() {
		err := billingDB.CostCenters.Update(toChangeCostCenter)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to update costCenter: %v", err), http.StatusBadGateway)
			return
		}
		result.Changed()
	}
	result.AddResult(*toChangeCostCenter).Send(w)
}

func deleteCostCenter(w http.ResponseWriter, r *http.Request) {
	result := NewResult[model.CostCenter]()

	params := mux.Vars(r)
	ccId := params["id"]

	if !billingDB.CostCenters.Exits(ccId) {
		result.Error("CostCenter not found", http.StatusNotFound).Send(w)
		return
	}

	err := billingDB.CostCenters.DeleteById(ccId)
	if err != nil {
		result.Error("Failed to delete CostCenter", http.StatusBadGateway).Send(w)
		log.Errorf("Failed to delete CostCenter: %s", err.Error())
		return
	}
	result.Msg("Delete costType ID=%s", ccId)
	log.Debugf("Delete costType ID=%s", ccId)
	result.Changed().Send(w)
}
