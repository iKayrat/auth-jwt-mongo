package api

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// NewToken
func (s *Server) NewToken(ctx *gin.Context) {
	// id from parameters
	id := ctx.Query("id")

	fmt.Println("get id:", id)

	// parse id into a uuid
	uid, err := uuid.Parse(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"Parsing UUID err": err.Error(),
		})
		return
	}

	// generate new access token
	tokens, err := s.Store.CreateTokens(ctx, uid.String(), time.Minute*15)
	if err != nil {
		fmt.Println("create token err:", err)
		ctx.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	// response tokens map[]string
	ctx.JSON(http.StatusOK, tokens)
}

// @ RefreshToken
// @ gets refresh token from body, finds in db,
// @ if its exp valid, creates new pair of tokens
func (s *Server) RefreshToken(ctx *gin.Context) {

	type Body struct {
		Token string `json:"refresh_token" binding:"required"`
	}

	body := Body{}

	// gets body request json refresh_token string
	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	// checks if its exist
	user, err := s.Store.FindUser(ctx, "refresh_token", body.Token)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errors.New("not verified refersh token"))
		return
	}

	// checks expiration time od refresh token
	// if expired
	if err := user.Valid(); err != nil {
		_, err := s.Store.UserCollection.DeleteOne(ctx, user)
		if err != nil {
			log.Println(err)
		}
		ctx.JSON(http.StatusForbidden, err.Error())
		return
	}

	// // creates new pair of jwt tokens
	tokens, err := s.Store.CreateTokens(ctx, user.UserID, time.Minute*15)
	if err != nil {
		fmt.Println("json err:", err)
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, tokens)
}
