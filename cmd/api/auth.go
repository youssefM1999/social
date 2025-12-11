package main

import (
	"net/http"

	"github.com/youssefM1999/social/internal/store"
)

type RegisterUserPayload struct {
	Username string `json:"username" validate:"required,max=100"`
	Email    string `json:"email" validate:"required,email,max=255"`
	Password string `json:"password" validate:"required,min=3,max=72"`
}

// RegisterUser godoc
//
//	@Summary		Register user
//	@Description	Register a new user account
//	@Tags			authentication
//	@Accept			json
//	@Produce		json
//	@Param			payload	body	RegisterUserPayload	true	"User registration payload"
//	@Success		201
//	@Failure		400	{object}	error
//	@Failure		500	{object}	error
//	@Router			/authentication/user [post]
func (app *application) registerUserHandler(w http.ResponseWriter, r *http.Request) {
	var payload RegisterUserPayload
	if err := readJSON(w, r, payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if err := Validate.Struct(payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	user := &store.User{
		Username: payload.Username,
		Email:    payload.Email,
	}
	// Hash the user password
	if err := user.Password.Set(payload.Password); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	// Store the user
	if err := app.store.Users.CreateAndInvite(r.Context(), user, "uuidv4"); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	// Send email

	if err := app.jsonResponse(w, http.StatusCreated, nil); err != nil {
		app.internalServerError(w, r, err)
	}
}
