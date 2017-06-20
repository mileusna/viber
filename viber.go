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

// Sender structure
type Sender struct {
	Name   string `json:"name"`
	Avatar string `json:"avatar,omitempty"`
}

type Event struct {
	Event        string    `json:"event"`
	Timestamp    Timestamp `json:"timestamp"`
	UserID       string    `json:"user_id"`
	MessageToken uint64    `json:"message_token"`
	Descr        string    `json:"descr"`
}

// Viber app
type Viber struct {
	AppKey string
	Sender Sender

	Subscribed          func(u User, msgToken string, t time.Time)
	ConversationStarted func()
	Message             func()
	Unsubscribed        func(userID string, msgToken uint64, t time.Time)
	Delivered           func(userID string, msgToken uint64, t time.Time)
	Seen                func(userID string, msgToken uint64, t time.Time)
	Failed              func(userID string, msgToken uint64, descr string, t time.Time)
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

	var e Event
	if err := json.Unmarshal(body, &e); err != nil {
		//return
	}

	switch e.Event {
	case "subscribed":
		if v.Subscribed != nil {

		}

	case "unsubscribed":
		if v.Unsubscribed != nil {
			go v.Unsubscribed(e.UserID, e.MessageToken, e.Timestamp.Time)
		}

	case "conversation_started":
		if v.ConversationStarted != nil {

		}

	case "delivered":
		if v.Delivered != nil {
			go v.Delivered(e.UserID, e.MessageToken, e.Timestamp.Time)
		}

	case "seen":
		if v.Seen != nil {
			go v.Seen(e.UserID, e.MessageToken, e.Timestamp.Time)
		}

	case "failed":
		if v.Failed != nil {
			go v.Failed(e.UserID, e.MessageToken, e.Descr, e.Timestamp.Time)
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
