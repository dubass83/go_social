package main

import (
	"fmt"
	"math/rand/v2"
	"net/http"
	"time"

	"github.com/dubass83/go_social/internal/mailer"
	"github.com/dubass83/go_social/internal/store"
	"github.com/dubass83/go_social/internal/util"
	"github.com/golang-jwt/jwt/v5"
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

	// always create new user with basic permissions
	userRole, err := app.store.Role.GetByName(r.Context(), "user")
	if err != nil {
		internalServerError(w, r, err)
		return
	}

	user := &store.User{
		Username:        payload.Username,
		Email:           payload.Email,
		Password:        hashedPassword,
		ActivationToken: util.GenerateToken(rand.Int64N(100)),
		RoleID:          int(userRole.ID),
	}

	if err := app.store.User.CreateAndInviteTx(r.Context(), user); err != nil {
		internalServerError(w, r, err)
		return
	}

	activationURL := fmt.Sprintf("%s/confirm/%s", app.config.frontendURL, user.ActivationToken)

	// send email
	if err := app.mailer.Send(mailer.Message{
		To:       []string{user.Email},
		Subject:  "Welcome to Go Social!",
		Data:     activationURL,
		Template: "confirmation-email",
	}); err != nil {
		log.Error().Err(err).Msg("failed to send email")

		// rollback user creation if email sending fails (SAGA patern)
		log.Debug().Msgf("clean activation token and remove inactive user: %s from database", user.Email)
		err := app.store.User.DeleteByID(r.Context(), user.ID)
		if err != nil {
			log.Error().Err(err).Msg("failed to delete user")
		}
		err = app.store.Invitation.CleanByID(r.Context(), user.ID)
		if err != nil {
			log.Error().Err(err).Msg("failed to clean invitation")
		}
		return
	}

	if err := app.jsonResponse(w, http.StatusCreated, user); err != nil {
		internalServerError(w, r, err)
		return
	}
}

// createTokenHandler godoc
//
//	@Summary		Generate JWT token for existing user
//	@Description	Generate JWT token for existing user
//	@Tags			AUTH
//	@Accept			json
//	@Produce		json
//	@Param			payload	body		createTokenPayload	true	"User credentials"
//	@Success		200		{string}	string			"Token"
//	@Failure		400		{object}	error
//	@Failure		401		{object}	error
//	@Failure		500		{object}	error
//	@Router			/authentication/token [post]
func (app *application) createTokenHandler(w http.ResponseWriter, r *http.Request) {
	var payload createTokenPayload
	errRJ := readJSON(w, r, &payload)
	if errRJ != nil {
		badRequestResponse(w, r, errRJ)
		return
	}

	errVal := validate.Struct(payload)
	if errVal != nil {
		badRequestResponse(w, r, errVal)
		return
	}

	user, err := app.store.User.GetByEmail(r.Context(), payload.Email)
	if err != nil {
		switch err {
		case store.ErrNotFound:
			unAuthorizedResponse(w, r, err)
			return
		default:
			internalServerError(w, r, err)
			return
		}
	}

	errPass := util.CheckPassword(payload.Password, user.Password)
	if errPass != nil {
		unAuthorizedResponse(w, r, errPass)
		return
	}

	claims := jwt.MapClaims{
		"sub": user.ID,
		"iss": tokenHost,
		"aud": []string{tokenHost},
		"exp": time.Now().Add(app.config.auth.jwt.expiry).Unix(),
		"iat": time.Now().Unix(),
		"nbf": time.Now().Unix(),
	}

	token, errTk := app.authenticator.GenerateToken(claims)
	if errTk != nil {
		internalServerError(w, r, errTk)
		return
	}

	errResp := app.jsonResponse(w, http.StatusOK, token)
	if errResp != nil {
		internalServerError(w, r, errResp)
		return
	}
}
