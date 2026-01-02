package main

import (
	"fmt"
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
		w.Option(app.Title("LinnUI ScrollView Example"))
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

	// Create a list of items for the ListView example
	var listItems []Widget
	for i := 1; i <= 50; i++ {
		itemText := fmt.Sprintf("List Item %d", i)
		// Alternate background colors
		var bgColor color.NRGBA
		if i%2 == 0 {
			bgColor = color.NRGBA{R: 240, G: 240, B: 255, A: 255}
		} else {
			bgColor = color.NRGBA{R: 255, G: 255, B: 255, A: 255}
		}
		listItems = append(listItems, Container(
			Padding(InsetsAll(16), Text(itemText)),
			Background(bgColor),
		))
	}

	for {
		switch e := w.Event().(type) {
		case app.DestroyEvent:
			return e.Err
		case app.FrameEvent:
			gtx := app.NewContext(&ops, e)

			// Use Row to show two scroll examples side by side
			Row([]any{
				// Left side: ScrollView with a Column (single child scrolling)
				Container(
					Column([]any{
						Text("ScrollView Example", Style(H5)),
						ScrollView(
							Column([]any{
								Text("This is a ScrollView wrapping a Column."),
								Text("Scroll down to see more content..."),
								Container(Text("Item 1"), Background(color.NRGBA{R: 255, G: 200, B: 200, A: 255}), BorderRadius(8)),
								Container(Text("Item 2"), Background(color.NRGBA{R: 200, G: 255, B: 200, A: 255}), BorderRadius(8)),
								Container(Text("Item 3"), Background(color.NRGBA{R: 200, G: 200, B: 255, A: 255}), BorderRadius(8)),
								Container(Text("Item 4"), Background(color.NRGBA{R: 255, G: 255, B: 200, A: 255}), BorderRadius(8)),
								Container(Text("Item 5"), Background(color.NRGBA{R: 255, G: 200, B: 255, A: 255}), BorderRadius(8)),
								Container(Text("Item 6"), Background(color.NRGBA{R: 200, G: 255, B: 255, A: 255}), BorderRadius(8)),
								Container(Text("Item 7"), Background(color.NRGBA{R: 255, G: 220, B: 180, A: 255}), BorderRadius(8)),
								Container(Text("Item 8"), Background(color.NRGBA{R: 180, G: 220, B: 255, A: 255}), BorderRadius(8)),
								Text("End of ScrollView content"),
							}, Spacing(12)),
							ScrollID("scroll-1"),
						),
					}, Spacing(8)),
					Background(color.NRGBA{R: 250, G: 250, B: 250, A: 255}),
					BorderRadius(12),
				),

				// Right side: ListView (optimized for many items)
				Container(
					Column([]any{
						Text("ListView Example", Style(H5)),
						Text("50 items, efficiently rendered:"),
						ListView(listItems, ScrollID("listview-1")),
					}, Spacing(8)),
					Background(color.NRGBA{R: 250, G: 250, B: 250, A: 255}),
					BorderRadius(12),
				),
			}, RowSpacing(16))(gtx, &th)

			e.Frame(gtx.Ops)
		}
	}
}
