package token

import "time"

type Maker interface {
	CreateAccessToken(uid string, duration time.Duration) (string, error)
	CreateRefreshToken() (string, error)

	VerifyToken(token string) (*Payload, error)
}
