package textarea

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/help" // Provides help view for key bindings
	"github.com/charmbracelet/bubbles/key"  // Manages key bindings
	"github.com/charmbracelet/bubbles/textarea"

	// Provides text textarea model
	tea "github.com/charmbracelet/bubbletea" // Framework for building terminal applications
	// Styles terminal UI components
	"github.com/nmeilick/go-ui"
)

var (
	defaultTextareaStyle = textarea.Style{}
)

type errMsg error

// Model is the model handling user textarea.
type Model struct {
	textInput  textarea.Model // textInput is the text textarea model.
	help       help.Model     // help is the help model for displaying key bindings.
	keymap     keymap         // keymap is for managing key bindings.
	cancelable bool           // cancelable determines if selection can be canceled with escape key
	quitable   bool           // quitable determines if execution can be quit via ctrl+c

	canceled bool // canceled indicates whether the selection was canceled
	quit     bool // quit indicates whether the selection was quit
}

type keymap struct{}

// ShortHelp returns a list of key bindings for short help.
func (k keymap) ShortHelp() []key.Binding {
	return []key.Binding{
		key.NewBinding(key.WithKeys("esc"), key.WithHelp("esc", "quit")),
	}
}

// FullHelp returns a list of key bindings for full help.
func (k keymap) FullHelp() [][]key.Binding {
	return [][]key.Binding{k.ShortHelp()}
}

// New creates and returns a new Model with default settings.
func New(prompt, value string) *Model {
	ti := textarea.New()
	ti.Prompt = prompt
	ti.SetValue(value)
	//ti.FocusedStyle = defaultTextareaStyle
	//ti.BlurredStyle = defaultTextareaStyle
	ti.Focus()
	ti.CharLimit = 100
	ti.MaxWidth = 40
	ti.MaxHeight = 10
	ti.ShowLineNumbers = true
	h := help.New()
	km := keymap{}

	return &Model{
		textInput:  ti,
		help:       h,
		keymap:     km,
		cancelable: true,
		quitable:   true,

		canceled: false,
		quit:     false,
	}
}

// WithPrompt sets the prompt for the text textarea model and returns a new Model with the updated prompt.
func (m *Model) WithPrompt(s string) *Model {
	newModel := *m
	newModel.textInput.Prompt = s
	return &newModel
}

// WithPlaceholder sets the placeholder for the text textarea model and returns a new Model with the updated placeholder.
func (m *Model) WithPlaceholder(s string) *Model {
	newModel := *m
	newModel.textInput.Placeholder = s
	return &newModel
}

// WithCharLimit sets the maximum allowed number of textarea characters and returns a new Model with the updated
// character limit.
func (m *Model) WithCharLimit(n int) *Model {
	newModel := *m
	newModel.textInput.CharLimit = n
	return &newModel
}

// WithMaxWidth sets the width of the text textarea model and returns a new Model with the updated width.
func (m *Model) WithMaxWidth(n int) *Model {
	newModel := *m
	newModel.textInput.MaxWidth = n
	return &newModel
}

// WithMaxHeight sets the height of the text textarea model and returns a new Model with the updated height.
func (m *Model) WithMaxHeight(n int) *Model {
	newModel := *m
	newModel.textInput.MaxHeight = n
	return &newModel
}

// WithCancel sets the cancelable flag and returns a new Model with the updated flag.
func (m *Model) WithCancel(cancelable bool) *Model {
	newModel := *m
	newModel.cancelable = cancelable
	return &newModel
}

// WithQuit sets the quitable flag and returns a new Model with the updated flag.
func (m *Model) WithQuit(quitable bool) *Model {
	newModel := *m
	newModel.quitable = quitable
	return &newModel
}

// Value returns the current textarea.
func (m *Model) Value() string {
	return m.textInput.Value()
}

// Canceled returns the canceled flag.
func (m *Model) Canceled() bool {
	return m.canceled
}

// Quit returns the quit flag.
func (m *Model) Quit() bool {
	return m.quit
}

// Init initializes the Model.
func (m *Model) Init() tea.Cmd {
	return textarea.Blink
}

// Update handles user textarea and updates the textarea state by processing key messages and updating the text textarea model
// accordingly.
func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			lines := strings.Split(m.textInput.Value(), "\n")
			for i := range lines {
				lines[i] = strings.TrimSpace(lines[i])
			}
			m.textInput.SetValue(strings.Join(lines, "\n"))
			if len(lines) > 0 && lines[len(lines)-1] == "" {
				m.canceled, m.quit = false, false
				return m, tea.Quit
			}
		case "esc":
			if m.textInput.Focused() {
				m.textInput.Blur()
			}
			m.canceled, m.quit = true, false
			return m, tea.Quit
		case "ctrl+c":
			m.canceled, m.quit = true, true
			return m, tea.Quit
		}
	// We handle errors just like any other message
	case errMsg:
		//m.err = msg
		return m, nil
	}

	m.textInput, cmd = m.textInput.Update(msg)
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)

	/*
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "enter":
				m.canceled, m.quit = false, false
				return m, tea.Quit
			case "esc":
				m.canceled, m.quit = true, false
				return m, tea.Quit
			case "ctrl+c":
				m.canceled, m.quit = true, true
				return m, tea.Quit
			}
		}

		var cmd tea.Cmd
		m.textInput, cmd = m.textInput.Update(msg)
		return m, cmd
	*/
}

// View renders the textarea widget as a string, displaying the prompt, text textarea, and help view for key bindings.
func (m *Model) View() string {
	return fmt.Sprintf(
		"%s\n%s",
		m.textInput.View(),
		m.help.View(m.keymap),
	)
}

// Showcase demonstrates all features of the Model component by creating an textarea model with autocomplete
// suggestions and running an interactive example in the terminal.
func Showcase() {
	m := New("", "")
	// Run interactive examples
	fmt.Println("=== Model Showcase ===")

	fmt.Println("\nDefault Style Input (Type to see suggestions, Enter to select):")
	err := ui.Run(m)
	switch {
	case errors.Is(err, ui.QuitError):
		fmt.Println("Quit")
		os.Exit(0)
	case errors.Is(err, ui.CanceledError):
		fmt.Println("Canceled")
	case err != nil:
		fmt.Printf("Error running program: %v", err)
	default:
		fmt.Printf("Final textarea: %s\n", m.textInput.Value())
	}
}
