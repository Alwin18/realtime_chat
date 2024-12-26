package api

import (
	"github.com/google/uuid"
)

type BaseResponse[T any] struct {
	Status  bool   `json:"status"`
	Message string `json:"message"`
	Code    int64  `json:"code"`
	Data    T      `json:"data"`
}

func NewBaseResponse[T any](status bool, message string, code int64, data T) BaseResponse[T] {
	return BaseResponse[T]{
		Status:  status,
		Message: message,
		Code:    code,
		Data:    data,
	}
}

type ResponseWithoutData struct {
	Status  bool   `json:"status"`
	Message string `json:"message"`
	Code    int64  `json:"code"`
}

func NewResponseWithoutData(message string, code int64, status bool) ResponseWithoutData {
	return ResponseWithoutData{
		Status:  status,
		Message: message,
		Code:    code,
	}
}

type GetContactByCakupanRequest struct {
	GedungID int64  `query:"gedung_id"`
	Role     string `query:"role"`
}

type GetContactByCakupanResponse struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Email       *string   `json:"email,omitempty"`
	PhoneNumber string    `json:"phone_number"`
	AvatarURL   *string   `json:"avatar_url,omitempty"`
	KotaID      int64     `json:"kota_id"`
	GedungID    int64     `json:"gedung_id"`
	IsOnline    bool      `json:"is_online"`
}

type GetChatHistoryRequest struct {
	SenderID   string `query:"sender_id"`
	ReceiverID string `query:"receiver_id"`
}
