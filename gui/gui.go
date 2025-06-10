package gui

import (
	"log"

	"github.com/awesome-gocui/gocui"
	dbcontroller "github.com/haochend413/mantis/controllers/db_controller"
	"github.com/haochend413/mantis/models"
)

// main Gui struct
type Gui struct {
	g                *gocui.Gui
	windows          []*models.Window
	first_init_check bool
	DBManager        *dbcontroller.DBManager
}

// need to use a map to hande quick window search

// This function inits a new Gui object;
func (gui *Gui) GuiInit() {
	// setup the new gui instance
	g, err := gocui.NewGui(gocui.Output256, false)
	if err != nil {
		//check startup err
		log.Panicln(err)
	}
	gui.g = g
	defer gui.g.Close()

	//
	//set configs for layout functions
	gui.first_init_check = true
	// Set layout manager function (called every frame to layout views) (set windows)
	gui.g.SetManagerFunc(gui.layout)
	//set up db manager
	gui.DBManager = &dbcontroller.DBManager{}
	gui.DBManager.InitManager()
	//fetch data into DB_Data stored;
	DB_Data = gui.DBManager.FetchAll()
	//defer close;
	// defer gui.DBManager.CloseManager()

	//init keybindings
	if err := gui.InitKeyBindings(); err != nil {
		log.Panicln(err)
	}

	if err := gui.g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}

func (gui *Gui) GuiClose() {
	gui.DBManager.RefreshAll(DB_Data)
}
