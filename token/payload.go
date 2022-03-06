package token

import (
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/x/mongo/driver/uuid"
)

var (
	ErrInvalidToken = errors.New("token is invalid")
	ErrExpiredToken = errors.New("token is expired")
)

type Payload struct {
	TokenID   uuid.UUID `json:"token_id"`
	User_uuid string    `json:"user_uid"`
	IssuedAt  time.Time `json:"issued_at"`
	ExpiredAt time.Time `json:"expired_at"`
}

func NewPayload(uid string, duration time.Duration) (*Payload, error) {

	tokenid, err := uuid.New()
	if err != nil {
		return nil, err
	}

	payload := &Payload{
		TokenID:   tokenid,
		User_uuid: uid,
		IssuedAt:  time.Now(),
		ExpiredAt: time.Now().Add(duration),
	}

	return payload, err
}

func (payload *Payload) Valid() error {
	if time.Now().After(payload.ExpiredAt) {
		return ErrExpiredToken
	}
	return nil
}
