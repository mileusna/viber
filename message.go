package viber

import (
	"encoding/json"
	"fmt"
)

/*
{
	"receiver": "01234567890A=",
	"min_api_version": 1,
	"sender": {
		"name": "John McClane",
		"avatar": "http://avatar.example.com"
	},
	"tracking_data": "tracking data",
	"type": "text",
	"text": "a message from pa"
}
*/

// MessageType for viber messaging
type MessageType string

type messageResponse struct {
	Status        int    `json:"status"`
	StatusMessage string `json:"status_message"`
	MessageToken  uint64 `json:"message_token"`
}

// Message interface for all types of viber messages
type Message interface {
	SetReceiver(r string)
	SetFrom(from string)
}

// TextMessage for Viber
type TextMessage struct {
	Receiver      string      `json:"receiver,omitempty"`
	From          string      `json:"from,omitempty"`
	MinAPIVersion uint        `json:"min_api_version,omitempty"`
	Sender        Sender      `json:"sender"`
	Type          MessageType `json:"type"`
	TrackingData  string      `json:"tracking_data,omitempty"`
	Text          string      `json:"text"`
	//    "media": "http://www.images.com/img.jpg",
	//    "thumbnail": "http://www.images.com/thumb.jpg"
	// 	"size": 10000,
	// 	"duration": 10
}

// URLMessage structure
type URLMessage struct {
	TextMessage
	Media string `json:"media"`
}

// PictureMessage structure
type PictureMessage struct {
	TextMessage
	Media     string `json:"media"`
	Thumbnail string `json:"thumbnail,omitempty"`
}

// VideoMessage structure
type VideoMessage struct {
	TextMessage
	Media     string `json:"media"`
	Thumbnail string `json:"thumbnail,omitempty"`
	Size      uint   `json:"size"`
	Duration  uint   `json:"duration,omitempty"`
}

// Button for carousel
type Button struct {
	Columns    int    `json:"Columns"`
	Rows       int    `json:"Rows"`
	ActionType string `json:"ActionType"`
	ActionBody string `json:"ActionBody"`
	Image      string `json:"Image,omitempty"`
	Text       string `json:"Text,omitempty"`
	TextSize   string `json:"TextSize,omitempty"`
	TextVAlign string `json:"TextVAlign,omitempty"`
	TextHAlign string `json:"TextHAlign,omitempty"`
}

// RichMedia for carousel
type RichMedia struct {
	Type                string   `json:"Type"`
	ButtonsGroupColumns int      `json:"ButtonsGroupColumns"`
	ButtonsGroupRows    int      `json:"ButtonsGroupRows"`
	BgColor             string   `json:"BgColor"`
	TrackingData        string   `json:"tracking_data,omitempty"`
	Buttons             []Button `json:"Buttons"`
}

// RichMediaMessage / Carousel
type RichMediaMessage struct {
	AuthToken     string      `json:"auth_token"`
	Receiver      string      `json:"receiver"`
	Type          MessageType `json:"type"`
	MinAPIVersion int         `json:"min_api_version"`
	RichMedia     RichMedia   `json:"rich_media"`
}

const (
	TypeTextMessage      = MessageType("text")
	TypeURLMessage       = MessageType("url")
	TypePictureMessage   = MessageType("picture")
	TypeVideoMessage     = MessageType("video")
	TypeFileMessage      = MessageType("file")
	TypeLocationMessage  = MessageType("location")
	TypeContactMessage   = MessageType("contact")
	TypeStickerMessage   = MessageType("sticker")
	TypeRichMediaMessage = MessageType("rich_media")
)

//video, file, location, contact, sticker, carousel content

func parseMsgResponse(b []byte) (msgToken uint64, err error) {
	var resp messageResponse
	if err := json.Unmarshal(b, &resp); err != nil {
		return 0, err
	}

	if resp.Status != 0 {
		return resp.MessageToken, Error{Status: resp.Status, StatusMessage: resp.StatusMessage}
	}

	return resp.MessageToken, nil
}

func (v *Viber) sendMessage(url string, m interface{}) (msgToken uint64, err error) {
	b, err := v.PostData(url, m)
	if err != nil {
		return 0, err
	}
	return parseMsgResponse(b)
}

// NewTextMessage viber
func (v *Viber) NewTextMessage(msg string) *TextMessage {
	return &TextMessage{
		Sender: v.Sender,
		Type:   TypeTextMessage,
		Text:   msg,
	}
}

// SendTextMessage to reciever, returns message token
func (v *Viber) SendTextMessage(receiver string, msg string) (msgToken uint64, err error) {
	m := v.NewTextMessage(msg)
	m.Receiver = receiver
	return v.sendMessage("https://chatapi.viber.com/pa/send_message", m)
}

// NewURLMessage for viber
func (v *Viber) NewURLMessage(msg string, url string) *URLMessage {
	return &URLMessage{
		TextMessage: TextMessage{
			Sender: v.Sender,
			Type:   TypeURLMessage,
			Text:   msg,
		},
		Media: url,
	}
}

// SendURLMessage to receiver, return message token
func (v *Viber) SendURLMessage(receiver string, s Sender, msg string, url string) (msgToken uint64, err error) {
	m := v.NewURLMessage(msg, url)
	return v.sendMessage("https://chatapi.viber.com/pa/send_message", m)
}

// NewPictureMessage for viber
func (v *Viber) NewPictureMessage(msg string, url string, thumbURL string) *PictureMessage {
	return &PictureMessage{
		TextMessage: TextMessage{
			Sender: v.Sender,
			Type:   TypePictureMessage,
			Text:   msg,
		},
		Media:     url,
		Thumbnail: thumbURL,
	}
}

// SendPictureMessage to receiver, returns message token
// func (v *Viber) SendPictureMessage(receiver string, s Sender, msg string, url string, thumbURL string) (token int64, err error) {
// 	m := PictureMessage{
// 		Message: Message{
// 			Receiver: receiver,
// 			Sender:   s,
// 			Type:     TypePictureMessage,
// 			Text:     msg,
// 		},
// 		Media:     url,
// 		Thumbnail: thumbURL,
// 	}

// 	return v.sendMessage("https://chatapi.viber.com/pa/send_message", m)
// }

func (v *Viber) SendCarousel(receiver string) {
	r := RichMediaMessage{
		MinAPIVersion: 2,
		Receiver:      receiver,
		AuthToken:     v.AppKey,
		Type:          TypeRichMediaMessage,
		RichMedia: RichMedia{
			Type:                "rich_media",
			ButtonsGroupColumns: 6,
			ButtonsGroupRows:    6,
			BgColor:             "#FFFFFF",
		},
	}

	b1 := Button{
		Columns:    6,
		Rows:       3,
		ActionType: "open-url",
		ActionBody: "https://aviokarte.rs/",
		Image:      "http://nstatic.net/beta/2b5b3ff1972f61d9bcfaaddd061aa1b9.jpg",
	}
	r.RichMedia.Buttons = append(r.RichMedia.Buttons, b1)

	b1 = Button{
		Columns:    6,
		Rows:       2,
		ActionType: "open-url",
		ActionBody: "https://aviokarte.rs/",
		Text:       "Košarkaši Crvene zvezde odbranili titulu šampiona Srbije",
	}
	r.RichMedia.Buttons = append(r.RichMedia.Buttons, b1)

	b1 = Button{
		Columns:    3,
		Rows:       1,
		ActionType: "reply",
		ActionBody: "ID: 21432323",
		Text:       "<font color=#ffffff>Otvori</font>",
		TextSize:   "large",
		TextVAlign: "middle",
		TextHAlign: "middle",
		Image:      "https://s14.postimg.org/4mmt4rw1t/Button.png",
	}
	r.RichMedia.Buttons = append(r.RichMedia.Buttons, b1)

	// b2 := Button{

	// 	Columns:    6,
	// 	Rows:       6,
	// 	ActionType: "reply",
	// 	ActionBody: "https://aviokarte.rs/",
	// 	Image:      "https://aviokarte.rs/images/logo.png",
	// 	Text:       "Drugi tekst",
	// 	TextSize:   "large",
	// 	TextVAlign: "middle",
	// 	TextHAlign: "left",
	// }
	// r.RichMedia.Buttons = append(r.RichMedia.Buttons, b2)

	resp, err := v.PostData("https://chatapi.viber.com/pa/send_message", r)
	fmt.Println(string(resp), err)

}

// SendPublicMessage from public account
func (v *Viber) SendPublicMessage(from string, m Message) (msgToken uint64, err error) {
	m.SetFrom(from)
	return v.sendMessage("https://chatapi.viber.com/pa/post", m)
}

// SendMessage to receiver
func (v *Viber) SendMessage(to string, m Message) (msgToken uint64, err error) {
	m.SetReceiver(to)
	return v.sendMessage("https://chatapi.viber.com/pa/send_message", m)
}

// SetReceiver for text message
func (m *TextMessage) SetReceiver(r string) {
	m.Receiver = r
}

// SetFrom to text message for public account message
func (m *TextMessage) SetFrom(from string) {
	m.From = from
}
