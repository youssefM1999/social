package main

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/youssefM1999/social/internal/store"
)

type userKey string

const userCtx userKey = "user"

// GetUser godoc
//
//	@Summary		Get user
//	@Description	get user
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			userID	path		int64	true	"User ID"
//	@Success		200		{object}	store.User
//	@Failure		400		{object}	error
//	@Failure		404		{object}	error
//	@Failure		500		{object}	error
//	@Router			/users/{userID} [get]
func (app *application) getUserHandler(w http.ResponseWriter, r *http.Request) {
	user := getUserFromCtx(r)

	if err := app.jsonResponse(w, http.StatusOK, user); err != nil {
		app.internalServerError(w, r, err)
	}
}

type FollowUserPayload struct {
	UserID int64 `json:"user_id"`
}

// FollowUser godoc
//
//	@Summary		Follow user
//	@Description	Follow a user
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			userID	path	int64				true	"User ID to follow"
//	@Param			payload	body	FollowUserPayload	true	"Follow payload"
//	@Success		204
//	@Failure		400	{object}	error
//	@Failure		404	{object}	error
//	@Failure		409	{object}	error
//	@Failure		500	{object}	error
//	@Router			/users/{userID}/follow [put]
func (app *application) followUserHandler(w http.ResponseWriter, r *http.Request) {
	followerUser := getUserFromCtx(r)

	// TODO: Revert back to auth userID from ctx
	var payload FollowUserPayload
	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	ctx := r.Context()
	if err := app.store.Followers.Follow(ctx, followerUser.ID, payload.UserID); err != nil {
		switch err {
		case store.ErrConflict:
			app.conflictResponse(w, r, err)
			return
		default:
			app.internalServerError(w, r, err)
			return
		}
	}

	if err := app.jsonResponse(w, http.StatusNoContent, nil); err != nil {
		app.internalServerError(w, r, err)
	}
}

// UnfollowUser godoc
//
//	@Summary		Unfollow user
//	@Description	Unfollow a user
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			userID	path	int64				true	"User ID to unfollow"
//	@Param			payload	body	FollowUserPayload	true	"Unfollow payload"
//	@Success		204
//	@Failure		400	{object}	error
//	@Failure		404	{object}	error
//	@Failure		500	{object}	error
//	@Router			/users/{userID}/unfollow [put]
func (app *application) unfollowUserHandler(w http.ResponseWriter, r *http.Request) {
	unfollowedUser := getUserFromCtx(r)

	if err := app.jsonResponse(w, http.StatusNoContent, nil); err != nil {
		app.internalServerError(w, r, err)
	}

	// TODO: Revert back to auth userID from ctx
	var payload FollowUserPayload
	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	ctx := r.Context()
	app.store.Followers.Unfollow(ctx, unfollowedUser.ID, payload.UserID)

	if err := app.jsonResponse(w, http.StatusNoContent, nil); err != nil {
		app.internalServerError(w, r, err)
	}

}

func (app *application) userContextMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		idParam := chi.URLParam(r, "userID")
		id, err := strconv.ParseInt(idParam, 10, 64)
		if err != nil {
			app.internalServerError(w, r, err)
		}

		ctx := r.Context()
		user, err := app.store.Users.GetByID(ctx, id)
		if err != nil {
			switch {
			case errors.Is(err, store.ErrNotFound):
				app.notFoundResponse(w, r, err)
			default:
				app.internalServerError(w, r, err)
			}
			return
		}

		ctx = context.WithValue(ctx, userCtx, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (app *application) activateUserHandler(w http.ResponseWriter, r *http.Request) {
	token := chi.URLParam(r, "token")
	app.logger.Infow("activating user", "token", token)
	err := app.store.Users.Activate(r.Context(), token)
	if err != nil {
		switch err {
		case store.ErrNotFound:
			app.logger.Errorw("user not found", "error", err)
			app.notFoundResponse(w, r, err)
		default:
			app.logger.Errorw("failed to activate user", "error", err)
			app.internalServerError(w, r, err)
		}
		return
	}

	app.logger.Infow("user activated successfully", "token", token)

	if err := app.jsonResponse(w, http.StatusNoContent, nil); err != nil {
		app.internalServerError(w, r, err)
	}
}

func getUserFromCtx(r *http.Request) *store.User {
	user, _ := r.Context().Value(userCtx).(*store.User)
	return user
}
