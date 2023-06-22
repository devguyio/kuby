package views

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/knipferrc/teacup/code"
)

// CodeModel represents the properties of the UI.
type CodeModel struct {
	codeModels []code.Model
	Height     int
	Width      int
}

// NewCodeModel creates a new instance of the UI.
func NewCodeModel(Height int, Width int) CodeModel {
	var models []code.Model
	models = append(models, code.New(true, true, lipgloss.AdaptiveColor{Light: "#000000", Dark: "#ffffff"}))
	models = append(models, code.New(false, true, lipgloss.AdaptiveColor{Light: "#000000", Dark: "#ffffff"}))
	models = append(models, code.New(true, false, lipgloss.AdaptiveColor{Light: "#000000", Dark: "#ffffff"}))

	return CodeModel{
		codeModels: models,
		Height:     Height,
		Width:      Width,
	}
}

// Init intializes the UI.
func (m CodeModel) Init() tea.Cmd {
	return tea.Batch(m.codeModels[0].SetFileName("utils/delegates.go"), m.codeModels[1].SetFileName("views/listPods.go"), m.codeModels[2].SetFileName("utils/styles.go"))
}

// Update handles all UI interactions.
func (m CodeModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height

		return m, nil
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc", "q":
			cmds = append(cmds, tea.Quit)
		case "g":
			fmt.Println("hello")
		}
	}

	sizeCodeInputs(&m)

	for i := range m.codeModels {
		newModel, cmd := m.codeModels[i].Update(msg)
		m.codeModels[i] = newModel
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func sizeCodeInputs(m *CodeModel) {
	for i := range m.codeModels {
		m.codeModels[i].SetSize(m.Width/len(m.codeModels), m.Height)
	}
}

// View returns a string representation of the UI.
func (m CodeModel) View() string {
	//return m.code.View()
	var views []string
	for i := range m.codeModels {
		views = append(views, m.codeModels[i].View())
	}

	return lipgloss.JoinHorizontal(lipgloss.Top, views...)
}
