package viber

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"
)

// AppKey for you app provided by Viber
//const AppKey = "46160a5f87b294eb-9502de2bc1cf5ddb-5b70d84954155377"

// Sender structure
type Sender struct {
	Name   string `json:"name"`
	Avatar string `json:"avatar,omitempty"`
}

type event struct {
	Event        string    `json:"event"`
	Timestamp    Timestamp `json:"timestamp"`
	UserID       string    `json:"user_id"`
	MessageToken string    `json:"message_token"`
	Descr        string    `json:"descr"`
}

// Viber app
type Viber struct {
	AppKey string
	Sender Sender

	Subscribed          func(u User, msgToken string, t time.Time)
	ConversationStarted func()
	Message             func()
	Unsubscribed        func(userID, msgToken string, t time.Time)
	Delivered           func(userID, msgToken string, t time.Time)
	Seen                func(userID, msgToken string, t time.Time)
	Failed              func(userID, msgToken, descr string, t time.Time)
}

func (v *Viber) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return
	}
	r.Body.Close()

	if !v.checkHMAC(body, r.Header.Get("X-Viber-Content-Signature")) {
		return
	}

	var e event
	if err := json.Unmarshal(body, &e); err != nil {
		return
	}

	switch e.Event {
	case "subscribed":
		if v.Subscribed != nil {

		}

	case "unsubscribed":
		if v.Unsubscribed != nil {
			v.Unsubscribed(e.UserID, e.MessageToken, e.Timestamp.Time)
		}

	case "conversation_started":
		if v.ConversationStarted != nil {

		}

	case "delivered":
		if v.Delivered != nil {
			v.Delivered(e.UserID, e.MessageToken, e.Timestamp.Time)
		}

	case "seen":
		if v.Seen != nil {
			v.Seen(e.UserID, e.MessageToken, e.Timestamp.Time)
		}

	case "failed":
		if v.Failed != nil {
			v.Failed(e.UserID, e.MessageToken, e.Descr, e.Timestamp.Time)
		}

	case "message":
		if v.Message != nil {

		}

	}
}

// checkHMAC reports whether messageMAC is a valid HMAC tag for message.
func (v *Viber) checkHMAC(message []byte, messageMAC string) bool {
	hmac := hmac.New(sha256.New, []byte(v.AppKey))
	hmac.Write(message)
	return messageMAC == hex.EncodeToString(hmac.Sum(nil))
}
