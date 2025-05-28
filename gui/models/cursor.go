package models

type Direction int

const (
	Up = iota
	Down
	Left
	Right
)

type Cursor struct {
	OnDisplay bool
	Move      func(d Direction) error
}
