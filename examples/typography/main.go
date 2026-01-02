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
		w.Option(app.Title("LinnUI Typography Example"))
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
					Text("H1 Typography", Style(H1)),
					Text("H2 Typography", Style(H2)),
					Text("H3 Typography", Style(H3)),
					Text("H4 Typography", Style(H4)),
					Text("H5 Typography", Style(H5)),
					Text("H6 Typography", Style(H6)),
					Text("Caption Typography", Style(Caption)),
					Text("Overline Typography", Style(Overline)),
					Text("12 point Typography", Size(12)),
					Text("10 point Typography", Size(10)),
				}),
			)(gtx, &th)

			e.Frame(gtx.Ops)
		}
	}
}
