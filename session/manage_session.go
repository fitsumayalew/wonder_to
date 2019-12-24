package session

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"net/http"
	"time"
	"xcut/entity"
	"xcut/rtoken"
)

const SessionKey = "session_key"

//create new signing key and sessionId
func CreateNewSession(id uint) *entity.Session {
	tokenExpires := time.Now().AddDate(0, 1, 0).Unix()
	signingString, err := rtoken.GenerateRandomString(32)
	sessionId, err := rtoken.GenerateRandomString(32)
	if err != nil {
		panic(err)
	}
	return &entity.Session{
		SessionId:  sessionId,
		Expires:    tokenExpires,
		SigningKey: []byte(signingString),
		UUID:       id,
	}
}

// Set session cookie
func SetCookie(claims jwt.Claims, expire int64, signingKey []byte, w http.ResponseWriter) {
	signedString, err := rtoken.GenerateJwtToken(signingKey, claims)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
	}
	cookie := http.Cookie{
		Name:     SessionKey,
		Value:    signedString,
		HttpOnly: true,
		Expires:  time.Unix(expire, 0),
	}

	http.SetCookie(w, &cookie)
}


// expire existing session
func RemoveCookie(w http.ResponseWriter) {
	cookie := http.Cookie{
		Name:    SessionKey,
		Value:   "",
		Expires: time.Unix(1, 0),
		MaxAge:  -1,
	}
	http.SetCookie(w, &cookie)
}