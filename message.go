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
	Columns    int        `json:"Columns"`
	Rows       int        `json:"Rows"`
	ActionType ActionType `json:"ActionType"`
	ActionBody string     `json:"ActionBody"`
	Image      string     `json:"Image,omitempty"`
	Text       string     `json:"Text,omitempty"`
	TextSize   TextSize   `json:"TextSize,omitempty"`
	TextVAlign TextVAlign `json:"TextVAlign,omitempty"`
	TextHAlign TextHAlign `json:"TextHAlign,omitempty"`
}

// RichMedia for carousel
type RichMedia struct {
	Type                MessageType `json:"Type"`
	ButtonsGroupColumns int         `json:"ButtonsGroupColumns"`
	ButtonsGroupRows    int         `json:"ButtonsGroupRows"`
	BgColor             string      `json:"BgColor"`
	TrackingData        string      `json:"tracking_data,omitempty"`
	Buttons             []Button    `json:"Buttons"`
}

// RichMediaMessage / Carousel
type RichMediaMessage struct {
	AuthToken     string      `json:"auth_token"`
	Receiver      string      `json:"receiver,omitempty"`
	Type          MessageType `json:"type"`
	MinAPIVersion int         `json:"min_api_version"`
	RichMedia     RichMedia   `json:"rich_media"`
}

// MessageType for viber messaging
type MessageType string

// Message types
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

// NewURLMessage creates new message with global sender and common params set
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

// SendTextMessage to reciever, returns message token
func (v *Viber) SendTextMessage(receiver string, msg string) (msgToken uint64, err error) {
	return v.SendMessage(receiver, v.NewTextMessage(msg))
}

// SendURLMessage to easily send url messages as global sender
func (v *Viber) SendURLMessage(receiver string, msg string, url string) (msgToken uint64, err error) {
	return v.SendMessage(receiver, v.NewURLMessage(msg, url))
}

// SendPictureMessage to receiver, returns message token
func (v *Viber) SendPictureMessage(receiver string, msg string, url string, thumbURL string) (token uint64, err error) {
	return v.SendMessage(receiver, v.NewPictureMessage(msg, url, thumbURL))
}

// NewRichMediaMessage creates new empty carousel message
func (v *Viber) NewRichMediaMessage(cols, rows int, bgColor string) *RichMediaMessage {
	return &RichMediaMessage{
		MinAPIVersion: 2,
		AuthToken:     v.AppKey,
		Type:          TypeRichMediaMessage,
		RichMedia: RichMedia{
			Type:                TypeRichMediaMessage,
			ButtonsGroupColumns: cols,
			ButtonsGroupRows:    rows,
			BgColor:             bgColor,
		},
	}
}

// NewButton helper function for creating button with text and image
func (v *Viber) NewButton(cols, rows int, typ ActionType, actionBody string, text, image string) Button {
	return Button{
		Columns:    cols,
		Rows:       rows,
		ActionType: typ,
		ActionBody: actionBody,
		Text:       text,
		Image:      image,
	}
}

// NewImageButton helper function for creating image button struct with common params
func (v *Viber) NewImageButton(cols, rows int, typ ActionType, actionBody string, image string) Button {
	return Button{
		Columns:    cols,
		Rows:       rows,
		ActionType: typ,
		ActionBody: actionBody,
		Image:      image,
	}
}

// NewTextButton helper function for creating image button struct with common params
func (v *Viber) NewTextButton(cols, rows int, t ActionType, actionBody, text string) Button {
	return Button{
		Columns:    cols,
		Rows:       rows,
		ActionType: t,
		ActionBody: actionBody,
		Text:       text,
	}
}

// AddButton to rich media message
func (rm *RichMediaMessage) AddButton(b Button) {
	rm.RichMedia.Buttons = append(rm.RichMedia.Buttons, b)
}

func (v *Viber) sendCarousel(receiver string) {
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
		ActionType: OpenURL,
		ActionBody: "https://aviokarte.rs/",
		Image:      "http://nstatic.net/beta/2b5b3ff1972f61d9bcfaaddd061aa1b9.jpg",
	}
	r.RichMedia.Buttons = append(r.RichMedia.Buttons, b1)

	b1 = Button{
		Columns:    6,
		Rows:       2,
		ActionType: OpenURL,
		ActionBody: "https://aviokarte.rs/",
		Text:       "Košarkaši Crvene zvezde odbranili titulu šampiona Srbije",
	}
	r.RichMedia.Buttons = append(r.RichMedia.Buttons, b1)

	b1 = Button{
		Columns:    3,
		Rows:       1,
		ActionType: Reply,
		ActionBody: "ID: 21432323",
		Text:       "<font color=#ffffff>Otvori</font>",
		TextSize:   Large,
		TextVAlign: Middle,
		TextHAlign: Center,
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
	// text, picture, video, file, location, contact, sticker and url
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

// SetReceiver for text message
func (rm *RichMediaMessage) SetReceiver(r string) {
	rm.Receiver = r
}

// SetFrom to satisfy interface although RichMedia messages can't be sent to publich chat and don't have From
func (rm *RichMediaMessage) SetFrom(from string) {}

// TextSize for carousel buttons
// viber.Small
// viber.Medium
// viber.Large
type TextSize string

// TextSize values
const (
	Small  = TextSize("small")
	Medium = TextSize("medium")
	Large  = TextSize("large")
)

// ActionType for carousel buttons
// viber.Reply
// viber.OpenURL
type ActionType string

// ActionType values
const (
	Reply   = ActionType("reply")
	OpenURL = ActionType("open-url")
)

// TextVAlign for carousel buttons
// viber.Top
// viber.Middle
// viber.Bottom
type TextVAlign string

// TextVAlign values
const (
	Top    = TextVAlign("top")
	Middle = TextVAlign("middle")
	Bottom = TextVAlign("bottom")
)

// TextHAlign for carousel buttons
// viber.Left
// viber.Center
// viber.Middle
type TextHAlign string

// TextHAlign values
const (
	Left   = TextHAlign("left")
	Center = TextHAlign("middle")
	Right  = TextHAlign("right")
)
