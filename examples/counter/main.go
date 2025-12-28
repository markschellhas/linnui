package main

import (
	"log"
	"os"
	"strconv"

	"gioui.org/app"
	"gioui.org/op"

	. "github.com/markschellhas/linnui/ui"
)

func main() {
	go func() {
		w := new(app.Window)
		w.Option(app.Title("LinnUI Reactive Counter"))

		// Create reactive state - binds to window for auto-invalidation
		count := NewState(0).Bind(w)

		if err := run(w, count); err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}()
	app.Main()
}

func run(w *app.Window, count *State[int]) error {
	var ops op.Ops
	th := Light

	for {
		switch e := w.Event().(type) {
		case app.DestroyEvent:
			return e.Err
		case app.FrameEvent:
			gtx := app.NewContext(&ops, e)

			Scaffold(
				AppBar(TitleBar("Reactive Counter")),
				Body(
					Column([]any{
						Spacer(),
						Center(
							Column([]any{
								Text("Count: "+strconv.Itoa(count.Get()), Style(H1)),
								SizedBox(Height(32)),
								Row([]any{
									Button("- Decrement",
										OnClick(func() { count.Set(count.Get() - 1) }),
										Variant(Outlined),
									),
									SizedBox(Width(16)),
									Button("+ Increment",
										OnClick(func() { count.Set(count.Get() + 1) }),
										Variant(Filled),
									),
								}, RowSpacing(0)),
								SizedBox(Height(16)),
								Button("Reset",
									OnClick(func() { count.Set(0) }),
									Variant(TextButton),
								),
							}, Spacing(0)),
						),
						Spacer(),
					}, Spacing(0)),
				),
			)(gtx, &th)

			e.Frame(gtx.Ops)
		}
	}
}
