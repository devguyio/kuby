package utils

import tea "github.com/charmbracelet/bubbletea"

type ChangeModelMsg struct {
	NewModel     *tea.Model
	CurrentModel *tea.Model
}

func CreateChangeModel(newModel tea.Model, currentModel tea.Model) func() tea.Msg {
	return func() tea.Msg {
		return ChangeModelMsg{NewModel: &newModel, CurrentModel: &currentModel}
	}
}

type BackToMainMenuMsg struct{}

func BackToMainMenu() tea.Msg {
	return BackToMainMenuMsg{}
}

type BackInHistoryMsg struct{}

func BackInHistory() tea.Msg {
	return BackInHistoryMsg{}
}

func WindowSizeCmd(Width int, Height int) tea.Cmd {
	return func() tea.Msg {
		return tea.WindowSizeMsg{
			Width:  Width,
			Height: Height,
		}
	}
}
