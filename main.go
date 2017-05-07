// Simon is an implementation of the memory game for children.
package main

import (
	"image"
	"log"
	"math/rand"
	"time"

	"golang.org/x/exp/shiny/driver"
	"golang.org/x/exp/shiny/imageutil"
	"golang.org/x/exp/shiny/materialdesign/colornames"
	"golang.org/x/exp/shiny/screen"
	"golang.org/x/mobile/event/key"
	"golang.org/x/mobile/event/lifecycle"
	"golang.org/x/mobile/event/mouse"
	"golang.org/x/mobile/event/paint"
	"golang.org/x/mobile/event/size"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

var (
	inset = 1

	sz   image.Rectangle
	game *Game
	pict *Picture
)

func main() {
	driver.Main(func(s screen.Screen) {
		options := &screen.NewWindowOptions{Width: 350, Height: 350}
		window, err := s.NewWindow(options)
		if err != nil {
			log.Fatal(err)
		}
		defer window.Release()

		for {
			e := window.NextEvent()

			switch e := e.(type) {
			case mouse.Event:
				if e.Button != mouse.ButtonLeft {
					continue
				}

				if game == nil || game.Player() {
					pt := image.Pt(int(e.X)-inset, int(e.Y)-inset)
					if pe, ok := pict.PlayerEvent(pt); ok {
						window.Send(pe)
					}
				}

			case PlayerEvent:
				pict.DrawEvent(e)
				window.Send(paint.Event{})

				if game == nil {
					if e.Center() {
						game = new(Game)
						window.Send(ErrGameTurn)
					}

				} else if e.Up {
					if err := game.Play(e.Color()); err != nil {
						window.Send(err)
					}
					if !game.Player() {
						window.Send(ErrGameTurn)
					}
				}

			case ErrBadColor:
				window.Send(NewLostEvent(e.Want))

			case LostEvent:
				time.Sleep(100 * time.Millisecond)

				pict.DrawEvent(e)
				window.Send(paint.Event{})

				if e.Up {
					e.Count++
				}
				if e.Count < 3 {
					e.Up = !e.Up
					window.Send(e)
				}

			case error:
				switch e {
				case ErrGameTurn:
					if c, err := game.Display(); err != nil {
						window.Send(err)
					} else {
						window.Send(NewGameEvent(c))
					}
				case ErrGameOver:
					game = nil
				}

			case GameEvent:
				time.Sleep(300 * time.Millisecond)

				pict.DrawEvent(e)
				window.Send(paint.Event{})

				if e.Down() {
					e.Up = true
					window.Send(e)
				} else {
					window.Send(ErrGameTurn)
				}

			case size.Event:
				if pict != nil {
					pict.Release()
				}

				sz = e.Bounds()
				var err error

				pict, err = NewPicture(sz.Inset(inset), s)
				if err != nil {
					log.Fatal(err)
				}
				pict.Draw()

			case paint.Event:
				for _, r := range imageutil.Border(sz, inset) {
					window.Fill(r, colornames.White, screen.Src)
				}

				rect := sz.Inset(inset)
				dp := image.Pt(rect.Min.X, rect.Min.Y)

				window.Copy(dp, pict.Texture(), pict.Bounds(), screen.Src, nil)
				window.Publish()

			case key.Event:
				if e.Code == key.CodeEscape {
					return
				}

			case lifecycle.Event:
				if e.To == lifecycle.StageDead {
					return
				}
			}
		}
	})
}
