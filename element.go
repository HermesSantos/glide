package glide

// Shape defines how an element is rendered.
type Shape int

const (
	// ShapeBox is a rounded bordered box (default).
	ShapeBox Shape = iota
)

// Element is a draggable item on the canvas.
type Element struct {
	// Label is the text shown inside the element.
	Label string

	// Shape controls the visual style. Defaults to ShapeBox.
	Shape Shape

	// X, Y are the top-left cell coordinates. If both are 0 and
	// Center is true (or neither X/Y were set via Centered), the
	// element is placed at the center on first layout.
	X, Y int

	// Width and Height in terminal cells. Zero uses defaults (14x5).
	Width, Height int

	// Center places the element in the middle of the terminal on first layout.
	Center bool
}

func (e Element) width() int {
	if e.Width > 0 {
		return e.Width
	}
	return defaultW
}

func (e Element) height() int {
	if e.Height > 0 {
		return e.Height
	}
	return defaultH
}

func (e Element) label() string {
	if e.Label != "" {
		return e.Label
	}
	return "drag me"
}
