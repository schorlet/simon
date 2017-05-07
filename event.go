package main

import (
	"image"
)

type (
	// PlayerEvent is a player click event.
	PlayerEvent struct {
		// Index is the index of the target button.
		Index int
		// Up is true on mouse release.
		Up bool
	}

	// GameEvent is a game click event.
	GameEvent struct {
		// Index is the index of the target button.
		Index int
		// Up is true on mouse release.
		Up bool
	}

	// LostEvent is a lost click event.
	LostEvent struct {
		// Index is the index of the target button.
		Index int
		// Up is true on mouse release.
		Up bool
		// Count is a state counter.
		Count int
	}
)

// PlayerEvent creates a PlayerEvent from src point in the given picture.
// PlayerEvent returns false if src is beyond the picture bounds.
func (p *Picture) PlayerEvent(src image.Point) (PlayerEvent, bool) {
	color := p.At(src.X, src.Y)

	for i, sc := range upColors {
		if color == sc {
			return PlayerEvent{Index: i}, true
		}
	}
	for i, sc := range downColors {
		if color == sc {
			return PlayerEvent{Index: i, Up: true}, true
		}
	}
	return PlayerEvent{}, false
}

// Center returns true if the event targets the control button.
func (e PlayerEvent) Center() bool {
	return e.Up && e.Index == 4
}

// Down returns true on mouse press.
func (e PlayerEvent) Down() bool {
	return !e.Up
}

// Color gives the button color.
func (e PlayerEvent) Color() Color {
	return Color(e.Index + 1)
}

// NewGameEvent creates a GameEvent for the specified color.
func NewGameEvent(c Color) GameEvent {
	return GameEvent{Index: int(c) - 1}
}

// Down returns true on mouse press.
func (e GameEvent) Down() bool {
	return !e.Up
}

// NewLostEvent creates a LostEvent for the specified color.
func NewLostEvent(c Color) LostEvent {
	return LostEvent{Index: int(c) - 1}
}
