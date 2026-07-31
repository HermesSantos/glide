package main

import (
	"fmt"
	"os"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

const (
	boxW = 14
	boxH = 5
)

var (
	boxIdle = lipgloss.NewStyle().
		Width(boxW-2).
		Height(boxH-2).
		Align(lipgloss.Center, lipgloss.Center).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("39")).
		Foreground(lipgloss.Color("39"))

	boxDrag = lipgloss.NewStyle().
		Width(boxW-2).
		Height(boxH-2).
		Align(lipgloss.Center, lipgloss.Center).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("213")).
		Foreground(lipgloss.Color("213")).
		Bold(true)

	hintStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("245")).
			Padding(0, 1)
)

type model struct {
	width, height int
	x, y          int
	placed        bool
	dragging      bool
	offsetX       int
	offsetY       int
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
		if !m.placed {
			m.x = (m.width - boxW) / 2
			m.y = (m.height - boxH) / 2
			m.placed = true
		}
		m.clamp()

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
		if m.hit(mouse.X, mouse.Y) {
			m.dragging = true
			m.offsetX = mouse.X - m.x
			m.offsetY = mouse.Y - m.y
		}

	case tea.MouseMotionMsg:
		if !m.dragging {
			break
		}
		mouse := msg.Mouse()
		m.x = mouse.X - m.offsetX
		m.y = mouse.Y - m.offsetY
		m.clamp()

	case tea.MouseReleaseMsg:
		m.dragging = false
	}

	return m, nil
}

func (m *model) hit(mx, my int) bool {
	return mx >= m.x && mx < m.x+boxW && my >= m.y && my < m.y+boxH
}

func (m *model) clamp() {
	maxX := max(0, m.width-boxW)
	maxY := max(0, m.height-boxH)
	m.x = clamp(m.x, 0, maxX)
	m.y = clamp(m.y, 0, maxY)
}

func (m model) View() tea.View {
	label := "drag me"
	style := boxIdle
	if m.dragging {
		label = "dragging"
		style = boxDrag
	}

	hint := "Arraste o quadrado com o mouse  ·  q para sair"
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

	box := style.Render(label)

	root := lipgloss.NewLayer(bg)
	root.AddLayers(
		lipgloss.NewLayer(box).X(m.x).Y(m.y).Z(1),
	)

	v := tea.NewView(lipgloss.NewCompositor(root).Render())
	v.AltScreen = true
	v.MouseMode = tea.MouseModeAllMotion
	return v
}

func clamp(n, lo, hi int) int {
	return max(lo, min(n, hi))
}

func main() {
	if _, err := tea.NewProgram(model{}).Run(); err != nil {
		fmt.Fprintln(os.Stderr, "erro:", err)
		os.Exit(1)
	}
}
