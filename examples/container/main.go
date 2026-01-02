package main

import (
	"image/color"
	"log"
	"os"

	"gioui.org/app"
	"gioui.org/op"
	. "github.com/markschellhas/linnui/ui"
)

func main() {
	go func() {
		w := new(app.Window)
		w.Option(app.Title("LinnUI Container Example"))
		if err := run(w); err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}()
	app.Main()
}

func run(w *app.Window) error {
	var ops op.Ops
	th := Light

	for {
		switch e := w.Event().(type) {
		case app.DestroyEvent:
			return e.Err
		case app.FrameEvent:
			gtx := app.NewContext(&ops, e)

			Center(
				Column([]any{
					// Basic container with background
					Container(
						Text("Background Color"),
						Background(color.NRGBA{R: 200, G: 230, B: 255, A: 255}),
					),

					// Container with border radius
					Container(
						Text("Rounded Corners"),
						Background(color.NRGBA{R: 255, G: 220, B: 200, A: 255}),
						BorderRadius(16),
					),

					// Container with border
					Container(
						Text("With Border"),
						Border(BorderAll(2, color.NRGBA{R: 100, G: 100, B: 200, A: 255})),
						BorderRadius(8),
					),

					// Container with shadow/elevation
					Container(
						Text("With Shadow"),
						Background(color.NRGBA{R: 255, G: 255, B: 255, A: 255}),
						BorderRadius(12),
						Shadow(8),
					),

					// Container with all options combined
					Container(
						Text("All Options Combined"),
						Background(color.NRGBA{R: 240, G: 255, B: 240, A: 255}),
						BorderRadius(20),
						Border(BorderAll(3, color.NRGBA{R: 50, G: 150, B: 50, A: 255})),
						Shadow(12),
					),

					// Nested containers
					Container(
						Container(
							Text("Nested Container"),
							Background(color.NRGBA{R: 255, G: 255, B: 200, A: 255}),
							BorderRadius(8),
						),
						Background(color.NRGBA{R: 200, G: 200, B: 255, A: 255}),
						BorderRadius(12),
					),
				}),
			)(gtx, &th)

			e.Frame(gtx.Ops)
		}
	}
}
