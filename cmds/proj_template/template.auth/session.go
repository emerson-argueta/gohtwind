package auth

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/volatiletech/authboss/v3"
	"net/http"
	"os"
	"{{PROJECT_NAME}}/infra"
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
	cookie, err := request.Cookie(authboss.SessionKey)
	if err != nil && err == http.ErrNoCookie {
		return &ClientState{d: make(map[string]string)}, nil
	}
	if err != nil {
		return nil, err
	}
	sk := os.Getenv("GOHTWIND_SECRET_KEY")
	hsid := sha256.Sum256([]byte(fmt.Sprintf("%s%s", cookie.Value, sk)))
	cs := s.getFromDatabase(hsid[:])
	return cs, nil
}

func (s SessionState) WriteState(writer http.ResponseWriter, state authboss.ClientState, events []authboss.ClientStateEvent) error {
	sid := make([]byte, 32)
	_, err := rand.Read(sid)
	if err != nil {
		return err
	}
	isSecure := os.Getenv("ENV") == "production"
	val := base64.StdEncoding.EncodeToString(sid)
	http.SetCookie(writer, &http.Cookie{
		Name:     authboss.SessionKey,
		Value:    val,
		HttpOnly: true,
		Secure:   isSecure,
		SameSite: http.SameSiteStrictMode,
	})
	sk := os.Getenv("GOHTWIND_SECRET_KEY")
	hsid := sha256.Sum256([]byte(fmt.Sprintf("%s%s", val, sk)))
	nsm := state.(*ClientState).d
	for _, event := range events {
		nsm[event.Key] = event.Value
	}
	err = s.saveToDatabase(hsid[:], nsm)
	if err != nil {
		return err
	}
	return nil
}

func (s SessionState) getFromDatabase(hashedSessionId []byte) *ClientState {
	// TODO: get from database
	d := make(map[string]string)
	res := &ClientState{d: d}
	return res
}

func (s SessionState) saveToDatabase(hashedSessionId []byte, newStateMap map[string]string) error {
	fmt.Println(newStateMap)
	//then to json
	j, err := json.Marshal(newStateMap)
	if err != nil {
		return err
	}
	//then the json should be encrypted, then the encrypted json should be saved to database
	ej, err := infra.Encrypt(j)
	if err != nil {
		return err
	}
	fmt.Println(ej)

	return nil
}
