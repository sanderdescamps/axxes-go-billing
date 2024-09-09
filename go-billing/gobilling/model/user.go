package model

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

const (
	PARAM_USER_ID = "id"
)

type User struct {
	Id           string `json:"id"`
	Username     string `json:"username"`
	PasswordHash string `json:"password_hash"`
	ApiTokenHash string `json:"api_token_hash"`
	CreationTime string `json:"creation_time"`
	RoleId       string `json:"role_id"`
}

func (u User) ID() (jsonField string, value interface{}) {
	value = u.Id
	jsonField = PARAM_USER_ID
	return
}

func (u *User) ToString() string {
	return fmt.Sprintf("{ID=%s, Name=%s, RoleId=%s}", u.Id, u.Username, u.RoleId)
}

func NewUser(username string, pasword string, roleId string) *User {
	id := uuid.New().String()
	now := time.Now()
	creationTime := now.Format(time.RFC822)
	u := User{Username: username, Id: id, CreationTime: creationTime, RoleId: roleId}
	u.SetPassword(pasword)
	return &u
}

func (u *User) SetPassword(password string) {
	hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	u.PasswordHash = string(hash)
}

func (u *User) MatchPassword(password string) error {
	return bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password))
}

func (u *User) SetNewApiToken() (string, error) {
	tokenId := uuid.New().String()
	tokenBody := struct {
		Username string `json:"username"`
		TokenId  string `json:"token_id"`
		UserId   string `json:"user_id"`
	}{
		TokenId:  tokenId,
		Username: u.Username,
		UserId:   u.Id,
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(tokenId), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	u.ApiTokenHash = string(hash)

	jsonBytes, err := json.Marshal(tokenBody)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(jsonBytes), nil
}

func (u *User) MatchApiToken(token string) error {
	jsonBytes, _ := base64.StdEncoding.DecodeString(token)

	tokenBody := struct {
		Username string `json:"username"`
		TokenId  string `json:"token_id"`
		UserId   string `json:"user_id"`
	}{}

	err := json.Unmarshal(jsonBytes, &tokenBody)
	if err != nil {
		return fmt.Errorf("invalid token")
	}

	return bcrypt.CompareHashAndPassword([]byte(u.ApiTokenHash), []byte(tokenBody.TokenId))
}
