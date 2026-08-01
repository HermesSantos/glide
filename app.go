package glide

const (
	defaultW = 14
	defaultH = 5
)

// App is a terminal drag-and-drop canvas.
type App struct {
	elements []Element
	hint     string
}

// New creates an empty glide app.
func New() *App {
	return &App{
		hint: "Arraste com o mouse  ·  q para sair",
	}
}

// AddElement registers a draggable element and returns the app for chaining.
//
//	glide.New().
//	    AddElement(glide.Element{Shape: glide.ShapeBox, X: 10, Y: 5, Label: "a"}).
//	    AddElement(glide.Element{Shape: glide.ShapeBox, Center: true, Label: "b"}).
//	    Run()
func (a *App) AddElement(e Element) *App {
	a.elements = append(a.elements, e)
	return a
}

// SetHint overrides the bottom/top hint text.
func (a *App) SetHint(hint string) *App {
	a.hint = hint
	return a
}

// Run starts the interactive terminal UI. Blocks until the user quits.
func (a *App) Run() error {
	m := newModel(a.elements, a.hint)
	_, err := newProgram(m).Run()
	return err
}
