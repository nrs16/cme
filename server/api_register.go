package server

import (
	"encoding/json"
	"io"
	"net/http"
	"nrs16/cme/metrics"
	"nrs16/cme/middleware"
	"nrs16/cme/repository/entities"
	"nrs16/cme/requests"
	"nrs16/cme/responses"

	"github.com/golang-jwt/jwt"
	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
)

func (app *App) Register(w http.ResponseWriter, r *http.Request) {
	metrics.HttpRequestsTotal.WithLabelValues(r.URL.Path).Inc()
	t := prometheus.NewTimer(metrics.HttpRequestDuration.WithLabelValues(r.URL.Path))
	defer t.ObserveDuration()
	ctx := r.Context()

	var payload requests.RegisterBody

	/// read body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		er := Error("internal_server_error", "Could not read body")
		_ = Reply(w, r, er, http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	/// unpack into struct
	err = json.Unmarshal(body, &payload)
	if err != nil {
		er := Error("bad_request", "Could not unmarshal body")
		_ = Reply(w, r, er, http.StatusBadRequest)
		return
	}

	//// validate body
	err = payload.Validate()
	if err != nil {
		er := Error("bad_request", err.Error())
		_ = Reply(w, r, er, http.StatusBadRequest)
		return
	}

	//// check username is unique through redis

	usernames, err := GetUsersFromCache(ctx, app.Redis)
	if err != nil {
		log.Errorf("fetch usernames error: %s", err.Error())
		er := Error("internal_server_error", err.Error())
		_ = Reply(w, r, er, http.StatusInternalServerError)
		return
	}

	if !UsernameAllowed(payload.Username, usernames) {
		er := Error("bad_request", "username already taken")
		_ = Reply(w, r, er, http.StatusBadRequest)
		return
	}

	//// generate password salt and hash
	salt, hash, err := middleware.HashPassword(payload.Password)
	if err != nil {
		log.Errorf("password error: %s", err.Error())
		er := Error("internal_server_error", err.Error())
		_ = Reply(w, r, er, http.StatusInternalServerError)
		return
	}

	//// map payload to user entities
	user := entities.User{Username: payload.Username,
		FirstName:    payload.FirstName,
		LastName:     payload.LastName,
		PasswordSalt: &salt,
		PasswordHash: &hash,
		Email:        payload.EmailAddress}

	/// insert user in DB:
	err = app.AuthenticationDatabase.Register(user)
	if err != nil {
		log.Errorf("insert user error: %s", err.Error())
		er := Error("internal_server_error", "Could not register user")
		_ = Reply(w, r, er, http.StatusInternalServerError)
		return
	}
	///// add username to cache list

	err = AddUsernameToCache(ctx, app.Redis, payload.Username)
	if err != nil {
		log.Warnf("could not add usename to cache: %s", err.Error())
	}
	//// get token
	claims := jwt.MapClaims{"username": user.Username}
	token, err := middleware.GenerateJWT(claims)
	if err != nil {
		log.Errorf("jwt error: %s", err.Error())
		er := Error("internal_server_error", "could not create token")
		_ = Reply(w, r, er, http.StatusInternalServerError)
		return
	}
	resp := responses.RegistrationSuccess{Token: token}
	_ = Reply(w, r, resp, http.StatusCreated)

}
