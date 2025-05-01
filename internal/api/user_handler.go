package api

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"regexp"

	"fem-go-crud/internal/store"
	"fem-go-crud/internal/utils"
)

type UserHandler struct {
	userStore store.UserStore
	logger    *log.Logger
}

func NewUserHandler(us store.UserStore, l *log.Logger) *UserHandler {
	return &UserHandler{
		userStore: us,
		logger:    l,
	}
}

type registerUserPayload struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (uh *UserHandler) validateRegisterUserPayload(payload *registerUserPayload) error {
	if payload.Username == "" {
		return errors.New("missing username")
	}
	if len(payload.Username) < 3 || len(payload.Username) > 50 {
		return errors.New("invalid username length")
	}

	if payload.Email == "" {
		return errors.New("missing email")
	}
	if len(payload.Email) < 5 || len(payload.Email) > 100 {
		return errors.New("invalid email length")
	}
	emailRegex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	emailRegexCheck := regexp.MustCompile(emailRegex)
	if !emailRegexCheck.MatchString(payload.Email) {
		return errors.New("invalid email")
	}

	if payload.Password == "" {
		return errors.New("missing password")
	}
	if len(payload.Password) < 8 || len(payload.Password) > 100 {
		return errors.New("invalid password length")
	}
	// simple password validation, for demo purposes only
	passwordRegex := `^(.*[0-9])`
	passwordRegexCheck := regexp.MustCompile(passwordRegex)
	if !passwordRegexCheck.MatchString(payload.Password) {
		return errors.New("invalid password")
	}

	return nil
}

func (uh *UserHandler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	var payload registerUserPayload

	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		uh.logger.Printf("ERROR: %v", err)
		_ = utils.WriteJSONResponse(w, http.StatusBadRequest, utils.Envelope{"error": "invalid payload"})
		return
	}

	err = uh.validateRegisterUserPayload(&payload)
	if err != nil {
		uh.logger.Printf("ERROR: %v", err)
		_ = utils.WriteJSONResponse(w, http.StatusBadRequest, utils.Envelope{"error": err.Error()})
	}

	user := store.User{
		Username: payload.Username,
		Email:    payload.Email,
	}

	err = user.Password.Set(payload.Password)
	if err != nil {
		uh.logger.Printf("ERROR: %v", err)
		_ = utils.WriteJSONResponse(w, http.StatusInternalServerError, utils.Envelope{"error": "failed to persist resource"})
		return
	}

	err = uh.userStore.PersistUser(&user)
	if err != nil {
		uh.logger.Printf("ERROR: %v", err)
		_ = utils.WriteJSONResponse(w, http.StatusInternalServerError, utils.Envelope{"error": "failed to persist resource"})
		return
	}

	_ = utils.WriteJSONResponse(w, http.StatusCreated, utils.Envelope{"user": user})
}

func (uh *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	userID, err := utils.ParseIDParamFromURL(r, "userId")
	if err != nil {
		uh.logger.Printf("ERROR: %v", err)
		_ = utils.WriteJSONResponse(w, http.StatusBadRequest, utils.Envelope{"error": "invalid ID param"})
	}

	user, err := uh.userStore.GetUserByIdOrUsername(userID, "")
	if err != nil {
		uh.logger.Printf("ERROR: %v", err)
		_ = utils.WriteJSONResponse(w, http.StatusInternalServerError, utils.Envelope{"error": "failed to retrieve resource"})
		return
	}

	if user == nil {
		uh.logger.Printf("ERROR: user %d not found", userID)
		_ = utils.WriteJSONResponse(w, http.StatusNotFound, utils.Envelope{"error": "resource not found"})
		return
	}

	_ = utils.WriteJSONResponse(w, http.StatusOK, utils.Envelope{"user": user})
}
