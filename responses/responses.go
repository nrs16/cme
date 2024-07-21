package responses

import (
	"time"

	"github.com/gocql/gocql"
)

type ErrorResponse struct {
	Code    string `json:"error_code"`
	Message string `json:"message"`
}

type RegistrationSuccess struct {
	Token string `json:"token"`
}

type OKResponse struct {
	Message string `json:"message"`
}

type ChatResponse struct {
	Id           gocql.UUID `json:"chat_id"`
	TsCreated    time.Time  `json:"ts_created"`
	Participant1 string     `json:"participant_1"`
	Participant2 string     `json:"participant_2"`
}

type MessageResponse struct {
	Id        gocql.UUID `json:"id"`
	ChatID    gocql.UUID `json:"chat_id"`
	FromID    string     `json:"from_id"`
	Message   string     `json:"message"`
	TsCreated time.Time  `json:"ts_created"`
}
