package viber

type InputFieldStateType string

// InputFieldState Values
const (
	RegularState   = InputFieldStateType("regular")
	MinimizedState = InputFieldStateType("minimized")
	HiddenState    = InputFieldStateType("hidden")
)

// Keyboard struct
type Keyboard struct {
	Type            string              `json:"Type"`
	DefaultHeight   bool                `json:"DefaultHeight,omitempty"`
	BgColor         string              `json:"BgColor,omitempty"`
	Buttons         []Button            `json:"Buttons"`
	InputFieldState InputFieldStateType `json:"InputFieldState,omitempty"`
}

// AddButton to keyboard
func (k *Keyboard) AddButton(b *Button) {
	k.Buttons = append(k.Buttons, *b)
}

func (k *Keyboard) SetInputFieldState(state InputFieldStateType) *Keyboard {
	k.InputFieldState = state
	return k
}

// NewKeyboard struct with attribs init
func (v *Viber) NewKeyboard(bgcolor string, defaultHeight bool) *Keyboard {
	return &Keyboard{
		Type:          "keyboard",
		DefaultHeight: defaultHeight,
		BgColor:       bgcolor,
	}
}
