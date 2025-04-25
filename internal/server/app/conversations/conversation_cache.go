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
	firstMessageKey     = "fm"
	promptKey           = "prm"
	dynamicVariablesKey = "dv"
	callSidKey          = "sid"
)

func (s *Service) GenerateHash() string {
	return strings.Replace(uuid.NewString(), "-", "", -1)
}

func (s *Service) CreateCache(conversationHash, firstMessage, prompt, callSid string, dynamicVariables map[string]any) {
	jsonData, _ := json.Marshal(dynamicVariables)
	s.conversationCache(conversationHash, map[string]any{
		callSidKey:          callSid,
		firstMessageKey:     firstMessage,
		promptKey:           prompt,
		dynamicVariablesKey: string(jsonData),
	})
}

func (s *Service) GetByHashFromCache(conversationHash string) (string, string, string, map[string]any, error) {
	callSid := s.App.Cache.Prefix(conversationHash).Get(callSidKey)
	if callSid == nil {
		return "", "", "", nil, errors.New("conversation not found by hash")
	}
	firstMessage := s.App.Cache.Prefix(conversationHash).Get(firstMessageKey)
	prompt := s.App.Cache.Prefix(conversationHash).Get(promptKey)
	dynamicVariables := s.App.Cache.Prefix(conversationHash).Get(dynamicVariablesKey)

	var dynamicVariablesMap map[string]any
	fm, _ := firstMessage.(string)
	prm, _ := prompt.(string)
	sid, _ := callSid.(string)
	dvStr, _ := dynamicVariables.(string)
	if dvStr != "" {
		if err := json.Unmarshal([]byte(dvStr), &dynamicVariablesMap); err != nil {
			return "", "", "", nil, err
		}
	}
	return fm, prm, sid, dynamicVariablesMap, nil
}

func (s *Service) Delete(conversationHash string) {
	s.clearConversationCache(conversationHash, []string{
		callSidKey,
		firstMessageKey,
		promptKey,
		dynamicVariablesKey,
	})
}

func (s *Service) conversationCache(conversationHash string, data map[string]any) {
	s.App.Cache.Prefix(conversationHash).SetManyWithTTL(data, time.Minute)

}

func (s *Service) clearConversationCache(conversationHash string, keys []string) {
	s.App.Cache.Prefix(conversationHash).DeleteMany(keys)

}
