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
		tree := NewTree(w)
		defer tree.Close()
		if err := run(w, tree); err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}()
	app.Main()
}

func run(w *app.Window, tree *Tree) error {
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
					tree.Button("Click me!"),
					tree.Button("Filled Button (Default)", Variant(Filled)),
					tree.Button("Outlined Button", Variant(Outlined)),
					tree.Button("Text Button", Variant(TextButton)),
					tree.Button("Elevated Button", Variant(Elevated)),
					tree.Button("With OnClick", OnClick(func() {
						println("Button clicked!")
					})),
					tree.Button("With Custom ID", ButtonID("custom-id-button")),
					tree.Button("Duplicate Label", ButtonID("button-1")),
					tree.Button("Duplicate Label", ButtonID("button-2")),
				}),
			)(gtx, &th)

			e.Frame(gtx.Ops)
		}
	}
}
