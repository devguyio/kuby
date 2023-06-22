package views

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/muesli/termenv"
	"kuby/utils"
)

var (
	color     = termenv.EnvColorProfile().Color
	keyword   = termenv.Style{}.Foreground(color("204")).Background(color("235")).Styled
	helpStyle = termenv.Style{}.Foreground(color("241")).Styled
)

type ConnectErrorModel struct {
	Error error
}

func (m ConnectErrorModel) Init() tea.Cmd {
	return nil
}

func (m ConnectErrorModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q":
			return m, utils.BackInHistory
		}
	}

	return m, nil
}

func (m ConnectErrorModel) View() string {
	return utils.AppStyle.Render(fmt.Sprint(keyword(m.Error.Error()), "\n\n", "Is Docker/K8s running?", "\n\n\n", helpStyle("space: switch modes • q: go back\n")))
}
