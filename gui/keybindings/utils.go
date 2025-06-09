package keybindings

import (
	"github.com/jroimartin/gocui"
)

// Keybinding Handlers
type KeyType struct {
	Key    gocui.Key
	Rune   rune
	IsRune bool
	Valid  bool
}

// This function deals with parsing the keystring into the real keyvalue for gocui setkeybinder function
func Parsor(key string) KeyType {
	//non-rune cases;
	switch key {
	case "enter":
		return KeyType{Key: gocui.KeyEnter, Valid: true}
	case "ct-x":
		return KeyType{Key: gocui.KeyCtrlX, Valid: true}
	case "ct-a":
		return KeyType{Key: gocui.KeyCtrlA, Valid: true}
	case "ct-c":
		return KeyType{Key: gocui.KeyCtrlC, Valid: true}
	case "ct-e":
		return KeyType{Key: gocui.KeyCtrlE, Valid: true}
	case "ct-s":
		return KeyType{Key: gocui.KeyCtrlS, Valid: true}
	case "up":
		return KeyType{Key: gocui.KeyArrowUp, Valid: true}
	case "down":
		return KeyType{Key: gocui.KeyArrowDown, Valid: true}
	case "left":
		return KeyType{Key: gocui.KeyArrowLeft, Valid: true}
	case "right":
		return KeyType{Key: gocui.KeyArrowRight, Valid: true}
	case "ct-space":
		return KeyType{Key: gocui.KeyCtrlSpace, Valid: true}
	case "tab":
		return KeyType{Key: gocui.KeyTab, Valid: true}
	case "ct-q":
		return KeyType{Key: gocui.KeyCtrlQ, Valid: true}

	}

	//rune case: length of string is 1
	if len(key) == 1 {
		return KeyType{Rune: rune(key[0]), IsRune: true, Valid: true}
	}

	//default: valid = false (0)
	return KeyType{}

}
