package conversations

import (
	"strings"

	"github.com/google/uuid"
)

const (
	firstMessageKey = "fm"
	promptKey       = "prm"
	deleteCache     = "delete"
	callSid         = "sid"
)

/*
	type CurrentConversation struct {
		firstMessageKey     string
		prompt              string
		callSid             string
		incomingPhoneNumber string
		outgoingPhoneNumber string
	}
*/
func (s *Service) GenerateHash() string {
	return strings.Replace(uuid.NewString(), "-", "", -1)
}

func (s *Service) CreateCache(conversationHash, firstMessage, prompt string) {
	s.conversationCache(conversationHash, map[string]string{
		firstMessageKey: firstMessage,
		promptKey:       prompt,
	})
}

func (s *Service) GetByHashFromCache(conversationHash string) (string, string) {
	firstMessage := s.App.Cache.Prefix(conversationHash).Get(firstMessageKey)
	prompt := s.App.Cache.Prefix(conversationHash).Get(promptKey)
	return firstMessage.(string), prompt.(string)
}

func (s *Service) Delete(conversationHash string) {
	s.conversationCache(conversationHash, map[string]string{
		firstMessageKey: deleteCache,
		promptKey:       deleteCache,
	})
}

func (s *Service) conversationCache(conversationHash string, data map[string]string) {
	for k, v := range data {
		if v == deleteCache {
			s.App.Cache.Prefix(conversationHash).Delete(k)
		} else {
			s.App.Cache.Prefix(conversationHash).Set(k, v)
		}
	}
}
