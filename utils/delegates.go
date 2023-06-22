package utils

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
)

type SliceKeyMap struct {
	KeyBindings []key.Binding
}

func NewDefaultListDelegate(keys *SliceKeyMap) list.DefaultDelegate {
	d := list.NewDefaultDelegate()

	d.ShortHelpFunc = func() []key.Binding {
		return keys.KeyBindings
	}
	d.FullHelpFunc = func() [][]key.Binding {
		return [][]key.Binding{keys.KeyBindings}
	}

	return d
}
