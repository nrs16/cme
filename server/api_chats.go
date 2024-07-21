package server

import (
	"net/http"
	"nrs16/cme/metrics"
	"nrs16/cme/middleware"
	"nrs16/cme/responses"

	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
)

func (app *App) GetChats(w http.ResponseWriter, r *http.Request) {
	metrics.HttpRequestsTotal.WithLabelValues(r.URL.Path).Inc()
	t := prometheus.NewTimer(metrics.HttpRequestDuration.WithLabelValues(r.URL.Path))
	defer t.ObserveDuration()
	ctx := r.Context()
	c := ctx.Value("claims").(middleware.Claims)
	username := c.Username

	chatsList, err := GetChatsFromCache(ctx, app.Redis, username)
	if err != nil && err.Error() != "not_found" {
		log.Errorf("error getting chats: %s", err.Error())
		er := Error("internal_server_error", "Error getting chats")
		Reply(w, r, er, http.StatusInternalServerError)
		return
	}
	chats, err := app.ChatDatabase.SelectChatsByListIds(chatsList)
	if err != nil {
		log.Errorf("error selecting chats: %s", err.Error())
		er := Error("internal_server_error", "error selecting chats")
		Reply(w, r, er, http.StatusInternalServerError)
		return
	}

	//// mapping entities struct to responses struct should be done in a separate function
	response := make([]responses.ChatResponse, 0)
	for _, c := range chats {
		response = append(response, responses.ChatResponse{Id: *c.Id, Participant1: c.Participant1, Participant2: c.Participant2, TsCreated: c.TsCreated})
	}
	Reply(w, r, response, http.StatusOK)
}
