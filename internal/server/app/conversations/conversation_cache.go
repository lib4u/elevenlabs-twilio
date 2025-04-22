package conversations

import (
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
	jsoniter "github.com/json-iterator/go"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

const (
	conversationExistKey = "ch"
	firstMessageKey      = "fm"
	promptKey            = "prm"
	deleteCache          = "delete"
	dynamicVariablesKey  = "dv"
	callSid              = "sid"
)

/*
	type CurrentConversation struct {
		firstMessageKey     string
		prompt              string
		callSid             string
		incomingPhoneNumber string
		outgoingPhoneNumber string
		dynamicVariablesKey string
	}
*/

func (s *Service) GenerateHash() string {
	return strings.Replace(uuid.NewString(), "-", "", -1)
}

func (s *Service) CreateCache(conversationHash, firstMessage, prompt string, DynamicVariables map[string]any) {
	jsonData, _ := json.Marshal(DynamicVariables)
	s.conversationCache(conversationHash, map[string]string{
		conversationExistKey: "ok",
		firstMessageKey:      firstMessage,
		promptKey:            prompt,
		dynamicVariablesKey:  string(jsonData),
	})
}

func (s *Service) GetByHashFromCache(conversationHash string) (string, string, map[string]any, error) {
	conversationExist := s.App.Cache.Prefix(conversationHash).Get(conversationExistKey)
	if conversationExist == nil {
		return "", "", nil, errors.New("conversation not found by hash")
	}
	firstMessage := s.App.Cache.Prefix(conversationHash).Get(firstMessageKey)
	prompt := s.App.Cache.Prefix(conversationHash).Get(promptKey)
	dynamicVariables := s.App.Cache.Prefix(conversationHash).Get(dynamicVariablesKey)
	var dynamicVariablesMap map[string]any

	_ = json.Unmarshal([]byte(dynamicVariables.(string)), &dynamicVariablesMap)
	return firstMessage.(string), prompt.(string), dynamicVariablesMap, nil
}

func (s *Service) Delete(conversationHash string) {
	s.conversationCache(conversationHash, map[string]string{
		firstMessageKey:     deleteCache,
		promptKey:           deleteCache,
		dynamicVariablesKey: deleteCache,
	})
}

func (s *Service) conversationCache(conversationHash string, data map[string]string) {
	for k, v := range data {
		if v == deleteCache {
			s.App.Cache.Prefix(conversationHash).Delete(k)
		} else {
			s.App.Cache.Prefix(conversationHash).SetWithTTL(k, v, time.Minute)
		}
	}
}
