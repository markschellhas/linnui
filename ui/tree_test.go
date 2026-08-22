package ui

import "testing"

func TestTreeIsolatesWidgetState(t *testing.T) {
	first := NewTree(nil)
	second := NewTree(nil)
	defer first.Close()
	defer second.Close()

	firstButton := first.clickable(buttonWidget, "save")
	if got := first.clickable(buttonWidget, "save"); got != firstButton {
		t.Fatal("same tree and key did not reuse state")
	}
	if got := second.clickable(buttonWidget, "save"); got == firstButton {
		t.Fatal("different trees shared button state")
	}
	if got := first.Scope("dialog").clickable(buttonWidget, "save"); got == firstButton {
		t.Fatal("child scope shared parent button state")
	}
}

func TestTreeSeparatesWidgetKinds(t *testing.T) {
	tree := NewTree(nil)
	defer tree.Close()

	button := tree.clickable(buttonWidget, "shared")
	card := tree.clickable(cardWidget, "shared")
	if button == card {
		t.Fatal("different widget kinds shared interaction state")
	}
}

func TestTreeResetAndDelete(t *testing.T) {
	tree := NewTree(nil)
	defer tree.Close()
	child := tree.Scope("settings")

	button := child.clickable(buttonWidget, "save")
	child.Delete("save")
	if got := child.clickable(buttonWidget, "save"); got == button {
		t.Fatal("Delete did not release widget state")
	}

	field := child.textField("name")
	child.Reset()
	if got := child.textField("name"); got == field {
		t.Fatal("Reset did not release scoped widget state")
	}
}

func TestTreeTextFieldAccessors(t *testing.T) {
	tree := NewTree(nil)
	defer tree.Close()
	tree.textField("name")

	if ok := tree.SetTextFieldValue("name", "Lin"); !ok {
		t.Fatal("SetTextFieldValue did not find field")
	}
	if got, ok := tree.TextFieldValue("name"); !ok || got != "Lin" {
		t.Fatalf("TextFieldValue = %q, %v; want Lin, true", got, ok)
	}
	if tree.SetTextFieldValue("missing", "value") {
		t.Fatal("SetTextFieldValue unexpectedly found missing field")
	}
}

func TestTextFieldSynchronizesBoundState(t *testing.T) {
	tree := NewTree(nil)
	defer tree.Close()
	value := NewState("initial")
	field := tree.TextField(TextFieldID("name"), BindText(value))

	field(testContext(300, 56), &Light)
	if got, ok := tree.TextFieldValue("name"); !ok || got != "initial" {
		t.Fatalf("initial field value = %q, %v; want initial, true", got, ok)
	}

	tree.SetTextFieldValue("name", "updated")
	field(testContext(300, 56), &Light)
	if got := value.Get(); got != "updated" {
		t.Fatalf("bound value = %q, want updated", got)
	}
}

func TestClosedTreePanicsOnUse(t *testing.T) {
	tree := NewTree(nil)
	tree.Close()
	defer func() {
		if recover() == nil {
			t.Fatal("using a closed Tree did not panic")
		}
	}()
	tree.Button("Save")
}
