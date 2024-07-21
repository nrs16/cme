package server

import (
	"encoding/json"
	"io"
	"net/http"
	"nrs16/cme/metrics"
	"nrs16/cme/middleware"
	"nrs16/cme/requests"
	"nrs16/cme/responses"

	"github.com/golang-jwt/jwt"
	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
)

func (app *App) Login(w http.ResponseWriter, r *http.Request) {
	metrics.HttpRequestsTotal.WithLabelValues(r.URL.Path).Inc()
	t := prometheus.NewTimer(metrics.HttpRequestDuration.WithLabelValues(r.URL.Path))
	defer t.ObserveDuration()
	var payload requests.LoginBody

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

	/// get salt and hash from B
	user, err := app.AuthenticationDatabase.GetUserbyUsername(payload.Username)
	if err != nil {
		log.Errorf("error getting user: %s", err.Error())
		er := Error("internal_server_error", "Could not retrieve password information")
		_ = Reply(w, r, er, http.StatusInternalServerError)
		return
	}
	if user.PasswordHash == nil || user.PasswordSalt == nil {
		er := Error("bad_request", "Please set password first")
		_ = Reply(w, r, er, http.StatusBadRequest)
		return
	}
	/// verify password
	err = middleware.VerifyPassword(payload.Password, *user.PasswordSalt, *user.PasswordHash)
	if err != nil {
		er := Error("bad_request", "Wrong password")
		_ = Reply(w, r, er, http.StatusBadRequest)
		return
	}

	//// generate token and return
	claims := jwt.MapClaims{"username": user.Username}
	token, err := middleware.GenerateJWT(claims)
	if err != nil {
		log.Error(err.Error())
		er := Error("internal_server_error", "could not create token")
		_ = Reply(w, r, er, http.StatusInternalServerError)
		return
	}
	resp := responses.RegistrationSuccess{Token: token}
	_ = Reply(w, r, resp, http.StatusCreated)
}
