package db

// transfer input data to database;
func AddNote(content string) error {
	//init note struct
	note := &Note{Content: content}
	//pass the string to database;
	result := NoteDB.Create(note)
	return result.Error
}

// // Display all notes added;
// func Display() error {

// }
