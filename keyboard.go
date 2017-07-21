package viber

// Keyboard struct
type Keyboard struct {
	DefaultHeight bool     `json:"DefaultHeight"`
	BgColor       string   `json:"BgColor"`
	Buttons       []Button `json:"Buttons"`
}

// AddButton to keyboard
func (k *Keyboard) AddButton(b *Button) {
	k.Buttons = append(k.Buttons, *b)
}

// NewKeyboard struct with attribs init
func NewKeyboard(bgcolor string, defaultHeight bool) *Keyboard {
	return &Keyboard{
		DefaultHeight: defaultHeight,
		BgColor:       bgcolor,
	}
}
