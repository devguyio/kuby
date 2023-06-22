package views

import (
	"context"
	"fmt"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/samber/lo"
	lop "github.com/samber/lo/parallel"
	"gopkg.in/yaml.v3"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"kuby/k8s"
	"kuby/tables"
	"kuby/utils"
	"strconv"
	"strings"
)

type ListPodsModel struct {
	Help      help.Model
	Table     table.Model
	Altscreen bool
	Height    int
	Width     int
}

func (m ListPodsModel) Init() tea.Cmd { return nil }

func (m ListPodsModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, tables.DefaultTableKeyMap.Quit):
			return m, utils.BackToMainMenu
		case key.Matches(msg, tables.DefaultTableKeyMap.Help):
			m.Help.ShowAll = !m.Help.ShowAll
		case key.Matches(msg, tables.DefaultTableKeyMap.Choose):
			selectedRow := m.Table.SelectedRow()

			clientset := k8s.GetClientset()
			pod, err := clientset.CoreV1().Pods(selectedRow[2]).Get(context.TODO(), selectedRow[1], metav1.GetOptions{})
			if err != nil {
				panic(err)
			}

			pod.ObjectMeta.SetManagedFields(nil) // Really large and causes issues, not included in `kubectl get pods -o json` anyway

			metadataJson, _ := yaml.Marshal(pod.ObjectMeta)
			specJson, _ := yaml.Marshal(pod.Spec)
			statusJson, _ := yaml.Marshal(pod.Status)

			newModel := NewInspectJsonModel(m.Height, m.Width, string(metadataJson), string(specJson), string(statusJson))
			return m, utils.CreateChangeModel(&newModel, &m)
		}
	}
	m.Table, cmd = m.Table.Update(msg)
	return m, cmd
}

func (m ListPodsModel) View() string {
	var b strings.Builder

	tableView := utils.BaseStyle.Render(m.Table.View())
	helpView := list.DefaultStyles().HelpStyle.Render(m.Help.View(tables.DefaultTableKeyMap))

	b.WriteString(tableView)

	fmt.Fprint(&b, strings.Repeat("\n", lo.Max([]int{m.Height - lipgloss.Height(tableView) - lipgloss.Height(helpView) + 1, 0})))

	b.WriteString(helpView)

	return b.String()
}

func NewPodsTable(clientset *kubernetes.Clientset) (*table.Model, error) {
	// TODO: Maybe move the api call somewhere else and add pods as param e.g. separate concerns
	pods, err := clientset.CoreV1().Pods("").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	rows := lop.Map(pods.Items, func(item v1.Pod, index int) table.Row {
		return table.Row{strconv.Itoa(index), item.ObjectMeta.Name, item.ObjectMeta.Namespace, string(item.Status.Phase), item.Status.PodIP, item.Status.HostIP}
	})

	columns := []table.Column{
		{Title: "Index", Width: 5},
		{Title: "Name", Width: tables.LongestInColumn(&rows, 1)},
		{Title: "Namespace", Width: tables.LongestInColumn(&rows, 2)},
		{Title: "Phase", Width: tables.LongestInColumn(&rows, 3)},
		{Title: "Pod IP", Width: tables.LongestInColumn(&rows, 4)},
		{Title: "Host IP", Width: tables.LongestInColumn(&rows, 5)},
		//{Title: "Created", Width: LongestInColumn(&rows, 6)},
		//{Title: "Started", Width: LongestInColumn(&rows, 7)},
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(lo.Clamp(len(rows), 5, 30)),
	)

	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(false)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(false)
	t.SetStyles(s)

	return &t, err
}
