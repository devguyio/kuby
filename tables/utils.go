package tables

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/table"
	"github.com/samber/lo"
)

func LongestInColumn(rows *[]table.Row, columnIndex int) int {
	columnSlice := lo.Map(*rows, func(item table.Row, _ int) string {
		return item[columnIndex]
	})

	longestColumnValue := lo.MaxBy(columnSlice, func(item string, max string) bool {
		return len(item) > len(max)
	})

	return len(longestColumnValue)
}

// TableKeyMap defines a set of keybindings. To work for help it must satisfy
// key.Map. It could also very easily be a map[string]key.Binding.
type TableKeyMap struct {
	table.KeyMap
	Choose key.Binding
	Help   key.Binding
	Quit   key.Binding
}

// ShortHelp returns keybindings to be shown in the mini help view. It's part
// of the key.Map interface.
func (k TableKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.LineUp, k.LineDown, k.Choose, k.Help, k.Quit}
}

// FullHelp returns keybindings for the expanded help view. It's part of the
// key.Map interface.
func (k TableKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.LineUp, k.LineDown, k.Choose, k.Help, k.Quit}, // first column
		{k.PageUp, k.PageDown, k.GotoTop, k.GotoBottom},  // second column
	}
}

var DefaultTableKeyMap = TableKeyMap{KeyMap: table.DefaultKeyMap(), Choose: key.NewBinding(
	key.WithKeys("enter"),
	key.WithHelp("enter", "choose"),
), Help: key.NewBinding(
	key.WithKeys("?"),
	key.WithHelp("?", "toggle help"),
), Quit: key.NewBinding(
	key.WithKeys("q", "esc"),
	key.WithHelp("esc, q", "quit"),
),
}
