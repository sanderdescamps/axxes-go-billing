package api

import (
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/context"
	"github.com/sanderdescamps/go-billing-api/gobilling/log"
	"golang.org/x/crypto/bcrypt"
)

const (
	PERM_BASIC_LOGIN       = "basic:login"
	PERM_RESOURCE_EDIT     = "resource:write"
	PERM_RESOURCE_VIEW     = "resource:read"
	PERM_RESOURCE_NEW      = "resource:new"
	PERM_RESOURCE_DELETE   = "resource:delete"
	PERM_COSTTYPE_EDIT     = "cost_type:write"
	PERM_COSTTYPE_VIEW     = "cost_type:read"
	PERM_COSTTYPE_NEW      = "cost_type:new"
	PERM_COSTTYPE_DELETE   = "cost_type.:delete"
	PERM_COSTCENTER_EDIT   = "cost_center:write"
	PERM_COSTCENTER_VIEW   = "cost_center:read"
	PERM_COSTCENTER_NEW    = "cost_center:new"
	PERM_COSTCENTER_DELETE = "cost_center.:delete"
	PERM_USER_GET_TOKEN    = "token:get"
	// PERM_USER_REVOKE_TOKEN = "token:revoke"
)

// type RegisterUserRequestBody struct {
// 	Username string `json:"username"`
// 	Password string `json:"password"`
// 	Role     string `json:"role"`
// }

// func RegisterHandler(w http.ResponseWriter, r *http.Request) {
// 	w.Header().Set("Content-Type", "application/json")
// 	isAdmin := context.Get(r, "isAdmin").(string)
// 	username := context.Get(r, "username").(string)

// 	var body RegisterUserRequestBody
// 	err := json.NewDecoder(r.Body).Decode(&body)
// 	if err != nil {
// 		http.Error(w, "Invalid request body", http.StatusBadRequest)
// 		return
// 	}

// 	hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), 5)
// 	if err != nil {
// 		http.Error(w, fmt.Sprintf("Failed to create user: %s", err), http.StatusBadRequest)
// 		return
// 	}
// 	user, err := treeDB.Users.Create(body.Username, string(hash), body.Role)
// 	if err != nil {
// 		http.Error(w, fmt.Sprintf("Failed to create user: %s", err), http.StatusBadRequest)
// 		return
// 	}
// 	resp := map[string]string{
// 		"username": user.Username,
// 		"id":       user.Id,
// 	}

// 	json.NewEncoder(w).Encode(resp)
// }

type LoginUserRequestBody struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// func getToken(w http.ResponseWriter, r *http.Request) {
// 	w.Header().Set("Content-Type", "application/json")
// 	var body LoginUserRequestBody
// 	err := json.NewDecoder(r.Body).Decode(&body)
// 	if err != nil {
// 		http.Error(w, "Invalid request body", http.StatusBadRequest)
// 		return
// 	}

// 	user, err := billingDB.Users.GetByUsername(body.Username)
// 	if err != nil {
// 		http.Error(w, "Invalid user", http.StatusBadRequest)
// 		return
// 	}

// 	expiresAt := time.Now().Add(time.Minute * 100000).Unix()

// 	errf := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(body.Password))
// 	if errf != nil && errf == bcrypt.ErrMismatchedHashAndPassword {
// 		var res = map[string]interface{}{"status": false, "message": "Invalid login credential. Please try again"}
// 		json.NewEncoder(w).Encode(res)
// 		return
// 	}

// 	token, err := NewToken(user, expiresAt)
// 	if err != nil {
// 		http.Error(w, fmt.Sprintf("Failed to create token: %s", err), http.StatusBadRequest)
// 		return
// 	}
// 	var resp = map[string]interface{}{"status": 200, "message": "logged in"}
// 	resp["token"] = token

// 	json.NewEncoder(w).Encode(resp)
// }

type Token struct {
	Username       string
	Id             string
	StandardClaims *jwt.StandardClaims
}

func (t *Token) Valid() error {
	return t.StandardClaims.Valid()
}

func NewToken(username string, userId string, expiresAt int64) (string, error) {
	tk := &Token{
		Username: username,
		Id:       userId,
		StandardClaims: &jwt.StandardClaims{
			ExpiresAt: expiresAt,
		},
	}

	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tk)

	tokenString, err := token.SignedString([]byte(jwtSecret))
	return tokenString, err

}

// func authMW(next http.HandlerFunc, reqPermisions []string) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		username, password, ok := r.BasicAuth()

// 		if !ok {
// 			errMsg := "Authentication error"
// 			http.Error(w, errMsg, http.StatusForbidden)
// 			log.Errorf("%s", errMsg)
// 			return
// 		}

// 		user, err := billingDB.Users.GetByUsername(username)
// 		if err != nil {
// 			http.Error(w, "Authentication error", http.StatusForbidden)
// 			return
// 		}
// 		errf := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
// 		if errf != nil && errf == bcrypt.ErrMismatchedHashAndPassword {
// 			http.Error(w, "Authentication error", http.StatusForbidden)
// 			return
// 		}

// 		if !billingDB.Users.HasPermissions(user, reqPermisions) {
// 			http.Error(w, "Authentication error", http.StatusMethodNotAllowed)
// 			return
// 		}

// 		// context.Set(r, "isAdmin", true)

// 		log.Debugf("User %s logged in.", username)
// 		next(w, r)
// 	}
// }

func authMW(next http.HandlerFunc, reqPermisions []string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		result := NewResult[string]()
		username, password, ok := r.BasicAuth()

		if !ok {
			result.Error("Authentication error: Invalid Basic Auth", http.StatusForbidden).Send(w)
			return
		}

		user, err := billingDB.Users.GetByUsername(username)
		if err != nil {
			result.Error("Authentication error: Invalid user", http.StatusForbidden).Send(w)
			return
		}

		errf := user.MatchPassword(password)
		if errf != nil && errf == bcrypt.ErrMismatchedHashAndPassword {
			result.Error("Authentication error: Incorrect password", http.StatusForbidden).Send(w)
			return
		}

		if allowed, err := billingDB.Users.HasPermissions(user.Id, reqPermisions); !allowed {
			result.Error("Authentication error: Operation not allowed", http.StatusMethodNotAllowed).Send(w)
			return
		} else if err != nil {

		}

		context.Set(r, "user_id", user.Id)
		context.Set(r, "username", user.Username)

		log.Debugf("User %s logged in.", username)
		next(w, r)
	}
}
