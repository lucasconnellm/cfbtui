package keybindings

import "github.com/charmbracelet/bubbles/key"

const spacebar = " "

// KeyMap defines the keybindings for the viewport. Note that you don't
// necessary need to use keybindings at all; the viewport can be controlled
// programmatically with methods like Model.LineDown(1). See the GoDocs for
// details.
type KeyMap struct {
	PageDown     key.Binding
	PageUp       key.Binding
	HalfPageUp   key.Binding
	HalfPageDown key.Binding
	Down         key.Binding
	Up           key.Binding
	Back         key.Binding
	Quit         key.Binding
	Select       key.Binding
}

// DefaultKeyMap returns a set of pager-like default keybindings.
func DefaultKeyMap() KeyMap {
	return KeyMap{
		PageDown: key.NewBinding(
			key.WithKeys("pgdown", spacebar, "f"),
			key.WithHelp("f/pgdn", "page down"),
		),
		PageUp: key.NewBinding(
			key.WithKeys("pgup", "b"),
			key.WithHelp("b/pgup", "page up"),
		),
		HalfPageUp: key.NewBinding(
			key.WithKeys("u", "ctrl+u"),
			key.WithHelp("u", "½ page up"),
		),
		HalfPageDown: key.NewBinding(
			key.WithKeys("d", "ctrl+d"),
			key.WithHelp("d", "½ page down"),
		),
		Up: key.NewBinding(
			key.WithKeys("up", "k"),
			key.WithHelp("↑/k", "up"),
		),
		Down: key.NewBinding(
			key.WithKeys("down", "j"),
			key.WithHelp("↓/j", "down"),
		),
		Back: key.NewBinding(
			key.WithKeys("backspace", "left"),
			key.WithHelp("backspace", "Go back"),
		),
		Quit: key.NewBinding(
			key.WithKeys("q", "ctrl+c", "esc"),
			key.WithHelp("q", "Exit program"),
		),
		Select: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "select obj"),
		),
	}
}
