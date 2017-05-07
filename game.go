package main

import "math/rand"

// Color is a game color.
type Color int

// Game handles the state of the party between the player and the game.
// The game creates a series of colors and waits for the player to repeat the series.
// You get the colors via Display and repeat them via Play.
type Game struct {
	colors []Color
	score  int
	over   bool
	player bool
}

// Score returns the game score.
func (g *Game) Score() int {
	return g.score
}

// Over returns true when the game is over.
func (g *Game) Over() bool {
	return g.over
}

// Player returns true when the player can play.
// This assumes that the game is not over.
func (g *Game) Player() bool {
	return !g.over && g.player
}

// Display is the way to ask the game a random color to remember.
// The number of time you have to invoke Display depends on the score.
// You should call this method until the returned error is ErrPlayerTurn.
func (g *Game) Display() (Color, error) {
	if g.over {
		return 0, ErrGameOver
	}
	if g.player {
		return 0, ErrPlayerTurn
	}

	c := Color(rand.Intn(4) + 1)
	g.colors = append(g.colors, c)

	if len(g.colors) == g.score+3 {
		g.player = true
	}

	return c, nil
}

// Play is the way to play the game for the player.
// You should play all the colors given by Display in the same order.
// ErrBadColor is returned when the played color is not the one expected.
// Once ErrBadColor is returned then the game is over.
func (g *Game) Play(c Color) error {
	if g.over {
		return ErrGameOver
	}
	if !g.player {
		return ErrGameTurn
	}

	if c != g.colors[0] {
		g.over = true
		return ErrBadColor{Got: c, Want: g.colors[0]}
	}

	g.colors = g.colors[1:]
	if len(g.colors) == 0 {
		g.player = false
		g.score++
	}

	return nil
}
