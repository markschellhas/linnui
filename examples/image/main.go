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
		w.Option(app.Title("LinnUI Image Example"))
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
					// Basic image (natural size, constrained by window)
					Text("Basic Image:", Style(H5)),
					Image("../../images/linnui.png"),

					// Image with fixed width (height auto-calculated to maintain aspect ratio)
					Text("Fixed Width (200dp):", Style(H5)),
					Image("../../images/linnui.png", ImageWidth(200)),

					// Image with fixed height
					Text("Fixed Height (100dp):", Style(H5)),
					Image("../../images/linnui.png", ImageHeight(100)),

					// Image with both dimensions and FitFill (may distort)
					Text("Fixed Size with FitFill:", Style(H5)),
					Image("../../images/linnui.png", ImageWidth(150), ImageHeight(80), Fit(FitFill)),

					// Image with rounded corners
					Text("Rounded Corners:", Style(H5)),
					Image("../../images/linnui.png", ImageWidth(150), ImageRadius(20)),

					// Image with FitCover (crops to fill)
					Text("FitCover (crops to fill):", Style(H5)),
					Image("../../images/linnui.png", ImageWidth(150), ImageHeight(80), Fit(FitCover)),
				}),
			)(gtx, &th)

			e.Frame(gtx.Ops)
		}
	}
}
