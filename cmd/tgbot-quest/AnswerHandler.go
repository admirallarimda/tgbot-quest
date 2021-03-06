package main

import (
	"regexp"

	"github.com/admirallarimda/tgbot-quest/internal/pkg/quest"
	"github.com/admirallarimda/tgbotbase"
	log "github.com/sirupsen/logrus"
	"gopkg.in/telegram-bot-api.v4"
)

type answerHandler struct {
	tgbotbase.BaseHandler
	engine quest.QuestEngine
}

func (h *answerHandler) Name() string {
	return "answer handler"
}

func (h *answerHandler) HandleOne(msg tgbotapi.Message) {
	userID := tgbotbase.UserID(msg.From.ID)
	chatID := msg.Chat.ID
	logger := log.WithFields(log.Fields{"userID": userID, "userName": msg.From.UserName, "message": msg.Text})
	logger.Debug("Incoming answer")
	res := h.engine.CheckAnswer(userID, msg.Text)
	if !res.Active {
		logger.Debug("No Active quests, skipping")
		return
	}

	if !res.Correct {
		h.OutMsgCh <- tgbotapi.NewMessage(chatID, "Ответ неверный!")
	} else {
		h.OutMsgCh <- tgbotapi.NewMessage(chatID, "Правильно!")
		if res.Finished {
			h.OutMsgCh <- tgbotapi.NewMessage(chatID, "Это был последний вопрос. Ты молодец!")
		} else {
			h.OutMsgCh <- h.engine.GetCurrentQuestion(userID)
		}
	}
}

func (h *answerHandler) Init(outCh chan<- tgbotapi.Chattable, srvCh chan<- tgbotbase.ServiceMsg) tgbotbase.HandlerTrigger {
	h.OutMsgCh = outCh
	return tgbotbase.NewHandlerTrigger(regexp.MustCompile("^[^/].*"), nil)
}

func NewAnswerHandler(engine quest.QuestEngine) tgbotbase.IncomingMessageHandler {
	return &answerHandler{engine: engine}
}
