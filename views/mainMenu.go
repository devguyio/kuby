package views

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/common-nighthawk/go-figure"
	"kuby/utils"
	"strings"
)

// Update loop for the first view where you're choosing a task.
func updateChoices(msg tea.Msg, m MainMenuModel) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			if i, ok := m.List.SelectedItem().(MainMenuItem); ok {
				var cmds []tea.Cmd
				m.Chosen = true
				m.Choice = i.GetModel(&m)
				cmds = append(cmds, m.Choice.Init())
				cmds = append(cmds, utils.WindowSizeCmd(m.Width, m.Height))
				return m, tea.Batch(cmds...)
				//fmt.Println(i.Title())
			}
		case "q":
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		h, v := utils.AppStyle.GetFrameSize()
		m.Width = msg.Width
		m.Height = msg.Height
		m.List.SetSize(msg.Width-h, msg.Height-v)
	}

	var cmd tea.Cmd
	m.List, cmd = m.List.Update(msg)

	return m, cmd
}

// Update loop for the second view after a choice has been made
func updateChosen(msg tea.Msg, m MainMenuModel) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
	//if msg.String() == "q" {
	//	m.Chosen = false
	//	return m, nil
	//}
	case utils.BackInHistoryMsg:
		if len(m.ViewHistory) == 0 {
			m.Chosen = false
			return m, nil
		}
		m.Choice = *m.ViewHistory[len(m.ViewHistory)-1]
		m.ViewHistory = m.ViewHistory[:len(m.ViewHistory)-1]
		return m, nil
	case utils.BackToMainMenuMsg:
		m.Chosen = false
		m.ViewHistory = []*tea.Model{}
		return m, nil
	case utils.ChangeModelMsg:
		m.ViewHistory = append(m.ViewHistory, msg.CurrentModel)
		m.Choice = *msg.NewModel
		return m, m.Choice.Init()
	}

	var cmd tea.Cmd
	m.Choice, cmd = m.Choice.Update(msg)

	return m, cmd
}

type MainMenuItem struct {
	TitleString, DescString string
	GetModel                func(currentModel *MainMenuModel) tea.Model
}

func (i MainMenuItem) Title() string       { return i.TitleString }
func (i MainMenuItem) Description() string { return i.DescString }
func (i MainMenuItem) FilterValue() string { return i.TitleString }

type MainMenuModel struct {
	List        list.Model
	Choice      tea.Model
	Chosen      bool
	ViewHistory []*tea.Model
	Width       int
	Height      int
}

func (m MainMenuModel) Init() tea.Cmd {
	return nil
}

func (m MainMenuModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
	}

	// Hand off the message and model to the appropriate update function for the
	// appropriate view based on the current state.
	if !m.Chosen {
		return updateChoices(msg, m)
	} else {
		return updateChosen(msg, m)
	}
}

func (m MainMenuModel) View() string {
	var s string

	if !m.Chosen {
		s = m.List.View()
	} else {
		s = m.Choice.View()
	}

	return utils.AppStyle.Render(s)
}

func NewMainMenuModel(items *[]list.Item) MainMenuModel {
	delegateKeys := utils.SliceKeyMap{KeyBindings: []key.Binding{
		key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "choose"),
		),
	}}

	// Setup list
	delegate := utils.NewDefaultListDelegate(&delegateKeys)
	mainMenuList := list.New(*items, delegate, 0, 0)
	logo := figure.NewFigure("Kuby", "", true)
	longLogo := strings.Join(logo.Slicify(), "\n")
	mainMenuList.Title = longLogo
	mainMenuList.Styles.Title = utils.TitleStyle

	return MainMenuModel{
		List:        mainMenuList,
		Choice:      nil,
		Chosen:      false,
		ViewHistory: []*tea.Model{},
	}
}
