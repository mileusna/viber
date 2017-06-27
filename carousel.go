package viber

// RichMediaMessage / Carousel
type RichMediaMessage struct {
	AuthToken     string      `json:"auth_token"`
	Receiver      string      `json:"receiver,omitempty"`
	Type          MessageType `json:"type"`
	MinAPIVersion int         `json:"min_api_version"`
	RichMedia     RichMedia   `json:"rich_media"`
	AltText       string      `json:"alt_text,omitempty"`
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

// AddButton to rich media message
func (rm *RichMediaMessage) AddButton(b *Button) {
	rm.RichMedia.Buttons = append(rm.RichMedia.Buttons, *b)
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
func (v *Viber) NewButton(cols, rows int, typ ActionType, actionBody string, text, image string) *Button {
	return &Button{
		Columns:    cols,
		Rows:       rows,
		ActionType: typ,
		ActionBody: actionBody,
		Text:       text,
		Image:      image,
	}
}

// NewImageButton helper function for creating image button struct with common params
func (v *Viber) NewImageButton(cols, rows int, typ ActionType, actionBody string, image string) *Button {
	return &Button{
		Columns:    cols,
		Rows:       rows,
		ActionType: typ,
		ActionBody: actionBody,
		Image:      image,
	}
}

// NewTextButton helper function for creating image button struct with common params
func (v *Viber) NewTextButton(cols, rows int, t ActionType, actionBody, text string) *Button {
	return &Button{
		Columns:    cols,
		Rows:       rows,
		ActionType: t,
		ActionBody: actionBody,
		Text:       text,
	}
}

// SetReceiver for text message
func (rm *RichMediaMessage) SetReceiver(r string) {
	rm.Receiver = r
}

// SetFrom to satisfy interface although RichMedia messages can't be sent to publich chat and don't have From
func (rm *RichMediaMessage) SetFrom(from string) {}

func (b *Button) TextSizeSmall() *Button {
	b.TextSize = Small
	return b
}

func (b *Button) TextSizeMedium() *Button {
	b.TextSize = Medium
	return b
}

func (b *Button) TextSizeLarge() *Button {
	b.TextSize = Large
	return b
}

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

func (b *Button) TextVAlignTop() *Button {
	b.TextVAlign = Top
	return b
}

func (b *Button) TextVAlignMiddle() *Button {
	b.TextVAlign = Middle
	return b
}

func (b *Button) TextVAlignBottom() *Button {
	b.TextVAlign = Bottom
	return b
}

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

func (b *Button) TextHAlignLeft() *Button {
	b.TextHAlign = Left
	return b
}

func (b *Button) TextHAlignMiddle() *Button {
	b.TextHAlign = Center
	return b
}

func (b *Button) TextHAlignRight() *Button {
	b.TextHAlign = Right
	return b
}
