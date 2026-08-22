package main

import (
	"fmt"
	"log"
	"os"

	"gioui.org/app"
	"gioui.org/op"

	. "github.com/markschellhas/linnui/ui"
)

type galleryState struct {
	dark          *State[bool]
	name          *State[string]
	notes         *State[string]
	notifications *State[bool]
	analytics     *State[bool]
	plan          *State[string]
	volume        *State[float32]
	dialog        *DialogState
	snackbar      *SnackbarState
}

func main() {
	go func() {
		window := new(app.Window)
		window.Option(app.Title("LinnUI Component Gallery"))
		tree := NewTree(window)
		defer tree.Close()

		state := &galleryState{
			dark:          NewState(false).Bind(window),
			name:          NewState("Ada").Bind(window),
			notes:         NewState("").Bind(window),
			notifications: NewState(true).Bind(window),
			analytics:     NewState(false).Bind(window),
			plan:          NewState("pro").Bind(window),
			volume:        NewState(float32(65)).Bind(window),
			dialog:        NewDialogState().Bind(window),
			snackbar:      NewSnackbarState().Bind(window),
		}
		if err := run(window, tree, state); err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}()
	app.Main()
}

func run(window *app.Window, tree *Tree, state *galleryState) error {
	var ops op.Ops

	for {
		switch event := window.Event().(type) {
		case app.DestroyEvent:
			return event.Err
		case app.FrameEvent:
			gtx := app.NewContext(&ops, event)
			theme := Light
			themeLabel := "Dark mode"
			if state.dark.Get() {
				theme = Dark
				themeLabel = "Light mode"
			}

			dialog := Dialog(
				state.dialog,
				DialogTitle("Quality by default"),
				DialogContent(Text("LinnUI combines concise APIs with theming, accessibility, and predictable state.")),
				DialogActions(
					DialogAction{Label: "Cancel", Variant: TextButton},
					DialogAction{
						Label:   "Looks good",
						Variant: Filled,
						OnClick: func() { state.snackbar.Show("Thanks for trying LinnUI") },
					},
				),
			)
			snackbar := Snackbar(
				state.snackbar,
				SnackbarAction("Undo", func() {}),
			)

			Scaffold(
				TopBar(AppBar(
					"LinnUI",
					AppBarSubtitle("Component gallery"),
					AppBarActions(tree.Button(
						themeLabel,
						ButtonID("theme_toggle"),
						Variant(TextButton),
						OnClick(func() { state.dark.Update(func(value bool) bool { return !value }) }),
					)),
				)),
				Body(tree.ScrollView(
					Padding(
						InsetsAll(24),
						Column([]any{
							Text("Simple, polished Go UI", Style(H3)),
							Text("A compact component set with semantic themes and reactive state.", Style(BodyText)),
							section("Buttons", Row([]any{
								tree.Button("Filled", ButtonID("gallery_filled")),
								tree.Button("Outlined", ButtonID("gallery_outlined"), Variant(Outlined)),
								tree.Button("Text", ButtonID("gallery_text"), Variant(TextButton)),
								tree.Button("Disabled", ButtonID("gallery_disabled"), ButtonDisabled(true)),
							}, RowSpacing(12))),
							section("Text fields", Column([]any{
								tree.TextField(
									TextFieldID("gallery_name"),
									Hint("Your name"),
									BindText(state.name),
								),
								tree.TextField(
									TextFieldID("gallery_notes"),
									Hint("Notes"),
									MultiLine(),
									BindText(state.notes),
								),
								Text("Hello, "+state.name.Get(), Style(Caption)),
							}, Spacing(12))),
							section("Selection controls", Column([]any{
								tree.Checkbox("Product updates", state.notifications, CheckboxID("updates")),
								tree.Switch("Anonymous analytics", state.analytics, SwitchID("analytics")),
								Text("Plan", Style(H6)),
								tree.Radio(state.plan, []RadioChoice{
									{Value: "free", Label: "Free"},
									{Value: "pro", Label: "Pro"},
									{Value: "team", Label: "Team"},
								}, RadioID("plan")),
								tree.Slider(
									"Volume",
									state.volume,
									0,
									100,
									SliderID("volume"),
									SliderStep(5),
									SliderValueText(func(value float32) string {
										return fmt.Sprintf("%.0f%%", value)
									}),
								),
							}, Spacing(12))),
							Text("Cards", Style(H5)),
							Row([]any{
								Expanded(tree.Card(
									Column([]any{
										Text("Elevated", Style(H6)),
										Text("Depth for important groups.", Style(Caption)),
									}, Spacing(6)),
									CardID("elevated_card"),
								)),
								Expanded(tree.Card(
									Column([]any{
										Text("Filled", Style(H6)),
										Text("Subtle grouping without a border.", Style(Caption)),
									}, Spacing(6)),
									CardID("filled_card"),
									CardStyle(CardFilled),
								)),
								Expanded(tree.Card(
									Column([]any{
										Text("Outlined", Style(H6)),
										Text("Clear boundaries with low emphasis.", Style(Caption)),
									}, Spacing(6)),
									CardID("outlined_card"),
									CardStyle(CardOutlined),
								)),
							}, RowSpacing(16)),
							section("Feedback", Row([]any{
								tree.Button(
									"Open dialog",
									ButtonID("open_dialog"),
									OnClick(state.dialog.Show),
								),
								tree.Button(
									"Show snackbar",
									ButtonID("show_snackbar"),
									Variant(Outlined),
									OnClick(func() { state.snackbar.Show("Settings saved") }),
								),
							}, RowSpacing(12))),
						}, Spacing(20)),
					),
					ScrollID("gallery_scroll"),
				)),
				FAB(tree.Button(
					"+ Create",
					ButtonID("gallery_fab"),
					Variant(Elevated),
					OnClick(func() { state.snackbar.Show("Created a new item") }),
				)),
				ScaffoldOverlays(snackbar, dialog),
			)(gtx, &theme)

			event.Frame(gtx.Ops)
		}
	}
}

func section(title string, content Widget) Widget {
	return Card(
		Column([]any{
			Text(title, Style(H5)),
			content,
		}, Spacing(16)),
		CardStyle(CardOutlined),
	)
}
