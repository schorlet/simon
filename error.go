package main

import (
	"errors"
	"fmt"
)

var (
	// ErrGameOver is the error returned when the game is over.
	ErrGameOver = errors.New("game over")

	// ErrPlayerTurn is the error returned by Display when you should invoke Play.
	ErrPlayerTurn = errors.New("player turn")

	// ErrGameTurn is the error returned by Play when you should invoke Display.
	ErrGameTurn = errors.New("game turn")
)

// ErrBadColor is the error returned by Play when you played a bad color.
type ErrBadColor struct {
	Got  Color
	Want Color
}

func (e ErrBadColor) Error() string {
	return fmt.Sprintf("bad color:%v, want:%v", e.Got, e.Want)
}
