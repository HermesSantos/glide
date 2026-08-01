package glide

import (
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

var (
	boxIdle = lipgloss.NewStyle().
		Align(lipgloss.Center, lipgloss.Center).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("39")).
		Foreground(lipgloss.Color("39"))

	boxDrag = lipgloss.NewStyle().
		Align(lipgloss.Center, lipgloss.Center).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("213")).
		Foreground(lipgloss.Color("213")).
		Bold(true)

	hintStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("245")).
			Padding(0, 1)
)

type item struct {
	label         string
	shape         Shape
	x, y          int
	w, h          int
	center        bool
	placed        bool
}

type model struct {
	width, height int
	items         []item
	hint          string
	dragIdx       int // -1 = none
	offsetX       int
	offsetY       int
}

func newModel(elements []Element, hint string) model {
	items := make([]item, len(elements))
	for i, e := range elements {
		items[i] = item{
			label:  e.label(),
			shape:  e.Shape,
			x:      e.X,
			y:      e.Y,
			w:      e.width(),
			h:      e.height(),
			center: e.Center,
			placed: !e.Center,
		}
	}
	return model{
		items:   items,
		hint:    hint,
		dragIdx: -1,
	}
}

func newProgram(m model) *tea.Program {
	return tea.NewProgram(m)
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
		for i := range m.items {
			it := &m.items[i]
			if !it.placed {
				if it.center {
					it.x = (m.width - it.w) / 2
					it.y = (m.height - it.h) / 2
				}
				it.placed = true
			}
			m.clamp(i)
		}

	case tea.KeyPressMsg:
		switch msg.String() {
		case "q", "ctrl+c", "esc":
			return m, tea.Quit
		}

	case tea.MouseClickMsg:
		mouse := msg.Mouse()
		if mouse.Button != tea.MouseLeft {
			break
		}
		// Topmost hit wins (last in slice = higher z).
		for i := len(m.items) - 1; i >= 0; i-- {
			if m.hit(i, mouse.X, mouse.Y) {
				m.dragIdx = i
				m.offsetX = mouse.X - m.items[i].x
				m.offsetY = mouse.Y - m.items[i].y
				break
			}
		}

	case tea.MouseMotionMsg:
		if m.dragIdx < 0 {
			break
		}
		mouse := msg.Mouse()
		i := m.dragIdx
		m.items[i].x = mouse.X - m.offsetX
		m.items[i].y = mouse.Y - m.offsetY
		m.clamp(i)

	case tea.MouseReleaseMsg:
		m.dragIdx = -1
	}

	return m, nil
}

func (m model) hit(i, mx, my int) bool {
	it := m.items[i]
	return mx >= it.x && mx < it.x+it.w && my >= it.y && my < it.y+it.h
}

func (m *model) clamp(i int) {
	it := &m.items[i]
	maxX := max(0, m.width-it.w)
	maxY := max(0, m.height-it.h)
	it.x = clamp(it.x, 0, maxX)
	it.y = clamp(it.y, 0, maxY)
}

func (m model) View() tea.View {
	hint := m.hint
	bg := hintStyle.Render(hint)
	if m.width > 0 && m.height > 0 {
		bg = lipgloss.Place(
			m.width,
			m.height,
			lipgloss.Left,
			lipgloss.Top,
			hintStyle.Render(hint),
			lipgloss.WithWhitespaceStyle(lipgloss.NewStyle()),
		)
	}

	root := lipgloss.NewLayer(bg)
	layers := make([]*lipgloss.Layer, 0, len(m.items))
	for i, it := range m.items {
		label := it.label
		idle, drag := stylesFor(it.shape)
		style := idle.Width(it.w - 2).Height(it.h - 2)
		if i == m.dragIdx {
			label = "dragging"
			style = drag.Width(it.w - 2).Height(it.h - 2)
		}
		box := style.Render(label)
		layers = append(layers, lipgloss.NewLayer(box).X(it.x).Y(it.y).Z(i+1))
	}
	root.AddLayers(layers...)

	v := tea.NewView(lipgloss.NewCompositor(root).Render())
	v.AltScreen = true
	v.MouseMode = tea.MouseModeAllMotion
	return v
}

func stylesFor(shape Shape) (idle, drag lipgloss.Style) {
	switch shape {
	default: // ShapeBox
		return boxIdle, boxDrag
	}
}

func clamp(n, lo, hi int) int {
	return max(lo, min(n, hi))
}
