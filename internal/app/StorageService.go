package app

import (
	"github.com/google/uuid"
	spec "password-storage/internal/pkg/models"
	"password-storage/internal/pkg/models/db"
	repository "password-storage/internal/pkg/repository"
	"time"
)

var Con = repository.NewDBConn()

func CreateNewPasswordData(request spec.NewPasswordDataRequest) spec.NewPasswordDataResponse {
	key := *request.Key
	newPasswordData := db.PasswordData{
		Key:   key,
		Value: *request.Value,
	}
	err := repository.CreateNewPasswordData(Con, newPasswordData)
	if err != nil {
		return spec.NewPasswordDataResponse{}
	}
	response := spec.NewPasswordDataResponse{
		Key: &key,
	}
	return response
}

func CreateNewSessionWithTelegramId(telegramId int64) string {
	sessionId := uuid.New().String()
	newSession := db.Session{
		SessionId:  sessionId,
		TelegramId: telegramId,
		ExpiredAt:  time.Now().Add(15 * time.Minute),
	}
	err := repository.CreateNewSession(Con, newSession)
	if err != nil {
		return ""
	}
	return sessionId
}

func CreateNewSession() string {
	sessionId := uuid.New().String()
	newSession := db.Session{
		SessionId: sessionId,
		ExpiredAt: time.Now().Add(15 * time.Minute),
	}
	err := repository.CreateNewSession(Con, newSession)
	if err != nil {
		return ""
	}
	return sessionId
}

func IsSessionValid(sessionId string) bool {
	expiredAt := repository.GetSessionExpiredAt(Con, sessionId)
	if expiredAt.Before(time.Now()) {
		return false
	} else {
		return true
	}
}

func IsSessionValidByTelegramId(telegramId int64) bool {
	expiredAt := repository.GetSessionExpiredAtByTelegramId(Con, telegramId)
	if expiredAt.Before(time.Now()) {
		return false
	} else {
		return true
	}
}

func GetAllPasswordData() []spec.AllPasswordData {
	pd, err := repository.GetAllPasswordData(Con)
	if err != nil {
		return nil
	}
	response := make([]spec.AllPasswordData, 0)
	for i := range pd {
		response = append(response, spec.AllPasswordData{
			{
				Key:   &pd[i].Key,
				Value: &pd[i].Value,
			},
		})
	}
	return response
}

func GetPasswordDataByKey(key interface{}) spec.PasswordDataByKey {
	byKey, err := repository.GetPasswordDataByKey(Con, key)
	if err != nil {
		return spec.PasswordDataByKey{}
	}
	pwdByKey := spec.PasswordDataByKey{
		Key:   &byKey.Key,
		Value: &byKey.Value,
	}
	return pwdByKey
}
