package keybindings

// //This defines keybindings for different input states;

// /* global keybindings */
// func globalKeys(g *gocui.Gui) {
// 	// [ctrl-C] for exit ; Keybinding
// 	// g.DeleteKeybindings("commandInput")
// 	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, Quit); err != nil {
// 		log.Panicln(err)
// 	}
// }

// // Keybindings when focusing on cmdbars
// func cmdinputKeys(g *gocui.Gui, prev string) {
// 	// [ctrl-C] for exit ; Keybinding
// 	removeKeyBinding(g, prev)
// 	if err := g.SetKeybinding("commandInput", gocui.KeyCtrlX, gocui.ModNone, QuitCommandInput); err != nil {
// 		log.Panicln(err)
// 	}

// 	//run cobra command
// 	if err := g.SetKeybinding("commandInput", gocui.KeyEnter, gocui.ModNone, SendCmd); err != nil {
// 		log.Panicln(err)
// 	}
// }

// // Keybindings for notekeys
// func noteKeys(g *gocui.Gui, prev string) {
// 	// [ctrl-x] for input bar setup;
// 	// g.DeleteKeybindings("")
// 	removeKeyBinding(g, prev)
// 	//pull up command bar
// 	if err := g.SetKeybinding("note", gocui.KeyCtrlX, gocui.ModNone, SetCommandInput); err != nil {
// 		log.Panicln(err)
// 	}

// 	if err := g.SetKeybinding("note", gocui.KeyCtrlA, gocui.ModNone, SetNoteHistory); err != nil {
// 		log.Panicln(err)
// 	}

// 	//send note to noteDB, display on new view;
// 	if err := g.SetKeybinding("note", gocui.KeyEnter, gocui.ModNone, SendNote); err != nil {
// 		log.Panicln(err)
// 	}

// }

// // Keybindings for note-history
// func noteHistoryKeys(g *gocui.Gui, prev string) {
// 	removeKeyBinding(g, prev)

// 	if err := g.SetKeybinding("noteHistory", gocui.KeyCtrlA, gocui.ModNone, QuitNoteHistory); err != nil {
// 		log.Panicln(err)
// 	}
// 	if err := g.SetKeybinding("noteHistory", 'h', gocui.ModNone, CursorLeft); err != nil {
// 		log.Panicln(err)
// 	}
// 	if err := g.SetKeybinding("noteHistory", 'l', gocui.ModNone, CursorRight); err != nil {
// 		log.Panicln(err)
// 	}
// 	if err := g.SetKeybinding("noteHistory", 'j', gocui.ModNone, CursorUp); err != nil {
// 		log.Panicln(err)
// 	}
// 	if err := g.SetKeybinding("noteHistory", 'k', gocui.ModNone, CursorDown); err != nil {
// 		log.Panicln(err)
// 	}

// }

//	func removeKeyBinding(g *gocui.Gui, prev string) {
//		if prev == "none" {
//			return
//		}
//		//else:
//		g.DeleteKeybindings(prev)
//	}

// // Keybinder for note
// type NoteKeyBinder struct{}

// // try to not use pointers when unnecessary
// var _ ui_models.KeyBinder = (NoteKeyBinder{})

// func (kb NoteKeyBinder) BindKeys(g *gocui.Gui, prev string) {
// 	//remove previous keybindings
// 	if prev == "None" {
// 		return
// 	}
// 	g.DeleteKeybindings(prev)

// 	//pull up command bar
// 	if err := g.SetKeybinding("note", gocui.KeyCtrlX, gocui.ModNone, SetCommandInput); err != nil {
// 		log.Panicln(err)
// 	}

// 	if err := g.SetKeybinding("note", gocui.KeyCtrlA, gocui.ModNone, SetNoteHistory); err != nil {
// 		log.Panicln(err)
// 	}

// 	//send note to noteDB, display on new view;
// 	if err := g.SetKeybinding("note", gocui.KeyEnter, gocui.ModNone, SendNote); err != nil {
// 		log.Panicln(err)
// 	}
// }

// Maybe just do ordinary handlers, and then do abstractions later;
