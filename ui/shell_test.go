package ui

import (
	"image"
	"testing"

	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/unit"
)

func testContext(width, height int) layout.Context {
	return layout.Context{
		Ops:         new(op.Ops),
		Constraints: layout.Exact(image.Pt(width, height)),
		Metric:      unit.Metric{PxPerDp: 1, PxPerSp: 1},
	}
}

func TestScaffoldFABDoesNotReserveBodySpaceByDefault(t *testing.T) {
	bodyHeight := 0
	var body Widget = func(gtx layout.Context, _ *Theme) layout.Dimensions {
		bodyHeight = gtx.Constraints.Max.Y
		return layout.Dimensions{Size: gtx.Constraints.Min}
	}
	var fab Widget = func(_ layout.Context, _ *Theme) layout.Dimensions {
		return layout.Dimensions{Size: image.Pt(56, 56)}
	}

	Scaffold(Body(body), FAB(fab))(testContext(300, 300), &Light)
	if bodyHeight != 300 {
		t.Fatalf("body height = %d, want 300", bodyHeight)
	}
}

func TestScaffoldCanReserveFABSpace(t *testing.T) {
	bodyHeight := 0
	var body Widget = func(gtx layout.Context, _ *Theme) layout.Dimensions {
		bodyHeight = gtx.Constraints.Max.Y
		return layout.Dimensions{Size: gtx.Constraints.Min}
	}
	var fab Widget = func(_ layout.Context, _ *Theme) layout.Dimensions {
		return layout.Dimensions{Size: image.Pt(56, 56)}
	}

	Scaffold(Body(body), FAB(fab), BodyAvoidsFAB(true))(testContext(300, 300), &Light)
	if bodyHeight != 220 {
		t.Fatalf("body height = %d, want 220", bodyHeight)
	}
}

func TestAppBarMinimumHeight(t *testing.T) {
	gtx := testContext(400, 200)
	gtx.Constraints.Min.Y = 0
	if got := AppBar("Title")(gtx, &Light).Size.Y; got < 64 {
		t.Fatalf("AppBar height = %d, want at least 64", got)
	}
}

func TestDialogAndSnackbarStateInvalidate(t *testing.T) {
	invalidator := new(testInvalidator)
	dialog := NewDialogState().Bind(invalidator)
	dialog.Show()
	dialog.Show()
	dialog.Dismiss()
	if got := invalidator.count.Load(); got != 2 {
		t.Fatalf("dialog invalidations = %d, want 2", got)
	}

	snackbar := NewSnackbarState().Bind(invalidator)
	snackbar.Show("Saved")
	if !snackbar.Visible() {
		t.Fatal("Snackbar should be visible after Show")
	}
	snackbar.Dismiss()
	if snackbar.Visible() {
		t.Fatal("Snackbar should be hidden after Dismiss")
	}
	if got := invalidator.count.Load(); got != 4 {
		t.Fatalf("total invalidations = %d, want 4", got)
	}
}
