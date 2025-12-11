package main

import (
	"net/http"

	"github.com/youssefM1999/social/internal/store"
)

// GetUserFeed godoc
//
//	@Summary		Get user feed
//	@Description	Get paginated feed of posts for a user
//	@Tags			feed
//	@Accept			json
//	@Produce		json
//	@Param			limit	query		int		false	"Limit"		default(20)
//	@Param			offset	query		int		false	"Offset"	default(0)
//	@Param			sort	query		string	false	"Sort"		default(desc)
//	@Param			tags	query		string	false	"Tags (comma separated)"
//	@Param			search	query		string	false	"Search"
//	@Success		200		{object}	[]store.Post
//	@Failure		400		{object}	error
//	@Failure		500		{object}	error
//	@Router			/users/feed [get]
func (app *application) getUserFeedHandler(w http.ResponseWriter, r *http.Request) {
	// pagination, filters, sort
	fq := store.PaginatedFeedQuery{
		Limit:  20,
		Offset: 0,
		Sort:   "desc",
	}

	fq, err := fq.Parser(r)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if err := Validate.Struct(fq); err != nil {
		app.badRequestResponse(w, r, err)
	}

	ctx := r.Context()

	feed, err := app.store.Posts.GetUserFeed(ctx, int64(2), fq)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.jsonResponse(w, http.StatusOK, feed); err != nil {
		app.internalServerError(w, r, err)
	}
}
