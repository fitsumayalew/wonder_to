package session

import (
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"net/http"
	"time"
	"xCut/entity"
	"xCut/rtoken"
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



// Create creates and sets session cookie
func Create(claims jwt.Claims, sessionID string, signingKey []byte, w http.ResponseWriter) {

	signedString, err := rtoken.GenerateJwtToken(signingKey, claims)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
	}
	c := http.Cookie{
		Name:     sessionID,
		Value:    signedString,
		HttpOnly: true,
	}
	http.SetCookie(w, &c)
}

// Valid validates client cookie value
func Valid(cookieValue string, signingKey []byte) (bool, error) {
	valid := rtoken.IsCSRFValid(cookieValue, signingKey)
	if valid {
		return true, nil
	}
	return false, errors.New("Invalid Session Cookie")

}

// Remove expires existing session
func Remove(sessionID string, w http.ResponseWriter) {
	c := http.Cookie{
		Name:    sessionID,
		MaxAge:  -1,
		Expires: time.Unix(1, 0),
		Value:   "",
	}
	http.SetCookie(w, &c)
}
