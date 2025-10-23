package main

import (
	"math/rand/v2"
	"net/http"

	"github.com/dubass83/go_social/internal/mailer"
	"github.com/dubass83/go_social/internal/store"
	"github.com/dubass83/go_social/internal/util"
	"github.com/rs/zerolog/log"
)

// registerUserHandler godoc
//
//	@Summary		Registers a new user
//	@Description	Registers a new user
//	@Tags			AUTH
//	@Accept			json
//	@Produce		json
//	@Param			payload	body		registerUserPayload	true	"User credentials"
//	@Success		201		{object}	store.User			"user registered"
//	@Failure		400		{object}	error
//	@Failure		500		{object}	error
//	@Router			/authentication/user [post]
func (app *application) registerUserHandler(w http.ResponseWriter, r *http.Request) {
	var payload registerUserPayload
	err := readJSON(w, r, &payload)
	if err != nil {
		badRequestResponse(w, r, err)
		return
	}
	if err = validate.Struct(payload); err != nil {
		badRequestResponse(w, r, err)
		return
	}

	hashedPassword, err := util.HashPassword(payload.Password)
	if err != nil {
		badRequestResponse(w, r, err)
		return
	}

	user := &store.User{
		Username:        payload.Username,
		Email:           payload.Email,
		Password:        hashedPassword,
		ActivationToken: util.GenerateToken(rand.Int64N(100)),
	}

	if err := app.store.User.CreateAndInviteTx(r.Context(), user); err != nil {
		internalServerError(w, r, err)
		return
	}

	// send email
	if err := app.mailer.Send(mailer.Message{
		To:       []string{user.Email},
		Subject:  "Welcome to Go Social!",
		Data:     user.ActivationToken,
		Template: "confirmation-email",
	}); err != nil {
		log.Error().Err(err).Msg("failed to send email")
	}

	if err := app.jsonResponse(w, http.StatusCreated, user); err != nil {
		internalServerError(w, r, err)
		return
	}
}
