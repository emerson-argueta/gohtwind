package auth

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"github.com/volatiletech/authboss/v3"
	"net/http"
	"os"
	"time"
)

type SessionState struct {
	authboss.ClientStateReadWriter
}

type ClientState struct {
	authboss.ClientState
	d map[string]string
}

func (c *ClientState) Get(key string) (string, bool) {
	v, ok := c.d[key]
	return v, ok
}

func (s SessionState) ReadState(request *http.Request) (authboss.ClientState, error) {
	cookie, err := request.Cookie("sessionId")
	if err != nil {
		return nil, err
	}
	sk := os.Getenv("GOHTWIND_SECRET_KEY")
	hsid := sha256.Sum256([]byte(fmt.Sprintf("%s%s", cookie.Value, sk)))
	sd, exp := s.getFromDatabase(hsid[:])
	if time.Now().After(exp) {
		return nil, fmt.Errorf("session expired")
	}
	res := &ClientState{d: sd}
	return res, nil
}

func (s SessionState) WriteState(writer http.ResponseWriter, state authboss.ClientState, events []authboss.ClientStateEvent) error {
	sid := make([]byte, 32)
	_, err := rand.Read(sid)
	if err != nil {
		return err
	}
	isSecure := os.Getenv("ENV") == "production"
	sdur := SessionDuration * time.Hour
	val := base64.StdEncoding.EncodeToString(sid)
	http.SetCookie(writer, &http.Cookie{
		Name:     "sessionId",
		Value:    val,
		HttpOnly: true,
		Secure:   isSecure,
		SameSite: http.SameSiteStrictMode,
		MaxAge:   int(sdur.Seconds()),
	})
	sk := os.Getenv("GOHTWIND_SECRET_KEY")
	hsid := sha256.Sum256([]byte(fmt.Sprintf("%s%s", val, sk)))
	exp := time.Now().Add(sdur)
	err = s.saveToDatabase(hsid[:], state, exp)
	if err != nil {
		return err
	}
	return nil
}

func (s SessionState) getFromDatabase(bytes []byte) (map[string]string, time.Time) {
	// TODO: get from database
	res := make(map[string]string)
	return res, time.Now()
}

func (s SessionState) saveToDatabase(bytes []byte, state authboss.ClientState, exp time.Time) error {
	// TODO: save to database
	return nil
}
