package entities

import (
	"time"

	"github.com/gocql/gocql"
)

type GroupChat struct {
	Id           gocql.UUID `json:"chat_id"`
	TsCreated    time.Time  `json:"ts_created"`
	Participants []string   `json:"participants"`
}

type Chat struct {
	Id           *gocql.UUID `json:"chat_id"`
	TsCreated    time.Time   `json:"ts_created"`
	Participant1 string      `json:"participant_1"`
	Participant2 string      `json:"participant_2"`
}

type Message struct {
	Id        *gocql.UUID `json:"message_id"`
	ChatId    gocql.UUID  `json:"chatId"`
	Message   string      `json:"message"`
	FromID    string      `json:"from"`
	TsCreated time.Time   `json:"ts_created"`
}
