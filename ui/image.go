package ui

import (
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"os"

	"gioui.org/f32"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
)

// ImageFit defines how the image should be scaled to fit its container
type ImageFit int

const (
	// FitContain scales the image to fit within the bounds while preserving aspect ratio
	FitContain ImageFit = iota
	// FitCover scales the image to cover the bounds while preserving aspect ratio (may crop)
	FitCover
	// FitFill stretches the image to fill the bounds (may distort)
	FitFill
	// FitNone displays the image at its original size
	FitNone
)

// ImageOption configures the Image widget
type ImageOption func(*imageModel)

// ImageWidth sets the display width in dp
func ImageWidth(dp float32) ImageOption {
	return func(m *imageModel) { m.width = dp; m.hasWidth = true }
}

// ImageHeight sets the display height in dp
func ImageHeight(dp float32) ImageOption {
	return func(m *imageModel) { m.height = dp; m.hasHeight = true }
}

// Fit sets how the image should be scaled
func Fit(f ImageFit) ImageOption {
	return func(m *imageModel) { m.fit = f }
}

// ImageRadius sets the corner radius for rounded images
func ImageRadius(dp float32) ImageOption {
	return func(m *imageModel) { m.radius = dp }
}

// imageModel holds image configuration (internal)
type imageModel struct {
	src       image.Image
	width     float32
	height    float32
	hasWidth  bool
	hasHeight bool
	fit       ImageFit
	radius    float32
}

// Image creates an image widget from a local file path
// Usage: Image("path/to/image.png", ImageWidth(200), ImageHeight(150), Fit(FitCover), ImageRadius(8))
func Image(path string, opts ...ImageOption) Widget {
	m := &imageModel{
		fit: FitContain, // sensible default
	}
	for _, opt := range opts {
		opt(m)
	}

	// Load the image from file
	file, err := os.Open(path)
	if err != nil {
		// Return an empty widget if file can't be opened
		return func(gtx layout.Context, th *Theme) layout.Dimensions {
			return layout.Dimensions{}
		}
	}

	img, _, err := image.Decode(file)
	file.Close()
	if err != nil {
		// Return an empty widget if image can't be decoded
		return func(gtx layout.Context, th *Theme) layout.Dimensions {
			return layout.Dimensions{}
		}
	}
	m.src = img

	// Pre-create the image op (cached)
	imgOp := paint.NewImageOp(m.src)

	return func(gtx layout.Context, th *Theme) layout.Dimensions {
		if m.src == nil {
			return layout.Dimensions{}
		}

		imgSize := m.src.Bounds().Size()

		// Calculate display size
		var displayWidth, displayHeight int
		if m.hasWidth && m.hasHeight {
			displayWidth = gtx.Dp(unit.Dp(m.width))
			displayHeight = gtx.Dp(unit.Dp(m.height))
		} else if m.hasWidth {
			displayWidth = gtx.Dp(unit.Dp(m.width))
			// Maintain aspect ratio
			displayHeight = int(float32(displayWidth) * float32(imgSize.Y) / float32(imgSize.X))
		} else if m.hasHeight {
			displayHeight = gtx.Dp(unit.Dp(m.height))
			// Maintain aspect ratio
			displayWidth = int(float32(displayHeight) * float32(imgSize.X) / float32(imgSize.Y))
		} else {
			// Use image's natural size, constrained by available space
			displayWidth = imgSize.X
			displayHeight = imgSize.Y
			if displayWidth > gtx.Constraints.Max.X {
				displayWidth = gtx.Constraints.Max.X
				displayHeight = int(float32(displayWidth) * float32(imgSize.Y) / float32(imgSize.X))
			}
			if displayHeight > gtx.Constraints.Max.Y {
				displayHeight = gtx.Constraints.Max.Y
				displayWidth = int(float32(displayHeight) * float32(imgSize.X) / float32(imgSize.Y))
			}
		}

		// Apply fit mode
		var scaleX, scaleY float32
		switch m.fit {
		case FitContain:
			scaleW := float32(displayWidth) / float32(imgSize.X)
			scaleH := float32(displayHeight) / float32(imgSize.Y)
			if scaleW < scaleH {
				scaleX, scaleY = scaleW, scaleW
			} else {
				scaleX, scaleY = scaleH, scaleH
			}
		case FitCover:
			scaleW := float32(displayWidth) / float32(imgSize.X)
			scaleH := float32(displayHeight) / float32(imgSize.Y)
			if scaleW > scaleH {
				scaleX, scaleY = scaleW, scaleW
			} else {
				scaleX, scaleY = scaleH, scaleH
			}
		case FitFill:
			scaleX = float32(displayWidth) / float32(imgSize.X)
			scaleY = float32(displayHeight) / float32(imgSize.Y)
		case FitNone:
			scaleX, scaleY = 1, 1
		}

		// Calculate final size
		finalWidth := int(float32(imgSize.X) * scaleX)
		finalHeight := int(float32(imgSize.Y) * scaleY)

		// For FitContain, use the scaled size as display size
		if m.fit == FitContain && !m.hasWidth && !m.hasHeight {
			displayWidth = finalWidth
			displayHeight = finalHeight
		}

		// Clip to display bounds with optional rounded corners
		radius := gtx.Dp(unit.Dp(m.radius))
		rect := image.Rect(0, 0, displayWidth, displayHeight)
		var clipStack clip.Stack
		if m.radius > 0 {
			clipStack = clip.UniformRRect(rect, radius).Push(gtx.Ops)
		} else {
			clipStack = clip.Rect(rect).Push(gtx.Ops)
		}

		// Center the image if it's larger than display (for FitCover)
		offsetX := float32(displayWidth-finalWidth) / 2
		offsetY := float32(displayHeight-finalHeight) / 2

		// Apply transformation and draw
		transform := f32.Affine2D{}.Offset(f32.Pt(offsetX, offsetY)).Scale(f32.Point{}, f32.Pt(scaleX, scaleY))
		defer op.Affine(transform).Push(gtx.Ops).Pop()
		imgOp.Add(gtx.Ops)
		paint.PaintOp{}.Add(gtx.Ops)

		clipStack.Pop()

		return layout.Dimensions{
			Size: image.Pt(displayWidth, displayHeight),
		}
	}
}

// ImageFromImage creates an image widget from an existing image.Image
// Useful when you've already loaded the image or generated it programmatically
func ImageFromImage(img image.Image, opts ...ImageOption) Widget {
	m := &imageModel{
		src: img,
		fit: FitContain,
	}
	for _, opt := range opts {
		opt(m)
	}

	if m.src == nil {
		return func(gtx layout.Context, th *Theme) layout.Dimensions {
			return layout.Dimensions{}
		}
	}

	// Pre-create the image op (cached)
	imgOp := paint.NewImageOp(m.src)

	return func(gtx layout.Context, th *Theme) layout.Dimensions {
		imgSize := m.src.Bounds().Size()

		// Calculate display size
		var displayWidth, displayHeight int
		if m.hasWidth && m.hasHeight {
			displayWidth = gtx.Dp(unit.Dp(m.width))
			displayHeight = gtx.Dp(unit.Dp(m.height))
		} else if m.hasWidth {
			displayWidth = gtx.Dp(unit.Dp(m.width))
			displayHeight = int(float32(displayWidth) * float32(imgSize.Y) / float32(imgSize.X))
		} else if m.hasHeight {
			displayHeight = gtx.Dp(unit.Dp(m.height))
			displayWidth = int(float32(displayHeight) * float32(imgSize.X) / float32(imgSize.Y))
		} else {
			displayWidth = imgSize.X
			displayHeight = imgSize.Y
			if displayWidth > gtx.Constraints.Max.X {
				displayWidth = gtx.Constraints.Max.X
				displayHeight = int(float32(displayWidth) * float32(imgSize.Y) / float32(imgSize.X))
			}
			if displayHeight > gtx.Constraints.Max.Y {
				displayHeight = gtx.Constraints.Max.Y
				displayWidth = int(float32(displayHeight) * float32(imgSize.X) / float32(imgSize.Y))
			}
		}

		// Apply fit mode
		var scaleX, scaleY float32
		switch m.fit {
		case FitContain:
			scaleW := float32(displayWidth) / float32(imgSize.X)
			scaleH := float32(displayHeight) / float32(imgSize.Y)
			if scaleW < scaleH {
				scaleX, scaleY = scaleW, scaleW
			} else {
				scaleX, scaleY = scaleH, scaleH
			}
		case FitCover:
			scaleW := float32(displayWidth) / float32(imgSize.X)
			scaleH := float32(displayHeight) / float32(imgSize.Y)
			if scaleW > scaleH {
				scaleX, scaleY = scaleW, scaleW
			} else {
				scaleX, scaleY = scaleH, scaleH
			}
		case FitFill:
			scaleX = float32(displayWidth) / float32(imgSize.X)
			scaleY = float32(displayHeight) / float32(imgSize.Y)
		case FitNone:
			scaleX, scaleY = 1, 1
		}

		finalWidth := int(float32(imgSize.X) * scaleX)
		finalHeight := int(float32(imgSize.Y) * scaleY)

		if m.fit == FitContain && !m.hasWidth && !m.hasHeight {
			displayWidth = finalWidth
			displayHeight = finalHeight
		}

		radius := gtx.Dp(unit.Dp(m.radius))
		rect := image.Rect(0, 0, displayWidth, displayHeight)
		var clipStack clip.Stack
		if m.radius > 0 {
			clipStack = clip.UniformRRect(rect, radius).Push(gtx.Ops)
		} else {
			clipStack = clip.Rect(rect).Push(gtx.Ops)
		}

		offsetX := float32(displayWidth-finalWidth) / 2
		offsetY := float32(displayHeight-finalHeight) / 2

		transform := f32.Affine2D{}.Offset(f32.Pt(offsetX, offsetY)).Scale(f32.Point{}, f32.Pt(scaleX, scaleY))
		defer op.Affine(transform).Push(gtx.Ops).Pop()
		imgOp.Add(gtx.Ops)
		paint.PaintOp{}.Add(gtx.Ops)

		clipStack.Pop()

		return layout.Dimensions{
			Size: image.Pt(displayWidth, displayHeight),
		}
	}
}
