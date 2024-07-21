package server

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"nrs16/cme/repository/entities"
	"nrs16/cme/responses"

	"github.com/gocql/gocql"
	"github.com/redis/go-redis/v9"
	log "github.com/sirupsen/logrus"
)

func Error(code string, message string) responses.ErrorResponse {
	return responses.ErrorResponse{Code: code, Message: message}
}
func Marshal(a interface{}) ([]byte, error) {
	return json.Marshal(a)
}

func Reply(w http.ResponseWriter, r *http.Request, reply interface{}, status int) error {
	b, err := Marshal(reply)
	if err != nil {
		return err
	}
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(b)
	return nil
}

func SetChatInCache(ctx context.Context, redisCon *redis.Client, user string, chatId gocql.UUID) error {

	//// we keep track if user chats in cache we they are faster to access on the get /request
	userChats := make([]gocql.UUID, 0)

	result, err := redisCon.Get(ctx, ChatRedisKey(user)).Result()
	if err != nil && err != redis.Nil {
		log.Errorf("error reading chat from redis: %s", err.Error())
		return err
	}
	if err == nil { //skip the below if not found, nothing to unmarshal
		err = json.Unmarshal([]byte(result), &userChats)
		if err != nil {
			log.Errorf("error unmarshaling: %s", err.Error())
			return err
		}
	}

	userChats = append(userChats, chatId)
	chatsJSON, err := json.Marshal(userChats)
	if err != nil {
		log.Errorf("error marshaling chats: %v", err)
		return err
	}
	status := redisCon.Set(ctx, ChatRedisKey(user), chatsJSON, 0)
	if status.Err() != nil {
		log.Errorf("error caching one to one chat: %s", status.Err())
		return status.Err()

	}
	return nil
}

func GetChatsFromCache(ctx context.Context, redisCon *redis.Client, username string) ([]gocql.UUID, error) {
	var chats []gocql.UUID
	key := ChatRedisKey(username)

	userChats, err := redisCon.Get(ctx, key).Result()

	if err != nil {
		if err == redis.Nil {
			return chats, errors.New("not_found")
		}
		log.Errorf("error check chat: %s", err.Error())
		return chats, err
	}
	err = json.Unmarshal([]byte(userChats), &chats)
	if err != nil {
		log.Errorf("invalid error unmarshaling: %s", err.Error())
		return chats, err
	}

	return chats, nil
}

func ChatRedisKey(username string) string {
	return fmt.Sprintf("%s_chats", username)
}

func IsParticipant(username string, chat *entities.Chat) bool {
	if username == chat.Participant1 || username == chat.Participant2 {
		return true
	}
	return false
}

func AddUsernameToCache(ctx context.Context, redisCon *redis.Client, user string) error {

	//// we keep track if user chats in cache we they are faster to access on the get /request
	users := make([]string, 0)

	result, err := redisCon.Get(ctx, "users").Result()
	if err != nil && err != redis.Nil {
		log.Errorf("error reading chat from redis: %s", err.Error())
		return err
	}
	if err == nil { //skip the below if not found, nothing to unmarshal
		err = json.Unmarshal([]byte(result), &users)
		if err != nil {
			log.Errorf("error unmarshaling: %s", err.Error())
			return err
		}
	}

	users = append(users, user)
	usersJSON, err := json.Marshal(users)
	if err != nil {
		log.Errorf("error marshaling chats: %v", err)
		return err
	}
	status := redisCon.Set(ctx, "users", usersJSON, 0)
	if status.Err() != nil {
		log.Errorf("error caching one to one chat: %s", status.Err())
		return status.Err()

	}
	return nil
}

func GetUsersFromCache(ctx context.Context, redisCon *redis.Client) ([]string, error) {
	users := make([]string, 0)

	usernames, err := redisCon.Get(ctx, "users").Result()

	if err != nil {
		if err == redis.Nil {
			return users, nil
		}
		log.Errorf("error check chat: %s", err.Error())
		return users, err
	}
	err = json.Unmarshal([]byte(usernames), &users)
	if err != nil {
		log.Errorf("invalid error unmarshaling: %s", err.Error())
		return users, err
	}

	return users, nil
}

func UsernameAllowed(username string, usernames []string) bool {
	for _, u := range usernames {
		if u == username {
			return false
		}
	}
	return true
}
