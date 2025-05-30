package gui

import "github.com/haochend413/mantis/gui/keybindings"

// This function inits all the keybindings
func (gui *Gui) InitKeyBindings() error {
	for _, k := range keybindings.CreateAllKeybinders() {
		//use parsor function to get key
		s := keybindings.Parsor(k.Key)
		if s.Valid {
			// kCopy := k
			if s.IsRune {
				if err := gui.g.SetKeybinding(k.ViewName, s.Rune, k.Modifier, k.Handler); err != nil {
					return err
				}

			} else {
				if err := gui.g.SetKeybinding(k.ViewName, s.Key, k.Modifier, k.Handler); err != nil {
					return err
				}

			}
		}
	}
	return nil
}
