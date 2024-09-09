package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/context"
	"github.com/sanderdescamps/go-billing-api/gobilling/log"
	"github.com/sanderdescamps/go-billing-api/gobilling/model"
)

type RegisterUserRequestBody struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Role     string `json:"role"`
	RoleId   string `json:"role_id"`
}

func NewUser(w http.ResponseWriter, r *http.Request) {
	result := NewResult[map[string]string]()

	var body RegisterUserRequestBody
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		result.Error("Invalid request body", http.StatusBadRequest).Send(w)
		return
	}

	var userRole *model.Role
	if body.RoleId == "" && body.Role != "" {
		userRole, err = billingDB.Roles.GetByName(body.Role)
		if err != nil {
			result.Error("Role does not exist", http.StatusBadRequest).Send(w)
			return
		}
		body.RoleId = userRole.RoleId
	} else if body.RoleId != "" {
		if !billingDB.Roles.Exits(body.RoleId) {
			result.Error("Role does not exist", http.StatusBadRequest).Send(w)
			return
		}
	}

	user, err := billingDB.Users.Create(body.Username, body.Password, body.RoleId)
	if err != nil {
		result.Error(fmt.Sprintf("Failed to create user (%s)", body.Username), http.StatusBadRequest).Send(w)
		log.Errorf(fmt.Sprintf("Failed to create user (%s): %v", body.Username, err))
		return
	}
	userShort := map[string]string{
		"username": user.Username,
		"id":       user.Id,
	}
	result.AddResult(userShort).Changed().Msg("user %s created", user.Username)
	result.Send(w)
}

func getBearerToken(w http.ResponseWriter, r *http.Request) {
	username := context.Get(r, "username").(string)
	userId := context.Get(r, "user_id").(string)

	expiresAt := time.Now().Add(time.Minute * 60)
	token, err := NewToken(username, userId, expiresAt.Unix())
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to create token: %s", err), http.StatusBadRequest)
		return
	}
	var resp = map[string]interface{}{
		"status": 200,
		"msg":    fmt.Sprintf("Token will be valid till %s", expiresAt.Format(time.RFC850)),
		"token":  token,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(resp)
}

func getApiToken(w http.ResponseWriter, r *http.Request) {
	userId := context.Get(r, "user_id").(string)

	token, err := billingDB.Users.SetNewApiToken(userId)
	if err != nil {
		http.Error(w, "Failed to create token", http.StatusBadRequest)
		log.Errorf("Failed to create token: %s", err.Error())
		return
	}
	var resp = map[string]interface{}{
		"status": 200,
		"msg":    "",
		"token":  token,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(resp)
}
