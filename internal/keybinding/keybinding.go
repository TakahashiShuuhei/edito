package keybinding

import (
	"github.com/nsf/termbox-go"
)

type KeyHandler func()

type KeyBinding struct {
	Key      termbox.Key
	Ch       rune
	Modifier termbox.Modifier
	Handler  KeyHandler
}

type KeyMap struct {
	bindings []KeyBinding
}

func NewKeyMap() *KeyMap {
	return &KeyMap{
		bindings: make([]KeyBinding, 0),
	}
}

func (km *KeyMap) Bind(key termbox.Key, ch rune, modifier termbox.Modifier, handler KeyHandler) {
	binding := KeyBinding{
		Key:      key,
		Ch:       ch,
		Modifier: modifier,
		Handler:  handler,
	}
	km.bindings = append(km.bindings, binding)
}

func (km *KeyMap) BindKey(key termbox.Key, handler KeyHandler) {
	km.Bind(key, 0, 0, handler)
}

func (km *KeyMap) BindChar(ch rune, handler KeyHandler) {
	km.Bind(0, ch, 0, handler)
}

func (km *KeyMap) BindCtrl(ch rune, handler KeyHandler) {
	km.Bind(0, ch, termbox.ModAlt, handler)
}

func (km *KeyMap) Handle(ev termbox.Event) bool {
	for _, binding := range km.bindings {
		if binding.matches(ev) {
			binding.Handler()
			return true
		}
	}
	return false
}

func (b *KeyBinding) matches(ev termbox.Event) bool {
	if b.Key != 0 && ev.Key != b.Key {
		return false
	}
	if b.Ch != 0 && ev.Ch != b.Ch {
		return false
	}
	if b.Modifier != 0 && ev.Mod != b.Modifier {
		return false
	}
	return true
}

func CreateEmacsKeyMap() *KeyMap {
	return NewKeyMap()
}