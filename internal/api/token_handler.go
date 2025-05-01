package api

import (
	"encoding/json"
	"log"
	"net/http"

	"fem-go-crud/internal/auth"
	"fem-go-crud/internal/store"
	"fem-go-crud/internal/utils"
)

type TokenHandler struct {
	tokenStore store.TokenStore
	userStore  store.UserStore
	logger     *log.Logger
}

func NewTokenHandler(ts store.TokenStore, us store.UserStore, l *log.Logger) *TokenHandler {
	return &TokenHandler{
		tokenStore: ts,
		userStore:  us,
		logger:     l,
	}
}

type createTokenPayload struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (th *TokenHandler) CreateToken(w http.ResponseWriter, r *http.Request) {
	var payload createTokenPayload

	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		th.logger.Printf("ERROR: %v", err)
		_ = utils.WriteJSONResponse(w, http.StatusBadRequest, utils.Envelope{"error": "invalid payload"})
		return
	}

	user, err := th.userStore.GetUserByIdOrUsername(0, payload.Username)
	if err != nil {
		th.logger.Printf("ERROR: %v", err)
		_ = utils.WriteJSONResponse(w, http.StatusInternalServerError, utils.Envelope{"error": "failed to retrieve user"})
		return
	}
	if user == nil {
		th.logger.Printf("ERROR: user %s not found", payload.Username)
		_ = utils.WriteJSONResponse(w, http.StatusNotFound, utils.Envelope{"error": "user not found"})
		return
	}

	passwordMatches, err := user.Password.Matches(payload.Password)
	if err != nil {
		th.logger.Printf("ERROR: %v", err)
		_ = utils.WriteJSONResponse(w, http.StatusInternalServerError, utils.Envelope{"error": "failed to create token"})
		return
	}

	if !passwordMatches {
		th.logger.Printf("ERROR: invalid password for user %s", payload.Username)
		_ = utils.WriteJSONResponse(w, http.StatusUnauthorized, utils.Envelope{"error": "invalid password"})
		return
	}

	newToken, err := auth.MakeToken(user.ID, auth.TokenTTL, auth.TokenScopeAuth)
	if err != nil {
		th.logger.Printf("ERROR: %v", err)
		_ = utils.WriteJSONResponse(w, http.StatusInternalServerError, utils.Envelope{"error": "failed to create token"})
		return
	}

	err = th.tokenStore.PersistToken(newToken)
	if err != nil {
		th.logger.Printf("ERROR: %v", err)
		_ = utils.WriteJSONResponse(w, http.StatusInternalServerError, utils.Envelope{"error": "failed to create token"})
		return
	}

	_ = utils.WriteJSONResponse(w, http.StatusCreated, utils.Envelope{"token": newToken})
}
