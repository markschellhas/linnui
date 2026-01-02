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
		w.Option(app.Title("LinnUI Buttons Example"))
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
					Button("Click me!"),
					Button("Filled Button (Default)", Variant(Filled)),
					Button("Outlined Button", Variant(Outlined)),
					Button("Text Button", Variant(TextButton)),
					Button("Elevated Button", Variant(Elevated)),
					Button("With OnClick", OnClick(func() {
						println("Button clicked!")
					})),
					Button("With Custom ID", ButtonID("custom-id-button")),
					Button("Duplicate Label", ButtonID("button-1")),
					Button("Duplicate Label", ButtonID("button-2")),
				}),
			)(gtx, &th)

			e.Frame(gtx.Ops)
		}
	}
}
