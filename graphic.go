package main

import (
	"image"
	"image/color"

	"golang.org/x/exp/shiny/materialdesign/colornames"
	"golang.org/x/exp/shiny/screen"
)

var (
	upColors = [5]color.RGBA{
		colornames.Green900,
		colornames.Red900,
		colornames.Yellow900,
		colornames.Blue900,
		colornames.BlueGrey900,
	}

	downColors = [5]color.RGBA{
		colornames.Green600,
		colornames.Red600,
		colornames.Yellow600,
		colornames.Blue600,
		colornames.BlueGrey600,
	}
)

// Picture is the image to draw the game.
type Picture struct {
	buffer  screen.Buffer
	texture screen.Texture
}

// NewPicture create a new picture on the specified screen.
func NewPicture(r image.Rectangle, s screen.Screen) (*Picture, error) {
	b, err := s.NewBuffer(r.Size())
	if err != nil {
		return nil, err
	}

	t, err := s.NewTexture(r.Size())
	if err != nil {
		return nil, err
	}

	return &Picture{buffer: b, texture: t}, nil
}

// Draw designs the game image.
func (p *Picture) Draw() {
	rgba := p.buffer.RGBA()
	size := p.buffer.Size()

	// draw buttons color
	for i, button := range rectangles(size) {
		sb := button.Bounds()

		for x := sb.Min.X; x < sb.Max.X; x++ {
			for y := sb.Min.Y; y < sb.Max.Y; y++ {
				rgba.SetRGBA(x, y, upColors[i])
			}
		}
	}

	// draw buttons separators
	mx := size.X / 2
	my := size.Y / 2

	for x := 0; x < size.X; x++ {
		for y := 0; y < size.Y; y++ {
			if x == mx || y == my {
				rgba.SetRGBA(x, y, colornames.White)
			} else if x == mx-1 || y == my-1 || x == mx+1 || y == my+1 {
				rgba.SetRGBA(x, y, colornames.Black)
			}
		}
	}

	// draw center
	rx := min(mx/5, my/5)
	cr := image.Rect(mx-rx, my-rx, mx+rx, my+rx)

	for x := cr.Min.X; x < cr.Max.X; x++ {
		for y := cr.Min.Y; y < cr.Max.Y; y++ {
			if x == cr.Min.X || y == cr.Min.Y || x == cr.Max.X-1 || y == cr.Max.Y-1 {
				// black border external
				if x != mx && y != my {
					rgba.SetRGBA(x, y, colornames.Black)
				}
			} else if x == cr.Min.X+1 || y == cr.Min.Y+1 || x == cr.Max.X-2 || y == cr.Max.Y-2 {
				// white border
				rgba.SetRGBA(x, y, colornames.White)
			} else if x == cr.Min.X+2 || y == cr.Min.Y+2 || x == cr.Max.X-3 || y == cr.Max.Y-3 {
				// black border internal
				rgba.SetRGBA(x, y, colornames.Black)
			} else {
				// center
				rgba.SetRGBA(x, y, upColors[4])
			}
		}
	}

	// upload buffer to texture
	p.texture.Upload(image.ZP, p.buffer, p.buffer.Bounds())
}

// DrawEvent designs the event e.
func (p *Picture) DrawEvent(e interface{}) {
	var index int
	var up bool

	switch e := e.(type) {
	case PlayerEvent:
		index, up = e.Index, e.Up
	case GameEvent:
		index, up = e.Index, e.Up
	case LostEvent:
		index, up = e.Index, e.Up
	default:
		return
	}

	if index == 4 {
		p.drawCenter(up)
	} else {
		p.drawIndex(index, up)
	}

	// upload buffer to texture
	p.texture.Upload(image.ZP, p.buffer, p.buffer.Bounds())
}

func (p *Picture) drawIndex(index int, up bool) {
	rgba := p.buffer.RGBA()
	size := p.buffer.Size()

	mx := size.X / 2
	my := size.Y / 2

	// center
	rx := min(mx/5, my/5)
	cr := image.Rect(mx-rx, my-rx, mx+rx, my+rx)

	// button
	button := rectangles(size)[index]
	sb := button.Bounds()

	// color
	bg := downColors[index]
	if up {
		bg = upColors[index]
	}

	// draw button and separators
	for x := sb.Min.X; x < sb.Max.X; x++ {
		for y := sb.Min.Y; y < sb.Max.Y; y++ {
			p := image.Pt(x, y)
			if !p.In(cr) {
				if x == mx || y == my {
					rgba.SetRGBA(x, y, colornames.White)
				} else if x == mx-1 || y == my-1 || x == mx+1 || y == my+1 {
					rgba.SetRGBA(x, y, colornames.Black)
				} else {
					rgba.SetRGBA(x, y, bg)
				}
			}
		}
	}
}

func (p *Picture) drawCenter(up bool) {
	rgba := p.buffer.RGBA()
	size := p.buffer.Size()

	mx := size.X / 2
	my := size.Y / 2

	// center
	rx := min(mx/5, my/5)
	cr := image.Rect(mx-rx, my-rx, mx+rx, my+rx)

	for x := cr.Min.X + 3; x < cr.Max.X-3; x++ {
		for y := cr.Min.Y + 3; y < cr.Max.Y-3; y++ {
			if up {
				rgba.SetRGBA(x, y, upColors[4])
			} else {
				rgba.SetRGBA(x, y, downColors[4])
			}
		}
	}
}

// Release releases the Picture's resources.
func (p *Picture) Release() {
	p.buffer.Release()
	p.texture.Release()
}

// Texture returns the Texture of the picture.
func (p *Picture) Texture() screen.Texture {
	return p.texture
}

// Bounds returns the domain for which At can return non-zero color.
func (p *Picture) Bounds() image.Rectangle {
	return p.buffer.RGBA().Bounds()
}

// At returns the color of the pixel at (x, y).
func (p *Picture) At(x, y int) color.Color {
	return p.buffer.RGBA().At(x, y)
}

func rectangles(p image.Point) [4]image.Rectangle {
	mx := p.X / 2
	my := p.Y / 2

	return [4]image.Rectangle{{
		// NW.
		Min: image.Point{
			X: 0,
			Y: 0,
		},
		Max: image.Point{
			X: mx,
			Y: my,
		},
	}, {
		// NE.
		Min: image.Point{
			X: mx,
			Y: 0,
		},
		Max: image.Point{
			X: p.X,
			Y: my,
		},
	}, {
		// SW.
		Min: image.Point{
			X: 0,
			Y: my,
		},
		Max: image.Point{
			X: mx,
			Y: p.Y,
		},
	}, {
		// SE.
		Min: image.Point{
			X: mx,
			Y: my,
		},
		Max: image.Point{
			X: p.X,
			Y: p.Y,
		},
	}}
}

func min(i, j int) int {
	if i < j {
		return i
	}
	return j
}
