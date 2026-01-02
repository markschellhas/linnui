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
		w.Option(app.Title("LinnUI Center Example"))
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

			Center(Text("Oh, hi Mark"))(gtx, &th)
			e.Frame(gtx.Ops)
		}
	}
}
