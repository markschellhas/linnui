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
		w.Option(app.Title("LinnUI Themes"))
		darkMode := NewState(false).Bind(w)
		tree := NewTree(w)
		defer tree.Close()

		if err := run(w, tree, darkMode); err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}()
	app.Main()
}

func run(w *app.Window, tree *Tree, darkMode *State[bool]) error {
	var ops op.Ops

	for {
		switch e := w.Event().(type) {
		case app.DestroyEvent:
			return e.Err
		case app.FrameEvent:
			gtx := app.NewContext(&ops, e)
			th := Light
			modeLabel := "Switch to dark mode"
			if darkMode.Get() {
				th = Dark
				modeLabel = "Switch to light mode"
			}

			Scaffold(
				TopBar(AppBar("Semantic themes")),
				Body(Padding(
					InsetsAll(24),
					Column([]any{
						Text("Theme-aware components", Style(H4)),
						Text(
							"Colors come from semantic palette roles, so every component updates together.",
							Style(BodyText),
						),
						Container(
							Padding(
								InsetsAll(20),
								Column([]any{
									Text("Primary container", Style(H5), TextColor(th.Palette.OnPrimaryContainer)),
									Text(
										"Text, surfaces, outlines, and controls retain contrast in either mode.",
										Style(Caption),
										TextColor(th.Palette.OnPrimaryContainer),
									),
								}, Spacing(8)),
							),
							PrimaryContainerBackground(),
							BorderRadius(16),
						),
						Container(
							Padding(
								InsetsAll(20),
								Column([]any{
									Text("Surface card", Style(H5)),
									tree.TextField(
										TextFieldID("theme_example_input"),
										Hint("Theme-aware text field"),
									),
									Row([]any{
										tree.Button("Filled", Variant(Filled)),
										tree.Button("Outlined", Variant(Outlined)),
										tree.Button("Elevated", Variant(Elevated)),
									}, RowSpacing(12)),
								}, Spacing(16)),
							),
							SurfaceBackground(),
							OutlineBorder(1),
							BorderRadius(16),
							Shadow(4),
						),
						tree.Button(
							modeLabel,
							ButtonID("toggle_theme"),
							OnClick(func() { darkMode.Update(func(dark bool) bool { return !dark }) }),
						),
					}, Spacing(20)),
				)),
			)(gtx, &th)

			e.Frame(gtx.Ops)
		}
	}
}
