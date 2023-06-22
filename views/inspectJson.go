package views

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"kuby/utils"
)

const (
	initialInputs = 3
	maxInputs     = 3
	minInputs     = 3
	helpHeight    = 5
)

var (
	cursorStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("212"))

	cursorLineStyle = lipgloss.NewStyle().
			Background(lipgloss.Color("57")).
			Foreground(lipgloss.Color("230"))

	placeholderStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("238"))

	endOfBufferStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("235"))

	focusedPlaceholderStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("99"))

	focusedBorderStyle = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(lipgloss.Color("238"))

	blurredBorderStyle = lipgloss.NewStyle().
				Border(lipgloss.HiddenBorder())
)

type keymap = struct {
	next, prev, add, remove, moveToBegin, quit key.Binding
}

func newTextarea(textValue string) textarea.Model {
	t := textarea.New()
	t.Prompt = ""
	t.Placeholder = ""
	t.ShowLineNumbers = true
	t.Cursor.Style = cursorStyle
	t.FocusedStyle.Placeholder = focusedPlaceholderStyle
	t.BlurredStyle.Placeholder = placeholderStyle
	t.FocusedStyle.CursorLine = cursorLineStyle
	t.FocusedStyle.Base = focusedBorderStyle
	t.BlurredStyle.Base = blurredBorderStyle
	t.FocusedStyle.EndOfBuffer = endOfBufferStyle
	t.BlurredStyle.EndOfBuffer = endOfBufferStyle
	t.KeyMap.DeleteWordBackward.SetEnabled(false)
	t.KeyMap.LineNext = key.NewBinding(key.WithKeys("down"))
	t.KeyMap.LinePrevious = key.NewBinding(key.WithKeys("up"))
	t.CharLimit = 0
	t.SetValue(textValue)
	t.Blur()
	return t
}

type InspectJsonModel struct {
	width  int
	height int
	keymap keymap
	help   help.Model
	inputs []textarea.Model
	focus  int
}

func NewInspectJsonModel(newHeight int, newWidth int, metadata string, spec string, status string) InspectJsonModel {
	m := InspectJsonModel{
		//inputs: make([]textarea.Model, initialInputs),
		inputs: make([]textarea.Model, 3),
		help:   help.New(),
		keymap: keymap{
			next: key.NewBinding(
				key.WithKeys("tab"),
				key.WithHelp("tab", "next"),
			),
			prev: key.NewBinding(
				key.WithKeys("shift+tab"),
				key.WithHelp("shift+tab", "prev"),
			),
			//add: key.NewBinding(
			//	key.WithKeys("ctrl+n"),
			//	key.WithHelp("ctrl+n", "add an editor"),
			//),
			moveToBegin: key.NewBinding(
				key.WithKeys("ctrl+home"),
				key.WithHelp("ctrl+home", "move cursor to top"),
			),
			remove: key.NewBinding(
				key.WithKeys("ctrl+w"),
				key.WithHelp("ctrl+w", "remove an editor"),
			),
			quit: key.NewBinding(
				key.WithKeys("esc", "q"),
				key.WithHelp("esc, q", "go back"),
			),
		},
		height: newHeight,
		width:  newWidth,
	}
	//for i := 0; i < initialInputs; i++ {
	//	m.inputs[i] = newTextarea()
	//}

	//metadataJson, _ := json.MarshalIndent(pod.ObjectMeta, "", "    ")
	//specJson, _ := json.MarshalIndent(pod.Spec, "", "    ")
	//statusJson, _ := json.MarshalIndent(pod.Status, "", "    ")
	//
	//m.inputs[0] = newTextarea(string(metadataJson)) // metadata
	//m.inputs[1] = newTextarea(string(specJson))     // spec
	//m.inputs[2] = newTextarea(string(statusJson))   // status

	//metadataJson, _ := yaml.Marshal(pod.ObjectMeta)
	//specJson, _ := yaml.Marshal(pod.Spec)
	//statusJson, _ := yaml.Marshal(pod.Status)
	//
	//m.inputs[0] = newTextarea(string(metadataJson)) // metadata
	//m.inputs[1] = newTextarea(string(specJson))     // spec
	//m.inputs[2] = newTextarea(string(statusJson))   // status

	m.inputs[0] = newTextarea(metadata) // metadata
	m.inputs[1] = newTextarea(spec)     // spec
	m.inputs[2] = newTextarea(status)   // status

	m.inputs[m.focus].Focus()
	updateKeybindings(&m)
	return m
}

func (m InspectJsonModel) Init() tea.Cmd {
	return textarea.Blink
}

func (m InspectJsonModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keymap.quit):
			for i := range m.inputs {
				m.inputs[i].Blur()
			}
			return m, utils.BackInHistory
		case key.Matches(msg, m.keymap.next):
			m.inputs[m.focus].Blur()
			m.focus++
			if m.focus > len(m.inputs)-1 {
				m.focus = 0
			}
			cmd := m.inputs[m.focus].Focus()
			cmds = append(cmds, cmd)
		case key.Matches(msg, m.keymap.prev):
			m.inputs[m.focus].Blur()
			m.focus--
			if m.focus < 0 {
				m.focus = len(m.inputs) - 1
			}
			cmd := m.inputs[m.focus].Focus()
			cmds = append(cmds, cmd)
		//case key.Matches(msg, m.keymap.add):
		//	m.inputs = append(m.inputs, newTextarea())
		case key.Matches(msg, m.keymap.remove):
			m.inputs = m.inputs[:len(m.inputs)-1]
			if m.focus > len(m.inputs)-1 {
				m.focus = len(m.inputs) - 1
			}
		}
	case tea.WindowSizeMsg:
		m.height = msg.Height
		m.width = msg.Width
	}

	updateKeybindings(&m)
	sizeInputs(&m)

	// Update all textareas
	for i := range m.inputs {
		newModel, cmd := m.inputs[i].Update(msg)
		m.inputs[i] = newModel
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func sizeInputs(m *InspectJsonModel) {
	for i := range m.inputs {
		m.inputs[i].SetWidth(m.width / len(m.inputs))
		m.inputs[i].SetHeight(m.height - helpHeight)
	}
}

func updateKeybindings(m *InspectJsonModel) {
	m.keymap.add.SetEnabled(len(m.inputs) < maxInputs)
	m.keymap.remove.SetEnabled(len(m.inputs) > minInputs)
}

func (m InspectJsonModel) View() string {
	shortHelp := m.help.ShortHelpView([]key.Binding{
		m.keymap.next,
		m.keymap.prev,
		//m.keymap.add,
		m.keymap.moveToBegin,
		m.keymap.remove,
		m.keymap.quit,
	})

	var views []string
	for i := range m.inputs {
		views = append(views, m.inputs[i].View())
	}

	return lipgloss.JoinHorizontal(lipgloss.Top, views...) + "\n\n" + shortHelp
}
