![LinnUI](images/linnui.png "LinnUI")

# LinnUI

**LinnUI** is a modern, pure-Go UI framework inspired by declaritive UI frameworks like Flutter, SwiftUI and Jetpack Compose, built on [Gio](https://gioui.org) for creating beautiful, reactive, cross-platform desktop, mobile, and web applications from a single codebase.

## Features

- **Flutter-like API**: Composable widgets (`Scaffold`, `Column`, `Row`, `Container`, etc.) with concise, declarative syntax.
- **Reactive state**: Svelte-inspired stores that automatically update the UI on changes.
- **Beautiful defaults**: Modern Material 3-inspired themes with light/dark mode and smooth animations.
- **Pure Go**: No webviews, no cgo — tiny static binaries with GPU acceleration.
- **True cross-platform**: Native desktop (Windows/macOS/Linux), mobile (via Gomobile), and WebAssembly.

## Quick Example

```go
package main

import (
	"gioui.org/app"
	ui "github.com/markschellhas/linnui/ui" 
)

func main() {
	// Reactive counter state
	count := ui.NewState(0)

	go func() {
		w := new(app.Window)
		w.Option(app.Title("LinnUI Counter"))
		count.Bind(w) // enable auto-redraw

		// UI loop...
		app.Main()
	}()
}
```

(See `examples/counter` for the full reactive counter demo.)

## Installation

```bash
go get github.com/markschellhas/linnui/ui
```

## Why LinnUI?

Go developers deserve modern, joyful UI tooling without compromises. LinnUI fills the gap between low-level Gio and bloated webview solutions.

## Status

Early development (v0.1-alpha) — API subject to change. Contributions welcome!

---

Copyright (c) 2025-2026 Mark Schellhas. All Rights Reserved.