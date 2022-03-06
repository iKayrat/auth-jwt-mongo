package token

import (
	"encoding/base64"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/iKayrat/auth-jwt-mongo/util"
	"golang.org/x/crypto/bcrypt"
)

const minSeceretKeySize = 32

type JwtMaker struct {
	secretkey []byte
}

func NewJWTMaker() (Maker, error) {

	// load config parametres from env
	config, err := util.LoadConfig("../")
	if err != nil {
		log.Fatal("Cannot load config:", err)
	}

	if len(config.SECRET_KEY) < minSeceretKeySize {
		return nil, fmt.Errorf("key size must be at least %d characters", minSeceretKeySize)
	}

	return &JwtMaker{[]byte(config.SECRET_KEY)}, nil
}

//
func (maker *JwtMaker) CreateAccessToken(uid string, duration time.Duration) (string, error) {

	accessPayload, err := NewPayload(uid, duration)
	if err != nil {
		return "", err
	}

	accessjwt := jwt.NewWithClaims(jwt.SigningMethodHS512, accessPayload)
	access, err := accessjwt.SignedString(maker.secretkey)
	if err != nil {
		return "", err
	}

	return access, nil
}

func (m *JwtMaker) VerifyToken(token string) (*Payload, error) {
	keyfunc := func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, ErrInvalidToken
		}
		return m.secretkey, nil
	}
	jwtToken, err := jwt.ParseWithClaims(token, &Payload{}, keyfunc)
	if err != nil {
		verr, ok := err.(*jwt.ValidationError)
		if ok && errors.Is(verr.Inner, ErrExpiredToken) {
			return nil, ErrExpiredToken
		}
		return nil, ErrInvalidToken
	}

	payload, ok := jwtToken.Claims.(*Payload)
	if !ok {
		return nil, ErrInvalidToken
	}
	return payload, nil

}

func (maker *JwtMaker) CreateRefreshToken() (string, error) {
	// empty slice of byte
	b := make([]byte, 32)

	// new random time parameters
	s := rand.NewSource(time.Now().Unix())
	r := rand.New(s)

	// generate into []byte
	_, err := r.Read(b)
	if err != nil {
		return "", err
	}

	// encode into base64 format
	str := base64.StdEncoding.EncodeToString(b)
	fmt.Println(str)

	// hash
	refresh, err := bcrypt.GenerateFromPassword([]byte(str), 10)
	if err != nil {
		log.Fatal("Generate secret key err:", err)
	}

	return string(refresh), nil
}
