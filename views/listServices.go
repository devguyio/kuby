package views

import (
	"context"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"gopkg.in/yaml.v3"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"kuby/k8s"
	"kuby/utils"
)

type ListServicesModel struct {
	Table     table.Model
	Altscreen bool
	Height    int
	Width     int
}

func (m ListServicesModel) Init() tea.Cmd { return nil }

func (m ListServicesModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			if m.Table.Focused() {
				m.Table.Blur()
			} else {
				m.Table.Focus()
			}
		case "q":
			return m, utils.BackToMainMenu
		case "enter":
			selectedRow := m.Table.SelectedRow()

			clientset := k8s.GetClientset()
			service, err := clientset.CoreV1().Services(selectedRow[2]).Get(context.TODO(), selectedRow[1], metav1.GetOptions{})
			if err != nil {
				panic(err)
			}

			service.ObjectMeta.SetManagedFields(nil) // Really large and causes issues, not included in `kubectl get pods -o json` anyway

			metadataJson, _ := yaml.Marshal(service.ObjectMeta)
			specJson, _ := yaml.Marshal(service.Spec)
			statusJson, _ := yaml.Marshal(service.Status)

			newModel := NewInspectJsonModel(m.Height, m.Width, string(metadataJson), string(specJson), string(statusJson))
			return m, utils.CreateChangeModel(&newModel, &m)
		}
	}
	m.Table, cmd = m.Table.Update(msg)
	return m, cmd
}

func (m ListServicesModel) View() string {
	return utils.BaseStyle.Render(m.Table.View()) + "\n" + utils.Subtle("up/down: select") + utils.Dot + utils.Subtle("enter: choose") + utils.Dot + utils.Subtle("q: go back") + "\n"
}
