package requests

import (
	"github.com/gocql/gocql"
)

type Message struct {
	ChatId *gocql.UUID `json:"chat_id"`
	/// ToId is actually the username, I am not using Ids as PK, username is PK
	ToId    *string `json:"to_id"`
	Message string  `json:"message"`
}
