package server

import (
	"database/sql"
	"encoding/json"
	"io"
	"net/http"
	"nrs16/cme/middleware"
	"nrs16/cme/repository/entities"
	"nrs16/cme/requests"
	"nrs16/cme/responses"

	"github.com/gocql/gocql"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

func (app *App) SendMessage(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	c := ctx.Value("claims").(middleware.Claims)
	username := c.Username

	var message requests.Message

	/// read body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		er := Error("internal_server_error", "Could not read body")
		Reply(w, r, er, http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	/// unpack into struct
	err = json.Unmarshal(body, &message)
	if err != nil {
		er := Error("bad_request", "Could not unmarshal body")
		Reply(w, r, er, http.StatusBadRequest)
		return
	}

	////  toId XOR chatId
	if message.ChatId == nil && message.ToId == nil {
		er := Error("bad_request", "empty chat_id and empty to_id")
		Reply(w, r, er, http.StatusBadRequest)
		return
	}

	if message.ChatId != nil && message.ToId != nil {
		er := Error("bad_request", "chat_id and to_id are provided, need to pass only one")
		Reply(w, r, er, http.StatusBadRequest)
		return
	}

	if message.ToId != nil {
		/// This is one to one message
		_, err = app.AuthenticationDatabase.GetUserbyUsername(*message.ToId)
		if err != nil {
			if err == sql.ErrNoRows {
				er := Error("bad_request", "to id valid")
				Reply(w, r, er, http.StatusBadRequest)
				return
			}
			er := Error("internal_server_error", "could not check users")
			Reply(w, r, er, http.StatusInternalServerError)
			return
		}

		existingChat, err := app.ChatDatabase.SelectChatbyUsers(username, *message.ToId)
		ChatId := existingChat.Id
		if err != nil {
			///chat does not exist
			c, err := gocql.RandomUUID()
			if err != nil {
				log.Errorf("could not generate uuid: %s", err.Error())
				Reply(w, r, Error("internal_server_error", "Could not create chat"), http.StatusInternalServerError)
				return
			}
			ChatId = &c
			chat := entities.Chat{Id: ChatId, Participant1: username, Participant2: *message.ToId}
			//// insert new chat in DB
			err = app.ChatDatabase.InsertChat(chat)
			if err != nil {
				log.Errorf("error inserting chat: %s", err.Error())
				Reply(w, r, Error("internal_server_error", "Could not create chat"), http.StatusInternalServerError)
				return
			}

			//// cache chat in redis
			err = SetChatInCache(ctx, app.Redis, username, *ChatId)
			if err != nil {
				log.Errorf("error caching chat: %s", err.Error())
				Reply(w, r, Error("internal_server_error", "Could not cache chat"), http.StatusInternalServerError)
				return
			}

			/// we are allowing self chat , so we need to check the below
			if *message.ToId != username {
				err = SetChatInCache(ctx, app.Redis, *message.ToId, *ChatId)
				if err != nil {
					log.Errorf("error caching chat: %s", err.Error())
					Reply(w, r, Error("internal_server_error", "Could not cache chat"), http.StatusInternalServerError)
					return
				}
			}
		}
		messageId, err := gocql.RandomUUID()
		if err != nil {
			log.Errorf("could not generate uuid: %s", err.Error())
			Reply(w, r, Error("internal_server_error", "Could not create chat"), http.StatusInternalServerError)
			return
		}
		messageNew := entities.Message{Id: &messageId, ChatId: *ChatId, FromID: username, Message: message.Message}
		err = app.ChatDatabase.InsertMessage(messageNew)
		if err != nil {
			log.Errorf("error insert message: %s", err.Error())
			Reply(w, r, Error("internal_server_error", "Could not create message"), http.StatusInternalServerError)
			return
		}

	} else {
		/////chatid is sent- this should be used for group message, but also works for one to one as well
		chat, err := app.ChatDatabase.SelectChatbyID(*message.ChatId)
		if err != nil {
			log.Errorf("error selecting chat: %s", err.Error())
			/// even if err is sqlnorows, let's not expose the error because we still don't know if the user belongs to the chat
			Reply(w, r, Error("bad_request", "Bad Request"), http.StatusBadRequest)
			return
		}
		if IsParticipant(username, &chat) {
			messageId, err := gocql.RandomUUID()
			if err != nil {
				log.Errorf("could not generate uuid: %s", err.Error())
				Reply(w, r, Error("internal_server_error", "Could not create chat"), http.StatusInternalServerError)
				return
			}
			messageNew := entities.Message{Id: &messageId, ChatId: *message.ChatId, FromID: username, Message: message.Message}
			err = app.ChatDatabase.InsertMessage(messageNew)
			if err != nil {
				log.Errorf("error insert message on existing chat: %s", err.Error())
				er := Error("internal_server_error", "Could not create message")
				Reply(w, r, er, http.StatusInternalServerError)
				return
			}
		} else {
			er := Error("forbidden", "Forbidden")
			Reply(w, r, er, http.StatusForbidden)
			return
		}
	}
	resp := responses.OKResponse{Message: "success"}
	Reply(w, r, resp, http.StatusCreated)
}

func (app *App) GetChatMessages(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	c := ctx.Value("claims").(middleware.Claims)
	username := c.Username
	vars := mux.Vars(r)

	chatId := vars["chatId"]

	chatUuid, err := gocql.ParseUUID(chatId)
	if err != nil {
		er := Error("Bad_request", "malformed chat id")
		Reply(w, r, er, http.StatusBadRequest)
		return
	}

	chat, err := app.ChatDatabase.SelectChatbyID(chatUuid)
	if err != nil {
		er := Error("Bad_request", "error getting chat")
		Reply(w, r, er, http.StatusBadRequest)
		return
	}

	if username != chat.Participant1 && username != chat.Participant2 {
		er := Error("Forbidden", "Forbidden")
		Reply(w, r, er, http.StatusForbidden)
		return
	}

	messages, err := app.ChatDatabase.SelectMessagesByChatId(chatUuid)
	if err != nil {
		er := Error("internal_server_error", "Error getting chats")
		Reply(w, r, er, http.StatusInternalServerError)
		return
	}

	response := make([]responses.MessageResponse, 0)
	for _, m := range messages {
		response = append(response, responses.MessageResponse{Id: *m.Id, FromID: m.FromID, ChatID: m.ChatId, Message: m.Message, TsCreated: m.TsCreated})
	}

	Reply(w, r, response, http.StatusOK)
}
