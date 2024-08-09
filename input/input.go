package input

import (
	"errors"
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/help"      // Provides help view for key bindings
	"github.com/charmbracelet/bubbles/key"       // Manages key bindings
	"github.com/charmbracelet/bubbles/textinput" // Provides text input model
	tea "github.com/charmbracelet/bubbletea"     // Framework for building terminal applications
	"github.com/charmbracelet/lipgloss"          // Styles terminal UI components
	"github.com/nmeilick/go-ui"
)

// Model is the model handling user input.
type Model struct {
	textInput  textinput.Model // textInput is the text input model.
	help       help.Model      // help is the help model for displaying key bindings.
	keymap     keymap          // keymap is for managing key bindings.
	abort      bool            // abort indicates if the input operation was aborted.
	cancelable bool            // cancelable determines if selection can be canceled with escape key
	quitable   bool            // quitable determines if execution can be quit via ctrl+c

	canceled bool // canceled indicates whether the selection was canceled
	quit     bool // quit indicates whether the selection was quit
}

type keymap struct{}

// ShortHelp returns a list of key bindings for short help.
func (k keymap) ShortHelp() []key.Binding {
	return []key.Binding{
		key.NewBinding(key.WithKeys("tab"), key.WithHelp("tab", "complete")),
		key.NewBinding(key.WithKeys("ctrl+n"), key.WithHelp("ctrl+n", "next")),
		key.NewBinding(key.WithKeys("ctrl+p"), key.WithHelp("ctrl+p", "prev")),
		key.NewBinding(key.WithKeys("esc"), key.WithHelp("esc", "quit")),
	}
}

// FullHelp returns a list of key bindings for full help.
func (k keymap) FullHelp() [][]key.Binding {
	return [][]key.Binding{k.ShortHelp()}
}

// New creates and returns a new Model with default settings.
func New(prompt, value string, suggestions ...string) *Model {
	ti := textinput.New()
	ti.Prompt = prompt
	ti.SetValue(value)
	if len(suggestions) > 0 {
		ti.SetSuggestions(suggestions)
	}
	ti.Placeholder = ""
	ti.PromptStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("63"))
	ti.Cursor.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("63"))
	ti.Focus()
	ti.CharLimit = 100
	ti.Width = 40
	ti.ShowSuggestions = true
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

// WithPrompt sets the prompt for the text input model and returns a new Model with the updated prompt.
func (m *Model) WithPrompt(s string) *Model {
	newModel := *m
	newModel.textInput.Prompt = s
	return &newModel
}

// WithPlaceholder sets the placeholder for the text input model and returns a new Model with the updated placeholder.
func (m *Model) WithPlaceholder(s string) *Model {
	newModel := *m
	newModel.textInput.Placeholder = s
	return &newModel
}

// WithPromptStyle sets the style of the prompt for the text input model and returns a new Model with the updated prompt style.
func (m *Model) WithPromptStyle(style lipgloss.Style) *Model {
	newModel := *m
	newModel.textInput.PromptStyle = style
	return &newModel
}

// WithCursorStyle sets the style of the cursor for the text input model and returns a new Model with the updated cursor style.
func (m *Model) WithCursorStyle(style lipgloss.Style) *Model {
	newModel := *m
	newModel.textInput.Cursor.Style = style
	return &newModel
}

// WithCharLimit sets the maximum allowed number of input characters and returns a new Model with the updated
// character limit.
func (m *Model) WithCharLimit(n int) *Model {
	newModel := *m
	newModel.textInput.CharLimit = n
	return &newModel
}

// WithWidth sets the width of the text input model and returns a new Model with the updated width.
func (m *Model) WithWidth(n int) *Model {
	newModel := *m
	newModel.textInput.Width = n
	return &newModel
}

// WithSuggestion sets the autocomplete suggestions for the text input model and returns a new Model with the
// updated suggestions.
func (m *Model) WithSuggestion(suggestions []string) *Model {
	newModel := *m
	newModel.textInput.SetSuggestions(suggestions)
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

// Value returns the current input.
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

// Init initializes the Model, resets the abort flag, and returns a nil command.
func (m *Model) Init() tea.Cmd {
	m.abort = false
	return nil
}

// Update handles user input and updates the input state by processing key messages and updating the text input model
// accordingly.
func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
}

// View renders the input widget as a string, displaying the prompt, text input, and help view for key bindings.
func (m *Model) View() string {
	return fmt.Sprintf(
		"%s\n%s",
		m.textInput.View(),
		m.help.View(m.keymap),
	)
}

// Showcase demonstrates all features of the Model component by creating an input model with autocomplete
// suggestions and running an interactive example in the terminal.
func Showcase() {
	autocomplete := []string{"Apple", "Aardvark", "Banana", "Cherry", "Date", "Elderberry", "Fig", "Grape"}

	m := New("Default Style Input: ", "", autocomplete...)
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
		fmt.Printf("Final input: %s\n", m.textInput.Value())
	}
}
