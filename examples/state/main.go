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
		w.Option(app.Title("LinnUI State Examples"))

		// Create reactive state - binds to window for auto-invalidation
		count := NewState(0).Bind(w)
		inputText := NewState("").Bind(w)

		if err := run(w, count, inputText); err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}()
	app.Main()
}

func run(w *app.Window, count *State[int], inputText *State[string]) error {
	var ops op.Ops
	th := Light

	for {
		switch e := w.Event().(type) {
		case app.DestroyEvent:
			return e.Err
		case app.FrameEvent:
			gtx := app.NewContext(&ops, e)

			Scaffold(
				AppBar(TitleBar("LinnUI State Management")),
				Body(
					Column([]any{
						// Text input section
						Text("Type something:", Style(H5)),
						SizedBox(Height(8)),
						TextField(
							TextFieldID("name_input"),
							Hint("Enter your name..."),
							OnChange(func(s string) {
								inputText.Set(s)
							}),
						),
						SizedBox(Height(8)),
						Text("You typed: "+inputText.Get(), Style(BodyText)),

						SizedBox(Height(32)),

						// Counter section
						Center(
							Column([]any{
								Text("Count: "+strconv.Itoa(count.Get()), Style(H5)),
								SizedBox(Height(32)),
								Row([]any{
									Button("- Decrement",
										ButtonID("decrement"),
										OnClick(func() { count.Set(count.Get() - 1) }),
										Variant(Outlined),
									),
									SizedBox(Width(16)),
									Button("+ Increment",
										ButtonID("increment"),
										OnClick(func() { count.Set(count.Get() + 1) }),
										Variant(Filled),
									),
								}, RowSpacing(0)),
								SizedBox(Height(16)),
								Button("Reset",
									ButtonID("reset"),
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
