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

			Scaffold(
				AppBar(TitleBar("LinnUI Simple")),
				Body(
					Column([]any{
						Text("Welcome to LinnUI", Style(H3)),
						Row([]any{
							Container(
								Column([]any{
									Padding(InsetsAll(30),
										Text("Hi"),
									),
								}),
								Background(Gray100),
							),
						}),
						SizedBox(Height(20)),
						Margin(
							InsetsAll(50),
							Container(
								Padding(InsetsAll(50),
									Column([]any{
										Text("Card content", Style(H4)),
										SizedBox(Height(8)),
										Text("This is a decorated container"),
									}),
								),
								Background(White),
								BorderRadius(16),
								Border(BorderStyle{Width: 1, Color: Black}),
								Shadow(8),
							),
						),
						SizedBox(Height(20)),
						Padding(
							Insets{Left: 20, Right: 20},
							Row([]any{
								Button("Left", Variant(Filled)),
								Spacer(), // Fills space between buttons
								Button("Right", Variant(Outlined)),
							}),
						),
						Spacer(), // Pushes content below to bottom
						Center(
							Text("Footer at bottom"),
						),
					}, Spacing(0)),
				),
			)(gtx, &th)

			e.Frame(gtx.Ops)
		}
	}
}
