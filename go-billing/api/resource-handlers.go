package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sanderdescamps/go-billing-api/gobilling/log"
	"github.com/sanderdescamps/go-billing-api/gobilling/model"
)

func getAllResources(w http.ResponseWriter, r *http.Request) {
	result := NewResult[model.Resource]()
	if r.URL.Query().Has("name") {
		name := r.URL.Query().Get("name")
		log.Debugf("Search for resources with name: %s", name)
		r, err := billingDB.Resources.GetByName(name)
		if err != nil {
			result.Error(err.Error(), http.StatusBadRequest).Send(w)
			log.Debugf("%v", err)
			return
		}
		log.Debugf("Resource found ID=%s", r.ResourceID)
		result.AddResult(*r)
	} else {
		var err error
		resources, err := billingDB.Resources.GetAll()
		if err != nil {
			result.Error("failed to get all resources", http.StatusBadGateway).Send(w)
			log.Errorf(err.Error())
			return
		}
		result.AddResult(resources...)
	}
	result.Send(w)
}

func getResource(w http.ResponseWriter, r *http.Request) {
	result := NewResult[model.Resource]()

	params := mux.Vars(r)
	rId := params["id"]

	if !billingDB.Resources.Exits(rId) {
		result.Error("Resource not found", http.StatusNotFound).Send(w)
		return
	}

	resource, err := billingDB.Resources.Get(rId)
	if err != nil {
		result.Error("failed to get resource by ID", http.StatusBadGateway).Send(w)
		return
	}
	result.AddResult(*resource).Send(w)
}

type CreateResourceRequestBody struct {
	Name          string   `json:"name"`
	Description   string   `json:"description"`
	CostTypesIds  []string `json:"cost_type_ids"`
	CostTypeNames []string `json:"cost_type_names"`
	CostCenter    string   `json:"cost_center"`
	CostCenterId  string   `json:"cost_center_id"`
	Value         float64  `json:"value"`
}

func createResource(w http.ResponseWriter, r *http.Request) {
	result := NewResult[model.Resource]()

	var body CreateResourceRequestBody
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		result.Error("Invalid request body", http.StatusBadRequest).Send(w)
		return
	}

	if billingDB.Resources.ExitsByName(body.Name) {
		result.Error("Resource already exists", http.StatusBadRequest).Send(w)
		return
	}

	if body.CostCenterId != "" {
		if !billingDB.CostCenters.Exits(body.CostCenterId) {
			result.Error("costCenter does not exist", http.StatusBadRequest).Send(w)
			return
		}
	} else if body.CostCenter != "" {
		if costCenter, err := billingDB.CostCenters.GetByName(body.CostCenter); err != nil {
			result.Error("costCenter does not exist", http.StatusBadRequest).Send(w)
			return
		} else {
			body.CostCenterId = costCenter.Id
		}
	}

	for _, i := range body.CostTypeNames {
		ct, err := billingDB.CostTypes.GetByName(i)
		if err != nil {
			result.Error(fmt.Sprintf("invalid costType name: %s", i), http.StatusBadRequest).Send(w)
			return
		}
		body.CostTypesIds = append(body.CostTypesIds, ct.TypeID)
	}

	newResource, err := billingDB.Resources.Create(body.Name, body.Description, body.CostCenterId, body.CostTypesIds, body.Value)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to create resource: %s", err), http.StatusBadRequest)
		return
	}
	log.Debugf("New Resource created ID=%s, %v", newResource.ResourceID, newResource)
	result.Changed().Status(http.StatusCreated).AddResult(*newResource)
	result.Send(w)
}

type UpdateResourceRequestBody struct {
	Name         string   `json:"name"`
	Description  string   `json:"description"`
	Value        float64  `json:"value"`
	CostTypesIds []string `json:"cost_type_ids"`
	CostCenter   string   `json:"cost_center"`
	CostCenterId string   `json:"cost_center_id"`
}

func updateResource(w http.ResponseWriter, r *http.Request) {
	result := NewResult[model.Resource]()

	params := mux.Vars(r)
	rId := params["id"]

	toChangeResource, err := billingDB.Resources.Get(rId)
	if err != nil {
		result.Error("Resource not found", http.StatusBadRequest).Send(w)
		return
	}
	hash1 := toChangeResource.ToHash()

	var body UpdateResourceRequestBody
	err = json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		result.Error("Invalid request body", http.StatusBadRequest).Send(w)
		return
	}

	if body.Description != "" {
		toChangeResource.Description = body.Description
	}
	if body.Name != "" {
		log.Warningf("TODO: check for name conflicts")
		toChangeResource.Name = body.Name
	}
	if body.CostTypesIds != nil {
		setCostIds := []string{}
		for _, i := range body.CostTypesIds {
			if billingDB.CostTypes.Exits(i) {
				setCostIds = append(setCostIds)
			} else {
				result.Error(fmt.Sprintf("Invalid costType id %s", i), http.StatusBadRequest).Send(w)
				return
			}
		}
		toChangeResource.CostTypesIds = setCostIds
	}
	if body.CostCenterId != "" {
		if !billingDB.CostCenters.Exits(body.CostCenterId) {
			result.Error("costCenter does not exist", http.StatusBadRequest).Send(w)
			return
		}
		toChangeResource.CostCenterId = body.CostCenterId
	} else if body.CostCenter != "" {
		if costCenter, err := billingDB.CostCenters.GetByName(body.CostCenter); err != nil {
			result.Error("costCenter does not exist", http.StatusBadRequest).Send(w)
			return
		} else {
			toChangeResource.CostCenterId = costCenter.Id
		}
	}

	if hash1 != toChangeResource.ToHash() {
		err := billingDB.Resources.Update(toChangeResource)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to update resource: %v", err), http.StatusBadGateway)
			return
		}
		result.Changed()
	}

	result.AddResult(*toChangeResource).Send(w)
}

func updateResourceAddCostType(w http.ResponseWriter, r *http.Request) {
	result := NewResult[model.Resource]()

	params := mux.Vars(r)
	rId := params["id"]
	costTypeID := params["costTypeID"]

	toChangeResource, err := billingDB.Resources.Get(rId)
	if err != nil {
		result.Error("Resource not found", http.StatusBadRequest).Send(w)
		return
	}
	hash1 := toChangeResource.ToHash()

	skip := false
	for _, i := range toChangeResource.CostTypesIds {
		if i == costTypeID {
			skip = true
			break
		}
	}
	if !skip {
		toChangeResource.CostTypesIds = append(toChangeResource.CostTypesIds, costTypeID)
	}

	if hash1 != toChangeResource.ToHash() {
		err := billingDB.Resources.Update(toChangeResource)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to update resource: %v", err), http.StatusBadGateway)
			return
		}
		result.Changed().AddResult(*toChangeResource)
	}

	result.Send(w)
}

func updateResourceDeleteCostType(w http.ResponseWriter, r *http.Request) {
	result := NewResult[model.Resource]()

	params := mux.Vars(r)
	rId := params["id"]
	costTypeID := params["costTypeID"]

	toChangeResource, err := billingDB.Resources.Get(rId)
	if err != nil {
		result.Error("Resource not found", http.StatusBadRequest).Send(w)
		return
	}
	hash1 := toChangeResource.ToHash()

	keepCostTypes := []string{}
	for _, i := range toChangeResource.CostTypesIds {
		if i != costTypeID {
			keepCostTypes = append(keepCostTypes, i)
		}
	}
	toChangeResource.CostTypesIds = keepCostTypes

	if hash1 != toChangeResource.ToHash() {
		err := billingDB.Resources.Update(toChangeResource)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to update resource: %v", err), http.StatusBadGateway)
			return
		}
		result.Changed().AddResult(*toChangeResource)
	}

	result.Send(w)
}

func deleteResource(w http.ResponseWriter, r *http.Request) {
	result := NewResult[model.Resource]()

	params := mux.Vars(r)
	rId := params["id"]

	if !billingDB.Resources.Exits(rId) {
		result.Error("Resource not found", http.StatusBadRequest).Send(w)
		return
	}
	err := billingDB.Resources.DeleteById(rId)
	if err != nil {
		result.Error("Failed to delete resource", http.StatusBadGateway).Send(w)
		log.Errorf("Failed to delete resource: %s", err.Error())
		return
	}
	result.Msg("Delete resource ID=%s", rId)
	log.Debugf("Delete resource ID=%s", rId)
	result.Changed().Send(w)
}
