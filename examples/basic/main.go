package main

import (
	"log"
	"os"

	"gioui.org/app"
	"gioui.org/op"

	. "github.com/markschellhas/linnui/ui"
)

func main() {
	go func() {
		w := new(app.Window)
		w.Option(app.Title("LinnUI Simple Example"))
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

			// The entire UI in one clean Scaffold - Flutter-like!
			Scaffold(
				AppBar(TitleBar("LinnUI Simple")),
				Body(
					Column([]Widget{
						Text("Welcome to LinnUI", Style(H4)),
						Text("A Flutter-like experience for Go"),
						Row([]Widget{
							Button("Hello LinnUI", Variant(Filled)),
							Column([]Widget{
								Text("Some more:"),
								Button("Hello Mark", Variant(Filled)),
							}),
						}),
					}, Spacing(30)),
				),
			)(gtx, &th)

			e.Frame(gtx.Ops)
		}
	}
}
