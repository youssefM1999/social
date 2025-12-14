package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/youssefM1999/social/internal/mailer"
	"github.com/youssefM1999/social/internal/store"
)

type UserWithToken struct {
	User  *store.User `json:"user"`
	Token string      `json:"token"`
}

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
	if err := readJSON(w, r, &payload); err != nil {
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

	ctx := r.Context()

	plainToken := uuid.New().String()

	// Store the hashed token in the database
	hash := sha256.Sum256([]byte(plainToken))
	hashToken := hex.EncodeToString(hash[:])

	// Store the user
	if err := app.store.Users.CreateAndInvite(ctx, user, hashToken, app.config.mail.exp); err != nil {
		switch err {
		case store.ErrDuplicateEmail:
			app.badRequestResponse(w, r, err)
		case store.ErrDuplicateUsername:
			app.badRequestResponse(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}

	userWithToken := UserWithToken{
		User:  user,
		Token: plainToken,
	}

	activationURL := fmt.Sprintf("%s/confirm/%s", app.config.frontendURL, plainToken)

	isProdEnv := app.config.env == "production"
	app.logger.Infow("sending invitation email",
		"env", app.config.env,
		"isProduction", isProdEnv,
		"isSandbox", !isProdEnv,
		"to", user.Email,
		"from", app.config.mail.fromEmail,
	)

	vars := struct {
		Username      string
		ActivationURL string
	}{
		Username:      user.Username,
		ActivationURL: activationURL,
	}

	// Send email
	statusCode, err := app.mailer.Send(
		mailer.UserInvitationTemplateFile,
		user.Username,
		user.Email,
		vars,
		!isProdEnv,
	)
	if err != nil {
		app.logger.Errorw("failed to send welcome email", "error", err, "statusCode", statusCode)

		// rollback user creation if email fails (SAGA pattern)
		if err := app.store.Users.Delete(ctx, user.ID); err != nil {
			app.logger.Errorw("failed to rollback user creation", "error", err)

		}

		app.internalServerError(w, r, err)
		return
	}
	app.logger.Infow("email sent successfully", "statusCode", statusCode, "to", user.Email)

	if err := app.jsonResponse(w, http.StatusCreated, userWithToken); err != nil {
		app.internalServerError(w, r, err)
	}
}
