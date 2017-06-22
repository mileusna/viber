package viber

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
	"time"
)

// Sender structure
type Sender struct {
	Name   string `json:"name"`
	Avatar string `json:"avatar,omitempty"`
}

type event struct {
	Event        string    `json:"event"`
	Timestamp    Timestamp `json:"timestamp"`
	MessageToken uint64    `json:"message_token,omitempty"`
	UserID       string    `json:"user_id,omitempty"`

	// failed event
	Descr string `json:"descr,omitempty"`

	//conversation_started event
	Type       string          `json:"type,omitempty"`
	Context    string          `json:"context,omitempty"`
	Subscribed bool            `json:"subscribed,omitempty"`
	User       json.RawMessage `json:"user,omitempty"`

	// message event
	Sender  json.RawMessage `json:"sender,omitempty"`
	Message json.RawMessage `json:"message,omitempty"`
}

// Viber app
type Viber struct {
	AppKey string
	Sender Sender

	// event methods
	Subscribed          func(u User, token uint64, t time.Time)
	ConversationStarted func(u User, conversationType, context string, subscribed bool, token uint64, t time.Time) Message
	Message             func(u User, m Message, token uint64, t time.Time)
	Unsubscribed        func(userID string, token uint64, t time.Time)
	Delivered           func(userID string, token uint64, t time.Time)
	Seen                func(userID string, token uint64, t time.Time)
	Failed              func(userID string, token uint64, descr string, t time.Time)
}

var regexpPeekMsgType = regexp.MustCompile("\"type\":\\s*\"(.*)\"")

// ServeHTTP
// https://developers.viber.com/docs/api/rest-bot-api/#callbacks
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
			go v.Unsubscribed(e.UserID, e.MessageToken, e.Timestamp.Time)
		}

	case "conversation_started":
		if v.ConversationStarted != nil {
			var u User
			if err := json.Unmarshal(e.User, &u); err != nil {
				return
			}
			if msg := v.ConversationStarted(u, e.Type, e.Context, e.Subscribed, e.MessageToken, e.Timestamp.Time); msg != nil {
				msg.SetReceiver("")
				msg.SetFrom("")
				b, _ := json.Marshal(msg)
				w.Write(b)
			}
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
			var u User
			if err := json.Unmarshal(e.Sender, &u); err != nil {
				return
			}

			msgType := peakMessageType(e.Message)
			switch msgType {
			case "text":
				var m TextMessage
				if err := json.Unmarshal(e.Message, &m); err != nil {
					return
				}
				go v.Message(u, &m, e.MessageToken, e.Timestamp.Time)

			case "picture":
				var m PictureMessage
				if err := json.Unmarshal(e.Message, &m); err != nil {
					return
				}
				go v.Message(u, &m, e.MessageToken, e.Timestamp.Time)

			case "video":
				var m VideoMessage
				if err := json.Unmarshal(e.Message, &m); err != nil {
					return
				}
				go v.Message(u, &m, e.MessageToken, e.Timestamp.Time)

			case "url":
				var m URLMessage
				if err := json.Unmarshal(e.Message, &m); err != nil {
					return
				}
				go v.Message(u, &m, e.MessageToken, e.Timestamp.Time)

			case "contact":
				// TODO
			case "location":
				// TODO
			default:
				return
			}
		}
	}
}

// checkHMAC reports whether messageMAC is a valid HMAC tag for message.
func (v *Viber) checkHMAC(message []byte, messageMAC string) bool {
	hmac := hmac.New(sha256.New, []byte(v.AppKey))
	hmac.Write(message)
	return messageMAC == hex.EncodeToString(hmac.Sum(nil))
}

// peakMessageType uses regexp to determin message type for unmarshaling
func peakMessageType(b []byte) string {
	matches := regexpPeekMsgType.FindAllSubmatch(b, -1)
	if len(matches) == 0 {
		return ""
	}

	return strings.ToLower(string(matches[0][1]))
}
