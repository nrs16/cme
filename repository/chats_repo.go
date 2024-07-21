package repository

import (
	"errors"
	"nrs16/cme/repository/entities"

	"github.com/gocql/gocql"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

type Chat interface {
	InsertChat(chat entities.Chat) error
	SelectChat(Id uuid.UUID) (entities.Chat, error)
	InsertMessage(message entities.Message) error
}

type ChatsDatabase struct {
	Database *gocql.Session
}

func (c *ChatsDatabase) InsertChat(chat entities.Chat) error {
	return c.Database.Query(`insert into chat (id, participant_1, participant_2, ts_created) values (?,?,?, toTimestamp(now()))`, chat.Id, chat.Participant1, chat.Participant2).Exec()

}

func (c *ChatsDatabase) SelectChatbyID(Id gocql.UUID) (entities.Chat, error) {
	var chat entities.Chat
	iter := c.Database.Query(`select id, participant_1, participant_2 from chat where id = ?`, Id).Iter()

	for iter.Scan(&chat.Id, &chat.Participant1, &chat.Participant2) {
		if chat.Id != nil {
			return chat, nil
		}
	}

	return entities.Chat{}, errors.New("not_found")
}

func (c *ChatsDatabase) SelectChatbyUsers(user1 string, user2 string) (entities.Chat, error) {
	var chat entities.Chat
	chat.Participant1 = user1
	chat.Participant2 = user2
	err := c.Database.Query(`select id from chat where participant_1 = ? and participant_2 = ? ALLOW FILTERING `, user1, user2).Scan(&chat.Id)
	if err != nil {
		log.Warn(err.Error())
		err = c.Database.Query(`select id from chat where participant_1 = ? and participant_2 = ? ALLOW FILTERING `, user2, user1).Scan(&chat.Id)
	}
	return chat, err
}

func (c *ChatsDatabase) InsertMessage(message entities.Message) error {
	return c.Database.Query(`insert into message (id, chat_id, message, from_id, ts_created) values (?,?,?,?, toTimestamp(now()))`, message.Id, message.ChatId, message.Message, message.FromID).Exec()
}

func (c *ChatsDatabase) SelectChatsByListIds(chatIds []gocql.UUID) ([]entities.Chat, error) {
	chats := make([]entities.Chat, 0)
	results := c.Database.Query(`select id, participant_1, participant_2, ts_created from chat where id in ? `, chatIds).Iter()

	var chat entities.Chat

	for results.Scan(&chat.Id, &chat.Participant1, &chat.Participant2, &chat.TsCreated) {
		if chat.Id != nil {
			chats = append(chats, chat)
		}

	}
	_ = results.Close() /*err != nil {
		log.Errorf("error in scanning message results: %s", err.Error())
		return chats, err
	}*/

	return chats, nil
}

func (c *ChatsDatabase) SelectMessagesByChatId(chatId gocql.UUID) ([]entities.Message, error) {
	messages := make([]entities.Message, 0)
	results := c.Database.Query(`select id, from_id, message, chat_id, ts_created from message where chat_id = ? `, chatId).Iter()

	var message entities.Message

	for results.Scan(&message.Id, &message.FromID, &message.Message, &message.ChatId, &message.TsCreated) {
		if message.Id != nil {
			messages = append(messages, message)
		}

	}

	_ = results.Close() /*; err != nil {
		log.Errorf("error in scanning results: %s", err.Error())
		return messages, err
	}*/
	return messages, nil
}
